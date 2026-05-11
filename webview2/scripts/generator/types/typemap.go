package types

// IDLTypePattern documents one IDL → Go type mapping pattern.
// The full mapping depends on three axes: IDL base type, pointer depth, and
// parameter direction ([in] vs [out]/[out,retval]).
type IDLTypePattern struct {
	// IDL shows the declaration as it appears in the .idl file.
	IDL string
	// GoType is the resulting Go type seen in the method signature.
	GoType string
	// VtableArg is the expression passed to ComProc.Call for this parameter.
	VtableArg string
	// Notes explains any non-obvious conversion.
	Notes string
}

// TypePatterns lists all 17 IDL→Go patterns handled by the code generator.
// The table is the authoritative reference for param.go + template behaviour.
//
// Pointer depth key:
//   (none) = no '*' on the IDL type
//   *      = one '*' on the IDL type
//   **     = two '*'s on the IDL type
var TypePatterns = []IDLTypePattern{
	// ── Strings ─────────────────────────────────────────────────────────────
	{
		IDL:       "[in] LPWSTR param",
		GoType:    "string",
		VtableArg: "uintptr(unsafe.Pointer(_param))",
		Notes:     "inputStringSetup.tmpl converts Go string → *uint16 (_param).",
	},
	{
		IDL:       "[in] LPCWSTR param",
		GoType:    "string",
		VtableArg: "uintptr(unsafe.Pointer(_param))",
		Notes:     "Same as LPWSTR input; const qualifier is irrelevant in Go.",
	},
	{
		IDL:       "[out, retval] LPWSTR* param",
		GoType:    "string",
		VtableArg: "uintptr(unsafe.Pointer(&_param))",
		Notes:     "outputStringSetup.tmpl declares var _param *uint16; address is passed so COM writes the pointer back. outputStringCleanup.tmpl converts to string and calls CoTaskMemFree.",
	},
	{
		IDL:       "[in] LPCWSTR* param",
		GoType:    "[]string",
		VtableArg: "uintptr(unsafe.Pointer(_param))",
		Notes:     "inputStringArraySetup.tmpl converts []string → []*uint16 slice then exposes **uint16 (_param) pointing at slice[0].",
	},
	{
		IDL:       "[out] LPWSTR** param",
		GoType:    "*string",
		VtableArg: "uintptr(unsafe.Pointer(&param))",
		Notes:     "Double-pointer COM string array output; caller must CoTaskMemFree each element and the array itself.",
	},

	// ── Unsigned integers ────────────────────────────────────────────────────
	{
		IDL:       "[in] UINT32 param",
		GoType:    "uint32",
		VtableArg: "uintptr(param)",
		Notes:     "Scalar passed by value; no setup code needed.",
	},
	{
		IDL:       "[out] UINT32* param",
		GoType:    "uint32",
		VtableArg: "uintptr(unsafe.Pointer(&param))",
		Notes:     "outputDefaultSetup.tmpl declares var param uint32.",
	},
	{
		IDL:       "[in] UINT64 param",
		GoType:    "uint64",
		VtableArg: "uintptr(param)",
		Notes:     "Scalar passed by value.",
	},

	// ── Signed integers ──────────────────────────────────────────────────────
	{
		IDL:       "[in] INT32 param",
		GoType:    "int32",
		VtableArg: "uintptr(param)",
		Notes:     "Scalar passed by value.",
	},
	{
		IDL:       "[out, retval] INT32* param",
		GoType:    "int32",
		VtableArg: "uintptr(unsafe.Pointer(&param))",
		Notes:     "outputDefaultSetup.tmpl declares var param int32.",
	},
	{
		IDL:       "[out, retval] int* param",
		GoType:    "int",
		VtableArg: "uintptr(unsafe.Pointer(&param))",
		Notes:     "Native Go int (platform-sized). outputDefaultSetup.tmpl declares var param int.",
	},

	// ── Booleans ─────────────────────────────────────────────────────────────
	{
		IDL:       "[in] BOOL param",
		GoType:    "bool",
		VtableArg: "uintptr(_param)",
		Notes:     "COM BOOL is int32 (4 bytes). inputBoolSetup.tmpl declares var _param int32 and sets it to 1 if param is true.",
	},
	{
		IDL:       "[out, retval] BOOL* param",
		GoType:    "bool",
		VtableArg: "uintptr(unsafe.Pointer(&_param))",
		Notes:     "outputBoolSetup.tmpl declares var _param int32. outputBoolCleanup.tmpl assigns param := _param != 0.",
	},

	// ── Floating point ───────────────────────────────────────────────────────
	{
		IDL:       "[in] double param",
		GoType:    "float64",
		VtableArg: "uintptr(param)",
		Notes:     "Scalar passed by value.",
	},

	// ── COM interface pointers ───────────────────────────────────────────────
	{
		IDL:       "[in] IInterface* param",
		GoType:    "*IInterface",
		VtableArg: "uintptr(unsafe.Pointer(param))",
		Notes:     "Pointer value passed directly; no address-of.",
	},
	{
		IDL:       "[out] IInterface** param",
		GoType:    "*IInterface",
		VtableArg: "uintptr(unsafe.Pointer(&param))",
		Notes:     "outputDefaultSetup.tmpl declares var param *IInterface. COM writes the pointer into &param.",
	},

	// ── Struct / token types ─────────────────────────────────────────────────
	{
		IDL:       "[out] EventRegistrationToken* param",
		GoType:    "EventRegistrationToken",
		VtableArg: "uintptr(unsafe.Pointer(&param))",
		Notes:     "outputDefaultSetup.tmpl declares var param EventRegistrationToken.",
	},
}

// ResolveGoType returns the Go type string for an IDL type given pointer depth
// and direction.  It is the single-source-of-truth called from Param.Process().
// direction should be "in" or "out".
func ResolveGoType(idlType, pointer, direction string) string {
	isOut := direction == "out"

	switch idlType {
	case "LPWSTR", "LPCWSTR":
		// [in] LPCWSTR* = array of strings; all other string variants = string.
		if !isOut && pointer == "*" {
			return "[]string"
		}
		return "string"

	case "HRESULT":
		return "uintptr"
	case "UINT64":
		return "uint64"
	case "UINT32", "DWORD":
		return "uint32"
	case "UINT":
		return "uint"
	case "INT64":
		return "int64"
	case "INT32":
		return "int32"
	case "INT":
		return "int"
	case "BOOL":
		return "bool"
	case "BYTE":
		return "uint8"
	case "double":
		return "float64"
	case "IUnknown":
		return "IUnknown"
	case "EventRegistrationToken":
		return "EventRegistrationToken"
	}

	// Unknown IDL type: pass through unchanged (interface name, enum, struct …).
	return idlType
}
