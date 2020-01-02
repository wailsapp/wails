package interfaces

// CallbackFunc defines the signature of a function required to be provided to the
// Dispatch function so that the response may be returned
type CallbackFunc func(string) error

// IPCManager is the event manager interface
type IPCManager interface {
	BindRenderer(Renderer)
	Dispatch(message string, f CallbackFunc)
	Start(eventManager EventManager, bindingManager BindingManager)
	Shutdown()
}
