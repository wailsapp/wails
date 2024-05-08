package render

import (
	"embed"
	"text/template"
)

//go:embed templates/*.tmpl
var templates embed.FS

type tmplLanguage bool

const tmplJS, tmplTS tmplLanguage = false, true

var tmplBindings = map[tmplLanguage]*template.Template{
	tmplJS: template.Must(template.New("bindings.js.tmpl").Funcs(tmplFunctions).ParseFS(templates, "templates/bindings.js.tmpl")),
	tmplTS: template.Must(template.New("bindings.ts.tmpl").Funcs(tmplFunctions).ParseFS(templates, "templates/bindings.ts.tmpl")),
}

var tmplModels = map[tmplLanguage]*template.Template{
	tmplJS: template.Must(template.New("models.js.tmpl").Funcs(tmplFunctions).ParseFS(templates, "templates/models.js.tmpl")),
	tmplTS: template.Must(template.New("models.ts.tmpl").Funcs(tmplFunctions).ParseFS(templates, "templates/models.ts.tmpl")),
}

var tmplIndex = template.Must(template.New("index.tmpl").Funcs(tmplFunctions).ParseFS(templates, "templates/index.tmpl"))
var tmplGlobalIndex = template.Must(template.New("global.tmpl").Funcs(tmplFunctions).ParseFS(templates, "templates/global.tmpl"))
var tmplShortcut = template.Must(template.New("shortcut.tmpl").Funcs(tmplFunctions).ParseFS(templates, "templates/shortcut.tmpl"))
