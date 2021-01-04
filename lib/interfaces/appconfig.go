package interfaces

// AppConfig is the application config interface
type AppConfig interface {
	GetWidth() int
	GetHeight() int
	GetTitle() string
	GetResizable() bool
	GetHTML() string
	GetDisableInspector() bool
	GetColour() string
	GetCSS() string
	GetJS() string
}
