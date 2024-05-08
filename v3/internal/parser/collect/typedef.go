package collect

import (
	"go/ast"
	"go/token"
	"go/types"
	"slices"
	"sync"
)

type (
	// TypeDefInfo records information about a single type declaration.
	TypeDefInfo struct {
		Name    string
		Doc     *ast.CommentGroup
		Group   *GroupInfo
		Alias   bool
		Methods map[string]*MethodInfo

		fields func() map[string]*FieldDefInfo

		// The following fields are used by the Rhs retrieval algorithm.
		mu sync.Mutex

		pkg  *PackageInfo
		file *ast.File
		def  ast.Expr

		obj *types.TypeName
		rhs types.Type
	}

	// FieldDefInfo records information about a struct field declaration.
	FieldDefInfo struct {
		// Do not store the name here to avoid collisions with [FieldInfo.Name].

		Pos   token.Pos
		Doc   *ast.CommentGroup
		Group *GroupInfo
	}

	// MethodInfo records information about a method declaration.
	MethodInfo struct {
		Name string
		Doc  *ast.CommentGroup
	}
)

// newTypeDefInfo initialises a descriptor for the given type spec.
func newTypeDefInfo(pkg *PackageInfo, file *ast.File, group *GroupInfo, spec *ast.TypeSpec) *TypeDefInfo {
	info := &TypeDefInfo{
		Name:  spec.Name.Name,
		Doc:   spec.Doc,
		Group: group,
		Alias: spec.Assign.IsValid(),

		pkg:  pkg,
		file: file,
		def:  ast.Unparen(spec.Type),
	}

	info.fields = sync.OnceValue(info.collectFields)

	return info
}

// Fields returns information about the fields of a struct type.
// If the receiver does not describe a struct type, Fields returns nil.
func (info *TypeDefInfo) Fields() map[string]*FieldDefInfo {
	return info.fields()
}

// collectFields parses a struct definition to extract field comments.
// It is meant to be wrapped by [sync.OnceValue], for result caching.
func (info *TypeDefInfo) collectFields() map[string]*FieldDefInfo {
	def, ok := info.def.(*ast.StructType)
	if !ok {
		return nil
	}

	result := make(map[string]*FieldDefInfo)

	for _, field := range def.Fields.List {

		doc := field.Doc
		group := &GroupInfo{
			Doc:   doc,
			Group: nil,
		}

		if len(field.Names) > 1 {
			doc = nil
		} else {
			group.Doc = nil
		}

		if len(field.Names) == 0 {
			// Embedded field, do not ignore unexported names.
			var name *ast.Ident

			// Unwrap pointer expression.
			typ := field.Type
			if ptr, ok := field.Type.(*ast.StarExpr); ok {
				typ = ptr.X
			}

			// Unwrap generic type instantiation.
			switch t := typ.(type) {
			case *ast.IndexExpr:
				typ = t.X
			case *ast.IndexListExpr:
				typ = t.X
			}

			switch t := typ.(type) {
			case *ast.Ident:
				name = t
			case *ast.SelectorExpr:
				name = t.Sel
			}

			if name == nil {
				// Invalid embedded field.
				continue
			}

			if _, present := result[name.Name]; present {
				// Ignore redefinitions.
				continue
			}

			result[name.Name] = &FieldDefInfo{
				Pos:   name.Pos(),
				Doc:   doc,
				Group: group,
			}

			continue
		}

		// Named field.

		for _, name := range field.Names {
			if !name.IsExported() {
				// Ignore unexported fields.
				continue
			}

			if _, present := result[name.Name]; present {
				// Ignore redefinitions.
				continue
			}

			result[name.Name] = &FieldDefInfo{
				Pos:   name.Pos(),
				Doc:   doc,
				Group: group,
			}
		}
	}

	return result
}

// Rhs computes the immediate denotation of a type definition,
// given the object that represents that definition.
//
// This is unfortunate but necessary, because the Go type checker
// remembers only the underlying type,
// forgetting the actual chain of aliases and named types.
//
// If the given object does not match the receiver,
// or the definition has incorrect syntax, Rhs returns nil.
func (info *TypeDefInfo) Rhs(obj *types.TypeName) types.Type {
	if obj == nil || obj.Name() != info.Name || obj.Pkg() == nil || obj.Pkg().Path() != info.pkg.Path {
		// Invalid object.
		return nil
	}

	// Parse definition.
	def := info.def
	var ident, pkgName *ast.Ident

	// Unwrap generic type instantiations.
	switch d := def.(type) {
	case *ast.IndexExpr:
		def = d.X
	case *ast.IndexListExpr:
		def = d.X
	}

	switch d := def.(type) {
	case *ast.Ident:
		// We have to compute the denotation for an identifier.
		ident = d
	case *ast.SelectorExpr:
		// We have to compute the denotation for a qualified identifier.
		ident = d.Sel
		pkgName, _ = d.X.(*ast.Ident)
		if pkgName == nil {
			// Invalid syntax.
			return nil
		}
	default:
		// The underlying type _is_ the immediate denotation
		return obj.Type().Underlying()
	}

	if ident == nil || ident.Name == "" {
		// Invalid syntax.
		return nil
	}

	if pkgName == nil {
		// Attempt package-local lookup first.
		if _, rhs := obj.Pkg().Scope().LookupParent(ident.Name, token.NoPos); rhs != nil {
			if _, isType := rhs.(*types.TypeName); !isType {
				// Object is not a type def.
				return nil
			}

			return rhs.Type()
		}
	}

	// We have to look into other packages.

	if !ident.IsExported() {
		// Unexported identifiers cannot be provided by other packages.
		return nil
	}

	// Acquire cache mutex.
	info.mu.Lock()
	defer info.mu.Unlock()

	// If we cached lookup results for the same object, return them.
	if info.obj == obj {
		return info.rhs
	}

	// Reset cache.
	info.obj = obj
	info.rhs = nil

	// Lookup object among dependencies.
	// No need to call info.pkg.Collect: the receiver has been created there.
	rhs := resolveImportedObject(obj.Pkg().Imports(), info.pkg.Imports[info.file], pkgName.Name, ident.Name)

	if _, isType := rhs.(*types.TypeName); !isType {
		// Object is not a type def.
		return nil
	}

	info.rhs = rhs.Type()
	return info.rhs
}

// resolveImportedObject resolves a file-local qualified or unqualified name
// to the object it denotes. If no matching package or object is found,
// resolveImportedObject returns nil.
func resolveImportedObject(imports []*types.Package, fileImports *FileImports, pkgName string, objectName string) types.Object {
	var paths []string

	if pkgName == "" {
		paths = fileImports.Dot
	} else if path, ok := fileImports.Named[pkgName]; ok {
		paths = []string{path}
	} else {
		paths = fileImports.Unnamed
	}

	if len(paths) == 0 {
		// No matching package.
		return nil
	}

	// Visit candidate packages.
	// TODO: Can we do any better than linear search?

	if len(paths) == 1 {
		// Look for a package with a specific path.
		for _, pkg := range imports {
			if pkg.Path() == paths[0] {
				return pkg.Scope().Lookup(objectName)
			}
		}

		return nil
	}

	// Look for a package with matching package name or object,
	// preferring earlier ones in source order.
	var earliest types.Object

	for _, pkg := range imports {
		// Match path against candidates.
		i, ok := slices.BinarySearch(paths, pkg.Path())
		if !ok {
			continue
		}

		// Candidate package found.
		// If we have no package name to match, lookup object.
		var obj types.Object
		if pkgName == "" {
			obj = pkg.Scope().Lookup(objectName)
		}

		if obj != nil {
			// Object match: shrink search space.
			earliest = obj
			paths = paths[:i]
		} else if pkgName != "" && pkg.Name() == pkgName {
			// Package name match: lookup object and shrink search space.
			earliest = pkg.Scope().Lookup(objectName)
			paths = paths[:i]
		}

		if len(paths) == 0 {
			// We've tested all relevant packages.
			break
		}
	}

	return earliest
}
