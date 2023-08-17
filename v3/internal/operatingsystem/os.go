package operatingsystem

// OS contains information about the operating system
type OS struct {
	ID      string
	Name    string
	Version string
}

// Info retrieves information about the current platform
func Info() (*OS, error) {
	return platformInfo()
}
