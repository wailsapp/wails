package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_findNewCalls(t *testing.T) {
	tests := []struct {
		name                string
		dir                 string
		expectedExpressions int
	}{
		{
			name:                "should find single call to application.New",
			dir:                 "testdata/boundstructs/struct_literal_single",
			expectedExpressions: 1,
		},
		{
			name:                "should find single call to application.New in multiple files",
			dir:                 "testdata/boundstructs/struct_literal_multiple_files",
			expectedExpressions: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imports, err := findFilesImportingPackage(tt.dir, "github.com/wailsapp/wails/exp/pkg/application")
			require.NoError(t, err)
			got := findNewCalls(imports)
			require.NotNil(t, got.Options)
		})
	}
}
