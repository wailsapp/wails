package cmd

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/leaanthony/slicer"
)

// FSHelper - Wrapper struct for File System utility commands
type FSHelper struct {
}

// NewFSHelper - Returns a new FSHelper
func NewFSHelper() *FSHelper {
	result := &FSHelper{}
	return result
}

// DirExists - Returns true if the given path resolves to a directory on the filesystem
func (fs *FSHelper) DirExists(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}

	return fi.Mode().IsDir()
}

// FileExists returns a boolean value indicating whether
// the given file exists
func (fs *FSHelper) FileExists(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}

	return fi.Mode().IsRegular()
}

// FindFile returns the first occurrence of match inside path.
func (fs *FSHelper) FindFile(path, match string) (string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}

	for _, f := range files {
		if !f.IsDir() && strings.Contains(f.Name(), match) {
			return f.Name(), nil
		}
	}

	return "", fmt.Errorf("file not found")
}

// CreateFile creates a file at the given filename location with the contents
// set to the given data. It will create intermediary directories if needed.
func (fs *FSHelper) CreateFile(filename string, data []byte) error {
	// Ensure directory exists
	fs.MkDirs(filepath.Dir(filename))
	return ioutil.WriteFile(filename, data, 0644)
}

// MkDirs creates the given nested directories.
// Returns error on failure
func (fs *FSHelper) MkDirs(fullPath string, mode ...os.FileMode) error {
	var perms os.FileMode
	perms = 0700
	if len(mode) == 1 {
		perms = mode[0]
	}
	return os.MkdirAll(fullPath, perms)
}

// CopyFile from source to target
func (fs *FSHelper) CopyFile(source, target string) error {
	s, err := os.Open(source)
	if err != nil {
		return err
	}
	defer s.Close()
	d, err := os.Create(target)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}

// Cwd returns the current working directory
// Aborts on Failure
func (fs *FSHelper) Cwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to get working directory!")
	}
	return cwd
}

// RemoveFile removes the given filename
func (fs *FSHelper) RemoveFile(filename string) error {
	return os.Remove(filename)
}

// RemoveFiles removes the given filenames
func (fs *FSHelper) RemoveFiles(files []string, continueOnError bool) error {
	for _, filename := range files {
		err := os.Remove(filename)
		if err != nil && !continueOnError {
			return err
		}
	}
	return nil
}

// Dir holds information about a directory
type Dir struct {
	localPath string
	fullPath  string
}

// Directory creates a new Dir struct with the given directory path
func (fs *FSHelper) Directory(dir string) (*Dir, error) {
	fullPath, err := filepath.Abs(dir)
	return &Dir{fullPath: fullPath}, err
}

// LocalDir creates a new Dir struct based on a path relative to the caller
func (fs *FSHelper) LocalDir(dir string) (*Dir, error) {
	_, filename, _, _ := runtime.Caller(1)
	fullPath, err := filepath.Abs(filepath.Join(path.Dir(filename), dir))
	return &Dir{
		localPath: dir,
		fullPath:  fullPath,
	}, err
}

// LoadRelativeFile loads the given file relative to the caller's directory
func (fs *FSHelper) LoadRelativeFile(relativePath string) ([]byte, error) {
	_, filename, _, _ := runtime.Caller(0)
	fullPath, err := filepath.Abs(filepath.Join(path.Dir(filename), relativePath))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(fullPath)
}

// GetSubdirs will return a list of FQPs to subdirectories in the given directory
func (d *Dir) GetSubdirs() (map[string]string, error) {

	// Read in the directory information
	fileInfo, err := ioutil.ReadDir(d.fullPath)
	if err != nil {
		return nil, err
	}

	// Allocate space for the list
	subdirs := make(map[string]string)

	// Pull out the directories and store in the map as
	// map["directoryName"] = "path/to/directoryName"
	for _, file := range fileInfo {
		if file.IsDir() {
			subdirs[file.Name()] = filepath.Join(d.fullPath, file.Name())
		}
	}
	return subdirs, nil
}

// GetAllFilenames returns all filename in and below this directory
func (d *Dir) GetAllFilenames() (*slicer.StringSlicer, error) {
	result := slicer.String()
	err := filepath.Walk(d.fullPath, func(dir string, info os.FileInfo, err error) error {
		if dir == d.fullPath {
			return nil
		}
		if err != nil {
			return err
		}

		// Don't copy template metadata
		result.Add(dir)

		return nil
	})
	return result, err
}

// MkDir creates the given directory.
// Returns error on failure
func (fs *FSHelper) MkDir(dir string) error {
	return os.Mkdir(dir, 0700)
}

// SaveAsJSON saves the JSON representation of the given data to the given filename
func (fs *FSHelper) SaveAsJSON(data interface{}, filename string) error {

	var buf bytes.Buffer
	e := json.NewEncoder(&buf)
	e.SetEscapeHTML(false)
	e.SetIndent("", "  ")
	e.Encode(data)

	err := ioutil.WriteFile(filename, buf.Bytes(), 0755)
	if err != nil {
		return err
	}
	return nil
}

// LoadAsString will attempt to load the given file and return
// its contents as a string
func (fs *FSHelper) LoadAsString(filename string) (string, error) {
	bytes, err := fs.LoadAsBytes(filename)
	return string(bytes), err
}

// LoadAsBytes returns the contents of the file as a byte slice
func (fs *FSHelper) LoadAsBytes(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

// FileMD5 returns the md5sum of the given file
func (fs *FSHelper) FileMD5(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
