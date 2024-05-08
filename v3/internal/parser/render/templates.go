package render

import (
	"embed"
	"text/template"
)

//go:embed templates/*.tmpl
var templates embed.FS

type tmplLanguage bool

const tmplJS, tmplTS tmplLanguage = false, true

type tmplMode bool

const tmplClasses, tmplInterfaces tmplMode = false, true

var tmplBindings = map[tmplLanguage]*template.Template{
	tmplJS: template.Must(template.New("bindings.js.tmpl").Funcs(tmplFunctions).ParseFS(templates, "templates/bindings.js.tmpl")),
	tmplTS: template.Must(template.New("bindings.ts.tmpl").Funcs(tmplFunctions).ParseFS(templates, "templates/bindings.ts.tmpl")),
}

var tmplModels = map[tmplLanguage]map[tmplMode]*template.Template{
	tmplJS: {
		tmplClasses: template.Must(template.New("models.js.tmpl").Funcs(tmplFunctions).ParseFS(templates, "templates/models.js.tmpl")),
	},
	tmplTS: {
		tmplClasses:    template.Must(template.New("models.ts.tmpl").Funcs(tmplFunctions).ParseFS(templates, "templates/models.ts.tmpl")),
		tmplInterfaces: template.Must(template.New("interfaces.ts.tmpl").Funcs(tmplFunctions).ParseFS(templates, "templates/interfaces.ts.tmpl")),
	},
}

var tmplIndex = map[tmplLanguage]*template.Template{
	tmplJS: template.Must(template.New("index.js.tmpl").Funcs(tmplFunctions).ParseFS(templates, "templates/index.js.tmpl")),
	tmplTS: template.Must(template.New("index.ts.tmpl").Funcs(tmplFunctions).ParseFS(templates, "templates/index.ts.tmpl")),
}

var tmplStruct *template.Template

func init() {
	tmplModels[tmplJS][tmplInterfaces] = tmplModels[tmplJS][tmplClasses]

	// Init struct template here to break initialisation cycle.
	tmplStruct = template.Must(template.New("struct.tmpl").Funcs(tmplFunctions).ParseFS(templates, "templates/struct.tmpl"))
}
