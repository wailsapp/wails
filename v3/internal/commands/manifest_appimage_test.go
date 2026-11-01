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
	spec.Project.Kind = "cli"
	if err := os.WriteFile(filepath.Join(root, spec.Binary), fixtureELF(spec.TargetArch, false), 0755); err != nil {
		t.Fatal(err)
	}
	previous, previousIdentity := manifestAppImageCreate, manifestAppImageIdentity
	t.Cleanup(func() { manifestAppImageCreate, manifestAppImageIdentity = previous, previousIdentity })
	manifestAppImageIdentity = func(appimage.Options, context.Context) (string, error) {
		return "runtime fixture", nil
	}
	manifestAppImageCreate = func(ctx context.Context, options appimage.Options) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		image := fixtureELF(options.Arch, true)
		if err := os.MkdirAll(filepath.Dir(options.Output), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(options.Output, image, 0755); err != nil {
			return err
		}
		if options.UpdateInfo != "" {
			return os.WriteFile(options.Output+".zsync", []byte(fmt.Sprintf("zsync: 0.6.2\nFilename: %s\nLength: %d\nSHA-1: %x\n\n", filepath.Base(options.Output), len(image), sha1.Sum(image))), 0644)
		}
		return nil
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

func TestAppImageRuntimeIdentityIncludesGUICompanion(t *testing.T) {
	previous := manifestAppImageIdentity
	t.Cleanup(func() { manifestAppImageIdentity = previous })
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"linux-root", "runtime"} {
		if err := os.MkdirAll(filepath.Join(root, path), 0755); err != nil {
			t.Fatal(err)
		}
	}
	fingerprint := "first"
	manifestAppImageIdentity = func(options appimage.Options, ctx context.Context) (string, error) {
		if !options.WebKit || options.GTKVersion != 3 || len(options.Binaries) != 2 || options.RuntimeRoot != filepath.Join(root, "linux-root") || options.RuntimePath != filepath.Join(root, "runtime") {
			t.Fatalf("lost runtime options: %+v", options)
		}
		return fingerprint, nil
	}
	handler := &manifestHandler{root: root}
	spec := pipeline.PackageSpec{Format: "appimage", TargetArch: "amd64", GTKVersion: 3,
		Project: manifest.Project{Kind: "cli", BinaryName: "cli"}, Binary: "bin/cli",
		Companions: []pipeline.NamedBinary{{Name: "gui", Path: "bin/gui", Kind: "gui"}},
		Config:     manifest.PackageFormat{RuntimeRoot: "linux-root", RuntimePath: "runtime"}}
	node := pipeline.Node{Kind: pipeline.PackageArtifact, Spec: spec}
	first, err := handler.Identity(t.Context(), node)
	if err != nil {
		t.Fatal(err)
	}
	fingerprint = "changed content"
	second, err := handler.Identity(t.Context(), node)
	if err != nil || first == second {
		t.Fatal("runtime content change did not invalidate package", err)
	}
	paths, err := handler.packageToolPaths(t.Context(), spec)
	if err != nil || len(paths) != 0 {
		t.Fatal("native AppImage unexpectedly requires external tools", paths, err)
	}
}
