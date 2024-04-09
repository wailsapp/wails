package parser

import (
	"go/ast"
	"go/doc"
	"path"
	"strings"

	"golang.org/x/tools/go/packages"
)

func NewDoc(pkg *packages.Package) *doc.Package {
	files := make(map[string]*ast.File)
	for i, f := range pkg.Syntax {
		files[path.Base(pkg.CompiledGoFiles[i])] = f
	}

	return doc.New(&ast.Package{
		Name:  pkg.Name,
		Files: files,
	}, pkg.PkgPath, doc.PreserveAST)
}

func (m *BoundMethod) DocComment(pkg *Package) string {
	for _, t := range pkg.doc.Types {
		if t.Name == m.Service.Name() {
			for _, f := range t.Methods {
				if f.Name == m.Name() {
					return strings.TrimSpace(f.Doc)
				}
			}
			break
		}
	}
	return ""
}

func (e *EnumDef) DocComment(pkg *Package) string {
	for _, t := range pkg.doc.Types {
		if t.Name == e.Name {
			return strings.TrimSpace(t.Doc)
		}
	}
	return ""
}

func (s *StructDef) DocComment(pkg *Package) string {
	for _, t := range pkg.doc.Types {
		if t.Name == s.Name {
			return strings.TrimSpace(t.Doc)
		}
	}
	return ""
}

func (c *ConstDef) DocComment(pkg *Package, enum *EnumDef) string {
	for _, t := range pkg.doc.Types {
		if t.Name == enum.Name {
			for _, value := range t.Consts {
				// comment of a single const declaration
				if len(value.Names) == 1 && c.Name == value.Names[0] {
					return strings.TrimSpace(value.Doc)
				}

				// comments inside a grouped const declaration
				specs := value.Decl.Specs
				for i, name := range value.Names {
					if spec, ok := specs[i].(*ast.ValueSpec); ok && name == c.Name {
						return strings.TrimSpace(spec.Doc.Text())
					}
				}
			}
			break
		}
	}
	return ""
}

func (p *Parameter) DocComment(pkg *Package) string {
	//TODO
	return ""
}
