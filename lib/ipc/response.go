package ipc

import (
	"encoding/hex"
	"encoding/json"
)

// ipcResponse contains the response data from an RPC call
type ipcResponse struct {
	CallbackID   string      `json:"callbackid"`
	ErrorMessage string      `json:"error,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

// newErrorResponse returns the given error message to the frontend with the callbackid
func newErrorResponse(callbackID string, errorMessage string) *ipcResponse {
	// Create response object
	result := &ipcResponse{
		CallbackID:   callbackID,
		ErrorMessage: errorMessage,
	}
	return result
}

// newSuccessResponse returns the given data to the frontend with the callbackid
func newSuccessResponse(callbackID string, data interface{}) *ipcResponse {

	// Create response object
	result := &ipcResponse{
		CallbackID: callbackID,
		Data:       data,
	}

	return result
}

// Serialise formats the response to a string
func (i *ipcResponse) Serialise() (string, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	result := hex.EncodeToString(b)
	return result, err
}
