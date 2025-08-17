package application

import (
	"encoding/json"
	"io"
)

type runtimeRequest struct {
	Object int         `json:"object,omitempty"`
	Method int         `json:"method,omitempty"`
	Params interface{} `json:"params,omitempty"`
}

func (mp *MessageProcessor) decodeBody(body io.ReadCloser) (*runtimeRequest, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		body.Close()
		return nil, err
	}
	body.Close()

	if len(data) == 0 {
		return nil, io.EOF
	}

	var req runtimeRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	return &req, nil
}
