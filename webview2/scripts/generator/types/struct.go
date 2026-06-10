package types

import (
	"fmt"
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

// sizeOf computes the C layout size of the struct (field sizes, natural
// alignment, trailing padding). Field types are restricted by the parser to
// UINT32 | BOOL | BYTE so every size here is exact, not estimated.
func (d *StructDeclaration) sizeOf() (int, error) {
	offset, maxAlign := 0, 1
	for _, f := range d.Fields {
		var size, align int
		switch f.Type {
		case "UINT32", "BOOL":
			size, align = 4, 4
		case "BYTE":
			size, align = 1, 1
		default:
			return 0, fmt.Errorf("struct %s: unknown size for field %s of type %s", d.Name, f.Name, f.Type)
		}
		if align > maxAlign {
			maxAlign = align
		}
		offset = (offset + align - 1) / align * align
		offset += size
	}
	return (offset + maxAlign - 1) / maxAlign * maxAlign, nil
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
