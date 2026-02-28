package runtime

import (
	"context"
	"testing"

	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
)

type mockFrontend struct {
	frontend.Frontend
	trayMenu *menu.TrayMenu
}

func (m *mockFrontend) TraySetSystemTray(trayMenu *menu.TrayMenu) {
	m.trayMenu = trayMenu
}

func (m *mockFrontend) Run(ctx context.Context) error { return nil }
func (m *mockFrontend) RunMainLoop()                  {}
func (m *mockFrontend) ExecJS(js string)              {}
func (m *mockFrontend) Hide()                         {}
func (m *mockFrontend) Show()                         {}
func (m *mockFrontend) Quit()                         {}

func (m *mockFrontend) OpenFileDialog(dialogOptions frontend.OpenDialogOptions) (string, error) {
	return "", nil
}
func (m *mockFrontend) OpenMultipleFilesDialog(dialogOptions frontend.OpenDialogOptions) ([]string, error) {
	return nil, nil
}
func (m *mockFrontend) OpenDirectoryDialog(dialogOptions frontend.OpenDialogOptions) (string, error) {
	return "", nil
}
func (m *mockFrontend) SaveFileDialog(dialogOptions frontend.SaveDialogOptions) (string, error) {
	return "", nil
}
func (m *mockFrontend) MessageDialog(dialogOptions frontend.MessageDialogOptions) (string, error) {
	return "", nil
}

func (m *mockFrontend) WindowSetTitle(title string)                 {}
func (m *mockFrontend) WindowShow()                                 {}
func (m *mockFrontend) WindowHide()                                 {}
func (m *mockFrontend) WindowCenter()                               {}
func (m *mockFrontend) WindowToggleMaximise()                       {}
func (m *mockFrontend) WindowMaximise()                             {}
func (m *mockFrontend) WindowUnmaximise()                           {}
func (m *mockFrontend) WindowMinimise()                             {}
func (m *mockFrontend) WindowUnminimise()                           {}
func (m *mockFrontend) WindowSetAlwaysOnTop(b bool)                 {}
func (m *mockFrontend) WindowSetPosition(x int, y int)              {}
func (m *mockFrontend) WindowGetPosition() (int, int)               { return 0, 0 }
func (m *mockFrontend) WindowSetSize(width int, height int)         {}
func (m *mockFrontend) WindowGetSize() (int, int)                   { return 0, 0 }
func (m *mockFrontend) WindowSetMinSize(width int, height int)      {}
func (m *mockFrontend) WindowSetMaxSize(width int, height int)      {}
func (m *mockFrontend) WindowFullscreen()                           {}
func (m *mockFrontend) WindowUnfullscreen()                         {}
func (m *mockFrontend) WindowSetBackgroundColour(col *options.RGBA) {}
func (m *mockFrontend) WindowReload()                               {}
func (m *mockFrontend) WindowReloadApp()                            {}
func (m *mockFrontend) WindowSetSystemDefaultTheme()                {}
func (m *mockFrontend) WindowSetLightTheme()                        {}
func (m *mockFrontend) WindowSetDarkTheme()                         {}
func (m *mockFrontend) WindowIsMaximised() bool                     { return false }
func (m *mockFrontend) WindowIsMinimised() bool                     { return false }
func (m *mockFrontend) WindowIsNormal() bool                        { return false }
func (m *mockFrontend) WindowIsFullscreen() bool                    { return false }
func (m *mockFrontend) WindowClose()                                {}
func (m *mockFrontend) WindowPrint()                                {}
func (m *mockFrontend) ScreenGetAll() ([]frontend.Screen, error)    { return nil, nil }
func (m *mockFrontend) MenuSetApplicationMenu(menu *menu.Menu)      {}
func (m *mockFrontend) MenuUpdateApplicationMenu()                  {}
func (m *mockFrontend) Notify(name string, data ...interface{})     {}
func (m *mockFrontend) BrowserOpenURL(url string)                   {}
func (m *mockFrontend) ClipboardGetText() (string, error)           { return "", nil }
func (m *mockFrontend) ClipboardSetText(text string) error          { return nil }

func TestTraySetSystemTray(t *testing.T) {
	mock := &mockFrontend{}
	ctx := context.WithValue(context.Background(), "frontend", mock)

	trayMenu := &menu.TrayMenu{
		Label: "Test",
	}

	TraySetSystemTray(ctx, trayMenu)

	if mock.trayMenu != trayMenu {
		t.Errorf("Expected tray menu to be set")
	}
}
