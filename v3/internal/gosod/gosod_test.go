package gosod

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
	"time"
)

// ---- helpers ----

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}

// ---- Basic extraction ----

func TestExtract_StandardFile(t *testing.T) {
	fsys := fstest.MapFS{
		"hello.txt": {Data: []byte("hello world")},
	}
	td := New(fsys)
	dir := t.TempDir()
	if err := td.Extract(dir, nil); err != nil {
		t.Fatal(err)
	}
	got := readFile(t, filepath.Join(dir, "hello.txt"))
	if got != "hello world" {
		t.Errorf("expected 'hello world', got %q", got)
	}
}

// ---- Template file (.tmpl suffix removed, content rendered) ----

func TestExtract_TemplateFile(t *testing.T) {
	fsys := fstest.MapFS{
		"greeting.txt.tmpl": {Data: []byte("Hello, {{.Name}}!")},
	}
	td := New(fsys)
	dir := t.TempDir()
	data := struct{ Name string }{"World"}
	if err := td.Extract(dir, data); err != nil {
		t.Fatal(err)
	}
	// .tmpl suffix removed
	got := readFile(t, filepath.Join(dir, "greeting.txt"))
	if got != "Hello, World!" {
		t.Errorf("expected 'Hello, World!', got %q", got)
	}
}

// ---- Subdirectory creation ----

func TestExtract_Subdirectory(t *testing.T) {
	fsys := fstest.MapFS{
		"sub/file.txt": {Data: []byte("subfile"), Mode: 0644},
	}
	td := New(fsys)
	dir := t.TempDir()
	if err := td.Extract(dir, nil); err != nil {
		t.Fatal(err)
	}
	got := readFile(t, filepath.Join(dir, "sub", "file.txt"))
	if got != "subfile" {
		t.Errorf("expected 'subfile', got %q", got)
	}
}

// ---- IgnoreFile ----

func TestExtract_IgnoreFile(t *testing.T) {
	fsys := fstest.MapFS{
		"keep.txt":   {Data: []byte("keep")},
		"ignore.txt": {Data: []byte("should not appear")},
	}
	td := New(fsys)
	td.IgnoreFile("ignore.txt")
	dir := t.TempDir()
	if err := td.Extract(dir, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "ignore.txt")); !os.IsNotExist(err) {
		t.Error("ignored file should not be extracted")
	}
	if _, err := os.Stat(filepath.Join(dir, "keep.txt")); err != nil {
		t.Error("non-ignored file should be extracted")
	}
}

// ---- SetTemplateFilters ----

func TestExtract_SetTemplateFilters(t *testing.T) {
	fsys := fstest.MapFS{
		"config.go.tpl": {Data: []byte("package {{.Pkg}}")},
	}
	td := New(fsys)
	td.SetTemplateFilters([]string{".tpl"})
	dir := t.TempDir()
	data := struct{ Pkg string }{"main"}
	if err := td.Extract(dir, data); err != nil {
		t.Fatal(err)
	}
	// .tpl suffix should be removed
	got := readFile(t, filepath.Join(dir, "config.go"))
	if got != "package main" {
		t.Errorf("expected 'package main', got %q", got)
	}
}

// ---- RenameFiles (standard file) ----

func TestExtract_RenameStandardFile(t *testing.T) {
	fsys := fstest.MapFS{
		"original.txt": {Data: []byte("renamed")},
	}
	td := New(fsys)
	td.RenameFiles(map[string]string{"original.txt": "new.txt"})
	dir := t.TempDir()
	if err := td.Extract(dir, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "original.txt")); !os.IsNotExist(err) {
		t.Error("original filename should not exist")
	}
	got := readFile(t, filepath.Join(dir, "new.txt"))
	if got != "renamed" {
		t.Errorf("expected 'renamed', got %q", got)
	}
}

// ---- RenameFiles (template file) ----

func TestExtract_RenameTemplateFile(t *testing.T) {
	fsys := fstest.MapFS{
		"src.txt.tmpl": {Data: []byte("content")},
	}
	td := New(fsys)
	// After .tmpl is stripped, name is "src.txt"; rename that
	td.RenameFiles(map[string]string{"src.txt": "dest.txt"})
	dir := t.TempDir()
	if err := td.Extract(dir, nil); err != nil {
		t.Fatal(err)
	}
	got := readFile(t, filepath.Join(dir, "dest.txt"))
	if got != "content" {
		t.Errorf("expected 'content', got %q", got)
	}
}

// ---- Path templating: directory name contains {{.Field}} ----

func TestExtract_PathTemplating(t *testing.T) {
	fsys := fstest.MapFS{
		"{{.AppName}}/main.go": {Data: []byte("package main")},
	}
	td := New(fsys)
	dir := t.TempDir()
	data := struct{ AppName string }{"myapp"}
	if err := td.Extract(dir, data); err != nil {
		t.Fatal(err)
	}
	got := readFile(t, filepath.Join(dir, "myapp", "main.go"))
	if got != "package main" {
		t.Errorf("expected 'package main', got %q", got)
	}
}

// ---- Target dir created if not exists ----

func TestExtract_CreatesTargetDir(t *testing.T) {
	fsys := fstest.MapFS{
		"f.txt": {Data: []byte("x")},
	}
	td := New(fsys)
	dir := filepath.Join(t.TempDir(), "newdir", "nested")
	if err := td.Extract(dir, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("target dir not created: %v", err)
	}
}

// ---- Empty FS: no error ----

func TestExtract_EmptyFS(t *testing.T) {
	fsys := fstest.MapFS{}
	td := New(fsys)
	dir := t.TempDir()
	if err := td.Extract(dir, nil); err != nil {
		t.Fatalf("expected no error for empty FS, got: %v", err)
	}
}

// ---- resolveTarget with nil data (no template expansion) ----

func TestExtract_NilData_NoTemplateExpansion(t *testing.T) {
	fsys := fstest.MapFS{
		"file.txt": {Data: []byte("content")},
	}
	td := New(fsys)
	dir := t.TempDir()
	if err := td.Extract(dir, nil); err != nil {
		t.Fatal(err)
	}
}

// ---- Error: processTemplates parse error (bad template syntax) ----

func TestExtract_TemplateParseFails(t *testing.T) {
	fsys := fstest.MapFS{
		"bad.txt.tmpl": {Data: []byte("{{unterminated")},
	}
	td := New(fsys)
	dir := t.TempDir()
	err := td.Extract(dir, nil)
	if err == nil {
		t.Fatal("expected error for bad template syntax")
	}
}

// ---- Error: template Execute fails (accessing missing field) ----

func TestExtract_TemplateExecuteFails(t *testing.T) {
	fsys := fstest.MapFS{
		// Option: a template that calls a nil function
		"fail.txt.tmpl": {Data: []byte("{{call .Nil}}")},
	}
	td := New(fsys)
	dir := t.TempDir()
	data := struct{ Nil func() }{nil}
	err := td.Extract(dir, data)
	if err == nil {
		t.Fatal("expected error when template Execute fails")
	}
}

// ---- Error: os.Create fails for template output (target is a directory) ----

func TestExtract_TemplateCreateFails(t *testing.T) {
	fsys := fstest.MapFS{
		"out.txt.tmpl": {Data: []byte("hello")},
	}
	td := New(fsys)
	dir := t.TempDir()
	// Pre-create a directory where the output file would go
	if err := os.MkdirAll(filepath.Join(dir, "out.txt"), 0755); err != nil {
		t.Fatal(err)
	}
	err := td.Extract(dir, nil)
	if err == nil {
		t.Fatal("expected error when target path is a directory")
	}
}

// ---- Error: createDirs MkdirAll fails (target path blocked by file) ----

func TestExtract_CreateDirsFails(t *testing.T) {
	fsys := fstest.MapFS{
		"subdir/file.txt": {Data: []byte("x")},
	}
	td := New(fsys)
	dir := t.TempDir()
	// Block "subdir" from being a directory by creating it as a file
	if err := os.WriteFile(filepath.Join(dir, "subdir"), []byte("blocked"), 0644); err != nil {
		t.Fatal(err)
	}
	err := td.Extract(dir, nil)
	if err == nil {
		t.Fatal("expected error when dir path blocked by file")
	}
}

// ---- Error: copyFile Open fails ----

type openErrFS struct {
	inner   fs.FS
	failFor string
}

func (f *openErrFS) Open(name string) (fs.File, error) {
	if name == f.failFor {
		return nil, os.ErrPermission
	}
	return f.inner.Open(name)
}

func TestExtract_CopyFileOpenFails(t *testing.T) {
	inner := fstest.MapFS{
		"secret.txt": {Data: []byte("top secret")},
	}
	fsys := &openErrFS{inner: inner, failFor: "secret.txt"}
	td := New(fsys)
	dir := t.TempDir()
	err := td.Extract(dir, nil)
	if err == nil {
		t.Fatal("expected error when source file cannot be opened")
	}
}

// ---- Error: copyFile os.Create fails (target is a directory) ----

func TestExtract_CopyFileCreateFails(t *testing.T) {
	fsys := fstest.MapFS{
		"out.txt": {Data: []byte("content")},
	}
	td := New(fsys)
	dir := t.TempDir()
	// Pre-create a directory where "out.txt" would land
	if err := os.MkdirAll(filepath.Join(dir, "out.txt"), 0755); err != nil {
		t.Fatal(err)
	}
	err := td.Extract(dir, nil)
	if err == nil {
		t.Fatal("expected error when target path is a directory")
	}
}

// ---- Error: copyFile io.Copy fails (custom failing reader) ----

type badReadFile struct {
	fs.File
}

func (b *badReadFile) Read(_ []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

type badReadFS struct {
	inner   fs.FS
	failFor string
}

func (f *badReadFS) Open(name string) (fs.File, error) {
	file, err := f.inner.Open(name)
	if err != nil {
		return nil, err
	}
	if name == f.failFor {
		return &badReadFile{File: file}, nil
	}
	return file, nil
}

func TestExtract_CopyFileReadFails(t *testing.T) {
	inner := fstest.MapFS{
		"data.txt": {Data: []byte("data")},
	}
	fsys := &badReadFS{inner: inner, failFor: "data.txt"}
	td := New(fsys)
	dir := t.TempDir()
	err := td.Extract(dir, nil)
	if err == nil {
		t.Fatal("expected error when reading source file fails")
	}
}

// ---- Error: copyFile close logs (source Close fails) ----
// This exercises the log.Println path in copyFile's deferred close.
// We can't assert the error (it's swallowed), but we exercise the branch.

type badCloseFile struct {
	fs.File
}

func (b *badCloseFile) Close() error {
	return io.ErrUnexpectedEOF
}

type badCloseFS struct {
	inner   fs.FS
	failFor string
}

func (f *badCloseFS) Open(name string) (fs.File, error) {
	file, err := f.inner.Open(name)
	if err != nil {
		return nil, err
	}
	if name == f.failFor {
		return &badCloseFile{File: file}, nil
	}
	return file, nil
}

func TestExtract_CopyFileCloseFails(t *testing.T) {
	inner := fstest.MapFS{
		"data.txt": {Data: []byte("data")},
	}
	fsys := &badCloseFS{inner: inner, failFor: "data.txt"}
	td := New(fsys)
	dir := t.TempDir()
	// Should not return error (close error is logged, not returned)
	if err := td.Extract(dir, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}


// ---- resolveTarget: template parse error returns original path ----

func TestResolveTarget_TemplateParseError(t *testing.T) {
	fsys := fstest.MapFS{}
	td := New(fsys)
	// Path with invalid template syntax — resolveTarget should return the raw joined path.
	// Must pass non-nil data so we don't short-circuit at the nil check.
	result := td.resolveTarget("{{invalid", "/base", struct{}{})
	if result == "" {
		t.Error("expected non-empty result")
	}
}

// ---- resolveTarget: template execute error returns original path ----

func TestResolveTarget_TemplateExecuteError(t *testing.T) {
	fsys := fstest.MapFS{}
	td := New(fsys)
	// Template that calls a nil function — execute will fail
	data := struct{ Nil func() }{nil}
	result := td.resolveTarget("{{call .Nil}}", "/base", data)
	if result == "" {
		t.Error("expected non-empty result")
	}
}

// ---- Multiple template filters: ensure all are stripped ----

func TestExtract_MultipleTemplateFilters(t *testing.T) {
	fsys := fstest.MapFS{
		"a.txt.tmpl": {Data: []byte("A")},
		"b.go.tpl":   {Data: []byte("B")},
	}
	td := New(fsys)
	td.SetTemplateFilters([]string{".tmpl", ".tpl"})
	dir := t.TempDir()
	if err := td.Extract(dir, nil); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, filepath.Join(dir, "a.txt")); got != "A" {
		t.Errorf("a.txt: expected 'A', got %q", got)
	}
	if got := readFile(t, filepath.Join(dir, "b.go")); got != "B" {
		t.Errorf("b.go: expected 'B', got %q", got)
	}
}

// ---- WalkDir error: custom FS that returns an error during walk ----

type walkErrFS struct {
	inner fs.FS
}

// walkErrDirEntry wraps a real DirEntry to return errors when its dir is read.
type walkErrDir struct {
	fs.File
	entries []fs.DirEntry
}

func (w *walkErrDir) ReadDir(n int) ([]fs.DirEntry, error) {
	return nil, io.ErrClosedPipe
}

func (f *walkErrFS) Open(name string) (fs.File, error) {
	file, err := f.inner.Open(name)
	if err != nil {
		return nil, err
	}
	// Wrap the root "." dir to fail ReadDir
	if name == "." {
		info, err := file.Stat()
		if err != nil {
			_ = file.Close()
			return nil, err
		}
		if info.IsDir() {
			return &walkErrDir{File: file}, nil
		}
	}
	return file, nil
}

func TestExtract_WalkDirError(t *testing.T) {
	inner := fstest.MapFS{
		"file.txt": {Data: []byte("content")},
	}
	fsys := &walkErrFS{inner: inner}
	td := New(fsys)
	dir := t.TempDir()
	err := td.Extract(dir, nil)
	if err == nil {
		t.Fatal("expected error when WalkDir fails")
	}
}

// ---- Ensure we cover the case where WalkDir callback receives a non-nil err ----
// We need to simulate a WalkDir error callback. We do this with a custom FS
// that makes WalkDir pass an error into the callback.

type callbackErrFS struct{}

func (f *callbackErrFS) Open(name string) (fs.File, error) {
	if name == "." {
		return &callbackErrDir{}, nil
	}
	return nil, os.ErrNotExist
}

type callbackErrDir struct{}

func (d *callbackErrDir) Stat() (fs.FileInfo, error) {
	return &callbackErrFileInfo{}, nil
}
func (d *callbackErrDir) Read([]byte) (int, error) { return 0, io.EOF }
func (d *callbackErrDir) Close() error             { return nil }
func (d *callbackErrDir) ReadDir(n int) ([]fs.DirEntry, error) {
	// Return an entry that will fail Info()
	return []fs.DirEntry{&callbackErrEntry{}}, nil
}

type callbackErrFileInfo struct{}

func (i *callbackErrFileInfo) Name() string      { return "." }
func (i *callbackErrFileInfo) Size() int64       { return 0 }
func (i *callbackErrFileInfo) Mode() fs.FileMode { return fs.ModeDir | 0755 }
func (i *callbackErrFileInfo) ModTime() time.Time { return time.Time{} }
func (i *callbackErrFileInfo) IsDir() bool       { return true }
func (i *callbackErrFileInfo) Sys() interface{}  { return nil }

type callbackErrEntry struct{}

func (e *callbackErrEntry) Name() string               { return "fail.txt" }
func (e *callbackErrEntry) IsDir() bool                { return false }
func (e *callbackErrEntry) Type() fs.FileMode          { return 0 }
func (e *callbackErrEntry) Info() (fs.FileInfo, error) { return nil, io.ErrUnexpectedEOF }

func TestExtract_CategoriseCallbackError(t *testing.T) {
	td := New(&callbackErrFS{})
	dir := t.TempDir()
	err := td.Extract(dir, nil)
	if err == nil {
		t.Fatal("expected error from WalkDir callback")
	}
}

// ---- hook overrides for unreachable OS-level error paths ----

// TestExtract_AbsPathError covers the filepath.Abs error branch by swapping
// the absPath hook with a failing stub.
func TestExtract_AbsPathError(t *testing.T) {
	orig := absPath
	absPath = func(_ string) (string, error) { return "", io.ErrUnexpectedEOF }
	t.Cleanup(func() { absPath = orig })

	td := New(fstest.MapFS{})
	err := td.Extract("anything", nil)
	if err != io.ErrUnexpectedEOF {
		t.Fatalf("expected ErrUnexpectedEOF, got %v", err)
	}
}

// failingWriteCloser is an io.WriteCloser whose Close always returns an error.
type failingWriteCloser struct{ io.Writer }

func (f failingWriteCloser) Close() error { return io.ErrClosedPipe }

// TestProcessTemplates_CloseError covers the w.Close() error branch by
// swapping the osCreateFile hook with a stub that returns a failing closer.
func TestProcessTemplates_CloseError(t *testing.T) {
	orig := osCreateFile
	osCreateFile = func(_ string) (io.WriteCloser, error) {
		return failingWriteCloser{Writer: io.Discard}, nil
	}
	t.Cleanup(func() { osCreateFile = orig })

	fsys := fstest.MapFS{
		"out.tmpl": {Data: []byte("hello")},
	}
	td := New(fsys)
	err := td.Extract(t.TempDir(), nil)
	if err != io.ErrClosedPipe {
		t.Fatalf("expected ErrClosedPipe from w.Close(), got %v", err)
	}
}
