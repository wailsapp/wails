package commands

import (
	"bytes"
	"fmt"
	"os"
)

type DotDesktopOptions struct {
	OutputFile    string `description:"The output file to write to"`
	Type          string `description:"The type of the desktop entry" default:"Application"`
	Name          string `description:"The name of the application"`
	Exec          string `description:"The binary name + args to execute"`
	Icon          string `description:"The icon name or path for the application"`
	Categories    string `description:"Categories in which the application should be shown e.g. 'Development;IDE;'"`
	Comment       string `description:"A brief description of the application"`
	Terminal      bool   `description:"Whether the application runs in a terminal" default:"false"`
	Keywords      string `description:"Keywords associated with the application e.g. 'Editor;Image;'" default:"wails"`
	Version       string `description:"The version of the Desktop Entry Specification" default:"1.0"`
	GenericName   string `description:"A generic name for the application"`
	StartupNotify bool   `description:"If true, the app will send a notification when starting" default:"false"`
	MimeType      string `description:"The MIME types the application can handle e.g. 'image/gif;image/jpeg;'"`
	//Actions       []string `description:"Additional actions offered by the application"`
}

func (d *DotDesktopOptions) asBytes() []byte {
	var buf bytes.Buffer
	// Mandatory fields
	buf.WriteString("[Desktop Entry]\n")
	buf.WriteString(fmt.Sprintf("Type=%s\n", d.Type))
	buf.WriteString(fmt.Sprintf("Name=%s\n", d.Name))
	buf.WriteString(fmt.Sprintf("Exec=%s\n", d.Exec))

	// Optional fields with checks
	if d.Icon != "" {
		buf.WriteString(fmt.Sprintf("Icon=%s\n", d.Icon))
	}
	buf.WriteString(fmt.Sprintf("Categories=%s\n", d.Categories))
	if d.Comment != "" {
		buf.WriteString(fmt.Sprintf("Comment=%s\n", d.Comment))
	}
	buf.WriteString(fmt.Sprintf("Terminal=%t\n", d.Terminal))
	if d.Keywords != "" {
		buf.WriteString(fmt.Sprintf("Keywords=%s\n", d.Keywords))
	}
	if d.Version != "" {
		buf.WriteString(fmt.Sprintf("Version=%s\n", d.Version))
	}
	if d.GenericName != "" {
		buf.WriteString(fmt.Sprintf("GenericName=%s\n", d.GenericName))
	}
	buf.WriteString(fmt.Sprintf("StartupNotify=%t\n", d.StartupNotify))
	if d.MimeType != "" {
		buf.WriteString(fmt.Sprintf("MimeType=%s\n", d.MimeType))
	}
	return buf.Bytes()
}

func GenerateDotDesktop(options *DotDesktopOptions) error {
	DisableFooter = true

	if options.Name == "" {
		return fmt.Errorf("name is required")
	}

	options.Name = normaliseName(options.Name)

	if options.Exec == "" {
		return fmt.Errorf("exec is required")
	}

	if options.OutputFile == "" {
		options.OutputFile = options.Name + ".desktop"
	}

	// Write to file
	err := os.WriteFile(options.OutputFile, options.asBytes(), 0755)

	return err
}
