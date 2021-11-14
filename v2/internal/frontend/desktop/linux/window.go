//go:build linux
// +build linux

package linux

import (
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"sync"
)

type Window struct {
	frontendOptions *options.App
	applicationMenu *menu.Menu
	m               sync.Mutex
	//dispatchq       []func()
}

func NewWindow(options *options.App) *Window {
	result := new(Window)
	result.frontendOptions = options
	//result.SetIsForm(true)
	//
	//var exStyle int
	//if options.Windows != nil {
	//	exStyle = w32.WS_EX_CONTROLPARENT | w32.WS_EX_APPWINDOW
	//	if options.Windows.WindowIsTranslucent {
	//		exStyle |= w32.WS_EX_NOREDIRECTIONBITMAP
	//	}
	//}
	//if options.AlwaysOnTop {
	//	exStyle |= w32.WS_EX_TOPMOST
	//}
	//
	//var dwStyle = w32.WS_OVERLAPPEDWINDOW
	//if options.Frameless {
	//	dwStyle = w32.WS_POPUP
	//}
	//
	//winc.RegClassOnlyOnce("wailsWindow")
	//result.SetHandle(winc.CreateWindow("wailsWindow", parent, uint(exStyle), uint(dwStyle)))
	//result.SetParent(parent)
	//
	//loadIcon := true
	//if options.Windows != nil && options.Windows.DisableWindowIcon == true {
	//	loadIcon = false
	//}
	//if loadIcon {
	//	if ico, err := winc.NewIconFromResource(winc.GetAppInstance(), uint16(winc.AppIconID)); err == nil {
	//		result.SetIcon(0, ico)
	//	}
	//}
	//
	//result.SetSize(options.Width, options.Height)
	//result.SetText(options.Title)
	//if options.Frameless == false && !options.Fullscreen {
	//	result.EnableMaxButton(!options.DisableResize)
	//	result.EnableSizable(!options.DisableResize)
	//	result.SetMinSize(options.MinWidth, options.MinHeight)
	//	result.SetMaxSize(options.MaxWidth, options.MaxHeight)
	//}
	//
	//if options.Windows != nil {
	//	if options.Windows.WindowIsTranslucent {
	//		result.SetTranslucentBackground()
	//	}
	//
	//	if options.Windows.DisableWindowIcon {
	//		result.DisableIcon()
	//	}
	//}
	//
	//// Dlg forces display of focus rectangles, as soon as the user starts to type.
	//w32.SendMessage(result.Handle(), w32.WM_CHANGEUISTATE, w32.UIS_INITIALIZE, 0)
	//winc.RegMsgHandler(result)
	//
	//result.SetFont(winc.DefaultFont)
	//
	//if options.Menu != nil {
	//	result.SetApplicationMenu(options.Menu)
	//}

	return result
}

func (w *Window) Run() {

}

func (w *Window) Dispatch(f func()) {
	//w.m.Lock()
	//w.dispatchq = append(w.dispatchq, f)
	//w.m.Unlock()
	//w32.PostMainThreadMessage(w32.WM_APP, 0, 0)
}

func (w *Window) Fullscreen() {

}

func (w *Window) Close() {

}

func (w *Window) Center() {

}

func (w *Window) SetPos(x int, y int) {

}

func (w *Window) Pos() (int, int) {
	return 0, 0
}

func (w *Window) SetSize(width int, height int) {

}

func (w *Window) Size() (int, int) {
	return 0, 0

}

func (w *Window) SetText(title string) {

}

func (w *Window) SetMaxSize(maxWidth int, maxHeight int) {

}

func (w *Window) SetMinSize(minWidth int, minHeight int) {

}

func (w *Window) UnFullscreen() {

}

func (w *Window) Show() {

}

func (w *Window) Hide() {

}

func (w *Window) Maximise() {

}

func (w *Window) Restore() {

}

func (w *Window) Minimise() {

}

func (w *Window) IsFullScreen() bool {
	return false
}
