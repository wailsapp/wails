package cmd

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/jackmordaunt/icns"
)

// BundleHelper helps with the 'wails bundle' command
type BundleHelper struct {
	fs     *FSHelper
	log    *Logger
	system *SystemHelper
}

// NewBundleHelper creates a new BundleHelper!
func NewBundleHelper() *BundleHelper {
	return &BundleHelper{
		fs:     NewFSHelper(),
		log:    NewLogger(),
		system: NewSystemHelper(),
	}
}

// var assetsBox packr.Box

type plistData struct {
	Title    string
	Exe      string
	BundleID string
	Version  string
	Author   string
	Date     string
}

func newPlistData(title, exe, bundleID, version, author string) *plistData {
	now := time.Now().Format(time.RFC822)
	return &plistData{
		Title:    title,
		Exe:      exe,
		Version:  version,
		BundleID: bundleID,
		Author:   author,
		Date:     now,
	}
}

func defaultString(val string, defaultVal string) string {
	if val != "" {
		return val
	}
	return defaultVal
}

func (b *BundleHelper) getBundleFileBaseDir() string {
	return filepath.Join(b.system.homeDir, "go", "src", "github.com", "wailsapp", "wails", "cmd", "bundle", runtime.GOOS)
}

// Bundle the application into a platform specific package
func (b *BundleHelper) Bundle(po *ProjectOptions) error {
	// Check we have the exe
	if !b.fs.FileExists(po.BinaryName) {
		return fmt.Errorf("cannot bundle non-existant binary file '%s'. Please build with 'wails build' first", po.BinaryName)
	}
	switch runtime.GOOS {
	case "darwin":
		return b.bundleOSX(po)
	default:
		return fmt.Errorf("platform '%s' not supported for bundling yet", runtime.GOOS)
	}
}

// Bundle the application
func (b *BundleHelper) bundleOSX(po *ProjectOptions) error {

	system := NewSystemHelper()
	config, err := system.LoadConfig()
	if err != nil {
		return err
	}

	name := defaultString(po.Name, "WailsTest")
	exe := defaultString(po.BinaryName, name)
	version := defaultString(po.Version, "0.1.0")
	author := defaultString(config.Name, "Anonymous")
	bundleID := strings.Join([]string{"wails", name, version}, ".")
	plistData := newPlistData(name, exe, bundleID, version, author)
	appname := po.Name + ".app"

	// Check binary exists
	source := path.Join(b.fs.Cwd(), exe)
	if !b.fs.FileExists(source) {
		// We need to build!
		return fmt.Errorf("Target '%s' not available. Has it been compiled yet?", exe)
	}

	// REmove the existing bundle
	os.RemoveAll(appname)

	exeDir := path.Join(b.fs.Cwd(), appname, "/Contents/MacOS")
	b.fs.MkDirs(exeDir, 0755)
	resourceDir := path.Join(b.fs.Cwd(), appname, "/Contents/Resources")
	b.fs.MkDirs(resourceDir, 0755)
	tmpl := template.New("infoPlist")
	plistFile := filepath.Join(b.getBundleFileBaseDir(), "info.plist")
	infoPlist, err := ioutil.ReadFile(plistFile)
	if err != nil {
		return err
	}
	tmpl.Parse(string(infoPlist))

	// Write the template to a buffer
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, plistData)
	if err != nil {
		return err
	}
	filename := path.Join(b.fs.Cwd(), appname, "Contents", "Info.plist")
	err = ioutil.WriteFile(filename, tpl.Bytes(), 0644)
	if err != nil {
		return err
	}

	// Copy executable
	target := path.Join(exeDir, exe)
	err = b.fs.CopyFile(source, target)
	if err != nil {
		return err
	}

	err = os.Chmod(target, 0755)
	if err != nil {
		return err
	}
	err = b.bundleIcon(resourceDir)
	return err
}

func (b *BundleHelper) bundleIcon(resourceDir string) error {

	// TODO: Read this from project.json
	const appIconFilename = "appicon.png"

	srcIcon := path.Join(b.fs.Cwd(), appIconFilename)

	// Check if appicon.png exists
	if !b.fs.FileExists(srcIcon) {

		// Install default icon
		iconfile := filepath.Join(b.getBundleFileBaseDir(), "icon.png")
		iconData, err := ioutil.ReadFile(iconfile)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(srcIcon, iconData, 0644)
		if err != nil {
			return err
		}
	}

	tgtBundle := path.Join(resourceDir, "iconfile.icns")
	imageFile, err := os.Open(srcIcon)
	if err != nil {
		return err
	}
	defer imageFile.Close()
	srcImg, _, err := image.Decode(imageFile)
	if err != nil {
		return err

	}
	dest, err := os.Create(tgtBundle)
	if err != nil {
		return err

	}
	defer dest.Close()
	if err := icns.Encode(dest, srcImg); err != nil {
		return err

	}
	return nil
}
