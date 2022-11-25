package win32

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

const (
	NIF_MESSAGE             = 0x00000001
	NIF_ICON                = 0x00000002
	NIF_TIP                 = 0x00000004
	NIF_STATE               = 0x00000008
	NIF_INFO                = 0x00000010
	NIF_GUID                = 0x00000020
	NIF_REALTIME            = 0x00000040
	NIF_SHOWTIP             = 0x00000080
	NIM_ADD                 = 0x00000000
	NIM_MODIFY              = 0x00000001
	NIM_DELETE              = 0x00000002
	NIM_SETFOCUS            = 0x00000003
	NIM_SETVERSION          = 0x00000004
	NIS_HIDDEN              = 0x00000001
	NIS_SHAREDICON          = 0x00000002
	NIN_BALLOONSHOW         = 0x0402
	NIN_BALLOONTIMEOUT      = 0x0404
	NIN_BALLOONUSERCLICK    = 0x0405
	NIIF_NONE               = 0x00000000
	NIIF_INFO               = 0x00000001
	NIIF_WARNING            = 0x00000002
	NIIF_ERROR              = 0x00000003
	NIIF_USER               = 0x00000004
	NIIF_NOSOUND            = 0x00000010
	NIIF_LARGE_ICON         = 0x00000020
	NIIF_RESPECT_QUIET_TIME = 0x00000080
	NIIF_ICON_MASK          = 0x0000000F
)

type NOTIFYICONDATA struct {
	CbSize           uint32
	HWnd             uintptr
	UID              uint32
	UFlags           uint32
	UCallbackMessage uint32
	HIcon            uintptr
	SzTip            [128]uint16
	DwState          uint32
	DwStateMask      uint32
	SzInfo           [256]uint16
	UVersion         uint32
	SzInfoTitle      [64]uint16
	DwInfoFlags      uint32
	GUIDItem         windows.GUID
	HBalloonIcon     uintptr
}

func NotifyIcon(msg uint32, lpData *NOTIFYICONDATA) (int32, error) {
	r, _, err := procShellNotifyIcon.Call(
		uintptr(msg),
		uintptr(unsafe.Pointer(lpData)))
	if r == 0 {
		return 0, err
	}
	return int32(r), nil
}

func LookupIconIdFromDirectoryEx(presbits uintptr, isIcon bool, cxDesired int, cyDesired int, flags uint) (int32, error) {
	var icon uint32 = 0
	if isIcon {
		icon = 1
	}
	r, _, err := procLookupIconIdFromDirectoryEx.Call(
		presbits,
		uintptr(icon),
		uintptr(cxDesired),
		uintptr(cyDesired),
		uintptr(flags),
	)
	if r == 0 {
		return 0, err
	}
	return int32(r), nil
}

func CreateIconIndirect(data uintptr) (uintptr, error) {
	r, _, err := procCreateIconIndirect.Call(
		data,
	)

	if r == 0 {
		return 0, err
	}
	return r, nil
}

func CreateIconFromResourceEx(presbits uintptr, dwResSize uint32, isIcon bool, version uint32, cxDesired int, cyDesired int, flags uint) (uintptr, error) {
	icon := 0
	if isIcon {
		icon = 1
	}
	r, _, err := procCreateIconFromResourceEx.Call(
		presbits,
		uintptr(dwResSize),
		uintptr(icon),
		uintptr(version),
		uintptr(cxDesired),
		uintptr(cyDesired),
		uintptr(flags),
	)

	if r == 0 {
		return 0, err
	}
	return r, nil
}

func LoadImage(
	hInst uintptr,
	name *uint16,
	type_ uint32,
	cx, cy int32,
	fuLoad uint32) (uintptr, error) {
	r, _, err := procLoadImageW.Call(
		hInst,
		uintptr(unsafe.Pointer(name)),
		uintptr(type_),
		uintptr(cx),
		uintptr(cy),
		uintptr(fuLoad))
	if r == 0 {
		return 0, err
	}
	return r, nil
}
