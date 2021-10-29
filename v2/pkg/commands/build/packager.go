package build

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

	"github.com/leaanthony/winicon"
	"github.com/tc-hib/winres"

	"github.com/jackmordaunt/icns"
	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/pkg/buildassets"

	"github.com/wailsapp/wails/v2/internal/fs"
)

// PackageProject packages the application
func packageProject(options *Options, platform string) error {

	var err error
	switch platform {
	case "darwin":
		err = packageApplicationForDarwin(options)
	case "windows":
		err = packageApplicationForWindows(options)
	default:
		err = fmt.Errorf("packing not supported for %s yet", platform)
	}

	if err != nil {
		return err
	}

	return nil
}

// cleanBuildDirectory will remove an existing build directory and recreate it
func cleanBuildDirectory(options *Options) error {

	buildDirectory := options.BuildDirectory

	// Clear out old builds
	if fs.DirExists(buildDirectory) {
		err := os.RemoveAll(buildDirectory)
		if err != nil {
			return err
		}
	}

	// Create clean directory
	err := os.MkdirAll(buildDirectory, 0700)
	if err != nil {
		return err
	}

	return nil
}

// Gets (and creates) the build base directory
func getBuildBaseDirectory(options *Options) (string, error) {
	buildDirectory := filepath.Join(options.ProjectData.Path, "build")
	if !fs.DirExists(buildDirectory) {
		err := os.MkdirAll(buildDirectory, 0700)
		if err != nil {
			return "", err
		}
	}
	return buildDirectory, nil
}

// Gets the platform dependent package assets directory
func getPackageAssetsDirectory() string {
	return fs.RelativePath("internal/packager", runtime.GOOS)
}

func packageApplicationForDarwin(options *Options) error {

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

	// Generate Info.plist
	err = processPList(options, contentsDirectory)
	if err != nil {
		return err
	}

	// Generate Icons
	err = processApplicationIcon(resourceDir, options.ProjectData.BuildDir)
	if err != nil {
		return err
	}

	options.CompiledBinary = packedBinaryPath

	return nil
}

func processPList(options *Options, contentsDirectory string) error {

	// Check if plist already exists in project dir
	plistFile := filepath.Join(options.ProjectData.BuildDir, "darwin", "Info.plist")

	// If the file doesn't exist, generate it
	if !fs.FileExists(plistFile) {
		err := generateDefaultPlist(options, plistFile)
		if err != nil {
			return err
		}
	}

	// Copy it to the contents directory
	targetFile := filepath.Join(contentsDirectory, "Info.plist")
	return fs.CopyFile(plistFile, targetFile)
}

func generateDefaultPlist(options *Options, targetPlistFile string) error {
	name := defaultString(options.ProjectData.Name, "WailsTest")
	exe := defaultString(options.OutputFile, name)
	version := "1.0.0"
	author := defaultString(options.ProjectData.Author.Name, "Anonymous")
	packageID := strings.Join([]string{"wails", name}, ".")
	plistData := newPlistData(name, exe, packageID, version, author)

	tmpl := template.New("infoPlist")
	plistTemplate := fs.RelativePath("./internal/packager/darwin/Info.plist")
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

	// Create the directory if it doesn't exist
	err = fs.MkDirs(filepath.Dir(targetPlistFile))
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
		err = buildassets.RegenerateAppIcon(appIcon)
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

func packageApplicationForWindows(options *Options) error {
	// Generate icon
	var err error
	err = generateIcoFile(options)
	if err != nil {
		return err
	}

	// Ensure Manifest is present
	err = generateManifest(options)
	if err != nil {
		return err
	}

	// Create syso file
	err = compileResources(options)
	if err != nil {
		return err
	}

	return nil
}

func generateManifest(options *Options) error {
	filename := options.ProjectData.Name + ".exe.manifest"
	manifestFile := filepath.Join(options.ProjectData.Path, "build", "windows", filename)
	if !fs.FileExists(manifestFile) {
		return buildassets.RegenerateManifest(manifestFile)
	}
	return nil
}

func generateIcoFile(options *Options) error {
	// Check ico file exists already
	icoFile := filepath.Join(options.ProjectData.Path, "build", "windows", "icon.ico")
	if !fs.FileExists(icoFile) {
		// Check icon exists
		appicon := filepath.Join(options.ProjectData.Path, "build", "appicon.png")
		if !fs.FileExists(appicon) {
			return fmt.Errorf("application icon missing: %s", appicon)
		}
		// Load icon
		input, err := os.Open(appicon)
		if err != nil {
			return err
		}
		output, err := os.OpenFile(icoFile, os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		err = winicon.GenerateIcon(input, output, []int{256, 128, 64, 48, 32, 16})
		if err != nil {
			return err
		}
	}
	return nil
}

func compileResources(options *Options) error {

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer func() {
		os.Chdir(currentDir)
	}()
	windowsDir := filepath.Join(options.ProjectData.Path, "build", "windows")
	err = os.Chdir(windowsDir)
	if err != nil {
		return err
	}
	rs := winres.ResourceSet{}
	icon := filepath.Join(windowsDir, "icon.ico")
	iconFile, err := os.Open(icon)
	if err != nil {
		return err
	}
	defer iconFile.Close()
	ico, err := winres.LoadICO(iconFile)
	if err != nil {
		return err
	}
	err = rs.SetIcon(winres.RT_ICON, ico)
	if err != nil {
		return err
	}

	ManifestFilename := options.ProjectData.Name + ".exe.manifest"
	manifestData, err := os.ReadFile(ManifestFilename)
	xmlData, err := winres.AppManifestFromXML(manifestData)
	if err != nil {
		return err
	}
	rs.SetManifest(xmlData)

	targetFile := filepath.Join(options.ProjectData.Path, options.ProjectData.Name+"-res.syso")
	fout, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer fout.Close()

	archs := map[string]winres.Arch{
		"amd64": winres.ArchAMD64,
	}
	targetArch, supported := archs[options.Arch]
	if !supported {
		return fmt.Errorf("arch '%s' not supported", options.Arch)
	}

	err = rs.WriteObject(fout, targetArch)
	if err != nil {
		return err
	}
	return nil
}
