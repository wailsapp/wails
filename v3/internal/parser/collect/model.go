package collect

import (
	"cmp"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"slices"
	"strings"
	"sync"
)

type (
	// ModelInfo records all information that is required
	// to render JS/TS code for a model type.
	//
	// Read accesses to exported fields are only safe
	// if a call to [ModelInfo.Collect] has been completed before the access,
	// for example by calling it in the accessing goroutine
	// or before spawning the accessing goroutine.
	ModelInfo struct {
		*TypeDefInfo

		// Imports records dependencies for this model.
		Imports *ImportMap

		// Type records the target type for an alias or derived model,
		// the underlying type for an enum.
		Type types.Type

		// Fields records the property list for a class or struct alias model,
		// in order of declaration and grouped by their declaring [ast.Field].
		Fields [][]*ModelFieldInfo

		// Values records the value list for an enum model,
		// in order of declaration and grouped
		// by their declaring [ast.GenDecl] and [ast.ValueSpec].
		Values [][][]*EnumValueInfo

		// TypeParams records type parameter names for generic models.
		TypeParams []string

		obj  *types.TypeName
		pkg  *PackageInfo
		once sync.Once
	}

	ModelFieldInfo struct {
		*FieldInfo
		*FieldDefInfo
	}

	EnumValueInfo struct {
		*ConstInfo
		Value any
	}
)

// Model retrieves the the unique [ModelInfo] instance
// associated to the given model type within a Collector.
// If none is present, a new one is initialised.
//
// If the model's declaring package fails to load, Model returns nil.
// Errors are reported through the controller associated to the collector.
//
// Model is safe for concurrent use.
func (collector *Collector) Model(typ *types.TypeName) *ModelInfo {
	return collector.Package(typ.Pkg().Path()).recordModel(typ)
}

// Collect gathers information for the model described by its receiver.
// It can be called concurrently by multiple goroutines;
// the computation will be performed just once.
//
// After Collect returns, the calling goroutine and all goroutines
// it might spawn afterwards are free to access
// the receiver's fields indefinitely.
func (info *ModelInfo) Collect() {
	// Changes in the following logic must be reflected adequately
	// by the predicates in properties.go, by [ImportMap.AddType]
	// and by [render.RenderType].

	info.once.Do(func() {
		pkg := info.pkg
		pkg.Collect()

		obj := info.obj
		typ := obj.Type()

		// Retrieve type def information.
		info.TypeDefInfo = pkg.Types[obj.Name()]

		// Check type def information.
		if info.TypeDefInfo == nil {
			pkg.collector.controller.Errorf(
				"package %s: type %s not found; try cleaning the build cache (go clean -cache)",
				pkg.Path,
				obj.Name(),
			)
			return
		}

		// Initialise import map.
		info.Imports = NewImportMap(pkg)

		// Setup fallback type.
		info.Type = types.Universe.Lookup("any").Type()

		// Check marshalers and detect enums.
		var constants []*types.Const

		// Retrieve type denotation.
		var def types.Type
		if obj.IsAlias() {
			def = info.TypeDefInfo.Rhs(obj)
			if def == nil {
				def = types.Unalias(typ)
			}
		} else {
			// This is a named type.
			// Check whether it implements marshaler interfaces
			// or has defined constants.

			if IsAny(typ) {
				// Type marshals to a custom value of unknown shape.
				return
			} else if MaybeTextMarshaler(typ) {
				// Type marshals to a custom string of unknown shape.
				info.Type = types.Typ[types.String]
				return
			}

			// Store type parameter names.
			tp := typ.(*types.Named).TypeParams()
			if tp != nil && tp.Len() > 0 {
				info.TypeParams = make([]string, tp.Len())
				for i, length := 0, tp.Len(); i < length; i++ {
					info.TypeParams[i] = tp.At(i).Obj().Name()
				}
			}

			def = info.TypeDefInfo.Rhs(obj)
			if def == nil {
				def = typ.Underlying()
			}
			if IsClass(typ) {
				// Skip alias chains for classes.
				def = types.Unalias(def)
			}

			// Test for enums (excluding generic types).
			basic, ok := typ.Underlying().(*types.Basic)
			if ok && tp == nil && basic.Info()&types.IsConstType != 0 && basic.Info()&types.IsComplex == 0 {
				// Named type is defined as a representable constant type:
				// look for defined constants of that named type.
				for _, name := range obj.Pkg().Scope().Names() {
					if cnst, ok := obj.Pkg().Scope().Lookup(name).(*types.Const); ok {
						if cnst.Val().Kind() != constant.Unknown && types.Identical(cnst.Type(), typ) {
							constants = append(constants, cnst)
						}
					}
				}
			}
		}

		// Record required imports.
		info.Imports.AddType(def, pkg.collector)

		// Handle enum types.
		// constants slice is always empty for aliases.
		if len(constants) > 0 {
			// Collect information about enum values.
			info.collectEnum(constants)
			info.Type = def
			return
		}

		// Test for structs.
		strct, isStruct := def.(*types.Struct)
		if !isStruct {
			// That's all, folks. Render as a TS alias.
			info.Type = def
			return
		}

		// Collect information about struct fields.
		info.collectStruct(strct)
		info.Type = nil

		info.pkg, info.obj = nil, nil
	})
}

// collectEnum collects information about enum values and their declarations.
func (info *ModelInfo) collectEnum(constants []*types.Const) {
	pkg := info.pkg
	dummyGroup := &GroupInfo{
		Group: &GroupInfo{},
	}

	names := make(map[string]bool, len(constants))
	values := make([]*EnumValueInfo, len(constants))

	for i, cnst := range constants {
		names[cnst.Name()] = true
		value := &EnumValueInfo{
			ConstInfo: pkg.Consts[cnst.Name()],
			Value:     constant.Val(cnst.Val()),
		}

		if value.ConstInfo == nil {
			value.ConstInfo = &ConstInfo{
				Name:  cnst.Name(),
				Group: dummyGroup,
			}
			pkg.collector.controller.Warningf(
				"package %s: could not retrieve definition for constant %s; try cleaning the build cache (go clean -cache)",
				pkg.Path,
				cnst.Name(),
			)
		}

		values[i] = value
	}

	// Sort values by grouping and source order.
	slices.SortFunc(values, func(v1 *EnumValueInfo, v2 *EnumValueInfo) int {
		// Sort first by source order of declaration group.
		if g1, g2 := v1.Group.Group, v2.Group.Group; g1 != g2 {
			return cmp.Compare(g1.Pos, g2.Pos)
		}

		// Then by source order of spec.
		if sg1, sg2 := v1.Group, v2.Group; sg1 != sg2 {
			return cmp.Compare(sg1.Pos, sg2.Pos)
		}

		// Then by source order of identifiers.
		if v1.Pos != v2.Pos {
			return cmp.Compare(v1.Pos, v2.Pos)
		}

		// Finally by name (for constants whose source position is unknown).
		return strings.Compare(v1.Name, v2.Name)
	})

	// Split value list into groups and subgroups.
	var decl, spec *GroupInfo
	decli, speci := -1, -1

	for _, value := range values {
		if value.Group != spec {
			spec = value.Group

			if spec.Group == decl {
				speci++
			} else {
				decl = spec.Group
				decli++
				speci = 0
				info.Values = append(info.Values, nil)
			}

			info.Values[decli] = append(info.Values[decli], nil)
		}

		info.Values[decli][speci] = append(info.Values[decli][speci], value)
	}
}

// collectStruct collects information about struct fields and their declarations.
func (info *ModelInfo) collectStruct(strct *types.Struct) {
	pkg := info.pkg
	dummyFieldDef := &FieldDefInfo{
		Group: &GroupInfo{},
	}

	// Retrieve struct info.
	structInfo := pkg.collector.Struct(strct)
	structInfo.Collect()

	// Cache resolved TypeDefInfo for embedded struct types.
	rootTypeInfo := info.resolveTypeInfo(info.obj)
	if rootTypeInfo == nil {
		pkg.collector.controller.Warningf(
			"package %s: could not resolve definition for type %s; try cleaning the build cache (go clean -cache)",
			info.pkg.Path,
			info.Name,
		)
	}

	embeddedInfo := map[*types.TypeName]*TypeDefInfo{
		nil:      rootTypeInfo,
		info.obj: rootTypeInfo,
	}

	// Allocate result slice.
	fields := make([]*ModelFieldInfo, len(structInfo.Fields))

	// Collect fields.
	for i, field := range structInfo.Fields {
		mfield := &ModelFieldInfo{
			FieldInfo:    field,
			FieldDefInfo: dummyFieldDef,
		}

		// Lookup field definition.
		typeInfo, ok := embeddedInfo[field.Parent]
		if !ok {
			// Resolve and cache.
			typeInfo = info.resolveTypeInfo(field.Parent)
			embeddedInfo[field.Parent] = typeInfo

			// Report errors
			if typeInfo == nil {
				pkg.collector.controller.Warningf(
					"package %s: could not resolve definition for type %s; try cleaning the build cache (go clean -cache)",
					field.Parent.Pkg().Path(),
					field.Parent.Name(),
				)
			}
		}

		if typeInfo != nil {
			mfield.FieldDefInfo = typeInfo.Fields()[field.Field.Name()]
		}

		if mfield.FieldDefInfo == nil {
			mfield.FieldDefInfo = dummyFieldDef

			parent := field.Parent
			if parent == nil {
				parent = info.obj
			}
			pkg.collector.controller.Warningf(
				"package %s: type %s: could not retrieve definition for field %s; loading package %s explicitly might solve the issue",
				parent.Pkg().Path(),
				parent.Name(),
				field.Field.Name(),
				parent.Pkg().Path(),
			)
		}

		fields[i] = mfield
	}

	// Split field list into groups, preserving the original order.
	var decl *GroupInfo
	decli := -1

	for _, field := range fields {
		if field.Group != decl {
			decl = field.Group
			decli++
			info.Fields = append(info.Fields, nil)
		}

		info.Fields[decli] = append(info.Fields[decli], field)
	}
}

// resolveTypeInfo follows the alias/named type chain
// for the given defined type to find the syntax
// that defines its fields.
//
// It returns nil on failure.
func (info *ModelInfo) resolveTypeInfo(obj *types.TypeName) *TypeDefInfo {
	pkg := info.pkg

	if obj == nil {
		obj = info.obj
	}

	for {
		// Fast path for aliases of named types.
		if obj.IsAlias() {
			if named, ok := types.Unalias(obj.Type()).(*types.Named); ok {
				obj = named.Obj()
			}
		}

		var typeInfo *TypeDefInfo
		if obj == info.obj {
			typeInfo = info.TypeDefInfo
		} else {
			tpkg := pkg.collector.Package(obj.Pkg().Path())
			if tpkg.Collect() {
				typeInfo = tpkg.Types[obj.Name()]
			}
		}

		if typeInfo == nil {
			// Lookup failed.
			return nil
		}

		// Follow aliases and named types, stop when there are no more.
		switch rhs := typeInfo.Rhs(obj).(type) {
		case *types.Alias:
			obj = rhs.Obj()
		case *types.Named:
			obj = rhs.Obj()
		case nil:
			// Lookup failed.
			// One last desperate attempt:
			// we might be dealing with an unexported named type.
			def := typeInfo.def

			// Unwrap generic type instantiations.
			switch d := def.(type) {
			case *ast.IndexExpr:
				def = d.X
			case *ast.IndexListExpr:
				def = d.X
			}

			ident, ok := def.(*ast.Ident)
			if !ok || ident.IsExported() {
				// We can't do anything more: give up.
				return typeInfo
			}

			// Feed a fake object to the next iteration.
			fake := types.NewTypeName(token.NoPos, obj.Pkg(), ident.Name, nil)
			types.NewNamed(fake, obj.Type().Underlying(), nil)
			obj = fake
		default:
			return typeInfo
		}
	}
}
