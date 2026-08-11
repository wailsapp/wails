package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/generator/config"
)

const projectionFixture = "github.com/wailsapp/wails/v3/internal/generator/testdata/service_projection"

func TestFindProjectedService(t *testing.T) {
	pkgs, err := LoadPackages(nil, projectionFixture)
	if err != nil {
		t.Fatal(err)
	}
	systemPaths, err := ResolveSystemPaths(nil)
	if err != nil {
		t.Fatal(err)
	}

	services, _, err := FindServices(pkgs, systemPaths, config.NullLogger)
	if err != nil {
		t.Fatal(err)
	}

	var found []*ServiceBinding
	for service := range services {
		found = append(found, service)
	}
	if len(found) != 1 {
		t.Fatalf("found %d services, want 1", len(found))
	}
	if got := found[0].Type.Name(); got != "Service" {
		t.Fatalf("service type = %s, want Service", got)
	}
	if got := found[0].Type.Pkg().Path(); got != projectionFixture+"/backend" {
		t.Fatalf("service package = %s, want %s/backend", got, projectionFixture)
	}
	if found[0].Projection == nil || found[0].Projection.Name() != "FrontendAPI" {
		t.Fatalf("service projection = %v, want FrontendAPI", found[0].Projection)
	}
}

func TestGenerateProjectedService(t *testing.T) {
	bindingsDir := t.TempDir()
	metadataDir := filepath.Join(t.TempDir(), "metadata")
	if err := os.MkdirAll(metadataDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(metadataDir, "package.go"), []byte("package metadata\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	options := &flags.GenerateBindingsOptions{
		OutputDirectory:   bindingsDir,
		ModelsFilename:    "models",
		IndexFilename:     "index",
		TimeType:          "string",
		TS:                true,
		UseBundledRuntime: true,
		Obfuscated:        true,
		ObfuscatedOutput:  metadataDir,
		NoEvents:          true,
	}

	generator := NewGenerator(options, config.DirCreator(bindingsDir), config.NullLogger)
	if _, err := generator.Generate(projectionFixture); err != nil {
		t.Fatal(err)
	}

	serviceFile := filepath.Join(bindingsDir, projectionFixture, "backend", "service.ts")
	generated, err := os.ReadFile(serviceFile)
	if err != nil {
		t.Fatal(err)
	}
	code := string(generated)
	for _, method := range []string{"Echo", "Version"} {
		if !strings.Contains(code, "export function "+method+"(") {
			t.Errorf("generated binding does not contain projected method %s", method)
		}
	}
	if count := strings.Count(code, "export function "); count != 2 {
		t.Errorf("generated binding contains %d methods, want 2", count)
	}
	if strings.Contains(code, "BackendOnly") {
		t.Error("generated binding contains a method outside the projection")
	}
	if strings.Contains(code, "ServiceStartup") {
		t.Error("generated binding contains a lifecycle method")
	}

	metadata, err := os.ReadFile(filepath.Join(metadataDir, bindingIDMetadataFile))
	if err != nil {
		t.Fatal(err)
	}
	metadataCode := string(metadata)
	for _, method := range []string{"Echo", "Version"} {
		if !strings.Contains(metadataCode, ")."+method+",") {
			t.Errorf("obfuscation metadata does not contain projected method %s", method)
		}
	}
	if count := strings.Count(metadataCode, "application.RegisterBindingMethodID("); count != 2 {
		t.Errorf("obfuscation metadata contains %d method registrations, want 2", count)
	}
	if strings.Contains(metadataCode, "BackendOnly") {
		t.Error("obfuscation metadata contains a method outside the projection")
	}
}
