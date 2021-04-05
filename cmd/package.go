package cmd

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/jackmordaunt/icns"
	"golang.org/x/image/draw"
)

// PackageHelper helps with the 'wails package' command
type PackageHelper struct {
	platform string
	fs       *FSHelper
	log      *Logger
	system   *SystemHelper
}

// NewPackageHelper creates a new PackageHelper!
func NewPackageHelper(platform string) *PackageHelper {
	return &PackageHelper{
		platform: platform,
		fs:       NewFSHelper(),
		log:      NewLogger(),
		system:   NewSystemHelper(),
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

type windowsIcoHeader struct {
	_          uint16
	imageType  uint16
	imageCount uint16
}

type windowsIcoDescriptor struct {
	width   uint8
	height  uint8
	colours uint8
	_       uint8
	planes  uint16
	bpp     uint16
	size    uint32
	offset  uint32
}

type windowsIcoContainer struct {
	Header windowsIcoDescriptor
	Data   []byte
}

func generateWindowsIcon(pngFilename string, iconfile string) error {
	sizes := []int{256, 128, 64, 48, 32, 16}

	pngfile, err := os.Open(pngFilename)
	if err != nil {
		return err
	}
	defer pngfile.Close()

	pngdata, err := png.Decode(pngfile)
	if err != nil {
		return err
	}

	icons := []windowsIcoContainer{}

	for _, size := range sizes {
		rect := image.Rect(0, 0, int(size), int(size))
		rawdata := image.NewRGBA(rect)
		scale := draw.CatmullRom
		scale.Scale(rawdata, rect, pngdata, pngdata.Bounds(), draw.Over, nil)

		icondata := new(bytes.Buffer)
		writer := bufio.NewWriter(icondata)
		err = png.Encode(writer, rawdata)
		if err != nil {
			return err
		}
		writer.Flush()

		imgSize := size
		if imgSize >= 256 {
			imgSize = 0
		}

		data := icondata.Bytes()

		icn := windowsIcoContainer{
			Header: windowsIcoDescriptor{
				width:  uint8(imgSize),
				height: uint8(imgSize),
				planes: 1,
				bpp:    32,
				size:   uint32(len(data)),
			},
			Data: data,
		}
		icons = append(icons, icn)
	}

	outfile, err := os.Create(iconfile)
	if err != nil {
		return err
	}
	defer outfile.Close()

	ico := windowsIcoHeader{
		imageType:  1,
		imageCount: uint16(len(sizes)),
	}
	err = binary.Write(outfile, binary.LittleEndian, ico)
	if err != nil {
		return err
	}

	offset := uint32(6 + 16*len(sizes))
	for _, icon := range icons {
		icon.Header.offset = offset
		err = binary.Write(outfile, binary.LittleEndian, icon.Header)
		if err != nil {
			return err
		}
		offset += icon.Header.size
	}
	for _, icon := range icons {
		_, err = outfile.Write(icon.Data)
		if err != nil {
			return err
		}
	}
	return nil
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
	return filepath.Join(path.Dir(filename), "packages", b.platform)
}

// Package the application into a platform specific package
func (b *PackageHelper) Package(po *ProjectOptions) error {
	switch b.platform {
	case "darwin":
		return b.packageOSX(po)
	case "windows":
		return b.PackageWindows(po, true)
	case "linux":
		return b.packageLinux(po)
	default:
		return fmt.Errorf("platform '%s' not supported for bundling yet", b.platform)
	}
}

func (b *PackageHelper) packageLinux(po *ProjectOptions) error {
	return nil
}

// Package the application for OSX
func (b *PackageHelper) packageOSX(po *ProjectOptions) error {
	build := path.Join(b.fs.Cwd(), "build")

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
	plistFilename := path.Join(build, appname, "Contents", "Info.plist")
	customPlist := path.Join(b.fs.Cwd(), "info.plist")

	// Check binary exists
	source := path.Join(build, exe)
	if po.CrossCompile == true {
		file, err := b.fs.FindFile(build, "darwin")
		if err != nil {
			return err
		}
		source = path.Join(build, file)
	}

	if !b.fs.FileExists(source) {
		// We need to build!
		return fmt.Errorf("Target '%s' not available. Has it been compiled yet?", source)
	}
	// Remove the existing package
	os.RemoveAll(appname)

	// Create directories
	exeDir := path.Join(build, appname, "/Contents/MacOS")
	b.fs.MkDirs(exeDir, 0755)
	resourceDir := path.Join(build, appname, "/Contents/Resources")
	b.fs.MkDirs(resourceDir, 0755)

	// Do we have a custom plist in the project directory?
	if !fs.FileExists(customPlist) {

		// No - create a new plist from our defaults
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

		// Save to the package
		err = ioutil.WriteFile(plistFilename, tpl.Bytes(), 0644)
		if err != nil {
			return err
		}

		// Also write to project directory for customisation
		err = ioutil.WriteFile(customPlist, tpl.Bytes(), 0644)
		if err != nil {
			return err
		}
	} else {
		// Yes - we have a plist. Copy it to the package verbatim
		err = fs.CopyFile(customPlist, plistFilename)
		if err != nil {
			return err
		}
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

// CleanWindows removes any windows related files found in the directory
func (b *PackageHelper) CleanWindows(po *ProjectOptions) {
	pdir := b.fs.Cwd()
	basename := strings.TrimSuffix(po.BinaryName, ".exe")
	exts := []string{".ico", ".exe.manifest", ".rc", "-res.syso"}
	rsrcs := []string{}
	for _, ext := range exts {
		rsrcs = append(rsrcs, filepath.Join(pdir, basename+ext))
	}
	b.fs.RemoveFiles(rsrcs, true)
}

// PackageWindows packages the application for windows platforms
func (b *PackageHelper) PackageWindows(po *ProjectOptions, cleanUp bool) error {
	outputDir := b.fs.Cwd()
	basename := strings.TrimSuffix(po.BinaryName, ".exe")

	// Copy default icon if needed
	icon, err := b.copyIcon()
	if err != nil {
		return err
	}

	// Generate icon from PNG
	err = generateWindowsIcon(icon, basename+".ico")
	if err != nil {
		return err
	}

	// Copy manifest
	tgtManifestFile := filepath.Join(outputDir, basename+".exe.manifest")
	if !b.fs.FileExists(tgtManifestFile) {
		srcManifestfile := filepath.Join(b.getPackageFileBaseDir(), "wails.exe.manifest")
		err := b.fs.CopyFile(srcManifestfile, tgtManifestFile)
		if err != nil {
			return err
		}
	}

	// Copy rc file
	tgtRCFile := filepath.Join(outputDir, basename+".rc")
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
	sysofile := filepath.Join(outputDir, basename+"-res.syso")

	// cross-compile
	if b.platform != runtime.GOOS {
		args := []string{
			"docker", "run", "--rm",
			"-v", outputDir + ":/build",
			"--entrypoint", "/bin/sh",
			"wailsapp/xgo:1.16.2",
			"-c", "/usr/bin/x86_64-w64-mingw32-windres -o /build/" + basename + "-res.syso /build/" + basename + ".rc",
		}
		if err := NewProgramHelper().RunCommandArray(args); err != nil {
			return err
		}
	} else {
		batfile, err := fs.LocalDir(".")
		if err != nil {
			return err
		}

		windresBatFile := filepath.Join(batfile.fullPath, "windres.bat")
		windresCommand := []string{windresBatFile, sysofile, tgtRCFile}
		err = NewProgramHelper().RunCommandArray(windresCommand)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *PackageHelper) copyIcon() (string, error) {

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

	srcIcon, err := b.copyIcon()
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
