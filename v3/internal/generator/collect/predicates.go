package collect

// This file gathers functions that test useful properties of model types.
// The rationale for the way things are handled here
// is given in the example file found at ./_reference/json_marshaler_behaviour.go

import (
	"go/token"
	"go/types"
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

// IsTextMarshaler tests whether the given type implements
// the encoding.TextMarshaler interface.
func IsTextMarshaler(typ types.Type) bool {
	return types.Implements(typ, ifaceTextMarshaler)
}

// MaybeTextMarshaler tests whether the given type implements
// the encoding.TextMarshaler interface for at least one receiver form.
func MaybeTextMarshaler(typ types.Type) bool {
	if _, ok := types.Unalias(typ).(*types.Pointer); !ok {
		typ = types.NewPointer(typ)
	}
	return IsTextMarshaler(typ)
}

// IsJSONMarshaler tests whether the given type implements
// the json.Marshaler interface.
func IsJSONMarshaler(typ types.Type) bool {
	return types.Implements(typ, ifaceJSONMarshaler)
}

// MaybeJSONMarshaler tests whether the given type implements
// the json.Marshaler interface for at least one receiver form.
func MaybeJSONMarshaler(typ types.Type) bool {
	if _, ok := types.Unalias(typ).(*types.Pointer); !ok {
		typ = types.NewPointer(typ)
	}
	return IsJSONMarshaler(typ)
}

// IsMapKey returns true if the given type
// is accepted as a map key by encoding/json.
func IsMapKey(typ types.Type) bool {
	if basic, ok := typ.Underlying().(*types.Basic); ok {
		return basic.Info()&(types.IsInteger|types.IsString) != 0
	}

	// Other types are only accepted
	// if they implement encoding.TextMarshaler strictly as they are.
	return IsTextMarshaler(typ)
}

// IsString returns true if the given type (or element type for pointers)
// will be rendered as an alias for the TS string type.
func IsString(typ types.Type) bool {
	// Unwrap at most one pointer.
	// NOTE: do not unalias typ before testing:
	// aliases whose underlying type is a pointer
	// are _never_ rendered as strings.
	if ptr, ok := typ.(*types.Pointer); ok {
		typ = ptr.Elem()
	}

	switch typ.(type) {
	case *types.Alias, *types.Named:
		// Aliases and named types might be rendered as string aliases.
	default:
		// Not a model type.
		return false
	}

	// Follow alias chain.
	typ = types.Unalias(typ)

	if basic, ok := typ.(*types.Basic); ok {
		// Test whether basic type is a string.
		return basic.Info()&types.IsString != 0
	}

	// JSONMarshalers can only be rendered as any.
	// TextMarshalers that aren't JSONMarshalers render as strings.
	if MaybeJSONMarshaler(typ) {
		return false
	} else if MaybeTextMarshaler(typ) {
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
	return isStruct && !MaybeJSONMarshaler(typ) && !MaybeTextMarshaler(typ)
}

// IsAny returns true if the given type will be rendered as a TS any type.
func IsAny(typ types.Type) bool {
	// Follow alias chain.
	typ = types.Unalias(typ)

	if MaybeJSONMarshaler(typ) {
		return true
	}

	if MaybeTextMarshaler(typ) {
		return false
	}

	// Retrieve underlying type
	switch t := typ.Underlying().(type) {
	case *types.Basic:
		// Complex types are not supported.
		return t.Info()&types.IsComplex != 0
	case *types.Chan, *types.Signature, *types.Interface:
		return true
	}

	return false
}
