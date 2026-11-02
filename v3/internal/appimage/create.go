// Package appimage stages Linux application runtimes and publishes type-2 AppImages.
package appimage

import (
	"context"
	"debug/elf"
	"fmt"
	"io"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/internal/stagefile"
)

type Binary struct{ Name, Path string }
type Command struct {
	Program   string
	Args      []string
	Directory string
	Env       map[string]string
}
type Runner func(context.Context, Command) error

type Tools struct{ LinuxDeploy, AppRun string }
type Options struct {
	Root, Name, Arch, Output, Icon, IconName, DesktopFile, WorkDir, CacheDir string
	Binaries                                                                 []Binary
	LibraryDirs                                                              []string
	RuntimeRoot                                                              string
	RuntimePath                                                              string
	GTKVersion                                                               int
	// UpdateInfo is passed to linuxdeploy's AppImage plugin; it must emit a zsync sidecar.
	UpdateInfo string
	WebKit     bool
	// GTKPlugin is the existing external plugin artifact, not a native replacement.
	GTKPlugin  []byte
	Tools      Tools
	ForceTools bool
	Env        map[string]string
	Runner     Runner
	Client     *http.Client
	Progress   func(string)
}

func OutputName(name, arch string) (string, error) {
	if name == "" || name == "." || name == ".." || strings.ContainsAny(name, "/\\\x00") {
		return "", fmt.Errorf("invalid AppImage name %q", name)
	}
	label, err := architectureLabel(arch)
	if err != nil {
		return "", err
	}
	return strings.ToLower(strings.ReplaceAll(name, " ", "-")) + "-" + label + ".AppImage", nil
}
func architectureLabel(arch string) (string, error) {
	switch arch {
	case "amd64":
		return "x86_64", nil
	case "arm64":
		return "aarch64", nil
	}
	return "", fmt.Errorf("unsupported AppImage architecture %q", arch)
}

func Create(ctx context.Context, options Options) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	var err error
	options, err = normalizeOptions(options)
	if err != nil {
		return err
	}
	var inputs []string
	for _, binary := range options.Binaries {
		inputs = append(inputs, binary.Path)
	}
	var analysis elfAnalysis
	err = options.runStage("Inspect ELF dependencies", func() error {
		var inspectErr error
		analysis, inspectErr = analyzeELF(ctx, inputs, options.libraryDirs())
		return inspectErr
	})
	if err != nil {
		return err
	}
	if options.WebKit && analysis.gtk == "" {
		return fmt.Errorf("WebKit AppImage has no GTK dependency")
	}
	if analysis.gtk != "" && len(options.GTKPlugin) == 0 {
		return fmt.Errorf("GTK AppImage requires the GTK deployment plugin")
	}
	if err = os.MkdirAll(filepath.Dir(options.Output), 0755); err != nil {
		return err
	}
	base := filepath.Dir(options.Output)
	if options.WorkDir != "" {
		base = options.path(options.WorkDir)
	}
	if err = os.MkdirAll(base, 0755); err != nil {
		return err
	}
	stage, err := os.MkdirTemp(base, ".wails-appimage-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(stage)
	var tools Tools
	err = options.runStage("Prepare AppImage tools", func() error {
		var toolErr error
		tools, toolErr = options.PrepareTools(ctx)
		return toolErr
	})
	if err != nil {
		return err
	}
	var appdir string
	var helperDirs []string
	err = options.runStage("Stage AppImage files", func() error {
		var stageErr error
		appdir, helperDirs, stageErr = options.prepareDirectory(ctx, stage, tools.AppRun, analysis.gtk)
		return stageErr
	})
	if err != nil {
		return err
	}
	if err = options.runStage("Deploy dependencies and compress AppImage", func() error {
		return options.runLinuxdeploy(ctx, stage, tools.LinuxDeploy, appdir, analysis, helperDirs)
	}); err != nil {
		return err
	}
	return options.runStage("Verify and publish AppImage", func() error {
		result := filepath.Join(stage, options.resultName())
		if err := verify(result, options.Arch); err != nil {
			return err
		}
		if err := os.Chmod(result, 0755); err != nil {
			return err
		}
		return options.publishResults(ctx, result)
	})
}

func (options Options) runStage(name string, run func() error) error {
	start := time.Now()
	if options.Progress != nil {
		options.Progress(name)
	}
	err := run()
	status := "completed"
	if err != nil {
		status = "failed"
	}
	if options.Progress != nil {
		options.Progress(fmt.Sprintf("%s %s in %s", name, status, time.Since(start).Round(time.Millisecond)))
	}
	if err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	return nil
}

func normalizeOptions(options Options) (Options, error) {
	if runtime.GOOS != "linux" {
		return options, fmt.Errorf("AppImage tooling requires a Linux host")
	}
	root, err := filepath.Abs(options.Root)
	if err != nil {
		return options, err
	}
	options.Root = root
	if options.Arch == "" {
		options.Arch = runtime.GOARCH
	}
	if options.Arch != runtime.GOARCH {
		return options, fmt.Errorf("AppImage tooling for %s requires a matching Linux host or container", options.Arch)
	}
	if options.Runner == nil || len(options.Binaries) == 0 || options.Icon == "" || options.DesktopFile == "" {
		return options, fmt.Errorf("AppImage needs a runner, binaries, icon and desktop file")
	}
	if options.Name == "" {
		options.Name = options.Binaries[0].Name
	}
	options, err = normalizeOutputOptions(options)
	if err != nil {
		return options, err
	}
	options.Binaries = slices.Clone(options.Binaries)
	seen := map[string]bool{}
	for i, binary := range options.Binaries {
		if binary.Name == "" || binary.Name == "." || binary.Name == ".." || strings.ContainsAny(binary.Name, "/\\\x00") || seen[binary.Name] {
			return options, fmt.Errorf("invalid or duplicate AppImage executable %q", binary.Name)
		}
		seen[binary.Name] = true
		options.Binaries[i].Path = options.path(binary.Path)
		if slices.Contains(options.outputPaths(), options.Binaries[i].Path) {
			return options, fmt.Errorf("AppImage output conflicts with executable %q", binary.Name)
		}
		image, err := elf.Open(options.Binaries[i].Path)
		if err != nil {
			return options, err
		}
		matches := machineForArch(options.Arch) == image.Machine
		image.Close()
		if !matches {
			return options, fmt.Errorf("executable %q does not match AppImage architecture %s", binary.Name, options.Arch)
		}
	}
	return options, nil
}

func normalizeOutputOptions(options Options) (Options, error) {
	name, err := OutputName(options.Name, options.Arch)
	if err != nil {
		return options, err
	}
	if options.Output == "" {
		options.Output = name
	}
	options.Output = options.path(options.Output)
	if strings.ContainsAny(options.UpdateInfo, "\x00\r\n") {
		return options, fmt.Errorf("AppImage updateInfo must be a single line")
	}
	return options, nil
}

func (options Options) runLinuxdeploy(ctx context.Context, stage, tool, appdir string, analysis elfAnalysis, helperDirs []string) error {
	env := maps.Clone(options.Env)
	if env == nil {
		env = map[string]string{}
	}
	env["OUTPUT"] = options.resultName()
	env["LDAI_OUTPUT"] = env["OUTPUT"]
	if options.UpdateInfo != "" {
		env["LDAI_UPDATE_INFORMATION"] = options.UpdateInfo
	}
	if analysis.relr {
		env["NO_STRIP"] = "1"
	}
	args := []string{"--appimage-extract-and-run", "--appdir", appdir}
	if !options.WebKit {
		args = append(args, "--output", "appimage")
	}
	if analysis.gtk != "" {
		plugin := filepath.Join(stage, "linuxdeploy-plugin-gtk.sh")
		if err := os.WriteFile(plugin, options.GTKPlugin, 0755); err != nil {
			return err
		}
		env["DEPLOY_GTK_VERSION"] = analysis.gtk
		// linuxdeploy discovers plugins beside its binary or on PATH, not in cwd.
		env["PATH"] = stage + string(os.PathListSeparator) + options.environment("PATH")
		args = append(args, "--plugin", "gtk")
	}
	if err := options.Runner(ctx, Command{Program: tool, Args: args, Directory: stage, Env: env}); err != nil {
		return err
	}
	if !options.WebKit {
		return nil
	}
	if err := RelocateWebKitHelpers(appdir, helperDirs); err != nil {
		return err
	}
	return options.Runner(ctx, Command{Program: tool, Args: []string{"--appimage-extract-and-run", "--appdir", appdir, "--output", "appimage"}, Directory: stage, Env: env})
}

func (options Options) libraryDirs() []string {
	dirs := slices.Clone(options.LibraryDirs)
	for i, dir := range dirs {
		dirs[i] = options.path(dir)
	}
	dirs = append(dirs, filepath.SplitList(options.environment("LD_LIBRARY_PATH"))...)
	return append(dirs, defaultLibraryDirs(options.Arch)...)
}

func (options Options) path(value string) string {
	if filepath.IsAbs(value) {
		return filepath.Clean(value)
	}
	return filepath.Join(options.Root, value)
}
func (options Options) environment(key string) string {
	if value, ok := options.Env[key]; ok {
		return value
	}
	return os.Getenv(key)
}
func machineForArch(arch string) elf.Machine {
	if arch == "arm64" {
		return elf.EM_AARCH64
	}
	return elf.EM_X86_64
}

// PrepareTools resolves and downloads the external tools before action identity is calculated.
func (options Options) PrepareTools(ctx context.Context) (Tools, error) {
	tools := options.Tools
	label, err := architectureLabel(options.Arch)
	if err != nil {
		return Tools{}, err
	}
	cache := options.CacheDir
	if cache == "" {
		cache = filepath.Join(".wake", "appimage-tools")
	}
	cache = filepath.Join(options.path(cache), options.Arch)
	if tools.AppRun != "" {
		tools.AppRun = options.path(tools.AppRun)
	}
	if tools.LinuxDeploy != "" {
		tools.LinuxDeploy = options.path(tools.LinuxDeploy)
		return tools, nil
	}
	tools.LinuxDeploy = filepath.Join(cache, "linuxdeploy.AppImage")
	url := "https://github.com/linuxdeploy/linuxdeploy/releases/download/continuous/linuxdeploy-" + label + ".AppImage"
	return tools, ensureTool(ctx, options.Client, url, tools.LinuxDeploy, options.ForceTools)
}

func (options Options) prepareDirectory(ctx context.Context, stage, appRun, gtk string) (string, []string, error) {
	buffer := make([]byte, 32<<10)
	appdir := filepath.Join(stage, "AppDir")
	if err := os.MkdirAll(appdir, 0755); err != nil {
		return "", nil, err
	}
	for _, binary := range options.Binaries {
		if err := stagefile.Copy(ctx, binary.Path, filepath.Join(appdir, "usr", "bin", binary.Name), 0755, buffer); err != nil {
			return "", nil, err
		}
	}
	if err := os.MkdirAll(filepath.Join(appdir, "usr", "lib"), 0755); err != nil {
		return "", nil, err
	}
	launcherPath := filepath.Join(appdir, "AppRun")
	if appRun == "" {
		if err := os.WriteFile(launcherPath, Launcher(options.Binaries[0].Name), 0755); err != nil {
			return "", nil, err
		}
	} else if err := stagefile.Copy(ctx, appRun, launcherPath, 0755, buffer); err != nil {
		return "", nil, err
	}
	if err := stagefile.Copy(ctx, options.path(options.Icon), filepath.Join(appdir, ".DirIcon"), 0644, buffer); err != nil {
		return "", nil, err
	}
	iconName := options.IconName
	if iconName == "" {
		iconName = filepath.Base(options.Icon)
	}
	if iconName != filepath.Base(iconName) || iconName == "." || iconName == ".." {
		return "", nil, fmt.Errorf("invalid bundled icon name %q", iconName)
	}
	if err := os.Symlink(".DirIcon", filepath.Join(appdir, iconName)); err != nil {
		return "", nil, err
	}
	if err := stagefile.Copy(ctx, options.path(options.DesktopFile), filepath.Join(appdir, filepath.Base(options.DesktopFile)), 0644, buffer); err != nil {
		return "", nil, err
	}
	if options.WebKit {
		helperDirs, err := options.copyWebKitFiles(ctx, appdir, gtk, buffer)
		return appdir, helperDirs, err
	}
	return appdir, nil, nil
}

func (options Options) copyWebKitFiles(ctx context.Context, appdir, gtk string, buffer []byte) ([]string, error) {
	injected := "libwebkit2gtkinjectedbundle.so"
	if gtk == "4" {
		injected = "libwebkitgtkinjectedbundle.so"
	}
	root := options.RuntimeRoot
	if root == "" {
		root = "/usr"
	}
	root = options.path(root)
	var files []string
	err := options.runStage("Locate GTK-compatible WebKit runtime files", func() error {
		var findErr error
		files, findErr = FindRuntimeFiles(ctx, root, []string{"WebKitWebProcess", "WebKitNetworkProcess", injected})
		if findErr != nil {
			return findErr
		}
		files, findErr = matchingRuntimeFiles(ctx, files, gtk, options.libraryDirs())
		return findErr
	})
	if err != nil {
		return nil, err
	}
	var helperDirs []string
	for _, file := range files {
		if strings.HasPrefix(filepath.Base(file), "WebKit") && !slices.Contains(helperDirs, filepath.Dir(file)) {
			helperDirs = append(helperDirs, filepath.Dir(file))
		}
		if filepath.Base(file) == "WebKitWebProcess" {
			gpu := filepath.Join(filepath.Dir(file), "WebKitGPUProcess")
			if _, err := os.Stat(gpu); err == nil {
				files = append(files, gpu)
			} else if !os.IsNotExist(err) {
				return nil, err
			}
		}
	}
	for _, source := range files {
		destination := filepath.Join(appdir, strings.TrimLeft(source, "/"))
		if err = stagefile.Copy(ctx, source, destination, 0, buffer); err != nil {
			return nil, err
		}
	}
	return helperDirs, nil
}

func verify(filename, arch string) error {
	input, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer input.Close()
	var header [11]byte
	if _, err = io.ReadFull(input, header[:]); err != nil {
		return err
	}
	if string(header[:4]) != "\x7fELF" || string(header[8:11]) != "AI\x02" {
		return fmt.Errorf("invalid type-2 AppImage")
	}
	image, err := elf.Open(filename)
	if err != nil {
		return err
	}
	defer image.Close()
	if image.Machine != machineForArch(arch) {
		return fmt.Errorf("AppImage output architecture does not match %s", arch)
	}
	return nil
}
