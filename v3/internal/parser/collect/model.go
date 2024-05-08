package collect

import (
	"cmp"
	"go/constant"
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
	// if a call to [ModelInfo.Collect] has completed before the access,
	// for example by calling it in the accessing goroutine
	// or before spawning the accessing goroutine.
	ModelInfo struct {
		*TypeInfo

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
		Values [][][]*ConstInfo

		// TypeParams records type parameter names for generic models.
		TypeParams []string

		collector *Collector
		once      sync.Once
	}

	// ModelFieldInfo holds extended information
	// about a struct field in a model type.
	ModelFieldInfo struct {
		*StructField
		*FieldInfo
	}
)

func newModelInfo(collector *Collector, obj *types.TypeName) *ModelInfo {
	return &ModelInfo{
		TypeInfo:  collector.Type(obj),
		collector: collector,
	}
}

// Model retrieves the the unique [ModelInfo] instance
// associated to the given type object within a Collector.
// If none is present, Model initialises a new one
// registers it for code generation
// and schedules background collection activity.
//
// Model is safe for concurrent use.
func (collector *Collector) Model(obj *types.TypeName) *ModelInfo {
	pkg := collector.Package(obj.Pkg())
	if pkg == nil {
		return nil
	}

	model, present := pkg.recordModel(obj)
	if !present {
		collector.scheduler.Schedule(func() { model.Collect() })
	}

	return model
}

// Collect gathers information for the model described by its receiver.
// It can be called concurrently by multiple goroutines;
// the computation will be performed just once.
//
// Collect returns the receiver for chaining.
// It is safe to call Collect with nil receiver.
//
// After Collect returns, the calling goroutine and all goroutines
// it might spawn afterwards are free to access
// the receiver's fields indefinitely.
func (info *ModelInfo) Collect() *ModelInfo {
	if info == nil {
		return nil
	}

	// Changes in the following logic must be reflected adequately
	// by the predicates in properties.go, by ImportMap.AddType
	// and by all render.Module methods.

	info.once.Do(func() {
		collector := info.collector
		obj := info.Object().(*types.TypeName)

		typ := obj.Type()

		// Collect type information.
		info.TypeInfo.Collect()

		// Initialise import map.
		info.Imports = NewImportMap(collector.Package(obj.Pkg()))

		// Setup fallback type.
		info.Type = types.Universe.Lookup("any").Type()

		// Retrieve type denotation and skip alias chains.
		def := info.TypeInfo.Def

		// Check marshalers and detect enums.
		var constants []*types.Const

		var isGeneric bool
		if generic, ok := obj.Type().(interface{ TypeParams() *types.TypeParamList }); ok {
			// Record type parameter names.
			tparams := generic.TypeParams()
			isGeneric = tparams != nil

			if isGeneric && tparams.Len() > 0 {
				info.TypeParams = make([]string, tparams.Len())
				for i := range tparams.Len() {
					info.TypeParams[i] = tparams.At(i).Obj().Name()
				}
			}
		}

		if _, isNamed := obj.Type().(*types.Named); isNamed {
			// Model is a named type.
			// Check whether it implements marshaler interfaces
			// or has defined constants.

			if IsAny(typ) {
				// Type marshals to a custom value of unknown shape.
				return
			} else if MaybeTextMarshaler(typ) {
				// Type marshals to a custom string of unknown shape.
				info.Type = types.Typ[types.String]
				return
			} else if isGeneric && !collector.options.UseInterfaces && IsClass(typ) {
				// Generic classes cannot be defined in terms of other generic classes.
				// That would break class creation code,
				// and I (@fbbdev) couldn't find any other satisfying workaround.
				def = typ.Underlying()
			}

			// Test for enums (excluding generic types).
			basic, ok := typ.Underlying().(*types.Basic)
			if ok && !isGeneric && basic.Info()&types.IsConstType != 0 && basic.Info()&types.IsComplex == 0 {
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
		info.Imports.AddType(def)

		// Handle enum types.
		// constants slice is always empty for aliases.
		if len(constants) > 0 {
			// Collect information about enum values.
			info.collectEnum(constants)
			info.Type = def
			return
		}

		// Handle struct types.
		strct, isStruct := def.(*types.Struct)
		if isStruct {
			// Collect information about struct fields.
			info.collectStruct(strct)
			info.Type = nil
			return
		}

		// That's all, folks. Render as a TS alias.
		info.Type = def
	})

	return info
}

// collectEnum collects information about enum values and their declarations.
func (info *ModelInfo) collectEnum(constants []*types.Const) {
	// Collect information about each constant object.
	values := make([]*ConstInfo, len(constants))
	for i, cnst := range constants {
		values[i] = info.collector.Const(cnst).Collect()
	}

	// Sort values by grouping and source order.
	slices.SortFunc(values, func(v1 *ConstInfo, v2 *ConstInfo) int {
		// Skip comparisons for identical pointers.
		if v1 == v2 {
			return 0
		}

		// Sort first by source order of declaration group.
		if v1.Decl != v2.Decl {
			return cmp.Compare(v1.Decl.Pos, v2.Decl.Pos)
		}

		// Then by source order of spec.
		if v1.Spec != v2.Spec {
			return cmp.Compare(v1.Spec.Pos, v2.Spec.Pos)
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
		if value.Spec != spec {
			spec = value.Spec

			if value.Decl == decl {
				speci++
			} else {
				decl = value.Decl
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
	collector := info.collector

	// Retrieve struct info.
	structInfo := collector.Struct(strct).Collect()

	// Allocate result slice.
	fields := make([]*ModelFieldInfo, len(structInfo.Fields))

	// Collect fields.
	for i, field := range structInfo.Fields {
		fields[i] = &ModelFieldInfo{
			StructField: field,
			FieldInfo:   collector.Field(field.Object).Collect(),
		}
	}

	// Split field list into groups, preserving the original order.
	var decl *GroupInfo
	decli := -1

	for _, field := range fields {
		if field.Decl != decl {
			decl = field.Decl
			decli++
			info.Fields = append(info.Fields, nil)
		}

		info.Fields[decli] = append(info.Fields[decli], field)
	}
}
