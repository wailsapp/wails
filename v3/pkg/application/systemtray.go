package application

import (
	"errors"
	"runtime"
	"strings"
	"sync"
	"time"
)

// SystemTrayLabelPart is a segment of a system tray label with optional colour information.
// FgColor and BgColor are hex strings (e.g. "#ff0000") or empty strings when not set.
type SystemTrayLabelPart struct {
	Text    string
	FgColor string
	BgColor string
}

// SystemTrayLabelParser is a function that splits a raw label string into styled parts.
// It is called by the macOS system tray implementation when the label contains ANSI
// escape codes. The default implementation strips escape codes and returns plain text.
//
// Replace this variable with a custom function to render colours in your system tray
// labels. For example, wire in go-ansi-parser or any other ANSI library here.
var SystemTrayLabelParser func(label string) ([]SystemTrayLabelPart, error) = stripANSI

// stripANSI removes ANSI escape sequences from s and returns a single unstyled part.
func stripANSI(s string) ([]SystemTrayLabelPart, error) {
	// Remove all sequences matching \033[...m
	var buf strings.Builder
	for i := 0; i < len(s); {
		if s[i] == '\033' && i+1 < len(s) && s[i+1] == '[' {
			end := strings.IndexByte(s[i:], 'm')
			if end == -1 {
				buf.WriteByte(s[i])
				i++
				continue
			}
			i += end + 1
			continue
		}
		buf.WriteByte(s[i])
		i++
	}
	return []SystemTrayLabelPart{{Text: buf.String()}}, nil
}

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
	setTooltip(tooltip string)
	run()
	setIcon(icon []byte)
	setMenu(menu *Menu)
	setIconPosition(position IconPosition)
	setTemplateIcon(icon []byte)
	destroy()
	setDarkModeIcon(icon []byte)
	bounds() (*Rect, error)
	getScreen() (*Screen, error)
	positionWindow(window Window, offset int) error
	openMenu()
	Show()
	Hide()
}

// systemTrayDarwinExtras is implemented by the macOS system tray only. The
// cross-platform SystemTray methods type-assert it so systemTrayImpl stays
// stable for the other platforms.
type systemTrayDarwinExtras interface {
	setSymbol(name string, pointSize float64, weight MacSymbolWeight)
	setTooltip(tooltip string)
	setRemovable(allowed bool, autosaveName string)
	isVisible() bool
}

// MacSymbolWeight is the weight applied to an SF Symbol rendered in the
// status bar (NSFontWeight). MacSymbolWeightUnspecified keeps the default.
type MacSymbolWeight int

const (
	MacSymbolWeightUnspecified MacSymbolWeight = iota
	MacSymbolWeightUltraLight
	MacSymbolWeightThin
	MacSymbolWeightLight
	MacSymbolWeightRegular
	MacSymbolWeightMedium
	MacSymbolWeightSemibold
	MacSymbolWeightBold
	MacSymbolWeightHeavy
	MacSymbolWeightBlack
)

// String returns the AppKit name of the weight.
func (w MacSymbolWeight) String() string {
	switch w {
	case MacSymbolWeightUltraLight:
		return "ultraLight"
	case MacSymbolWeightThin:
		return "thin"
	case MacSymbolWeightLight:
		return "light"
	case MacSymbolWeightRegular:
		return "regular"
	case MacSymbolWeightMedium:
		return "medium"
	case MacSymbolWeightSemibold:
		return "semibold"
	case MacSymbolWeightBold:
		return "bold"
	case MacSymbolWeightHeavy:
		return "heavy"
	case MacSymbolWeightBlack:
		return "black"
	default:
		return "unspecified"
	}
}

type SystemTray struct {
	id           uint
	label        string
	tooltip      string
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

	// macOS extras (see SetSymbol, SetRemovable). The model is kept here so
	// values set before Run are applied when the native item is created and
	// so the state is observable on platforms without native support.
	symbolName        string
	symbolPointSize   float64
	symbolWeight      MacSymbolWeight
	removable         bool
	autosaveName      string
	hidden            bool
	visibilityHandler func(visible bool)
}

func newSystemTray(id uint) *SystemTray {
	result := &SystemTray{
		id:           id,
		label:        "",
		tooltip:      "",
		iconPosition: NSImageLeading,
		attachedWindow: WindowAttachConfig{
			Window:   nil,
			Offset:   0,
			Debounce: 200 * time.Millisecond,
		},
	}
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
	if globalApplication == nil || globalApplication.running == false {
		return
	}

	s.applySmartDefaults()
	s.impl = newSystemTrayImpl(s)
	InvokeSync(s.impl.run)
}

func (s *SystemTray) applySmartDefaults() {
	hasWindow := s.attachedWindow.Window != nil
	hasMenu := s.menu != nil

	if s.clickHandler == nil && hasWindow {
		s.clickHandler = s.ToggleWindow
	}

	if s.rightClickHandler == nil && hasMenu {
		s.rightClickHandler = s.ShowMenu
	}
}

func (s *SystemTray) ToggleWindow() {
	if s.attachedWindow.Window == nil {
		return
	}

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

func (s *SystemTray) ShowMenu() {
	s.OpenMenu()
}

func (s *SystemTray) ShowWindow() {
	if s.attachedWindow.Window == nil {
		return
	}
	s.attachedWindow.hasBeenShown = true
	_ = s.PositionWindow(s.attachedWindow.Window, s.attachedWindow.Offset)
	s.attachedWindow.Window.Show().Focus()
}

func (s *SystemTray) HideWindow() {
	if s.attachedWindow.Window == nil {
		return
	}
	s.attachedWindow.Window.Hide()
}

func (s *SystemTray) PositionWindow(window Window, offset int) error {
	if s.impl == nil {
		return errors.New("system tray not running")
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

func (s *SystemTray) SetTooltip(tooltip string) {
	s.tooltip = tooltip
	if s.impl == nil {
		return
	}
	InvokeSync(func() {
		s.impl.setTooltip(tooltip)
	})
}

func (s *SystemTray) Destroy() {
	globalApplication.SystemTray.destroy(s)
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
	s.hidden = false
	if s.impl == nil {
		return
	}
	InvokeSync(func() {
		s.impl.Show()
	})
}

func (s *SystemTray) Hide() {
	s.hidden = true
	if s.impl == nil {
		return
	}
	InvokeSync(func() {
		s.impl.Hide()
	})
}

// SetVisible shows or hides the tray item. It is equivalent to Show and Hide.
func (s *SystemTray) SetVisible(visible bool) {
	if visible {
		s.Show()
		return
	}
	s.Hide()
}

// IsVisible reports whether the tray item is currently visible.
//
// On macOS this reads NSStatusItem.visible, so it is false after the user
// drags a removable item out of the menu bar (see SetRemovable) or when a
// removable item with an autosave name was removed in a previous session.
// On other platforms it reflects the last Show/Hide call.
func (s *SystemTray) IsVisible() bool {
	if s.impl == nil {
		return !s.hidden
	}
	extras, ok := s.impl.(systemTrayDarwinExtras)
	if !ok {
		return !s.hidden
	}
	return InvokeSyncWithResult(extras.isVisible)
}

// SetSymbol sets the tray icon to the SF Symbol with the given name, for
// example "star.fill" or "cloud.sun". Symbols are rendered as template
// images so they follow the menu bar appearance.
//
// macOS 11 or later only; it is ignored elsewhere and on older versions. The
// symbol coexists with SetIcon, SetTemplateIcon and SetDarkModeIcon: the last
// call wins. Optional sizing is applied with SetSymbolConfiguration.
func (s *SystemTray) SetSymbol(name string) *SystemTray {
	s.symbolName = name
	if s.impl == nil {
		return s
	}
	if extras, ok := s.impl.(systemTrayDarwinExtras); ok {
		InvokeSync(func() {
			extras.setSymbol(name, s.symbolPointSize, s.symbolWeight)
		})
	}
	return s
}

// Symbol returns the SF Symbol name set with SetSymbol, or "" when none is set.
func (s *SystemTray) Symbol() string {
	return s.symbolName
}

// SetSymbolConfiguration sets the point size and weight used to render the
// SF Symbol set with SetSymbol. A pointSize of 0 keeps the status bar default
// and MacSymbolWeightUnspecified keeps the default weight.
//
// macOS 11 or later only; ignored elsewhere.
func (s *SystemTray) SetSymbolConfiguration(pointSize float64, weight MacSymbolWeight) *SystemTray {
	if pointSize < 0 {
		pointSize = 0
	}
	s.symbolPointSize = pointSize
	s.symbolWeight = weight
	if s.impl == nil || s.symbolName == "" {
		return s
	}
	if extras, ok := s.impl.(systemTrayDarwinExtras); ok {
		InvokeSync(func() {
			extras.setSymbol(s.symbolName, pointSize, weight)
		})
	}
	return s
}

// SetRemovable controls whether the user may drag the tray item out of the
// menu bar (Command-drag), the way built-in status items work. autosaveName
// is the key under which macOS remembers the item's position and whether it
// was removed; it must be a non-empty, stable string when allowed is true so
// the item stays removed across launches. Pass an empty name to keep the
// default (an automatic name is used, which does not survive relaunch).
//
// A removed item can be brought back with Show or SetVisible(true). Use
// OnVisibilityChange to react when the user removes it.
//
// macOS only; the flag is stored but has no effect elsewhere.
func (s *SystemTray) SetRemovable(allowed bool, autosaveName string) *SystemTray {
	s.removable = allowed
	s.autosaveName = autosaveName
	if s.impl == nil {
		return s
	}
	if extras, ok := s.impl.(systemTrayDarwinExtras); ok {
		InvokeSync(func() {
			extras.setRemovable(allowed, autosaveName)
		})
	}
	return s
}

// IsRemovable reports whether SetRemovable allowed user removal.
func (s *SystemTray) IsRemovable() bool {
	return s.removable
}

// AutosaveName returns the autosave name set with SetRemovable.
func (s *SystemTray) AutosaveName() string {
	return s.autosaveName
}

// OnVisibilityChange registers a handler called when the tray item's
// visibility changes, including when the user removes a removable item from
// the menu bar or the item is shown or hidden programmatically.
//
// macOS only; the handler is never called elsewhere.
func (s *SystemTray) OnVisibilityChange(handler func(visible bool)) *SystemTray {
	s.visibilityHandler = handler
	return s
}

// Tooltip returns the tooltip set with SetTooltip.
func (s *SystemTray) Tooltip() string {
	return s.tooltip
}

type WindowAttachConfig struct {
	// Window is the window to attach to the system tray. If it's null, the request to attach will be ignored.
	Window Window

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
func (s *SystemTray) AttachWindow(window Window) *SystemTray {
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

func (s *SystemTray) OpenMenu() {
	if s.menu == nil {
		return
	}
	if s.impl == nil {
		return
	}
	InvokeSync(s.impl.openMenu)
}
