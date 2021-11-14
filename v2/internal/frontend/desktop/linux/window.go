//go:build linux
// +build linux

package linux

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"log"
	"math"
	"sync"
)

type Window struct {
	frontendOptions *options.App
	applicationMenu *menu.Menu
	m               sync.Mutex
	application     *gtk.Application
	gtkWindow       *gtk.ApplicationWindow

	//dispatchq       []func()
}

func NewWindow(options *options.App) *Window {
	result := new(Window)
	result.frontendOptions = options

	var linuxOptions linux.Options
	if options.Linux != nil {
		linuxOptions = *options.Linux
	}
	appID := linuxOptions.AppID
	if appID == "" {
		appID = "io.wails"
	}

	println("AppID =", appID)

	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Fatal("Could not create application:", err)
	}

	result.application = application

	application.Connect("activate", func() {

		window, err := gtk.ApplicationWindowNew(application)
		if err != nil {
			log.Fatal("Could not create application window:", err)
		}
		window.Connect("delete-event", func() {
			if options.HideWindowOnClose {
				result.gtkWindow.Hide()
				return
			}
			result.gtkWindow.Close()
		})
		result.gtkWindow = window
		window.SetTitle(options.Title)
		window.SetDecorated(!options.Frameless)
		window.SetDefaultSize(600, 300)
		window.SetResizable(!options.DisableResize)
		window.SetKeepAbove(options.AlwaysOnTop)
		window.SetPosition(gtk.WIN_POS_CENTER)
		if !options.StartHidden {
			window.ShowAll()
		}
	})

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
	w.application.Run(nil)
}

func (w *Window) Dispatch(f func()) {
	glib.IdleAdd(f)
}

func (w *Window) Fullscreen() {
	w.gtkWindow.Fullscreen()
}

func (w *Window) UnFullscreen() {
	w.gtkWindow.Unfullscreen()
}

func (w *Window) Close() {
	w.application.Quit()
}

func (w *Window) Center() {
	w.gtkWindow.SetPosition(gtk.WIN_POS_CENTER)
}

func (w *Window) SetPos(x int, y int) {
	display, err := w.gtkWindow.GetDisplay()
	if err != nil {
		w.gtkWindow.Move(x, y)
		return
	}
	window, err := w.gtkWindow.GetWindow()
	if err != nil {
		w.gtkWindow.Move(x, y)
		return
	}

	monitor, err := display.GetMonitorAtWindow(window)
	if err != nil {
		w.gtkWindow.Move(x, y)
		return
	}

	geom := monitor.GetGeometry()
	w.gtkWindow.Move(geom.GetX()+x, geom.GetY()+y)
}

func (w *Window) Pos() (int, int) {
	return w.gtkWindow.GetPosition()
}

func (w *Window) SetSize(width int, height int) {
	w.gtkWindow.SetDefaultSize(width, height)
}

func (w *Window) Size() (int, int) {
	return w.gtkWindow.GetSize()
}

func (w *Window) SetText(title string) {
	w.gtkWindow.SetTitle(title)
}

func (w *Window) SetMaxSize(maxWidth int, maxHeight int) {
	var geom gdk.Geometry
	if maxWidth == 0 {
		maxWidth = math.MaxInt
	}
	if maxHeight == 0 {
		maxHeight = math.MaxInt
	}
	geom.SetMaxWidth(maxWidth)
	geom.SetMaxHeight(maxHeight)
	w.gtkWindow.SetGeometryHints(w.gtkWindow, geom, gdk.HINT_MAX_SIZE)
}

func (w *Window) SetMinSize(minWidth int, minHeight int) {
	var geom gdk.Geometry
	geom.SetMinWidth(minWidth)
	geom.SetMinHeight(minHeight)
	w.gtkWindow.SetGeometryHints(w.gtkWindow, geom, gdk.HINT_MIN_SIZE)
}

func (w *Window) Show() {
	w.gtkWindow.ShowAll()
}

func (w *Window) Hide() {
	w.gtkWindow.Hide()
}

func (w *Window) Maximise() {
	w.gtkWindow.Maximize()
}

func (w *Window) UnMaximise() {
	w.gtkWindow.Unmaximize()
}

func (w *Window) Minimise() {
	w.gtkWindow.Iconify()
}

func (w *Window) UnMinimise() {
	w.gtkWindow.Present()
}

func (w *Window) IsFullScreen() bool {
	return false
}
