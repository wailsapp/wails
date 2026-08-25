package buildinfo

import (
	"testing"

	wdebug "github.com/wailsapp/wails/v3/internal/debug"
)

func TestGet(t *testing.T) {
	result, err := Get()
	if err != nil {
		t.Error(err)
	}
	_ = result
}

// Development must reflect whether a local Wails source tree is resolvable
// at runtime — not whether the binary happens to have been built inside a
// git checkout. CI-built release artefacts carry `vcs=git` in their build
// metadata, so the old "Development = (vcs == git)" heuristic produced a
// false positive on every user's machine and broke `wails3 init` by
// emitting a stray `replace ... => /v3` in the scaffolded go.mod.
//
// These tests pin the new semantics in place: Development tracks
// wdebug.LocalModulePath, and nothing else.
func TestDevelopmentFalseWhenNoLocalSource(t *testing.T) {
	saved := wdebug.LocalModulePath
	t.Cleanup(func() { wdebug.LocalModulePath = saved })

	wdebug.LocalModulePath = ""

	info, err := Get()
	if err != nil {
		t.Fatalf("buildinfo.Get() returned error: %v", err)
	}
	if info.Development {
		t.Fatalf("Development = true with empty LocalModulePath; want false (this is the released-binary case where vcs=git would have false-positived)")
	}
}

func TestDevelopmentTrueWhenLocalSourceResolves(t *testing.T) {
	saved := wdebug.LocalModulePath
	t.Cleanup(func() { wdebug.LocalModulePath = saved })

	wdebug.LocalModulePath = t.TempDir()

	info, err := Get()
	if err != nil {
		t.Fatalf("buildinfo.Get() returned error: %v", err)
	}
	if !info.Development {
		t.Fatalf("Development = false with LocalModulePath=%q; want true", wdebug.LocalModulePath)
	}
}
