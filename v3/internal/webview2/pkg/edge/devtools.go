//go:build windows

package edge

import (
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v3/internal/webview2/pkg/combridge"
	"golang.org/x/sys/windows"
)

// These callbacks use combridge so COM reference ownership also keeps their Go
// closures alive. All methods and callbacks run on the WebView's UI thread.
type devToolsCompletion interface{ Complete(int32, *uint16) }
type devToolsEvent interface{ Receive(*devToolsEventArgs) }
type devToolsCompletionFunc func(int32, *uint16)

func (f devToolsCompletionFunc) Complete(hr int32, result *uint16) { f(hr, result) }

type devToolsEventFunc func(*devToolsEventArgs)

func (f devToolsEventFunc) Receive(args *devToolsEventArgs) { f(args) }

func init() {
	combridge.RegisterVTable[combridge.IUnknown, devToolsCompletion]("{5c4889f0-5ef6-4c5a-952c-d8f1b92d0574}",
		func(this uintptr, hr int32, result *uint16) uintptr {
			combridge.Resolve[devToolsCompletion](this).Complete(hr, result)
			return 0
		})
	combridge.RegisterVTable[combridge.IUnknown, devToolsEvent]("{e2fda4be-5456-406c-a261-3d452138362c}",
		func(this uintptr, _ *ICoreWebView2, args *devToolsEventArgs) uintptr {
			combridge.Resolve[devToolsEvent](this).Receive(args)
			return 0
		})
}

type devToolsEventArgs struct {
	vtbl *struct {
		_IUnknownVtbl
		GetJSON ComProc
	}
}
type devToolsEventReceiver struct {
	vtbl *struct {
		_IUnknownVtbl
		Add, Remove ComProc
	}
}

func (a *devToolsEventArgs) json() (string, error) {
	var value *uint16
	hr, _, _ := a.vtbl.GetJSON.Call(uintptr(unsafe.Pointer(a)), uintptr(unsafe.Pointer(&value)))
	if int32(hr) < 0 {
		return "", syscall.Errno(hr)
	}
	defer windows.CoTaskMemFree(unsafe.Pointer(value))
	return windows.UTF16PtrToString(value), nil
}

// CallDevTools invokes CDP in process; it does not enable a remote debugging port.
// done receives asynchronous protocol errors as well as the result JSON.
func (e *Chromium) CallDevTools(method, params string, done func(string, error)) error {
	return e.CallDevToolsForSession("", method, params, done)
}

// CallDevToolsForSession addresses a worker or frame session attached with Target.setAutoAttach.
func (e *Chromium) CallDevToolsForSession(session, method, params string, done func(string, error)) error {
	m, err := windows.UTF16PtrFromString(method)
	if err != nil {
		return err
	}
	p, err := windows.UTF16PtrFromString(params)
	if err != nil {
		return err
	}
	callback := combridge.New[devToolsCompletion](devToolsCompletionFunc(func(hr int32, result *uint16) {
		if hr < 0 {
			done("", syscall.Errno(uint32(hr)))
			return
		}
		done(windows.UTF16PtrToString(result), nil)
	}))
	defer callback.Close()
	var hr uintptr
	if session == "" {
		hr, _, _ = e.webview.vtbl.CallDevToolsProtocolMethod.Call(uintptr(unsafe.Pointer(e.webview)), uintptr(unsafe.Pointer(m)), uintptr(unsafe.Pointer(p)), callback.Ref())
	} else {
		sid, err := windows.UTF16PtrFromString(session)
		if err != nil {
			return err
		}
		guid := windows.GUID{Data1: 0x0be78e56, Data2: 0xc193, Data3: 0x4051, Data4: [8]byte{0xb9, 0x43, 0x23, 0xb4, 0x60, 0xc0, 0x8b, 0xdb}}
		var view *devToolsWebView11
		hr, _, _ = e.webview.vtbl.QueryInterface.Call(uintptr(unsafe.Pointer(e.webview)), uintptr(unsafe.Pointer(&guid)), uintptr(unsafe.Pointer(&view)))
		if int32(hr) < 0 {
			return syscall.Errno(hr)
		}
		defer view.vtbl.Release.Call(uintptr(unsafe.Pointer(view)))
		hr, _, _ = view.vtbl.CallForSession.Call(uintptr(unsafe.Pointer(view)), uintptr(unsafe.Pointer(sid)), uintptr(unsafe.Pointer(m)), uintptr(unsafe.Pointer(p)), callback.Ref())
	}
	if int32(hr) < 0 {
		return syscall.Errno(hr)
	}
	return nil
}

// OnDevTools subscribes to one CDP event. The returned cleanup must run on the UI thread.
func (e *Chromium) OnDevTools(event string, receive func(string, string, error)) (func() error, error) {
	name, err := windows.UTF16PtrFromString(event)
	if err != nil {
		return nil, err
	}
	var receiver *devToolsEventReceiver
	hr, _, _ := e.webview.vtbl.GetDevToolsProtocolEventReceiver.Call(uintptr(unsafe.Pointer(e.webview)), uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(&receiver)))
	if int32(hr) < 0 {
		return nil, syscall.Errno(hr)
	}
	callback := combridge.New[devToolsEvent](devToolsEventFunc(func(args *devToolsEventArgs) {
		session, err := args.session()
		if err != nil {
			receive("", "", err)
			return
		}
		data, err := args.json()
		receive(session, data, err)
	}))
	defer callback.Close()
	var token _EventRegistrationToken
	hr, _, _ = receiver.vtbl.Add.Call(uintptr(unsafe.Pointer(receiver)), callback.Ref(), uintptr(unsafe.Pointer(&token)))
	if int32(hr) < 0 {
		receiver.vtbl.Release.Call(uintptr(unsafe.Pointer(receiver)))
		return nil, syscall.Errno(hr)
	}
	return func() error {
		if receiver == nil {
			return nil
		}
		var hr uintptr
		if unsafe.Sizeof(uintptr(0)) == 4 {
			hr, _, _ = receiver.vtbl.Remove.Call(uintptr(unsafe.Pointer(receiver)), uintptr(uint32(token.value)), uintptr(uint64(token.value)>>32))
		} else {
			hr, _, _ = receiver.vtbl.Remove.Call(uintptr(unsafe.Pointer(receiver)), uintptr(token.value))
		}
		receiver.vtbl.Release.Call(uintptr(unsafe.Pointer(receiver)))
		receiver = nil
		if int32(hr) < 0 {
			return syscall.Errno(hr)
		}
		return nil
	}, nil
}

// ICoreWebView2_2 through _10 add 38 methods in WebView2.1.0.2903.40.idl.
// _11 starts with CallDevToolsProtocolMethodForSession. This layout is shared by
// x86, x64 and arm64; every slot is one COM procedure pointer.
type devToolsWebView11 struct {
	vtbl *struct {
		iCoreWebView2Vtbl
		versions2Through10 [38]ComProc
		CallForSession     ComProc
	}
}

type devToolsEventArgs2 struct {
	vtbl *struct {
		_IUnknownVtbl
		GetJSON, GetSession ComProc
	}
}

func (a *devToolsEventArgs) session() (string, error) {
	guid := windows.GUID{Data1: 0x2dc4959d, Data2: 0x1494, Data3: 0x4393, Data4: [8]byte{0x95, 0xba, 0xbe, 0xa4, 0xcb, 0x9e, 0xbd, 0x1b}}
	var extended *devToolsEventArgs2
	hr, _, _ := a.vtbl.QueryInterface.Call(uintptr(unsafe.Pointer(a)), uintptr(unsafe.Pointer(&guid)), uintptr(unsafe.Pointer(&extended)))
	if uint32(hr) == uint32(windows.E_NOINTERFACE) {
		return "", nil
	}
	if int32(hr) < 0 {
		return "", syscall.Errno(hr)
	}
	defer extended.vtbl.Release.Call(uintptr(unsafe.Pointer(extended)))
	var value *uint16
	hr, _, _ = extended.vtbl.GetSession.Call(uintptr(unsafe.Pointer(extended)), uintptr(unsafe.Pointer(&value)))
	if int32(hr) < 0 {
		return "", syscall.Errno(hr)
	}
	defer windows.CoTaskMemFree(unsafe.Pointer(value))
	return windows.UTF16PtrToString(value), nil
}
