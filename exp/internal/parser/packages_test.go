package parser

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestParseDirectory(t *testing.T) {
	tests := []struct {
		name    string
		dir     string
		want    []string
		wantErr bool
	}{
		{
			name:    "should find single bound service",
			dir:     "testdata/struct_literal_single",
			want:    []string{"GreetService"},
			wantErr: false,
		},
		{
			name:    "should find multiple bound services",
			dir:     "testdata/struct_literal_multiple",
			want:    []string{"main.GreetService", "main.OtherService"},
			wantErr: false,
		},
		{
			name:    "should find multiple bound services over multiple files",
			dir:     "testdata/struct_literal_multiple_files",
			want:    []string{"main.GreetService", "main.OtherService"},
			wantErr: false,
		},
		{
			name:    "should find bound services from other packages",
			dir:     "../../examples/binding",
			want:    []string{"main.localStruct", "services.GreetService", "models.Person"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDirectory(tt.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for name, pkg := range got.packages {
				println("Got package", name)
				for _, boundStruct := range pkg.boundStructs {
					println("Got bound struct", boundStruct.Name.Name)
					require.True(t, lo.Contains(tt.want, name+"."+boundStruct.Name.Name))
				}
			}
		})
	}
}
