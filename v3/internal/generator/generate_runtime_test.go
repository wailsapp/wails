package generator

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/generator/config"
)

func TestGeneratedJavaScriptCanBeImportedWithGenericModels(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is required for executable generated JavaScript tests")
	}

	for name, assertions := range map[string]string{
		"type_mapping": `
const $$instance = AllTypes.createFrom({
    GenericType: { ID: 1, Data: 2, List: [3] },
});
if (!($$instance.GenericType instanceof GenericType)) {
    throw new Error("generic model field was not created");
}
`,
		"generic_class_type_deps": `
for (let i = 0; i < 2; i++) {
    const instance = Alpha.createFrom({
        Box: { IntItems: [1], Zarray: [{ Items: ["one"], Ref: { IntItems: [2], Zarray: [] } }] },
    });
    if (!(instance.Box instanceof Epsilon) ||
        !(instance.Box.Zarray[0] instanceof Zeta) ||
        !(instance.Box.Zarray[0].Ref instanceof Epsilon)) {
        throw new Error("mutually dependent generic models were not created");
    }
}
for (const [convert, input, expected] of [[Number, "42", 42], [String, 42, "42"]]) {
    const instance = Zeta.createFrom(convert)({
        Items: [input], Ref: { Zarray: [{ Items: [input], Ref: null }] },
    });
    if (instance.Items[0] !== expected || instance.Ref.Zarray[0].Items[0] !== expected) {
        throw new Error("generic creators reused another instantiation's converter");
    }
}
`,
	} {
		t.Run(name, func(t *testing.T) {
			testGeneratedJavaScript(t, node, name, assertions)
		})
	}
}

// testGeneratedJavaScript imports generated models and executes their creators.
func testGeneratedJavaScript(t *testing.T, node, name, assertions string) {
	t.Helper()
	outputDir := t.TempDir()
	options := &flags.GenerateBindingsOptions{
		ModelsFilename:    "models",
		IndexFilename:     "index",
		UseBundledRuntime: true,
		TimeType:          "Date",
	}

	generator := NewGenerator(options, config.DirCreator(outputDir), config.NullLogger)
	packagePath := "github.com/wailsapp/wails/v3/internal/generator/testcases/" + name
	_, err := generator.Generate(packagePath)
	if err != nil {
		t.Fatal(err)
	}

	modelsPath := filepath.Join(
		outputDir,
		packagePath, "models.js",
	)
	source, err := os.ReadFile(modelsPath)
	if err != nil {
		t.Fatal(err)
	}

	const runtimeImport = `import { Create as $Create } from "/wails/runtime.js";`
	const runtimeStub = `
const $$identity = value => value;
const $Create = new Proxy({}, {
    get: (_, name) => {
        if (name === "Array") return create => source => (source ?? []).map(create);
        if (name === "Nullable") return create => source => source === null ? null : create(source);
        return name === "Any" || name === "ByteSlice" || name === "DateFromTime"
            ? $$identity
            : () => $$identity;
    },
});`
	if !strings.Contains(string(source), runtimeImport) {
		t.Fatalf("generated models do not contain expected runtime import %q", runtimeImport)
	}

	testSource := strings.Replace(string(source), runtimeImport, runtimeStub, 1) + assertions
	testModelsPath := filepath.Join(filepath.Dir(modelsPath), "models.test.mjs")
	if err := os.WriteFile(testModelsPath, []byte(testSource), 0600); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, node, testModelsPath)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("generated JavaScript could not be imported: %v\n%s", err, output)
	}
}
