package collect

import (
	"cmp"
	"go/ast"
	"go/token"
	"go/types"
	"slices"
)

// findDeclaration returns the AST spec or declaration
// that defines the given _global_ type-checker object.
//
// Specifically, the first element in the returned slice
// is the relevant spec or declaration, followed by its chain
// of parent nodes up to the declaring [ast.File].
//
// If no corresponding declaration can be found within
// the set of registered packages, the returned slice is nil.
//
// Resulting node types are as follows:
//   - global functions and concrete methods (*types.Func)
//     map to *ast.FuncDecl nodes;
//   - interface methods from global interfaces (*types.Func)
//     map to *ast.Field nodes within their interface expression;
//   - struct fields from global structs (*types.Var)
//     map to *ast.Field nodes within their struct expression;
//   - global constants and variables map to *ast.ValueSpec nodes;
//   - global named types map to *ast.TypeSpec nodes;
//   - for type parameters, the result is always nil;
//   - for local objects defined within functions,
//     field types, variable types or field values,
//     the result is always nil;
//
// findDeclaration supports unsynchronised concurrent calls.
func (collector *Collector) findDeclaration(obj types.Object) (path []ast.Node) {
	pkg := collector.Package(obj.Pkg()).Collect()
	if pkg == nil {
		return nil
	}

	// Perform a binary search to find the file enclosing the node.
	// We can't use findEnclosingNode here because it is less accurate and less efficient with files.
	fileIndex, exact := slices.BinarySearchFunc(pkg.Files, obj.Pos(), func(f *ast.File, p token.Pos) int {
		return cmp.Compare(f.FileStart, p)
	})

	// If exact is true, pkg.Files[fileIndex] is the file we are looking for;
	// otherwise, it is the first file whose start position is _after_ obj.Pos().
	if !exact {
		fileIndex--
	}

	// When exact is false, the position might lie within an empty segment in between two files.
	if fileIndex < 0 || pkg.Files[fileIndex].FileEnd <= obj.Pos() {
		return nil
	}

	file := pkg.Files[fileIndex]

	// Find enclosing declaration.
	decl := findEnclosingNode(file.Decls, obj.Pos())
	if decl == nil {
		// Invalid position.
		return nil
	}

	var gen *ast.GenDecl

	switch d := decl.(type) {
	case *ast.FuncDecl:
		if obj.Pos() == d.Name.Pos() {
			// Object is function.
			return []ast.Node{decl, file}
		}

		// Ignore local objects defined within function bodies.
		return nil

	case *ast.BadDecl:
		// What's up??
		return nil

	case *ast.GenDecl:
		gen = d
	}

	// Handle *ast.GenDecl

	// Find enclosing ast.Spec
	spec := findEnclosingNode(gen.Specs, obj.Pos())
	if spec == nil {
		// Invalid position.
		return nil
	}

	var def ast.Expr

	switch s := spec.(type) {
	case *ast.ValueSpec:
		if s.Names[0].Pos() <= obj.Pos() && obj.Pos() < s.Names[len(s.Names)-1].End() {
			// Object is variable or constant.
			return []ast.Node{spec, decl, file}
		}

		// Ignore local objects defined within variable types/values.
		return nil

	case *ast.TypeSpec:
		if obj.Pos() == s.Name.Pos() {
			// Object is named type.
			return []ast.Node{spec, decl, file}
		}

		if obj.Pos() < s.Type.Pos() || s.Type.End() <= obj.Pos() {
			// Type param or invalid position.
			return nil
		}

		// Struct or interface field?
		def = s.Type
	}

	// Handle struct or interface field.

	var iface *ast.InterfaceType

	switch d := def.(type) {
	case *ast.StructType:
		// Find enclosing field
		field := findEnclosingNode(d.Fields.List, obj.Pos())
		if field == nil {
			// Invalid position.
			return nil
		}

		if len(field.Names) == 0 {
			// Handle embedded field.
			ftype := ast.Unparen(field.Type)

			// Unwrap pointer.
			if ptr, ok := ftype.(*ast.StarExpr); ok {
				ftype = ast.Unparen(ptr.X)
			}

			// Unwrap generic instantiation.
			switch t := field.Type.(type) {
			case *ast.IndexExpr:
				ftype = ast.Unparen(t.X)
			case *ast.IndexListExpr:
				ftype = ast.Unparen(t.X)
			}

			// Unwrap selector.
			if sel, ok := ftype.(*ast.SelectorExpr); ok {
				ftype = sel.Sel
			}

			// ftype must now be an identifier.
			if obj.Pos() == ftype.Pos() {
				// Object is this embedded field.
				return []ast.Node{field, d.Fields, def, spec, decl, file}
			}
		} else if field.Names[0].Pos() <= obj.Pos() && obj.Pos() < field.Names[len(field.Names)-1].End() {
			// Object is one of these fields.
			return []ast.Node{field, d.Fields, def, spec, decl, file}
		}

		// Ignore local objects defined within field types.
		return nil

	case *ast.InterfaceType:
		iface = d

	default:
		// Other local object or invalid position.
		return nil
	}

	path = []ast.Node{file, decl, spec, def, iface.Methods}

	// Handle interface method.
	for {
		field := findEnclosingNode(iface.Methods.List, obj.Pos())
		if field == nil {
			// Invalid position.
			return nil
		}

		path = append(path, field)

		if len(field.Names) == 0 {
			// Handle embedded interface.
			var ok bool
			iface, ok = ast.Unparen(field.Type).(*ast.InterfaceType)
			if !ok {
				// Not embedded interface, ignore.
				return nil
			}

			path = append(path, iface, iface.Methods)
			// Explore embedded interface.

		} else if field.Names[0].Pos() <= obj.Pos() && obj.Pos() < field.Names[len(field.Names)-1].End() {
			// Object is one of these fields.
			slices.Reverse(path)
			return path
		} else {
			// Ignore local objects defined within interface method signatures.
			return nil
		}
	}
}

// findEnclosingNode finds the unique node in nodes, if any,
// that encloses the given position.
//
// It uses binary search and therefore expects
// the nodes slice to be sorted in source order.
func findEnclosingNode[S ~[]E, E ast.Node](nodes S, pos token.Pos) (node E) {
	// Perform a binary search to find the nearest node.
	index, exact := slices.BinarySearchFunc(nodes, pos, func(n E, p token.Pos) int {
		return cmp.Compare(n.Pos(), p)
	})

	// If exact is true, nodes[index] is the node we are looking for;
	// otherwise, it is the first node whose start position is _after_ pos.
	if !exact {
		index--
	}

	// When exact is false, the position might lie within an empty segment in between two nodes.
	if index < 0 || nodes[index].End() <= pos {
		return // zero value, nil in practice.
	}

	return nodes[index]
}
