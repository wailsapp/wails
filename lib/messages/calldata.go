package messages

// CallData represents a call to a Go function/method
type CallData struct {
	BindingName string `json:"bindingName"`
	Data        string `json:"data,omitempty"`
}
