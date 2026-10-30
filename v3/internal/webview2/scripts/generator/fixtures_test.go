package generator

import (
	"testing"
	"updater/generator/idlreader"
	"updater/generator/types"
)

func generateFixture(t testing.TB, source []byte) []*types.GeneratedFile {
	t.Helper()
	idl, err := idlreader.Parse(source)
	if err != nil {
		t.Fatal(err)
	}
	if err := idl.Process(); err != nil {
		t.Fatal(err)
	}
	files, err := idl.Generate()
	if err != nil {
		t.Fatal(err)
	}
	// The fixtures cover declarations; com.go is the shared runtime preamble.
	return files[1:]
}

func checkGeneratedFiles(t testing.TB, got, want []*types.GeneratedFile) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("generated %d files, want %d", len(got), len(want))
	}
	type output struct{ name, packageName, content string }
	remaining := make(map[output]int, len(want))
	for _, file := range want {
		remaining[output{file.FileName, file.Package, file.Content.String()}]++
	}
	for _, file := range got {
		key := output{file.FileName, file.Package, file.Content.String()}
		if remaining[key] == 0 {
			t.Fatalf("unexpected generated file %s (package %s):\n%s", file.FileName, file.Package, key.content)
		}
		remaining[key]--
	}
}
