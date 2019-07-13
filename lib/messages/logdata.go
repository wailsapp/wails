package messages

// LogData represents a call to log from the frontend
type LogData struct {
	Level   string `json:"level"`
	Message string `json:"string"`
}
