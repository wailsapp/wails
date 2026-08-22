package manifest

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errInjectedWrite = errors.New("injected write failure")

type faultTemporaryFile struct {
	temporaryFile
	writeErr, chmodErr, syncErr, closeErr error
}

func (f *faultTemporaryFile) Write(data []byte) (int, error) {
	if f.writeErr != nil {
		return 0, f.writeErr
	}
	return f.temporaryFile.Write(data)
}
func (f *faultTemporaryFile) Chmod(mode os.FileMode) error {
	if f.chmodErr != nil {
		return f.chmodErr
	}
	return f.temporaryFile.Chmod(mode)
}
func (f *faultTemporaryFile) Sync() error {
	if f.syncErr != nil {
		return f.syncErr
	}
	return f.temporaryFile.Sync()
}
func (f *faultTemporaryFile) Close() error {
	err := f.temporaryFile.Close()
	if f.closeErr != nil {
		return f.closeErr
	}
	return err
}

func TestAtomicAndExclusiveWriteFaultsNeverPublishPartialData(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*writeOperations)
	}{
		{"directory", func(ops *writeOperations) { ops.mkdirAll = func(string, os.FileMode) error { return errInjectedWrite } }},
		{"temporary create", func(ops *writeOperations) {
			ops.createTemp = func(string, string) (temporaryFile, error) { return nil, errInjectedWrite }
		}},
		{"temporary chmod", wrapTemporaryFault(func(file *faultTemporaryFile) { file.chmodErr = errInjectedWrite })},
		{"temporary write", wrapTemporaryFault(func(file *faultTemporaryFile) { file.writeErr = errInjectedWrite })},
		{"temporary sync", wrapTemporaryFault(func(file *faultTemporaryFile) { file.syncErr = errInjectedWrite })},
		{"temporary close", wrapTemporaryFault(func(file *faultTemporaryFile) { file.closeErr = errInjectedWrite })},
		{"atomic publish", func(ops *writeOperations) { ops.replace = func(string, string) error { return errInjectedWrite } }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			destination := filepath.Join(root, "manifest.hcl")
			require.NoError(t, os.WriteFile(destination, []byte("old"), 0o644))
			ops := osWriteOperations()
			test.mutate(&ops)
			err := atomicWriteWithOperations(destination, []byte("new"), 0o644, ops)
			require.ErrorIs(t, err, errInjectedWrite)
			data, readErr := os.ReadFile(destination)
			require.NoError(t, readErr)
			assert.Equal(t, "old", string(data))
		})
	}

	root := t.TempDir()
	destination := filepath.Join(root, "manifest.hcl")
	ops := osWriteOperations()
	ops.link = func(string, string) error { return errInjectedWrite }
	require.ErrorIs(t, exclusiveWriteWithOperations(destination, []byte("new"), 0o644, ops), errInjectedWrite)
	assert.NoFileExists(t, destination)
	ops = osWriteOperations()
	ops.mkdirAll = func(string, os.FileMode) error { return errInjectedWrite }
	require.ErrorIs(t, exclusiveWriteWithOperations(destination, []byte("new"), 0o644, ops), errInjectedWrite)
	ops = osWriteOperations()
	ops.createTemp = func(string, string) (temporaryFile, error) { return nil, errInjectedWrite }
	require.ErrorIs(t, exclusiveWriteWithOperations(destination, []byte("new"), 0o644, ops), errInjectedWrite)
}

func wrapTemporaryFault(mutate func(*faultTemporaryFile)) func(*writeOperations) {
	return func(ops *writeOperations) {
		create := ops.createTemp
		ops.createTemp = func(directory, pattern string) (temporaryFile, error) {
			file, err := create(directory, pattern)
			if err != nil {
				return nil, err
			}
			fault := &faultTemporaryFile{temporaryFile: file}
			mutate(fault)
			return fault, nil
		}
	}
}

func TestEjectWriterSelectionAndPublicationErrors(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, WriteMinimal(root, Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"}))

	called := ""
	exclusive := func(string, []byte, os.FileMode) error { called = "exclusive"; return nil }
	replace := func(string, []byte, os.FileMode) error { called = "replace"; return nil }
	require.NoError(t, ejectWithWriters(root, "", "v3", false, EncodeEjectedHCL, exclusive, replace))
	assert.Equal(t, "exclusive", called)
	require.NoError(t, ejectWithWriters(root, "", "v3", true, EncodeEjectedHCL, exclusive, replace))
	assert.Equal(t, "replace", called)

	err := ejectWithWriters(root, "", "v3", false, EncodeEjectedHCL, func(string, []byte, os.FileMode) error { return os.ErrExist }, replace)
	require.ErrorContains(t, err, "--force")
	err = ejectWithWriters(root, "", "v3", false, EncodeEjectedHCL, func(string, []byte, os.FileMode) error { return errInjectedWrite }, replace)
	require.ErrorIs(t, err, errInjectedWrite)
	err = ejectWithWriters(t.TempDir(), "", "v3", false, EncodeEjectedHCL, exclusive, replace)
	require.Error(t, err)
	err = ejectWithWriters(root, "", "v3", false, func(Config, string) ([]byte, error) { return nil, errInjectedWrite }, exclusive, replace)
	require.ErrorIs(t, err, errInjectedWrite)
}

func TestMigrationDraftWriterPreservesPublicationErrors(t *testing.T) {
	document := NewDocument(Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"})
	err := writeMigrationDraftAt(t.TempDir(), MigratedFilename, document, nil, func(string, []byte, os.FileMode) error {
		return errInjectedWrite
	})
	require.ErrorIs(t, err, errInjectedWrite)
}
