package types

import (
	"io"
	"log"
	"text/template"
)

type StructDeclaration struct {
	Name   string         `parser:"'typedef' 'struct' @Ident '{' "`
	Fields []*StructField `parser:" (@@)+ '}' Ident ';'"`

	// private
	decl *Declaration
}

func (d *StructDeclaration) Process(decl *Declaration) error {
	d.decl = decl
	for _, f := range d.Fields {
		f.Process()
	}
	return nil
}

func (d *StructDeclaration) Generate(packageName string, w io.Writer) error {
	data := struct {
		PackageName string
		Name        string
		Fields      []*StructField
	}{
		PackageName: packageName,
		Name:        d.Name,
		Fields:      d.Fields,
	}
	templateData, err := templates.ReadFile("templates/struct.tmpl")
	if err != nil {
		return err
	}
	tmpl, err := template.New("Struct").Parse(string(templateData))
	if err != nil {
		log.Fatalln(err)
	}
	return tmpl.Execute(w, &data)
}

type StructField struct {
	Type string `parser:"@('UINT32' | 'BOOL' | 'BYTE')"`
	Name string `parser:"@Ident ';'"`

	// process
	GoType string
}

func (s *StructField) Process() {
	s.GoType = IdlTypeToGoType(s.Type)
}
