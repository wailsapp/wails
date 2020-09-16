// +build linux

package build

import (
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/xyproto/xpm"
)

// compileIcon will compile the icon found at <projectdir>/icon.png into the application
func (d *DesktopBuilder) compileIcon(assetDir string, iconFile string) error {

	// Load icon into a databuffer
	targetFilename := "icon"
	targetFile := filepath.Join(assetDir, targetFilename+".c")

	d.addFileToDelete(targetFile)

	// Create a new XPM encoder
	enc := xpm.NewEncoder(targetFilename)

	// Open the PNG file
	f, err := os.Open(iconFile)
	if err != nil {
		return err
	}
	m, err := png.Decode(f)
	if err != nil {
		return err
	}
	f.Close()

	var buf strings.Builder

	// Generate and output the XPM data
	err = enc.Encode(&buf, m)
	if err != nil {
		return err
	}

	// Massage the output so we can extern reference it
	output := buf.String()
	output = strings.Replace(output, "static char", "const char", 1)

	// save icon.c
	err = ioutil.WriteFile(targetFile, []byte(output), 0755)

	return err
}
