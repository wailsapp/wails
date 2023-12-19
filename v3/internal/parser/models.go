package parser

import (
	"bytes"
	"embed"
	"github.com/wailsapp/wails/v3/internal/flags"
	"io"
	"sort"
	"strings"
	"text/template"
)

//go:embed templates
var templates embed.FS

type ModelDefinitions struct {
	Package string
	Models  map[string]*StructDef
	Enums   map[string]*TypeDef
}

func GenerateModel(wr io.Writer, def *ModelDefinitions, options *flags.GenerateBindingsOptions) error {
	templateName := "model.js.tmpl"
	if options.TS {
		templateName = "model.ts.tmpl"
		if options.UseInterfaces {
			templateName = "interfaces.ts.tmpl"
		}
	}

	// Fix up TS names
	for _, model := range def.Models {
		model.Name = options.TSPrefix + model.Name + options.TSSuffix
	}

	tmpl, err := template.New(templateName).ParseFS(templates, "templates/"+templateName)
	if err != nil {
		println("Unable to create class template: " + err.Error())
		return err
	}

	err = tmpl.ExecuteTemplate(wr, templateName, def)
	if err != nil {
		println("Problem executing template: " + err.Error())
		return err
	}
	return nil
}

const modelsHeader = `// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
`

func pkgAlias(fullPkg string) string {
	pkgParts := strings.Split(fullPkg, "/")
	return pkgParts[len(pkgParts)-1]
}

func GenerateModels(models map[packagePath]map[structName]*StructDef, enums map[packagePath]map[string]*TypeDef, options *flags.GenerateBindingsOptions) (string, error) {
	if models == nil {
		return "", nil
	}

	var buffer bytes.Buffer
	buffer.WriteString(modelsHeader)

	// sort pkgs by alias (e.g. services) instead of full pkg name (e.g. github.com/wailsapp/wails/somedir/services)
	// and then sort resulting list by the alias
	var keys []string
	for pkg := range models {
		keys = append(keys, pkg)
	}

	sort.Slice(keys, func(i, j int) bool {
		return pkgAlias(keys[i]) < pkgAlias(keys[j])
	})

	for _, pkg := range keys {
		err := GenerateModel(&buffer, &ModelDefinitions{
			Package: pkgAlias(pkg),
			Models:  models[pkg],
			Enums:   enums[pkg],
		}, options)
		if err != nil {
			return "", err
		}
	}
	return buffer.String(), nil
}
