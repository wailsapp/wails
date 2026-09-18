package application

import (
	"errors"
	"strings"
	"unsafe"
)

var (
	// ErrDialogNotSupported is returned by the native panel helpers (text
	// prompt, colour panel, font panel) on platforms that have no native
	// implementation for them.
	ErrDialogNotSupported = errors.New("this dialog is not supported on this platform")
	// ErrDialogInProgress is returned when a shared system panel (colour or
	// font) is already open on behalf of an earlier call. AppKit has one
	// shared instance of each, so only one pick can run at a time.
	ErrDialogInProgress = errors.New("a native panel of this kind is already open")
)

// PromptOptions configures DialogManager.Prompt, a native alert with a text
// field as its accessory view.
type PromptOptions struct {
	Title        string
	Message      string
	Placeholder  string
	DefaultValue string
	// Secure renders the field as a password field.
	Secure      bool
	OKLabel     string
	CancelLabel string
	// Window, when set, presents the prompt as a sheet on that window.
	Window Window
}

// ColorPickerOptions configures DialogManager.PickColor.
type ColorPickerOptions struct {
	Initial    RGBA
	ShowsAlpha bool
	Title      string
	// OnChange, when set, is called with every colour the user selects while
	// the panel is open. The final colour is returned by PickColor.
	OnChange func(RGBA)
}

// FontDescriptor identifies a font chosen from the native font panel.
type FontDescriptor struct {
	Family         string
	Face           string
	PostScriptName string
	Size           float64
}

// FontPickerOptions configures DialogManager.PickFont.
type FontPickerOptions struct {
	// Family is the family (or PostScript) name of the font initially selected
	// in the panel. Empty selects the system font.
	Family string
	// Size is the initial point size. Zero uses the system font size.
	Size float64
	// OnChange, when set, is called with every font the user selects while
	// the panel is open. The final font is returned by PickFont.
	OnChange func(FontDescriptor)
}

// DialogFormat is one entry of the "Format:" pop-up shown as the accessory
// view of a save dialog (see SaveFileDialogStruct.SetFormats).
type DialogFormat struct {
	// Label is the text shown in the pop-up.
	Label string
	// Extension is applied to the name field when the format is chosen, and
	// used as the allowed type when UTI is empty. No leading dot.
	Extension string
	// UTI is an optional uniform type identifier used as the allowed content
	// type when the format is chosen (EG: "com.adobe.pdf").
	UTI string
}

// SetSuppression shows a "Do not show this message again" checkbox on the
// alert. When label is non-empty it replaces the default system title.
// Read the result with Suppressed or OnSuppression. Supported on macOS;
// ignored elsewhere.
func (d *MessageDialog) SetSuppression(label string) *MessageDialog {
	d.showsSuppression = true
	d.suppressionLabel = label
	return d
}

// OnSuppression registers a callback invoked with the state of the
// suppression checkbox when the alert is dismissed. It runs before the
// pressed button's callback.
func (d *MessageDialog) OnSuppression(callback func(bool)) *MessageDialog {
	d.onSuppression = callback
	return d
}

// Suppressed reports whether the suppression checkbox was ticked when the
// alert was dismissed. It is false until the alert has been dismissed, so
// read it from a button callback or from OnSuppression.
func (d *MessageDialog) Suppressed() bool {
	return d.suppressed
}

// SetHelp shows the alert's help button and calls callback when it is
// pressed. Supported on macOS; ignored elsewhere.
func (d *MessageDialog) SetHelp(callback func()) *MessageDialog {
	d.helpCallback = callback
	return d
}

// SetAccessoryView installs a native view as the alert's accessory view.
// On macOS the pointer must be an NSView; the alert retains it. This is an
// escape hatch for native integrations; ignored elsewhere.
func (d *MessageDialog) SetAccessoryView(view unsafe.Pointer) *MessageDialog {
	d.accessoryView = view
	return d
}

func (d *MessageDialog) setSuppressed(suppressed bool) {
	d.suppressed = suppressed
	if d.onSuppression != nil {
		d.onSuppression(suppressed)
	}
}

// AddContentType restricts the dialog to files conforming to a uniform type
// identifier (EG: "public.image", "com.adobe.pdf"). On macOS this maps
// directly to the panel's allowed content types and coexists with the
// extension filters added by AddFilter. Elsewhere the identifier is
// translated to an extension filter where a translation is known and
// otherwise ignored.
func (d *OpenFileDialogStruct) AddContentType(uti string) *OpenFileDialogStruct {
	uti = strings.TrimSpace(uti)
	if uti == "" {
		return d
	}
	d.contentTypes = append(d.contentTypes, uti)
	if !dialogContentTypesNative {
		if filter, ok := contentTypeFilter(uti); ok {
			d.filters = append(d.filters, filter)
		} else if globalApplication != nil {
			globalApplication.debug("OpenFileDialog: no extension translation for content type, ignoring", "uti", uti)
		}
	}
	return d
}

// AddContentType restricts the dialog to a uniform type identifier. See
// OpenFileDialogStruct.AddContentType.
func (d *SaveFileDialogStruct) AddContentType(uti string) *SaveFileDialogStruct {
	uti = strings.TrimSpace(uti)
	if uti == "" {
		return d
	}
	d.contentTypes = append(d.contentTypes, uti)
	if !dialogContentTypesNative {
		if filter, ok := contentTypeFilter(uti); ok {
			d.filters = append(d.filters, filter)
		} else if globalApplication != nil {
			globalApplication.debug("SaveFileDialog: no extension translation for content type, ignoring", "uti", uti)
		}
	}
	return d
}

// SetFormats adds a "Format:" pop-up to the save dialog. Choosing an entry
// swaps the panel's allowed type and the extension of the name field and
// calls onChange (which may be nil) with the new index. selected is the
// initially selected index; out of range values select the first entry.
// After the prompt returns, SelectedFormat reports the final choice.
// Supported on macOS; elsewhere the formats become extension filters.
func (d *SaveFileDialogStruct) SetFormats(formats []DialogFormat, selected int, onChange func(index int)) *SaveFileDialogStruct {
	d.formats = append([]DialogFormat(nil), formats...)
	if selected < 0 || selected >= len(d.formats) {
		selected = 0
	}
	d.selectedFormat = selected
	d.onFormatChange = onChange
	if !dialogContentTypesNative {
		for _, format := range d.formats {
			ext := strings.TrimPrefix(strings.TrimSpace(format.Extension), ".")
			if ext == "" {
				continue
			}
			d.filters = append(d.filters, FileFilter{
				DisplayName: format.Label,
				Pattern:     "*." + ext,
			})
		}
	}
	return d
}

// SelectedFormat returns the index into the formats passed to SetFormats
// that was selected when the dialog closed (or the initial selection while
// it is open). It is -1 when SetFormats was not called.
func (d *SaveFileDialogStruct) SelectedFormat() int {
	if len(d.formats) == 0 {
		return -1
	}
	return d.selectedFormat
}

// SetNameFieldLabel sets the label shown before the file name field
// ("Save As:" by default). Supported on macOS; ignored elsewhere.
func (d *SaveFileDialogStruct) SetNameFieldLabel(label string) *SaveFileDialogStruct {
	d.nameFieldLabel = label
	return d
}

// SetTags shows the Finder tags field pre-populated with the given tag
// names. Supported on macOS; ignored elsewhere.
func (d *SaveFileDialogStruct) SetTags(tags []string) *SaveFileDialogStruct {
	d.tags = nil
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			d.tags = append(d.tags, tag)
		}
	}
	return d
}

func (d *SaveFileDialogStruct) setSelectedFormat(index int) {
	if index < 0 || index >= len(d.formats) {
		return
	}
	d.selectedFormat = index
	if d.onFormatChange != nil {
		go func() {
			defer handlePanic()
			d.onFormatChange(index)
		}()
	}
}

// contentTypeExtensions translates common uniform type identifiers to file
// extensions for platforms without native UTI support.
var contentTypeExtensions = map[string][]string{
	"public.png":                         {"png"},
	"public.jpeg":                        {"jpg", "jpeg"},
	"com.compuserve.gif":                 {"gif"},
	"public.tiff":                        {"tif", "tiff"},
	"com.microsoft.bmp":                  {"bmp"},
	"public.heic":                        {"heic"},
	"public.svg-image":                   {"svg"},
	"org.webmproject.webp":               {"webp"},
	"com.microsoft.ico":                  {"ico"},
	"public.image":                       {"png", "jpg", "jpeg", "gif", "tif", "tiff", "bmp", "heic", "webp", "svg"},
	"com.adobe.pdf":                      {"pdf"},
	"public.plain-text":                  {"txt"},
	"public.utf8-plain-text":             {"txt"},
	"public.text":                        {"txt", "md", "rtf", "html", "htm", "csv", "json", "xml"},
	"public.html":                        {"html", "htm"},
	"public.json":                        {"json"},
	"public.xml":                         {"xml"},
	"public.yaml":                        {"yaml", "yml"},
	"public.comma-separated-values-text": {"csv"},
	"public.rtf":                         {"rtf"},
	"net.daringfireball.markdown":        {"md", "markdown"},
	"public.zip-archive":                 {"zip"},
	"org.gnu.gnu-zip-archive":            {"gz"},
	"public.tar-archive":                 {"tar"},
	"public.mp3":                         {"mp3"},
	"public.mpeg-4":                      {"mp4", "m4v"},
	"public.mpeg-4-audio":                {"m4a"},
	"com.microsoft.waveform-audio":       {"wav"},
	"public.aiff-audio":                  {"aif", "aiff"},
	"org.xiph.flac":                      {"flac"},
	"public.audio":                       {"mp3", "m4a", "wav", "aac", "flac", "ogg", "aif", "aiff"},
	"com.apple.quicktime-movie":          {"mov"},
	"public.avi":                         {"avi"},
	"org.webmproject.webm":               {"webm"},
	"public.movie":                       {"mp4", "m4v", "mov", "avi", "mkv", "webm"},
	"org.openxmlformats.wordprocessingml.document":   {"docx"},
	"org.openxmlformats.spreadsheetml.sheet":         {"xlsx"},
	"org.openxmlformats.presentationml.presentation": {"pptx"},
	"com.microsoft.word.doc":                         {"doc"},
	"com.microsoft.excel.xls":                        {"xls"},
	"com.microsoft.powerpoint.ppt":                   {"ppt"},
	"public.source-code":                             {"c", "h", "cpp", "hpp", "go", "js", "ts", "py", "rs", "java", "swift", "m"},
	"public.shell-script":                            {"sh"},
	"public.python-script":                           {"py"},
	"public.c-source":                                {"c"},
	"public.c-header":                                {"h"},
	"com.netscape.javascript-source":                 {"js"},
	"public.go-source":                               {"go"},
	"public.swift-source":                            {"swift"},
}

// contentTypeFilter translates a uniform type identifier into an extension
// filter. It reports false when no translation is known; "public.item",
// "public.data" and "public.content" match every file and are reported as
// a wildcard filter.
func contentTypeFilter(uti string) (FileFilter, bool) {
	uti = strings.ToLower(strings.TrimSpace(uti))
	switch uti {
	case "":
		return FileFilter{}, false
	case "public.item", "public.data", "public.content":
		return FileFilter{DisplayName: "All Files", Pattern: "*"}, true
	}
	extensions, ok := contentTypeExtensions[uti]
	if !ok || len(extensions) == 0 {
		return FileFilter{}, false
	}
	patterns := make([]string, len(extensions))
	for i, ext := range extensions {
		patterns[i] = "*." + ext
	}
	return FileFilter{
		DisplayName: uti + " (" + strings.Join(patterns, ", ") + ")",
		Pattern:     strings.Join(patterns, ";"),
	}, true
}
