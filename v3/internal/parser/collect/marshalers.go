package collect

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
// the encoding.TextMarshaler interface with non-pointer receiver.
func IsAlwaysTextMarshaler(typ types.Type) bool {
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
// the json.Marshaler interface with pointer receiver.
func MaybeJSONMarshaler(typ types.Type) bool {
	if _, ok := types.Unalias(typ).(*types.Pointer); !ok {
		typ = types.NewPointer(typ)
	}
	return types.Implements(typ, ifaceJSONMarshaler)
}

// IsMapKey returns true if the given type
// is acceptable as a map key by encoding/json.
func IsMapKey(typ types.Type) bool {
	if basic, ok := types.Unalias(typ).(*types.Basic); ok {
		return basic.Info()&(types.IsInteger|types.IsString) != 0
	}

	return IsTextMarshaler(typ)
}
