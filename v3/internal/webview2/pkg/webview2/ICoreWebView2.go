//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Vtbl struct {
	IUnknownVtbl
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
	PostWebMessageAsJson                   ComProc
	PostWebMessageAsString                 ComProc
	AddWebMessageReceived                  ComProc
	RemoveWebMessageReceived               ComProc
	CallDevToolsProtocolMethod             ComProc
	GetBrowserProcessId                    ComProc
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
	Vtbl *ICoreWebView2Vtbl
}

func (i *ICoreWebView2) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetSettings() (*ICoreWebView2Settings, error) {

	var settings *ICoreWebView2Settings

	hr, _, _ := i.Vtbl.GetSettings.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&settings)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return settings, nil
}

func (i *ICoreWebView2) GetSource() (string, error) {
	// Create *uint16 to hold result
	var _uri *uint16

	hr, _, _ := i.Vtbl.GetSource.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_uri)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	uri := UTF16PtrToString(_uri)
	CoTaskMemFree(unsafe.Pointer(_uri))
	return uri, nil
}

func (i *ICoreWebView2) Navigate(uri string) error {

	// Convert string 'uri' to *uint16
	_uri, err := UTF16PtrFromString(uri)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.Navigate.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_uri)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) NavigateToString(htmlContent string) error {

	// Convert string 'htmlContent' to *uint16
	_htmlContent, err := UTF16PtrFromString(htmlContent)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.NavigateToString.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_htmlContent)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) AddNavigationStarting(eventHandler *ICoreWebView2NavigationStartingEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddNavigationStarting.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2) RemoveNavigationStarting(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveNavigationStarting.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) AddContentLoading(eventHandler *ICoreWebView2ContentLoadingEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddContentLoading.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2) RemoveContentLoading(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveContentLoading.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) AddSourceChanged(eventHandler *ICoreWebView2SourceChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddSourceChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2) RemoveSourceChanged(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveSourceChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) AddHistoryChanged(eventHandler *ICoreWebView2HistoryChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddHistoryChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2) RemoveHistoryChanged(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveHistoryChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) AddNavigationCompleted(eventHandler *ICoreWebView2NavigationCompletedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddNavigationCompleted.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2) RemoveNavigationCompleted(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveNavigationCompleted.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) AddFrameNavigationStarting(eventHandler *ICoreWebView2NavigationStartingEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddFrameNavigationStarting.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2) RemoveFrameNavigationStarting(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveFrameNavigationStarting.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) AddFrameNavigationCompleted(eventHandler *ICoreWebView2NavigationCompletedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddFrameNavigationCompleted.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2) RemoveFrameNavigationCompleted(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveFrameNavigationCompleted.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) AddScriptDialogOpening(eventHandler *ICoreWebView2ScriptDialogOpeningEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddScriptDialogOpening.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2) RemoveScriptDialogOpening(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveScriptDialogOpening.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) AddPermissionRequested(eventHandler *ICoreWebView2PermissionRequestedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddPermissionRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2) RemovePermissionRequested(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemovePermissionRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) AddProcessFailed(eventHandler *ICoreWebView2ProcessFailedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddProcessFailed.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2) RemoveProcessFailed(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveProcessFailed.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) AddScriptToExecuteOnDocumentCreated(javaScript string, handler *ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler) error {

	// Convert string 'javaScript' to *uint16
	_javaScript, err := UTF16PtrFromString(javaScript)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.AddScriptToExecuteOnDocumentCreated.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_javaScript)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) RemoveScriptToExecuteOnDocumentCreated(id string) error {

	// Convert string 'id' to *uint16
	_id, err := UTF16PtrFromString(id)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.RemoveScriptToExecuteOnDocumentCreated.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_id)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) ExecuteScript(javaScript string, handler *ICoreWebView2ExecuteScriptCompletedHandler) error {

	// Convert string 'javaScript' to *uint16
	_javaScript, err := UTF16PtrFromString(javaScript)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.ExecuteScript.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_javaScript)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) CapturePreview(imageFormat COREWEBVIEW2_CAPTURE_PREVIEW_IMAGE_FORMAT, imageStream *IStream, handler *ICoreWebView2CapturePreviewCompletedHandler) error {

	hr, _, _ := i.Vtbl.CapturePreview.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(imageFormat),
		uintptr(unsafe.Pointer(imageStream)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) Reload() error {

	hr, _, _ := i.Vtbl.Reload.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) PostWebMessageAsJson(webMessageAsJson string) error {

	// Convert string 'webMessageAsJson' to *uint16
	_webMessageAsJson, err := UTF16PtrFromString(webMessageAsJson)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PostWebMessageAsJson.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_webMessageAsJson)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) PostWebMessageAsString(webMessageAsString string) error {

	// Convert string 'webMessageAsString' to *uint16
	_webMessageAsString, err := UTF16PtrFromString(webMessageAsString)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PostWebMessageAsString.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_webMessageAsString)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) AddWebMessageReceived(handler *ICoreWebView2WebMessageReceivedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddWebMessageReceived.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(handler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2) RemoveWebMessageReceived(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveWebMessageReceived.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) CallDevToolsProtocolMethod(methodName string, parametersAsJson string, handler *ICoreWebView2CallDevToolsProtocolMethodCompletedHandler) error {

	// Convert string 'methodName' to *uint16
	_methodName, err := UTF16PtrFromString(methodName)
	if err != nil {
		return err
	}
	// Convert string 'parametersAsJson' to *uint16
	_parametersAsJson, err := UTF16PtrFromString(parametersAsJson)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.CallDevToolsProtocolMethod.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_methodName)),
		uintptr(unsafe.Pointer(_parametersAsJson)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) GetBrowserProcessId() (uint32, error) {

	var value uint32

	hr, _, _ := i.Vtbl.GetBrowserProcessId.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2) GetCanGoBack() (bool, error) {
	// Create int32 to hold bool result
	var _canGoBack int32

	hr, _, _ := i.Vtbl.GetCanGoBack.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_canGoBack)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	canGoBack := _canGoBack != 0
	return canGoBack, nil
}

func (i *ICoreWebView2) GetCanGoForward() (bool, error) {
	// Create int32 to hold bool result
	var _canGoForward int32

	hr, _, _ := i.Vtbl.GetCanGoForward.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_canGoForward)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	canGoForward := _canGoForward != 0
	return canGoForward, nil
}

func (i *ICoreWebView2) GoBack() error {

	hr, _, _ := i.Vtbl.GoBack.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) GoForward() error {

	hr, _, _ := i.Vtbl.GoForward.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) GetDevToolsProtocolEventReceiver(eventName string) (*ICoreWebView2DevToolsProtocolEventReceiver, error) {

	// Convert string 'eventName' to *uint16
	_eventName, err := UTF16PtrFromString(eventName)
	if err != nil {
		return nil, err
	}
	var receiver *ICoreWebView2DevToolsProtocolEventReceiver

	hr, _, _ := i.Vtbl.GetDevToolsProtocolEventReceiver.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_eventName)),
		uintptr(unsafe.Pointer(&receiver)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return receiver, nil
}

func (i *ICoreWebView2) Stop() error {

	hr, _, _ := i.Vtbl.Stop.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) AddNewWindowRequested(eventHandler *ICoreWebView2NewWindowRequestedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddNewWindowRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2) RemoveNewWindowRequested(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveNewWindowRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) AddDocumentTitleChanged(eventHandler *ICoreWebView2DocumentTitleChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddDocumentTitleChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2) RemoveDocumentTitleChanged(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveDocumentTitleChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) GetDocumentTitle() (string, error) {
	// Create *uint16 to hold result
	var _title *uint16

	hr, _, _ := i.Vtbl.GetDocumentTitle.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_title)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	title := UTF16PtrToString(_title)
	CoTaskMemFree(unsafe.Pointer(_title))
	return title, nil
}

func (i *ICoreWebView2) AddHostObjectToScript(name string, object *VARIANT) error {

	// Convert string 'name' to *uint16
	_name, err := UTF16PtrFromString(name)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.AddHostObjectToScript.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_name)),
		uintptr(unsafe.Pointer(object)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) RemoveHostObjectFromScript(name string) error {

	// Convert string 'name' to *uint16
	_name, err := UTF16PtrFromString(name)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.RemoveHostObjectFromScript.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_name)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) OpenDevToolsWindow() error {

	hr, _, _ := i.Vtbl.OpenDevToolsWindow.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) AddContainsFullScreenElementChanged(eventHandler *ICoreWebView2ContainsFullScreenElementChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddContainsFullScreenElementChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2) RemoveContainsFullScreenElementChanged(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveContainsFullScreenElementChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) GetContainsFullScreenElement() (bool, error) {
	// Create int32 to hold bool result
	var _containsFullScreenElement int32

	hr, _, _ := i.Vtbl.GetContainsFullScreenElement.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_containsFullScreenElement)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	containsFullScreenElement := _containsFullScreenElement != 0
	return containsFullScreenElement, nil
}

func (i *ICoreWebView2) AddWebResourceRequested(eventHandler *ICoreWebView2WebResourceRequestedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddWebResourceRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2) RemoveWebResourceRequested(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveWebResourceRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) AddWebResourceRequestedFilter(uri string, resourceContext COREWEBVIEW2_WEB_RESOURCE_CONTEXT) error {

	// Convert string 'uri' to *uint16
	_uri, err := UTF16PtrFromString(uri)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.AddWebResourceRequestedFilter.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_uri)),
		uintptr(resourceContext),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) RemoveWebResourceRequestedFilter(uri string, resourceContext COREWEBVIEW2_WEB_RESOURCE_CONTEXT) error {

	// Convert string 'uri' to *uint16
	_uri, err := UTF16PtrFromString(uri)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.RemoveWebResourceRequestedFilter.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_uri)),
		uintptr(resourceContext),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2) AddWindowCloseRequested(eventHandler *ICoreWebView2WindowCloseRequestedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddWindowCloseRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2) RemoveWindowCloseRequested(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveWindowCloseRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
