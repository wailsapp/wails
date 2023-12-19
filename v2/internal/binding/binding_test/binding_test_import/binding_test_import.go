package binding_test_import

import "github.com/wailsapp/wails/v2/internal/binding/binding_test/binding_test_import/binding_test_nestedimport"

type AWrapper struct {
	AWrapper binding_test_nestedimport.A `json:"AWrapper"`
}

type ASliceWrapper struct {
	ASlice []binding_test_nestedimport.A `json:"ASlice"`
}

type AMapWrapper struct {
	AMap map[string]binding_test_nestedimport.A `json:"AMap"`
}

type ImportedEnum string

const (
	ImportedEnumValue1 ImportedEnum = "value1"
	ImportedEnumValue2 ImportedEnum = "value2"
	ImportedEnumValue3 ImportedEnum = "value3"
)

var AllImportedEnumValues = []struct {
	Value  ImportedEnum
	TSName string
}{
	{ImportedEnumValue1, "Value1"},
	{ImportedEnumValue2, "Value2"},
	{ImportedEnumValue3, "Value3"},
}
