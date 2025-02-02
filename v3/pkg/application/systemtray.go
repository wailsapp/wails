package application

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/events"
)

type IconPosition int

const (
	NSImageNone = iota
	NSImageOnly
	NSImageLeft
	NSImageRight
	NSImageBelow
	NSImageAbove
	NSImageOverlaps
	NSImageLeading
	NSImageTrailing
)

type systemTrayImpl interface {
	setLabel(label string)
	run()
	setIcon(icon []byte)
	setMenu(menu *Menu)
	setIconPosition(position IconPosition)
	setTemplateIcon(icon []byte)
	destroy()
	setDarkModeIcon(icon []byte)
	bounds() (*Rect, error)
	getScreen() (*Screen, error)
	positionWindow(window *WebviewWindow, offset int) error
	openMenu()
	Show()
	Hide()
}

type SystemTray struct {
	id           uint
	label        string
	icon         []byte
	darkModeIcon []byte
	iconPosition IconPosition

	clickHandler            func()
	rightClickHandler       func()
	doubleClickHandler      func()
	rightDoubleClickHandler func()
	mouseEnterHandler       func()
	mouseLeaveHandler       func()
	onMenuOpen              func()
	onMenuClose             func()

	// Platform specific implementation
	impl           systemTrayImpl
	menu           *Menu
	isTemplateIcon bool
	attachedWindow WindowAttachConfig
}

func newSystemTray(id uint) *SystemTray {
	result := &SystemTray{
		id:           id,
		label:        "",
		iconPosition: NSImageLeading,
		attachedWindow: WindowAttachConfig{
			Window:   nil,
			Offset:   0,
			Debounce: 200 * time.Millisecond,
		},
	}
	result.clickHandler = result.defaultClickHandler
	return result
}

func (s *SystemTray) SetLabel(label string) {
	if s.impl == nil {
		s.label = label
		return
	}
	InvokeSync(func() {
		s.impl.setLabel(label)
	})
}

func (s *SystemTray) Label() string {
	return s.label
}

func (s *SystemTray) Run() {
	s.impl = newSystemTrayImpl(s)

	if s.attachedWindow.Window != nil {
		// Setup listener
		s.attachedWindow.Window.OnWindowEvent(events.Common.WindowLostFocus, func(event *WindowEvent) {
			s.attachedWindow.Window.Hide()
			// Special handler for Windows
			if runtime.GOOS == "windows" {
				// We don't do this unless the window has already been shown
				if s.attachedWindow.hasBeenShown == false {
					return
				}
				s.attachedWindow.justClosed = true
				go func() {
					defer handlePanic()
					time.Sleep(s.attachedWindow.Debounce)
					s.attachedWindow.justClosed = false
				}()
			}
		})
	}

	InvokeSync(s.impl.run)
}

func (s *SystemTray) PositionWindow(window *WebviewWindow, offset int) error {
	if s.impl == nil {
		return fmt.Errorf("system tray not running")
	}
	return InvokeSyncWithError(func() error {
		return s.impl.positionWindow(window, offset)
	})
}

func (s *SystemTray) SetIcon(icon []byte) *SystemTray {
	if s.impl == nil {
		s.icon = icon
	} else {
		InvokeSync(func() {
			s.impl.setIcon(icon)
		})
	}
	return s
}

func (s *SystemTray) SetDarkModeIcon(icon []byte) *SystemTray {
	if s.impl == nil {
		s.darkModeIcon = icon
	} else {
		InvokeSync(func() {
			s.impl.setDarkModeIcon(icon)
		})
	}
	return s
}

func (s *SystemTray) SetMenu(menu *Menu) *SystemTray {
	if s.impl == nil {
		s.menu = menu
	} else {
		InvokeSync(func() {
			s.impl.setMenu(menu)
		})
	}
	return s
}

func (s *SystemTray) SetIconPosition(iconPosition IconPosition) *SystemTray {
	if s.impl == nil {
		s.iconPosition = iconPosition
	} else {
		InvokeSync(func() {
			s.impl.setIconPosition(iconPosition)
		})
	}
	return s
}

func (s *SystemTray) SetTemplateIcon(icon []byte) *SystemTray {
	if s.impl == nil {
		s.icon = icon
		s.isTemplateIcon = true
	} else {
		InvokeSync(func() {
			s.impl.setTemplateIcon(icon)
		})
	}
	return s
}

func (s *SystemTray) Destroy() {
	globalApplication.destroySystemTray(s)
}

func (s *SystemTray) destroy() {
	if s.impl == nil {
		return
	}
	s.impl.destroy()
}

func (s *SystemTray) OnClick(handler func()) *SystemTray {
	s.clickHandler = handler
	return s
}

func (s *SystemTray) OnRightClick(handler func()) *SystemTray {
	s.rightClickHandler = handler
	return s
}

func (s *SystemTray) OnDoubleClick(handler func()) *SystemTray {
	s.doubleClickHandler = handler
	return s
}

func (s *SystemTray) OnRightDoubleClick(handler func()) *SystemTray {
	s.rightDoubleClickHandler = handler
	return s
}

func (s *SystemTray) OnMouseEnter(handler func()) *SystemTray {
	s.mouseEnterHandler = handler
	return s
}

func (s *SystemTray) OnMouseLeave(handler func()) *SystemTray {
	s.mouseLeaveHandler = handler
	return s
}

func (s *SystemTray) Show() {
	if s.impl == nil {
		return
	}
	InvokeSync(func() {
		s.impl.Show()
	})
}

func (s *SystemTray) Hide() {
	if s.impl == nil {
		return
	}
	InvokeSync(func() {
		s.impl.Hide()
	})
}

type WindowAttachConfig struct {
	// Window is the window to attach to the system tray. If it's null, the request to attach will be ignored.
	Window *WebviewWindow

	// Offset indicates the gap in pixels between the system tray and the window
	Offset int

	// Debounce is used by Windows to indicate how long to wait before responding to a mouse
	// up event on the notification icon. See https://stackoverflow.com/questions/4585283/alternate-showing-hiding-window-when-notify-icon-is-clicked
	Debounce time.Duration

	// Indicates that the window has just been closed
	justClosed bool

	// Indicates that the window has been shown a first time
	hasBeenShown bool

	// Used to ensure that the window state is read on first click
	initialClick sync.Once
}

// AttachWindow attaches a window to the system tray. The window will be shown when the system tray icon is clicked.
// The window will be hidden when the system tray icon is clicked again, or when the window loses focus.
func (s *SystemTray) AttachWindow(window *WebviewWindow) *SystemTray {
	s.attachedWindow.Window = window
	return s
}

// WindowOffset sets the gap in pixels between the system tray and the window
func (s *SystemTray) WindowOffset(offset int) *SystemTray {
	s.attachedWindow.Offset = offset
	return s
}

// WindowDebounce is used by Windows to indicate how long to wait before responding to a mouse
// up event on the notification icon. This prevents the window from being hidden and then immediately
// shown when the user clicks on the system tray icon.
// See https://stackoverflow.com/questions/4585283/alternate-showing-hiding-window-when-notify-icon-is-clicked
func (s *SystemTray) WindowDebounce(debounce time.Duration) *SystemTray {
	s.attachedWindow.Debounce = debounce
	return s
}

func (s *SystemTray) defaultClickHandler() {
	if s.attachedWindow.Window == nil {
		s.OpenMenu()
		return
	}

	// Check the initial visibility state
	s.attachedWindow.initialClick.Do(func() {
		s.attachedWindow.hasBeenShown = s.attachedWindow.Window.IsVisible()
	})

	if runtime.GOOS == "windows" && s.attachedWindow.justClosed {
		return
	}

	if s.attachedWindow.Window.IsVisible() {
		s.attachedWindow.Window.Hide()
	} else {
		s.attachedWindow.hasBeenShown = true
		_ = s.PositionWindow(s.attachedWindow.Window, s.attachedWindow.Offset)
		s.attachedWindow.Window.Show().Focus()
	}
}

func (s *SystemTray) OpenMenu() {
	if s.menu == nil {
		return
	}
	if s.impl == nil {
		return
	}
	InvokeSync(s.impl.openMenu)
}
