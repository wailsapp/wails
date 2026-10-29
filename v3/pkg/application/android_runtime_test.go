package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type androidRuntimeTestError struct{ payload []byte }

func (*androidRuntimeTestError) Error() string { return "service failed" }

type androidRuntimeTestService struct{}

func (*androidRuntimeTestService) Fail() error { return &androidRuntimeTestError{} }
func (*androidRuntimeTestService) Wrapped() error {
	return fmt.Errorf("wrapped: %w", &androidRuntimeTestError{})
}
func (*androidRuntimeTestService) InvalidJSON() error {
	return &androidRuntimeTestError{payload: []byte(`invalid JSON`)}
}
func (*androidRuntimeTestService) Echo(value string) string { return value }
func (*androidRuntimeTestService) Data() map[string]any     { return map[string]any{"value": 42} }
func (*androidRuntimeTestService) Void()                    {}
func (*androidRuntimeTestService) Unencodable() chan int    { return make(chan int) }

func TestAndroidRuntimeCall(t *testing.T) {
	app := New(Options{})
	previous := app.bindings
	t.Cleanup(func() { app.bindings = previous })
	app.bindings = NewBindings(nil, nil)
	const cause = `{"detail":{"message":"try again","values":[1,true,null]},"retry":false}`
	err := app.bindings.Add(NewServiceWithOptions(&androidRuntimeTestService{}, ServiceOptions{
		MarshalError: func(err error) []byte {
			var serviceErr *androidRuntimeTestError
			if errors.As(err, &serviceErr) {
				if serviceErr.payload != nil {
					return serviceErr.payload
				}
				return []byte(cause)
			}
			return nil
		},
	}))
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range []struct {
		method string
		args   string
		want   string
	}{
		{"Fail", `[]`, `{"ok":false,"error":{"kind":"RuntimeError","message":"service failed","cause":` + cause + `}}`},
		{"Wrapped", `[]`, `{"ok":false,"error":{"kind":"RuntimeError","message":"wrapped: service failed","cause":` + cause + `}}`},
		{"InvalidJSON", `[]`, `{"ok":false,"error":"Binding call failed: Bound method returned an error: service failed"}`},
		{"Echo", `[]`, `{"ok":false,"error":{"kind":"TypeError","message":"github.com/wailsapp/wails/v3/pkg/application.androidRuntimeTestService.Echo expects 1 arguments, got 0"}}`},
		{"Echo", `["hello"]`, `{"ok":true,"text":"hello"}`},
		{"Data", `[]`, `{"ok":true,"data":{"value":42}}`},
		{"Void", `[]`, `{"ok":true,"data":null}`},
		{"Unencodable", `[]`, `{"ok":false,"error":"unable to marshal response: json: unsupported type: chan int"}`},
	} {
		t.Run(tt.method, func(t *testing.T) {
			var methodID uint32
			found := false
			for _, method := range app.bindings.boundMethods {
				if method.Name == tt.method {
					methodID = method.ID
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("bound method %q not found", tt.method)
			}
			payload := fmt.Sprintf(`{"object":0,"method":0,"args":{"call-id":"android-test","methodID":%d,"args":%s}}`, methodID, tt.args)
			assertAndroidEnvelope(t, handleRuntimeCallForAndroid(app, payload), tt.want)
		})
	}
}

func TestAndroidRuntimeCallLegacyErrors(t *testing.T) {
	app := New(Options{})
	for _, tt := range []struct{ name, payload, want string }{
		{"malformed JSON", `{`, `{"ok":false,"error":"unable to parse runtime call: unexpected end of JSON input"}`},
		{"missing object", `{"method":0}`, `{"ok":false,"error":"missing object value"}`},
		{"missing method", `{"object":0}`, `{"ok":false,"error":"missing method value"}`},
		{"invalid binding call", `{"object":0,"method":0,"args":{}}`, `{"ok":false,"error":"Invalid binding call: missing argument 'call-id'"}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assertAndroidEnvelope(t, handleRuntimeCallForAndroid(app, tt.payload), tt.want)
		})
	}
}

func assertAndroidEnvelope(t *testing.T, got, want string) {
	t.Helper()
	var gotJSON, wantJSON any
	if err := json.Unmarshal([]byte(got), &gotJSON); err != nil {
		t.Fatalf("invalid response %q: %v", got, err)
	}
	if err := json.Unmarshal([]byte(want), &wantJSON); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotJSON, wantJSON) {
		t.Fatalf("response = %s, want %s", got, want)
	}
}
