package operatingsystem

// OS contains information about the operating system
type OS struct {
	ID       string
	Name     string
	Version  string
	Branding string
}

func (o *OS) AsLogSlice() []any {
	return []any{
		"ID", o.ID,
		"Name", o.Name,
		"Version", o.Version,
		"Branding", o.Branding,
	}
}

// Info retrieves information about the current platform
func Info() (*OS, error) {
	return platformInfo()
}
