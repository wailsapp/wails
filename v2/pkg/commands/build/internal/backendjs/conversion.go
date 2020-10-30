package backendjs

import "reflect"

// JSType represents a javascript type
type JSType string

const (
	// JsString is a JS string
	JsString JSType = "string"
	// JsBoolean is a JS bool
	JsBoolean = "boolean"
	// JsInt is a JS number
	JsInt = "number"
	// JsFloat is a JS number
	JsFloat = "number"
	// JsArray is a JS array
	JsArray = "Array"
	// JsObject is a JS object
	JsObject = "Object"
	// JsUnsupported represents a type that cannot be converted
	JsUnsupported = "*"
)

func goTypeToJS(input reflect.Kind) JSType {
	switch input {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return JsInt
	case reflect.String:
		return JsString
	case reflect.Float32, reflect.Float64, reflect.Complex64:
		return JsFloat
	case reflect.Bool:
		return JsBoolean
	case reflect.Array, reflect.Slice:
		return JsArray
	case reflect.Ptr, reflect.Struct, reflect.Map, reflect.Interface:
		return JsObject
	default:
		return JsUnsupported
	}
}
