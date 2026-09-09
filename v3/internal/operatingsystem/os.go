package operatingsystem

// OS contains information about the operating system
type OS struct {
	ID       string `json:"ID"`
	Name     string `json:"Name"`
	Version  string `json:"Version"`
	Branding string `json:"Branding"`
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
