package application

var blankWindowEventContext = &WindowEventContext{}

type WindowEventContext struct {
	// contains filtered or unexported fields
	data map[string]any
}

func newWindowEventContext() *Context {
	return &Context{
		data: make(map[string]any),
	}
}
