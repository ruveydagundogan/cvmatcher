package tunnel

import "encoding/json"

type Request struct {
	ID      string            `json:"id"`
	Type    string            `json:"type"`
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type Response struct {
	ID      string            `json:"id"`
	Type    string            `json:"type"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

func (r *Request) Encode() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Response) Encode() ([]byte, error) {
	return json.Marshal(r)
}

func DecodeRequest(data []byte) (*Request, error) {
	var r Request
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func DecodeResponse(data []byte) (*Response, error) {
	var r Response
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
