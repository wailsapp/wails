//go:build windows
// +build windows

package edge

import (
	"fmt"
	"log"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/webview2/internal/w32"

	"golang.org/x/sys/windows"
)

func init() {
	runtime.LockOSThread()

	r, _, _ := w32.Ole32CoInitializeEx.Call(0, uintptr(w32.COINIT_APARTMENTTHREADED))
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

func (i *ICoreWebView2) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2) AddContainsFullScreenElementChanged(eventHandler *ICoreWebView2ContainsFullScreenElementChangedEventHandler, token *_EventRegistrationToken) error {
	hr, _, _ := i.vtbl.AddContainsFullScreenElementChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}

	return nil
}

func (i *ICoreWebView2) AddNavigationCompleted(eventHandler *ICoreWebView2NavigationCompletedEventHandler, token *_EventRegistrationToken) error {
	hr, _, _ := i.vtbl.AddNavigationCompleted.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}

	return nil
}

func (i *ICoreWebView2) AddPermissionRequested(handler *iCoreWebView2PermissionRequestedEventHandler, token *_EventRegistrationToken) error {
	hr, _, _ := i.vtbl.AddPermissionRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(handler)),
		uintptr(unsafe.Pointer(token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}

	return nil
}

func (i *ICoreWebView2) AddProcessFailed(eventHandler *ICoreWebView2ProcessFailedEventHandler, token *_EventRegistrationToken) error {
	hr, _, _ := i.vtbl.AddProcessFailed.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}

	return nil
}

func (i *ICoreWebView2) AddWebMessageReceived(handler *iCoreWebView2WebMessageReceivedEventHandler, token *_EventRegistrationToken) error {
	hr, _, _ := i.vtbl.AddWebMessageReceived.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(handler)),
		uintptr(unsafe.Pointer(token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}

	return nil
}

func (i *ICoreWebView2) AddWebResourceRequested(handler *iCoreWebView2WebResourceRequestedEventHandler, token *_EventRegistrationToken) error {
	hr, _, _ := i.vtbl.AddWebResourceRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(handler)),
		uintptr(unsafe.Pointer(token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}

	return nil
}

func (i *ICoreWebView2) AddScriptToExecuteOnDocumentCreated(javaScript string, handler *iCoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler) error {
	u16js, err := windows.UTF16PtrFromString(javaScript)
	if err != nil {
		return err
	}

	hr, _, _ := i.vtbl.AddScriptToExecuteOnDocumentCreated.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(u16js)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}

	return nil
}

func (i *ICoreWebView2) ExecuteScript(javascript string, handler *iCoreWebView2ExecuteScriptCompletedHandler) error {
	u16js, err := windows.UTF16PtrFromString(javascript)
	if err != nil {
		return err
	}

	hr, _, _ := i.vtbl.ExecuteScript.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(u16js)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}

	return nil
}

func (i *ICoreWebView2) GetSettings() (*ICoreWebViewSettings, error) {

	var settings *ICoreWebViewSettings
	hr, _, _ := i.vtbl.GetSettings.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&settings)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, windows.Errno(hr)
	}
	return settings, nil
}

func (i *ICoreWebView2) GetSource() (string, error) {
	var _source *uint16
	hr, _, _ := i.vtbl.GetSource.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_source)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", windows.Errno(hr)
	}
	source := windows.UTF16PtrToString(_source)
	windows.CoTaskMemFree(unsafe.Pointer(_source))
	return source, nil
}

func (i *ICoreWebView2) GetContainsFullScreenElement() (bool, error) {
	var result bool
	hr, _, _ := i.vtbl.GetContainsFullScreenElement.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&result)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, windows.Errno(hr)
	}
	return result, nil
}

func (i *ICoreWebView2) Navigate(url string) error {
	u16url, err := windows.UTF16PtrFromString(url)
	if err != nil {
		return err
	}

	hr, _, _ := i.vtbl.Navigate.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(u16url)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}

	return nil
}

func (i *ICoreWebView2) NavigateToString(htmlContent string) error {
	u16Html, err := windows.UTF16PtrFromString(htmlContent)
	if err != nil {
		return err
	}

	hr, _, _ := i.vtbl.NavigateToString.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(u16Html)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}

	return nil
}

func (i *ICoreWebView2) PostWebMessageAsString(webMessageAsString string) error {
	u16msg, err := windows.UTF16PtrFromString(webMessageAsString)
	if err != nil {
		return err
	}

	hr, _, _ := i.vtbl.PostWebMessageAsString.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(u16msg)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}

	return nil
}

func (i *ICoreWebView2) QueryInterface2() (*ICoreWebView2_2, error) {
	var result *ICoreWebView2_2
	iid := windows.GUID{
		Data1: 0x9E8F0CF8,
		Data2: 0xE670,
		Data3: 0x4B5E,
		Data4: [8]byte{0xB2, 0xBC, 0x73, 0xE0, 0x61, 0xE3, 0x18, 0x4C},
	}
	hr, _, _ := i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&iid)),
		uintptr(unsafe.Pointer(&result)))
	if hr != 0 {
		return nil, windows.Errno(hr)
	}
	return result, nil
}

// ICoreWebView2EnvironmentVtbl
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

// CreateCoreWebView2Controller asynchronously creates a new WebView.
func (e *ICoreWebView2Environment) CreateCoreWebView2Controller(parentWindow uintptr, handler *iCoreWebView2CreateCoreWebView2ControllerCompletedHandler) error {
	hr, _, _ := e.vtbl.CreateCoreWebView2Controller.Call(
		uintptr(unsafe.Pointer(e)),
		parentWindow,
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}

	return nil
}

// CreateWebResourceResponse creates a new ICoreWebView2WebResourceResponse, it must be released after finishing using it.
func (e *ICoreWebView2Environment) CreateWebResourceResponse(content []byte, statusCode int, reasonPhrase string, headers string) (*ICoreWebView2WebResourceResponse, error) {

	var stream uintptr

	if len(content) > 0 {
		// Create stream for response
		stream, err := w32.SHCreateMemStream(content)
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
	hr, _, _ := e.vtbl.CreateWebResourceResponse.Call(
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
		return nil, fmt.Errorf("unknown error")
	}

	return response, nil
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

func (i *iCoreWebView2PermissionRequestedEventArgs) GetPermissionKind() (CoreWebView2PermissionKind, error) {
	var kind CoreWebView2PermissionKind

	hr, _, _ := i.vtbl.GetPermissionKind.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&kind)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, windows.Errno(hr)
	}

	return kind, nil
}

func (i *iCoreWebView2PermissionRequestedEventArgs) PutState(state CoreWebView2PermissionState) error {
	hr, _, _ := i.vtbl.PutState.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(state),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}

	return nil
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

	// Convert string 'uri' to *uint16
	_uri, err := windows.UTF16PtrFromString(uri)
	if err != nil {
		return err
	}
	hr, _, _ := i.vtbl.AddWebResourceRequestedFilter.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_uri)),
		uintptr(resourceContext),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) OpenDevToolsWindow() error {

	hr, _, _ := i.vtbl.OpenDevToolsWindow.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}
