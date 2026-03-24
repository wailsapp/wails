package build

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/leaanthony/winicon"
	"github.com/tc-hib/winres"
	"github.com/tc-hib/winres/version"
	"github.com/wailsapp/wails/v2/internal/project"

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
	case "linux":
		err = packageApplicationForLinux(options)
	default:
		err = fmt.Errorf("packing not supported for %s yet", platform)
	}

	if err != nil {
		return err
	}

	return nil
}

// cleanBinDirectory will remove an existing bin directory and recreate it
func cleanBinDirectory(options *Options) error {
	buildDirectory := options.BinDirectory

	// Clear out old builds
	if fs.DirExists(buildDirectory) {
		err := os.RemoveAll(buildDirectory)
		if err != nil {
			return err
		}
	}

	// Create clean directory
	err := os.MkdirAll(buildDirectory, 0o700)
	if err != nil {
		return err
	}

	return nil
}

func packageApplicationForDarwin(options *Options) error {
	var err error

	// Create directory structure
	bundlename := options.BundleName
	if bundlename == "" {
		bundlename = options.ProjectData.Name + ".app"
	}

	contentsDirectory := filepath.Join(options.BinDirectory, bundlename, "/Contents")
	exeDir := filepath.Join(contentsDirectory, "/MacOS")
	err = fs.MkDirs(exeDir, 0o755)
	if err != nil {
		return err
	}
	resourceDir := filepath.Join(contentsDirectory, "/Resources")
	err = fs.MkDirs(resourceDir, 0o755)
	if err != nil {
		return err
	}
	// Copy binary
	packedBinaryPath := filepath.Join(exeDir, options.ProjectData.OutputFilename)
	err = fs.MoveFile(options.CompiledBinary, packedBinaryPath)
	if err != nil {
		return errors.Wrap(err, "Cannot move file: "+options.CompiledBinary)
	}

	// Generate Info.plist
	err = processPList(options, contentsDirectory)
	if err != nil {
		return err
	}

	// Generate App Icon
	err = processDarwinIcon(options.ProjectData, "appicon", resourceDir, "iconfile")
	if err != nil {
		return err
	}

	// Generate FileAssociation Icons
	for _, fileAssociation := range options.ProjectData.Info.FileAssociations {
		err = processDarwinIcon(options.ProjectData, fileAssociation.IconName, resourceDir, "")
		if err != nil {
			return err
		}
	}

	options.CompiledBinary = packedBinaryPath

	return nil
}

func processPList(options *Options, contentsDirectory string) error {
	sourcePList := "Info.plist"
	if options.Mode == Dev {
		// Use Info.dev.plist if using build mode
		sourcePList = "Info.dev.plist"
	}

	// Read the resolved BuildAssets file and copy it to the destination
	content, err := buildassets.ReadFileWithProjectData(options.ProjectData, "darwin/"+sourcePList)
	if err != nil {
		return err
	}

	targetFile := filepath.Join(contentsDirectory, "Info.plist")
	return os.WriteFile(targetFile, content, 0o644)
}

func processDarwinIcon(projectData *project.Project, iconName string, resourceDir string, destIconName string) (err error) {
	appIcon, err := buildassets.ReadFile(projectData, iconName+".png")
	if err != nil {
		return err
	}

	srcImg, _, err := image.Decode(bytes.NewBuffer(appIcon))
	if err != nil {
		return err
	}

	if destIconName == "" {
		destIconName = iconName
	}

	tgtBundle := filepath.Join(resourceDir, destIconName+".icns")
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
	// Generate app icon
	var err error
	err = generateIcoFile(options, "appicon", "icon")
	if err != nil {
		return err
	}

	// Generate FileAssociation Icons
	for _, fileAssociation := range options.ProjectData.Info.FileAssociations {
		err = generateIcoFile(options, fileAssociation.IconName, "")
		if err != nil {
			return err
		}
	}

	// Create syso file
	err = compileResources(options)
	if err != nil {
		return err
	}

	return nil
}

func packageApplicationForLinux(_ *Options) error {
	return nil
}

func generateIcoFile(options *Options, iconName string, destIconName string) error {
	content, err := buildassets.ReadFile(options.ProjectData, iconName+".png")
	if err != nil {
		return err
	}

	if destIconName == "" {
		destIconName = iconName
	}

	// Check ico file exists already
	icoFile := buildassets.GetLocalPath(options.ProjectData, "windows/"+destIconName+".ico")
	if !fs.FileExists(icoFile) {
		if dir := filepath.Dir(icoFile); !fs.DirExists(dir) {
			if err := fs.MkDirs(dir, 0o755); err != nil {
				return err
			}
		}

		output, err := os.OpenFile(icoFile, os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}
		defer output.Close()

		err = winicon.GenerateIcon(bytes.NewBuffer(content), output, []int{256, 128, 64, 48, 32, 16})
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
		_ = os.Chdir(currentDir)
	}()
	windowsDir := filepath.Join(options.ProjectData.GetBuildDir(), "windows")
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
		return fmt.Errorf("couldn't load icon from icon.ico: %w", err)
	}
	err = rs.SetIcon(winres.RT_ICON, ico)
	if err != nil {
		return err
	}

	manifestData, err := buildassets.ReadFileWithProjectData(options.ProjectData, "windows/wails.exe.manifest")
	if err != nil {
		return err
	}

	xmlData, err := winres.AppManifestFromXML(manifestData)
	if err != nil {
		return err
	}
	rs.SetManifest(xmlData)

	versionInfo, err := buildassets.ReadFileWithProjectData(options.ProjectData, "windows/info.json")
	if err != nil {
		return err
	}

	if len(versionInfo) != 0 {
		var v version.Info
		if err := v.UnmarshalJSON(versionInfo); err != nil {
			return err
		}
		rs.SetVersionInfo(v)
	}

	// replace spaces with underscores as go build behaves weirdly with spaces in syso filename
	targetFile := filepath.Join(options.ProjectData.Path, strings.ReplaceAll(options.ProjectData.Name, " ", "_")+"-res.syso")
	fout, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer fout.Close()

	archs := map[string]winres.Arch{
		"amd64": winres.ArchAMD64,
		"arm64": winres.ArchARM64,
		"386":   winres.ArchI386,
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
