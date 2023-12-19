package binding

import (
	"testing"

	"github.com/leaanthony/slicer"
	"github.com/stretchr/testify/assert"
	"github.com/wailsapp/wails/v2/internal/logger"
)

type BindForTest struct {
}

func (b *BindForTest) GetA() A {
	return A{}
}

type A struct {
	B B `json:"B"`
}

type B struct {
	Name string `json:"name"`
}

func TestNestedStruct(t *testing.T) {
	bind := &BindForTest{}
	testBindings := NewBindings(logger.New(nil), []interface{}{bind}, []interface{}{}, false, []interface{}{})

	namesStrSlicer := testBindings.getAllStructNames()
	names := []string{}
	namesStrSlicer.Each(func(s string) {
		names = append(names, s)
	})
	assert.Contains(t, names, "binding.A")
	assert.Contains(t, names, "binding.B")
}

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
		{
			name:  "map",
			input: "map[string]float64",
			want:  "{[key: string]: number}",
		},
		{
			name:  "map",
			input: "map[string]map[string]float64",
			want:  "{[key: string]: {[key: string]: number}}",
		},
		{
			name:  "types",
			input: "main.SomeType",
			want:  "main.SomeType",
		},
	}
	var importNamespaces slicer.StringSlicer
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := goTypeToJSDocType(tt.input, &importNamespaces); got != tt.want {
				t.Errorf("goTypeToJSDocType() = %v, want %v", got, tt.want)
			}
		})
	}
}
