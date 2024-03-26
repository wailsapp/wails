package templates

import (
	"embed"
	"strings"
	"text/template"
)

//go:embed *.tmpl
var templates embed.FS

var functions = template.FuncMap{
	"jsdoc": func(comment string, indent string) string {
		return strings.ReplaceAll(strings.ReplaceAll(comment, "*/", "*\\/"), "\n", "\n"+indent+" * ")
	},
	"paramdoc": func(comment string) string {
		return strings.ReplaceAll(strings.ReplaceAll(comment, "*/", "*\\/"), "\n", " ")
	},
}

var BindingsJS = template.Must(template.New("bindings.js.tmpl").Funcs(functions).ParseFS(templates, "bindings.js.tmpl"))
var BindingsTS = template.Must(template.New("bindings.ts.tmpl").Funcs(functions).ParseFS(templates, "bindings.ts.tmpl"))

var InterfacesTS = template.Must(template.New("interfaces.ts.tmpl").Funcs(functions).ParseFS(templates, "interfaces.ts.tmpl"))
var ModelsTS = template.Must(template.New("models.ts.tmpl").Funcs(functions).ParseFS(templates, "models.ts.tmpl"))
var ModelsJS = template.Must(template.New("models.js.tmpl").Funcs(functions).ParseFS(templates, "models.js.tmpl"))
