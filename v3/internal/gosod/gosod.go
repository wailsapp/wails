// Package gosod extracts an fs.FS to a directory, processing .tmpl files as Go templates.
// Inlined from github.com/leaanthony/gosod (MIT licence).
package gosod

import (
	"bytes"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"
)

// hooks overridable in tests
var (
	absPath      = filepath.Abs
	osCreateFile = func(name string) (io.WriteCloser, error) { return os.Create(name) }
)

// TemplateDir processes an fs.FS and extracts it to a directory.
type TemplateDir struct {
	fs              fs.FS
	templateFilters []string
	dirs            []string
	standardFiles   []string
	templateFiles   []string
	ignoredFiles    map[string]struct{}
	renameFiles     map[string]string
}

// New creates a TemplateDir for the given filesystem.
func New(fsys fs.FS) *TemplateDir {
	return &TemplateDir{
		fs:              fsys,
		templateFilters: []string{".tmpl"},
		ignoredFiles:    make(map[string]struct{}),
		renameFiles:     make(map[string]string),
	}
}

// IgnoreFile adds a filename to the ignore list.
func (t *TemplateDir) IgnoreFile(filename string) {
	t.ignoredFiles[filename] = struct{}{}
}

// SetTemplateFilters sets the file suffixes treated as templates.
func (t *TemplateDir) SetTemplateFilters(filters []string) {
	t.templateFilters = filters
}

// RenameFiles sets a rename map applied to output file names.
func (t *TemplateDir) RenameFiles(renameFiles map[string]string) {
	t.renameFiles = renameFiles
}

// Extract processes the FS and writes files to targetDirectory using data as template input.
func (t *TemplateDir) Extract(targetDirectory string, data interface{}) error {
	abs, err := absPath(targetDirectory)
	if err != nil {
		return err
	}
	if _, err := os.Stat(abs); os.IsNotExist(err) {
		if err := os.MkdirAll(abs, 0755); err != nil {
			return err
		}
	}
	return t.processFiles(abs, data)
}

func (t *TemplateDir) processFiles(targetDir string, data interface{}) error {
	if err := t.categorise(); err != nil {
		return err
	}
	if err := t.createDirs(targetDir, data); err != nil {
		return err
	}
	if err := t.processTemplates(targetDir, data); err != nil {
		return err
	}
	return t.copyStandard(targetDir, data)
}

func (t *TemplateDir) categorise() error {
	t.dirs = nil
	t.standardFiles = nil
	t.templateFiles = nil
	return fs.WalkDir(t.fs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != "." {
				t.dirs = append(t.dirs, path)
			}
			return nil
		}
		name := filepath.Base(path)
		if _, ignored := t.ignoredFiles[name]; ignored {
			return nil
		}
		for _, f := range t.templateFilters {
			if strings.Contains(name, f) {
				t.templateFiles = append(t.templateFiles, path)
				return nil
			}
		}
		t.standardFiles = append(t.standardFiles, path)
		return nil
	})
}

func (t *TemplateDir) resolveTarget(path, targetDir string, data any) string {
	result := filepath.Join(targetDir, path)
	if data == nil {
		return result
	}
	tmpl, err := template.New("p").Parse(result)
	if err != nil {
		return result
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return result
	}
	return buf.String()
}

func (t *TemplateDir) createDirs(targetDir string, data any) error {
	for _, d := range t.dirs {
		target := t.resolveTarget(d, targetDir, data)
		if err := os.MkdirAll(target, 0755); err != nil && err != syscall.EEXIST {
			return err
		}
	}
	return nil
}

func (t *TemplateDir) processTemplates(targetDir string, data interface{}) error {
	for _, srcPath := range t.templateFiles {
		tmpl, err := template.ParseFS(t.fs, srcPath)
		if err != nil {
			return err
		}
		target := t.resolveTarget(srcPath, targetDir, data)
		baseDir := filepath.Dir(target)
		name := filepath.Base(target)
		for _, f := range t.templateFilters {
			name = strings.ReplaceAll(name, f, "")
		}
		if r := t.renameFiles[name]; r != "" {
			name = r
		}
		target = filepath.Join(baseDir, name)
		w, err := osCreateFile(target)
		if err != nil {
			return err
		}
		if err := tmpl.Execute(w, data); err != nil {
			_ = w.Close()
			return err
		}
		if err := w.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (t *TemplateDir) copyStandard(targetDir string, data any) error {
	for _, srcPath := range t.standardFiles {
		name := filepath.Base(srcPath)
		if r := t.renameFiles[name]; r != "" {
			name = r
		}
		target := t.resolveTarget(filepath.Join(filepath.Dir(srcPath), name), targetDir, data)
		if err := t.copyFile(srcPath, target); err != nil {
			return err
		}
	}
	return nil
}

func (t *TemplateDir) copyFile(src, dst string) error {
	s, err := t.fs.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := s.Close(); err != nil {
			log.Println("gosod: close source:", err)
		}
	}()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		_ = d.Close()
		return err
	}
	return d.Close()
}
