package menu

import (
	"context"
	"log"
	goruntime "runtime"

	"github.com/wailsapp/wails/v2/pkg/events"
)

type TrayMenuAdd interface {
	TrayMenuAdd(menu *TrayMenu) TrayMenuImpl
}

type TrayMenuImpl interface {
	SetLabel(string)
	SetImage(*TrayImage)
	SetMenu(*Menu)
}

type EventsImpl interface {
	On(eventName string, callback func(...interface{}))
}

type ImagePosition int

const (
	ImageLeading  ImagePosition = 0
	ImageOnly     ImagePosition = 1
	ImageLeft     ImagePosition = 2
	ImageRight    ImagePosition = 3
	ImageBelow    ImagePosition = 4
	ImageAbove    ImagePosition = 5
	ImageOverlaps ImagePosition = 6
	NoImage       ImagePosition = 7
	ImageTrailing ImagePosition = 8
)

type TraySizing int

const (
	Variable TraySizing = 0
	Square   TraySizing = 1
)

type TrayImage struct {
	// Bitmaps hold images for different scaling factors
	// First = 1x, Second = 2x, etc
	Bitmaps     [][]byte
	BitmapsDark [][]byte
	IsTemplate  bool
	Position    ImagePosition
}

func (t *TrayImage) getBestBitmap(scale int, isDarkMode bool) []byte {
	bitmapsToCheck := t.Bitmaps
	if isDarkMode {
		bitmapsToCheck = t.BitmapsDark
	}
	if scale < 1 || scale >= len(bitmapsToCheck) {
		return nil
	}
	for i := scale; i > 0; i-- {
		if bitmapsToCheck[i] != nil {
			return bitmapsToCheck[i]
		}
	}
	return nil
}

// GetBestBitmap will attempt to return the best bitmap for the theme
// If dark theme is used and no dark theme bitmap exists, then it will
// revert to light theme bitmaps
func (t *TrayImage) GetBestBitmap(scale int, isDarkMode bool) []byte {
	var result []byte
	if isDarkMode {
		result = t.getBestBitmap(scale, true)
		if result != nil {
			return result
		}
	}
	return t.getBestBitmap(scale, false)
}

// TrayMenu are the options
type TrayMenu struct {
	ctx context.Context

	// Label is the text we wish to display in the tray
	Label string

	Image *TrayImage

	// Text Colour
	RGBA string

	// Font
	FontSize int
	FontName string

	// Tooltip
	Tooltip string

	// Callback function when menu clicked
	Click Callback

	// Disabled makes the item unselectable
	Disabled bool

	// Menu is the initial menu we wish to use for the tray
	Menu *Menu

	// OnOpen is called when the Menu is opened
	OnOpen func()

	// OnClose is called when the Menu is closed
	OnClose func()

	/* Mac Options */
	Sizing TraySizing

	// This is the reference to the OS specific implementation
	impl TrayMenuImpl

	// Theme change callback
	themeChangeCallback func(data ...interface{})
}

func NewTrayMenu() *TrayMenu {
	return &TrayMenu{}
}

func (t *TrayMenu) Show(ctx context.Context) {
	if ctx == nil {
		log.Fatal("TrayMenu.Show() called before Run()")
	}
	t.ctx = ctx
	result := ctx.Value("frontend")
	if result == nil {
		pc, _, _, _ := goruntime.Caller(1)
		funcName := goruntime.FuncForPC(pc).Name()
		log.Fatalf("invalid context at '%s'", funcName)
	}
	t.impl = result.(TrayMenuAdd).TrayMenuAdd(t)

	if t.themeChangeCallback == nil {
		t.themeChangeCallback = func(data ...interface{}) {
			println("Update button image")
			if t.Image != nil {
				// Update the image
				t.SetImage(t.Image)
			}
		}
		result := ctx.Value("events")
		if result != nil {
			result.(EventsImpl).On(events.ThemeChanged, t.themeChangeCallback)
		}
	}

}

func (t *TrayMenu) SetLabel(label string) {
	t.Label = label
	if t.impl != nil {
		t.impl.SetLabel(label)
	}
}

func (t *TrayMenu) SetImage(image *TrayImage) {
	t.Image = image
	if t.impl != nil {
		t.impl.SetImage(image)
	}
}

func (t *TrayMenu) SetMenu(menu *Menu) {
	t.Menu = menu
	if t.impl != nil {
		t.impl.SetMenu(menu)
	}
}
