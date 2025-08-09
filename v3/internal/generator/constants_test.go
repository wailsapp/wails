package generator

import "testing"

func TestGenerateConstants(t *testing.T) {
	tests := []struct {
		name    string
		goData  []byte
		want    string
		wantErr bool
	}{
		{
			name: "int",
			goData: []byte(`package test
const one = 1`),
			want:    "export const one = 1;",
			wantErr: false,
		},
		{
			name: "float",
			goData: []byte(`package test
const one_point_five = 1.5`),
			want:    "export const one_point_five = 1.5;",
			wantErr: false,
		},
		{
			name: "string",
			goData: []byte(`package test
const one_as_a_string = "1"`),
			want:    `export const one_as_a_string = "1";`,
			wantErr: false,
		},
		{
			name: "nested",
			goData: []byte(`package test
const (
  one_as_a_string = "1"
)`),
			want:    `export const one_as_a_string = "1";`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateConstants(tt.goData)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateConstants() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateConstants() got = %v, want %v", got, tt.want)
			}
		})
	}
}
