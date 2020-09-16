package servicebus

// Message is a service bus message that contains a
// topic and data
type Message struct {
	topic  string
	data   interface{}
	target string
}

// NewMessage creates a new message with the given
// topic and data
func NewMessage(topic string, data interface{}) *Message {
	return &Message{
		topic: topic,
		data:  data,
	}
}

// NewMessageForTarget creates a new message with the given
// topic and data
func NewMessageForTarget(topic string, data interface{}, target string) *Message {
	return &Message{
		topic:  topic,
		data:   data,
		target: target,
	}
}

// Topic returns the message topic
func (m *Message) Topic() string {
	return m.topic
}

// Data returns the message data
func (m *Message) Data() interface{} {
	return m.data
}

// Target returns the message Target
func (m *Message) Target() string {
	return m.target
}
