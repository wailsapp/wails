package frontend

type Dispatcher interface {
	ProcessMessage(message string, sender Frontend) (string, error)
}
