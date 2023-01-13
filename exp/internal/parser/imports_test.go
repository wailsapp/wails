package parser

import (
	"testing"

	"github.com/samber/lo"
)

func Test_findFilesImportingPackage(t *testing.T) {
	tests := []struct {
		name    string
		dir     string
		pkg     string
		want    []string
		wantErr bool
	}{
		{
			name: "should identify single file importing package",
			dir:  "testdata/imports/single_file",
			pkg:  "github.com/wailsapp/wails/exp/pkg/application",
			want: []string{"testdata/imports/single_file/main.go"},
		},
		{
			name: "should identify multiple files importing package",
			dir:  "testdata/imports/multiple_files",
			pkg:  "github.com/wailsapp/wails/exp/pkg/application",
			want: []string{"testdata/imports/multiple_files/app.go"},
		},
		{
			name: "should identify aliases",
			dir:  "testdata/imports/alias",
			pkg:  "github.com/wailsapp/wails/exp/pkg/application",
			want: []string{"testdata/imports/alias/main.go"},
		},
		{
			name: "should identify packages",
			dir:  "testdata/imports/other_package",
			pkg:  "github.com/wailsapp/wails/exp/pkg/application",
			want: []string{"testdata/imports/other_package/subpackage/app.go"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findFilesImportingPackage(tt.dir, tt.pkg)
			if (err != nil) != tt.wantErr {
				t.Errorf("findFilesImportingPackage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, imp := range got {
				if !lo.Contains(tt.want, imp.FileName) {
					t.Errorf("findFilesImportingPackage() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
