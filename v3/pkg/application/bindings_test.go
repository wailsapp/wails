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

func (t *TestService) PositionalAndVariadic(a int, _ ...string) int {
	return a
}

func (t *TestService) Slice(a []int) []int {
	return a
}

func newArgs(jsonArgs ...string) (args []json.RawMessage) {
	for _, j := range jsonArgs {
		args = append(args, json.RawMessage(j))
	}
	return
}

func TestBoundMethodCall(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		args     []json.RawMessage
		err      string
		expected interface{}
	}{
		{
			name:     "nil",
			method:   "Nil",
			args:     []json.RawMessage{},
			err:      "",
			expected: nil,
		},
		{
			name:     "string",
			method:   "String",
			args:     newArgs(`"foo"`),
			err:      "",
			expected: "foo",
		},
		{
			name:     "multiple",
			method:   "Multiple",
			args:     newArgs(`"foo"`, "0", "false"),
			err:      "",
			expected: []interface{}{"foo", 0, false},
		},
		{
			name:     "struct",
			method:   "Struct",
			args:     newArgs(`{ "name": "alice" }`),
			err:      "",
			expected: Person{Name: "alice"},
		},
		{
			name:     "struct, nil error",
			method:   "StructNil",
			args:     newArgs(`{ "name": "alice" }`),
			err:      "",
			expected: Person{Name: "alice"},
		},
		{
			name:     "struct, error",
			method:   "StructError",
			args:     newArgs(`{ "name": "alice" }`),
			err:      "error",
			expected: nil,
		},
		{
			name:     "invalid argument count",
			method:   "Multiple",
			args:     newArgs(`"foo"`),
			err:      "expects 3 arguments, got 1",
			expected: nil,
		},
		{
			name:     "invalid argument type",
			method:   "String",
			args:     newArgs("1"),
			err:      "could not parse",
			expected: nil,
		},
		{
			name:     "variadic, no arguments",
			method:   "Variadic",
			args:     newArgs(`[]`), // variadic parameters are passed as arrays
			err:      "",
			expected: []string{},
		},
		{
			name:     "variadic",
			method:   "Variadic",
			args:     newArgs(`["foo", "bar"]`),
			err:      "",
			expected: []string{"foo", "bar"},
		},
		{
			name:     "positional and variadic",
			method:   "PositionalAndVariadic",
			args:     newArgs("42", `[]`),
			err:      "",
			expected: 42,
		},
		{
			name:     "slice",
			method:   "Slice",
			args:     newArgs(`[1,2,3]`),
			err:      "",
			expected: []int{1, 2, 3},
		},
	}

	// init globalApplication
	_ = application.New(application.Options{})

	bindings := application.NewBindings(nil, nil)

	err := bindings.Add(application.NewService(&TestService{}))
	if err != nil {
		t.Fatalf("bindings.Add() error = %v\n", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callOptions := &application.CallOptions{
				MethodName: "github.com/wailsapp/wails/v3/pkg/application_test.TestService." + tt.method,
			}

			method := bindings.Get(callOptions)
			if method == nil {
				t.Fatalf("bound method not found: %s", callOptions.MethodName)
			}

			result, err := method.Call(context.TODO(), tt.args)
			if (tt.err == "") != (err == nil) || (err != nil && !strings.Contains(err.Error(), tt.err)) {
				expected := tt.err
				if expected == "" {
					expected = "nil"
				}
				t.Fatalf("error: %#v, expected error: %v", err, expected)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Fatalf("result: %v, expected result: %v", result, tt.expected)
			}
		})
	}

}
