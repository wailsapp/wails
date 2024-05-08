package collect

import (
	"go/ast"
	"go/token"
	"go/types"
	"slices"
	"sync"

	"golang.org/x/tools/go/packages"
)

type (
	// TypeDefInfo records information about a single type declaration.
	TypeDefInfo struct {
		Name    string
		Doc     *ast.CommentGroup
		Group   *GroupInfo
		Alias   bool
		Methods map[string]*MethodInfo

		// def stores the type's defining syntax.
		def ast.Expr
		// rhs caches def's denotation.
		rhs types.Type
		// fields caches information about a struct type's field definitions.
		fields map[string]*FieldDefInfo

		file       *FileInfo
		onceFields sync.Once
		onceRhs    sync.Once
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
func newTypeDefInfo(pkg *packages.Package, file *FileInfo, group *GroupInfo, spec *ast.TypeSpec) *TypeDefInfo {
	info := &TypeDefInfo{
		Name:  spec.Name.Name,
		Doc:   spec.Doc,
		Group: group,
		Alias: spec.Assign.IsValid(),

		def: spec.Type,
	}

	// Retrieve rhs or store required information.
	if pkg.TypesInfo != nil {
		info.rhs = pkg.TypesInfo.TypeOf(spec.Type)
	} else {
		info.file = file
	}

	return info
}

// Fields returns information about the fields of a struct type.
// If the receiver does not describe a struct type, Fields returns nil.
func (info *TypeDefInfo) Fields() map[string]*FieldDefInfo {
	info.onceFields.Do(func() {
		def, ok := info.def.(*ast.StructType)
		if !ok {
			return
		}

		info.fields = make(map[string]*FieldDefInfo)

		for _, field := range def.Fields.List {
			doc := field.Doc
			if doc == nil {
				doc = field.Comment
			} else if field.Comment != nil {
				doc = &ast.CommentGroup{
					List: slices.Concat(field.Doc.List, field.Comment.List),
				}
			}

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

				if _, present := info.fields[name.Name]; present {
					// Ignore redefinitions.
					continue
				}

				info.fields[name.Name] = &FieldDefInfo{
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

				if _, present := info.fields[name.Name]; present {
					// Ignore redefinitions.
					continue
				}

				info.fields[name.Name] = &FieldDefInfo{
					Pos:   name.Pos(),
					Doc:   doc,
					Group: group,
				}
			}
		}
	})

	return info.fields
}

// Rhs returns the immediate denotation of a type definition,
// given the object that represents that definition.
//
// This is unfortunate but necessary, because the Go type checker
// remembers only the underlying type, forgetting the actual chain
// of aliases and named types that lead there.
//
// If the given object does not match the receiver,
// or the definition has incorrect syntax, Rhs returns nil.
func (info *TypeDefInfo) Rhs(obj *types.TypeName) types.Type {
	if obj == nil || obj.Name() != info.Name || obj.Pkg() == nil {
		// Invalid object.
		return nil
	}

	switch info.def.(type) {
	case *ast.Ident, *ast.SelectorExpr, *ast.IndexExpr, *ast.IndexListExpr:
		// Definition is either a qualified identifier,
		// or a generic type instantiation.
		// Resolve it in file scope.
	default:
		// The underlying type _is_ the immediate denotation
		return obj.Type().Underlying()
	}

	info.onceRhs.Do(func() {
		if info.rhs != nil || info.file == nil {
			return
		}

		// When dealing with generic types we have to create
		// an additional scope for typeparam lookup.
		if named, ok := obj.Type().(*types.Named); ok && named.TypeParams() != nil && named.TypeParams().Len() > 0 {
			fscope := info.file.Scope(obj.Pkg())
			if fscope.Innermost(info.def.Pos()) == fscope {
				tscope := types.NewScope(fscope, info.def.Pos(), info.def.End(), "")
				for i, length := 0, named.TypeParams().Len(); i < length; i++ {
					tscope.Insert(named.TypeParams().At(i).Obj())
				}
			}
		}

		info.rhs = info.file.TypeOf(obj.Pkg(), info.def)
		info.file = nil
	})

	return info.rhs
}
