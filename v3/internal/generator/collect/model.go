package collect

import (
	"cmp"
	"go/ast"
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

		// Internal records whether the model
		// should be exported by the index file.
		Internal bool

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

		// Predicates caches the value of all type predicates for this model.
		//
		// WARN: whenever working with a generic uninstantiated model type,
		// use these instead of invoking predicate functions,
		// which may incur a large performance penalty.
		Predicates Predicates

		collector *Collector
		once      sync.Once
	}

	// ModelFieldInfo holds extended information
	// about a struct field in a model type.
	ModelFieldInfo struct {
		*StructField
		*FieldInfo
	}

	// Predicates caches the value of all type predicates.
	Predicates struct {
		IsJSONMarshaler    MarshalerKind
		MaybeJSONMarshaler MarshalerKind
		IsTextMarshaler    MarshalerKind
		MaybeTextMarshaler MarshalerKind
		IsMapKey           bool
		IsTypeParam        bool
		IsStringAlias      bool
		IsClass            bool
		IsAny              bool
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

		// Record whether the model should be exported.
		info.Internal = !obj.Exported()

		// Parse directives.
		for _, doc := range []*ast.CommentGroup{info.Doc, info.Decl.Doc} {
			if doc == nil {
				continue
			}
			for _, comment := range doc.List {
				if IsDirective(comment.Text, "internal") {
					info.Internal = true
				}
			}
		}

		// Record type parameter names.
		var isGeneric bool
		if generic, ok := typ.(interface{ TypeParams() *types.TypeParamList }); ok {
			tparams := generic.TypeParams()
			isGeneric = tparams != nil

			if isGeneric && tparams.Len() > 0 {
				info.TypeParams = make([]string, tparams.Len())
				for i := range tparams.Len() {
					info.TypeParams[i] = tparams.At(i).Obj().Name()
				}
			}
		}

		// Precompute predicates.
		// Preinstantiate typ to avoid repeated instantiations in predicate code.
		ityp := instantiate(typ)
		info.Predicates = Predicates{
			IsJSONMarshaler:    IsJSONMarshaler(ityp),
			MaybeJSONMarshaler: MaybeJSONMarshaler(ityp),
			IsTextMarshaler:    IsTextMarshaler(ityp),
			MaybeTextMarshaler: MaybeTextMarshaler(ityp),
			IsMapKey:           IsMapKey(ityp),
			IsTypeParam:        IsTypeParam(ityp),
			IsStringAlias:      IsStringAlias(ityp),
			IsClass:            IsClass(ityp),
			IsAny:              IsAny(ityp),
		}

		var def types.Type
		var constants []*types.Const

		switch t := typ.(type) {
		case *types.Alias:
			// Model is an alias: take rhs as definition.
			// It is important not to skip alias chains with [types.Unalias]
			// because in doing so we could end up with a private type from another package.
			def = t.Rhs()

			// Test for constants with alias type,
			// but only when non-generic alias resolves to a basic type
			// (hence not to e.g. a named type).
			if basic, ok := types.Unalias(def).(*types.Basic); ok {
				if !isGeneric && basic.Info()&types.IsConstType != 0 && basic.Info()&types.IsComplex == 0 {
					// Non-generic alias resolves to a representable constant type:
					// look for defined constants whose type is exactly the alias typ.
					for _, name := range obj.Pkg().Scope().Names() {
						if cnst, ok := obj.Pkg().Scope().Lookup(name).(*types.Const); ok {
							alias, isAlias := cnst.Type().(*types.Alias)
							if isAlias && cnst.Val().Kind() != constant.Unknown && alias.Obj() == t.Obj() {
								constants = append(constants, cnst)
							}
						}
					}
				}
			}

		case *types.Named:
			// Model is a named type:
			// jump directly to underlying type to match go semantics,
			// i.e. do not render named types as aliases for other named types.
			def = typ.Underlying()

			// Check whether it implements marshaler interfaces or has defined constants.
			if info.Predicates.MaybeJSONMarshaler != NonMarshaler {
				// Type marshals to a custom value of unknown shape.
				// If it has explicit custom marshaling logic, render it as any;
				// otherwise, delegate to the underlying type that must be the actual [json.Marshaler].
				if info.Predicates.MaybeJSONMarshaler == ExplicitMarshaler {
					return
				}
			} else if info.Predicates.MaybeTextMarshaler != NonMarshaler {
				// Type marshals to a custom string of unknown shape.
				// If it has explicit custom marshaling logic, render it as string;
				// otherwise, delegate to the underlying type that must be the actual [encoding.TextMarshaler].
				//
				// One exception must be made for situations
				// where the underlying type is a [json.Marshaler] but the model is not:
				// in that case, we cannot delegate to the underlying type either.
				// Note that in such a case the underlying type is never a pointer or interface,
				// because those cannot have explicitly defined methods,
				// hence it would not possible for the model not to be a [json.Marshaler]
				// while the underlying type is.
				if info.Predicates.MaybeTextMarshaler == ExplicitMarshaler || MaybeJSONMarshaler(def) != NonMarshaler {
					info.Type = types.Typ[types.String]
					return
				}
			} else if basic, ok := def.Underlying().(*types.Basic); ok {
				// Test for enums (excluding marshalers and generic types).
				if !isGeneric && basic.Info()&types.IsConstType != 0 && basic.Info()&types.IsComplex == 0 {
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

		default:
			panic("model has unknown object kind (neither alias nor named type)")
		}

		// Handle struct types.
		strct, isStruct := def.(*types.Struct)
		if isStruct && info.Predicates.MaybeJSONMarshaler == NonMarshaler && info.Predicates.MaybeTextMarshaler == NonMarshaler {
			// Def is struct and model is not a marshaler:
			// collect information about struct fields.
			info.collectStruct(strct)
			info.Type = nil
			return
		}

		// Record required imports.
		info.Imports.AddType(def)

		// Handle enum types.
		// constants slice is always empty for structs, marshalers.
		if len(constants) > 0 {
			// Collect information about enum values.
			info.collectEnum(constants)
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
		// Record required imports.
		info.Imports.AddType(field.Type)

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
