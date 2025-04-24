//go:build windows

package badge

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/w32"
)

var (
	ole32            = syscall.NewLazyDLL("ole32.dll")
	shobjidl         = syscall.NewLazyDLL("shell32.dll")
	coCreateInstance = ole32.NewProc("CoCreateInstance")
)

const (
	CLSCTX_INPROC_SERVER = 0x1
)

var (
	CLSID_TaskbarList = syscall.GUID{0x56FDF344, 0xFD6D, 0x11D0, [8]byte{0x95, 0x8A, 0x00, 0x60, 0x97, 0xC9, 0xA0, 0x90}}
	IID_ITaskbarList3 = syscall.GUID{0xEA1AFB91, 0x9E28, 0x4B86, [8]byte{0x90, 0xE9, 0x9E, 0x9F, 0x8A, 0x5E, 0xEF, 0xAF}}
)

type ITaskbarList3 struct {
	lpVtbl *taskbarList3Vtbl
}

type taskbarList3Vtbl struct {
	QueryInterface        uintptr
	AddRef                uintptr
	Release               uintptr
	HrInit                uintptr
	AddTab                uintptr
	DeleteTab             uintptr
	ActivateTab           uintptr
	SetActiveAlt          uintptr
	MarkFullscreenWindow  uintptr
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
}

func newTaskbarList3() (*ITaskbarList3, error) {
	var taskbar *ITaskbarList3
	hr, _, _ := coCreateInstance.Call(
		uintptr(unsafe.Pointer(&CLSID_TaskbarList)),
		0,
		uintptr(CLSCTX_INPROC_SERVER),
		uintptr(unsafe.Pointer(&IID_ITaskbarList3)),
		uintptr(unsafe.Pointer(&taskbar)),
	)

	if hr != 0 {
		return nil, syscall.Errno(hr)
	}

	return taskbar, nil
}

func (t *ITaskbarList3) SetOverlayIcon(hwnd syscall.Handle, hIcon syscall.Handle, description *uint16) error {
	ret, _, _ := syscall.SyscallN(
		t.lpVtbl.SetOverlayIcon,
		uintptr(unsafe.Pointer(t)),
		uintptr(hwnd),
		uintptr(hIcon),
		uintptr(unsafe.Pointer(description)),
	)
	if ret != 0 {
		return syscall.Errno(ret)
	}
	return nil
}

type windowsBadge struct {
	taskbar *ITaskbarList3
}

func New() *Service {
	return &Service{
		impl: &windowsBadge{},
	}
}

func (d *windowsBadge) Startup(ctx context.Context, options application.ServiceOptions) error {
	taskbar, err := newTaskbarList3()
	if err != nil {
		return err
	}
	d.taskbar = taskbar

	// Don't try to get the window handle here - wait until SetBadge is called
	return nil
}

func (d *windowsBadge) Shutdown() error {
	return nil
}

func (d *windowsBadge) SetBadge(label string) error {
	if d.taskbar == nil {
		return nil
	}

	// Get the window handle when SetBadge is called, not during startup
	app := application.Get()
	if app == nil {
		return nil // App not initialized yet
	}

	window := app.CurrentWindow()
	if window == nil {
		return nil // No window available yet
	}

	hwnd, err := window.NativeWindowHandle()
	if err != nil {
		return err
	}

	if label == "" {
		return d.taskbar.SetOverlayIcon(syscall.Handle(hwnd), 0, nil)
	}

	hicon, err := createBadgeIcon()
	if err != nil {
		return err
	}
	defer w32.DestroyIcon(hicon)

	return d.taskbar.SetOverlayIcon(syscall.Handle(hwnd), syscall.Handle(hicon), nil)
}

func (d *windowsBadge) RemoveBadge() error {
	if d.taskbar == nil {
		return nil
	}

	// Get the window handle when SetBadge is called, not during startup
	app := application.Get()
	if app == nil {
		return nil // App not initialized yet
	}

	window := app.CurrentWindow()
	if window == nil {
		return nil // No window available yet
	}

	hwnd, err := window.NativeWindowHandle()
	if err != nil {
		return err
	}

	return d.taskbar.SetOverlayIcon(syscall.Handle(hwnd), 0, nil)
}

func createBadgeIcon() (uintptr, error) {
	const size = 32

	img := image.NewRGBA(image.Rect(0, 0, size, size))

	red := color.RGBA{255, 0, 0, 255}
	radius := size / 2
	centerX, centerY := radius, radius

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := float64(x - centerX)
			dy := float64(y - centerY)

			if dx*dx+dy*dy < float64(radius*radius) {
				img.Set(x, y, red)
			}
		}
	}

	white := color.RGBA{255, 255, 255, 255}
	innerRadius := size / 5

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := float64(x - centerX)
			dy := float64(y - centerY)

			if dx*dx+dy*dy < float64(innerRadius*innerRadius) {
				img.Set(x, y, white)
			}
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return 0, err
	}

	hicon, err := w32.CreateSmallHIconFromImage(buf.Bytes())
	return uintptr(hicon), err
}
