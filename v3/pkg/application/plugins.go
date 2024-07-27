package application

type Plugin interface {
	Name() string
	OnStartup() error
	OnShutdown() error
}
