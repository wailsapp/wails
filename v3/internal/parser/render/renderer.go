package render

import (
	"fmt"
	"io"
	"text/template"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/parser/collect"
)

// Renderer holds the template set for a given configuration.
// It provides methods for rendering various output modules.
type Renderer struct {
	options *flags.GenerateBindingsOptions

	bindings *template.Template
	ext      string

	models       *template.Template
	modelsFile   string
	internalFile string

	indexFile string
}

// NewRenderer initialises a renderer for the given configuration.
func NewRenderer(options *flags.GenerateBindingsOptions) *Renderer {
	ext := ".js"
	if options.TS {
		ext = ".ts"
	}

	return &Renderer{
		options: options,

		bindings: tmplBindings[tmplLanguage(options.TS)],
		ext:      ext,

		models:       tmplModels[tmplLanguage(options.TS)][tmplMode(options.UseInterfaces)],
		modelsFile:   "models" + ext,
		internalFile: "internal" + ext,

		indexFile: "index" + ext,
	}
}

// Bindings renders bindings for the given bound type to w.
func (renderer *Renderer) Bindings(w io.Writer, info *collect.BoundTypeInfo, collector *collect.Collector) error {
	return renderer.bindings.Execute(w, &struct {
		*collect.BoundTypeInfo
		*Renderer
		*flags.GenerateBindingsOptions
		Collector *collect.Collector
	}{
		info,
		renderer,
		renderer.options,
		collector,
	})
}

// BindingsFile returns the standard name of a bindings file
// for the given struct name, with the appropriate extension.
func (renderer *Renderer) BindingsFile(name string) string {
	return name + renderer.ext
}

// ModelsFile returns the standard name of a models file
// with the appropriate extension.
func (renderer *Renderer) ModelsFile() string {
	return renderer.modelsFile
}

// ModelsFile returns the standard name of an internal model file
// with the appropriate extension.
func (renderer *Renderer) InternalFile() string {
	return renderer.internalFile
}

// Index renders the given package index to w.
func (renderer *Renderer) Index(w io.Writer, index *collect.PackageIndex) error {
	return tmplIndex.Execute(w, &struct {
		*collect.PackageIndex
		*Renderer
		*flags.GenerateBindingsOptions
	}{
		index,
		renderer,
		renderer.options,
	})
}

// IndexFile returns the standard name of a package index file
// with the appropriate extension.
func (renderer *Renderer) IndexFile() string {
	return renderer.indexFile
}

// GlobalIndex renders the given import map as a global package index to w.
func (renderer *Renderer) GlobalIndex(w io.Writer, imports *collect.ImportMap) error {
	return tmplGlobalIndex.Execute(w, &struct {
		*collect.ImportMap
		*Renderer
		*flags.GenerateBindingsOptions
	}{
		imports,
		renderer,
		renderer.options,
	})
}

// GlobalIndex renders a shortcut file for the given import to w.
func (renderer *Renderer) Shortcut(w io.Writer, info collect.ImportInfo) error {
	return tmplShortcut.Execute(w, &struct {
		collect.ImportInfo
		*Renderer
		*flags.GenerateBindingsOptions
	}{
		info,
		renderer,
		renderer.options,
	})
}

// ShortcutFile returns the standard name of an import shortcut file
// with the appropriate extension.
func (renderer *Renderer) ShortcutFile(info collect.ImportInfo) string {
	if info.Index > 0 || info.Name == "index" {
		return fmt.Sprintf("%s.%d%s", info.Name, info.Index, renderer.ext)
	} else {
		return info.Name + renderer.ext
	}
}
