package options

// OpenDialog contains the options for the OpenDialog runtime method
type OpenDialog struct {
	Title                      string
	Filter                     string
	AllowFiles                 bool
	AllowDirectories           bool
	AllowMultiple              bool
	ShowHiddenFiles            bool
	CanCreateDirectories       bool
	ResolveAliases             bool
	TreatPackagesAsDirectories bool
}
