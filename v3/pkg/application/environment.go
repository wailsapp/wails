package application

// EnvironmentInfo represents information about the current environment.
//
// Fields:
// - OS: the operating system that the program is running on.
// - Arch: the architecture of the operating system.
// - Debug: indicates whether debug mode is enabled.
type EnvironmentInfo struct {
	OS    string
	Arch  string
	Debug bool
}
