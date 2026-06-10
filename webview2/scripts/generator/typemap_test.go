package generator

// typemap_test.go validates all 17 IDL→Go type-mapping patterns by generating
// code from minimal IDL snippets and checking the output.

import (
	"strings"
	"testing"
	"updater/generator/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// wrapIDL wraps method declarations inside a minimal library/interface block so
// the parser can process them.
func wrapIDL(methods string) []byte {
	return []byte(`[uuid(26d34152-879f-4065-bea2-3daa2cfadfb8), version(1.0)]
library WebView2 {
[uuid(aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa), object, pointer_default(unique)]
interface ITest : IUnknown {
` + methods + `
}
}`)
}

// generatedBody parses the IDL, processes it, generates code, then returns the
// content of the first non-com.go file as a string.
func generatedBody(t *testing.T, idl []byte) string {
	t.Helper()
	files, err := ParseIDL(idl)
	require.NoError(t, err)
	require.Greater(t, len(files), 1, "expected at least one generated file beyond com.go")
	return files[1].Content.String()
}

// ── Pattern 1: [in] LPWSTR ───────────────────────────────────────────────────

func TestTypePattern_InLPWSTR(t *testing.T) {
	body := generatedBody(t, wrapIDL(
		`HRESULT Navigate([in] LPWSTR uri);`,
	))
	assert.Contains(t, body, "Navigate(uri string) error")
	assert.Contains(t, body, "_uri, err := UTF16PtrFromString(uri)")
	assert.Contains(t, body, "uintptr(unsafe.Pointer(_uri))")
}

// ── Pattern 2: [in] LPCWSTR ──────────────────────────────────────────────────

func TestTypePattern_InLPCWSTR(t *testing.T) {
	body := generatedBody(t, wrapIDL(
		`HRESULT Post([in] LPCWSTR json);`,
	))
	assert.Contains(t, body, "Post(json string) error")
	assert.Contains(t, body, "_json, err := UTF16PtrFromString(json)")
	assert.Contains(t, body, "uintptr(unsafe.Pointer(_json))")
}

// ── Pattern 3: [out, retval] LPWSTR* ─────────────────────────────────────────

func TestTypePattern_OutRetvalLPWSTR(t *testing.T) {
	body := generatedBody(t, wrapIDL(
		`[propget] HRESULT Source([out, retval] LPWSTR* source);`,
	))
	assert.Contains(t, body, "GetSource() (string, error)")
	assert.Contains(t, body, "var _source *uint16")
	// Must pass address-of so COM can write back the pointer.
	assert.Contains(t, body, "uintptr(unsafe.Pointer(&_source))")
	assert.Contains(t, body, "UTF16PtrToString(_source)")
	assert.Contains(t, body, "CoTaskMemFree(unsafe.Pointer(_source))")
}

// ── Pattern 4: [in] LPCWSTR* (string array) ──────────────────────────────────

func TestTypePattern_InLPCWSTRStar(t *testing.T) {
	body := generatedBody(t, wrapIDL(
		`HRESULT SetOrigins([in] UINT32 count, [in] LPCWSTR* origins);`,
	))
	assert.Contains(t, body, "origins []string")
	assert.Contains(t, body, "_originsptrs := make([]*uint16, len(origins))")
	assert.Contains(t, body, "var _origins **uint16")
	assert.Contains(t, body, "uintptr(unsafe.Pointer(_origins))")
	// count should be passed by value
	assert.Contains(t, body, "uintptr(count)")
	assert.NotContains(t, body, "unsafe.Pointer(&count)")
}

// ── Pattern 5: [out] LPWSTR** (string array output) ──────────────────────────

func TestTypePattern_OutLPWSTRDoubleStar(t *testing.T) {
	body := generatedBody(t, wrapIDL(
		`HRESULT GetOrigins([out] UINT32* count, [out] LPWSTR** origins);`,
	))
	// The outer pointer is stripped; the Go type is *string.
	assert.Contains(t, body, "origins *string")
	assert.Contains(t, body, "uintptr(unsafe.Pointer(&origins))")
}

// ── Pattern 6: [in] UINT32 ───────────────────────────────────────────────────

func TestTypePattern_InUINT32(t *testing.T) {
	body := generatedBody(t, wrapIDL(
		`HRESULT SetCount([in] UINT32 count);`,
	))
	assert.Contains(t, body, "SetCount(count uint32) error")
	assert.Contains(t, body, "uintptr(count)")
	assert.NotContains(t, body, "unsafe.Pointer(&count)")
}

// ── Pattern 7: [out] UINT32* ─────────────────────────────────────────────────

func TestTypePattern_OutUINT32(t *testing.T) {
	body := generatedBody(t, wrapIDL(
		`[propget] HRESULT Count([out, retval] UINT32* count);`,
	))
	assert.Contains(t, body, "GetCount() (uint32, error)")
	assert.Contains(t, body, "var count uint32")
	assert.Contains(t, body, "uintptr(unsafe.Pointer(&count))")
}

// ── Pattern 8: [in] UINT64 ───────────────────────────────────────────────────

func TestTypePattern_InUINT64(t *testing.T) {
	body := generatedBody(t, wrapIDL(
		`HRESULT SetSize([in] UINT64 size);`,
	))
	assert.Contains(t, body, "SetSize(size uint64) error")
	assert.Contains(t, body, "uintptr(size)")
	// 386 splits 8-byte integers into two stack words, low word first.
	assert.Contains(t, body, "case archIs386:")
	assert.Contains(t, body, "uintptr(uint32(uint64(size)))")
	assert.Contains(t, body, "uintptr(uint32(uint64(size)>>32))")
	assert.NotContains(t, body, "case archIsARM64:")
}

// ── Pattern 9: [in] INT32 ────────────────────────────────────────────────────

func TestTypePattern_InINT32(t *testing.T) {
	body := generatedBody(t, wrapIDL(
		`HRESULT SetLevel([in] INT32 level);`,
	))
	assert.Contains(t, body, "SetLevel(level int32) error")
	assert.Contains(t, body, "uintptr(level)")
}

// ── Pattern 10: [out, retval] INT32* ─────────────────────────────────────────

func TestTypePattern_OutINT32(t *testing.T) {
	body := generatedBody(t, wrapIDL(
		`[propget] HRESULT Level([out, retval] INT32* level);`,
	))
	assert.Contains(t, body, "GetLevel() (int32, error)")
	assert.Contains(t, body, "var level int32")
	assert.Contains(t, body, "uintptr(unsafe.Pointer(&level))")
}

// ── Pattern 11: [out, retval] int* ───────────────────────────────────────────

func TestTypePattern_OutInt(t *testing.T) {
	body := generatedBody(t, wrapIDL(
		`[propget] HRESULT ExitCode([out, retval] int* exitCode);`,
	))
	assert.Contains(t, body, "GetExitCode() (int, error)")
	assert.Contains(t, body, "var exitCode int")
	assert.Contains(t, body, "uintptr(unsafe.Pointer(&exitCode))")
	// Must NOT pass by value (old bug).
	assert.NotContains(t, body, "uintptr(exitCode),")
}

// ── Pattern 12: [in] BOOL ────────────────────────────────────────────────────

func TestTypePattern_InBOOL(t *testing.T) {
	body := generatedBody(t, wrapIDL(
		`[propput] HRESULT Visible([in] BOOL value);`,
	))
	assert.Contains(t, body, "PutVisible(value bool) error")
	// Must convert to int32 before the vtable call.
	assert.Contains(t, body, "var _value int32")
	assert.Contains(t, body, "uintptr(_value)")
	// Must NOT pass a pointer to a Go bool.
	assert.NotContains(t, body, "unsafe.Pointer(&value)")
}

// ── Pattern 13: [out, retval] BOOL* ──────────────────────────────────────────

func TestTypePattern_OutBOOL(t *testing.T) {
	body := generatedBody(t, wrapIDL(
		`[propget] HRESULT Visible([out, retval] BOOL* visible);`,
	))
	assert.Contains(t, body, "GetVisible() (bool, error)")
	assert.Contains(t, body, "var _visible int32")
	assert.Contains(t, body, "uintptr(unsafe.Pointer(&_visible))")
	assert.Contains(t, body, "visible := _visible != 0")
}

// ── Pattern 14: [in] double ──────────────────────────────────────────────────

func TestTypePattern_InDouble(t *testing.T) {
	body := generatedBody(t, wrapIDL(
		`HRESULT SetZoom([in] double factor);`,
	))
	assert.Contains(t, body, "SetZoom(factor float64) error")
	// The IEEE-754 bit pattern must be passed, never a truncating integer
	// conversion: uintptr(factor) would turn 1.5 into 1.
	assert.NotContains(t, body, "uintptr(factor)")
	assert.Contains(t, body, "uintptr(math.Float64bits(factor))")
	// On 386 a double occupies two 4-byte stack slots, low word first.
	assert.Contains(t, body, "case archIs386:")
	assert.Contains(t, body, "uintptr(uint32(math.Float64bits(factor)))")
	assert.Contains(t, body, "uintptr(uint32(math.Float64bits(factor)>>32))")
	// amd64 and arm64 encode doubles identically — no arm64-only branch.
	assert.NotContains(t, body, "case archIsARM64:")
	assert.Contains(t, body, `"math"`)
}

// ── Pattern 14b: [in] EventRegistrationToken (by value) ──────────────────────

func TestTypePattern_InEventRegistrationToken(t *testing.T) {
	body := generatedBody(t, wrapIDL(
		`HRESULT remove_NavigationStarting([in] EventRegistrationToken token);`,
	))
	assert.Contains(t, body, "RemoveNavigationStarting(token EventRegistrationToken) error")
	// Win64 passes 8-byte aggregates by value in a register — never by address.
	assert.NotContains(t, body, "uintptr(unsafe.Pointer(&token))")
	assert.Contains(t, body, "uintptr(*(*uint64)(unsafe.Pointer(&token)))")
	// 386 pushes the two 4-byte halves as separate stack words.
	assert.Contains(t, body, "uintptr((*(*[2]uint32)(unsafe.Pointer(&token)))[0])")
	assert.Contains(t, body, "uintptr((*(*[2]uint32)(unsafe.Pointer(&token)))[1])")
	assert.NotContains(t, body, "case archIsARM64:")
}

// ── Pattern 14c: [in] 4-byte struct by value (COREWEBVIEW2_COLOR) ────────────

func TestTypePattern_InSmallStructByValue(t *testing.T) {
	idl := []byte(`[uuid(26d34152-879f-4065-bea2-3daa2cfadfb8), version(1.0)]
library WebView2 {
typedef struct COREWEBVIEW2_COLOR {
	BYTE A;
	BYTE R;
	BYTE G;
	BYTE B;
} COREWEBVIEW2_COLOR;
[uuid(aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa), object, pointer_default(unique)]
interface ITest : IUnknown {
HRESULT put_DefaultBackgroundColor([in] COREWEBVIEW2_COLOR value);
}
}`)
	files, err := ParseIDL(idl)
	require.NoError(t, err)
	var body string
	for _, f := range files {
		if f.FileName == "ITest.go" {
			body = f.Content.String()
		}
	}
	require.NotEmpty(t, body)
	// A 4-byte struct is a single register/stack word on every architecture,
	// packed by value — never passed by address.
	assert.NotContains(t, body, "uintptr(unsafe.Pointer(&value))")
	assert.Contains(t, body, "uintptr(*(*uint32)(unsafe.Pointer(&value)))")
	// Same encoding everywhere — no arch switch needed.
	assert.NotContains(t, body, "case archIs386:")
}

// ── Pattern 14d: [in] POINT (8-byte struct by value) ─────────────────────────

func TestTypePattern_InPointByValue(t *testing.T) {
	body := generatedBody(t, wrapIDL(
		`HRESULT DragEnter([in] POINT point);`,
	))
	assert.NotContains(t, body, "uintptr(unsafe.Pointer(&point))")
	assert.Contains(t, body, "uintptr(*(*uint64)(unsafe.Pointer(&point)))")
	assert.Contains(t, body, "case archIs386:")
	assert.NotContains(t, body, "case archIsARM64:")
}

// ── Pattern 14e: [in] RECT (16-byte struct by value) ─────────────────────────

func TestTypePattern_InRectByValue(t *testing.T) {
	body := generatedBody(t, wrapIDL(
		`[propput] HRESULT Bounds([in] RECT bounds);`,
	))
	// amd64 passes >8-byte aggregates via pointer to the (already private) copy.
	assert.Contains(t, body, "uintptr(unsafe.Pointer(&bounds))")
	// arm64 packs 9-16 byte composites into a register pair.
	assert.Contains(t, body, "case archIsARM64:")
	assert.Contains(t, body, "(*(*[2]uintptr)(unsafe.Pointer(&bounds)))[0]")
	assert.Contains(t, body, "(*(*[2]uintptr)(unsafe.Pointer(&bounds)))[1]")
	// 386 pushes four 4-byte stack words.
	assert.Contains(t, body, "uintptr((*(*[4]uint32)(unsafe.Pointer(&bounds)))[3])")
}

// ── Pattern 15: [in] IInterface* ─────────────────────────────────────────────

func TestTypePattern_InInterfacePointer(t *testing.T) {
	body := generatedBody(t, wrapIDL(
		`HRESULT add_Event([in] ICoreWebView2NavigationStartingEventHandler* handler, [out] EventRegistrationToken* token);`,
	))
	assert.Contains(t, body, "handler *ICoreWebView2NavigationStartingEventHandler")
	assert.Contains(t, body, "uintptr(unsafe.Pointer(handler))")
}

// ── Pattern 16: [out] IInterface** ───────────────────────────────────────────

func TestTypePattern_OutInterfaceDoublePointer(t *testing.T) {
	body := generatedBody(t, wrapIDL(
		`[propget] HRESULT Settings([out, retval] ICoreWebView2Settings** settings);`,
	))
	assert.Contains(t, body, "GetSettings() (*ICoreWebView2Settings, error)")
	assert.Contains(t, body, "var settings *ICoreWebView2Settings")
	assert.Contains(t, body, "uintptr(unsafe.Pointer(&settings))")
}

// ── Pattern 17: EventRegistrationToken ───────────────────────────────────────

func TestTypePattern_EventRegistrationToken(t *testing.T) {
	body := generatedBody(t, wrapIDL(
		`HRESULT add_Nav([in] ICoreWebView2NavigationStartingEventHandler* h, [out] EventRegistrationToken* token);
HRESULT remove_Nav([in] EventRegistrationToken token);`,
	))
	// [out] EventRegistrationToken* — output via pointer
	assert.Contains(t, body, "var token EventRegistrationToken")
	assert.Contains(t, body, "uintptr(unsafe.Pointer(&token))")
	// [in] EventRegistrationToken — input as value (passed by ref in vtable)
	assert.Contains(t, body, "RemoveNav(token EventRegistrationToken) error")
}

// ── ResolveGoType unit tests ──────────────────────────────────────────────────

func TestResolveGoType(t *testing.T) {
	cases := []struct {
		idlType   string
		pointer   string
		direction string
		want      string
	}{
		{"LPWSTR", "", "in", "string"},
		{"LPCWSTR", "", "in", "string"},
		{"LPWSTR", "*", "out", "string"},
		{"LPCWSTR", "*", "in", "[]string"},   // string array
		{"LPWSTR", "*", "in", "[]string"},    // string array (writable variant)
		{"UINT32", "", "in", "uint32"},
		{"UINT32", "*", "out", "uint32"},
		{"INT32", "", "in", "int32"},
		{"UINT64", "", "in", "uint64"},
		{"BOOL", "", "in", "bool"},
		{"BOOL", "*", "out", "bool"},
		{"double", "", "in", "float64"},
		{"HRESULT", "", "in", "uintptr"},
		{"BYTE", "", "in", "uint8"},
		{"DWORD", "", "in", "uint32"},
		{"IUnknown", "", "in", "IUnknown"},
		{"EventRegistrationToken", "", "in", "EventRegistrationToken"},
	}

	for _, tc := range cases {
		tc := tc
		name := strings.Join([]string{tc.direction, tc.idlType, tc.pointer}, "_")
		t.Run(name, func(t *testing.T) {
			got := types.ResolveGoType(tc.idlType, tc.pointer, tc.direction)
			assert.Equal(t, tc.want, got)
		})
	}
}
