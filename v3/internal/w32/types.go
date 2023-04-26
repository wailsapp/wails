// Copyright 2010-2012 The W32 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package w32

import (
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"
)

// From MSDN: Windows Data Types
// http://msdn.microsoft.com/en-us/library/s3f49ktz.aspx
// http://msdn.microsoft.com/en-us/library/windows/desktop/aa383751.aspx
type (
	ATOM            uint16
	BOOL            int32
	COLORREF        uint32
	DWM_FRAME_COUNT uint64
	DWORD           uint32
	HACCEL          HANDLE
	HANDLE          uintptr
	HBITMAP         HANDLE
	HBRUSH          HANDLE
	HCURSOR         HANDLE
	HDC             HANDLE
	HDROP           HANDLE
	HDWP            HANDLE
	HENHMETAFILE    HANDLE
	HFONT           HANDLE
	HGDIOBJ         HANDLE
	HGLOBAL         HANDLE
	HGLRC           HANDLE
	HHOOK           HANDLE
	HICON           HANDLE
	HIMAGELIST      HANDLE
	HINSTANCE       HANDLE
	HKEY            HANDLE
	HKL             HANDLE
	HMENU           HANDLE
	HMODULE         HANDLE
	HMONITOR        HANDLE
	HPEN            HANDLE
	HRESULT         int32
	HRGN            HANDLE
	HRSRC           HANDLE
	HTHUMBNAIL      HANDLE
	HWND            HANDLE
	LPARAM          uintptr
	LPCVOID         unsafe.Pointer
	LRESULT         uintptr
	PVOID           unsafe.Pointer
	QPC_TIME        uint64
	ULONG_PTR       uintptr
	WPARAM          uintptr
	HRAWINPUT       HANDLE
)

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd162805.aspx
type POINT struct {
	X, Y int32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd162897.aspx
type RECT struct {
	Left, Top, Right, Bottom int32
}

func (r *RECT) Width() int32 {
	return r.Right - r.Left
}

func (r *RECT) Height() int32 {
	return r.Bottom - r.Top
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms633577.aspx
type WNDCLASSEX struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   HINSTANCE
	Icon       HICON
	Cursor     HCURSOR
	Background HBRUSH
	MenuName   *uint16
	ClassName  *uint16
	IconSm     HICON
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms644958.aspx
type MSG struct {
	Hwnd    HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd145037.aspx
type LOGFONT struct {
	Height         int32
	Width          int32
	Escapement     int32
	Orientation    int32
	Weight         int32
	Italic         byte
	Underline      byte
	StrikeOut      byte
	CharSet        byte
	OutPrecision   byte
	ClipPrecision  byte
	Quality        byte
	PitchAndFamily byte
	FaceName       [LF_FACESIZE]uint16
}

func toString(s []uint16) string {
	for i, c := range s {
		if c == 0 {
			return string(utf16.Decode(s[:i]))
		}
	}
	return string(utf16.Decode(s))
}

func (f *LOGFONT) GetFaceName() string {
	return toString(f.FaceName[:])
}

func (f *LOGFONT) SetFaceName(name string) {
	s := utf16.Encode([]rune(name))
	max := len(f.FaceName) - 1
	if len(s) > max {
		s = s[:max]
	}
	copy(f.FaceName[:], s)
	f.FaceName[len(s)] = 0
}

type ENUMLOGFONTEX struct {
	LOGFONT
	FullName [LF_FULLFACESIZE]uint16
	Style    [LF_FACESIZE]uint16
	Script   [LF_FACESIZE]uint16
}

func (f *ENUMLOGFONTEX) GetFullName() string {
	return toString(f.FullName[:])
}

func (f *ENUMLOGFONTEX) GetStyle() string {
	return toString(f.Style[:])
}

func (f *ENUMLOGFONTEX) GetScript() string {
	return toString(f.Script[:])
}

type ENUMTEXTMETRIC struct {
	NEWTEXTMETRICEX
	AXESLIST
}

type NEWTEXTMETRICEX struct {
	NEWTEXTMETRIC
	FONTSIGNATURE
}

type NEWTEXTMETRIC struct {
	Height           int32
	Ascent           int32
	Descent          int32
	InternalLeading  int32
	ExternalLeading  int32
	AveCharWidth     int32
	MaxCharWidth     int32
	Weight           int32
	Overhang         int32
	DigitizedAspectX int32
	DigitizedAspectY int32
	FirstChar        uint16
	LastChar         uint16
	DefaultChar      uint16
	BreakChar        uint16
	Italic           byte
	Underlined       byte
	StruckOut        byte
	PitchAndFamily   byte
	CharSet          byte
	Flags            uint32
	SizeEM           uint32
	CellHeight       uint32
	AvgWidth         uint32
}

type FONTSIGNATURE struct {
	Usb [4]uint32
	Csb [2]uint32
}

type AXESLIST struct {
	Reserved uint32
	NumAxes  uint32
	AxisInfo [MM_MAX_NUMAXES]AXISINFO
}

type AXISINFO struct {
	MinValue int32
	MaxValue int32
	AxisName [MM_MAX_AXES_NAMELEN]uint16
}

func (i *AXISINFO) GetAxisName() string {
	return toString(i.AxisName[:])
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms646839.aspx
type OPENFILENAME struct {
	StructSize      uint32
	Owner           HWND
	Instance        HINSTANCE
	Filter          *uint16
	CustomFilter    *uint16
	MaxCustomFilter uint32
	FilterIndex     uint32
	File            *uint16
	MaxFile         uint32
	FileTitle       *uint16
	MaxFileTitle    uint32
	InitialDir      *uint16
	Title           *uint16
	Flags           uint32
	FileOffset      uint16
	FileExtension   uint16
	DefExt          *uint16
	CustData        uintptr
	FnHook          uintptr
	TemplateName    *uint16
	PvReserved      unsafe.Pointer
	DwReserved      uint32
	FlagsEx         uint32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/bb773205.aspx
type BROWSEINFO struct {
	Owner        HWND
	Root         *uint16
	DisplayName  *uint16
	Title        *uint16
	Flags        uint32
	CallbackFunc uintptr
	LParam       uintptr
	Image        int32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/aa373931.aspx
type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms221627.aspx
type VARIANT struct {
	VT         uint16 //  2
	WReserved1 uint16 //  4
	WReserved2 uint16 //  6
	WReserved3 uint16 //  8
	Val        int64  // 16
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms221416.aspx
type DISPPARAMS struct {
	Rgvarg            uintptr
	RgdispidNamedArgs uintptr
	CArgs             uint32
	CNamedArgs        uint32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms221133.aspx
type EXCEPINFO struct {
	WCode             uint16
	WReserved         uint16
	BstrSource        *uint16
	BstrDescription   *uint16
	BstrHelpFile      *uint16
	DwHelpContext     uint32
	PvReserved        uintptr
	PfnDeferredFillIn uintptr
	Scode             int32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd145035.aspx
type LOGBRUSH struct {
	LbStyle uint32
	LbColor COLORREF
	LbHatch uintptr
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd183565.aspx
type DEVMODE struct {
	DmDeviceName       [CCHDEVICENAME]uint16
	DmSpecVersion      uint16
	DmDriverVersion    uint16
	DmSize             uint16
	DmDriverExtra      uint16
	DmFields           uint32
	DmOrientation      int16
	DmPaperSize        int16
	DmPaperLength      int16
	DmPaperWidth       int16
	DmScale            int16
	DmCopies           int16
	DmDefaultSource    int16
	DmPrintQuality     int16
	DmColor            int16
	DmDuplex           int16
	DmYResolution      int16
	DmTTOption         int16
	DmCollate          int16
	DmFormName         [CCHFORMNAME]uint16
	DmLogPixels        uint16
	DmBitsPerPel       uint32
	DmPelsWidth        uint32
	DmPelsHeight       uint32
	DmDisplayFlags     uint32
	DmDisplayFrequency uint32
	DmICMMethod        uint32
	DmICMIntent        uint32
	DmMediaType        uint32
	DmDitherType       uint32
	DmReserved1        uint32
	DmReserved2        uint32
	DmPanningWidth     uint32
	DmPanningHeight    uint32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd183376.aspx
type BITMAPINFOHEADER struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd162938.aspx
type RGBQUAD struct {
	RgbBlue     byte
	RgbGreen    byte
	RgbRed      byte
	RgbReserved byte
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd183375.aspx
type BITMAPINFO struct {
	BmiHeader BITMAPINFOHEADER
	BmiColors *RGBQUAD
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd183371.aspx
type BITMAP struct {
	BmType       int32
	BmWidth      int32
	BmHeight     int32
	BmWidthBytes int32
	BmPlanes     uint16
	BmBitsPixel  uint16
	BmBits       unsafe.Pointer
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/dd183374(v=vs.85).aspx
type BITMAPFILEHEADER struct {
	BfType      uint16
	BfSize      uint32
	BfReserved1 uint16
	BfReserved2 uint16
	BfOffBits   uint32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd183567.aspx
type DIBSECTION struct {
	DsBm        BITMAP
	DsBmih      BITMAPINFOHEADER
	DsBitfields [3]uint32
	DshSection  HANDLE
	DsOffset    uint32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd162607.aspx
type ENHMETAHEADER struct {
	IType          uint32
	NSize          uint32
	RclBounds      RECT
	RclFrame       RECT
	DSignature     uint32
	NVersion       uint32
	NBytes         uint32
	NRecords       uint32
	NHandles       uint16
	SReserved      uint16
	NDescription   uint32
	OffDescription uint32
	NPalEntries    uint32
	SzlDevice      SIZE
	SzlMillimeters SIZE
	CbPixelFormat  uint32
	OffPixelFormat uint32
	BOpenGL        uint32
	SzlMicrometers SIZE
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd145106.aspx
type SIZE struct {
	CX, CY int32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd145132.aspx
type TEXTMETRIC struct {
	TmHeight           int32
	TmAscent           int32
	TmDescent          int32
	TmInternalLeading  int32
	TmExternalLeading  int32
	TmAveCharWidth     int32
	TmMaxCharWidth     int32
	TmWeight           int32
	TmOverhang         int32
	TmDigitizedAspectX int32
	TmDigitizedAspectY int32
	TmFirstChar        uint16
	TmLastChar         uint16
	TmDefaultChar      uint16
	TmBreakChar        uint16
	TmItalic           byte
	TmUnderlined       byte
	TmStruckOut        byte
	TmPitchAndFamily   byte
	TmCharSet          byte
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd183574.aspx
type DOCINFO struct {
	CbSize       int32
	LpszDocName  *uint16
	LpszOutput   *uint16
	LpszDatatype *uint16
	FwType       uint32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/bb775514.aspx
type NMHDR struct {
	HwndFrom HWND
	IdFrom   uintptr
	Code     uint32
}

type NMUPDOWN struct {
	Hdr   NMHDR
	Pos   int32
	Delta int32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/bb774743.aspx
type LVCOLUMN struct {
	Mask       uint32
	Fmt        int32
	Cx         int32
	PszText    *uint16
	CchTextMax int32
	ISubItem   int32
	IImage     int32
	IOrder     int32
	CxMin      int32
	CxDefault  int32
	CxIdeal    int32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/bb774760.aspx
type LVITEM struct {
	Mask       uint32
	IItem      int32
	ISubItem   int32
	State      uint32
	StateMask  uint32
	PszText    *uint16
	CchTextMax int32
	IImage     int32
	LParam     uintptr
	IIndent    int32
	IGroupId   int32
	CColumns   uint32
	PuColumns  uint32
	PiColFmt   *int32
	IGroup     int32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/bb774754.aspx
type LVHITTESTINFO struct {
	Pt       POINT
	Flags    uint32
	IItem    int32
	ISubItem int32
	IGroup   int32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/bb774771.aspx
type NMITEMACTIVATE struct {
	Hdr       NMHDR
	IItem     int32
	ISubItem  int32
	UNewState uint32
	UOldState uint32
	UChanged  uint32
	PtAction  POINT
	LParam    uintptr
	UKeyFlags uint32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/bb774773.aspx
type NMLISTVIEW struct {
	Hdr       NMHDR
	IItem     int32
	ISubItem  int32
	UNewState uint32
	UOldState uint32
	UChanged  uint32
	PtAction  POINT
	LParam    uintptr
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/bb774780.aspx
type NMLVDISPINFO struct {
	Hdr  NMHDR
	Item LVITEM
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/bb775507.aspx
type INITCOMMONCONTROLSEX struct {
	size uint32
	ICC  uint32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/bb760256.aspx
type TOOLINFO struct {
	CbSize     uint32
	UFlags     uint32
	Hwnd       HWND
	UId        uintptr
	Rect       RECT
	Hinst      HINSTANCE
	LpszText   *uint16
	LParam     uintptr
	LpReserved unsafe.Pointer
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms645604.aspx
type TRACKMOUSEEVENT struct {
	CbSize      uint32
	DwFlags     uint32
	HwndTrack   HWND
	DwHoverTime uint32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms534067.aspx
type GdiplusStartupInput struct {
	GdiplusVersion           uint32
	DebugEventCallback       uintptr
	SuppressBackgroundThread BOOL
	SuppressExternalCodecs   BOOL
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms534068.aspx
type GdiplusStartupOutput struct {
	NotificationHook   uintptr
	NotificationUnhook uintptr
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd162768.aspx
type PAINTSTRUCT struct {
	Hdc         HDC
	FErase      BOOL
	RcPaint     RECT
	FRestore    BOOL
	FIncUpdate  BOOL
	RgbReserved [32]byte
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/aa363646.aspx
type EVENTLOGRECORD struct {
	Length              uint32
	Reserved            uint32
	RecordNumber        uint32
	TimeGenerated       uint32
	TimeWritten         uint32
	EventID             uint32
	EventType           uint16
	NumStrings          uint16
	EventCategory       uint16
	ReservedFlags       uint16
	ClosingRecordNumber uint32
	StringOffset        uint32
	UserSidLength       uint32
	UserSidOffset       uint32
	DataLength          uint32
	DataOffset          uint32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms685996.aspx
type SERVICE_STATUS struct {
	DwServiceType             uint32
	DwCurrentState            uint32
	DwControlsAccepted        uint32
	DwWin32ExitCode           uint32
	DwServiceSpecificExitCode uint32
	DwCheckPoint              uint32
	DwWaitHint                uint32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms684225.aspx
type MODULEENTRY32 struct {
	Size         uint32
	ModuleID     uint32
	ProcessID    uint32
	GlblcntUsage uint32
	ProccntUsage uint32
	ModBaseAddr  *uint8
	ModBaseSize  uint32
	HModule      HMODULE
	SzModule     [MAX_MODULE_NAME32 + 1]uint16
	SzExePath    [MAX_PATH]uint16
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms724284.aspx
type FILETIME struct {
	DwLowDateTime  uint32
	DwHighDateTime uint32
}

func (t FILETIME) Uint64() uint64 {
	return uint64(t.DwHighDateTime)<<32 | uint64(t.DwLowDateTime)
}

func (t FILETIME) Time() time.Time {
	// https://docs.microsoft.com/en-us/windows/win32/api/minwinbase/ns-minwinbase-filetime
	ref := time.Date(1601, time.January, 1, 0, 0, 0, 0, time.UTC)
	const tick = 100 * time.Nanosecond
	// The FILETIME is a uint64 of 100-nanosecond intervals since 1601.
	// Unfortunately time.Duration is really an int64 so if we cast our uint64
	// to a time.Duration it becomes negative. Thus we do it in 2 steps, adding
	// half the time each step to avoid overflow.
	return ref.
		Add(time.Duration(t.Uint64()) * (tick / 2)).
		Add(time.Duration(t.Uint64()) * (tick / 2))
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms682119.aspx
type COORD struct {
	X, Y int16
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms686311.aspx
type SMALL_RECT struct {
	Left, Top, Right, Bottom int16
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms682093.aspx
type CONSOLE_SCREEN_BUFFER_INFO struct {
	DwSize              COORD
	DwCursorPosition    COORD
	WAttributes         uint16
	SrWindow            SMALL_RECT
	DwMaximumWindowSize COORD
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/bb773244.aspx
type MARGINS struct {
	CxLeftWidth, CxRightWidth, CyTopHeight, CyBottomHeight int32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/aa969500.aspx
type DWM_BLURBEHIND struct {
	DwFlags                uint32
	fEnable                BOOL
	hRgnBlur               HRGN
	fTransitionOnMaximized BOOL
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/aa969501.aspx
type DWM_PRESENT_PARAMETERS struct {
	cbSize             uint32
	fQueue             BOOL
	cRefreshStart      DWM_FRAME_COUNT
	cBuffer            uint32
	fUseSourceRate     BOOL
	rateSource         UNSIGNED_RATIO
	cRefreshesPerFrame uint32
	eSampling          DWM_SOURCE_FRAME_SAMPLING
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/aa969502.aspx
type DWM_THUMBNAIL_PROPERTIES struct {
	dwFlags               uint32
	rcDestination         RECT
	rcSource              RECT
	opacity               byte
	fVisible              BOOL
	fSourceClientAreaOnly BOOL
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/aa969503.aspx
type DWM_TIMING_INFO struct {
	cbSize                 uint32
	rateRefresh            UNSIGNED_RATIO
	qpcRefreshPeriod       QPC_TIME
	rateCompose            UNSIGNED_RATIO
	qpcVBlank              QPC_TIME
	cRefresh               DWM_FRAME_COUNT
	cDXRefresh             uint32
	qpcCompose             QPC_TIME
	cFrame                 DWM_FRAME_COUNT
	cDXPresent             uint32
	cRefreshFrame          DWM_FRAME_COUNT
	cFrameSubmitted        DWM_FRAME_COUNT
	cDXPresentSubmitted    uint32
	cFrameConfirmed        DWM_FRAME_COUNT
	cDXPresentConfirmed    uint32
	cRefreshConfirmed      DWM_FRAME_COUNT
	cDXRefreshConfirmed    uint32
	cFramesLate            DWM_FRAME_COUNT
	cFramesOutstanding     uint32
	cFrameDisplayed        DWM_FRAME_COUNT
	qpcFrameDisplayed      QPC_TIME
	cRefreshFrameDisplayed DWM_FRAME_COUNT
	cFrameComplete         DWM_FRAME_COUNT
	qpcFrameComplete       QPC_TIME
	cFramePending          DWM_FRAME_COUNT
	qpcFramePending        QPC_TIME
	cFramesDisplayed       DWM_FRAME_COUNT
	cFramesComplete        DWM_FRAME_COUNT
	cFramesPending         DWM_FRAME_COUNT
	cFramesAvailable       DWM_FRAME_COUNT
	cFramesDropped         DWM_FRAME_COUNT
	cFramesMissed          DWM_FRAME_COUNT
	cRefreshNextDisplayed  DWM_FRAME_COUNT
	cRefreshNextPresented  DWM_FRAME_COUNT
	cRefreshesDisplayed    DWM_FRAME_COUNT
	cRefreshesPresented    DWM_FRAME_COUNT
	cRefreshStarted        DWM_FRAME_COUNT
	cPixelsReceived        uint64
	cPixelsDrawn           uint64
	cBuffersEmpty          DWM_FRAME_COUNT
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd389402.aspx
type MilMatrix3x2D struct {
	S_11, S_12, S_21, S_22 float64
	DX, DY                 float64
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/aa969505.aspx
type UNSIGNED_RATIO struct {
	uiNumerator   uint32
	uiDenominator uint32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms632603.aspx
type CREATESTRUCT struct {
	CreateParams uintptr
	Instance     HINSTANCE
	Menu         HMENU
	Parent       HWND
	Cy, Cx       int32
	Y, X         int32
	Style        int32
	Name         *uint16
	Class        *uint16
	dwExStyle    uint32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd145065.aspx
type MONITORINFO struct {
	CbSize    uint32
	RcMonitor RECT
	RcWork    RECT
	DwFlags   uint32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd145066.aspx
type MONITORINFOEX struct {
	MONITORINFO
	SzDevice [CCHDEVICENAME]uint16
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd368826.aspx
type PIXELFORMATDESCRIPTOR struct {
	Size                   uint16
	Version                uint16
	DwFlags                uint32
	IPixelType             byte
	ColorBits              byte
	RedBits, RedShift      byte
	GreenBits, GreenShift  byte
	BlueBits, BlueShift    byte
	AlphaBits, AlphaShift  byte
	AccumBits              byte
	AccumRedBits           byte
	AccumGreenBits         byte
	AccumBlueBits          byte
	AccumAlphaBits         byte
	DepthBits, StencilBits byte
	AuxBuffers             byte
	ILayerType             byte
	Reserved               byte
	DwLayerMask            uint32
	DwVisibleMask          uint32
	DwDamageMask           uint32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms724950(v=vs.85).aspx
type SYSTEMTIME struct {
	Year         uint16
	Month        uint16
	DayOfWeek    uint16
	Day          uint16
	Hour         uint16
	Minute       uint16
	Second       uint16
	Milliseconds uint16
}

func (t SYSTEMTIME) Time() time.Time {
	return time.Date(
		int(t.Year),
		time.Month(t.Month),
		int(t.Day),
		int(t.Hour),
		int(t.Minute),
		int(t.Second),
		int(t.Milliseconds)*int(time.Millisecond/time.Nanosecond),
		time.UTC,
	)
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms644967(v=vs.85).aspx
type KBDLLHOOKSTRUCT struct {
	VkCode      DWORD
	ScanCode    DWORD
	Flags       DWORD
	Time        DWORD
	DwExtraInfo ULONG_PTR
}

type HOOKPROC func(int, WPARAM, LPARAM) LRESULT

type WINDOWPLACEMENT struct {
	Length           uint32
	Flags            uint32
	ShowCmd          uint32
	PtMinPosition    POINT
	PtMaxPosition    POINT
	RcNormalPosition RECT
}

type MINMAXINFO struct {
	PtReserved     POINT
	PtMaxSize      POINT
	PtMaxPosition  POINT
	PtMinTrackSize POINT
	PtMaxTrackSize POINT
}

type RAWINPUTHEADER struct {
	Type   uint32
	Size   uint32
	Device HANDLE
	WParam WPARAM
}

type RAWINPUT struct {
	Header RAWINPUTHEADER
	// NOTE that there is no support for C unions in Go, this would actually be
	// a union of RAWMOUSE, RAWKEYBOARD and RAWHID. Since RAWMOUSE is the
	// largest of those three, use it here and cast unsafely to get the other
	// types.
	mouse RAWMOUSE
}

// GetMouse returns the raw input as a RAWMOUSE. Make sure to check the Header's
// Type flag so this is valid.
func (i *RAWINPUT) GetMouse() RAWMOUSE {
	return i.mouse
}

// GetKeyboard returns the raw input as a RAWKEYBOARD. Make sure to check the
// Header's Type flag so this is valid.
func (i *RAWINPUT) GetKeyboard() RAWKEYBOARD {
	return *((*RAWKEYBOARD)(unsafe.Pointer(&i.mouse)))
}

// GetHid returns the raw input as a RAWHID. Make sure to check the Header's
// Type flag so this is valid.
func (i *RAWINPUT) GetHid() RAWHID {
	return *((*RAWHID)(unsafe.Pointer(&i.mouse)))
}

type RAWKEYBOARD struct {
	MakeCode         uint16
	Flags            uint16
	Reserved         uint16
	VKey             uint16
	Message          uint32
	ExtraInformation uint32
}

type RAWHID struct {
	SizeHid uint32
	Count   uint32
	RawData [1]byte
}

type RAWMOUSE struct {
	Flags            uint16
	Buttons          uint32
	RawButtons       uint32
	LastX            int32
	LastY            int32
	ExtraInformation uint32
}

func (m *RAWMOUSE) ButtonFlags() uint16 {
	return uint16(m.Buttons & 0xFFFF)
}

func (m *RAWMOUSE) ButtonData() uint16 {
	return uint16((m.Buttons & 0xFFFF0000) >> 16)
}

type RAWINPUTDEVICE struct {
	UsagePage uint16
	Usage     uint16
	Flags     uint32
	Target    HWND
}

// INPUT is used in SendInput. To create a concrete INPUT type, use the helper
// functions MouseInput, KeyboardInput and HardwareInput. These are necessary
// because the C API uses a union here, which Go does not provide.
type INPUT struct {
	Type uint32
	// use MOUSEINPUT for the union because it is the largest of all allowed
	// structures
	mouse MOUSEINPUT
}

type MOUSEINPUT struct {
	Dx        int32
	Dy        int32
	MouseData uint32
	Flags     uint32
	Time      uint32
	ExtraInfo uintptr
}

type KEYBDINPUT struct {
	Vk        uint16
	Scan      uint16
	Flags     uint32
	Time      uint32
	ExtraInfo uintptr
}

type HARDWAREINPUT struct {
	Msg    uint32
	ParamL uint16
	ParamH uint16
}

func MouseInput(input MOUSEINPUT) INPUT {
	return INPUT{
		Type:  INPUT_MOUSE,
		mouse: input,
	}
}

func KeyboardInput(input KEYBDINPUT) INPUT {
	return INPUT{
		Type:  INPUT_KEYBOARD,
		mouse: *((*MOUSEINPUT)(unsafe.Pointer(&input))),
	}
}

func HardwareInput(input HARDWAREINPUT) INPUT {
	return INPUT{
		Type:  INPUT_HARDWARE,
		mouse: *((*MOUSEINPUT)(unsafe.Pointer(&input))),
	}
}

type VS_FIXEDFILEINFO struct {
	Signature        uint32
	StrucVersion     uint32
	FileVersionMS    uint32
	FileVersionLS    uint32
	ProductVersionMS uint32
	ProductVersionLS uint32
	FileFlagsMask    uint32
	FileFlags        uint32
	FileOS           uint32
	FileType         uint32
	FileSubtype      uint32
	FileDateMS       uint32
	FileDateLS       uint32
}

// FileVersion concatenates FileVersionMS and FileVersionLS to a uint64 value.
func (fi VS_FIXEDFILEINFO) FileVersion() uint64 {
	return uint64(fi.FileVersionMS)<<32 | uint64(fi.FileVersionLS)
}

// FileDate concatenates FileDateMS and FileDateLS to a uint64 value.
func (fi VS_FIXEDFILEINFO) FileDate() uint64 {
	return uint64(fi.FileDateMS)<<32 | uint64(fi.FileDateLS)
}

type ACCEL struct {
	// Virt is a bit mask which may contain:
	//   FALT, FCONTROL, FSHIFT: keys to be held for the accelerator
	//   FVIRTKEY: means that Key is a virtual key code, if not set, Key is
	//             interpreted as a character code
	Virt byte
	// Key can either be a virtual key code VK_... or a character
	Key uint16
	// Cmd is the value passed to WM_COMMAND or WM_SYSCOMMAND when the
	// accelerator triggers
	Cmd uint16
}

type PHYSICAL_MONITOR struct {
	Monitor     HANDLE
	Description [128]uint16
}

type MENUITEMINFO struct {
	Size         uint32
	Mask         uint32
	Type         uint32
	State        uint32
	ID           uint32
	SubMenu      HMENU
	BmpChecked   HBITMAP
	BmpUnChecked HBITMAP
	ItemData     uintptr
	TypeData     uintptr // UTF-16 string
	CCH          uint32
	BmpItem      HBITMAP
}

type TPMPARAMS struct {
	Size    uint32
	Exclude RECT
}

type MENUINFO struct {
	size          uint32
	Mask          uint32
	Style         uint32
	YMax          uint32
	Back          HBRUSH
	ContextHelpID uint32
	MenuData      uintptr
}

type MENUBARINFO struct {
	size       uint32
	Bar        RECT
	Menu       HMENU
	Window     HWND
	BarFocused int32 // bool
	Focused    int32 // bool
}

type ACTCTX struct {
	size                  uint32
	Flags                 uint32
	Source                *uint16 // UTF-16 string
	ProcessorArchitecture uint16
	LangID                uint16
	AssemblyDirectory     *uint16 // UTF-16 string
	ResourceName          *uint16 // UTF-16 string
	ApplicationName       *uint16 // UTF-16 string
	Module                HMODULE
}

type DRAWITEMSTRUCT struct {
	CtlType    uint32
	CtlID      uint32
	ItemID     uint32
	ItemAction uint32
	ItemState  uint32
	HwndItem   HWND
	HDC        HDC
	RcItem     RECT
	ItemData   uintptr
}

type BLENDFUNC struct {
	BlendOp             byte
	BlendFlags          byte
	SourceConstantAlpha byte
	AlphaFormat         byte
}

type NETRESOURCE struct {
	Scope       uint32
	Type        uint32
	DisplayType uint32
	Usage       uint32
	LocalName   string
	RemoteName  string
	Comment     string
	Provider    string
}

func (n *NETRESOURCE) toInternal() *netresource {
	internal := &netresource{
		Scope:       n.Scope,
		Type:        n.Type,
		DisplayType: n.DisplayType,
		Usage:       n.Usage,
	}
	if n.LocalName != "" {
		internal.LocalName = syscall.StringToUTF16Ptr(n.LocalName)
	}
	if n.RemoteName != "" {
		internal.RemoteName = syscall.StringToUTF16Ptr(n.RemoteName)
	}
	if n.Comment != "" {
		internal.Comment = syscall.StringToUTF16Ptr(n.Comment)
	}
	if n.Provider != "" {
		internal.Provider = syscall.StringToUTF16Ptr(n.Provider)
	}
	return internal
}

type netresource struct {
	Scope       uint32
	Type        uint32
	DisplayType uint32
	Usage       uint32
	LocalName   *uint16
	RemoteName  *uint16
	Comment     *uint16
	Provider    *uint16
}

type NMLVODSTATECHANGE struct {
	Hdr      NMHDR
	From     int32
	To       int32
	NewState uint32
	OldState uint32
}

type SECURITY_ATTRIBUTES struct {
	Length             uint32
	SecurityDescriptor unsafe.Pointer
	InheritHandle      uint32 // bool value
}

type OVERLAPPED struct {
	Internal     uintptr
	InternalHigh uintptr
	Pointer      uintptr
	Event        HANDLE
}

type STORAGE_DEVICE_DESCRIPTOR struct {
	Version               uint32
	Size                  uint32
	DeviceType            byte
	DeviceTypeModifier    byte
	RemovableMedia        byte // bool value
	CommandQueueing       byte // bool value
	VendorIdOffset        uint32
	ProductIdOffset       uint32
	ProductRevisionOffset uint32
	SerialNumberOffset    uint32
	BusType               uint32 // STORAGE_BUS_TYPE
	RawPropertiesLength   uint32
	RawDeviceProperties   [1]byte
}

type STORAGE_PROPERTY_QUERY struct {
	PropertyId           uint32
	QueryType            uint32
	AdditionalParameters [1]byte
}

// https://docs.microsoft.com/en-us/windows/desktop/api/fileapi/ns-fileapi-_win32_find_stream_data
type WIN32_FIND_STREAM_DATA struct {
	Size int64
	Name [MAX_PATH + 36]uint16
}

type MSGBOXPARAMS struct {
	Size           uint32
	Owner          HWND
	Instance       HINSTANCE
	Text           *uint16
	Caption        *uint16
	Style          uint32
	Icon           *uint16
	ContextHelpId  *uint32
	MsgBoxCallback uintptr
	LanguageId     uint32
}

type POWERBROADCAST_SETTING struct {
	PowerSetting GUID
	DataLength   uint32
	Data         [1]byte
}

// https://docs.microsoft.com/en-us/windows-hardware/drivers/ddi/wdm/ns-wdm-_osversioninfoexw
type RTL_OSVERSIONINFOEXW struct {
	OSVersionInfoSize uint32
	MajorVersion      uint32
	MinorVersion      uint32
	BuildNumber       uint32
	PlatformId        uint32
	CSDVersion        [128]uint16
	ServicePackMajor  uint16
	ServicePackMinor  uint16
	SuiteMask         uint16
	ProductType       byte
	Reserved          byte
}

// https://docs.microsoft.com/en-us/windows/win32/api/sysinfoapi/ns-sysinfoapi-system_info
type SYSTEM_INFO struct {
	ProcessorArchitecture     uint16
	Reserved                  uint16
	PageSize                  uint32
	MinimumApplicationAddress LPCVOID
	MaximumApplicationAddress LPCVOID
	ActiveProcessorMask       *uint32
	NumberOfProcessors        uint32
	ProcessorType             uint32
	AllocationGranularity     uint32
	ProcessorLevel            uint16
	ProcessorRevision         uint16
}

type SP_DEVINFO_DATA struct {
	Size      uint32
	ClassGuid GUID
	DevInst   uint32
	Reserved  uintptr
}

type WINDOWPOS struct {
	Hwnd            HWND
	HwndInsertAfter HWND
	X               int32
	Y               int32
	Cx              int32
	Cy              int32
	Flags           uint32
}

type WINDOWINFO struct {
	Cbsize          uint32
	RcWindow        RECT
	RcClient        RECT
	DwStyle         uint32
	DwExStyle       uint32
	DwWindowStatus  uint32
	CxWindowBorders uint32
	CyWindowBorders uint32
	AtomWindowType  uint16
	WCreatorVersion uint16
}
