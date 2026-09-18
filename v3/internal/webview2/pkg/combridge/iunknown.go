//go:build windows

package combridge

import (
	"golang.org/x/sys/windows"
)

const iUnknownGUID = "{00000000-0000-0000-C000-000000000046}"

func init() {
	registerVTableInternal[IUnknown, IUnknown](
		iUnknownGUID,
		true,
		iUnknownQueryInterface,
		iUnknownAddRef,
		iUnknownRelease,
	)
}

type IUnknown interface{}

func iUnknownQueryInterface(this uintptr, refiid *windows.GUID, ppvObject *uintptr) uintptr {
	if refiid == nil || ppvObject == nil {
		return uintptr(windows.E_INVALIDARG)
	}

	comIfcePointersL.RLock()
	obj := comIfcePointers[this]
	comIfcePointersL.RUnlock()

	ref := obj.queryInterface(refiid.String(), true)
	if ref != 0 {
		*ppvObject = ref
		return windows.NO_ERROR
	}

	*ppvObject = 0
	return uintptr(windows.E_NOINTERFACE)
}

func iUnknownAddRef(this uintptr) uintptr {
	comIfcePointersL.RLock()
	obj := comIfcePointers[this]
	comIfcePointersL.RUnlock()

	return uintptr(obj.addRef())
}

func iUnknownRelease(this uintptr) uintptr {
	comIfcePointersL.RLock()
	obj := comIfcePointers[this]
	comIfcePointersL.RUnlock()

	return uintptr(obj.release())
}
