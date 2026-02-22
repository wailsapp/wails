//go:build windows
// +build windows

package cfd

type comDlgFilterSpec struct {
	pszName *int16
	pszSpec *int16
}

type iUnknownVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
}

type iModalWindowVtbl struct {
	iUnknownVtbl
	Show uintptr // func (hwndOwner HWND) HRESULT
}

type iFileDialogVtbl struct {
	iModalWindowVtbl
	SetFileTypes        uintptr // func (cFileTypes UINT, rgFilterSpec *COMDLG_FILTERSPEC) HRESULT
	SetFileTypeIndex    uintptr // func(iFileType UINT) HRESULT
	GetFileTypeIndex    uintptr
	Advise              uintptr
	Unadvise            uintptr
	SetOptions          uintptr // func (fos FILEOPENDIALOGOPTIONS) HRESULT
	GetOptions          uintptr // func (pfos *FILEOPENDIALOGOPTIONS) HRESULT
	SetDefaultFolder    uintptr // func (psi *IShellItem) HRESULT
	SetFolder           uintptr // func (psi *IShellItem) HRESULT
	GetFolder           uintptr
	GetCurrentSelection uintptr
	SetFileName         uintptr // func (pszName LPCWSTR) HRESULT
	GetFileName         uintptr
	SetTitle            uintptr // func(pszTitle LPCWSTR) HRESULT
	SetOkButtonLabel    uintptr
	SetFileNameLabel    uintptr
	GetResult           uintptr // func (ppsi **IShellItem) HRESULT
	AddPlace            uintptr
	SetDefaultExtension uintptr // func (pszDefaultExtension LPCWSTR) HRESULT
	// This can only be used from a callback.
	Close           uintptr
	SetClientGuid   uintptr // func (guid REFGUID) HRESULT
	ClearClientData uintptr
	SetFilter       uintptr
}
