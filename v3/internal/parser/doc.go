package parser

import (
	"go/ast"
	"go/doc"
	"path"
	"strings"

	"golang.org/x/tools/go/packages"
)

type Doc struct {
	*doc.Package
	Types map[string]*doc.Type

	// Methods map[*doc.Type]map[string]string
	// Fields  map[*doc.Type]map[string]string
	// Consts  map[*doc.Type]map[string]string
}

func NewDoc(pkg *packages.Package) *Doc {
	files := make(map[string]*ast.File)
	for i, f := range pkg.Syntax {
		files[path.Base(pkg.CompiledGoFiles[i])] = f
	}

	pkgDoc := doc.New(&ast.Package{
		Name:  pkg.Name,
		Files: files,
	}, pkg.PkgPath, doc.PreserveAST)

	types := make(map[string]*doc.Type)
	for _, t := range pkgDoc.Types {
		types[t.Name] = t
	}

	return &Doc{
		Package: pkgDoc,
		Types:   types,
	}
}

func (m *BoundMethod) DocComment(pkg *Package) string {
	serviceType, ok := pkg.doc.Types[m.Service.Name()]
	if !ok {
		return ""
	}
	for _, f := range serviceType.Methods {
		if f.Name == m.Name() {
			return strings.TrimSpace(f.Doc)
		}
	}
	return ""
}

func (e *EnumDef) DocComment(pkg *Package) string {
	if enumType, ok := pkg.doc.Types[e.Name]; ok {
		return strings.TrimSpace(enumType.Doc)
	}
	return ""
}

func (s *StructDef) DocComment(pkg *Package) string {
	if structType, ok := pkg.doc.Types[s.Name]; ok {
		return strings.TrimSpace(structType.Doc)
	}
	return ""
}

func (c *ConstDef) DocComment(pkg *Package, enum *EnumDef) string {
	enumType, ok := pkg.doc.Types[enum.Name]
	if !ok {
		return ""
	}

	for _, value := range enumType.Consts {
		// comment of a single const declaration
		if len(value.Names) == 1 && c.Name == value.Names[0] {
			return strings.TrimSpace(value.Doc)
		}

		// comments inside a grouped const declaration
		for i, spec := range value.Decl.Specs {
			if spec, ok := spec.(*ast.ValueSpec); ok && value.Names[i] == c.Name {
				return strings.TrimSpace(spec.Doc.Text())
			}
		}
	}
	return ""
}

func (f *Field) DocComment(pkg *Package, structDef *StructDef) string {
	structType, ok := pkg.doc.Types[structDef.Name]
	if !ok {
		return ""
	}
	for _, spec := range structType.Decl.Specs {
		if spec, ok := spec.(*ast.TypeSpec); ok {
			if t, ok := spec.Type.(*ast.StructType); ok {
				return f.docComment(t)
			}
		}
	}
	return ""
}

func (f *Field) docComment(structType *ast.StructType) string {
	for _, field := range structType.Fields.List {
		for _, ident := range field.Names {
			if ident.Name == f.Name() {
				return strings.TrimSpace(field.Doc.Text())
			}
		}
	}
	return ""
}
