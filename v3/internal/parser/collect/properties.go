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

// AlwaysTextMarshaler tests whether the given type implements
// the encoding.TextMarshaler interface for all receiver forms.
func AlwaysTextMarshaler(typ types.Type) bool {
	if ptr, ok := types.Unalias(typ).(*types.Pointer); ok {
		typ = ptr.Elem()
	}
	return types.Implements(typ, ifaceTextMarshaler)
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
// is acceptable as a map key by encoding/json.
func IsMapKey(typ types.Type) bool {
	if basic, ok := types.Unalias(typ).(*types.Basic); ok {
		return basic.Info()&(types.IsInteger|types.IsString) != 0
	}

	// Other types are only accepted
	// if they implement encoding.TextMarshaler strictly as they are.
	return IsTextMarshaler(typ)
}

// IsClass returns true if the given type (or element type for pointers)
// will be rendered as a JS/TS model class (or interface).
func IsClass(typ types.Type) bool {
	if ptr, _ := typ.(*types.Pointer); ptr != nil {
		// Unwrap at most one pointer.
		// NOTE: do not unalias typ before testing:
		// aliases whose underlying type is a pointer
		// are not rendered as classes.
		typ = ptr.Elem()
	}

	// Follow alias chain.
	typ = types.Unalias(typ)

	if _, isNamed := typ.(*types.Named); !isNamed {
		// Not a model type.
		return false
	}

	// JSONMarshalers can only be rendered as any.
	// TextMarshalers are rendered as strings if they are not JSONMarshalers.
	if MaybeJSONMarshaler(typ) || MaybeTextMarshaler(typ) {
		return false
	}

	// Struct types without custom marshaling are rendered as classes.
	_, isStruct := typ.Underlying().(*types.Struct)
	return isStruct
}

// IsString returns true if the given type (or element type for pointers)
// will be rendered as an alias for the JS string type.
func IsString(typ types.Type) bool {
	if ptr, ok := typ.(*types.Pointer); ok {
		// Unwrap at most one pointer.
		// NOTE: do not unalias typ before testing:
		// aliases whose underlying type is a pointer
		// are not rendered as strings.
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
	// TextMarshalers are rendered as strings if they are not JSONMarshalers.
	return !MaybeJSONMarshaler(typ) && MaybeTextMarshaler(typ)
}
