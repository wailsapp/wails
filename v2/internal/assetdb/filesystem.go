// +build !desktop
package assetdb

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

var errWhence = errors.New("Seek: invalid whence")
var errOffset = errors.New("Seek: invalid offset")

// Open implements the http.FileSystem interface for the AssetDB
func (a *AssetDB) Open(name string) (http.File, error) {
	if name == "/" || name == "" {
		return &Entry{name: "/", dir: true, db: a}, nil
	}

	if data, ok := a.db[name]; ok {
		return &Entry{name: name, data: data, size: len(data)}, nil
	} else {
		for n, _ := range a.db {
			if strings.HasPrefix(n, name+"/") {
				return &Entry{name: name, db: a, dir: true}, nil
			}
		}
	}
	return &Entry{}, os.ErrNotExist
}

// readdir returns the directory entries for the requested directory
func (a *AssetDB) readdir(name string) ([]os.FileInfo, error) {
	dir := name
	var ents []string

	fim := make(map[string]os.FileInfo)
	for fn, data := range a.db {
		if strings.HasPrefix(fn, dir) {
			pieces := strings.Split(fn[len(dir)+1:], "/")
			if len(pieces) == 1 {
				fim[pieces[0]] = FI{name: pieces[0], dir: false, size: len(data)}
				ents = append(ents, pieces[0])
			} else {
				fim[pieces[0]] = FI{name: pieces[0], dir: true, size: -1}
				ents = append(ents, pieces[0])
			}
		}
	}

	if len(ents) == 0 {
		return nil, os.ErrNotExist
	}

	sort.Strings(ents)
	var list []os.FileInfo
	for _, dir := range ents {
		list = append(list, fim[dir])
	}
	return list, nil
}

// Entry implements the http.File interface to allow for
// use in the http.FileSystem implementation of AssetDB
type Entry struct {
	name string
	data []byte
	dir  bool
	size int
	db   *AssetDB
	off  int
}

// Close is a noop
func (e Entry) Close() error {
	return nil
}

// Read fills the slice provided returning how many bytes were written
// and any errors encountered
func (e *Entry) Read(p []byte) (n int, err error) {
	if e.off >= e.size {
		return 0, io.EOF
	}
	numout := len(p)
	if numout > e.size {
		numout = e.size
	}
	for i := 0; i < numout; i++ {
		p[i] = e.data[e.off+i]
	}
	e.off += numout
	n = int(numout)
	err = nil
	return
}

// Seek seeks to the specified offset from the location specified by whence
func (e *Entry) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	default:
		return 0, errWhence
	case io.SeekStart:
		offset += 0
	case io.SeekCurrent:
		offset += int64(e.off)
	case io.SeekEnd:
		offset += int64(e.size)
	}

	if offset < 0 {
		return 0, errOffset
	}
	e.off = int(offset)
	return offset, nil
}

// Readdir returns the directory entries inside this entry if it is a directory
func (e Entry) Readdir(count int) ([]os.FileInfo, error) {
	ents := []os.FileInfo{}
	if !e.dir {
		return ents, errors.New("Not a directory")
	}
	return e.db.readdir(e.name)
}

// Stat returns information about this directory entry
func (e Entry) Stat() (os.FileInfo, error) {
	return FI{e.name, e.size, e.dir}, nil
}

// FI is the AssetDB implementation of os.FileInfo.
type FI struct {
	name string
	size int
	dir  bool
}

// IsDir returns true if this is a directory
func (fi FI) IsDir() bool {
	return fi.dir
}

// ModTime always returns now
func (fi FI) ModTime() time.Time {
	return time.Time{}
}

// Mode returns the file as readonly and directories
// as world writeable and executable
func (fi FI) Mode() os.FileMode {
	if fi.IsDir() {
		return 0755 | os.ModeDir
	}
	return 0444
}

// Name returns the name of this object without
// any leading slashes
func (fi FI) Name() string {
	return path.Base(fi.name)
}

// Size returns the size of this item
func (fi FI) Size() int64 {
	return int64(fi.size)
}

// Sys returns nil
func (fi FI) Sys() interface{} {
	return nil
}
