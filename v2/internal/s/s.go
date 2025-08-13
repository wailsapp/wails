package s

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitfield/script"
)

var (
	Output         io.Writer = io.Discard
	IndentSize     int
	originalOutput io.Writer
	currentIndent  int
)

func checkError(err error) {
	if err != nil {
		println("\nERROR:", err.Error())
		os.Exit(1)
	}
}

func mute() {
	originalOutput = Output
	Output = ioutil.Discard
}

func unmute() {
	Output = originalOutput
}

func indent() {
	currentIndent += IndentSize
}

func unindent() {
	currentIndent -= IndentSize
}

func log(message string, args ...interface{}) {
	indent := strings.Repeat(" ", currentIndent)
	_, err := fmt.Fprintf(Output, indent+message+"\n", args...)
	checkError(err)
}

// RENAME a file or directory
func RENAME(source string, target string) {
	log("RENAME %s -> %s", source, target)
	err := os.Rename(source, target)
	checkError(err)
}

// MUSTDELETE a file.
func MUSTDELETE(filename string) {
	log("DELETE %s", filename)
	err := os.Remove(filepath.Join(CWD(), filename))
	checkError(err)
}

// DELETE a file.
func DELETE(filename string) {
	log("DELETE %s", filename)
	_ = os.Remove(filepath.Join(CWD(), filename))
}

func CD(dir string) {
	err := os.Chdir(dir)
	checkError(err)
	log("CD %s [%s]", dir, CWD())
}

func MKDIR(path string, mode ...os.FileMode) {
	var perms os.FileMode
	perms = 0o755
	if len(mode) == 1 {
		perms = mode[0]
	}
	log("MKDIR %s (perms: %v)", path, perms)
	err := os.MkdirAll(path, perms)
	checkError(err)
}

// ENDIR ensures that the path gets created if it doesn't exist
func ENDIR(path string, mode ...os.FileMode) {
	var perms os.FileMode
	perms = 0o755
	if len(mode) == 1 {
		perms = mode[0]
	}
	_ = os.MkdirAll(path, perms)
}

// COPYDIR recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
// Credit: https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04
func COPYDIR(src string, dst string) {
	log("COPYDIR %s -> %s", src, dst)
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	checkError(err)
	if !si.IsDir() {
		checkError(fmt.Errorf("source is not a directory"))
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		checkError(err)
	}
	if err == nil {
		checkError(fmt.Errorf("destination already exists"))
	}

	indent()
	MKDIR(dst)

	entries, err := os.ReadDir(src)
	checkError(err)

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			COPYDIR(srcPath, dstPath)
		} else {
			// Skip symlinks.
			if entry.Type()&os.ModeSymlink != 0 {
				continue
			}

			COPY(srcPath, dstPath)
		}
	}
	unindent()
}

// COPY file from source to target
func COPY(source string, target string) {
	log("COPY %s -> %s", source, target)
	src, err := os.Open(source)
	checkError(err)
	defer closefile(src)
	d, err := os.Create(target)
	checkError(err)
	defer closefile(d)
	_, err = io.Copy(d, src)
	checkError(err)
}

func CWD() string {
	result, err := os.Getwd()
	checkError(err)
	log("CWD [%s]", result)
	return result
}

func RMDIR(target string) {
	log("RMDIR %s", target)
	err := os.RemoveAll(target)
	checkError(err)
}

func RM(target string) {
	log("RM %s", target)
	err := os.Remove(target)
	checkError(err)
}

func ECHO(message string) {
	println(message)
}

func TOUCH(filepath string) {
	log("TOUCH %s", filepath)
	f, err := os.Create(filepath)
	checkError(err)
	closefile(f)
}

func EXEC(command string) {
	log("EXEC %s", command)
	gen := script.Exec(command)
	gen.Wait()
	checkError(gen.Error())
}

// EXISTS - Returns true if the given path exists
func EXISTS(path string) bool {
	_, err := os.Lstat(path)
	log("EXISTS %s (%T)", path, err == nil)
	return err == nil
}

// ISDIR returns true if the given directory exists
func ISDIR(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}

	return fi.Mode().IsDir()
}

// ISDIREMPTY returns true if the given directory is empty
func ISDIREMPTY(dir string) bool {
	// CREDIT: https://stackoverflow.com/a/30708914/8325411
	f, err := os.Open(dir)
	checkError(err)
	defer closefile(f)

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	return err == io.EOF
}

// ISFILE returns true if the given file exists
func ISFILE(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}

	return fi.Mode().IsRegular()
}

// SUBDIRS returns a list of subdirectories for the given directory
func SUBDIRS(rootDir string) []string {
	var result []string

	// Iterate root dir
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		checkError(err)
		// If we have a directory, save it
		if info.IsDir() {
			result = append(result, path)
		}
		return nil
	})
	checkError(err)
	return result
}

// SAVESTRING will create a file with the given string
func SAVESTRING(filename string, data string) {
	log("SAVESTRING %s", filename)
	mute()
	SAVEBYTES(filename, []byte(data))
	unmute()
}

// LOADSTRING returns the contents of the given filename as a string
func LOADSTRING(filename string) string {
	log("LOADSTRING %s", filename)
	mute()
	data := LOADBYTES(filename)
	unmute()
	return string(data)
}

// SAVEBYTES will create a file with the given string
func SAVEBYTES(filename string, data []byte) {
	log("SAVEBYTES %s", filename)
	err := os.WriteFile(filename, data, 0o755)
	checkError(err)
}

// LOADBYTES returns the contents of the given filename as a string
func LOADBYTES(filename string) []byte {
	log("LOADBYTES %s", filename)
	data, err := os.ReadFile(filename)
	checkError(err)
	return data
}

func closefile(f *os.File) {
	err := f.Close()
	checkError(err)
}

// MD5FILE returns the md5sum of the given file
func MD5FILE(filename string) string {
	f, err := os.Open(filename)
	checkError(err)
	defer closefile(f)

	h := md5.New()
	_, err = io.Copy(h, f)
	checkError(err)

	return hex.EncodeToString(h.Sum(nil))
}

// Sub is the substitution type
type Sub map[string]string

// REPLACEALL replaces all substitution keys with associated values in the given file
func REPLACEALL(filename string, substitutions Sub) {
	log("REPLACEALL %s (%v)", filename, substitutions)
	data := LOADSTRING(filename)
	for old, newText := range substitutions {
		data = strings.ReplaceAll(data, old, newText)
	}
	SAVESTRING(filename, data)
}
