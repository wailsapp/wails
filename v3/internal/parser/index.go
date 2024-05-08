package parser

import (
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/parser/templates"
)

// generateIndex generates an index file from the given index information.
// A call to index.Info.Collect must complete before entering generateIndex.
func (generator *Generator) generateIndex(index PackageIndex) {
	file, err := generator.create(filepath.Join(generator.options.OutputDirectory, index.Info.Path, "index.js"))
	if err != nil {
		pterm.Error.Println(err)
		pterm.Error.Printfln("package %s: index generation failed", index.Info.Path)
		return
	}

	success := true

	template := templates.IndexJS
	if generator.options.TS {
		template = templates.IndexTS
	}

	if err := template.Execute(file, &index); err != nil {
		success = false
		pterm.Error.Println(err)
	}

	if err := file.Close(); err != nil {
		success = false
		pterm.Error.Println(err)
	}

	if !success {
		pterm.Error.Printfln("package %s: index generation failed", index.Info.Path)
	}
}
