package service

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/wailsapp/wails/v3/internal/flags"
)

func BenchmarkServiceModuleTemplate(b *testing.B) {
	source, err := serviceTemplate.ReadFile("template/go.mod.tmpl")
	if err != nil {
		b.Fatal(err)
	}
	data := struct{ Name string }{Name: "example_service"}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		module, err := template.New("go.mod").Parse(string(source))
		if err != nil {
			b.Fatal(err)
		}
		if err := module.Execute(io.Discard, data); err != nil {
			b.Fatal(err)
		}
	}
}

func TestInstallGeneratedServiceCompiles(t *testing.T) {
	directory := t.TempDir()
	if err := Install(&flags.ServiceInit{Name: "example_service", OutputDir: directory}); err != nil {
		t.Fatal(err)
	}
	moduleDir, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	// Compile against this checkout; the standalone scaffold pins an older release.
	cmd := exec.Command("go", "test", "-mod=readonly", "-tags=server", filepath.Join(directory, "service.go"))
	cmd.Dir = moduleDir
	cmd.Env = append(os.Environ(), "GOWORK=off")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("generated service does not compile: %v\n%s", err, output)
	}
}
