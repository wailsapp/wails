package options

// OpenDialog contains the options for the OpenDialog runtime method
type OpenDialog struct {
	DefaultDirectory           string
	DefaultFilename            string
	Title                      string
	Filters                    string
	AllowFiles                 bool
	AllowDirectories           bool
	AllowMultiple              bool
	ShowHiddenFiles            bool
	CanCreateDirectories       bool
	ResolvesAliases            bool
	TreatPackagesAsDirectories bool
}

// SaveDialog contains the options for the SaveDialog runtime method
type SaveDialog struct {
	DefaultDirectory           string
	DefaultFilename            string
	Title                      string
	Filters                    string
	ShowHiddenFiles            bool
	CanCreateDirectories       bool
	TreatPackagesAsDirectories bool
}
