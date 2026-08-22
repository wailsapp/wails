package binding_test

import (
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/logger"
)

type GenericPointerBindings struct{}

type Message[T any] struct {
	Value T `json:"value"`
}

func (h *GenericPointerBindings) RoundTrip(input Message[*string]) Message[*string] {
	return input
}

func TestGenericPointerBindingsUseSameTypeName(t *testing.T) {
	generationDir := t.TempDir()

	testLogger := &logger.Logger{}
	b := binding.NewBindings(testLogger, []interface{}{&GenericPointerBindings{}}, []interface{}{}, false, []interface{}{})

	err := b.GenerateGoBindings(generationDir)
	if err != nil {
		t.Fatalf("could not generate the Go bindings: %v", err)
	}

	rawGeneratedBindings, err := fs.ReadFile(os.DirFS(generationDir), "binding_test/GenericPointerBindings.d.ts")
	if err != nil {
		t.Fatalf("could not read the generated bindings: %v", err)
	}
	generatedBindings := string(rawGeneratedBindings)

	if !strings.Contains(generatedBindings, "binding_test.Message_string_") {
		t.Fatalf("generated bindings should reference Message_string_, got:\n%s", generatedBindings)
	}
	if strings.Contains(generatedBindings, "Message__string_") {
		t.Fatalf("generated bindings should not reference Message__string_, got:\n%s", generatedBindings)
	}

	rawGeneratedModels, err := fs.ReadFile(os.DirFS(generationDir), "models.ts")
	if err != nil {
		t.Fatalf("could not read the generated models: %v", err)
	}
	generatedModels := string(rawGeneratedModels)

	if !strings.Contains(generatedModels, "export class Message_string_") {
		t.Fatalf("generated models should define Message_string_, got:\n%s", generatedModels)
	}
	if strings.Contains(generatedModels, "Message__string_") {
		t.Fatalf("generated models should not define Message__string_, got:\n%s", generatedModels)
	}
}
