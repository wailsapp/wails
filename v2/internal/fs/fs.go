package fs

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"unsafe"

	"github.com/leaanthony/slicer"
)

// LocalDirectory gets the caller's file directory
// Equivalent to node's __DIRNAME
func LocalDirectory() string {
	_, thisFile, _, _ := runtime.Caller(1)
	return filepath.Dir(thisFile)
}

// RelativeToCwd returns an absolute path based on the cwd
// and the given relative path
func RelativeToCwd(relativePath string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(cwd, relativePath), nil
}

// Mkdir will create the given directory
func Mkdir(dirname string) error {
	return os.Mkdir(dirname, 0755)
}

// MkDirs creates the given nested directories.
// Returns error on failure
func MkDirs(fullPath string, mode ...os.FileMode) error {
	var perms os.FileMode
	perms = 0755
	if len(mode) == 1 {
		perms = mode[0]
	}
	return os.MkdirAll(fullPath, perms)
}

// MoveFile attempts to move the source file to the target
// Target is a fully qualified path to a file *name*, not a
// directory
func MoveFile(source string, target string) error {
	return os.Rename(source, target)
}

// DeleteFile will delete the given file
func DeleteFile(filename string) error {
	return os.Remove(filename)
}

// CopyFile from source to target
func CopyFile(source string, target string) error {
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

// DirExists - Returns true if the given path resolves to a directory on the filesystem
func DirExists(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}

	return fi.Mode().IsDir()
}

// FileExists returns a boolean value indicating whether
// the given file exists
func FileExists(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}

	return fi.Mode().IsRegular()
}

// RelativePath returns a qualified path created by joining the
// directory of the calling file and the given relative path.
//
// Example: RelativePath("..") in *this* file would give you '/path/to/wails2/v2/internal`
func RelativePath(relativepath string, optionalpaths ...string) string {
	_, thisFile, _, _ := runtime.Caller(1)
	localDir := filepath.Dir(thisFile)

	// If we have optional paths, join them to the relativepath
	if len(optionalpaths) > 0 {
		paths := []string{relativepath}
		paths = append(paths, optionalpaths...)
		relativepath = filepath.Join(paths...)
	}
	result, err := filepath.Abs(filepath.Join(localDir, relativepath))
	if err != nil {
		// I'm allowing this for 1 reason only: It's fatal if the path
		// supplied is wrong as it's only used internally in Wails. If we get
		// that path wrong, we should know about it immediately. The other reason is
		// that it cuts down a ton of unnecassary error handling.
		panic(err)
	}
	return result
}

// MustLoadString attempts to load a string and will abort with a fatal message if
// something goes wrong
func MustLoadString(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("FATAL: Unable to load file '%s': %s\n", filename, err.Error())
		os.Exit(1)
	}
	return *(*string)(unsafe.Pointer(&data))
}

// MD5File returns the md5sum of the given file
func MD5File(filename string) (string, error) {
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

// MustMD5File will call MD5File and abort the program on error
func MustMD5File(filename string) string {
	result, err := MD5File(filename)
	if err != nil {
		println("FATAL: Unable to MD5Sum file:", err.Error())
		os.Exit(1)
	}
	return result
}

// MustWriteString will attempt to write the given data to the given filename
// It will abort the program in the event of a failure
func MustWriteString(filename string, data string) {
	err := ioutil.WriteFile(filename, []byte(data), 0755)
	if err != nil {
		fatal("Unable to write file", filename, ":", err.Error())
		os.Exit(1)
	}
}

// fatal will print the optional messages and die
func fatal(message ...string) {
	if len(message) > 0 {
		print("FATAL:")
		for text := range message {
			print(text)
		}
	}
	os.Exit(1)
}

// GetSubdirectories returns a list of subdirectories for the given root directory
func GetSubdirectories(rootDir string) (*slicer.StringSlicer, error) {
	var result slicer.StringSlicer

	// Iterate root dir
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// If we have a directory, save it
		if info.IsDir() {
			result.Add(path)
		}
		return nil
	})
	return &result, err
}

func DirIsEmpty(dir string) (bool, error) {

	// CREDIT: https://stackoverflow.com/a/30708914/8325411
	f, err := os.Open(dir)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

// Credit: https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04
// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = MkDirs(dst)
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}
