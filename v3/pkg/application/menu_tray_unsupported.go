//go:build server || (openharmony && cgo)

package application

import "errors"

// This backend does not expose desktop application menus or a system tray.
func newMenuImpl(_ *Menu) menuImpl { return unsupportedMenu{} }

type unsupportedMenu struct{}

func (unsupportedMenu) update() {}

func newSystemTrayImpl(_ *SystemTray) systemTrayImpl { return unsupportedSystemTray{} }

var errSystemTrayUnsupported = errors.New("system tray not available on this backend")

type unsupportedSystemTray struct{}

func (unsupportedSystemTray) setLabel(string)                  {}
func (unsupportedSystemTray) setTooltip(string)                {}
func (unsupportedSystemTray) run()                             {}
func (unsupportedSystemTray) setIcon([]byte)                   {}
func (unsupportedSystemTray) setMenu(*Menu)                    {}
func (unsupportedSystemTray) setIconPosition(IconPosition)     {}
func (unsupportedSystemTray) setTemplateIcon([]byte)           {}
func (unsupportedSystemTray) destroy()                         {}
func (unsupportedSystemTray) setDarkModeIcon([]byte)           {}
func (unsupportedSystemTray) bounds() (*Rect, error)           { return nil, errSystemTrayUnsupported }
func (unsupportedSystemTray) getScreen() (*Screen, error)      { return nil, errSystemTrayUnsupported }
func (unsupportedSystemTray) positionWindow(Window, int) error { return errSystemTrayUnsupported }
func (unsupportedSystemTray) openMenu()                        {}
func (unsupportedSystemTray) Show()                            {}
func (unsupportedSystemTray) Hide()                            {}
