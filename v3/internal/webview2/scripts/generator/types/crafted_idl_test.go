package types_test

// Crafted IDLs driven through the real parser (updater/generator) to exercise
// the generator's error-propagation arms: parse-but-fail-Process inputs,
// parse-but-fail-Generate inputs, and template faults injected via the seam
// exported from the types package.

import (
	"strings"
	"testing"

	"updater/generator"
	"updater/generator/types"
)

const (
	libHeader = "[uuid(11111111-1111-1111-1111-111111111111), version(1.0)]\nlibrary WebView2 {\n"
	ifaceUUID = "[uuid(22222222-2222-2222-2222-222222222222), object, pointer_default(unique)]\n"
)

// richIDL declares an enum, a struct, and an interface with both an [out]
// string method (cleanup template) and an [in] string method (setup template).
const richIDL = libHeader +
	"[v1_enum] typedef enum E { A, B } E;\n" +
	"typedef struct S { UINT32 x; } S;\n" +
	ifaceUUID +
	"interface IFoo : IUnknown {\n" +
	"  [propget] HRESULT Name([out, retval] LPWSTR* name);\n" +
	"  HRESULT SetName([in] LPCWSTR name);\n" +
	"}\n}\n"

func TestRichIDLParsesCleanly(t *testing.T) {
	if _, err := generator.ParseIDLWithTests([]byte(richIDL)); err != nil {
		t.Fatalf("rich IDL should parse and generate cleanly: %v", err)
	}
}

// TestGenerateTemplateFaults breaks one template at a time and confirms the
// error propagates out of ParseIDL through every wrapper in the chain.
func TestGenerateTemplateFaults(t *testing.T) {
	cases := []struct {
		name, template, mode string
	}{
		{"default-files", "com.tmpl", "missing"},
		{"enum", "enum.tmpl", "missing"},
		{"struct", "struct.tmpl", "missing"},
		{"interface-vtbl", "interfacevtbl.tmpl", "missing"},
		{"interface-methods", "interfaceMethod.tmpl", "parse"},
		{"interface-methods-exec", "interfaceMethod.tmpl", "exec"},
		{"param-setup", "inputStringSetup.tmpl", "missing"},
		{"param-cleanup", "outputStringCleanup.tmpl", "missing"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			restore := types.SetTemplateSourceForTest(types.FaultFS(tc.template, tc.mode))
			defer restore()
			if _, err := generator.ParseIDL([]byte(richIDL)); err == nil {
				t.Errorf("breaking %s (%s) should fail generation", tc.template, tc.mode)
			}
		})
	}
}

// TestGenerateInvokeTemplateFault covers generateInvoke's error arm: it needs a
// handler interface (one with an Invoke method) and a broken invoke template.
func TestGenerateInvokeTemplateFault(t *testing.T) {
	const handlerIDL = libHeader + ifaceUUID +
		"interface IHandler : IUnknown {\n" +
		"  HRESULT Invoke([in] LPCWSTR arg);\n" +
		"}\n}\n"
	restore := types.SetTemplateSourceForTest(types.FaultFS("interfaceInvoke.tmpl", "missing"))
	defer restore()
	if _, err := generator.ParseIDL([]byte(handlerIDL)); err == nil {
		t.Error("breaking interfaceInvoke.tmpl should fail generation")
	}
}

// TestUndeclaredBaseClass: parses, but Generate fails because the base
// interface is never declared.
func TestUndeclaredBaseClass(t *testing.T) {
	const idl = libHeader + ifaceUUID +
		"interface IFoo : INotDeclared {\n  HRESULT Ping();\n}\n}\n"
	_, err := generator.ParseIDL([]byte(idl))
	if err == nil || !strings.Contains(err.Error(), "not declared") {
		t.Errorf("undeclared base = %v, want a 'not declared' error", err)
	}
}

// TestUnknownByValueInput: parses, but Process fails because an [in] parameter
// has an unknown by-value type.
func TestUnknownByValueInput(t *testing.T) {
	const idl = libHeader + ifaceUUID +
		"interface IFoo : IUnknown {\n  HRESULT Take([in] Mystery m);\n}\n}\n"
	if _, err := generator.ParseIDL([]byte(idl)); err == nil {
		t.Error("unknown by-value [in] parameter should fail Process")
	}
}

// TestUnknownByValueOutput: same, for an [out] parameter (exercises the
// output-param setup error arm).
func TestUnknownByValueOutput(t *testing.T) {
	const idl = libHeader + ifaceUUID +
		"interface IFoo : IUnknown {\n  HRESULT Take([out] Mystery m);\n}\n}\n"
	if _, err := generator.ParseIDL([]byte(idl)); err == nil {
		t.Error("unknown by-value [out] parameter should fail Process")
	}
}

// TestVoidMethod: a non-HRESULT return type exercises the GetHResultVariable /
// ErrorValuesHRESULT / SuccessValues branches for non-HRESULT methods.
func TestVoidMethod(t *testing.T) {
	const idl = libHeader + ifaceUUID +
		"interface IFoo : IUnknown {\n  void Ping([out, retval] UINT32* n);\n}\n}\n"
	if _, err := generator.ParseIDL([]byte(idl)); err != nil {
		t.Fatalf("void-returning method should generate: %v", err)
	}
}

// TestStructByValueInputTest: generation succeeds, but the test emitter has no
// input strategy for a by-value struct parameter, so GenerateTests fails.
func TestStructByValueInputTest(t *testing.T) {
	const idl = libHeader +
		"typedef struct S { UINT32 x; } S;\n" + ifaceUUID +
		"interface IFoo : IUnknown {\n  HRESULT Take([in] S s);\n}\n}\n"

	if _, err := generator.ParseIDL([]byte(idl)); err != nil {
		t.Fatalf("by-value struct input should generate bindings fine: %v", err)
	}
	if _, err := generator.ParseIDLWithTests([]byte(idl)); err == nil {
		t.Error("by-value struct input should have no test-emitter strategy")
	}
}

// TestInvokeStructArgTest: a handler whose Invoke takes a by-value struct has
// no invoke-dispatch test strategy, so GenerateTests fails.
func TestInvokeStructArgTest(t *testing.T) {
	const idl = libHeader +
		"typedef struct S { UINT32 x; } S;\n" + ifaceUUID +
		"interface IHandler : IUnknown {\n  HRESULT Invoke([in] S arg);\n}\n}\n"

	if _, err := generator.ParseIDL([]byte(idl)); err != nil {
		t.Fatalf("handler with struct arg should generate bindings fine: %v", err)
	}
	if _, err := generator.ParseIDLWithTests([]byte(idl)); err == nil {
		t.Error("invoke dispatch test should have no strategy for a by-value struct arg")
	}
}
