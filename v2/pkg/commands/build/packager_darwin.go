package build

import (
	"bytes"
	"image"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/jackmordaunt/icns"
	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/internal/fs"
)

func packageApplication(options *Options) error {

	var err error

	// Create directory structure
	bundlename := options.ProjectData.Name + ".app"

	contentsDirectory := filepath.Join(options.BuildDirectory, bundlename, "/Contents")
	exeDir := filepath.Join(contentsDirectory, "/MacOS")
	err = fs.MkDirs(exeDir, 0755)
	if err != nil {
		return err
	}
	resourceDir := filepath.Join(contentsDirectory, "/Resources")
	err = fs.MkDirs(resourceDir, 0755)
	if err != nil {
		return err
	}
	// Copy binary
	packedBinaryPath := filepath.Join(exeDir, options.ProjectData.Name)
	err = fs.MoveFile(options.CompiledBinary, packedBinaryPath)
	if err != nil {
		return errors.Wrap(err, "Cannot move file: "+options.ProjectData.OutputFilename)
	}

	// Generate info.plist
	err = processPList(options, contentsDirectory)
	if err != nil {
		return err
	}

	// Generate Icons
	err = processApplicationIcon(resourceDir, options.ProjectData.IconsDir)
	if err != nil {
		return err
	}

	return nil
}

func processPList(options *Options, contentsDirectory string) error {
	// Check if plist already exists in project dir
	plistFile, err := fs.RelativeToCwd("info.plist")
	if err != nil {
		return err
	}

	// If the file doesn't exist, generate it
	if !fs.FileExists(plistFile) {
		err = generateDefaultPlist(options, plistFile)
		if err != nil {
			return err
		}
	}

	// Copy it to the contents directory
	targetFile := filepath.Join(contentsDirectory, "info.plist")
	return fs.CopyFile(plistFile, targetFile)
}

func generateDefaultPlist(options *Options, targetPlistFile string) error {
	name := defaultString(options.ProjectData.Name, "WailsTest")
	exe := defaultString(options.OutputFile, name)
	version := "1.0.0"
	author := defaultString(options.ProjectData.Author.Name, "Anonymous")
	packageID := strings.Join([]string{"wails", name, version}, ".")
	plistData := newPlistData(name, exe, packageID, version, author)

	tmpl := template.New("infoPlist")
	plistTemplate := fs.RelativePath("./internal/packager/darwin/info.plist")
	infoPlist, err := ioutil.ReadFile(plistTemplate)
	if err != nil {
		return errors.Wrap(err, "Cannot open plist template")
	}
	_, err = tmpl.Parse(string(infoPlist))
	if err != nil {
		return err
	}
	// Write the template to a buffer
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, plistData)
	if err != nil {
		return err
	}

	// Save the file
	return ioutil.WriteFile(targetPlistFile, tpl.Bytes(), 0644)
}

func defaultString(val string, defaultVal string) string {
	if val != "" {
		return val
	}
	return defaultVal
}

type plistData struct {
	Title     string
	Exe       string
	PackageID string
	Version   string
	Author    string
}

func newPlistData(title, exe, packageID, version, author string) *plistData {
	return &plistData{
		Title:     title,
		Exe:       exe,
		Version:   version,
		PackageID: packageID,
		Author:    author,
	}
}

func processApplicationIcon(resourceDir string, iconsDir string) (err error) {

	appIcon := filepath.Join(iconsDir, "appicon.png")

	// Install default icon if one doesn't exist
	if !fs.FileExists(appIcon) {
		// No - Install default icon
		defaultIcon := fs.RelativePath("./internal/packager/icon1024.png")
		err = fs.CopyFile(defaultIcon, appIcon)
		if err != nil {
			return
		}
	}

	tgtBundle := path.Join(resourceDir, "iconfile.icns")
	imageFile, err := os.Open(appIcon)
	if err != nil {
		return err
	}

	defer func() {
		err = imageFile.Close()
		if err == nil {
			return
		}
	}()
	srcImg, _, err := image.Decode(imageFile)
	if err != nil {
		return err

	}
	dest, err := os.Create(tgtBundle)
	if err != nil {
		return err

	}
	defer func() {
		err = dest.Close()
		if err == nil {
			return
		}
	}()
	return icns.Encode(dest, srcImg)
}
