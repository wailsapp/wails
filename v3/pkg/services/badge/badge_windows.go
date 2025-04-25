//go:build windows

package badge

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"github.com/flopp/go-findfont"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/w32"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var (
	ole32            = syscall.NewLazyDLL("ole32.dll")
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
	const COINIT_APARTMENTTHREADED = 0x2

	coInit := ole32.NewProc("CoInitializeEx")
	if hr, _, _ := coInit.Call(0, COINIT_APARTMENTTHREADED); hr != 0 && hr != 0x1 {
		return nil, syscall.Errno(hr)
	}

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

	if r, _, _ := syscall.SyscallN(taskbar.lpVtbl.HrInit, uintptr(unsafe.Pointer(taskbar))); r != 0 {
		syscall.SyscallN(taskbar.lpVtbl.Release, uintptr(unsafe.Pointer(taskbar)))
		ole32.NewProc("CoUninitialize").Call()
		return nil, syscall.Errno(r)
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
	taskbar   *ITaskbarList3
	badgeImg  *image.RGBA
	badgeSize int
}

func New() *Service {
	return &Service{
		impl: &windowsBadge{},
	}
}

func (w *windowsBadge) Startup(ctx context.Context, options application.ServiceOptions) error {
	taskbar, err := newTaskbarList3()
	if err != nil {
		return err
	}
	w.taskbar = taskbar

	return nil
}

func (w *windowsBadge) Shutdown() error {
	if w.taskbar != nil {
		syscall.SyscallN(w.taskbar.lpVtbl.Release, uintptr(unsafe.Pointer(w.taskbar)))
		ole32.NewProc("CoUninitialize").Call()
	}

	return nil
}

func (w *windowsBadge) SetBadge(label string) error {
	if w.taskbar == nil {
		return nil
	}

	app := application.Get()
	if app == nil {
		return nil
	}

	window := app.CurrentWindow()
	if window == nil {
		return nil
	}

	hwnd, err := window.NativeWindowHandle()
	if err != nil {
		return err
	}

	w.createBadge()

	var hicon w32.HICON
	if label == "" {
		hicon, err = w.createBadgeIcon()
		if err != nil {
			return err
		}
	} else {
		hicon, err = w.createBadgeIconWithText(label)
		if err != nil {
			return err
		}
	}
	defer w32.DestroyIcon(hicon)

	return w.taskbar.SetOverlayIcon(syscall.Handle(hwnd), syscall.Handle(hicon), nil)
}

func (w *windowsBadge) RemoveBadge() error {
	if w.taskbar == nil {
		return nil
	}

	app := application.Get()
	if app == nil {
		return nil
	}

	window := app.CurrentWindow()
	if window == nil {
		return nil
	}

	hwnd, err := window.NativeWindowHandle()
	if err != nil {
		return err
	}

	return w.taskbar.SetOverlayIcon(syscall.Handle(hwnd), 0, nil)
}

func (w *windowsBadge) createBadgeIcon() (w32.HICON, error) {
	radius := w.badgeSize / 2
	centerX, centerY := radius, radius
	white := color.RGBA{255, 255, 255, 255}
	innerRadius := w.badgeSize / 5

	for y := 0; y < w.badgeSize; y++ {
		for x := 0; x < w.badgeSize; x++ {
			dx := float64(x - centerX)
			dy := float64(y - centerY)

			if dx*dx+dy*dy < float64(innerRadius*innerRadius) {
				w.badgeImg.Set(x, y, white)
			}
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, w.badgeImg); err != nil {
		return 0, err
	}

	hicon, err := w32.CreateSmallHIconFromImage(buf.Bytes())
	return hicon, err
}

func (w *windowsBadge) createBadgeIconWithText(label string) (w32.HICON, error) {
	fontPath := ""
	for _, path := range findfont.List() {
		if strings.Contains(strings.ToLower(path), "segoeuib.ttf") || // Segoe UI Bold
			strings.Contains(strings.ToLower(path), "arialbd.ttf") {
			fontPath = path
			break
		}
	}
	if fontPath == "" {
		return w.createBadgeIcon()
	}

	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		return 0, err
	}

	ttf, err := opentype.Parse(fontBytes)
	if err != nil {
		return 0, err
	}

	fontSize := 18.0
	if len(label) > 1 {
		fontSize = 14.0
	}

	face, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     96,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return 0, err
	}
	defer face.Close()

	d := &font.Drawer{
		Dst:  w.badgeImg,
		Src:  image.NewUniform(color.White),
		Face: face,
	}

	textWidth := d.MeasureString(label).Ceil()
	d.Dot = fixed.P((w.badgeSize-textWidth)/2, int(float64(w.badgeSize)/2+fontSize/2))
	d.DrawString(label)

	var buf bytes.Buffer
	if err := png.Encode(&buf, w.badgeImg); err != nil {
		return 0, err
	}

	return w32.CreateSmallHIconFromImage(buf.Bytes())
}

func (w *windowsBadge) createBadge() {
	w.badgeSize = 32

	img := image.NewRGBA(image.Rect(0, 0, w.badgeSize, w.badgeSize))

	red := color.RGBA{255, 0, 0, 255}
	radius := w.badgeSize / 2
	centerX, centerY := radius, radius

	for y := 0; y < w.badgeSize; y++ {
		for x := 0; x < w.badgeSize; x++ {
			dx := float64(x - centerX)
			dy := float64(y - centerY)

			if dx*dx+dy*dy < float64(radius*radius) {
				img.Set(x, y, red)
			}
		}
	}

	w.badgeImg = img
}
