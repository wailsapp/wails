package binding_test

import (
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/leaanthony/slicer"
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/typescriptify"
)

type Message[T any] struct {
	Value T
}

type Status[T any] int

func (s Status[T]) TSName() string { return "Ready" }

type Status_string_ struct {
	Name string
}

type GenericEnumCollisionFixture struct{}

func (g *GenericEnumCollisionFixture) GetStatus() Status_string_ {
	return Status_string_{}
}

func TestGenericEnumNameCollisions(t *testing.T) {
	tests := []struct {
		name    string
		structs []interface{}
		enums   []interface{}
	}{
		{
			name:    "enum and class",
			structs: []interface{}{&GenericEnumCollisionFixture{}},
			enums:   []interface{}{[]Status[*string]{1}},
		},
		{
			name:  "distinct enums",
			enums: []interface{}{[]Status[string]{1}, []Status[*string]{1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := binding.NewBindings(&logger.Logger{}, tt.structs, nil, false, tt.enums)
			models, err := b.GenerateModels()
			if err == nil {
				t.Fatalf("expected an error for the Status_string_ name collision, got models:\n%s", models)
			}
			if !strings.Contains(err.Error(), "Status_string_") {
				t.Errorf("expected the error to identify Status_string_, got: %v", err)
			}
		})
	}
}

func TestGenericTypeConversionRepeated(t *testing.T) {
	converter := typescriptify.New()
	converter.Namespace = "binding_test"
	converter.KnownEnums = &slicer.StringSlicer{}
	converter.Add(Message[string]{})
	converter.Add(Message[*string]{})
	converter.AddEnum([]Status[*string]{1})
	converter.AddEnum([]Status[*string]{1})

	for i := 0; i < 2; i++ {
		models, err := converter.Convert(nil)
		if err != nil {
			t.Fatalf("conversion %d failed: %v", i+1, err)
		}
		if count := strings.Count(models, "export class Message_string_"); count != 1 {
			t.Errorf("conversion %d: expected one Message_string_ class, got %d:\n%s", i+1, count, models)
		}
		if count := strings.Count(models, "export enum Status_string_"); count != 1 {
			t.Errorf("conversion %d: expected one Status_string_ enum, got %d:\n%s", i+1, count, models)
		}
	}
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
