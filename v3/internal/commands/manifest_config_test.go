package commands

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
