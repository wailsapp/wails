package application_test

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type TestService struct {
}

type Person struct {
	Name string `json:"name"`
}

func (t *TestService) Nil() {}

func (t *TestService) String(s string) string {
	return s
}

func (t *TestService) Multiple(s string, i int, b bool) (string, int, bool) {
	return s, i, b
}

func (t *TestService) Struct(p Person) Person {
	return p
}

func (t *TestService) StructNil(p Person) (Person, error) {
	return p, nil
}

func (t *TestService) StructError(p Person) (Person, error) {
	return p, errors.New("error")
}

func (t *TestService) Variadic(s ...string) []string {
	return s
}

func (t *TestService) PositionalAndVariadic(a int, b ...string) int {
	return a
}

func (t *TestService) Slice(a []int) []int {
	return a
}

func newArgs(jsonArgs ...string) []json.RawMessage {
	args := []json.RawMessage{}

	for _, j := range jsonArgs {
		args = append(args, json.RawMessage(j))
	}
	return args
}

func TestBoundMethodCall(t *testing.T) {

	tests := []struct {
		name     string
		method   string
		args     []json.RawMessage
		err      error
		expected interface{}
	}{
		{
			name:     "nil",
			method:   "Nil",
			args:     []json.RawMessage{},
			err:      nil,
			expected: nil,
		},
		{
			name:     "string",
			method:   "String",
			args:     newArgs(`"foo"`),
			err:      nil,
			expected: "foo",
		},
		{
			name:     "multiple",
			method:   "Multiple",
			args:     newArgs(`"foo"`, "0", "false"),
			err:      nil,
			expected: []interface{}{"foo", 0, false},
		},
		{
			name:     "struct",
			method:   "Struct",
			args:     newArgs(`{ "name": "alice" }`),
			err:      nil,
			expected: Person{Name: "alice"},
		},
		{
			name:     "struct, nil error",
			method:   "StructNil",
			args:     newArgs(`{ "name": "alice" }`),
			err:      nil,
			expected: Person{Name: "alice"},
		},
		{
			name:     "struct, error",
			method:   "StructError",
			args:     newArgs(`{ "name": "alice" }`),
			err:      errors.New("error"),
			expected: nil,
		},
		{
			name:     "invalid argument count",
			method:   "Multiple",
			args:     newArgs(`"foo"`),
			err:      errors.New("expects 3 arguments, received 1"),
			expected: nil,
		},
		{
			name:     "invalid argument type",
			method:   "String",
			args:     newArgs("1"),
			err:      errors.New("could not parse"),
			expected: nil,
		},
		{
			name:     "variadic, no arguments",
			method:   "Variadic",
			args:     newArgs(`[]`), // variadic parameters are passed as arrays
			err:      nil,
			expected: []string{},
		},
		{
			name:     "variadic",
			method:   "Variadic",
			args:     newArgs(`["foo", "bar"]`),
			err:      nil,
			expected: []string{"foo", "bar"},
		},
		{
			name:     "positional and variadic",
			method:   "PositionalAndVariadic",
			args:     newArgs("42", `[]`),
			err:      nil,
			expected: 42,
		},
		{
			name:     "slice",
			method:   "Slice",
			args:     newArgs(`[1,2,3]`),
			err:      nil,
			expected: []int{1, 2, 3},
		},
	}

	// init globalApplication
	_ = application.New(application.Options{})

	bindings, err := application.NewBindings(
		[]any{
			&TestService{},
		}, make(map[uint32]uint32),
	)
	if err != nil {
		t.Fatalf("application.NewBindings() error = %v\n", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callOptions := &application.CallOptions{
				PackageName: "application_test",
				StructName:  "TestService",
				MethodName:  tt.method,
			}

			method := bindings.Get(callOptions)
			if method == nil {
				t.Fatalf("bound method not found: %s", callOptions.Name())
			}

			result, err := method.Call(context.TODO(), tt.args)
			if tt.err != err && (tt.err == nil || err == nil || !strings.Contains(err.Error(), tt.err.Error())) {
				t.Fatalf("error: %v, expected error: %v", err, tt.err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Fatalf("result: %v, expected result: %v", result, tt.expected)
			}

		})
	}

}
