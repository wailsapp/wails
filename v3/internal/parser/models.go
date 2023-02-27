package parser

import (
	"io"
	"text/template"
)

type ModelDefinitions struct {
	Package string
	Models  map[string]*StructDef
}

func GenerateModel(wr io.Writer, def *ModelDefinitions) error {
	tmpl, err := template.New("model.ts.tmpl").ParseFiles("templates/model.ts.tmpl")
	if err != nil {
		println("Unable to create class template: " + err.Error())
		return err
	}

	err = tmpl.ExecuteTemplate(wr, "model.ts.tmpl", def)
	if err != nil {
		println("Problem executing template: " + err.Error())
		return err
	}
	return nil
}

//func GenerateClass(wr io.Writer, def *StructDef) error {
//	tmpl, err := template.New("class.ts.tmpl").ParseFiles("templates/class.ts.tmpl")
//	if err != nil {
//		println("Unable to create class template: " + err.Error())
//		return err
//	}
//
//	err = tmpl.ExecuteTemplate(wr, "class.ts.tmpl", def)
//	if err != nil {
//		println("Problem executing template: " + err.Error())
//		return err
//	}
//	return nil
//}

//
//import (
//	"bytes"
//	"fmt"
//	"go/ast"
//	"go/types"
//	"sort"
//	"strings"
//	"unicode"
//)
//
//func GenerateModels(context *Context) ([]byte, error) {
//	var buf bytes.Buffer
//	var pkgs []Package
//	specs := context.GetBoundStructs()
//	for pkg, pkgSpecs := range specs {
//		pkgs = append(pkgs, Package{Name: pkg, Specs: pkgSpecs})
//	}
//	knownStructs := newAllModels(specs)
//	sort.Slice(pkgs, func(i, j int) bool { return pkgs[i].Name < pkgs[j].Name })
//	for _, pkg := range pkgs {
//		if _, err := fmt.Fprintf(&buf, "namespace %s {\n", pkg.Name); err != nil {
//			return nil, err
//		}
//		sort.Slice(pkg.Specs, func(i, j int) bool { return pkg.Specs[i].Name.Name < pkg.Specs[j].Name.Name })
//		for _, spec := range pkg.Specs {
//			if structType, ok := spec.Type.(*ast.StructType); ok {
//				if _, err := fmt.Fprintf(&buf, "  class %s {\n", spec.Name.Name); err != nil {
//					return nil, err
//				}
//
//				for _, field := range structType.Fields.List {
//
//					// Ignore field names that have a lower case first letter
//					if !unicode.IsUpper(rune(field.Names[0].Name[0])) {
//						continue
//					}
//
//					// Get the Go type of the field
//					goType := types.ExprString(field.Type)
//					// Check if the type is an array
//					if arrayType, ok := field.Type.(*ast.ArrayType); ok {
//						// Get the element type of the array
//						elementType := types.ExprString(arrayType.Elt)
//						// Look up the corresponding TypeScript type
//						tsType, ok := goToTS[elementType]
//						if !ok {
//							// strip off the * prefix if it is there
//							if strings.HasPrefix(elementType, "*") {
//								elementType = elementType[1:]
//							}
//							if knownStructs.exists(elementType) {
//								tsType = elementType
//							} else {
//								tsType = "any"
//							}
//						}
//						// Output the field as an array of the corresponding TypeScript type
//						if _, err := fmt.Fprintf(&buf, "    %s: %s[];\n", field.Names[0].Name, tsType); err != nil {
//							return nil, err
//						}
//					} else {
//						// strip off the * prefix if it is there
//						if strings.HasPrefix(goType, "*") {
//							goType = goType[1:]
//						}
//						// Look up the corresponding TypeScript type
//						tsType, ok := goToTS[goType]
//						if !ok {
//							if knownStructs.exists(goType) {
//								tsType = goType
//							} else {
//								tsType = "any"
//							}
//						}
//						// Output the field as the corresponding TypeScript type
//						if _, err := fmt.Fprintf(&buf, "    %s: %s;\n", field.Names[0].Name, tsType); err != nil {
//							return nil, err
//						}
//					}
//				}
//
//				if _, err := fmt.Fprintf(&buf, "  }\n"); err != nil {
//					return nil, err
//				}
//			}
//		}
//
//		if _, err := fmt.Fprintf(&buf, "}\n\n"); err != nil {
//			return nil, err
//		}
//	}
//	return buf.Bytes(), nil
//}
