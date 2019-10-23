package interfaces

// IPCManager is the event manager interface
type IPCManager interface {
	BindRenderer(Renderer)
	Dispatch(message string)
	Start(eventManager EventManager, bindingManager BindingManager)
	Shutdown()
}
