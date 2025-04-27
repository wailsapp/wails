//go:build windows

package badge

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/w32"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type windowsBadge struct {
	taskbar     *w32.ITaskbarList3
	badgeImg    *image.RGBA
	badgeSize   int
	fontManager *FontManager
	options     Options
}

var defaultOptions = Options{
	TextColour:       color.RGBA{255, 255, 255, 255},
	BackgroundColour: color.RGBA{255, 0, 0, 255},
	FontName:         "segoeuib.ttf",
	FontSize:         18,
	SmallFontSize:    14,
}

// Creates a new Badge Service.
func New() *Service {
	return &Service{
		impl: &windowsBadge{
			options: defaultOptions,
		},
	}
}

// NewWithOptions creates a new badge service with the given options.
func NewWithOptions(options Options) *Service {
	return &Service{
		impl: &windowsBadge{
			options: options,
		},
	}
}

func (w *windowsBadge) Startup(ctx context.Context, options application.ServiceOptions) error {
	taskbar, err := w32.NewTaskbarList3()
	if err != nil {
		return err
	}
	w.taskbar = taskbar
	w.fontManager = NewFontManager()

	return nil
}

func (w *windowsBadge) Shutdown() error {
	if w.taskbar != nil {
		w.taskbar.Release()
		w32.CoUninitialize()
	}

	return nil
}

// SetBadge sets the badge label on the application icon.
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

	return w.taskbar.SetOverlayIcon(hwnd, hicon, nil)
}

// RemoveBadge removes the badge label from the application icon.
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

	return w.taskbar.SetOverlayIcon(hwnd, 0, nil)
}

// createBadgeIcon creates a badge icon with the specified size and color.
func (w *windowsBadge) createBadgeIcon() (w32.HICON, error) {
	radius := w.badgeSize / 2
	centerX, centerY := radius, radius
	innerRadius := w.badgeSize / 5

	for y := 0; y < w.badgeSize; y++ {
		for x := 0; x < w.badgeSize; x++ {
			dx := float64(x - centerX)
			dy := float64(y - centerY)

			if dx*dx+dy*dy < float64(innerRadius*innerRadius) {
				w.badgeImg.Set(x, y, w.options.TextColour)
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

// createBadgeIconWithText creates a badge icon with the specified text.
func (w *windowsBadge) createBadgeIconWithText(label string) (w32.HICON, error) {

	fontPath := w.fontManager.FindFontOrDefault(w.options.FontName)
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

	fontSize := float64(w.options.FontSize)
	if len(label) > 1 {
		fontSize = float64(w.options.SmallFontSize)
	}

	// Get DPI of the current screen
	screen := w32.GetDesktopWindow()
	dpi := w32.GetDpiForWindow(screen)

	face, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     float64(dpi),
		Hinting: font.HintingFull,
	})
	if err != nil {
		return 0, err
	}
	defer face.Close()

	d := &font.Drawer{
		Dst:  w.badgeImg,
		Src:  image.NewUniform(w.options.TextColour),
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

// createBadge creates a circular badge with the specified background color.
func (w *windowsBadge) createBadge() {
	w.badgeSize = 32

	img := image.NewRGBA(image.Rect(0, 0, w.badgeSize, w.badgeSize))

	backgroundColour := w.options.BackgroundColour
	radius := w.badgeSize / 2
	centerX, centerY := radius, radius

	for y := 0; y < w.badgeSize; y++ {
		for x := 0; x < w.badgeSize; x++ {
			dx := float64(x - centerX)
			dy := float64(y - centerY)

			if dx*dx+dy*dy < float64(radius*radius) {
				img.Set(x, y, backgroundColour)
			}
		}
	}

	w.badgeImg = img
}
