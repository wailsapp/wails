package binding_test

import (
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/logger"
)

type Message[T any] struct {
	Value T
}

type GenericPointerFixture struct{}

func (g *GenericPointerFixture) GetMessage() Message[*string] {
	return Message[*string]{}
}

// TestGenericPointerTypeNames is a regression test for
// https://github.com/wailsapp/wails/issues/5160
func TestGenericPointerTypeNames(t *testing.T) {
	// given
	generationDir := t.TempDir()

	// setup
	testLogger := &logger.Logger{}
	b := binding.NewBindings(testLogger, []interface{}{&GenericPointerFixture{}}, []interface{}{}, false, []interface{}{})

	// then
	err := b.GenerateGoBindings(generationDir)
	if err != nil {
		t.Fatalf("could not generate the Go bindings: %v", err)
	}

	rawBindings, err := fs.ReadFile(os.DirFS(generationDir), "binding_test/GenericPointerFixture.d.ts")
	if err != nil {
		t.Fatalf("could not read the generated bindings: %v", err)
	}
	rawModels, err := fs.ReadFile(os.DirFS(generationDir), "models.ts")
	if err != nil {
		t.Fatalf("could not read the generated models: %v", err)
	}

	generatedBindings := string(rawBindings)
	generatedModels := string(rawModels)

	if !strings.Contains(generatedBindings, "Message_string_") {
		t.Errorf("the generated bindings do not reference Message_string_:\n%s", generatedBindings)
	}
	if count := strings.Count(generatedModels, "export class Message_string_"); count != 1 {
		t.Errorf("expected exactly one Message_string_ class in models.ts, got %d:\n%s", count, generatedModels)
	}
	if strings.Contains(generatedBindings, "Message__string_") {
		t.Errorf("the generated bindings contain Message__string_:\n%s", generatedBindings)
	}
	if strings.Contains(generatedModels, "Message__string_") {
		t.Errorf("the generated models contain Message__string_:\n%s", generatedModels)
	}
}
