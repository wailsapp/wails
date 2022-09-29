package application

type EventType int

const (
	StartUp EventType = iota
	ShutDown
	DomReady
)
