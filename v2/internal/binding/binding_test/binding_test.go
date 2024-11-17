package binding_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/logger"
)

type BindingTest struct {
	name        string
	structs     []interface{}
	enums       []interface{}
	exemptions  []interface{}
	want        string
	shouldError bool
	TsGenerationOptionsTest
}

type TsGenerationOptionsTest struct {
	TsPrefix     string
	TsSuffix     string
	TsOutputType string
}

func TestBindings_GenerateModels(t *testing.T) {

	tests := []BindingTest{
		EscapedNameTest,
		ImportedStructTest,
		ImportedSliceTest,
		ImportedMapTest,
		ImportedEnumTest,
		NestedFieldTest,
		NonStringMapKeyTest,
		SingleFieldTest,
		MultistructTest,
		EmptyStructTest,
		GeneratedJsEntityTest,
		GeneratedJsEntityWithIntEnumTest,
		GeneratedJsEntityWithStringEnumTest,
		GeneratedJsEntityWithEnumTsName,
		GeneratedJsEntityWithNestedStructInterfacesTest,
		AnonymousSubStructTest,
		AnonymousSubStructMultiLevelTest,
		GeneratedJsEntityWithNestedStructTest,
		EntityWithDiffNamespacesTest,
		SpecialCharacterFieldTest,
		WithoutFieldsTest,
		NoFieldTagsTest,
	}

	testLogger := &logger.Logger{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := binding.NewBindings(testLogger, tt.structs, tt.exemptions, false, tt.enums)
			for _, s := range tt.structs {
				err := b.Add(s)
				require.NoError(t, err)
			}
			b.SetTsPrefix(tt.TsPrefix)
			b.SetTsSuffix(tt.TsSuffix)
			b.SetOutputType(tt.TsOutputType)
			got, err := b.GenerateModels()
			if (err != nil) != tt.shouldError {
				t.Errorf("GenerateModels() error = %v, shouldError %v", err, tt.shouldError)
				return
			}
			if !reflect.DeepEqual(strings.Fields(string(got)), strings.Fields(tt.want)) {
				t.Errorf("GenerateModels() got = %v, want %v", string(got), tt.want)
			}
		})
	}
}
