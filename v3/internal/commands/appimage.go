package commands

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/wailsapp/wails/v3/internal/term"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/s"
)

//go:embed linuxdeploy-plugin-gtk.sh
var gtkPlugin []byte

func log(p *pterm.ProgressbarPrinter, message string) {
	p.UpdateTitle(message)
	term.Infof(message)
	p.Increment()
}

type GenerateAppImageOptions struct {
	Binary      string `description:"The binary to package including path"`
	Icon        string `description:"Path to the icon"`
	DesktopFile string `description:"Path to the desktop file"`
	OutputDir   string `description:"Path to the output directory" default:"."`
	BuildDir    string `description:"Path to the build directory"`
}

func GenerateAppImage(options *GenerateAppImageOptions) error {
	DisableFooter = true

	defer func() {
		_ = pterm.DefaultSpinner.Stop()
	}()

	if options.Binary == "" {
		return fmt.Errorf("binary not provided")
	}
	if options.Icon == "" {
		return fmt.Errorf("icon path not provided")
	}
	if options.DesktopFile == "" {
		return fmt.Errorf("desktop file path not provided")
	}
	if options.BuildDir == "" {
		// Create temp directory
		var err error
		options.BuildDir, err = os.MkdirTemp("", "wails-appimage-*")
		if err != nil {
			return err
		}
	}
	// Resolve every input path to absolute form up-front. The bundler does
	// `s.CD` into the build directory partway through, so anything left as
	// a relative path would be interpreted relative to the wrong CWD by
	// downstream goroutines and shell-outs (e.g. `ldd <binary>`).
	for _, p := range []*string{&options.OutputDir, &options.BuildDir, &options.Binary, &options.Icon, &options.DesktopFile} {
		abs, err := filepath.Abs(*p)
		if err != nil {
			return err
		}
		*p = abs
	}

	term.Header("AppImage Generator")

	return generateAppImage(options)
}

func generateAppImage(options *GenerateAppImageOptions) error {
	numberOfSteps := 5
	p, _ := pterm.DefaultProgressbar.WithTotal(numberOfSteps).WithTitle("Generating AppImage").Start()

	// Get the last path of the binary and normalise the name
	name := normaliseName(filepath.Base(options.Binary))

	// Architecture-specific variables using a map
	archDetails := map[string]string{
		"arm64":  "aarch64",
		"amd64":  "x86_64",
		"x86_64": "x86_64",
	}

	arch, exists := archDetails[runtime.GOARCH]
	if !exists {
		return fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
	}

	appDir := filepath.Join(options.BuildDir, fmt.Sprintf("%s-%s.AppDir", name, arch))
	s.RMDIR(appDir)

	log(p, "Preparing AppImage Directory: "+appDir)

	usrBin := filepath.Join(appDir, "usr", "bin")
	s.MKDIR(options.BuildDir)
	s.MKDIR(usrBin)
	s.COPY(options.Binary, usrBin)
	s.CHMOD(filepath.Join(usrBin, filepath.Base(options.Binary)), 0755)
	dotDirIcon := filepath.Join(appDir, ".DirIcon")
	s.COPY(options.Icon, dotDirIcon)
	iconLink := filepath.Join(appDir, filepath.Base(options.Icon))
	s.DELETE(iconLink)
	s.SYMLINK(".DirIcon", iconLink)
	s.COPY(options.DesktopFile, appDir)

	// Download linuxdeploy and make it executable
	s.CD(options.BuildDir)

	// Download URLs using a map based on architecture
	urls := map[string]string{
		"linuxdeploy": fmt.Sprintf("https://github.com/linuxdeploy/linuxdeploy/releases/download/continuous/linuxdeploy-%s.AppImage", arch),
		"AppRun":      fmt.Sprintf("https://github.com/AppImage/AppImageKit/releases/download/continuous/AppRun-%s", arch),
	}

	// Download necessary files concurrently
	log(p, "Downloading AppImage tooling")
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		linuxdeployPath := filepath.Join(options.BuildDir, filepath.Base(urls["linuxdeploy"]))
		if !s.EXISTS(linuxdeployPath) {
			s.DOWNLOAD(urls["linuxdeploy"], linuxdeployPath)
		}
		s.CHMOD(linuxdeployPath, 0755)
		wg.Done()
	}()

	go func() {
		target := filepath.Join(appDir, "AppRun")
		if !s.EXISTS(target) {
			s.DOWNLOAD(urls["AppRun"], target)
		}
		s.CHMOD(target, 0755)
		wg.Done()
	}()

	wg.Wait()

	// Determine which GTK stack the binary links against by ldd'ing the
	// source binary. We need this before searching for runtime files,
	// because GTK3 (WebKit2GTK 4.x) and GTK4 (WebKitGTK 6.0) ship the
	// injected-bundle library under different names. The %q quoting keeps
	// `s.EXEC`'s shlex split intact when the binary path contains spaces.
	lddOutput, err := s.EXEC(fmt.Sprintf("ldd %q", options.Binary))
	if err != nil {
		println(string(lddOutput))
		return err
	}
	lddString := string(lddOutput)
	var DeployGtkVersion string
	switch {
	case s.CONTAINS(lddString, "libgtk-4.so"):
		DeployGtkVersion = "4"
	case s.CONTAINS(lddString, "libgtk-3.so"):
		DeployGtkVersion = "3"
	case s.CONTAINS(lddString, "libgtk-x11-2.0.so"):
		DeployGtkVersion = "2"
	default:
		snippet := lddString
		if len(snippet) > 200 {
			snippet = snippet[:200] + "..."
		}
		return fmt.Errorf("unable to determine GTK version (looked for libgtk-4.so, libgtk-3.so, libgtk-x11-2.0.so in ldd %q output): %s", options.Binary, snippet)
	}

	// Processing GTK files
	log(p, "Processing GTK files.")
	injectedBundle := "libwebkit2gtkinjectedbundle.so"
	if DeployGtkVersion == "4" {
		injectedBundle = "libwebkitgtkinjectedbundle.so"
	}
	filesNeeded := []string{"WebKitWebProcess", "WebKitNetworkProcess", injectedBundle}
	files, err := findGTKFiles(filesNeeded)
	if err != nil {
		return err
	}
	s.CD(appDir)
	for _, file := range files {
		targetDir := filepath.Dir(file)
		if targetDir[0] == '/' {
			targetDir = targetDir[1:]
		}
		targetDir, err = filepath.Abs(targetDir)
		if err != nil {
			return err
		}
		s.MKDIR(targetDir)
		s.COPY(file, targetDir)
	}

	// Copy GTK Plugin
	err = os.WriteFile(filepath.Join(options.BuildDir, "linuxdeploy-plugin-gtk.sh"), gtkPlugin, 0755)
	if err != nil {
		return err
	}

	// Run linuxdeploy to bundle the application
	s.CD(options.BuildDir)
	linuxdeployAppImage := filepath.Join(options.BuildDir, fmt.Sprintf("linuxdeploy-%s.AppImage", arch))

	// Quote the executable and --appdir args so `s.EXEC`'s shlex split
	// keeps them as single tokens when the user-supplied paths contain
	// spaces.
	cmd := fmt.Sprintf("%q --appimage-extract-and-run --appdir %q --output appimage --plugin gtk", linuxdeployAppImage, appDir)
	s.SETENV("DEPLOY_GTK_VERSION", DeployGtkVersion)

	// Force linuxdeploy's appimage plugin to write the AppImage to a known
	// filename. Without this it derives the name from the desktop file's
	// `Name=` field, which often doesn't match the binary basename and
	// causes the subsequent move-to-output step to fail.
	appImageName := fmt.Sprintf("%s-%s.AppImage", name, arch)
	targetFile := filepath.Join(options.BuildDir, appImageName)
	s.SETENV("OUTPUT", appImageName)

	// Check if system libraries use .relr.dyn sections (modern toolchains)
	// If so, disable stripping as linuxdeploy's bundled strip can't handle them
	if hasRelrDynSections() {
		term.Infof("Detected modern toolchain (.relr.dyn sections), disabling stripping for compatibility. See: https://v3.wails.io/guides/build/linux#appimage-strip-compatibility")
		s.SETENV("NO_STRIP", "1")
	}

	output, err := s.EXEC(cmd)
	if err != nil {
		fmt.Println(string(output))
		return err
	}

	// Move file to output directory
	s.MOVE(targetFile, options.OutputDir)

	log(p, "AppImage created: "+filepath.Join(options.OutputDir, appImageName))
	return nil
}

func findGTKFiles(files []string) ([]string, error) {
	notFound := []string{}
	found := []string{}
	err := filepath.Walk("/usr/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				return nil
			}
			return err
		}

		if info.IsDir() {
			return nil
		}

		for _, fileName := range files {
			if strings.HasSuffix(path, fileName) {
				found = append(found, path)
				break
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	for _, fileName := range files {
		fileFound := false
		for _, foundPath := range found {
			if strings.HasSuffix(foundPath, fileName) {
				fileFound = true
				break
			}
		}
		if !fileFound {
			notFound = append(notFound, fileName)
		}
	}
	if len(notFound) > 0 {
		return nil, errors.New("Unable to locate all required files: " + strings.Join(notFound, ", "))
	}
	return found, nil
}

// hasRelrDynSections checks if system libraries use .relr.dyn sections
// which are incompatible with linuxdeploy's bundled strip binary.
// This is common on modern Linux distributions (Arch, Fedora 39+, Ubuntu 24.04+).
func hasRelrDynSections() bool {
	// Check common GTK libraries that will be bundled. We probe both the
	// GTK4 (default since v3.0.0-alpha.93) and GTK3 (legacy `-tags gtk3`)
	// libraries because either may be present depending on the build.
	testLibs := []string{
		// GTK4
		"/usr/lib/libgtk-4.so.1",
		"/usr/lib64/libgtk-4.so.1",
		"/usr/lib/x86_64-linux-gnu/libgtk-4.so.1",
		"/usr/lib/aarch64-linux-gnu/libgtk-4.so.1",
		// GTK3
		"/usr/lib/libgtk-3.so.0",
		"/usr/lib64/libgtk-3.so.0",
		"/usr/lib/x86_64-linux-gnu/libgtk-3.so.0",
		"/usr/lib/aarch64-linux-gnu/libgtk-3.so.0",
	}

	for _, lib := range testLibs {
		if _, err := os.Stat(lib); err == nil {
			output, err := s.EXEC(fmt.Sprintf("readelf -S %q", lib))
			if err == nil && strings.Contains(string(output), ".relr.dyn") {
				return true
			}
		}
	}
	return false
}
