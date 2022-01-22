package binding

import (
	"testing"
)

func Test_goTypeToJSDocType(t *testing.T) {

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "string",
			input: "string",
			want:  "string",
		},
		{
			name:  "error",
			input: "error",
			want:  "Error",
		},
		{
			name:  "int",
			input: "int",
			want:  "number",
		},
		{
			name:  "int32",
			input: "int32",
			want:  "number",
		},
		{
			name:  "uint",
			input: "uint",
			want:  "number",
		},
		{
			name:  "uint32",
			input: "uint32",
			want:  "number",
		},
		{
			name:  "float32",
			input: "float32",
			want:  "number",
		},
		{
			name:  "float64",
			input: "float64",
			want:  "number",
		},
		{
			name:  "bool",
			input: "bool",
			want:  "boolean",
		},
		{
			name:  "interface{}",
			input: "interface{}",
			want:  "any",
		},
		{
			name:  "[]byte",
			input: "[]byte",
			want:  "string",
		},
		{
			name:  "[]int",
			input: "[]int",
			want:  "Array<number>",
		},
		{
			name:  "[]bool",
			input: "[]bool",
			want:  "Array<boolean>",
		},
		{
			name:  "anything else",
			input: "foo",
			want:  "any",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := goTypeToJSDocType(tt.input); got != tt.want {
				t.Errorf("goTypeToJSDocType() = %v, want %v", got, tt.want)
			}
		})
	}
}
