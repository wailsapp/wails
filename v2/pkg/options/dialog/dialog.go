package dialog

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
