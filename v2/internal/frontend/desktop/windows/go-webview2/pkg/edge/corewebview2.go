//go:build windows
// +build windows

package edge

import (
	"fmt"
	"log"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/go-webview2/internal/w32"

	"golang.org/x/sys/windows"
)

func init() {
	runtime.LockOSThread()

	r, _, _ := w32.Ole32CoInitializeEx.Call(0, 2)
	if int(r) < 0 {
		log.Printf("Warning: CoInitializeEx call failed: E=%08x", r)
	}
}

type _EventRegistrationToken struct {
	value int64
}

type CoreWebView2PermissionKind uint32

const (
	CoreWebView2PermissionKindUnknownPermission CoreWebView2PermissionKind = iota
	CoreWebView2PermissionKindMicrophone
	CoreWebView2PermissionKindCamera
	CoreWebView2PermissionKindGeolocation
	CoreWebView2PermissionKindNotifications
	CoreWebView2PermissionKindOtherSensors
	CoreWebView2PermissionKindClipboardRead
)

type CoreWebView2PermissionState uint32

const (
	CoreWebView2PermissionStateDefault CoreWebView2PermissionState = iota
	CoreWebView2PermissionStateAllow
	CoreWebView2PermissionStateDeny
)

// ComProc stores a COM procedure.
type ComProc uintptr

// NewComProc creates a new COM proc from a Go function.
func NewComProc(fn interface{}) ComProc {
	return ComProc(windows.NewCallback(fn))
}

// Call calls a COM procedure.
//
//go:uintptrescapes
func (p ComProc) Call(a ...uintptr) (r1, r2 uintptr, lastErr error) {
	// The magic uintptrescapes comment is needed to prevent moving uintptr(unsafe.Pointer(p)) so calls to .Call() also
	// satisfy the unsafe.Pointer rule "(4) Conversion of a Pointer to a uintptr when calling syscall.Syscall."
	// Otherwise it might be that pointers get moved, especially pointer onto the Go stack which might grow dynamically.
	// See https://pkg.go.dev/unsafe#Pointer and https://github.com/golang/go/issues/34474
	switch len(a) {
	case 0:
		return syscall.Syscall(uintptr(p), 0, 0, 0, 0)
	case 1:
		return syscall.Syscall(uintptr(p), 1, a[0], 0, 0)
	case 2:
		return syscall.Syscall(uintptr(p), 2, a[0], a[1], 0)
	case 3:
		return syscall.Syscall(uintptr(p), 3, a[0], a[1], a[2])
	case 4:
		return syscall.Syscall6(uintptr(p), 4, a[0], a[1], a[2], a[3], 0, 0)
	case 5:
		return syscall.Syscall6(uintptr(p), 5, a[0], a[1], a[2], a[3], a[4], 0)
	case 6:
		return syscall.Syscall6(uintptr(p), 6, a[0], a[1], a[2], a[3], a[4], a[5])
	case 7:
		return syscall.Syscall9(uintptr(p), 7, a[0], a[1], a[2], a[3], a[4], a[5], a[6], 0, 0)
	case 8:
		return syscall.Syscall9(uintptr(p), 8, a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], 0)
	case 9:
		return syscall.Syscall9(uintptr(p), 9, a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8])
	case 10:
		return syscall.Syscall12(uintptr(p), 10, a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], 0, 0)
	case 11:
		return syscall.Syscall12(uintptr(p), 11, a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], 0)
	case 12:
		return syscall.Syscall12(uintptr(p), 12, a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], a[11])
	case 13:
		return syscall.Syscall15(uintptr(p), 13, a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], a[11], a[12], 0, 0)
	case 14:
		return syscall.Syscall15(uintptr(p), 14, a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], a[11], a[12], a[13], 0)
	case 15:
		return syscall.Syscall15(uintptr(p), 15, a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], a[11], a[12], a[13], a[14])
	default:
		panic("too many arguments")
	}
}

// IUnknown

type _IUnknownVtbl struct {
	QueryInterface ComProc
	AddRef         ComProc
	Release        ComProc
}

func (i *_IUnknownVtbl) CallRelease(this unsafe.Pointer) error {
	_, _, err := i.Release.Call(
		uintptr(this),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

type _IUnknownImpl interface {
	QueryInterface(refiid, object uintptr) uintptr
	AddRef() uintptr
	Release() uintptr
}

// ICoreWebView2

type iCoreWebView2Vtbl struct {
	_IUnknownVtbl
	GetSettings                            ComProc
	GetSource                              ComProc
	Navigate                               ComProc
	NavigateToString                       ComProc
	AddNavigationStarting                  ComProc
	RemoveNavigationStarting               ComProc
	AddContentLoading                      ComProc
	RemoveContentLoading                   ComProc
	AddSourceChanged                       ComProc
	RemoveSourceChanged                    ComProc
	AddHistoryChanged                      ComProc
	RemoveHistoryChanged                   ComProc
	AddNavigationCompleted                 ComProc
	RemoveNavigationCompleted              ComProc
	AddFrameNavigationStarting             ComProc
	RemoveFrameNavigationStarting          ComProc
	AddFrameNavigationCompleted            ComProc
	RemoveFrameNavigationCompleted         ComProc
	AddScriptDialogOpening                 ComProc
	RemoveScriptDialogOpening              ComProc
	AddPermissionRequested                 ComProc
	RemovePermissionRequested              ComProc
	AddProcessFailed                       ComProc
	RemoveProcessFailed                    ComProc
	AddScriptToExecuteOnDocumentCreated    ComProc
	RemoveScriptToExecuteOnDocumentCreated ComProc
	ExecuteScript                          ComProc
	CapturePreview                         ComProc
	Reload                                 ComProc
	PostWebMessageAsJSON                   ComProc
	PostWebMessageAsString                 ComProc
	AddWebMessageReceived                  ComProc
	RemoveWebMessageReceived               ComProc
	CallDevToolsProtocolMethod             ComProc
	GetBrowserProcessID                    ComProc
	GetCanGoBack                           ComProc
	GetCanGoForward                        ComProc
	GoBack                                 ComProc
	GoForward                              ComProc
	GetDevToolsProtocolEventReceiver       ComProc
	Stop                                   ComProc
	AddNewWindowRequested                  ComProc
	RemoveNewWindowRequested               ComProc
	AddDocumentTitleChanged                ComProc
	RemoveDocumentTitleChanged             ComProc
	GetDocumentTitle                       ComProc
	AddHostObjectToScript                  ComProc
	RemoveHostObjectFromScript             ComProc
	OpenDevToolsWindow                     ComProc
	AddContainsFullScreenElementChanged    ComProc
	RemoveContainsFullScreenElementChanged ComProc
	GetContainsFullScreenElement           ComProc
	AddWebResourceRequested                ComProc
	RemoveWebResourceRequested             ComProc
	AddWebResourceRequestedFilter          ComProc
	RemoveWebResourceRequestedFilter       ComProc
	AddWindowCloseRequested                ComProc
	RemoveWindowCloseRequested             ComProc
}

type ICoreWebView2 struct {
	vtbl *iCoreWebView2Vtbl
}

func (i *ICoreWebView2) GetSettings() (*ICoreWebViewSettings, error) {
	var err error
	var settings *ICoreWebViewSettings
	_, _, err = i.vtbl.GetSettings.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&settings)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	return settings, nil
}

// ICoreWebView2Environment

type iCoreWebView2EnvironmentVtbl struct {
	_IUnknownVtbl
	CreateCoreWebView2Controller     ComProc
	CreateWebResourceResponse        ComProc
	GetBrowserVersionString          ComProc
	AddNewBrowserVersionAvailable    ComProc
	RemoveNewBrowserVersionAvailable ComProc
}

type ICoreWebView2Environment struct {
	vtbl *iCoreWebView2EnvironmentVtbl
}

// CreateWebResourceResponse creates a new ICoreWebView2WebResourceResponse, it must be released after finishing using it.
func (e *ICoreWebView2Environment) CreateWebResourceResponse(content []byte, statusCode int, reasonPhrase string, headers string) (*ICoreWebView2WebResourceResponse, error) {
	var err error
	var stream uintptr

	if len(content) > 0 {
		// Create stream for response
		stream, err = w32.SHCreateMemStream(content)
		if err != nil {
			return nil, err
		}

		// Release the IStream after we are finished, CreateWebResourceResponse Call will increase the reference
		// count on IStream and therefore it won't be freed until the reference count of the response is 0.
		defer (*IStream)(unsafe.Pointer(stream)).Release()
	}

	// Convert string 'uri' to *uint16
	_reason, err := windows.UTF16PtrFromString(reasonPhrase)
	if err != nil {
		return nil, err
	}
	// Convert string 'uri' to *uint16
	_headers, err := windows.UTF16PtrFromString(headers)
	if err != nil {
		return nil, err
	}
	var response *ICoreWebView2WebResourceResponse
	hr, _, err := e.vtbl.CreateWebResourceResponse.Call(
		uintptr(unsafe.Pointer(e)),
		stream,
		uintptr(statusCode),
		uintptr(unsafe.Pointer(_reason)),
		uintptr(unsafe.Pointer(_headers)),
		uintptr(unsafe.Pointer(&response)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}

	if response == nil {
		if err == nil {
			err = fmt.Errorf("unknown error")
		}
		return nil, err
	}
	return response, nil

}

// ICoreWebView2WebMessageReceivedEventArgs

type iCoreWebView2WebMessageReceivedEventArgsVtbl struct {
	_IUnknownVtbl
	GetSource                ComProc
	GetWebMessageAsJSON      ComProc
	TryGetWebMessageAsString ComProc
}

type iCoreWebView2WebMessageReceivedEventArgs struct {
	vtbl *iCoreWebView2WebMessageReceivedEventArgsVtbl
}

// ICoreWebView2PermissionRequestedEventArgs

type iCoreWebView2PermissionRequestedEventArgsVtbl struct {
	_IUnknownVtbl
	GetURI             ComProc
	GetPermissionKind  ComProc
	GetIsUserInitiated ComProc
	GetState           ComProc
	PutState           ComProc
	GetDeferral        ComProc
}

type iCoreWebView2PermissionRequestedEventArgs struct {
	vtbl *iCoreWebView2PermissionRequestedEventArgsVtbl
}

// ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler

type iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerImpl interface {
	_IUnknownImpl
	EnvironmentCompleted(res uintptr, env *ICoreWebView2Environment) uintptr
}

type iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVtbl struct {
	_IUnknownVtbl
	Invoke ComProc
}

type iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler struct {
	vtbl *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVtbl
	impl iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerImpl
}

func _ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerIUnknownQueryInterface(this *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerIUnknownAddRef(this *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerIUnknownRelease(this *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerInvoke(this *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler, res uintptr, env *ICoreWebView2Environment) uintptr {
	return this.impl.EnvironmentCompleted(res, env)
}

var iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerFn = iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVtbl{
	_IUnknownVtbl{
		NewComProc(_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerIUnknownQueryInterface),
		NewComProc(_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerIUnknownAddRef),
		NewComProc(_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerIUnknownRelease),
	},
	NewComProc(_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerInvoke),
}

func newICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler(impl iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerImpl) *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler {
	return &iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler{
		vtbl: &iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerFn,
		impl: impl,
	}
}

// ICoreWebView2WebMessageReceivedEventHandler

type iCoreWebView2WebMessageReceivedEventHandlerImpl interface {
	_IUnknownImpl
	MessageReceived(sender *ICoreWebView2, args *iCoreWebView2WebMessageReceivedEventArgs) uintptr
}

type iCoreWebView2WebMessageReceivedEventHandlerVtbl struct {
	_IUnknownVtbl
	Invoke ComProc
}

type iCoreWebView2WebMessageReceivedEventHandler struct {
	vtbl *iCoreWebView2WebMessageReceivedEventHandlerVtbl
	impl iCoreWebView2WebMessageReceivedEventHandlerImpl
}

func _ICoreWebView2WebMessageReceivedEventHandlerIUnknownQueryInterface(this *iCoreWebView2WebMessageReceivedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2WebMessageReceivedEventHandlerIUnknownAddRef(this *iCoreWebView2WebMessageReceivedEventHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2WebMessageReceivedEventHandlerIUnknownRelease(this *iCoreWebView2WebMessageReceivedEventHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2WebMessageReceivedEventHandlerInvoke(this *iCoreWebView2WebMessageReceivedEventHandler, sender *ICoreWebView2, args *iCoreWebView2WebMessageReceivedEventArgs) uintptr {
	return this.impl.MessageReceived(sender, args)
}

var iCoreWebView2WebMessageReceivedEventHandlerFn = iCoreWebView2WebMessageReceivedEventHandlerVtbl{
	_IUnknownVtbl{
		NewComProc(_ICoreWebView2WebMessageReceivedEventHandlerIUnknownQueryInterface),
		NewComProc(_ICoreWebView2WebMessageReceivedEventHandlerIUnknownAddRef),
		NewComProc(_ICoreWebView2WebMessageReceivedEventHandlerIUnknownRelease),
	},
	NewComProc(_ICoreWebView2WebMessageReceivedEventHandlerInvoke),
}

func newICoreWebView2WebMessageReceivedEventHandler(impl iCoreWebView2WebMessageReceivedEventHandlerImpl) *iCoreWebView2WebMessageReceivedEventHandler {
	return &iCoreWebView2WebMessageReceivedEventHandler{
		vtbl: &iCoreWebView2WebMessageReceivedEventHandlerFn,
		impl: impl,
	}
}

// ICoreWebView2PermissionRequestedEventHandler

type iCoreWebView2PermissionRequestedEventHandlerImpl interface {
	_IUnknownImpl
	PermissionRequested(sender *ICoreWebView2, args *iCoreWebView2PermissionRequestedEventArgs) uintptr
}

type iCoreWebView2PermissionRequestedEventHandlerVtbl struct {
	_IUnknownVtbl
	Invoke ComProc
}

type iCoreWebView2PermissionRequestedEventHandler struct {
	vtbl *iCoreWebView2PermissionRequestedEventHandlerVtbl
	impl iCoreWebView2PermissionRequestedEventHandlerImpl
}

func _ICoreWebView2PermissionRequestedEventHandlerIUnknownQueryInterface(this *iCoreWebView2PermissionRequestedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2PermissionRequestedEventHandlerIUnknownAddRef(this *iCoreWebView2PermissionRequestedEventHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2PermissionRequestedEventHandlerIUnknownRelease(this *iCoreWebView2PermissionRequestedEventHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2PermissionRequestedEventHandlerInvoke(this *iCoreWebView2PermissionRequestedEventHandler, sender *ICoreWebView2, args *iCoreWebView2PermissionRequestedEventArgs) uintptr {
	return this.impl.PermissionRequested(sender, args)
}

var iCoreWebView2PermissionRequestedEventHandlerFn = iCoreWebView2PermissionRequestedEventHandlerVtbl{
	_IUnknownVtbl{
		NewComProc(_ICoreWebView2PermissionRequestedEventHandlerIUnknownQueryInterface),
		NewComProc(_ICoreWebView2PermissionRequestedEventHandlerIUnknownAddRef),
		NewComProc(_ICoreWebView2PermissionRequestedEventHandlerIUnknownRelease),
	},
	NewComProc(_ICoreWebView2PermissionRequestedEventHandlerInvoke),
}

func newICoreWebView2PermissionRequestedEventHandler(impl iCoreWebView2PermissionRequestedEventHandlerImpl) *iCoreWebView2PermissionRequestedEventHandler {
	return &iCoreWebView2PermissionRequestedEventHandler{
		vtbl: &iCoreWebView2PermissionRequestedEventHandlerFn,
		impl: impl,
	}
}

func (i *ICoreWebView2) AddWebResourceRequestedFilter(uri string, resourceContext COREWEBVIEW2_WEB_RESOURCE_CONTEXT) error {
	var err error
	// Convert string 'uri' to *uint16
	_uri, err := windows.UTF16PtrFromString(uri)
	if err != nil {
		return err
	}
	_, _, err = i.vtbl.AddWebResourceRequestedFilter.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_uri)),
		uintptr(resourceContext),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}
func (i *ICoreWebView2) AddNavigationCompleted(eventHandler *ICoreWebView2NavigationCompletedEventHandler, token *_EventRegistrationToken) error {
	var err error
	_, _, err = i.vtbl.AddNavigationCompleted.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2) OpenDevToolsWindow() error {
	var err error
	_, _, err = i.vtbl.OpenDevToolsWindow.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}
