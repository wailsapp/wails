package generator

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/generator/config"
)

func TestGeneratedJavaScriptCanBeImportedWithGenericModels(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is required for executable generated JavaScript tests")
	}

	outputDir := t.TempDir()
	options := &flags.GenerateBindingsOptions{
		ModelsFilename:    "models",
		IndexFilename:     "index",
		UseBundledRuntime: true,
		TimeType:          "Date",
	}

	generator := NewGenerator(options, config.DirCreator(outputDir), config.NullLogger)
	_, err = generator.Generate("github.com/wailsapp/wails/v3/internal/generator/testcases/type_mapping")
	if err != nil {
		t.Fatal(err)
	}

	modelsPath := filepath.Join(
		outputDir,
		"github.com/wailsapp/wails/v3/internal/generator/testcases/type_mapping/models.js",
	)
	source, err := os.ReadFile(modelsPath)
	if err != nil {
		t.Fatal(err)
	}

	const runtimeImport = `import { Create as $Create } from "/wails/runtime.js";`
	const runtimeStub = `
const $$identity = value => value;
const $Create = new Proxy({}, {
    get: (_, name) => name === "Any" || name === "ByteSlice" || name === "DateFromTime"
        ? $$identity
        : () => $$identity,
});`
	if !strings.Contains(string(source), runtimeImport) {
		t.Fatalf("generated models do not contain expected runtime import %q", runtimeImport)
	}

	testSource := strings.Replace(string(source), runtimeImport, runtimeStub, 1) + `
const $$instance = AllTypes.createFrom({
    GenericType: { ID: 1, Data: 2, List: [3] },
});
if (!($$instance.GenericType instanceof GenericType)) {
    throw new Error("generic model field was not created");
}
`
	testModelsPath := filepath.Join(filepath.Dir(modelsPath), "models.test.mjs")
	if err := os.WriteFile(testModelsPath, []byte(testSource), 0600); err != nil {
		t.Fatal(err)
	}

	command := exec.Command(node, testModelsPath)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("generated JavaScript could not be imported: %v\n%s", err, output)
	}
}
