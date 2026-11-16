package commands

import (
	"context"
	"github.com/wailsapp/wails/v3/internal/wake/cache"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
	"os"
	"path/filepath"
	"testing"
)

func TestCompileCPUFeatureEnvironmentInvalidatesCache(t *testing.T) {
	node := pipeline.Node{Kind: pipeline.CompileApplication, Spec: pipeline.CompileSpec{TargetOS: "linux", TargetArch: "amd64", Production: true, Toolchain: "native"}}
	handler := &manifestHandler{}
	t.Setenv("GOAMD64", "v3")
	before, err := handler.Identity(context.Background(), node)
	if err != nil {
		t.Fatal(err)
	}
	beforeKey, err := cache.ActionKey(string(node.Kind), map[string]any{"spec": node.Spec, "tool": before}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("GOAMD64", "v1")
	after, err := handler.Identity(context.Background(), node)
	if err != nil {
		t.Fatal(err)
	}
	afterKey, err := cache.ActionKey(string(node.Kind), map[string]any{"spec": node.Spec, "tool": after}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if before == after || beforeKey == afterKey {
		t.Fatal("GOAMD64 change must invalidate compiler identity and action key")
	}

}

func TestCompileGoEnvironmentFileInvalidatesCache(t *testing.T) {
	config := filepath.Join(t.TempDir(), "go-env")
	t.Setenv("GOENV", config)
	node := pipeline.Node{Kind: pipeline.CompileApplication, Spec: pipeline.CompileSpec{Production: true, Toolchain: "native"}}
	handler := &manifestHandler{}
	if err := os.WriteFile(config, []byte("GOAMD64=v3\n"), 0600); err != nil {
		t.Fatal(err)
	}
	before, err := handler.Identity(t.Context(), node)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(config, []byte("GOAMD64=v1\n"), 0600); err != nil {
		t.Fatal(err)
	}
	after, err := handler.Identity(t.Context(), node)
	if err != nil {
		t.Fatal(err)
	}
	if before == after {
		t.Fatal("go env -w configuration must invalidate compiler identity")
	}
}
