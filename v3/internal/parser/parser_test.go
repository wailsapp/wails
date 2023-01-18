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
			want:    []string{"main.GreetService"},
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
			Debug = true
			got, err := ParseDirectory(tt.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for name, pkg := range got.packages {
				for structName := range pkg.boundStructs {
					require.True(t, lo.Contains(tt.want, name+"."+structName))
					tt.want = lo.Without(tt.want, name+"."+structName)
				}
			}
			require.Empty(t, tt.want)
		})
	}
}

func TestGenerateTypeScript(t *testing.T) {
	tests := []struct {
		name    string
		dir     string
		want    string
		wantErr bool
	}{
		{
			name: "should find single bound service",
			dir:  "testdata/struct_literal_single",
			want: `namespace main {
  class GreetService {
    SomeVariable: number;
  }
}
`,
			wantErr: false,
		},
		{
			name: "should find multiple bound services",
			dir:  "testdata/struct_literal_multiple",
			want: `namespace main {
  class GreetService {
    SomeVariable: number;
  }
  class OtherService {
  }
}
`,
			wantErr: false,
		},
		{
			name: "should find multiple bound services over multiple files",
			dir:  "testdata/struct_literal_multiple_files",
			want: `namespace main {
  class GreetService {
    SomeVariable: number;
  }
  class OtherService {
  }
}
`,
			wantErr: false,
		},
		{
			name: "should find bound services from other packages",
			dir:  "../../examples/binding",
			want: `namespace main {
  class localStruct {
  }
}
namespace models {
  class Person {
    Name: string;
  }
}
namespace services {
  class GreetService {
    SomeVariable: number;
    Parent: models.Person;
  }
}
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Debug = true
			context, err := ParseDirectory(tt.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			ts, err := GenerateModels(context)
			require.NoError(t, err)
			require.Equal(t, tt.want, string(ts))

		})
	}
}
