package commands

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/cache"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func executableFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"go.mod":             "module example.com/executables\n\ngo 1.26.0\n",
		"main.go":            "package main\nimport \"fmt\"\nconst AppVersion=\"1.2.3-beta.4\"\nfunc main(){fmt.Println(\"primary\",AppVersion)}\n",
		"cmd/helper/main.go": "package main\nimport \"fmt\"\nfunc main(){fmt.Println(\"companion one\")}\n",
		"wails.yaml": `version: 3
project:
  name: executables
  product_name: Executables
  company: Fixture
  description: Native executable integration fixture
  identifier: com.example.executables
  kind: cli
  version_file: main.go
  version_symbol: AppVersion
frontend:
  disabled: true
binaries:
  helper:
    package: cmd/helper
    kind: cli
`,
	}
	for name, contents := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(contents), 0644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestYAMLNamedExecutablesBuildRunCacheAndRebuild(t *testing.T) {
	root := executableFixture(t)
	t.Chdir(root)
	options := manifestRunOptions{Verb: "build", TargetOS: runtime.GOOS, TargetArch: runtime.GOARCH}
	if _, err := runManifestPipelineResult(options); err != nil {
		t.Fatal(err)
	}
	extension := ""
	if runtime.GOOS == "windows" {
		extension = ".exe"
	}
	check := func(name, want string) {
		t.Helper()
		data, err := exec.Command(filepath.Join(root, "bin", name+extension)).CombinedOutput()
		if err != nil || strings.TrimSpace(string(data)) != want {
			t.Fatalf("%s output=%s err=%v", name, data, err)
		}
	}
	check("executables", "primary 1.2.3-beta.4")
	check("helper", "companion one")
	receipt, err := pipeline.VerifyArtifactReceipt(root, ".wails/artifacts/receipt.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(receipt.Artifacts) != 2 {
		t.Fatalf("receipt=%+v", receipt)
	}
	warm, err := runManifestPipelineResult(options)
	if err != nil {
		t.Fatal(err)
	}
	for key, result := range warm.Results {
		if result.Status == cache.LookupMiss {
			t.Errorf("unchanged step %s was rebuilt", key)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "cmd/helper/main.go"), []byte("package main\nimport \"fmt\"\nfunc main(){fmt.Println(\"companion two\")}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	rebuilt, err := runManifestPipelineResult(options)
	if err != nil {
		t.Fatal(err)
	}
	if rebuilt.Results[pipeline.NodeKey("target:"+runtime.GOOS+"/"+runtime.GOARCH+":compile")].Status == cache.LookupMiss {
		t.Error("editing a companion rebuilt the independent primary executable")
	}
	check("helper", "companion two")
	if _, err := pipeline.VerifyArtifactReceipt(root, ".wails/artifacts/receipt.json"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "Taskfile.yml")); !os.IsNotExist(err) {
		t.Fatalf("build created a Taskfile: %v", err)
	}
}

func TestNativeMSIXAdapterContainsCompanionOnAnyHost(t *testing.T) {
	for _, arch := range []string{"amd64", "arm64"} {
		t.Run(arch, func(t *testing.T) {
			root := executableFixture(t)
			loaded, err := manifest.Load(root, "")
			if err != nil {
				t.Fatal(err)
			}
			plan, err := pipeline.PlanBuildForHost(loaded.Config, pipeline.Request{Verb: "build", TargetOS: "windows", TargetArch: arch, Formats: []string{"msix"}}, pipeline.NewHostCapabilities(runtime.GOOS, runtime.GOARCH, []string{"go"}, nil))
			if err != nil {
				t.Fatal(err)
			}
			handler := &manifestHandler{root: root, config: loaded.Config}
			_, err = (pipeline.Executor{Handler: handler}).Execute(context.Background(), plan, pipeline.ExecuteOptions{Root: root})
			if err != nil {
				t.Fatal(err)
			}
			archive, err := zip.OpenReader(filepath.Join(root, "bin", "executables.msix"))
			if err != nil {
				t.Fatal(err)
			}
			defer archive.Close()
			names := map[string]bool{}
			for _, file := range archive.File {
				names[file.Name] = true
				if file.Name == "AppxManifest.xml" {
					reader, err := file.Open()
					if err != nil {
						t.Fatal(err)
					}
					data, err := io.ReadAll(reader)
					closeErr := reader.Close()
					if err != nil || closeErr != nil {
						t.Fatal(err, closeErr)
					}
					if !strings.Contains(string(data), `Version="1.2.3.0"`) {
						t.Fatalf("non-numeric MSIX version: %s", data)
					}
				}
			}
			for _, name := range []string{"executables.exe", "helper.exe", "AppxBlockMap.xml", "[Content_Types].xml", "AppxManifest.xml"} {
				if !names[name] {
					t.Errorf("MSIX omitted %s", name)
				}
			}

			if sdk := os.Getenv("WAILS_TEST_MAKEMSIX"); sdk != "" {
				unpacked := filepath.Join(t.TempDir(), "unpacked")
				command := exec.CommandContext(t.Context(), sdk, "unpack", "/p", filepath.Join(root, "bin", "executables.msix"), "/d", unpacked)
				if output, err := command.CombinedOutput(); err != nil {
					t.Fatalf("Microsoft SDK rejected native YAML package: %v\n%s", err, output)
				}
				for _, file := range archive.File {
					if file.FileInfo().IsDir() || file.Name == "[Content_Types].xml" || file.Name == "AppxBlockMap.xml" {
						continue
					}
					input, err := file.Open()
					if err != nil {
						t.Fatal(err)
					}
					expected, readErr := io.ReadAll(input)
					closeErr := input.Close()
					if readErr != nil || closeErr != nil {
						t.Fatal(readErr, closeErr)
					}
					actual, err := os.ReadFile(filepath.Join(unpacked, filepath.FromSlash(file.Name)))
					if err != nil {
						t.Fatal(err)
					}
					if !bytes.Equal(expected, actual) {
						t.Fatalf("SDK extracted different bytes for %s", file.Name)
					}
				}
			}
		})
	}
}

func TestYAMLCommitInjectionUsesEachExecutablesSymbol(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("requires Git")
	}
	root := executableFixture(t)
	t.Chdir(root)
	for name, value := range map[string]string{
		"main.go":            "package main\nimport \"fmt\"\nconst AppVersion=\"1.2.3-beta.4\"\nvar Commit string\nfunc main(){fmt.Print(Commit)}\n",
		"cmd/helper/main.go": "package main\nimport \"fmt\"\nvar Commit string\nfunc main(){fmt.Print(Commit)}\n",
	} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(value), 0644); err != nil {
			t.Fatal(err)
		}
	}
	yamlPath := filepath.Join(root, "wails.yaml")
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		t.Fatal(err)
	}
	source := strings.Replace(string(data), "  version_file:", "  commit_symbol: main.Commit\n  version_file:", 1)
	source = strings.Replace(source, "    kind: cli", "    kind: cli\n    commit_symbol: main.Commit", 1)
	if err := os.WriteFile(yamlPath, []byte(source), 0644); err != nil {
		t.Fatal(err)
	}
	git := func(args ...string) string {
		t.Helper()
		command := exec.Command("git", args...)
		command.Dir = root
		out, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	git("init")
	git("add", ".")
	git("-c", "user.name=Fixture", "-c", "user.email=fixture@example.test", "commit", "-m", "fixture")
	revision := git("rev-parse", "HEAD")
	if _, err := runManifestPipelineResult(manifestRunOptions{Verb: "build", TargetOS: runtime.GOOS, TargetArch: runtime.GOARCH}); err != nil {
		t.Fatal(err)
	}
	extension := ""
	if runtime.GOOS == "windows" {
		extension = ".exe"
	}
	for _, name := range []string{"executables", "helper"} {
		out, err := exec.Command(filepath.Join(root, "bin", name+extension)).CombinedOutput()
		if err != nil || string(out) != revision {
			t.Fatalf("%s revision=%s want=%s err=%v", name, out, revision, err)
		}
	}
}

func TestYAMLUPXFailurePreservesPublishedExecutable(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux executable compression integration")
	}
	root := executableFixture(t)
	t.Chdir(root)
	tools := t.TempDir()
	upx := filepath.Join(tools, "upx")
	if err := os.WriteFile(upx, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", tools+string(os.PathListSeparator)+os.Getenv("PATH"))
	yamlPath := filepath.Join(root, "wails.yaml")
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		t.Fatal(err)
	}
	source := strings.Replace(string(data), "  version_file:", "  upx_level: 1\n  version_file:", 1)
	if err := os.WriteFile(yamlPath, []byte(source), 0644); err != nil {
		t.Fatal(err)
	}
	options := manifestRunOptions{Verb: "build", TargetOS: runtime.GOOS, TargetArch: runtime.GOARCH}
	if _, err := runManifestPipelineResult(options); err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(root, "bin", "executables")
	before, err := os.ReadFile(binary)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(upx, []byte("#!/bin/sh\necho compression-fixture-failure >&2\nexit 23\n"), 0755); err != nil {
		t.Fatal(err)
	}
	if _, err := runManifestPipelineResult(options); err == nil {
		t.Fatal("expected UPX failure to invalidate cached compile")
	}
	after, err := os.ReadFile(binary)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("UPX failure replaced the published executable")
	}
	out, err := exec.Command(binary).CombinedOutput()
	if err != nil || strings.TrimSpace(string(out)) != "primary 1.2.3-beta.4" {
		t.Fatalf("last executable unusable: %s %v", out, err)
	}
}
