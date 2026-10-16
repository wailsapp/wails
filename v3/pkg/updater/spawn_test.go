package updater

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAppImageOr(t *testing.T) {
	img := filepath.Join(t.TempDir(), "App.AppImage")
	writeFile(t, img, []byte("appimage"))
	mount := t.TempDir()
	inMount := filepath.Join(mount, "usr", "bin", "app")
	outside := filepath.Join(t.TempDir(), "other")

	for _, tc := range []struct {
		name, image, appdir, exe, want string
	}{
		{"not an AppImage", "", "", inMount, inMount},
		{"running from the mount", img, mount, inMount, img},
		{"inherited by another program", img, mount, outside, outside},
		{"mount name is only a prefix", img, mount, mount + "-x/app", mount + "-x/app"},
		{"AppImage missing", filepath.Join(t.TempDir(), "gone"), mount, inMount, inMount},
		{"AppImage is a directory", mount, mount, inMount, inMount},
		{"APPDIR unset", img, "", inMount, inMount},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("APPIMAGE", tc.image)
			t.Setenv("APPDIR", tc.appdir)
			if got := appImageOr(tc.exe); got != tc.want {
				t.Errorf("appImageOr(%q) = %q, want %q", tc.exe, got, tc.want)
			}
		})
	}
}

func TestSelfExecutable_AppImage(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	img := filepath.Join(t.TempDir(), "App.AppImage")
	writeFile(t, img, []byte("appimage"))
	t.Setenv("APPIMAGE", img)
	t.Setenv("APPDIR", filepath.Dir(exe))
	if got, err := resolveTarget(); err != nil || got != img {
		t.Errorf("resolveTarget() = %q, %v; want %q", got, err, img)
	}
}
