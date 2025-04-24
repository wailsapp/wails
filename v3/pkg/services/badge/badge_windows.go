//go:build windows

package badge

import (
	"context"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// COM GUIDs for taskbar interfaces
var (
	CLSID_TaskbarList = syscall.GUID{
		Data1: 0x56FDF344,
		Data2: 0xFD6D,
		Data3: 0x11D0,
		Data4: [8]byte{0x95, 0x8A, 0x00, 0x60, 0x97, 0xC9, 0xA0, 0x90},
	}
	IID_ITaskbarList3 = syscall.GUID{
		Data1: 0xEA1AFB91,
		Data2: 0x9E28,
		Data3: 0x4B86,
		Data4: [8]byte{0x90, 0xE9, 0x9E, 0x9F, 0x8A, 0x5E, 0xEF, 0xAF},
	}
)

// ITaskbarList3 COM interface
type ITaskbarList3Vtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
	// ITaskbarList methods
	HrInit       uintptr
	AddTab       uintptr
	DeleteTab    uintptr
	ActivateTab  uintptr
	SetActiveAlt uintptr
	// ITaskbarList2 methods
	MarkFullscreenWindow uintptr
	// ITaskbarList3 methods
	SetProgressValue      uintptr
	SetProgressState      uintptr
	RegisterTab           uintptr
	UnregisterTab         uintptr
	SetTabOrder           uintptr
	SetTabActive          uintptr
	ThumbBarAddButtons    uintptr
	ThumbBarUpdateButtons uintptr
	ThumbBarSetImageList  uintptr
	SetOverlayIcon        uintptr
	SetThumbnailTooltip   uintptr
	SetThumbnailClip      uintptr
	SetTabProperties      uintptr
}

type ITaskbarList3 struct {
	Vtbl *ITaskbarList3Vtbl
}

type windowsBadge struct {
	taskbarList  *ITaskbarList3
	mainWindow   *application.WebviewWindow
	initialized  bool
	redBadgeIcon syscall.Handle
}

func New() *Service {
	return &Service{
		impl: &windowsBadge{
			mainWindow:   application.Get().CurrentWindow(),
			redBadgeIcon: 0,
		},
	}
}

func (d *windowsBadge) Startup(ctx context.Context, options application.ServiceOptions) error {
	if d.initialized {
		return nil
	}

	// Initialize COM
	err := coInitialize()
	if err != nil {
		return fmt.Errorf("failed to initialize COM: %w", err)
	}

	// Create an instance of the TaskbarList COM object
	var taskbarList *ITaskbarList3
	hr, _, _ := procCoCreateInstance.Call(
		uintptr(unsafe.Pointer(&CLSID_TaskbarList)),
		0,
		21, // CLSCTX_INPROC_SERVER
		uintptr(unsafe.Pointer(&IID_ITaskbarList3)),
		uintptr(unsafe.Pointer(&taskbarList)))
	if hr != 0 {
		return fmt.Errorf("failed to create TaskbarList instance: %d", hr)
	}

	// Initialize the taskbar list
	hr, _, _ = syscall.Syscall(
		taskbarList.Vtbl.HrInit,
		1,
		uintptr(unsafe.Pointer(taskbarList)),
		0,
		0)
	if hr != 0 {
		return fmt.Errorf("failed to initialize TaskbarList: %d", hr)
	}

	// Create the red badge icon
	redBadge, err := d.createRedBadgeIcon()
	if err != nil {
		return fmt.Errorf("failed to create red badge icon: %w", err)
	}

	d.taskbarList = taskbarList
	d.redBadgeIcon = redBadge
	d.initialized = true
	return nil
}

func (d *windowsBadge) Shutdown() error {
	if !d.initialized {
		return nil
	}

	// Destroy the red badge icon
	if d.redBadgeIcon != 0 {
		destroyIcon(d.redBadgeIcon)
		d.redBadgeIcon = 0
	}

	// Release the taskbar list interface
	if d.taskbarList != nil {
		syscall.Syscall(
			d.taskbarList.Vtbl.Release,
			1,
			uintptr(unsafe.Pointer(d.taskbarList)),
			0,
			0)
		d.taskbarList = nil
	}

	// Uninitialize COM
	coUninitialize()
	d.initialized = false
	return nil
}

func (d *windowsBadge) SetBadge(label string) error {
	// If not initialized, initialize
	if !d.initialized {
		err := d.Startup(context.Background(), application.ServiceOptions{})
		if err != nil {
			return err
		}
	}

	// Ensure we have a window to work with
	if d.mainWindow == nil {
		d.mainWindow = application.Get().CurrentWindow()
		if d.mainWindow == nil {
			return fmt.Errorf("no window available for setting badge")
		}
	}

	// Get the window handle using NativeWindowHandle() method
	handle, err := d.mainWindow.NativeWindowHandle()
	if err != nil {
		return fmt.Errorf("failed to get window handle: %w", err)
	}
	hwnd := handle

	// If empty value, remove the overlay
	if label == "" {
		hr, _, _ := syscall.SyscallN(
			d.taskbarList.Vtbl.SetOverlayIcon,
			4,
			uintptr(unsafe.Pointer(d.taskbarList)),
			hwnd,
			0, // NULL icon handle
			0, // NULL description
			0, 0)
		if hr != 0 {
			return fmt.Errorf("failed to remove overlay icon: %d", hr)
		}
		return nil
	}

	// Set the overlay icon with a description
	description, err := syscall.UTF16PtrFromString("New notification")
	if err != nil {
		return fmt.Errorf("failed to convert description: %w", err)
	}

	hr, _, _ := syscall.SyscallN(
		d.taskbarList.Vtbl.SetOverlayIcon,
		4,
		uintptr(unsafe.Pointer(d.taskbarList)),
		hwnd,
		uintptr(d.redBadgeIcon),
		uintptr(unsafe.Pointer(description)),
		0, 0)
	if hr != 0 {
		return fmt.Errorf("failed to set overlay icon: %d", hr)
	}

	return nil
}

// ICONINFO structure for creating icons
type ICONINFO struct {
	FIcon    uint32
	XHotspot uint32
	YHotspot uint32
	HbmMask  syscall.Handle
	HbmColor syscall.Handle
}

// createRedBadgeIcon creates a simple red circle icon for the badge
func (d *windowsBadge) createRedBadgeIcon() (syscall.Handle, error) {
	// Create a 16x16 pixel bitmap
	hdc := getDC(0)
	if hdc == 0 {
		return 0, fmt.Errorf("failed to get DC")
	}
	defer releaseDC(0, hdc)

	memDC := createCompatibleDC(hdc)
	if memDC == 0 {
		return 0, fmt.Errorf("failed to create compatible DC")
	}
	defer deleteObject(memDC)

	// Create a bitmap
	bmp := createCompatibleBitmap(hdc, 16, 16)
	if bmp == 0 {
		return 0, fmt.Errorf("failed to create bitmap")
	}

	oldBmp := selectObject(memDC, bmp)

	// Create a solid red brush
	redBrush := createSolidBrush(0x0000FF) // BGR format (blue=0, green=0, red=255)
	defer deleteObject(redBrush)

	// Fill the circle with red
	selectObject(memDC, redBrush)
	ellipse(memDC, 0, 0, 16, 16)

	// Restore original bitmap
	selectObject(memDC, oldBmp)

	// Convert bitmap to icon
	iconInfo := ICONINFO{
		FIcon:    1, // TRUE for icon (vs. cursor)
		XHotspot: 0,
		YHotspot: 0,
		HbmMask:  bmp, // Use same bitmap for mask
		HbmColor: bmp,
	}

	icon := createIconIndirect(&iconInfo)
	if icon == 0 {
		deleteObject(bmp)
		return 0, fmt.Errorf("failed to create icon")
	}

	// Don't delete the bitmap here, as it's now owned by the icon
	return icon, nil
}

// Define Windows API functions
var (
	user32 = syscall.NewLazyDLL("user32.dll")
	gdi32  = syscall.NewLazyDLL("gdi32.dll")
	ole32  = syscall.NewLazyDLL("ole32.dll")

	procCoCreateInstance = ole32.NewProc("CoCreateInstance")
	procCoInitialize     = ole32.NewProc("CoInitialize")
	procCoUninitialize   = ole32.NewProc("CoUninitialize")

	procGetDC              = user32.NewProc("GetDC")
	procReleaseDC          = user32.NewProc("ReleaseDC")
	procDestroyIcon        = user32.NewProc("DestroyIcon")
	procCreateIconIndirect = user32.NewProc("CreateIconIndirect")

	procCreateCompatibleDC     = gdi32.NewProc("CreateCompatibleDC")
	procCreateCompatibleBitmap = gdi32.NewProc("CreateCompatibleBitmap")
	procSelectObject           = gdi32.NewProc("SelectObject")
	procDeleteObject           = gdi32.NewProc("DeleteObject")
	procCreateSolidBrush       = gdi32.NewProc("CreateSolidBrush")
	procEllipse                = gdi32.NewProc("Ellipse")
)

// GDI function wrappers
func getDC(hwnd syscall.Handle) syscall.Handle {
	ret, _, _ := procGetDC.Call(uintptr(hwnd))
	return syscall.Handle(ret)
}

func releaseDC(hwnd, hdc syscall.Handle) bool {
	ret, _, _ := procReleaseDC.Call(uintptr(hwnd), uintptr(hdc))
	return ret != 0
}

func createCompatibleDC(hdc syscall.Handle) syscall.Handle {
	ret, _, _ := procCreateCompatibleDC.Call(uintptr(hdc))
	return syscall.Handle(ret)
}

func createCompatibleBitmap(hdc syscall.Handle, width, height int32) syscall.Handle {
	ret, _, _ := procCreateCompatibleBitmap.Call(uintptr(hdc), uintptr(width), uintptr(height))
	return syscall.Handle(ret)
}

func selectObject(hdc, hgdiobj syscall.Handle) syscall.Handle {
	ret, _, _ := procSelectObject.Call(uintptr(hdc), uintptr(hgdiobj))
	return syscall.Handle(ret)
}

func deleteObject(hObject syscall.Handle) bool {
	ret, _, _ := procDeleteObject.Call(uintptr(hObject))
	return ret != 0
}

func createSolidBrush(color int32) syscall.Handle {
	ret, _, _ := procCreateSolidBrush.Call(uintptr(color))
	return syscall.Handle(ret)
}

func ellipse(hdc syscall.Handle, left, top, right, bottom int32) bool {
	ret, _, _ := procEllipse.Call(uintptr(hdc), uintptr(left), uintptr(top), uintptr(right), uintptr(bottom))
	return ret != 0
}

func createIconIndirect(iconInfo *ICONINFO) syscall.Handle {
	ret, _, _ := procCreateIconIndirect.Call(uintptr(unsafe.Pointer(iconInfo)))
	return syscall.Handle(ret)
}

func destroyIcon(hIcon syscall.Handle) bool {
	ret, _, _ := procDestroyIcon.Call(uintptr(hIcon))
	return ret != 0
}

func coInitialize() error {
	hr, _, _ := procCoInitialize.Call(0)
	if hr != 0 {
		return fmt.Errorf("CoInitialize failed with code: %d", hr)
	}
	return nil
}

func coUninitialize() {
	procCoUninitialize.Call()
}
