package templates

import (
	"bytes"
	"text/template"
)

// RootTaskfile renders the root Taskfile.yml that `wails3 init` places in a
// new project, for the given binary name. It is used by `wails3 migrate` to
// give a migrated v2 project the same task-based build system as a fresh v3
// project.
func RootTaskfile(binaryName string) ([]byte, error) {
	data, err := templates.ReadFile("_common/Taskfile.tmpl.yml")
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New("Taskfile.yml").Parse(string(data))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, struct {
		BinaryName string
		Opn        string
		Cls        string
	}{
		BinaryName: binaryName,
		Opn:        "{{",
		Cls:        "}}",
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
