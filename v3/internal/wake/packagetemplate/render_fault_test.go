package packagetemplate

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errInjectedCopy = errors.New("injected copy failure")

type faultReadCloser struct {
	reader   io.Reader
	readErr  error
	closeErr error
}

func (f *faultReadCloser) Read(data []byte) (int, error) {
	if f.readErr != nil {
		return 0, f.readErr
	}
	return f.reader.Read(data)
}
func (f *faultReadCloser) Close() error { return f.closeErr }

type faultStagedFile struct {
	stagedFile
	writeErr, chmodErr, closeErr error
}

func (f *faultStagedFile) Write(data []byte) (int, error) {
	if f.writeErr != nil {
		return 0, f.writeErr
	}
	return f.stagedFile.Write(data)
}
func (f *faultStagedFile) Chmod(mode os.FileMode) error {
	if f.chmodErr != nil {
		return f.chmodErr
	}
	return f.stagedFile.Chmod(mode)
}
func (f *faultStagedFile) Close() error {
	err := f.stagedFile.Close()
	if f.closeErr != nil {
		return f.closeErr
	}
	return err
}

type faultWriteCloser struct {
	io.WriteCloser
	writeErr, closeErr error
}

func (f *faultWriteCloser) Write(data []byte) (int, error) {
	if f.writeErr != nil {
		return 0, f.writeErr
	}
	return f.WriteCloser.Write(data)
}
func (f *faultWriteCloser) Close() error {
	err := f.WriteCloser.Close()
	if f.closeErr != nil {
		return f.closeErr
	}
	return err
}

type faultDirEntry struct{}

func (faultDirEntry) Name() string               { return "nested" }
func (faultDirEntry) IsDir() bool                { return true }
func (faultDirEntry) Type() fs.FileMode          { return fs.ModeDir }
func (faultDirEntry) Info() (fs.FileInfo, error) { return nil, errInjectedCopy }

func TestCopyFileFaultsAreTransactional(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*copyOperations)
	}{
		{"source stat", func(ops *copyOperations) {
			ops.stat = func(string) (fs.FileInfo, error) { return nil, errInjectedCopy }
		}},
		{"destination parent", func(ops *copyOperations) { ops.mkdirAll = func(string, os.FileMode) error { return errInjectedCopy } }},
		{"staging create", func(ops *copyOperations) {
			ops.createTemp = func(string, string) (stagedFile, error) { return nil, errInjectedCopy }
		}},
		{"source open", func(ops *copyOperations) {
			ops.open = func(string) (io.ReadCloser, error) { return nil, errInjectedCopy }
		}},
		{"source read", func(ops *copyOperations) {
			ops.open = func(string) (io.ReadCloser, error) { return &faultReadCloser{readErr: errInjectedCopy}, nil }
		}},
		{"source close", func(ops *copyOperations) {
			ops.open = func(string) (io.ReadCloser, error) {
				return &faultReadCloser{reader: bytes.NewReader([]byte("new")), closeErr: errInjectedCopy}, nil
			}
		}},
		{"staging write", func(ops *copyOperations) {
			create := ops.createTemp
			ops.createTemp = func(dir, pattern string) (stagedFile, error) {
				file, err := create(dir, pattern)
				return &faultStagedFile{stagedFile: file, writeErr: errInjectedCopy}, err
			}
		}},
		{"staging chmod", func(ops *copyOperations) {
			create := ops.createTemp
			ops.createTemp = func(dir, pattern string) (stagedFile, error) {
				file, err := create(dir, pattern)
				return &faultStagedFile{stagedFile: file, chmodErr: errInjectedCopy}, err
			}
		}},
		{"staging close", func(ops *copyOperations) {
			create := ops.createTemp
			ops.createTemp = func(dir, pattern string) (stagedFile, error) {
				file, err := create(dir, pattern)
				return &faultStagedFile{stagedFile: file, closeErr: errInjectedCopy}, err
			}
		}},
		{"replacement prepare", func(ops *copyOperations) {
			ops.mkdirTemp = func(string, string) (string, error) { return "", errInjectedCopy }
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			source := filepath.Join(root, "source")
			destination := filepath.Join(root, "destination")
			require.NoError(t, os.WriteFile(source, []byte("new"), 0o755))
			require.NoError(t, os.WriteFile(destination, []byte("old"), 0o644))
			ops := osCopyOperations()
			test.mutate(&ops)
			err := copyWithOperations(source, destination, ops)
			require.ErrorIs(t, err, errInjectedCopy)
			assert.Equal(t, "old", string(mustReadFile(t, destination)))
		})
	}
}

func TestCopyDirectoryFaultsAreTransactional(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*copyOperations, string, string)
	}{
		{"destination parent", func(ops *copyOperations, _, _ string) {
			ops.mkdirAll = func(string, os.FileMode) error { return errInjectedCopy }
		}},
		{"staging directory", func(ops *copyOperations, _, _ string) {
			ops.mkdirTemp = func(string, string) (string, error) { return "", errInjectedCopy }
		}},
		{"staging mode", func(ops *copyOperations, _, _ string) {
			ops.chmod = func(string, os.FileMode) error { return errInjectedCopy }
		}},
		{"walk", func(ops *copyOperations, _, _ string) {
			ops.walkDir = func(string, fs.WalkDirFunc) error { return errInjectedCopy }
		}},
		{"walk callback", func(ops *copyOperations, source, _ string) {
			ops.walkDir = func(_ string, walk fs.WalkDirFunc) error {
				return walk(filepath.Join(source, "broken"), nil, errInjectedCopy)
			}
		}},
		{"relative path", func(ops *copyOperations, _, _ string) {
			ops.rel = func(string, string) (string, error) { return "", errInjectedCopy }
		}},
		{"directory info", func(ops *copyOperations, source, _ string) {
			ops.walkDir = func(_ string, walk fs.WalkDirFunc) error {
				return walk(filepath.Join(source, "nested"), faultDirEntry{}, nil)
			}
		}},
		{"nested mkdir", func(ops *copyOperations, _, _ string) {
			ops.mkdir = func(string, os.FileMode) error { return errInjectedCopy }
		}},
		{"nested mode", func(ops *copyOperations, _, _ string) {
			chmod := ops.chmod
			ops.chmod = func(path string, mode os.FileMode) error {
				if filepath.Base(path) == "nested" {
					return errInjectedCopy
				}
				return chmod(path, mode)
			}
		}},
		{"file stat", func(ops *copyOperations, _, _ string) {
			stat := ops.stat
			ops.stat = func(path string) (fs.FileInfo, error) {
				if filepath.Base(path) == "payload" {
					return nil, errInjectedCopy
				}
				return stat(path)
			}
		}},
		{"file destination parent", func(ops *copyOperations, _, _ string) {
			mkdirAll := ops.mkdirAll
			calls := 0
			ops.mkdirAll = func(path string, mode os.FileMode) error {
				calls++
				if calls == 2 {
					return errInjectedCopy
				}
				return mkdirAll(path, mode)
			}
		}},
		{"file open", func(ops *copyOperations, _, _ string) {
			ops.open = func(string) (io.ReadCloser, error) { return nil, errInjectedCopy }
		}},
		{"file create", func(ops *copyOperations, _, _ string) {
			ops.openFile = func(string, int, os.FileMode) (io.WriteCloser, error) { return nil, errInjectedCopy }
		}},
		{"file write", func(ops *copyOperations, _, _ string) {
			open := ops.openFile
			ops.openFile = func(path string, flag int, mode os.FileMode) (io.WriteCloser, error) {
				file, err := open(path, flag, mode)
				return &faultWriteCloser{WriteCloser: file, writeErr: errInjectedCopy}, err
			}
		}},
		{"file close", func(ops *copyOperations, _, _ string) {
			open := ops.openFile
			ops.openFile = func(path string, flag int, mode os.FileMode) (io.WriteCloser, error) {
				file, err := open(path, flag, mode)
				return &faultWriteCloser{WriteCloser: file, closeErr: errInjectedCopy}, err
			}
		}},
		{"file mode", func(ops *copyOperations, _, _ string) {
			chmod := ops.chmod
			ops.chmod = func(path string, mode os.FileMode) error {
				if filepath.Base(path) == "payload" {
					return errInjectedCopy
				}
				return chmod(path, mode)
			}
		}},
		{"replacement prepare", func(ops *copyOperations, _, _ string) {
			mkdirTemp := ops.mkdirTemp
			calls := 0
			ops.mkdirTemp = func(dir, pattern string) (string, error) {
				calls++
				if calls == 2 {
					return "", errInjectedCopy
				}
				return mkdirTemp(dir, pattern)
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			source := filepath.Join(root, "source")
			destination := filepath.Join(root, "destination")
			require.NoError(t, os.MkdirAll(filepath.Join(source, "nested"), 0o755))
			require.NoError(t, os.WriteFile(filepath.Join(source, "nested", "payload"), []byte("new"), 0o644))
			require.NoError(t, os.MkdirAll(destination, 0o755))
			require.NoError(t, os.WriteFile(filepath.Join(destination, "previous"), []byte("old"), 0o644))
			ops := osCopyOperations()
			test.mutate(&ops, source, destination)
			err := copyWithOperations(source, destination, ops)
			require.ErrorIs(t, err, errInjectedCopy)
			assert.Equal(t, "old", string(mustReadFile(t, filepath.Join(destination, "previous"))))
		})
	}
}

func TestReplacePathFaultsRestoreThePreviousDestination(t *testing.T) {
	tests := []struct {
		name            string
		destination     bool
		mutate          func(*copyOperations, string, string)
		wantDestination bool
	}{
		{"backup directory", false, func(ops *copyOperations, _, _ string) {
			ops.mkdirTemp = func(string, string) (string, error) { return "", errInjectedCopy }
		}, false},
		{"destination stat", false, func(ops *copyOperations, _, destination string) {
			stat := ops.stat
			ops.stat = func(path string) (fs.FileInfo, error) {
				if path == destination {
					return nil, errInjectedCopy
				}
				return stat(path)
			}
		}, false},
		{"backup rename", true, func(ops *copyOperations, _, destination string) {
			rename := ops.rename
			ops.rename = func(source, target string) error {
				if source == destination {
					return errInjectedCopy
				}
				return rename(source, target)
			}
		}, true},
		{"publish without destination", false, func(ops *copyOperations, source, _ string) {
			rename := ops.rename
			ops.rename = func(from, target string) error {
				if from == source {
					return errInjectedCopy
				}
				return rename(from, target)
			}
		}, false},
		{"publish with destination", true, func(ops *copyOperations, source, _ string) {
			rename := ops.rename
			ops.rename = func(from, target string) error {
				if from == source {
					return errInjectedCopy
				}
				return rename(from, target)
			}
		}, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			source := filepath.Join(root, "staged")
			destination := filepath.Join(root, "destination")
			require.NoError(t, os.WriteFile(source, []byte("new"), 0o644))
			if test.destination {
				require.NoError(t, os.WriteFile(destination, []byte("old"), 0o644))
			}
			ops := osCopyOperations()
			test.mutate(&ops, source, destination)
			err := replacePath(source, destination, ops)
			require.ErrorIs(t, err, errInjectedCopy)
			if test.wantDestination {
				assert.Equal(t, "old", string(mustReadFile(t, destination)))
			} else {
				assert.NoFileExists(t, destination)
			}
		})
	}
}
