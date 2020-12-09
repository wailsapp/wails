// +build darwin

package build

import (
	"fmt"
	"github.com/wailsapp/wails/v2/internal/fs"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// desktop_linux.go will compile the tray icon found at <projectdir>/trayicon.png into the application
func (d *DesktopBuilder) processTrayIcons(assetDir string, options *Options) error {

	// Determine icon file
	iconFile := filepath.Join(options.ProjectData.Path, "trayicon.png")

	var err error

	// Setup target
	targetFilename := "trayicon"
	targetFile := filepath.Join(assetDir, targetFilename+".c")
	//d.addFileToDelete(targetFile)

	var dataBytes []byte

	// If the icon file exists, load it up
	if fs.FileExists(iconFile) {
		// Load the tray icon
		dataBytes, err = ioutil.ReadFile(iconFile)
		if err != nil {
			return err
		}
	}

	// Use a strings builder
	var cdata strings.Builder

	// Write header
	header := `// trayicon.c
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL.
// This file was auto-generated. DO NOT MODIFY.

`
	cdata.WriteString(header)
	cdata.WriteString(fmt.Sprintf("const unsigned int trayIconLength = %d;\n", len(dataBytes)))
	cdata.WriteString("const unsigned char trayIcon[] = { ")

	// Convert each byte to hex
	for _, b := range dataBytes {
		cdata.WriteString(fmt.Sprintf("0x%x, ", b))
	}

	cdata.WriteString("0x00 };\n")

	err = ioutil.WriteFile(targetFile, []byte(cdata.String()), 0600)
	if err != nil {
		return err
	}
	return nil
}
