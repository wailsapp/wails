package appimage

import (
	"context"
	"debug/elf"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/wailsapp/wails/v3/internal/stagefile"
)

func (options Options) prepareNativeDirectory(ctx context.Context, stage string) (string, error) {
	appdir := filepath.Join(stage, "AppDir")
	if err := os.MkdirAll(filepath.Join(appdir, "usr/lib"), 0755); err != nil {
		return "", err
	}
	buffer := make([]byte, 128<<10)
	for _, binary := range options.Binaries {
		source := options.path(binary.Path)
		image, err := elf.Open(source)
		if err != nil {
			return "", err
		}
		matches := image.Machine == machineForArch(options.Arch)
		image.Close()
		if !matches {
			return "", fmt.Errorf("executable %q does not match %s", binary.Name, options.Arch)
		}
		destination := filepath.Join(appdir, "usr/bin", binary.Name)
		if err := stagefile.Copy(ctx, source, destination, 0755, buffer); err != nil {
			return "", err
		}
		if err := preferBundledLibraries(destination); err != nil {
			return "", err
		}
	}
	deployment, err := options.collectRuntime(ctx)
	if err != nil {
		return "", err
	}
	if err := deployment.stage(ctx, appdir); err != nil {
		return "", err
	}
	if options.WebKit {
		if err := RelocateWebKitHelpers(appdir, deployment.helpers); err != nil {
			return "", err
		}
	}
	launcher := nativeLauncher(options.Binaries[0].Name)
	if options.Tools.AppRun != "" {
		if err := stagefile.Copy(ctx, options.path(options.Tools.AppRun), filepath.Join(appdir, "AppRun"), 0755, buffer); err != nil {
			return "", err
		}
	} else if err := os.WriteFile(filepath.Join(appdir, "AppRun"), launcher, 0755); err != nil {
		return "", err
	}
	iconName := options.IconName
	if iconName == "" {
		iconName = filepath.Base(options.Icon)
	}
	for _, name := range []string{".DirIcon", iconName} {
		if err := stagefile.Copy(ctx, options.path(options.Icon), filepath.Join(appdir, name), 0644, buffer); err != nil {
			return "", err
		}
	}
	if err := stagefile.Copy(ctx, options.path(options.DesktopFile), filepath.Join(appdir, filepath.Base(options.DesktopFile)), 0644, buffer); err != nil {
		return "", err
	}
	return appdir, nil
}

func nativeLauncher(binary string) []byte {
	environment := `export GTK_DATA_PREFIX="$APPDIR"
export GTK_EXE_PREFIX="$usr"
export GI_TYPELIB_PATH="$usr/lib/girepository-1.0"
export GIO_MODULE_DIR="$usr/lib/gio/modules"
export GDK_BACKEND="${GDK_BACKEND:-x11}"
for gtk in "$usr"/lib/gtk-*.0; do
 if [ -d "$gtk" ]; then export GTK_PATH="$gtk"; break; fi
done
for cache in "$usr"/lib/gdk-pixbuf-2.0/*/loaders.cache; do
 if [ -f "$cache" ]; then export GDK_PIXBUF_MODULE_FILE="$cache"; break; fi
done
for cache in "$usr"/lib/gtk-3.0/*/immodules.cache; do
 if [ -f "$cache" ]; then export GTK_IM_MODULE_FILE="$cache"; break; fi
done
`
	return []byte(strings.Replace(string(Launcher(binary)), "cd \"$usr/lib\"\n", environment+"cd \"$usr/lib\"\n", 1))
}

func (options Options) collectRuntime(ctx context.Context) (*runtimeDeployment, error) {
	var inputs []string
	for _, binary := range options.Binaries {
		inputs = append(inputs, options.path(binary.Path))
	}
	root := options.RuntimeRoot
	if root == "" {
		if runtime.GOOS == "linux" {
			root = "/"
		} else if options.WebKit {
			return nil, fmt.Errorf("GUI AppImage creation on %s requires a Linux runtime root", runtime.GOOS)
		} else {
			root = options.Root
		}
	} else {
		root = options.path(root)
	}
	deployment, err := newRuntimeDeployment(root, options.Arch, options.GTKVersion, options.LibraryDirs)
	if err != nil {
		return nil, err
	}
	if err := deployment.dependencies(ctx, inputs); err != nil {
		return nil, err
	}
	if options.WebKit {
		if deployment.gtk == 0 {
			deployment.gtk = 4
		}
		if err := deployment.profile(ctx); err != nil {
			return nil, err
		}
	}
	return deployment, nil
}

// RuntimeInputs lists the resolved target libraries and data that influence
// packaging, so an OS or supplied runtime update invalidates the package cache.
func (options Options) RuntimeInputs(ctx context.Context) ([]string, error) {
	deployment, err := options.collectRuntime(ctx)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	var result []string
	for _, source := range deployment.files {
		if !seen[source] {
			result = append(result, source)
			seen[source] = true
		}
	}
	slices.Sort(result)
	return result, nil
}
