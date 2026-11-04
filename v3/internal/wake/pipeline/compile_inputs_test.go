package pipeline

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/cache"
)

func TestExecutableInputsFollowImportsAndEmbeddedFiles(t *testing.T) {
	root := t.TempDir()
	write := func(name, value string) {
		t.Helper()
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(value), 0644); err != nil {
			t.Fatal(err)
		}
	}
	write("go.mod", "module example.test/app\ngo 1.26\n")
	write("cmd/first/main.go", "package main\nimport (_ \"example.test/app/shared\";_ \"embed\")\n//go:embed assets/*\nvar page string\nfunc main() {}\n")
	write("cmd/second/main.go", "package main\nimport _ \"example.test/app/shared\"\nfunc main() {}\n")
	write("shared/value.go", "package shared\nconst Value = 1\n")
	write("cmd/first/assets/page.txt", "one")
	write("unrelated/helper.go", "package unrelated\nconst Value = 1\n")
	original := []InputSpec{{Label: "go-sources", Root: ".", IncludeGoEmbed: true, IncludeExtensions: goSourceExtensions(), ExcludeDirs: []string{".git", ".wails"}, ExcludeSuffixes: []string{"_test.go"}}, {Label: "go-module", Files: []string{"go.mod"}}}
	snapshot := func(entry string) []string {
		t.Helper()
		store, err := cache.OpenCache(root)
		if err != nil {
			t.Fatal(err)
		}
		node := Node{Inputs: executableInputs(root, entry, original)}
		digest, err := snapshotNodeInputs(store, node)
		if err != nil {
			t.Fatal(err)
		}
		return digest
	}
	first, second := snapshot("cmd/first"), snapshot("cmd/second")
	write("unrelated/helper.go", "package unrelated\nconst Value = 2\n")
	if !reflect.DeepEqual(first, snapshot("cmd/first")) || !reflect.DeepEqual(second, snapshot("cmd/second")) {
		t.Fatal("unimported package invalidated executables")
	}
	write("cmd/first/assets/page.txt", "two")
	if reflect.DeepEqual(first, snapshot("cmd/first")) {
		t.Fatal("embedded file edit was missed")
	}
	if !reflect.DeepEqual(second, snapshot("cmd/second")) {
		t.Fatal("independent executable tracks another executable's embed")
	}
	first = snapshot("cmd/first")
	write("shared/value.go", "package shared\nconst Value = 2\n")
	if reflect.DeepEqual(first, snapshot("cmd/first")) || reflect.DeepEqual(second, snapshot("cmd/second")) {
		t.Fatal("shared import must invalidate both executables")
	}
	first = snapshot("cmd/first")
	write("cmd/first/extra.go", "package main\nimport _ \"example.test/app/extra\"\n")
	write("extra/value.go", "package extra\nconst Value = 1\n")
	if reflect.DeepEqual(first, snapshot("cmd/first")) {
		t.Fatal("new Go file or import was missed")
	}
	first = snapshot("cmd/first")
	write("extra/value.go", "package extra\nconst Value = 2\n")
	if reflect.DeepEqual(first, snapshot("cmd/first")) {
		t.Fatal("newly imported package edit was missed")
	}
	first = snapshot("cmd/first")
	write("shared/native.h", "#define VALUE 1\n")
	if reflect.DeepEqual(first, snapshot("cmd/first")) {
		t.Fatal("native header edit was missed")
	}
	first = snapshot("cmd/first")
	if err := os.Remove(filepath.Join(root, "cmd/first/extra.go")); err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(first, snapshot("cmd/first")) {
		t.Fatal("removed Go file or import was missed")
	}
}

func TestExecutableInputsRemainConservativeForIncompleteGraphs(t *testing.T) {
	root := t.TempDir()
	for name, value := range map[string]string{"go.mod": "module example.test/app\n", "main.go": "package main\nimport _ \"example.test/app/generated\"\n"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(value), 0644); err != nil {
			t.Fatal(err)
		}
	}
	original := []InputSpec{{Label: "go-sources", Root: ".", IncludeExtensions: goSourceExtensions()}}
	if !reflect.DeepEqual(original, executableInputs(root, ".", original)) {
		t.Fatal("incomplete generated package graph must retain broad inputs")
	}
}
