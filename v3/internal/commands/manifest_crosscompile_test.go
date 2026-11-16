package commands

import (
	"debug/buildinfo"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestManifestExplicitCGoBuildsWithNativeCompiler(t *testing.T) {
	if testing.Short() {
		t.Skip("builds and runs a real CGo executable")
	}
	compiler := "gcc"
	if runtime.GOOS == "darwin" {
		compiler = "clang"
	}
	path, err := exec.LookPath(compiler)
	if err != nil {
		t.Skipf("native C compiler is unavailable: %v", err)
	}
	root := t.TempDir()
	for name, content := range map[string]string{
		"go.mod":  "module example.com/cgobuild\n\ngo 1.26.0\n",
		"main.go": "package main\n/* int answer(void) { return 42; } */\nimport \"C\"\nimport \"fmt\"\nfunc main() { fmt.Println(C.answer()) }\n",
	} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", filepath.Join(runtime.GOROOT(), "bin")+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("GOWORK", "off")
	t.Setenv("CGO_ENABLED", "0")
	t.Setenv("CC", path)
	spec := pipeline.CompileSpec{Package: ".", Console: true, CGOEnabled: "1", TargetOS: runtime.GOOS, TargetArch: runtime.GOARCH, Toolchain: "native", Output: "bin/native"}
	if runtime.GOOS == "windows" {
		spec.Output += ".exe"
	}
	if _, err := (&manifestHandler{root: root}).compile(t.Context(), spec); err != nil {
		t.Fatal(err)
	}
	output, err := exec.CommandContext(t.Context(), filepath.Join(root, spec.Output)).CombinedOutput()
	if err != nil || string(output) != "42\n" {
		t.Fatalf("CGo application: output=%q error=%v", output, err)
	}
}

func TestManifestPureGoCrossCompilesDesktopBinaries(t *testing.T) {
	if testing.Short() {
		t.Skip("builds six real Go executables")
	}
	root := t.TempDir()
	for name, content := range map[string]string{
		"go.mod":  "module example.com/crosscompile\n\ngo 1.26.0\n",
		"main.go": "package main\nimport \"fmt\"\nfunc main() { fmt.Println(\"crosscompile\") }\n",
	} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", filepath.Join(runtime.GOROOT(), "bin")+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("GOWORK", "off")
	t.Setenv("CGO_ENABLED", "1")
	t.Setenv("CC", filepath.Join(root, "no-c-compiler"))
	t.Setenv("CXX", filepath.Join(root, "no-cxx-compiler"))
	t.Setenv("PKG_CONFIG", filepath.Join(root, "no-pkg-config"))
	for _, targetOS := range []string{"linux", "darwin", "windows"} {
		for _, targetArch := range []string{"amd64", "arm64"} {
			t.Run(targetOS+"/"+targetArch, func(t *testing.T) {
				spec := pipeline.CompileSpec{Package: ".", Console: true, CGOEnabled: "0", TargetOS: targetOS, TargetArch: targetArch, Toolchain: "native", Output: "bin/" + targetOS + "-" + targetArch}
				if _, err := (&manifestHandler{root: root}).compile(t.Context(), spec); err != nil {
					t.Fatal(err)
				}
				info, err := buildinfo.ReadFile(filepath.Join(root, spec.Output))
				if err != nil {
					t.Fatal(err)
				}
				settings := map[string]string{}
				for _, setting := range info.Settings {
					settings[setting.Key] = setting.Value
				}
				for key, want := range map[string]string{"GOOS": targetOS, "GOARCH": targetArch, "CGO_ENABLED": "0"} {
					if settings[key] != want {
						t.Errorf("%s=%q, want %q", key, settings[key], want)
					}
				}
			})
		}
	}
}
