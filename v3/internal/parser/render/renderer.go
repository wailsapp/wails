package render

import (
	"io"
	"text/template"

	"github.com/wailsapp/wails/v3/internal/flags"
)

// Renderer holds the template set for a given configuration.
// It provides methods for rendering various output modules.
type Renderer struct {
	bindings *template.Template
	ext      string

	models       *template.Template
	modelsFile   string
	internalFile string

	index     *template.Template
	indexFile string
}

// NewRenderer initialises a renderer for the given configuration.
func NewRenderer(options *flags.GenerateBindingsOptions) Renderer {
	ext := ".js"
	if options.TS {
		ext = ".ts"
	}

	return Renderer{
		bindings: tmplBindings[tmplLanguage(options.TS)],
		ext:      ext,

		models:       tmplModels[tmplLanguage(options.TS)][tmplMode(options.UseInterfaces)],
		modelsFile:   "models" + ext,
		internalFile: "internal" + ext,

		index:     tmplIndex[tmplLanguage(options.TS)],
		indexFile: "index" + ext,
	}
}

// Bindings renders the given binding data to w.
func (renderer *Renderer) Bindings(w io.Writer, data any) error {
	return renderer.bindings.Execute(w, data)
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

// Index renders the given package index data to w.
func (renderer *Renderer) Index(w io.Writer, data any) error {
	return renderer.index.Execute(w, data)
}

// IndexFile returns the standard name of a package index file
// with the appropriate extension.
func (renderer *Renderer) IndexFile() string {
	return renderer.indexFile
}
