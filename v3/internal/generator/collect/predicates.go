package collect

// This file gathers functions that test useful properties of model types.
// The rationale for the way things are handled here
// is given in the example file found at ./_reference/json_marshaler_behaviour.go

import (
	"go/token"
	"go/types"
	"iter"

	"golang.org/x/exp/typeparams"
)

// Cached interface types.
var (
	ifaceTextMarshaler = types.NewInterfaceType([]*types.Func{
		types.NewFunc(token.NoPos, nil, "MarshalText",
			types.NewSignatureType(nil, nil, nil, types.NewTuple(), types.NewTuple(
				types.NewParam(token.NoPos, nil, "", types.NewSlice(types.Universe.Lookup("byte").Type())),
				types.NewParam(token.NoPos, nil, "", types.Universe.Lookup("error").Type()),
			), false)),
	}, nil).Complete()

	ifaceJSONMarshaler = types.NewInterfaceType([]*types.Func{
		types.NewFunc(token.NoPos, nil, "MarshalJSON",
			types.NewSignatureType(nil, nil, nil, types.NewTuple(), types.NewTuple(
				types.NewParam(token.NoPos, nil, "", types.NewSlice(types.Universe.Lookup("byte").Type())),
				types.NewParam(token.NoPos, nil, "", types.Universe.Lookup("error").Type()),
			), false)),
	}, nil).Complete()
)

// MarshalerKind values describe
// whether and how a type implements a marshaler interface.
// For any one of the two marshaler interfaces, a type is
//   - a NonMarshaler if it does not implement it;
//   - an ImplicitMarshaler if it inherits the implementation from its underlying type;
//   - an ExplicitMarshaler if it defines the relevant method explicitly.
type MarshalerKind byte

const (
	NonMarshaler MarshalerKind = iota
	ImplicitMarshaler
	ExplicitMarshaler
)

// termlist returns an iterator over the normalised term list of the given type.
// If typ is invalid or has an empty type set, termlist returns the empty sequence.
// If typ has an empty term list
// then termlist returns a sequence with just one element: the type itself.
//
// TODO: replace with new term set API once Go 1.25 is out.
// See go.dev/issue/61013
func termlist(typ types.Type) iter.Seq[*typeparams.Term] {
	terms, err := typeparams.NormalTerms(types.Unalias(typ))
	return func(yield func(*typeparams.Term) bool) {
		if err == nil && len(terms) == 0 {
			yield(typeparams.NewTerm(false, typ))
		} else {
			for _, term := range terms {
				if !yield(term) {
					break
				}
			}
		}
	}
}

// instantiate instantiates typ if it is an uninstantiated generic type
// using its own type parameters as arguments in order to preserve genericity.
//
// If typ is not generic or already instantiated, it is returned as is.
// If typ is not an alias, then the returned type is not an alias either.
func instantiate(typ types.Type) types.Type {
	if t, ok := typ.(interface {
		TypeParams() *types.TypeParamList
		TypeArgs() *types.TypeList
	}); ok && t.TypeParams() != nil && t.TypeArgs() == nil {
		args := make([]types.Type, t.TypeParams().Len())
		for i := range args {
			args[i] = t.TypeParams().At(i)
		}

		typ, _ = types.Instantiate(nil, typ, args, false)
	}

	return typ
}

// isMarshaler checks whether the given type
// implements one of the two marshaler interfaces,
// and whether it implements it explicitly,
// i.e. by defining the relevant method directly
// instead of inheriting it from the underlying type.
//
// If addressable is true, it checks both pointer and non-pointer receivers.
//
// The behaviour of isMarshaler is unspecified
// if marshaler is not one of [json.Marshaler] or [encoding.TextMarshaler].
func isMarshaler(typ types.Type, marshaler *types.Interface, addressable bool, visited map[*types.TypeName]MarshalerKind) MarshalerKind {
	// Follow alias chain and instantiate if necessary.
	//
	// types.Implements does not handle generics,
	// hence when typ is generic it must be instantiated.
	//
	// Instantiation operations may incur a large performance penalty and are usually cached,
	// but doing so here would entail some complex global state and a potential memory leak.
	// Because typ should be generic only during model collection,
	// it should be enough to cache the result of marshaler queries for models.
	typ = instantiate(types.Unalias(typ))

	// Invariant: at this point, typ is not an alias.

	if typ == types.Typ[types.Invalid] {
		// Do not pass invalid types to [types.Implements].
		return NonMarshaler
	}

	result := types.Implements(typ, marshaler)

	ptr, isPtr := typ.Underlying().(*types.Pointer)

	if !result && addressable && !isPtr {
		result = types.Implements(types.NewPointer(typ), marshaler)
	}

	named, isNamed := typ.(*types.Named)

	if result {
		// Check whether marshaler method is implemented explicitly on a named type.
		if isNamed {
			method := marshaler.Method(0).Name()
			for i := range named.NumMethods() {
				if named.Method(i).Name() == method {
					return ExplicitMarshaler
				}
			}
		}

		return ImplicitMarshaler
	}

	// Fast path: named types that fail the [types.Implements] test cannot be marshalers.
	//
	// WARN: currently typeparams cannot be used on the rhs of a named type declaration.
	// If that changes in the future,
	// this guard will become essential for correctness,
	// not just a shortcut.
	if isNamed {
		return NonMarshaler
	}

	// Unwrap at most one pointer and follow alias chain.
	if isPtr {
		typ = types.Unalias(ptr.Elem())
	}

	// Invariant: at this point, typ is not an alias.

	// Type parameters require special handling:
	// iterate over their term list and treat them as marshalers
	// if so are all their potential instantiations.

	tp, ok := typ.(*types.TypeParam)
	if !ok {
		return NonMarshaler
	}

	// Init cycle detection/deduplication map.
	if visited == nil {
		visited = make(map[*types.TypeName]MarshalerKind)
	}

	// Type params cannot be embedded in constraints directly,
	// but they can be embedded as pointer terms.
	//
	// When we hit that kind of cycle,
	// we can err towards it being a marshaler:
	// such a constraint is meaningless anyways,
	// as no type can be simultaneously a pointer to itself.
	//
	// Therefore, we iterate the type set
	// only for unvisited pointers-to-typeparams,
	// and return the current best guess
	// for those we have already visited.
	//
	// WARN: there has been some talk
	// of allowing type parameters as embedded fields/terms.
	// That might make our lives miserable here.
	// The spec must be monitored for changes in that regard.
	if isPtr {
		if kind, ok := visited[tp.Obj()]; ok {
			return kind
		}
	}

	// Initialise kind to explicit marshaler, then decrease as needed.
	kind := ExplicitMarshaler

	if isPtr {
		// Pointers are never explicit marshalers.
		kind = ImplicitMarshaler
		// Mark pointer-to-typeparam as visited and init current best guess.
		visited[tp.Obj()] = kind
	}

	// Iterate term list.
	for term := range termlist(tp) {
		ttyp := types.Unalias(term.Type())

		// Reject if tp has a tilde or invalid element in its term list
		// or has a method-only constraint.
		//
		// Valid tilde terms
		// can always be satisfied by named types that hide their methods
		// hence fail in general to implement the required interface.
		if term.Tilde() || ttyp == types.Typ[types.Invalid] || ttyp == tp {
			kind = NonMarshaler
			break
		}

		// Propagate the presence of a wrapping pointer.
		if isPtr {
			ttyp = types.NewPointer(ttyp)
		}

		kind = min(kind, isMarshaler(ttyp, marshaler, addressable && !isPtr, visited))
		if kind == NonMarshaler {
			// We can stop here as we've reached the minimum [MarshalerKind].
			break
		}
	}

	// Store final response for pointer-to-typeparam.
	if isPtr {
		visited[tp.Obj()] = kind
	}

	return kind
}

// IsTextMarshaler queries whether and how the given type
// implements the [encoding.TextMarshaler] interface.
func IsTextMarshaler(typ types.Type) MarshalerKind {
	return isMarshaler(typ, ifaceTextMarshaler, false, nil)
}

// MaybeTextMarshaler queries whether and how the given type
// implements the [encoding.TextMarshaler] interface for at least one receiver form.
func MaybeTextMarshaler(typ types.Type) MarshalerKind {
	return isMarshaler(typ, ifaceTextMarshaler, true, nil)
}

// IsJSONMarshaler queries whether and how the given type
// implements the [json.Marshaler] interface.
func IsJSONMarshaler(typ types.Type) MarshalerKind {
	return isMarshaler(typ, ifaceJSONMarshaler, false, nil)
}

// MaybeJSONMarshaler queries whether and how the given type
// implements the [json.Marshaler] interface for at least one receiver form.
func MaybeJSONMarshaler(typ types.Type) MarshalerKind {
	return isMarshaler(typ, ifaceJSONMarshaler, true, nil)
}

// IsMapKey returns true if the given type
// is accepted as a map key by encoding/json.
func IsMapKey(typ types.Type) bool {
	// Iterate over type set and return true if all elements are valid.
	//
	// We cannot simply delegate to [IsTextMarshaler] here
	// because a union of some basic terms and some TextMarshalers
	// might still be acceptable.
	//
	// NOTE: If typ is not a typeparam or constraint, termlist returns just typ itself.
	// If typ has an empty type set, it's safe to return true
	// because the map cannot be instantiated anyways.
	for term := range termlist(typ) {
		ttyp := types.Unalias(term.Type())

		// Types whose underlying type is a signed/unsigned integer or a string
		// are always acceptable, whether they are marshalers or not.
		if basic, ok := ttyp.Underlying().(*types.Basic); ok {
			if basic.Info()&(types.IsInteger|types.IsUnsigned|types.IsString) != 0 {
				continue
			}
		}

		// Valid tilde terms
		// can always be satisfied by named types that hide their methods
		// hence fail in general to implement the required interface.
		// For example one could have:
		//
		//     type NotAKey struct{ encoding.TextMarshaler }
		//     func (NotAKey) MarshalText() int { ... }
		//
		// which satisfies ~struct{ encoding.TextMarshaler }
		// but is not itself a TextMarshaler.
		//
		// It might still be the case that the constraint
		// requires explicitly a marshaling method,
		// hence we perform one last check on typ.
		//
		// For example, we reject interface{ ~struct{ ... } }
		// but still accept interface{ ~struct{ ... }; MarshalText() ([]byte, error) }
		//
		// All other cases are only acceptable
		// if the type implements [encoding.TextMarshaler] in non-addressable mode.
		if term.Tilde() || IsTextMarshaler(ttyp) == NonMarshaler {
			// When some term fails, test the input typ itself,
			// but only if it has not been tested already.
			//
			// Note that when term.Tilde() is true
			// then it is always the case that typ != term.Type(),
			// because cyclic constraints are not allowed
			// and naked type parameters cannot occur in type unions.
			return typ != term.Type() && IsTextMarshaler(typ) != NonMarshaler
		}
	}

	return true
}

// IsTypeParam returns true when the given type
// is either a TypeParam or a pointer to a TypeParam.
func IsTypeParam(typ types.Type) bool {
	switch t := types.Unalias(typ).(type) {
	case *types.TypeParam:
		return true
	case *types.Pointer:
		_, ok := types.Unalias(t.Elem()).(*types.TypeParam)
		return ok
	default:
		return false
	}
}

// IsStringAlias returns true when
// either typ will be rendered to JS/TS as an alias for the TS type `string`,
// or typ itself (not its underlying type) is a pointer
// whose element type satisfies the property described above.
//
// This predicate is only safe to use either with map keys,
// where pointers are treated in an ad-hoc way by [json.Marshal],
// or when typ IS ALREADY KNOWN to be either [types.Alias] or [types.Named].
//
// Otherwise, the result might be incorrect:
// IsStringAlias MUST NOT be used to check
// whether an arbitrary instance of [types.Type]
// renders as a JS/TS string type.
//
// Notice that IsStringAlias returns false for all type parameters:
// detecting those that must be always instantiated as string aliases
// is technically possible, but very difficult.
func IsStringAlias(typ types.Type) bool {
	// Unwrap at most one pointer.
	// NOTE: do not unalias typ before testing:
	// aliases whose underlying type is a pointer
	// are never rendered as strings.
	if ptr, ok := typ.(*types.Pointer); ok {
		typ = ptr.Elem()
	}

	switch typ.(type) {
	case *types.Alias, *types.Named:
		// Aliases and named types might be rendered as string aliases.
	default:
		// Not a model type, hence not an alias.
		return false
	}

	// Skip pointer and interface types: they are always nullable
	// and cannot have any explicitly defined methods.
	// This takes care of rejecting type params as well,
	// since their underlying type is guaranteed to be an interface.
	switch typ.Underlying().(type) {
	case *types.Pointer, *types.Interface:
		return false
	}

	// Follow alias chain.
	typ = types.Unalias(typ)

	// Aliases of the basic string type are rendered as strings.
	if basic, ok := typ.(*types.Basic); ok {
		return basic.Info()&types.IsString != 0
	}

	// json.Marshalers can only be rendered as any.
	// TextMarshalers that aren't json.Marshalers render as strings.
	if MaybeJSONMarshaler(typ) != NonMarshaler {
		return false
	} else if MaybeTextMarshaler(typ) != NonMarshaler {
		return true
	}

	// Named types whose underlying type is a string are rendered as strings.
	basic, ok := typ.Underlying().(*types.Basic)
	return ok && basic.Info()&types.IsString != 0
}

// IsClass returns true if the given type will be rendered
// as a JS/TS model class (or interface).
func IsClass(typ types.Type) bool {
	// Follow alias chain.
	typ = types.Unalias(typ)

	if _, isNamed := typ.(*types.Named); !isNamed {
		// Unnamed types are never rendered as classes.
		return false
	}

	// Struct named types without custom marshaling are rendered as classes.
	_, isStruct := typ.Underlying().(*types.Struct)
	return isStruct && MaybeJSONMarshaler(typ) == NonMarshaler && MaybeTextMarshaler(typ) == NonMarshaler
}

// IsAny returns true if the given type
// is guaranteed to render as the TS any type or equivalent.
//
// It might return false negatives for generic aliases,
// hence should only be used with instantiated types
// or in contexts where false negatives are acceptable.
func IsAny(typ types.Type) bool {
	// Follow alias chain.
	typ = types.Unalias(typ)

	if MaybeJSONMarshaler(typ) != NonMarshaler {
		// If typ is either a named type, an interface, a pointer or a struct,
		// it will be rendered as (possibly an alias for) the TS any type.
		//
		// If it is a type parameter that implements json.Marshal,
		// every possible concrete instantiation will implement json.Marshal,
		// hence will be rendered as the TS any type.
		return true
	}

	if MaybeTextMarshaler(typ) != NonMarshaler {
		// If type is either a named type, an interface, a pointer or a struct,
		// it will be rendered as (possibly an alias for)
		// the (possibly nullable) TS string type.
		//
		// If typ is a type parameter, we know at this point
		// that it does not necessarily implement json.Marshaler,
		// hence it will be possible to instantiate it in a way
		// that renders as the (possibly nullable) TS string type.
		return false
	}

	if ptr, ok := typ.Underlying().(*types.Pointer); ok {
		// Pointers render as the union of their element type with null.
		// This is equivalent to the TS any type
		// if and only if so is the element type.
		return IsAny(ptr.Elem())
	}

	// All types listed below have rich TS equivalents,
	// hence won't be equivalent to the TS any type.
	//
	// WARN: it is important to keep these lists explicit and up to date
	// instead of listing the unsupported types (which would be much easier).
	//
	// By doing so, IsAny will keep working correctly
	// in case future updates to the Go spec introduce new type families,
	// thus buying the maintainers some time to patch the binding generator.

	// Retrieve underlying type.
	switch t := typ.Underlying().(type) {
	case *types.Basic:
		// Complex types are not supported.
		return t.Info()&(types.IsBoolean|types.IsInteger|types.IsUnsigned|types.IsFloat|types.IsString) == 0
	case *types.Array, *types.Slice, *types.Map, *types.Struct, *types.TypeParam:
		return false
	}

	return true
}
