package types

var idlTypeToGoType = map[string]string{
	"IUnknown":               "IUnknown",
	"EventRegistrationToken": "EventRegistrationToken",
	"LPWSTR":                 "string",
	"LPCWSTR":                "string",
	"HRESULT":                "uintptr",
	"UINT64":                 "uint64",
	"UINT32":                 "uint32",
	"UINT":                   "uint",
	"INT":                    "int",
	"INT32":                  "int32",
	"INT64":                  "int64",
	"BOOL":                   "bool",
	"BYTE":                   "uint8",
	"DWORD":                  "uint32",
	"double":                 "float64",
}

func IdlTypeToGoType(input string) string {
	result := idlTypeToGoType[input]
	if result == "" {
		return input
	}
	return result
}
