//go:build windows
// +build windows

package cfd

import (
	"fmt"
	"github.com/go-ole/go-ole"
	"strings"
	"syscall"
	"unsafe"
)

func hresultToError(hr uintptr) error {
	if hr < 0 {
		return ole.NewError(hr)
	}
	return nil
}

func (vtbl *iUnknownVtbl) release(objPtr unsafe.Pointer) error {
	ret, _, _ := syscall.Syscall(vtbl.Release,
		0,
		uintptr(objPtr),
		0,
		0)
	return hresultToError(ret)
}

func (vtbl *iModalWindowVtbl) show(objPtr unsafe.Pointer, hwnd uintptr) error {
	ret, _, _ := syscall.Syscall(vtbl.Show,
		1,
		uintptr(objPtr),
		hwnd,
		0)
	return hresultToError(ret)
}

func (vtbl *iFileDialogVtbl) setFileTypes(objPtr unsafe.Pointer, filters []FileFilter) error {
	cFileTypes := len(filters)
	if cFileTypes < 0 {
		return fmt.Errorf("must specify at least one filter")
	}
	comDlgFilterSpecs := make([]comDlgFilterSpec, cFileTypes)
	for i := 0; i < cFileTypes; i++ {
		filter := &filters[i]
		comDlgFilterSpecs[i] = comDlgFilterSpec{
			pszName: ole.SysAllocString(filter.DisplayName),
			pszSpec: ole.SysAllocString(filter.Pattern),
		}
	}
	ret, _, _ := syscall.Syscall(vtbl.SetFileTypes,
		2,
		uintptr(objPtr),
		uintptr(cFileTypes),
		uintptr(unsafe.Pointer(&comDlgFilterSpecs[0])))
	return hresultToError(ret)
}

// Options are:
// FOS_OVERWRITEPROMPT = 0x2,
// FOS_STRICTFILETYPES = 0x4,
// FOS_NOCHANGEDIR = 0x8,
// FOS_PICKFOLDERS = 0x20,
// FOS_FORCEFILESYSTEM = 0x40,
// FOS_ALLNONSTORAGEITEMS = 0x80,
// FOS_NOVALIDATE = 0x100,
// FOS_ALLOWMULTISELECT = 0x200,
// FOS_PATHMUSTEXIST = 0x800,
// FOS_FILEMUSTEXIST = 0x1000,
// FOS_CREATEPROMPT = 0x2000,
// FOS_SHAREAWARE = 0x4000,
// FOS_NOREADONLYRETURN = 0x8000,
// FOS_NOTESTFILECREATE = 0x10000,
// FOS_HIDEMRUPLACES = 0x20000,
// FOS_HIDEPINNEDPLACES = 0x40000,
// FOS_NODEREFERENCELINKS = 0x100000,
// FOS_OKBUTTONNEEDSINTERACTION = 0x200000,
// FOS_DONTADDTORECENT = 0x2000000,
// FOS_FORCESHOWHIDDEN = 0x10000000,
// FOS_DEFAULTNOMINIMODE = 0x20000000,
// FOS_FORCEPREVIEWPANEON = 0x40000000,
// FOS_SUPPORTSTREAMABLEITEMS = 0x80000000
func (vtbl *iFileDialogVtbl) setOptions(objPtr unsafe.Pointer, options uint32) error {
	ret, _, _ := syscall.Syscall(vtbl.SetOptions,
		1,
		uintptr(objPtr),
		uintptr(options),
		0)
	return hresultToError(ret)
}

func (vtbl *iFileDialogVtbl) getOptions(objPtr unsafe.Pointer) (uint32, error) {
	var options uint32
	ret, _, _ := syscall.Syscall(vtbl.GetOptions,
		1,
		uintptr(objPtr),
		uintptr(unsafe.Pointer(&options)),
		0)
	return options, hresultToError(ret)
}

func (vtbl *iFileDialogVtbl) addOption(objPtr unsafe.Pointer, option uint32) error {
	if options, err := vtbl.getOptions(objPtr); err == nil {
		return vtbl.setOptions(objPtr, options|option)
	} else {
		return err
	}
}

func (vtbl *iFileDialogVtbl) removeOption(objPtr unsafe.Pointer, option uint32) error {
	if options, err := vtbl.getOptions(objPtr); err == nil {
		return vtbl.setOptions(objPtr, options&^option)
	} else {
		return err
	}
}

func (vtbl *iFileDialogVtbl) setDefaultFolder(objPtr unsafe.Pointer, path string) error {
	shellItem, err := newIShellItem(path)
	if err != nil {
		return err
	}
	defer shellItem.vtbl.release(unsafe.Pointer(shellItem))
	ret, _, _ := syscall.Syscall(vtbl.SetDefaultFolder,
		1,
		uintptr(objPtr),
		uintptr(unsafe.Pointer(shellItem)),
		0)
	return hresultToError(ret)
}

func (vtbl *iFileDialogVtbl) setFolder(objPtr unsafe.Pointer, path string) error {
	shellItem, err := newIShellItem(path)
	if err != nil {
		return err
	}
	defer shellItem.vtbl.release(unsafe.Pointer(shellItem))
	ret, _, _ := syscall.Syscall(vtbl.SetFolder,
		1,
		uintptr(objPtr),
		uintptr(unsafe.Pointer(shellItem)),
		0)
	return hresultToError(ret)
}

func (vtbl *iFileDialogVtbl) setTitle(objPtr unsafe.Pointer, title string) error {
	titlePtr := ole.SysAllocString(title)
	ret, _, _ := syscall.Syscall(vtbl.SetTitle,
		1,
		uintptr(objPtr),
		uintptr(unsafe.Pointer(titlePtr)),
		0)
	return hresultToError(ret)
}

func (vtbl *iFileDialogVtbl) close(objPtr unsafe.Pointer) error {
	ret, _, _ := syscall.Syscall(vtbl.Close,
		1,
		uintptr(objPtr),
		0,
		0)
	return hresultToError(ret)
}

func (vtbl *iFileDialogVtbl) getResult(objPtr unsafe.Pointer) (*iShellItem, error) {
	var shellItem *iShellItem
	ret, _, _ := syscall.Syscall(vtbl.GetResult,
		1,
		uintptr(objPtr),
		uintptr(unsafe.Pointer(&shellItem)),
		0)
	return shellItem, hresultToError(ret)
}

func (vtbl *iFileDialogVtbl) getResultString(objPtr unsafe.Pointer) (string, error) {
	shellItem, err := vtbl.getResult(objPtr)
	if err != nil {
		return "", err
	}
	if shellItem == nil {
		return "", ErrCancelled
	}
	defer shellItem.vtbl.release(unsafe.Pointer(shellItem))
	return shellItem.vtbl.getDisplayName(unsafe.Pointer(shellItem))
}

func (vtbl *iFileDialogVtbl) setClientGuid(objPtr unsafe.Pointer, guid *ole.GUID) error {
	ret, _, _ := syscall.Syscall(vtbl.SetClientGuid,
		1,
		uintptr(objPtr),
		uintptr(unsafe.Pointer(guid)),
		0)
	return hresultToError(ret)
}

func (vtbl *iFileDialogVtbl) setDefaultExtension(objPtr unsafe.Pointer, defaultExtension string) error {
	if defaultExtension[0] == '.' {
		defaultExtension = strings.TrimPrefix(defaultExtension, ".")
	}
	defaultExtensionPtr := ole.SysAllocString(defaultExtension)
	ret, _, _ := syscall.Syscall(vtbl.SetDefaultExtension,
		1,
		uintptr(objPtr),
		uintptr(unsafe.Pointer(defaultExtensionPtr)),
		0)
	return hresultToError(ret)
}

func (vtbl *iFileDialogVtbl) setFileName(objPtr unsafe.Pointer, fileName string) error {
	fileNamePtr := ole.SysAllocString(fileName)
	ret, _, _ := syscall.Syscall(vtbl.SetFileName,
		1,
		uintptr(objPtr),
		uintptr(unsafe.Pointer(fileNamePtr)),
		0)
	return hresultToError(ret)
}

func (vtbl *iFileDialogVtbl) setSelectedFileFilterIndex(objPtr unsafe.Pointer, index uint) error {
	ret, _, _ := syscall.Syscall(vtbl.SetFileTypeIndex,
		1,
		uintptr(objPtr),
		uintptr(index+1), // SetFileTypeIndex counts from 1
		0)
	return hresultToError(ret)
}
