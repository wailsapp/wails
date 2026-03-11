package application

type theme int

// theme set to internal unexported enum
const (
	// systemDefault will use whatever the system theme is. The application will follow system theme changes.
	systemDefault theme = 0
	// dark Mode
	dark theme = 1
	// light Mode
	light theme = 2
)
