package cmd

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
func (fs *FSHelper) RemoveFiles(files []string) error {
	for _, filename := range files {
		err := os.Remove(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetSubdirs will return a list of FQPs to subdirectories in the given directory
func (fs *FSHelper) GetSubdirs(dir string) (map[string]string, error) {

	// Read in the directory information
	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// Allocate space for the list
	subdirs := make(map[string]string)

	// Pull out the directories and store in the map as
	// map["directoryName"] = "path/to/directoryName"
	for _, file := range fileInfo {
		if file.IsDir() {
			subdirs[file.Name()] = filepath.Join(dir, file.Name())
		}
	}
	return subdirs, nil
}

// MkDir creates the given directory.
// Returns error on failure
func (fs *FSHelper) MkDir(dir string) error {
	return os.Mkdir(dir, 0700)
}

// LoadAsString will attempt to load the given file and return
// its contents as a string
func (fs *FSHelper) LoadAsString(filename string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	return string(bytes), err
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
