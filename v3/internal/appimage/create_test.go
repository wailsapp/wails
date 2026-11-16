package appimage

import (
	"context"
	"debug/elf"
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"
)

func minimalELF(arch string, appImage bool) []byte {
	data := make([]byte, 64)
	copy(data, "\x7fELF")
	data[4] = 2
	data[5] = 1
	data[6] = 1
	if appImage {
		copy(data[8:], "AI\x02")
	}
	binary.LittleEndian.PutUint16(data[16:], uint16(elf.ET_EXEC))
	binary.LittleEndian.PutUint16(data[18:], uint16(machineForArch(arch)))
	binary.LittleEndian.PutUint32(data[20:], 1)
	binary.LittleEndian.PutUint16(data[52:], 64)
	binary.LittleEndian.PutUint16(data[54:], 56)
	binary.LittleEndian.PutUint16(data[58:], 64)
	return data
}

func appImageFixture(t testing.TB) Options {
	t.Helper()
	if runtime.GOOS != "linux" || (runtime.GOARCH != "amd64" && runtime.GOARCH != "arm64") {
		t.Skip("AppImage creation needs supported native Linux host")
	}
	root := t.TempDir()
	for name, data := range map[string][]byte{"app": minimalELF(runtime.GOARCH, false), "helper": minimalELF(runtime.GOARCH, false), "icon.png": []byte("icon"), "app.desktop": []byte("[Desktop Entry]\nName=App\nExec=app\nIcon=app\nType=Application\n"), "linuxdeploy": []byte("tool"), "AppRun": []byte("runtime")} {
		if err := os.WriteFile(filepath.Join(root, name), data, 0755); err != nil {
			t.Fatal(err)
		}
	}
	return Options{Root: root, Name: "My App", Arch: runtime.GOARCH, Output: "output with spaces/My.AppImage", Icon: "icon.png", IconName: "app.png", DesktopFile: "app.desktop", WorkDir: "work with spaces", Tools: Tools{LinuxDeploy: "linuxdeploy", AppRun: "AppRun"}, Binaries: []Binary{{"app", "app"}, {"helper", "helper"}}, Env: map[string]string{"OUTPUT": "unchanged"}}
}

func TestCreateStagesCompanionsAndPreservesProcessState(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	inherited := os.Getenv("OUTPUT")
	options := appImageFixture(t)
	options.Runner = func(ctx context.Context, command Command) error {
		want := []string{"--appimage-extract-and-run", "--appdir", filepath.Join(command.Directory, "AppDir"), "--output", "appimage"}
		if !reflect.DeepEqual(command.Args, want) {
			t.Fatalf("args=%q", command.Args)
		}
		if command.Program != filepath.Join(options.Root, "linuxdeploy") || !filepath.IsAbs(command.Directory) {
			t.Fatalf("command=%+v", command)
		}
		if command.Env["OUTPUT"] != "result.AppImage" {
			t.Fatalf("env=%v", command.Env)
		}
		for _, name := range []string{"app", "helper"} {
			source := filepath.Join(command.Directory, "AppDir", "usr", "bin", name)
			info, err := os.Stat(source)
			if err != nil || info.Mode().Perm() != 0755 {
				t.Fatalf("%s mode: %v", name, err)
			}
		}
		for path, want := range map[string]string{".DirIcon": "icon", "app.png": "icon", "AppRun": "runtime", "app.desktop": "[Desktop Entry]\nName=App\nExec=app\nIcon=app\nType=Application\n"} {
			data, err := os.ReadFile(filepath.Join(command.Directory, "AppDir", path))
			if err != nil || string(data) != want {
				t.Fatalf("%s: %q, %v", path, data, err)
			}
		}
		return os.WriteFile(filepath.Join(command.Directory, command.Env["OUTPUT"]), minimalELF(runtime.GOARCH, true), 0600)
	}
	if err := Create(context.Background(), options); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(options.Root, options.Output)
	if err := verify(output, runtime.GOARCH); err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(output); err != nil || info.Mode().Perm() != 0755 {
		t.Fatal("output is not executable")
	}
	after, err := os.Getwd()
	if err != nil || after != cwd || os.Getenv("OUTPUT") != inherited || options.Env["OUTPUT"] != "unchanged" {
		t.Fatal("process state or caller env changed")
	}
	entries, err := os.ReadDir(filepath.Join(options.Root, options.WorkDir))
	if err != nil || len(entries) != 0 {
		t.Fatalf("scratch leaked: %v, %v", entries, err)
	}
}

func TestCreateFailurePreservesPreviousImage(t *testing.T) {
	for _, failure := range []string{"tool", "invalid output", "cancel"} {
		t.Run(failure, func(t *testing.T) {
			options := appImageFixture(t)
			output := filepath.Join(options.Root, options.Output)
			if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(output, []byte("previous image"), 0755); err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			options.Runner = func(ctx context.Context, command Command) error {
				if failure == "tool" {
					return errors.New("tool failure")
				}
				data := []byte("invalid image")
				if failure == "cancel" {
					data = minimalELF(runtime.GOARCH, true)
					cancel()
				}
				return os.WriteFile(filepath.Join(command.Directory, "result.AppImage"), data, 0755)
			}
			if err := Create(ctx, options); err == nil {
				t.Fatal("expected failure")
			}
			data, err := os.ReadFile(output)
			if err != nil || string(data) != "previous image" {
				t.Fatalf("old output=%q, %v", data, err)
			}
			entries, err := os.ReadDir(filepath.Join(options.Root, options.WorkDir))
			if err != nil || len(entries) != 0 {
				t.Fatalf("scratch leaked: %v, %v", entries, err)
			}
		})
	}
}

func TestConcurrentImagesUseIndependentStages(t *testing.T) {
	options := appImageFixture(t)
	options.Runner = func(ctx context.Context, command Command) error {
		return os.WriteFile(filepath.Join(command.Directory, command.Env["OUTPUT"]), minimalELF(runtime.GOARCH, true), 0755)
	}
	var group sync.WaitGroup
	for _, name := range []string{"one.AppImage", "two.AppImage"} {
		group.Go(func() {
			copy := options
			copy.Output = filepath.Join("output with spaces", name)
			if err := Create(context.Background(), copy); err != nil {
				t.Error(err)
			}
		})
	}
	group.Wait()
	for _, name := range []string{"one.AppImage", "two.AppImage"} {
		if err := verify(filepath.Join(options.Root, "output with spaces", name), runtime.GOARCH); err != nil {
			t.Fatal(err)
		}
	}
}

func TestOutputNamesAndVerification(t *testing.T) {
	name, err := OutputName("My App", "arm64")
	if err != nil || name != "my-app-aarch64.AppImage" {
		t.Fatalf("%s, %v", name, err)
	}
	for _, pair := range [][2]string{{"../escape", "amd64"}, {"app", "386"}} {
		if _, err := OutputName(pair[0], pair[1]); err == nil {
			t.Fatal("accepted invalid output")
		}
	}
	root := t.TempDir()
	filename := filepath.Join(root, "image")
	if err := os.WriteFile(filename, minimalELF("arm64", true), 0755); err != nil {
		t.Fatal(err)
	}
	if err := verify(filename, "amd64"); err == nil {
		t.Fatal("accepted wrong architecture")
	}
}

func TestStageProgressPreservesFailure(t *testing.T) {
	var messages []string
	options := Options{Progress: func(message string) { messages = append(messages, message) }}
	failure := errors.New("deployment failed")
	err := options.runStage("Deploy dependencies and compress AppImage", func() error { return failure })
	if !errors.Is(err, failure) || len(messages) != 2 || !strings.Contains(messages[1], "failed in") {
		t.Fatalf("progress=%v err=%v", messages, err)
	}
}
