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

// PackageHelper helps with the 'wails package' command
type PackageHelper struct {
	fs     *FSHelper
	log    *Logger
	system *SystemHelper
}

// NewPackageHelper creates a new PackageHelper!
func NewPackageHelper() *PackageHelper {
	return &PackageHelper{
		fs:     NewFSHelper(),
		log:    NewLogger(),
		system: NewSystemHelper(),
	}
}

type plistData struct {
	Title     string
	Exe       string
	PackageID string
	Version   string
	Author    string
	Date      string
}

func newPlistData(title, exe, packageID, version, author string) *plistData {
	now := time.Now().Format(time.RFC822)
	return &plistData{
		Title:     title,
		Exe:       exe,
		Version:   version,
		PackageID: packageID,
		Author:    author,
		Date:      now,
	}
}

func defaultString(val string, defaultVal string) string {
	if val != "" {
		return val
	}
	return defaultVal
}

func (b *PackageHelper) getPackageFileBaseDir() string {
	// Calculate template base dir
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Join(path.Dir(filename), "packages", runtime.GOOS)
}

// Package the application into a platform specific package
func (b *PackageHelper) Package(po *ProjectOptions) error {
	switch runtime.GOOS {
	case "darwin":
		// Check we have the exe
		if !b.fs.FileExists(po.BinaryName) {
			return fmt.Errorf("cannot bundle non-existent binary file '%s'. Please build with 'wails build' first", po.BinaryName)
		}
		return b.packageOSX(po)
	case "windows":
		return b.PackageWindows(po, true)
	case "linux":
		return fmt.Errorf("linux is not supported at this time. Please see https://github.com/wailsapp/wails/issues/2")
	default:
		return fmt.Errorf("platform '%s' not supported for bundling yet", runtime.GOOS)
	}
}

// Package the application for OSX
func (b *PackageHelper) packageOSX(po *ProjectOptions) error {

	system := NewSystemHelper()
	config, err := system.LoadConfig()
	if err != nil {
		return err
	}

	name := defaultString(po.Name, "WailsTest")
	exe := defaultString(po.BinaryName, name)
	version := defaultString(po.Version, "0.1.0")
	author := defaultString(config.Name, "Anonymous")
	packageID := strings.Join([]string{"wails", name, version}, ".")
	plistData := newPlistData(name, exe, packageID, version, author)
	appname := po.Name + ".app"

	// Check binary exists
	source := path.Join(b.fs.Cwd(), exe)
	if !b.fs.FileExists(source) {
		// We need to build!
		return fmt.Errorf("Target '%s' not available. Has it been compiled yet?", exe)
	}

	// Remove the existing package
	os.RemoveAll(appname)

	exeDir := path.Join(b.fs.Cwd(), appname, "/Contents/MacOS")
	b.fs.MkDirs(exeDir, 0755)
	resourceDir := path.Join(b.fs.Cwd(), appname, "/Contents/Resources")
	b.fs.MkDirs(resourceDir, 0755)
	tmpl := template.New("infoPlist")
	plistFile := filepath.Join(b.getPackageFileBaseDir(), "info.plist")
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
	err = b.packageIconOSX(resourceDir)
	return err
}

// PackageWindows packages the application for windows platforms
func (b *PackageHelper) PackageWindows(po *ProjectOptions, cleanUp bool) error {
	basename := strings.TrimSuffix(po.BinaryName, ".exe")

	// Copy icon
	tgtIconFile := filepath.Join(b.fs.Cwd(), basename+".ico")
	if !b.fs.FileExists(tgtIconFile) {
		srcIconfile := filepath.Join(b.getPackageFileBaseDir(), "wails.ico")
		err := b.fs.CopyFile(srcIconfile, tgtIconFile)
		if err != nil {
			return err
		}
	}

	// Copy manifest
	tgtManifestFile := filepath.Join(b.fs.Cwd(), basename+".exe.manifest")
	if !b.fs.FileExists(tgtManifestFile) {
		srcManifestfile := filepath.Join(b.getPackageFileBaseDir(), "wails.exe.manifest")
		err := b.fs.CopyFile(srcManifestfile, tgtManifestFile)
		if err != nil {
			return err
		}
	}

	// Copy rc file
	tgtRCFile := filepath.Join(b.fs.Cwd(), basename+".rc")
	if !b.fs.FileExists(tgtRCFile) {
		srcRCfile := filepath.Join(b.getPackageFileBaseDir(), "wails.rc")
		rcfilebytes, err := ioutil.ReadFile(srcRCfile)
		if err != nil {
			return err
		}
		rcfiledata := strings.Replace(string(rcfilebytes), "$NAME$", basename, -1)
		err = ioutil.WriteFile(tgtRCFile, []byte(rcfiledata), 0755)
		if err != nil {
			return err
		}
	}

	// Build syso
	sysofile := filepath.Join(b.fs.Cwd(), basename+"-res.syso")
	windresCommand := []string{"windres", "-o", sysofile, tgtRCFile}
	err := NewProgramHelper().RunCommandArray(windresCommand)
	if err != nil {
		return err
	}

	// clean up
	if cleanUp {
		filesToDelete := []string{tgtIconFile, tgtManifestFile, tgtRCFile}
		err := b.fs.RemoveFiles(filesToDelete)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *PackageHelper) copyIcon(resourceDir string) (string, error) {

	// TODO: Read this from project.json
	const appIconFilename = "appicon.png"
	srcIcon := path.Join(b.fs.Cwd(), appIconFilename)

	// Check if appicon.png exists
	if !b.fs.FileExists(srcIcon) {

		// Install default icon
		iconfile := filepath.Join(b.getPackageFileBaseDir(), "icon.png")
		iconData, err := ioutil.ReadFile(iconfile)
		if err != nil {
			return "", err
		}
		err = ioutil.WriteFile(srcIcon, iconData, 0644)
		if err != nil {
			return "", err
		}
	}
	return srcIcon, nil
}

func (b *PackageHelper) packageIconOSX(resourceDir string) error {

	srcIcon, err := b.copyIcon(resourceDir)
	if err != nil {
		return err
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
	return icns.Encode(dest, srcImg)
}
