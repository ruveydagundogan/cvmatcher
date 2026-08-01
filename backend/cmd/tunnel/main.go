package main

import (
	"encoding/base64"
	"flag"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/infrastructure/tunnel"
)

func main() {
	serverURL := flag.String("server", "ws://localhost:8080/tunnel/ws", "tunnel server WebSocket URL")
	ollamaURL := flag.String("ollama", "http://localhost:11434", "local Ollama URL")
	flag.Parse()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	log.Info("starting tunnel client", "server", *serverURL, "ollama", *ollamaURL)

	httpClient := &http.Client{Timeout: 120 * time.Second}

	for {
		c, _, err := websocket.DefaultDialer.Dial(*serverURL, nil)
		if err != nil {
			log.Warn("websocket connection failed, retrying in 3s", "error", err)
			time.Sleep(3 * time.Second)
			continue
		}

		log.Info("connected to tunnel server")
		go pingLoop(c, log)
		readLoop(c, ollamaURL, httpClient, log)
		log.Warn("disconnected, reconnecting in 3s")
		time.Sleep(3 * time.Second)
	}
}

func pingLoop(conn *websocket.Conn, log *slog.Logger) {
	// Keep the websocket alive so proxies (e.g. Render) do not drop it while idle.
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(5*time.Second)); err != nil {
				return
			}
		}
	}
}

func readLoop(conn *websocket.Conn, ollamaURL *string, httpClient *http.Client, log *slog.Logger) {
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		req, err := tunnel.DecodeRequest(msg)
		if err != nil {
			log.Warn("invalid request message", "error", err)
			continue
		}

		go handleRequest(conn, req, *ollamaURL, httpClient, log)
	}
}

func handleRequest(conn *websocket.Conn, req *tunnel.Request, ollamaURL string, httpClient *http.Client, log *slog.Logger) {
	targetURL := ollamaURL + req.Path

	body, err := base64.StdEncoding.DecodeString(req.Body)
	if err != nil {
		log.Warn("base64 decode failed", "error", err)
		sendError(conn, req.ID, "base64 decode failed", log)
		return
	}

	httpReq, err := http.NewRequest(req.Method, targetURL, strings.NewReader(string(body)))
	if err != nil {
		log.Warn("create request failed", "error", err)
		sendError(conn, req.ID, "create request failed", log)
		return
	}

	for k, v := range req.Headers {
		if k == "Host" || k == "Connection" || k == "Upgrade" {
			continue
		}
		httpReq.Header.Set(k, v)
	}

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		log.Warn("ollama request failed", "error", err)
		sendError(conn, req.ID, "ollama request failed: "+err.Error(), log)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warn("read ollama response failed", "error", err)
		sendError(conn, req.ID, "read response failed", log)
		return
	}

	respHeaders := make(map[string]string)
	for k := range resp.Header {
		respHeaders[k] = resp.Header.Get(k)
	}

	tunnelResp := &tunnel.Response{
		ID:      req.ID,
		Type:    "response",
		Status:  resp.StatusCode,
		Headers: respHeaders,
		Body:    base64.StdEncoding.EncodeToString(respBody),
	}

	data, err := tunnelResp.Encode()
	if err != nil {
		log.Warn("encode response failed", "error", err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Warn("write response failed", "error", err)
	}
}

func sendError(conn *websocket.Conn, id, msg string, log *slog.Logger) {
	resp := &tunnel.Response{
		ID:     id,
		Type:   "response",
		Status: 502,
		Body:   base64.StdEncoding.EncodeToString([]byte(msg)),
	}
	data, err := resp.Encode()
	if err != nil {
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Warn("write error response failed", "error", err)
	}
}
