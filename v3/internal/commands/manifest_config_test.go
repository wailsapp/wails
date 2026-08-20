package commands

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestEjectCommandContractAndFailures(t *testing.T) {
	want := errors.New("injected failure")
	var output bytes.Buffer
	operations := ejectOperations{
		getwd:   func() (string, error) { return "/project", nil },
		write:   func(string, string, string, bool) error { return nil },
		version: "3.0.0",
		output:  &output,
	}

	require.ErrorContains(t, ejectWithOperations(&EjectOptions{}, []string{"release"}, operations), "usage")

	getwdFailure := operations
	getwdFailure.getwd = func() (string, error) { return "", want }
	require.ErrorIs(t, ejectWithOperations(&EjectOptions{}, nil, getwdFailure), want)

	writeFailure := operations
	writeFailure.write = func(string, string, string, bool) error { return want }
	require.ErrorIs(t, ejectWithOperations(&EjectOptions{}, nil, writeFailure), want)

	writerFailure := operations
	writerFailure.output = failingWriter{err: want}
	require.ErrorIs(t, ejectWithOperations(&EjectOptions{}, nil, writerFailure), want)

	var gotRoot, gotProfile, gotVersion string
	var gotForce bool
	operations.write = func(root, profile, version string, force bool) error {
		gotRoot, gotProfile, gotVersion, gotForce = root, profile, version, force
		return nil
	}
	require.NoError(t, ejectWithOperations(&EjectOptions{Force: true}, nil, operations))
	assert.Equal(t, "/project", gotRoot)
	assert.Empty(t, gotProfile)
	assert.Equal(t, "3.0.0", gotVersion)
	assert.True(t, gotForce)
	assert.Contains(t, output.String(), "wails.ejected.hcl")
}

type failingWriter struct{ err error }

func (w failingWriter) Write([]byte) (int, error) { return 0, w.err }

func BenchmarkNativeManifestRouting(b *testing.B) {
	root := b.TempDir()
	require.NoError(b, manifest.WriteMinimal(root, manifest.Project{Name: "bench", ProductName: "Bench", Identifier: "com.example.bench", Version: "1.0.0"}))
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		active, err := activeManifestProjectAt(root)
		if err != nil || !active {
			b.Fatalf("route active manifest: active=%v err=%v", active, err)
		}
	}
}

func BenchmarkEjectCommand(b *testing.B) {
	root := b.TempDir()
	require.NoError(b, manifest.WriteMinimal(root, manifest.Project{Name: "bench", ProductName: "Bench", Identifier: "com.example.bench", Version: "1.0.0"}))
	operations := ejectOperations{
		getwd:   func() (string, error) { return root, nil },
		write:   manifest.Eject,
		version: "3.0.0",
		output:  io.Discard,
	}
	require.NoError(b, ejectWithOperations(&EjectOptions{Force: true}, nil, operations))
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if err := ejectWithOperations(&EjectOptions{Force: true}, nil, operations); err != nil {
			b.Fatal(err)
		}
	}
}
