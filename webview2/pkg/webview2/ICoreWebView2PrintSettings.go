//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2PrintSettingsVtbl struct {
	IUnknownVtbl
	GetOrientation                ComProc
	PutOrientation                ComProc
	GetScaleFactor                ComProc
	PutScaleFactor                ComProc
	GetPageWidth                  ComProc
	PutPageWidth                  ComProc
	GetPageHeight                 ComProc
	PutPageHeight                 ComProc
	GetMarginTop                  ComProc
	PutMarginTop                  ComProc
	GetMarginBottom               ComProc
	PutMarginBottom               ComProc
	GetMarginLeft                 ComProc
	PutMarginLeft                 ComProc
	GetMarginRight                ComProc
	PutMarginRight                ComProc
	GetShouldPrintBackgrounds     ComProc
	PutShouldPrintBackgrounds     ComProc
	GetShouldPrintSelectionOnly   ComProc
	PutShouldPrintSelectionOnly   ComProc
	GetShouldPrintHeaderAndFooter ComProc
	PutShouldPrintHeaderAndFooter ComProc
	GetHeaderTitle                ComProc
	PutHeaderTitle                ComProc
	GetFooterUri                  ComProc
	PutFooterUri                  ComProc
}

type ICoreWebView2PrintSettings struct {
	Vtbl *ICoreWebView2PrintSettingsVtbl
}

func (i *ICoreWebView2PrintSettings) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2PrintSettings) GetOrientation() (COREWEBVIEW2_PRINT_ORIENTATION, error) {

	var orientation COREWEBVIEW2_PRINT_ORIENTATION

	hr, _, _ := i.Vtbl.GetOrientation.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&orientation)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return orientation, nil
}

func (i *ICoreWebView2PrintSettings) PutOrientation(orientation COREWEBVIEW2_PRINT_ORIENTATION) error {

	hr, _, _ := i.Vtbl.PutOrientation.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(orientation),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PrintSettings) GetScaleFactor() (float64, error) {

	var scaleFactor float64

	hr, _, _ := i.Vtbl.GetScaleFactor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&scaleFactor)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0.0, syscall.Errno(hr)
	}
	return scaleFactor, nil
}

func (i *ICoreWebView2PrintSettings) PutScaleFactor(scaleFactor float64) error {
	// The double parameter is passed BY VALUE; the per-arch appendDoubleArg
	// helpers pass it correctly for the target ABI (a pointer here reached
	// the callee as a garbage near-0.0 value).
	args, ok := appendDoubleArg([]uintptr{uintptr(unsafe.Pointer(i))}, scaleFactor)
	if !ok {
		// windows/arm64 cannot pass a by-value double (golang.org/issue/62583).
		return ErrDoubleArgUnsupported
	}
	hr, _, _ := i.Vtbl.PutScaleFactor.Call(args...)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PrintSettings) GetPageWidth() (float64, error) {

	var pageWidth float64

	hr, _, _ := i.Vtbl.GetPageWidth.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&pageWidth)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0.0, syscall.Errno(hr)
	}
	return pageWidth, nil
}

func (i *ICoreWebView2PrintSettings) PutPageWidth(pageWidth float64) error {
	// The double parameter is passed BY VALUE; the per-arch appendDoubleArg
	// helpers pass it correctly for the target ABI (a pointer here reached
	// the callee as a garbage near-0.0 value).
	args, ok := appendDoubleArg([]uintptr{uintptr(unsafe.Pointer(i))}, pageWidth)
	if !ok {
		// windows/arm64 cannot pass a by-value double (golang.org/issue/62583).
		return ErrDoubleArgUnsupported
	}
	hr, _, _ := i.Vtbl.PutPageWidth.Call(args...)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PrintSettings) GetPageHeight() (float64, error) {

	var pageHeight float64

	hr, _, _ := i.Vtbl.GetPageHeight.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&pageHeight)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0.0, syscall.Errno(hr)
	}
	return pageHeight, nil
}

func (i *ICoreWebView2PrintSettings) PutPageHeight(pageHeight float64) error {
	// The double parameter is passed BY VALUE; the per-arch appendDoubleArg
	// helpers pass it correctly for the target ABI (a pointer here reached
	// the callee as a garbage near-0.0 value).
	args, ok := appendDoubleArg([]uintptr{uintptr(unsafe.Pointer(i))}, pageHeight)
	if !ok {
		// windows/arm64 cannot pass a by-value double (golang.org/issue/62583).
		return ErrDoubleArgUnsupported
	}
	hr, _, _ := i.Vtbl.PutPageHeight.Call(args...)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PrintSettings) GetMarginTop() (float64, error) {

	var marginTop float64

	hr, _, _ := i.Vtbl.GetMarginTop.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&marginTop)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0.0, syscall.Errno(hr)
	}
	return marginTop, nil
}

func (i *ICoreWebView2PrintSettings) PutMarginTop(marginTop float64) error {
	// The double parameter is passed BY VALUE; the per-arch appendDoubleArg
	// helpers pass it correctly for the target ABI (a pointer here reached
	// the callee as a garbage near-0.0 value).
	args, ok := appendDoubleArg([]uintptr{uintptr(unsafe.Pointer(i))}, marginTop)
	if !ok {
		// windows/arm64 cannot pass a by-value double (golang.org/issue/62583).
		return ErrDoubleArgUnsupported
	}
	hr, _, _ := i.Vtbl.PutMarginTop.Call(args...)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PrintSettings) GetMarginBottom() (float64, error) {

	var marginBottom float64

	hr, _, _ := i.Vtbl.GetMarginBottom.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&marginBottom)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0.0, syscall.Errno(hr)
	}
	return marginBottom, nil
}

func (i *ICoreWebView2PrintSettings) PutMarginBottom(marginBottom float64) error {
	// The double parameter is passed BY VALUE; the per-arch appendDoubleArg
	// helpers pass it correctly for the target ABI (a pointer here reached
	// the callee as a garbage near-0.0 value).
	args, ok := appendDoubleArg([]uintptr{uintptr(unsafe.Pointer(i))}, marginBottom)
	if !ok {
		// windows/arm64 cannot pass a by-value double (golang.org/issue/62583).
		return ErrDoubleArgUnsupported
	}
	hr, _, _ := i.Vtbl.PutMarginBottom.Call(args...)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PrintSettings) GetMarginLeft() (float64, error) {

	var marginLeft float64

	hr, _, _ := i.Vtbl.GetMarginLeft.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&marginLeft)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0.0, syscall.Errno(hr)
	}
	return marginLeft, nil
}

func (i *ICoreWebView2PrintSettings) PutMarginLeft(marginLeft float64) error {
	// The double parameter is passed BY VALUE; the per-arch appendDoubleArg
	// helpers pass it correctly for the target ABI (a pointer here reached
	// the callee as a garbage near-0.0 value).
	args, ok := appendDoubleArg([]uintptr{uintptr(unsafe.Pointer(i))}, marginLeft)
	if !ok {
		// windows/arm64 cannot pass a by-value double (golang.org/issue/62583).
		return ErrDoubleArgUnsupported
	}
	hr, _, _ := i.Vtbl.PutMarginLeft.Call(args...)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PrintSettings) GetMarginRight() (float64, error) {

	var marginRight float64

	hr, _, _ := i.Vtbl.GetMarginRight.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&marginRight)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0.0, syscall.Errno(hr)
	}
	return marginRight, nil
}

func (i *ICoreWebView2PrintSettings) PutMarginRight(marginRight float64) error {
	// The double parameter is passed BY VALUE; the per-arch appendDoubleArg
	// helpers pass it correctly for the target ABI (a pointer here reached
	// the callee as a garbage near-0.0 value).
	args, ok := appendDoubleArg([]uintptr{uintptr(unsafe.Pointer(i))}, marginRight)
	if !ok {
		// windows/arm64 cannot pass a by-value double (golang.org/issue/62583).
		return ErrDoubleArgUnsupported
	}
	hr, _, _ := i.Vtbl.PutMarginRight.Call(args...)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PrintSettings) GetShouldPrintBackgrounds() (bool, error) {
	// Create int32 to hold bool result
	var _shouldPrintBackgrounds int32

	hr, _, _ := i.Vtbl.GetShouldPrintBackgrounds.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_shouldPrintBackgrounds)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	shouldPrintBackgrounds := _shouldPrintBackgrounds != 0
	return shouldPrintBackgrounds, nil
}

func (i *ICoreWebView2PrintSettings) PutShouldPrintBackgrounds(shouldPrintBackgrounds bool) error {
	// BOOL is a 4-byte by-value parameter: pass the value, not a pointer
	// (and not a 1-byte Go bool).
	var v int32
	if shouldPrintBackgrounds {
		v = 1
	}
	hr, _, _ := i.Vtbl.PutShouldPrintBackgrounds.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(v),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PrintSettings) GetShouldPrintSelectionOnly() (bool, error) {
	// Create int32 to hold bool result
	var _shouldPrintSelectionOnly int32

	hr, _, _ := i.Vtbl.GetShouldPrintSelectionOnly.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_shouldPrintSelectionOnly)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	shouldPrintSelectionOnly := _shouldPrintSelectionOnly != 0
	return shouldPrintSelectionOnly, nil
}

func (i *ICoreWebView2PrintSettings) PutShouldPrintSelectionOnly(shouldPrintSelectionOnly bool) error {
	// BOOL is a 4-byte by-value parameter: pass the value, not a pointer
	// (and not a 1-byte Go bool).
	var v int32
	if shouldPrintSelectionOnly {
		v = 1
	}
	hr, _, _ := i.Vtbl.PutShouldPrintSelectionOnly.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(v),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PrintSettings) GetShouldPrintHeaderAndFooter() (bool, error) {
	// Create int32 to hold bool result
	var _shouldPrintHeaderAndFooter int32

	hr, _, _ := i.Vtbl.GetShouldPrintHeaderAndFooter.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_shouldPrintHeaderAndFooter)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	shouldPrintHeaderAndFooter := _shouldPrintHeaderAndFooter != 0
	return shouldPrintHeaderAndFooter, nil
}

func (i *ICoreWebView2PrintSettings) PutShouldPrintHeaderAndFooter(shouldPrintHeaderAndFooter bool) error {
	// BOOL is a 4-byte by-value parameter: pass the value, not a pointer
	// (and not a 1-byte Go bool).
	var v int32
	if shouldPrintHeaderAndFooter {
		v = 1
	}
	hr, _, _ := i.Vtbl.PutShouldPrintHeaderAndFooter.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(v),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PrintSettings) GetHeaderTitle() (string, error) {
	// Create *uint16 to hold result
	var _headerTitle *uint16

	hr, _, _ := i.Vtbl.GetHeaderTitle.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_headerTitle)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	headerTitle := UTF16PtrToString(_headerTitle)
	CoTaskMemFree(unsafe.Pointer(_headerTitle))
	return headerTitle, nil
}

func (i *ICoreWebView2PrintSettings) PutHeaderTitle(headerTitle string) error {

	// Convert string 'headerTitle' to *uint16
	_headerTitle, err := UTF16PtrFromString(headerTitle)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PutHeaderTitle.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_headerTitle)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PrintSettings) GetFooterUri() (string, error) {
	// Create *uint16 to hold result
	var _footerUri *uint16

	hr, _, _ := i.Vtbl.GetFooterUri.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_footerUri)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	footerUri := UTF16PtrToString(_footerUri)
	CoTaskMemFree(unsafe.Pointer(_footerUri))
	return footerUri, nil
}

func (i *ICoreWebView2PrintSettings) PutFooterUri(footerUri string) error {

	// Convert string 'footerUri' to *uint16
	_footerUri, err := UTF16PtrFromString(footerUri)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PutFooterUri.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_footerUri)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
