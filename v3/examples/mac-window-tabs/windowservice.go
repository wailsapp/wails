package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type WindowService struct{}

func (w *WindowService) OpenTabbedWindow() {
	w.openWindow("Tabbed Window", application.MacWindowTabbingModePreferred)
}

func (w *WindowService) OpenNonTabbedWindow() {
	w.openWindow("Non-Tabbed Window", application.MacWindowTabbingModeDisallowed)
}

// AddTab opens an Automatic-mode window and attaches it to the key window's
// tab group with WebviewWindow.AddTab. Automatic windows do not join a group
// on their own unless the user's "Prefer tabs" setting is "always", which
// makes the explicit call visible.
func (w *WindowService) AddTab() error {
	return addTabToCurrentWindow()
}

// DetachTab moves the key window's tab into its own window.
func (w *WindowService) DetachTab() {
	if window := currentWebviewWindow(); window != nil {
		window.MoveTabToNewWindow()
	}
}

// TabSummary describes the key window's tab group for the frontend.
func (w *WindowService) TabSummary() string {
	window := currentWebviewWindow()
	if window == nil {
		return "no key window"
	}
	group := window.TabGroup()
	if group == nil {
		return fmt.Sprintf("%q is not in a tab group", window.Name())
	}
	selected := "none"
	if selectedWindow := group.SelectedWindow(); selectedWindow != nil {
		selected = selectedWindow.Name()
	}
	return fmt.Sprintf("%d tab(s), selected %q, tab bar visible: %t",
		group.Count(), selected, group.IsTabBarVisible())
}

func (w *WindowService) openWindow(titlePrefix string, tabbingMode application.MacWindowTabbingMode) *application.WebviewWindow {
	app := application.Get()
	if app == nil {
		return nil
	}

	timestamp := time.Now().Format("15:04:05")
	windowTitle := fmt.Sprintf("%s (%s)", titlePrefix, timestamp)

	return app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: windowTitle,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
			TabbingMode:             tabbingMode,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})
}

// currentWebviewWindow returns the key window when it is a WebviewWindow.
func currentWebviewWindow() *application.WebviewWindow {
	app := application.Get()
	if app == nil {
		return nil
	}
	window, _ := app.Window.Current().(*application.WebviewWindow)
	return window
}

// currentTabGroup returns the key window's tab group. The result may be nil;
// MacWindowTabGroup methods are nil-safe so callers can chain directly.
func currentTabGroup() *application.MacWindowTabGroup {
	return currentWebviewWindow().TabGroup()
}

func addTabToCurrentWindow() error {
	window := currentWebviewWindow()
	if window == nil {
		return errors.New("no key window to add a tab to")
	}
	newWindow := (&WindowService{}).openWindow("Added Tab", application.MacWindowTabbingModeAutomatic)
	if newWindow == nil {
		return errors.New("could not open a window")
	}
	if err := window.AddTab(newWindow, application.MacTabOrderAbove); err != nil {
		return err
	}
	newWindow.SetTabTooltip("Added with WebviewWindow.AddTab")
	return nil
}

func addNativeTabToCurrentWindow() error {
	window := currentWebviewWindow()
	if window == nil {
		return errors.New("no key window to add a tab to")
	}
	app := application.Get()
	if app == nil {
		return errors.New("application is not running")
	}
	native := app.NativeWindow.NewWithOptions(application.NativeWindowOptions{
		Title: fmt.Sprintf("Native Tab (%s)", time.Now().Format("15:04:05")),
	})
	return window.AddNativeTab(native, application.MacTabOrderAbove)
}
