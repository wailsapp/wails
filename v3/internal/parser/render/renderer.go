package render

import (
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

	ext string

	service *template.Template
	models  *template.Template
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

		ext: ext,

		service: tmplService[tmplLanguage(options.TS)],
		models:  tmplModels[tmplLanguage(options.TS)],
	}
}

// ServiceFile returns the standard name of a service file
// for the given struct name, with the appropriate extension.
func (renderer *Renderer) ServiceFile(name string) string {
	return name + renderer.ext
}

// ModelsFile returns the standard name of a models file
// with the appropriate extension.
func (renderer *Renderer) ModelsFile() string {
	return renderer.options.ModelsFilename + renderer.ext
}

// InternalFile returns the standard name of an internal model file
// with the appropriate extension.
func (renderer *Renderer) InternalFile() string {
	return renderer.options.InternalFilename + renderer.ext
}

// IndexFile returns the standard name of a package index file
// with the appropriate extension.
func (renderer *Renderer) IndexFile() string {
	return renderer.options.IndexFilename + renderer.ext
}

// ShortcutFile returns the standard name of a package shortcut file
// with the appropriate extension.
func (renderer *Renderer) ShortcutFile(name string) string {
	return name + renderer.ext
}

// Service renders binding code for the given service type to w.
func (renderer *Renderer) Service(w io.Writer, info *collect.ServiceInfo) error {
	return renderer.service.Execute(w, &struct {
		module
		Service *collect.ServiceInfo
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

// GlobalIndex renders a shortcut file for the given package name to w.
func (renderer *Renderer) Shortcut(w io.Writer, imports *collect.ImportMap, name string) error {
	return tmplShortcut.Execute(w, &struct {
		*collect.ImportMap
		*Renderer
		*flags.GenerateBindingsOptions
		Name string
	}{
		imports,
		renderer,
		renderer.options,
		name,
	})
}
