//go:build windows

package webview2

import (
	"runtime"
	"testing"
	"unsafe"
)

type stringCallbackProbe struct {
	code   uintptr
	result string
}

func (*stringCallbackProbe) QueryInterface(uintptr, uintptr) uintptr { return 0 }
func (*stringCallbackProbe) AddRef() uintptr                         { return 1 }
func (*stringCallbackProbe) Release() uintptr                        { return 1 }
func (p *stringCallbackProbe) AddScriptToExecuteOnDocumentCreatedCompleted(code uintptr, result string) uintptr {
	p.code, p.result = code, result
	return 0
}
func (p *stringCallbackProbe) ExecuteScriptCompleted(code uintptr, result string) uintptr {
	p.code, p.result = code, result
	return 0
}
func (p *stringCallbackProbe) CallDevToolsProtocolMethodCompleted(code uintptr, result string) uintptr {
	p.code, p.result = code, result
	return 0
}

// Exercise the native callback ABI, including UTF-16 decoding, rather than
// calling the Go implementation directly. A Go string is not a COM argument.
func TestStringCompletionCallbacks(t *testing.T) {
	for _, value := range []string{"", "result: 日本語 🐇"} {
		p := &stringCallbackProbe{}
		text, err := UTF16PtrFromString(value)
		if err != nil {
			t.Fatal(err)
		}
		add := NewICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler(p)
		script := NewICoreWebView2ExecuteScriptCompletedHandler(p)
		devtools := NewICoreWebView2CallDevToolsProtocolMethodCompletedHandler(p)
		for _, tc := range []struct {
			name     string
			callback ComProc
			handler  unsafe.Pointer
		}{
			{"add script", add.Vtbl.Invoke, unsafe.Pointer(add)},
			{"execute script", script.Vtbl.Invoke, unsafe.Pointer(script)},
			{"devtools", devtools.Vtbl.Invoke, unsafe.Pointer(devtools)},
		} {
			t.Run(tc.name+value, func(t *testing.T) {
				p.result = "unset"
				result, _, _ := tc.callback.Call(uintptr(tc.handler), 0x80004005, uintptr(unsafe.Pointer(text)))
				runtime.KeepAlive(text)
				runtime.KeepAlive(tc.handler)
				if result != 0 || p.code != 0x80004005 || p.result != value {
					t.Fatalf("result=%d code=%x text=%q", result, p.code, p.result)
				}
			})
		}
	}
}
