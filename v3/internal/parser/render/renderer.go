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
	options   *flags.GenerateBindingsOptions
	collector *collect.Collector

	bindings *template.Template
	ext      string

	models       *template.Template
	modelsFile   string
	internalFile string

	indexFile string
}

// NewRenderer initialises a code renderer
// for the given configuration and data collector.
func NewRenderer(options *flags.GenerateBindingsOptions, collector *collect.Collector) *Renderer {
	ext := ".js"
	if options.TS {
		ext = ".ts"
	}

	return &Renderer{
		options:   options,
		collector: collector,

		bindings: tmplBindings[tmplLanguage(options.TS)],
		ext:      ext,

		models:       tmplModels[tmplLanguage(options.TS)],
		modelsFile:   "models" + ext,
		internalFile: "internal" + ext,

		indexFile: "index" + ext,
	}
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

// InternalFile returns the standard name of an internal model file
// with the appropriate extension.
func (renderer *Renderer) InternalFile() string {
	return renderer.internalFile
}

// IndexFile returns the standard name of a package index file
// with the appropriate extension.
func (renderer *Renderer) IndexFile() string {
	return renderer.indexFile
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

// Bindings renders bindings code for the given bound type to w.
func (renderer *Renderer) Bindings(w io.Writer, info *collect.BoundTypeInfo) error {
	return renderer.bindings.Execute(w, &struct {
		module
		BoundType *collect.BoundTypeInfo
	}{
		module{
			Renderer:                renderer,
			GenerateBindingsOptions: renderer.options,
			Imports:                 info.Imports,
		},
		info,
	})
}

// Models renders models code for the given list of models.
func (renderer *Renderer) Models(w io.Writer, imports *collect.ImportMap, models []*collect.ModelInfo) error {
	return renderer.models.Execute(w, &struct {
		module
		Models []*collect.ModelInfo
	}{
		module{
			Renderer:                renderer,
			GenerateBindingsOptions: renderer.options,
			Imports:                 imports,
		},
		models,
	})
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
