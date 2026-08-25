package types

import (
	"io"
	"log"
	"strconv"
	"text/template"
)

type EnumDeclaration struct {
	Name   string       `parser:"'typedef' 'enum' @Ident '{'"`
	Values []*EnumValue `parser:" (@@)+ '}' Ident ';'"`

	// private
	decl *Declaration
}

type EnumValue struct {
	Key   string         `parser:"@Ident"`
	Value *EnumValueDecl `parser:"('=' @@)? ','?"`
}

type EnumValueDecl struct {
	Value     string `parser:"@Hex | @Int"`
	LeftShift *int   `parser:"('<' '<' @Int)?"`
}

func (e *EnumValueDecl) Process() {
	if e.LeftShift != nil {
		e.Value += " << " + strconv.Itoa(*e.LeftShift)
	}
}

func (d *EnumDeclaration) Process(decl *Declaration) error {
	d.decl = decl
	for index, value := range d.Values {
		if value.Value == nil {
			value.Value = &EnumValueDecl{
				Value: strconv.Itoa(index),
			}
		} else {
			value.Value.Process()
		}
	}
	decl.library.enums.Add(d.Name)
	return nil
}

func (d *EnumDeclaration) Generate(packageName string, w io.Writer) error {
	data := struct {
		PackageName string
		Name        string
		Values      []*EnumValue
	}{
		PackageName: packageName,
		Name:        d.Name,
		Values:      d.Values,
	}
	templateData, err := templates.ReadFile("templates/enum.tmpl")
	if err != nil {
		return err
	}
	tmpl, err := template.New("Enum").Parse(string(templateData))
	if err != nil {
		log.Fatalln(err)
	}
	return tmpl.Execute(w, &data)
}
