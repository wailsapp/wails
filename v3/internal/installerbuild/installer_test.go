package installerbuild

import (
	"context"
	"debug/pe"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/buildmeta"
)

func testOptions(t *testing.T) Options {
	t.Helper()
	root := filepath.Join(t.TempDir(), "path with spaces $cash")
	if err := os.MkdirAll(root, 0755); err != nil {
		t.Fatal(err)
	}
	options := Options{Script: filepath.Join(root, "project.nsi"), Output: filepath.Join(root, "dist", "installer.exe"), Arch: "amd64", Scope: "user", Primary: File{Name: "app.exe", Path: filepath.Join(root, "app.exe")}, Companions: []File{{Name: "helper.exe", Path: filepath.Join(root, "helper.exe")}}}
	for path, data := range map[string]string{
		options.Script: `Unicode true
Name "Fixture"
OutFile "unused.exe"
RequestExecutionLevel ${REQUEST_EXECUTION_LEVEL}
InstallDir "$LOCALAPPDATA\Fixture"
!ifndef INFO_PRODUCTVERSION
!define INFO_PRODUCTVERSION "1.0.0"
!endif
VIProductVersion "${INFO_PRODUCTVERSION}.0"
!include "relative.nsh"
Section "Main"
SetOutPath "$INSTDIR"
File "/oname=${PRODUCT_EXECUTABLE}" "${ARG_WAILS_AMD64_BINARY}"
SectionEnd
`,
		filepath.Join(root, "relative.nsh"): `BrandingText "Relative include retained"`,
		options.Primary.Path:                "primary payload", options.Companions[0].Path: "companion payload",
	} {
		if err := os.WriteFile(path, []byte(data), 0644); err != nil {
			t.Fatal(err)
		}
	}
	return options
}

func TestInstallerPublishesOnlySuccessfulCompilerOutput(t *testing.T) {
	options := testOptions(t)
	if err := os.MkdirAll(filepath.Dir(options.Output), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(options.Output, []byte("previous installer"), 0644); err != nil {
		t.Fatal(err)
	}
	for _, outcome := range []string{"error", "cancelled", "missing", "empty", "directory", "success"} {
		t.Run(outcome, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var stage string
			err := Create(ctx, options, func(ctx context.Context, directory string, args []string) error {
				if directory != filepath.Dir(options.Script) || !reflect.DeepEqual(args[:3], []string{"-NOCD", "-INPUTCHARSET", "UTF8"}) {
					t.Fatalf("compiler invocation: %s %v", directory, args)
				}
				stage = filepath.Dir(args[len(args)-1])
				output := filepath.Join(stage, "installer.exe")
				switch outcome {
				case "missing":
					return nil
				case "directory":
					return os.Mkdir(output, 0755)
				case "empty":
					return os.WriteFile(output, nil, 0644)
				}
				if err := os.WriteFile(output, []byte("new installer"), 0644); err != nil {
					return err
				}
				if outcome == "error" {
					return errors.New("compiler failure")
				}
				if outcome == "cancelled" {
					cancel()
				}
				return nil
			})
			if (err == nil) != (outcome == "success") {
				t.Fatalf("outcome=%s error=%v", outcome, err)
			}
			want := "previous installer"
			if outcome == "success" {
				want = "new installer"
			}
			data, readErr := os.ReadFile(options.Output)
			if readErr != nil || string(data) != want {
				t.Fatalf("published %q: %v", data, readErr)
			}
			if _, err := os.Stat(stage); !os.IsNotExist(err) {
				t.Fatalf("stage remains: %s, %v", stage, err)
			}
		})
	}
}

func TestInstallerWithRealCompiler(t *testing.T) {
	compiler := os.Getenv("WAILS_TEST_MAKENSIS")
	if compiler == "" {
		var err error
		compiler, err = exec.LookPath("makensis")
		if err != nil {
			t.Skip("makensis not installed")
		}
	}
	if runtime.GOOS == "darwin" {
		t.Setenv("LC_ALL", "C")
	}
	options := testOptions(t)
	version, err := buildmeta.Parse("1.2.3-beta.4+build.5")
	if err != nil {
		t.Fatal(err)
	}
	options.Version = &version
	source, err := os.ReadFile(options.Script)
	if err != nil {
		t.Fatal(err)
	}
	err = Create(context.Background(), options, func(ctx context.Context, directory string, args []string) error {
		cmd := exec.CommandContext(ctx, compiler, args...)
		cmd.Dir = directory
		cmd.Env = CompilerEnvironment(runtime.GOOS, os.Environ())
		output, err := cmd.CombinedOutput()
		if err != nil {
			return errors.New(string(output) + ": " + err.Error())
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	image, err := pe.Open(options.Output)
	if err != nil {
		t.Fatal(err)
	}
	image.Close()
	after, err := os.ReadFile(options.Script)
	if err != nil || string(source) != string(after) {
		t.Fatalf("source changed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(options.Script), "unused.exe")); !os.IsNotExist(err) {
		t.Fatalf("custom output was not redirected: %v", err)
	}
	if extractor := os.Getenv("WAILS_TEST_7Z"); extractor != "" {
		directory := t.TempDir()
		cmd := exec.Command(extractor, "x", "-y", "-o"+directory, options.Output)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("extract installer: %v\n%s", err, output)
		}
		for name, want := range map[string]string{"app.exe": "primary payload", "helper.exe": "companion payload"} {
			data, err := os.ReadFile(filepath.Join(directory, name))
			if err != nil || string(data) != want {
				t.Fatalf("extracted %s = %q: %v", name, data, err)
			}
		}
	}
}

func TestInvalidInstallerOptions(t *testing.T) {
	for _, change := range []func(*Options){
		func(o *Options) { o.Arch = "386" },
		func(o *Options) { o.Scope = "all" },
		func(o *Options) { o.Script = "" },
		func(o *Options) { o.Primary.Path = "" },
		func(o *Options) { o.Companions[0].Name = "APP.EXE" },
		func(o *Options) { o.Companions[0].Name = "../escape.exe" },
		func(o *Options) { o.Companions[0].Name = "UNINSTALL.exe" },
		func(o *Options) { o.Companions[0].Path = "${ARCH}/helper.exe" },
		func(o *Options) { o.Version = &buildmeta.Info{Core: "65536.0.0"} },
	} {
		options := testOptions(t)
		change(&options)
		if err := Create(context.Background(), options, func(context.Context, string, []string) error {
			t.Fatal("invalid options reached compiler")
			return nil
		}); err == nil {
			t.Fatal("accepted invalid installer options")
		}
	}
}

func BenchmarkInstallerWrapper(b *testing.B) {
	options := Options{Primary: File{Name: "app.exe", Path: "/workspace/bin/app.exe"}, Companions: []File{{Name: "helper.exe", Path: "/workspace/bin/helper.exe"}}, Arch: "amd64", Scope: "user"}
	b.ReportAllocs()
	for b.Loop() {
		script, err := wrapperScript("/workspace/build/windows/nsis/project.nsi", "/workspace/bin/stage/installer.exe", options)
		if err != nil || !strings.Contains(script, "helper.exe") {
			b.Fatal(err)
		}
	}
}
