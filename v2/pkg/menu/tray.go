package menu

import (
	"context"
	"log"
	goruntime "runtime"
)

type TrayMenuAdd interface {
	TrayMenuAdd(menu *TrayMenu) TrayMenuImpl
}

type TrayMenuImpl interface {
	SetLabel(string)
	SetImage(*TrayImage)
}

type ImagePosition int

const (
	NSImageLeading  ImagePosition = 0
	NSImageOnly     ImagePosition = 1
	NSImageLeft     ImagePosition = 2
	NSImageRight    ImagePosition = 3
	NSImageBelow    ImagePosition = 4
	NSImageAbove    ImagePosition = 5
	NSImageOverlaps ImagePosition = 6
	NSNoImage       ImagePosition = 7
	NSImageTrailing ImagePosition = 8
)

type TraySizing int

const (
	Variable TraySizing = 0
	Square   TraySizing = 1
)

type TrayImage struct {
	Image      []byte
	Image2x    []byte
	IsTemplate bool
	Position   ImagePosition
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
