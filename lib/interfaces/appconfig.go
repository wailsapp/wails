package interfaces

// AppConfig is the application config interface
type AppConfig interface {
	GetWidth() int
	GetHeight() int
	GetTitle() string
	GetMinWidth() int
	GetMinHeight() int
	GetMaxWidth() int
	GetMaxHeight() int
	GetResizable() bool
	GetHTML() string
	GetDisableInspector() bool
	GetColour() string
	GetCSS() string
	GetJS() string
}
