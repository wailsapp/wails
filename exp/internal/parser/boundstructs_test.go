package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_extractBindExprs(t *testing.T) {
	tests := []struct {
		name                 string
		dir                  string
		expectedBoundStructs []*structInfo
	}{
		{
			name:                 "should find single bound struct literal",
			dir:                  "testdata/boundstructs/struct_literal_single",
			expectedBoundStructs: []*structInfo{{structName: "GreetService"}},
		},
		{
			name:                 "should find multiple bound struct literal",
			dir:                  "testdata/boundstructs/struct_literal_multiple",
			expectedBoundStructs: []*structInfo{{structName: "GreetService"}, {structName: "OtherService"}},
		},
		{
			name:                 "should find multiple bound struct literals over multiple files",
			dir:                  "testdata/boundstructs/struct_literal_multiple_files",
			expectedBoundStructs: []*structInfo{{structName: "GreetService"}, {structName: "OtherService"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imports, err := findFilesImportingPackage(tt.dir, "github.com/wailsapp/wails/exp/pkg/application")
			require.NoError(t, err)
			imp := findNewCalls(imports)
			require.NotNil(t, imp)
			isOptions := isOptionsApplication(imp.Options)
			require.True(t, isOptions)
			extractBindExprs(imp)
			require.ElementsMatchf(t, tt.expectedBoundStructs, imp.BoundStructNames, "expected %v, got %v", tt.expectedBoundStructs, imp.BoundStructNames)

		})
	}
}
