package application

// ChildView is a native view that can be owned by a Window. It is intended for
// advanced integrations such as embedded browser surfaces. Implementations must
// make Attach and Detach safe to call once during the parent window lifecycle.
type ChildView interface {
	Attach(Window)
	Detach()
}
