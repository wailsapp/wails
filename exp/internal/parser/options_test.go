package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_isOptionsApplication(t *testing.T) {
	tests := []struct {
		name string
		dir  string
		want bool
	}{
		{
			name: "should return true when expr is a struct literal",
			dir:  "testdata/boundstructs/struct_literal_single",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imports, err := findFilesImportingPackage(tt.dir, "github.com/wailsapp/wails/exp/pkg/application")
			require.NoError(t, err)
			imp := findNewCalls(imports)
			require.NotNil(t, imp)
			isOptions := isOptionsApplication(imp.Options)
			require.Equal(t, tt.want, isOptions)
		})
	}
}
