package pipeline

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/report"
)

// obfuscatedIDsHandler mimics the bindings generator, which writes
// wails_obfuscated.gen.go beside the main package as well as the bindings.
type obfuscatedIDsHandler struct{ *fakeHandler }

func (h obfuscatedIDsHandler) Run(ctx context.Context, node Node) (RunResult, error) {
	result, err := h.fakeHandler.Run(ctx, node)
	if err == nil && node.Kind == GenerateBindings {
		err = os.WriteFile(filepath.Join(h.root, "wails_obfuscated.gen.go"), []byte("package main\n"), 0o644)
	}
	return result, err
}

func TestObfuscatedBindingsRegenerateWhenIDFileIsMissing(t *testing.T) {
	config := performanceConfig(t)
	plan, err := PlanBuild(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64", Obfuscated: true})
	require.NoError(t, err)
	handler := obfuscatedIDsHandler{&fakeHandler{root: config.Root}}
	executor := Executor{Handler: handler}
	options := ExecuteOptions{Root: config.Root, Reporter: report.Nop{}}
	_, err = executor.Execute(context.Background(), plan, options)
	require.NoError(t, err)
	idFile := filepath.Join(config.Root, "wails_obfuscated.gen.go")
	require.FileExists(t, idFile)

	_, err = executor.Execute(context.Background(), plan, options)
	require.NoError(t, err)
	assert.FileExists(t, idFile, "a warm run keeps the file")

	require.NoError(t, os.Remove(idFile))
	_, err = executor.Execute(context.Background(), plan, options)
	require.NoError(t, err)
	assert.FileExists(t, idFile, "bindings must regenerate the missing obfuscated ID file instead of restoring from cache")
}
