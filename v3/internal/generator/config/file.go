package config

import (
	"io"
	"os"
	"path/filepath"
)

// FileCreator abstracts away file and directory creation.
// We use this to implement tests cleanly.
//
// Paths are always relative to the output directory.
//
// A FileCreator must allow concurrent calls to Create transparently.
// Each [io.WriteCloser] instance returned by a call to Create
// will be used by one goroutine at a time; but distinct instances
// must support concurrent use by distinct goroutines.
type FileCreator interface {
	Create(path string) (io.WriteCloser, error)
}

// FileCreatorFunc is an adapter to allow
// the use of ordinary functions as file creators.
type FileCreatorFunc func(path string) (io.WriteCloser, error)

// Create calls f(path).
func (f FileCreatorFunc) Create(path string) (io.WriteCloser, error) {
	return f(path)
}

// NullCreator is a dummy file creator implementation.
// Calls to Create never fail and return
// a writer that discards all incoming data.
var NullCreator FileCreator = FileCreatorFunc(func(path string) (io.WriteCloser, error) {
	return nullWriteCloser{}, nil
})

// DirCreator returns a file creator that creates files
// relative to the given output directory.
//
// It joins the output directory and the file path,
// calls [os.MkdirAll] on the directory part of the result,
// then [os.Create] on the full file path.
func DirCreator(outputDir string) FileCreator {
	return FileCreatorFunc(func(path string) (io.WriteCloser, error) {
		path = filepath.Join(outputDir, path)

		if err := os.MkdirAll(filepath.Dir(path), 0o777); err != nil {
			return nil, err
		}

		return os.Create(path)
	})
}

type nullWriteCloser struct{}

func (nullWriteCloser) Write(data []byte) (int, error) {
	return len(data), nil
}

func (nullWriteCloser) Close() error {
	return nil
}
