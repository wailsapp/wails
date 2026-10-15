package types

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"text/template"
)

//go:embed templates/*
var templates embed.FS

// templateSource is the filesystem templates are read from. It defaults to the
// embedded FS but is a package variable (not the embed.FS directly) so tests
// can substitute a fault-injecting filesystem to exercise the read/parse/
// execute error paths that the embedded templates can never trigger at runtime.
var templateSource fs.FS = templates

// renderTemplate reads, parses, and executes the named template (relative to
// the templates/ directory) into w. It returns a descriptive error rather than
// terminating the process, so a malformed or missing template surfaces as a
// generation error the caller can handle.
func renderTemplate(name, file string, data any, w io.Writer) error {
	templateData, err := fs.ReadFile(templateSource, "templates/"+file)
	if err != nil {
		return fmt.Errorf("read template %s: %w", file, err)
	}
	tmpl, err := template.New(name).Parse(string(templateData))
	if err != nil {
		return fmt.Errorf("parse template %s: %w", file, err)
	}
	if err := tmpl.Execute(w, data); err != nil {
		return fmt.Errorf("execute template %s: %w", file, err)
	}
	return nil
}
