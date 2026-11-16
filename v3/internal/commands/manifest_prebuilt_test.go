package commands

import (
	"archive/tar"
	"bytes"
	"context"
	"crypto/sha256"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/universalbinary"
	"github.com/wailsapp/wails/v3/internal/wake/cache"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func buildPrebuiltFixture(t *testing.T, target pipeline.Target, value string) []byte {
	t.Helper()
	directory := t.TempDir()
	source := `package main
import "os"
func main() { if err := os.WriteFile("prebuilt-result", []byte("` + value + `:" + os.Getenv("PREBUILT_VALUE") + ":" + os.Args[1]), 0600); err != nil { panic(err) } }
`
	input, output := filepath.Join(directory, "main.go"), filepath.Join(directory, "app")
	if err := os.WriteFile(input, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	command := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", output, input)
	command.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS="+target.OS, "GOARCH="+target.Arch, "GOWORK=off", "GOFLAGS=-buildvcs=false")
	if data, err := command.CombinedOutput(); err != nil {
		t.Fatalf("build fixture: %v\n%s", err, data)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestPrebuiltTargetValidation(t *testing.T) {
	root := t.TempDir()
	darwin := map[string]string{}
	for _, platform := range []string{"linux", "windows", "darwin"} {
		for _, arch := range []string{"amd64", "arm64"} {
			target := pipeline.Target{OS: platform, Arch: arch}
			path := filepath.Join(root, platform+"-"+arch)
			if err := os.WriteFile(path, buildPrebuiltFixture(t, target, "target"), 0644); err != nil {
				t.Fatal(err)
			}
			if err := validatePrebuiltTarget(path, target); err != nil {
				t.Fatalf("%s/%s: %v", platform, arch, err)
			}
			other := "arm64"
			if arch == "arm64" {
				other = "amd64"
			}
			if err := validatePrebuiltTarget(path, pipeline.Target{OS: platform, Arch: other}); err == nil {
				t.Fatalf("accepted %s/%s as %s", platform, arch, other)
			}
			if platform == "darwin" {
				darwin[arch] = path
			}
		}
	}
	fat := filepath.Join(root, "universal")
	if err := universalbinary.Create(fat, []string{darwin["amd64"], darwin["arm64"]}); err != nil {
		t.Fatal(err)
	}
	if err := validatePrebuiltTarget(fat, pipeline.Target{OS: "darwin", Arch: "universal"}); err != nil {
		t.Fatal(err)
	}
	if err := validatePrebuiltTarget(darwin["arm64"], pipeline.Target{OS: "darwin", Arch: "universal"}); err == nil {
		t.Fatal("accepted a thin binary as universal")
	}
	if err := validatePrebuiltTarget(darwin["arm64"], pipeline.Target{OS: "linux", Arch: "arm64"}); err == nil {
		t.Fatal("accepted Mach-O as Linux ELF")
	}
}

func TestPrebuiltPackagesUseInputBytesWithoutBuildTools(t *testing.T) {
	first := buildPrebuiltFixture(t, pipeline.Target{OS: "linux", Arch: "amd64"}, "first")
	second := buildPrebuiltFixture(t, pipeline.Target{OS: "linux", Arch: "amd64"}, "second")
	handler, _ := linuxPackageFixture(t)
	root := handler.root
	if err := manifest.WriteMinimal(root, manifest.Project{Name: "example", ProductName: "Example", Identifier: "com.example.prebuilt", Version: "1.0.0", Kind: "cli"}); err != nil {
		t.Fatal(err)
	}
	loaded, err := manifest.Load(root, "")
	if err != nil {
		t.Fatal(err)
	}
	loaded.Config.Binaries = []manifest.Binary{{Name: "helper", Package: ".", Kind: "cli"}}
	handler.config = loaded.Config
	for _, name := range []string{"example", "helper"} {
		if err := os.WriteFile(filepath.Join(root, "bin", name), first, 0644); err != nil {
			t.Fatal(err)
		}
	}
	before, err := os.Stat(filepath.Join(root, "bin", "example"))
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", t.TempDir())
	plan, err := pipeline.PlanBuildForHost(loaded.Config, pipeline.Request{Verb: "package", TargetOS: "linux", TargetArch: "amd64", Formats: []string{"deb"}, SkipBuild: true}, pipeline.NewHostCapabilities(runtime.GOOS, runtime.GOARCH, nil, nil))
	if err != nil {
		t.Fatal(err)
	}
	executor := pipeline.Executor{Handler: handler}
	options := pipeline.ExecuteOptions{Root: root, Workers: 2}
	if _, err := executor.Execute(context.Background(), plan, options); err != nil {
		t.Fatal(err)
	}
	packagePath := filepath.Join(root, filepath.FromSlash(plan.Nodes[plan.Artifacts[0]].Output))
	data, err := os.ReadFile(packagePath)
	if err != nil {
		t.Fatal(err)
	}
	firstHash := sha256.Sum256(data)
	archive := tar.NewReader(bytes.NewReader(debTarMember(t, data, "data.tar.gz")))
	found := 0
	for {
		entry, err := archive.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		name := strings.TrimPrefix(entry.Name, "./")
		if name != "usr/local/bin/example" && name != "usr/bin/helper" {
			continue
		}
		payload, err := io.ReadAll(archive)
		if err != nil || !bytes.Equal(payload, first) || entry.Mode&0777 != 0755 {
			t.Fatalf("incorrect prebuilt payload %s: %v", name, err)
		}
		found++
	}
	if found != 2 {
		t.Fatalf("found %d executable payloads", found)
	}
	after, err := os.Stat(filepath.Join(root, "bin", "example"))
	if err != nil || !before.ModTime().Equal(after.ModTime()) || before.Mode() != after.Mode() {
		t.Fatalf("prebuilt input was changed: %v", err)
	}
	results, err := executor.Execute(context.Background(), plan, options)
	if err != nil || results["package:linux/amd64:deb"].Status != cache.LookupHit {
		t.Fatalf("unchanged package did not reuse cache: %v %v", results, err)
	}
	if err := os.WriteFile(filepath.Join(root, "bin", "example"), second, 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := executor.Execute(context.Background(), plan, options); err != nil {
		t.Fatal(err)
	}
	data, err = os.ReadFile(packagePath)
	if err != nil || sha256.Sum256(data) == firstHash {
		t.Fatalf("changed binary reused old package: %v", err)
	}
	if err := os.Remove(filepath.Join(root, "bin", "helper")); err != nil {
		t.Fatal(err)
	}
	if _, err := executor.Execute(context.Background(), plan, options); err == nil {
		t.Fatal("missing prebuilt input reused cached artifacts")
	}
}

func TestPrebuiltRunWithoutBuildTools(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "windows" && runtime.GOOS != "darwin" {
		t.Skip("desktop target required")
	}
	data := buildPrebuiltFixture(t, pipeline.Target{OS: runtime.GOOS, Arch: runtime.GOARCH}, "run")
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "bin"), 0700); err != nil {
		t.Fatal(err)
	}
	name := "prebuilt"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	if err := os.WriteFile(filepath.Join(root, "bin", name), data, 0644); err != nil {
		t.Fatal(err)
	}
	if err := manifest.WriteMinimal(root, manifest.Project{Name: "prebuilt", ProductName: "Prebuilt", Identifier: "com.example.prebuilt", Version: "1.0.0", Kind: "cli"}); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
	t.Setenv("PATH", t.TempDir())
	t.Setenv("WAILS_EXP_USE_WAKE", "1")
	t.Setenv("PREBUILT_VALUE", "environment")
	if err := RunApplication(&RunOptions{SkipBuild: true}, []string{"argument"}); err != nil {
		t.Fatal(err)
	}
	result, err := os.ReadFile(filepath.Join(root, "prebuilt-result"))
	if err != nil || string(result) != "run:environment:argument" {
		t.Fatalf("run result %q: %v", result, err)
	}
}

func TestPrebuiltFlagsRequireManifestAndReachPipeline(t *testing.T) {
	for _, active := range []bool{false, true} {
		called := false
		operations := manifestCommandOperations{
			active: func() (bool, error) { return active, nil },
			run: func(options manifestRunOptions) error {
				called = true
				if !options.SkipBuild {
					t.Fatal("skip-build lost before pipeline")
				}
				return nil
			},
			task: func(string, []string) error { t.Fatal("prebuilt command fell back to Task"); return nil },
		}
		for _, run := range []func() error{
			func() error {
				return buildWithOperations(&flags.Build{SkipBuild: true, Formats: "deb"}, nil, operations)
			},
			func() error { return packageWithOperations(&flags.Package{SkipBuild: true}, nil, operations) },
		} {
			called = false
			err := run()
			if active != (err == nil) || called != active {
				t.Fatalf("active=%v called=%v error=%v", active, called, err)
			}
		}
	}
}

func TestPrebuiltRejectsNonExecutableAndUnsafeInputs(t *testing.T) {
	root := t.TempDir()
	handler := &manifestHandler{root: root}
	spec := pipeline.PrebuiltSpec{Source: "input", Output: ".wails/imported", Target: pipeline.Target{OS: "linux", Arch: "amd64"}}
	if err := os.WriteFile(filepath.Join(root, "input"), []byte("not an executable"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := handler.importBinary(spec); err == nil {
		t.Fatal("accepted non-executable contents")
	}
	if err := os.Remove(filepath.Join(root, "input")); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, "input"), 0700); err != nil {
		t.Fatal(err)
	}
	if _, err := handler.importBinary(spec); err == nil || !strings.Contains(err.Error(), "regular file") {
		t.Fatalf("directory input: %v", err)
	}
	if err := os.Remove(filepath.Join(root, "input")); err != nil {
		t.Fatal(err)
	}
	t.Run("symlink", func(t *testing.T) {
		if err := os.Symlink(filepath.Join(root, "other"), filepath.Join(root, "input")); err != nil {
			t.Skipf("symlink unavailable: %v", err)
		}
		if _, err := handler.importBinary(spec); err == nil || !strings.Contains(err.Error(), "regular file") {
			t.Fatalf("symlink input: %v", err)
		}
	})
}
