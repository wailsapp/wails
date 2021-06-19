package dialog

// FileFilter defines a filter for dialog boxes
type FileFilter struct {
	DisplayName string // Filter information EG: "Image Files (*.jpg, *.png)"
	Pattern     string // semi-colon separated list of extensions, EG: "*.jpg;*.png"
}

// OpenDialog contains the options for the OpenDialog runtime method
type OpenDialog struct {
	DefaultDirectory           string
	DefaultFilename            string
	Title                      string
	Filters                    []FileFilter
	AllowFiles                 bool
	AllowDirectories           bool
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
	Filters                    []FileFilter
	ShowHiddenFiles            bool
	CanCreateDirectories       bool
	TreatPackagesAsDirectories bool
}

type DialogType string

const (
	InfoDialog     DialogType = "info"
	WarningDialog  DialogType = "warning"
	ErrorDialog    DialogType = "error"
	QuestionDialog DialogType = "question"
)

// MessageDialog contains the options for the Message dialogs, EG Info, Warning, etc runtime methods
type MessageDialog struct {
	Type          DialogType
	Title         string
	Message       string
	Buttons       []string
	DefaultButton string
	CancelButton  string
	Icon          string
}
