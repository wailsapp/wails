//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2PrintSettings2Vtbl struct {
	IUnknownVtbl
	GetPageRanges   ComProc
	PutPageRanges   ComProc
	GetPagesPerSide ComProc
	PutPagesPerSide ComProc
	GetCopies       ComProc
	PutCopies       ComProc
	GetCollation    ComProc
	PutCollation    ComProc
	GetColorMode    ComProc
	PutColorMode    ComProc
	GetDuplex       ComProc
	PutDuplex       ComProc
	GetMediaSize    ComProc
	PutMediaSize    ComProc
	GetPrinterName  ComProc
	PutPrinterName  ComProc
}

type ICoreWebView2PrintSettings2 struct {
	Vtbl *ICoreWebView2PrintSettings2Vtbl
}

func (i *ICoreWebView2PrintSettings2) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2PrintSettings2() *ICoreWebView2PrintSettings2 {
	var result *ICoreWebView2PrintSettings2

	iidICoreWebView2PrintSettings2 := NewGUID("{CA7F0E1F-3484-41D1-8C1A-65CD44A63F8D}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2PrintSettings2)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2PrintSettings2) GetPageRanges() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, _ := i.Vtbl.GetPageRanges.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, nil
}

func (i *ICoreWebView2PrintSettings2) PutPageRanges(value string) error {

	// Convert string 'value' to *uint16
	_value, err := UTF16PtrFromString(value)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PutPageRanges.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PrintSettings2) GetPagesPerSide() (int32, error) {

	var value int32

	hr, _, _ := i.Vtbl.GetPagesPerSide.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2PrintSettings2) PutPagesPerSide(value int32) error {

	hr, _, _ := i.Vtbl.PutPagesPerSide.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PrintSettings2) GetCopies() (int32, error) {

	var value int32

	hr, _, _ := i.Vtbl.GetCopies.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2PrintSettings2) PutCopies(value int32) error {

	hr, _, _ := i.Vtbl.PutCopies.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PrintSettings2) GetCollation() (COREWEBVIEW2_PRINT_COLLATION, error) {

	var value COREWEBVIEW2_PRINT_COLLATION

	hr, _, _ := i.Vtbl.GetCollation.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2PrintSettings2) PutCollation(value COREWEBVIEW2_PRINT_COLLATION) error {

	hr, _, _ := i.Vtbl.PutCollation.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PrintSettings2) GetColorMode() (COREWEBVIEW2_PRINT_COLOR_MODE, error) {

	var value COREWEBVIEW2_PRINT_COLOR_MODE

	hr, _, _ := i.Vtbl.GetColorMode.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2PrintSettings2) PutColorMode(value COREWEBVIEW2_PRINT_COLOR_MODE) error {

	hr, _, _ := i.Vtbl.PutColorMode.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PrintSettings2) GetDuplex() (COREWEBVIEW2_PRINT_DUPLEX, error) {

	var value COREWEBVIEW2_PRINT_DUPLEX

	hr, _, _ := i.Vtbl.GetDuplex.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2PrintSettings2) PutDuplex(value COREWEBVIEW2_PRINT_DUPLEX) error {

	hr, _, _ := i.Vtbl.PutDuplex.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PrintSettings2) GetMediaSize() (COREWEBVIEW2_PRINT_MEDIA_SIZE, error) {

	var value COREWEBVIEW2_PRINT_MEDIA_SIZE

	hr, _, _ := i.Vtbl.GetMediaSize.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2PrintSettings2) PutMediaSize(value COREWEBVIEW2_PRINT_MEDIA_SIZE) error {

	hr, _, _ := i.Vtbl.PutMediaSize.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PrintSettings2) GetPrinterName() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, _ := i.Vtbl.GetPrinterName.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, nil
}

func (i *ICoreWebView2PrintSettings2) PutPrinterName(value string) error {

	// Convert string 'value' to *uint16
	_value, err := UTF16PtrFromString(value)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PutPrinterName.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
