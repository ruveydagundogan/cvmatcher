package tunnel

import (
	"encoding/base64"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type TunnelResponse struct {
	Status  int
	Headers map[string]string
	Body    []byte
}

type Server struct {
	upgrader   websocket.Upgrader
	conn       *websocket.Conn
	connMu     sync.RWMutex
	pending    map[string]chan<- *TunnelResponse
	pendingMu  sync.Mutex
	connected  bool
}

func NewServer() *Server {
	return &Server{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		pending: make(map[string]chan<- *TunnelResponse),
	}
}

func (s *Server) IsConnected() bool {
	s.connMu.RLock()
	defer s.connMu.RUnlock()
	return s.connected
}

func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("tunnel: websocket upgrade failed", "error", err)
		return
	}

	s.connMu.Lock()
	if s.conn != nil {
		s.conn.Close()
	}
	s.conn = conn
	s.connected = true
	s.connMu.Unlock()

	slog.Info("tunnel: agent connected")

	defer func() {
		s.connMu.Lock()
		if s.conn == conn {
			s.conn = nil
			s.connected = false
		}
		s.connMu.Unlock()

		// Fail all in-flight requests instead of letting them hang until timeout.
		s.pendingMu.Lock()
		for id, ch := range s.pending {
			delete(s.pending, id)
			if ch != nil {
				ch <- &TunnelResponse{Status: http.StatusBadGateway, Body: []byte("tunnel: agent disconnected")}
			}
		}
		s.pendingMu.Unlock()

		conn.Close()
		slog.Warn("tunnel: agent disconnected")
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		resp, err := DecodeResponse(msg)
		if err != nil {
			slog.Warn("tunnel: invalid response message", "error", err)
			continue
		}

		s.pendingMu.Lock()
		ch, ok := s.pending[resp.ID]
		if ok {
			delete(s.pending, resp.ID)
		}
		s.pendingMu.Unlock()

		if ok && ch != nil {
			body, _ := base64.StdEncoding.DecodeString(resp.Body)
			ch <- &TunnelResponse{
				Status:  resp.Status,
				Headers: resp.Headers,
				Body:    body,
			}
		}
	}
}

func (s *Server) ProxyHandler(w http.ResponseWriter, r *http.Request) {
	s.connMu.RLock()
	conn := s.conn
	s.connMu.RUnlock()

	if conn == nil {
		http.Error(w, "tunnel: no agent connected", http.StatusServiceUnavailable)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "tunnel: read body failed", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	reqID := uuid.New().String()
	headers := make(map[string]string)
	for k := range r.Header {
		headers[k] = r.Header.Get(k)
	}

	targetPath := r.URL.Path
	if len(targetPath) >= 8 && targetPath[:8] == "/tunnel/" {
		targetPath = targetPath[7:]
	}

	req := &Request{
		ID:      reqID,
		Type:    "request",
		Method:  r.Method,
		Path:    targetPath,
		Headers: headers,
		Body:    base64.StdEncoding.EncodeToString(body),
	}

	data, err := req.Encode()
	if err != nil {
		http.Error(w, "tunnel: encode failed", http.StatusInternalServerError)
		return
	}

	respCh := make(chan *TunnelResponse, 1)

	s.pendingMu.Lock()
	s.pending[reqID] = respCh
	s.pendingMu.Unlock()

	s.connMu.RLock()
	if s.conn != conn {
		s.connMu.RUnlock()
		http.Error(w, "tunnel: agent reconnected", http.StatusServiceUnavailable)
		return
	}
	err = conn.WriteMessage(websocket.TextMessage, data)
	s.connMu.RUnlock()

	if err != nil {
		http.Error(w, "tunnel: write failed", http.StatusInternalServerError)
		return
	}

	select {
	case tunnelResp := <-respCh:
		s.pendingMu.Lock()
		delete(s.pending, reqID)
		s.pendingMu.Unlock()
		for k, v := range tunnelResp.Headers {
			w.Header().Set(k, v)
		}
		w.WriteHeader(tunnelResp.Status)
		w.Write(tunnelResp.Body)
	case <-time.After(180 * time.Second):
		s.pendingMu.Lock()
		delete(s.pending, reqID)
		s.pendingMu.Unlock()
		http.Error(w, "tunnel: request timed out", http.StatusGatewayTimeout)
	}
}
