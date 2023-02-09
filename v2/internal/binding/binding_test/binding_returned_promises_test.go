package binding_test

import (
	"io/fs"
	"os"
	"testing"

	"github.com/ciderapp/wails/v2/internal/binding"
	"github.com/ciderapp/wails/v2/internal/logger"
)

const expectedPromiseBindings = `// Copyright Cider Collective
// This file is automatically generated. DO NOT EDIT

export function ErrorReturn(arg1:number):Promise<void>;

export function NoReturn(arg1:string):Promise<void>;

export function SingleReturn(arg1:any):Promise<number>;

export function SingleReturnWithError(arg1:number):Promise<string>;

export function TwoReturn(arg1:any):Promise<string|number>;
`

type PromisesTest struct{}

func (h *PromisesTest) NoReturn(_ string)                           {}
func (h *PromisesTest) ErrorReturn(_ int) error                     { return nil }
func (h *PromisesTest) SingleReturn(_ interface{}) int              { return 0 }
func (h *PromisesTest) SingleReturnWithError(_ int) (string, error) { return "", nil }
func (h *PromisesTest) TwoReturn(_ interface{}) (string, int)       { return "", 0 }

func TestPromises(t *testing.T) {
	// given
	generationDir := t.TempDir()

	// setup
	testLogger := &logger.Logger{}
	b := binding.NewBindings(testLogger, []interface{}{&PromisesTest{}}, []interface{}{}, false)

	// then
	err := b.GenerateGoBindings(generationDir)
	if err != nil {
		t.Fatalf("could not generate the Go bindings: %v", err)
	}

	// then
	rawGeneratedBindings, err := fs.ReadFile(os.DirFS(generationDir), "binding_test/PromisesTest.d.ts")
	if err != nil {
		t.Fatalf("could not read the generated bindings: %v", err)
	}

	// then
	generatedBindings := string(rawGeneratedBindings)
	if generatedBindings != expectedPromiseBindings {
		t.Fatalf("the generated bindings does not match the expected ones.\nWanted:\n%s\n\nGot:\n%s", expectedPromiseBindings, generatedBindings)
	}
}
