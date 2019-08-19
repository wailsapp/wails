package messages

// EventData represents an event sent from the frontend
type EventData struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}
