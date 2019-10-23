package interfaces

import "github.com/wailsapp/wails/lib/messages"

// BindingManager is the binding manager interface
type BindingManager interface {
	Bind(object interface{})
	Start(renderer Renderer, runtime Runtime) error
	ProcessCall(callData *messages.CallData) (result interface{}, err error)
	Shutdown()
}
