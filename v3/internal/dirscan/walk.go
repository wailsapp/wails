// Package dirscan scans project trees with bounded filesystem concurrency.
// It does not follow symlinks. Callbacks may run concurrently and must protect
// shared state; all callbacks have finished when Walk returns.
package dirscan

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
)

type directory struct {
	path  string
	entry fs.DirEntry
}
type result struct {
	children []directory
	err      error
}

func Walk(root string, workers int, visit fs.WalkDirFunc) error {
	if workers <= 0 {
		workers = max(1, runtime.GOMAXPROCS(0))
	}
	info, err := os.Lstat(root)
	if err != nil {
		return walkResult(visit(root, nil, err))
	}
	entry := fs.FileInfoToDirEntry(info)
	if err := visit(root, entry, nil); err != nil {
		return walkResult(err)
	}
	if !entry.IsDir() {
		return nil
	}
	scan := scanner{visit: visit}
	jobs := make(chan directory)
	results := make(chan result, workers)
	var group sync.WaitGroup
	for range workers {
		group.Go(func() {
			for job := range jobs {
				results <- scan.directory(job)
			}
		})
	}
	pending := []directory{{root, entry}}
	active := 0
	var firstErr error
	for len(pending) > 0 || active > 0 {
		var ready chan directory
		var next directory
		if len(pending) > 0 {
			ready = jobs
			next = pending[len(pending)-1]
		}
		select {
		case ready <- next:
			pending = pending[:len(pending)-1]
			active++
		case result := <-results:
			active--
			if result.err != nil && !scan.stopped.Swap(true) {
				firstErr = walkResult(result.err)
				pending = nil
			}
			if !scan.stopped.Load() {
				pending = append(pending, result.children...)
			}
		}
	}
	close(jobs)
	group.Wait()
	return firstErr
}

type scanner struct {
	visit   fs.WalkDirFunc
	stopped atomic.Bool
}

func (scan *scanner) directory(current directory) result {
	if scan.stopped.Load() {
		return result{}
	}
	file, err := os.Open(current.path)
	if err != nil {
		return result{err: directoryError(scan.visit(current.path, current.entry, err))}
	}
	entries, readErr := file.ReadDir(-1)
	closeErr := file.Close()
	if err := errors.Join(readErr, closeErr); err != nil {
		return result{err: directoryError(scan.visit(current.path, current.entry, err))}
	}
	var children []directory
	for _, entry := range entries {
		if scan.stopped.Load() {
			return result{}
		}
		path := filepath.Join(current.path, entry.Name())
		err := scan.visit(path, entry, nil)
		if errors.Is(err, fs.SkipDir) {
			if entry.IsDir() {
				continue
			}
			break
		}
		if err != nil {
			return result{err: err}
		}
		if entry.IsDir() {
			children = append(children, directory{path, entry})
		}
	}
	return result{children: children}
}

func walkResult(err error) error {
	if errors.Is(err, fs.SkipDir) || errors.Is(err, fs.SkipAll) {
		return nil
	}
	return err
}

func directoryError(err error) error {
	if errors.Is(err, fs.SkipDir) {
		return nil
	}
	return err
}
