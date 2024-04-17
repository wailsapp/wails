package parser

import (
	"go/ast"
	"go/doc"
	"strings"
)

type Doc struct {
	*doc.Package
	Types map[string]*doc.Type
}

func NewDoc(pkgPath string, pkg *ast.Package) *Doc {
	pkgDoc := doc.New(pkg, pkgPath, doc.PreserveAST|doc.AllDecls)

	types := make(map[string]*doc.Type)
	for _, t := range pkgDoc.Types {
		types[t.Name] = t
	}

	return &Doc{
		Package: pkgDoc,
		Types:   types,
	}
}

func (m *BoundMethod) DocComment(pkg *Package, service *Service) string {
	serviceType, ok := pkg.doc.Types[service.Name()]
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

func (a *AliasDef) DocComment(pkg *Package) string {
	if basic, ok := pkg.doc.Types[a.Name]; ok {
		return strings.TrimSpace(basic.Doc)
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

func (f *Field) DocComment(pkg *Package) string {
	structType, ok := pkg.doc.Types[f.origin.Name]
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
		// comment of embedded basic type
		if fieldType, ok := field.Type.(*ast.Ident); ok && len(field.Names) == 0 && fieldType.Name == f.Name() {
			return strings.TrimSpace(field.Doc.Text())
		}

		for _, ident := range field.Names {
			if ident.Name == f.Name() {
				return strings.TrimSpace(field.Doc.Text())
			}
		}
	}
	return ""
}
