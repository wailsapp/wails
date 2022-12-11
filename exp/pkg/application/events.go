package application

var applicationEvents = make(chan uint)

type WindowEvent struct {
	WindowID uint
	EventID  uint
}

var windowEvents = make(chan *WindowEvent)
