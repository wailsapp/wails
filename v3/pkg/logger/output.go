package logger

type Output interface {
	Name() string
	Log(message *Message)
	Flush() error
}
