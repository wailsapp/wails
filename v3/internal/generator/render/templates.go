package render

import (
	"embed"
	"strings"
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

var tmplTypedefs = map[tmplLanguage]*template.Template{
	tmplJS: template.Must(template.New("internal.js.tmpl").Funcs(tmplFunctions).ParseFS(templates, "templates/internal.js.tmpl")),
	tmplTS: template.Must(template.New("internal.ts.tmpl").Funcs(tmplFunctions).ParseFS(templates, "templates/internal.ts.tmpl")),
}

var tmplModels = template.Must(template.New("models.tmpl").Funcs(tmplFunctions).ParseFS(templates, "templates/models.tmpl"))

var tmplIndex = template.Must(template.New("index.tmpl").Funcs(tmplFunctions).ParseFS(templates, "templates/index.tmpl"))

var newline string

func init() {
	var builder strings.Builder

	err := template.Must(template.New("newline.tmpl").ParseFS(templates, "templates/newline.tmpl")).Execute(&builder, nil)
	if err != nil {
		panic(err)
	}

	newline = builder.String()
}
