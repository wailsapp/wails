package frontend

type Dispatcher interface {
	ProcessMessage(message string) (string, error)
}
