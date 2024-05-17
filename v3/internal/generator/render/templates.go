package render

import (
	"embed"
	"text/template"
)

//go:embed templates/*.tmpl
var templates embed.FS

type tmplLanguage bool

const tmplJS, tmplTS tmplLanguage = false, true

var tmplService = map[tmplLanguage]*template.Template{
	tmplJS: template.Must(template.New("service.js.tmpl").Funcs(tmplFunctions).ParseFS(templates, "templates/service.js.tmpl")),
	tmplTS: template.Must(template.New("service.ts.tmpl").Funcs(tmplFunctions).ParseFS(templates, "templates/service.ts.tmpl")),
}

var tmplModels = map[tmplLanguage]*template.Template{
	tmplJS: template.Must(template.New("models.js.tmpl").Funcs(tmplFunctions).ParseFS(templates, "templates/models.js.tmpl")),
	tmplTS: template.Must(template.New("models.ts.tmpl").Funcs(tmplFunctions).ParseFS(templates, "templates/models.ts.tmpl")),
}

var tmplIndex = template.Must(template.New("index.tmpl").Funcs(tmplFunctions).ParseFS(templates, "templates/index.tmpl"))
