package commands

import (
	"context"
	"crypto/sha1"
	"debug/elf"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/wailsapp/wails/v3/internal/appimage"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func fixtureELF(arch string, image bool) []byte {
	data := make([]byte, 64)
	copy(data, "\x7fELF")
	data[4], data[5], data[6] = 2, 1, 1
	if image {
		copy(data[8:], "AI\x02")
	}
	machine := elf.EM_X86_64
	if arch == "arm64" {
		machine = elf.EM_AARCH64
	}
	binary.LittleEndian.PutUint16(data[16:], uint16(elf.ET_EXEC))
	binary.LittleEndian.PutUint16(data[18:], uint16(machine))
	binary.LittleEndian.PutUint32(data[20:], 1)
	binary.LittleEndian.PutUint16(data[52:], 64)
	binary.LittleEndian.PutUint16(data[54:], 56)
	binary.LittleEndian.PutUint16(data[58:], 64)
	return data
}

func configureNativeAppImageFixture(t *testing.T, root string, spec *pipeline.PackageSpec) {
	t.Helper()
	if runtime.GOOS != "linux" || runtime.GOARCH != spec.TargetArch {
		t.Skip("AppImage requires a matching Linux host")
	}
	spec.Project.Kind = "cli"
	if err := os.WriteFile(filepath.Join(root, spec.Binary), fixtureELF(spec.TargetArch, false), 0755); err != nil {
		t.Fatal(err)
	}
	tools := t.TempDir()
	for _, name := range []string{"linuxdeploy", "AppRun"} {
		if err := os.WriteFile(filepath.Join(tools, name), []byte("external tool fixture"), 0755); err != nil {
			t.Fatal(err)
		}
	}
	previous, previousTools := manifestAppImageCreate, manifestAppImageTools
	t.Cleanup(func() { manifestAppImageCreate, manifestAppImageTools = previous, previousTools })
	manifestAppImageTools = func(appimage.Options, context.Context) (appimage.Tools, error) {
		return appimage.Tools{LinuxDeploy: filepath.Join(tools, "linuxdeploy"), AppRun: filepath.Join(tools, "AppRun")}, nil
	}
	manifestAppImageCreate = func(ctx context.Context, options appimage.Options) error {
		options.Tools = appimage.Tools{LinuxDeploy: filepath.Join(tools, "linuxdeploy"), AppRun: filepath.Join(tools, "AppRun")}
		options.Runner = func(ctx context.Context, command appimage.Command) error {
			image := fixtureELF(options.Arch, true)
			output := filepath.Join(command.Directory, command.Env["OUTPUT"])
			if err := os.WriteFile(output, image, 0755); err != nil {
				return err
			}
			if options.UpdateInfo != "" {
				return os.WriteFile(output+".zsync", []byte(fmt.Sprintf("zsync: 0.6.2\nFilename: %s\nLength: %d\nSHA-1: %x\n\n", filepath.Base(output), len(image), sha1.Sum(image))), 0644)
			}
			return nil
		}
		return appimage.Create(ctx, options)
	}
}

type appImageIntegrationHandler struct{ *manifestHandler }

func (h appImageIntegrationHandler) Run(ctx context.Context, node pipeline.Node) (pipeline.RunResult, error) {
	if spec, ok := node.Spec.(pipeline.CompileSpec); ok {
		path := filepath.Join(h.root, spec.Output)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return pipeline.RunResult{}, err
		}
		return pipeline.RunResult{}, os.WriteFile(path, fixtureELF(spec.TargetArch, false), 0755)
	}
	return h.manifestHandler.Run(ctx, node)
}

func TestAppImageSidecarsAreRestoredAndCollected(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux AppImage adapter")
	}
	root := executableFixture(t)
	source := `version: 3
project: {name: image, product_name: Image, identifier: com.example.image, version: 1.0.0, kind: cli}
frontend: {disabled: true}
packages:
  appimage:
    checksum: true
    update_info: zsync|https://example.test/image.zsync
`
	if err := os.WriteFile(filepath.Join(root, "wails.yaml"), []byte(source), 0644); err != nil {
		t.Fatal(err)
	}
	loaded, err := manifest.Load(root, "")
	if err != nil {
		t.Fatal(err)
	}
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build", TargetOS: "linux", TargetArch: runtime.GOARCH, Formats: []string{"appimage"}})
	if err != nil {
		t.Fatal(err)
	}
	key := pipeline.NodeKey("package:linux/" + runtime.GOARCH + ":appimage")
	spec := plan.Nodes[key].Spec.(pipeline.PackageSpec)
	if err := os.MkdirAll(filepath.Dir(filepath.Join(root, spec.Binary)), 0755); err != nil {
		t.Fatal(err)
	}
	configureNativeAppImageFixture(t, root, &spec)
	executor := pipeline.Executor{Handler: appImageIntegrationHandler{&manifestHandler{root: root, config: loaded.Config}}}
	for iteration := 0; iteration < 3; iteration++ {
		if _, err := executor.Execute(t.Context(), plan, pipeline.ExecuteOptions{Root: root}); err != nil {
			t.Fatal(err)
		}
		receipt, err := pipeline.VerifyArtifactReceipt(root, ".wails/artifacts/receipt.json")
		if err != nil {
			t.Fatal(err)
		}
		if len(receipt.Artifacts) != 3 {
			t.Fatalf("sidecars missing from receipt: %+v", receipt)
		}
		if iteration == 0 {
			if err := os.RemoveAll(filepath.Join(root, plan.Nodes[key].Output)); err != nil {
				t.Fatal(err)
			}
		}
		if iteration == 1 {
			if err := os.RemoveAll(filepath.Join(root, "bin")); err != nil {
				t.Fatal(err)
			}
		}
	}
}
