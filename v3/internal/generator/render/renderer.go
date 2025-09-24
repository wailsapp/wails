package render

import (
	"go/types"
	"io"
	"slices"
	"strings"
	"text/template"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/generator/collect"
)

// Renderer holds the template set for a given configuration.
// It provides methods for rendering various output modules.
type Renderer struct {
	collector *collect.Collector
	options   *flags.GenerateBindingsOptions

	ext string

	service *template.Template
	models  *template.Template
}

// NewRenderer initialises a code renderer
// for the given configuration and data collector.
func NewRenderer(collector *collect.Collector, options *flags.GenerateBindingsOptions) *Renderer {
	ext := ".js"
	if options.TS {
		ext = ".ts"
	}

	return &Renderer{
		collector: collector,
		options:   options,

		ext: ext,

		service: tmplService[tmplLanguage(options.TS)],
		models:  tmplModels[tmplLanguage(options.TS)],
	}
}

// ServiceFile returns the standard name of a service file
// for the given struct name, with the appropriate extension.
func (renderer *Renderer) ServiceFile(name string) string {
	return strings.ToLower(name) + renderer.ext
}

// ModelsFile returns the standard name of a models file
// with the appropriate extension.
func (renderer *Renderer) ModelsFile() string {
	return renderer.options.ModelsFilename + renderer.ext
}

// EventDataFile returns the standard name of the event data definitions file
// with the appropriate extension.
func (renderer *Renderer) EventDataFile() string {
	return "eventdata.d.ts"
}

// EventCreateFile returns the standard name of the event data creation file
// with the appropriate extension.
func (renderer *Renderer) EventCreateFile() string {
	return "eventcreate" + renderer.ext
}

// IndexFile returns the standard name of a package index file
// with the appropriate extension.
func (renderer *Renderer) IndexFile() string {
	return renderer.options.IndexFilename + renderer.ext
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

// Models renders type definitions for the given list of models.
func (renderer *Renderer) Models(w io.Writer, imports *collect.ImportMap, models []*collect.ModelInfo) error {
	if !renderer.options.UseInterfaces {
		// Sort class aliases after the class they alias.
		// Works in amortized linear time thanks to an auxiliary map.

		// Track postponed class aliases and their dependencies.
		aliases := make(map[types.Object][]*collect.ModelInfo, len(models))

		models = slices.Clone(models)
		for i, j := 0, 0; i < len(models); i++ {
			if models[i].Type != nil && models[i].Predicates.IsClass {
				// models[i] is a class alias:
				// models[i].Type is guaranteed to be
				// either an alias or a named type
				obj := models[i].Type.(interface{ Obj() *types.TypeName }).Obj()
				if obj.Pkg().Path() == imports.Self {
					// models[i] aliases a type from the current module.
					if a, ok := aliases[obj]; !ok || len(a) > 0 {
						// The aliased type has not been visited already, postpone.
						aliases[obj] = append(a, models[i])
						continue
					}
				}
			}

			// Append models[i].
			models[j] = models[i]
			j++

			// Keep appending aliases whose aliased type has been just appended.
			for k := j - 1; k < j; k++ {
				a := aliases[models[k].Object()]
				aliases[models[k].Object()] = nil // Mark aliased model as visited
				j += copy(models[j:], a)
			}
		}
	}

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

// EventData renders the given event map to w as an event data table.
func (renderer *Renderer) EventData(w io.Writer, events *collect.EventMap) error {
	return tmplEventData.Execute(w, &struct {
		module
		Events *collect.EventMap
	}{
		module{
			Renderer:                renderer,
			GenerateBindingsOptions: renderer.options,
			Imports:                 events.Imports,
		},
		events,
	})
}

// EventCreate renders the given event map to w as event data creation code.
func (renderer *Renderer) EventCreate(w io.Writer, events *collect.EventMap) error {
	return tmplEventCreate.Execute(w, &struct {
		module
		Events *collect.EventMap
	}{
		module{
			Renderer:                renderer,
			GenerateBindingsOptions: renderer.options,
			Imports:                 events.Imports,
		},
		events,
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
