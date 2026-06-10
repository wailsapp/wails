package types

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/leaanthony/slicer"
	"log"
	"strings"
	"text/template"
)

type GeneratedFile struct {
	FileName string
	Package  string
	Content  *bytes.Buffer
}

type IDL struct {
	Imports   []*Import  `parser:"@@*"`
	Libraries []*Library `parser:"@@*"`
}

func (i *IDL) Process() error {
	for _, library := range i.Libraries {
		err := library.Process()
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *IDL) Generate() ([]*GeneratedFile, error) {
	for _, library := range i.Libraries {
		return library.Generate()
	}
	return nil, nil
}

type Import struct {
	Name string `parser:"'import' @(!';')* ';'"`
}

type LibraryHeader struct {
	UUID string `parser:"'uuid' '(' @UUID ')' ',' 'version' '(' Int ('.' Int)? ')'"`
}

type Library struct {
	Header       *LibraryHeader `parser:"'[' @@ ']'"`
	Name         string         `parser:"'library' @Ident"`
	Declarations []*Declaration `parser:"'{' @@* '}'"`

	// private
	forewardInterfaceDeclarations slicer.StringSlicer
	enums                         slicer.StringSlicer
	packageName                   string
	// structSizes maps IDL typedef struct names to their byte size so by-value
	// parameter marshalling can pick the right argument encoding.
	structSizes map[string]int
	// interfaceBases maps each declared interface to its immediate base
	// interface, used to lay out vtbls and choose QI helper receivers.
	interfaceBases map[string]string
}

func (l *Library) Process() error {
	l.packageName = strings.ToLower(l.Name)
	// Pre-register struct sizes and the interface inheritance graph so
	// references resolve regardless of declaration order.
	l.structSizes = map[string]int{}
	l.interfaceBases = map[string]string{}
	for _, declaration := range l.Declarations {
		if declaration.Struct != nil {
			size, err := declaration.Struct.sizeOf()
			if err != nil {
				return err
			}
			l.structSizes[declaration.Struct.Name] = size
		}
		if declaration.Interface != nil {
			l.interfaceBases[declaration.Interface.Name] = declaration.Interface.BaseClass
		}
	}
	for _, declaration := range l.Declarations {
		err := declaration.Process(l)
		if err != nil {
			return err
		}
	}
	return nil
}

// chainRoot returns the first interface in name's inheritance chain whose
// base is IUnknown — the object identity every other interface on that COM
// object can be queried from.
func (l *Library) chainRoot(name string) (string, error) {
	for depth := 0; depth < 32; depth++ {
		base, ok := l.interfaceBases[name]
		if !ok || base == "IUnknown" {
			return name, nil
		}
		name = base
	}
	return "", fmt.Errorf("interface inheritance chain for %s exceeds 32 levels — cycle?", name)
}

func (l *Library) Generate() ([]*GeneratedFile, error) {
	result := l.GenerateDefaultFiles()

	for _, declaration := range l.Declarations {
		generatedFile, err := declaration.Generate()
		if err != nil {
			return nil, err
		}
		if generatedFile != nil {
			result = append(result, generatedFile)
		}
	}

	return result, nil
}

func (l *Library) addInterfaceName(interfaceName string) {
	l.forewardInterfaceDeclarations.Add(interfaceName)
}

func (l *Library) GenerateDefaultFiles() []*GeneratedFile {
	data := struct {
		PackageName string
	}{
		PackageName: l.packageName,
	}

	var result []*GeneratedFile
	var buf bytes.Buffer

	templateData, err := templates.ReadFile("templates/com.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	tmpl, err := template.New("COM").Parse(string(templateData))
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(&buf, &data)
	if err != nil {
		log.Fatal(err)
	}

	result = append(result, &GeneratedFile{
		FileName: "com.go",
		Package:  l.packageName,
		Content:  &buf,
	})

	return result

}

type Declaration struct {
	InterfaceForewardDecl string                `parser:"'interface' @Ident ';'"`
	Enum                  *EnumDeclaration      `parser:"| '[' 'v1_enum' ']' @@"`
	Struct                *StructDeclaration    `parser:"| @@"`
	Interface             *InterfaceDeclaration `parser:"| @@"`
	CppQuote              string                `parser:"| 'cpp_quote' '(' @String ')'"`

	// Private
	library *Library
}

func (d *Declaration) Process(l *Library) error {
	d.library = l
	if d.Enum != nil {
		return d.Enum.Process(d)
	}
	if d.Struct != nil {
		return d.Struct.Process(d)
	}
	if d.Interface != nil {
		return d.Interface.Process(d)
	}
	if d.CppQuote != "" {
		return nil
	}
	if d.InterfaceForewardDecl != "" {
		l.addInterfaceName(d.InterfaceForewardDecl)
		return nil
	}
	return errors.New("unknown declaration to process")
}

func (d *Declaration) Generate() (*GeneratedFile, error) {

	var buffer bytes.Buffer
	var packageName = strings.ToLower(d.library.Name)
	var filename string

	if d.Enum != nil {
		err := d.Enum.Generate(packageName, &buffer)
		if err != nil {
			return nil, err
		}
		filename = d.Enum.Name + ".go"
	}
	if d.Struct != nil {
		err := d.Struct.Generate(packageName, &buffer)
		if err != nil {
			return nil, err
		}
		filename = d.Struct.Name + ".go"
	}
	if d.Interface != nil {
		err := d.Interface.Generate(packageName, &buffer)
		if err != nil {
			return nil, err
		}
		filename = d.Interface.Name + ".go"
	}
	if d.CppQuote != "" {
		return nil, nil
	}
	if d.InterfaceForewardDecl != "" {
		return nil, nil
	}
	//f := filepath.Join(packageDir, filename)
	//err := os.WriteFile(f, buffer.Bytes(), 0755)
	return &GeneratedFile{
		FileName: filename,
		Package:  packageName,
		Content:  &buffer,
	}, nil

}
