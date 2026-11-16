package appimage

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

// RelocateWebKitHelpers points the bundled WebKit library at the helper
// processes copied into the AppDir. WebKitGTK looks for them in a directory
// compiled into the library (only developer builds honour WEBKIT_EXEC_PATH),
// so as deployed the AppImage only starts where the build machine's path
// exists, e.g. /usr/lib/x86_64-linux-gnu on Debian and Ubuntu, and there it
// runs the system's helpers. Everywhere else WebKit fails with "Unable to spawn
// a new child process".
//
// Each helper directory moves to usr/lib/<name> and the library's copy of its
// path becomes "././<name>/", NUL-padded to the original length. WebKit hands
// that string to bubblewrap both as a bind destination, resolved against the
// sandbox root, and as the exec path, resolved against the working directory.
// bwrap can create a top-level destination in its read-only root where a nested
// one under /lib fails, and appRun starts the application in $APPDIR/usr/lib,
// which WebKit binds at the same absolute path, so the exec path resolves inside
// the sandbox too. Only the whole NUL-terminated string is rewritten: the
// injected-bundle directory shares the prefix on Debian-family systems and must
// stay absolute, because a relative bind destination there fails the same way;
// appRun redirects it to the bundled copy.
func RelocateWebKitHelpers(appDir string, helperDirs []string) error {
	libDir := filepath.Join(appDir, "usr", "lib")
	libs, err := filepath.Glob(filepath.Join(libDir, "libwebkit*gtk-*.so*"))
	if err != nil {
		return err
	}
	contents := make([][]byte, len(libs))
	for i, lib := range libs {
		if contents[i], err = os.ReadFile(lib); err != nil {
			return err
		}
	}
	patched := make([][]byte, len(libs))
	copy(patched, contents)
	for _, dir := range helperDirs {
		abs := []byte(dir + "\x00")
		rel := []byte("././" + filepath.Base(dir) + "/")
		if len(rel) >= len(abs) {
			return fmt.Errorf("WebKit helper directory %q is too short to relocate in place", dir)
		}
		rel = append(rel, make([]byte, len(abs)-len(rel))...)
		// Only the directory a library was built with is moved: the file search can also turn up other copies.
		referenced := false
		for i := range patched {
			if bytes.Contains(patched[i], abs) {
				patched[i] = bytes.ReplaceAll(patched[i], abs, rel)
				referenced = true
			}
		}
		if !referenced {
			continue
		}
		if target := filepath.Join(libDir, filepath.Base(dir)); target != filepath.Join(appDir, dir) {
			if err := os.Rename(filepath.Join(appDir, dir), target); err != nil {
				return err
			}
		}
	}
	relocated := false
	for i, lib := range libs {
		if bytes.Equal(patched[i], contents[i]) {
			continue
		}
		info, err := os.Stat(lib)
		if err != nil {
			return err
		}
		if err := os.WriteFile(lib, patched[i], info.Mode().Perm()); err != nil {
			return err
		}
		relocated = true
	}
	if !relocated {
		return fmt.Errorf("no WebKit library in %s references the helper directories %v", libDir, helperDirs)
	}
	return nil
}
