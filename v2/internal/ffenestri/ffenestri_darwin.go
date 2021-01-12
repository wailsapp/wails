package ffenestri

/*
#cgo darwin CFLAGS: -DFFENESTRI_DARWIN=1
#cgo darwin LDFLAGS: -framework WebKit -lobjc

#include "ffenestri.h"
#include "ffenestri_darwin.h"

extern void TitlebarAppearsTransparent(void *);
extern void HideTitle(void *);
extern void HideTitleBar(void *);
extern void FullSizeContent(void *);
extern void UseToolbar(void *);
extern void HideToolbarSeparator(void *);
extern void DisableFrame(void *);
extern void SetAppearance(void *, const char *);
extern void WebviewIsTransparent(void *);
extern void WindowBackgroundIsTranslucent(void *);
extern void SetTray(void *, const char *, const char *, const char *);
extern void SetContextMenus(void *, const char *);

*/
import "C"
import (
	"github.com/wailsapp/wails/v2/pkg/options"
	"os"
)

func (a *Application) processPlatformSettings() error {

	mac := a.config.Mac
	titlebar := mac.TitleBar

	// HideTitle
	if titlebar.HideTitle {
		C.HideTitle(a.app)
	}

	// HideTitleBar
	if titlebar.HideTitleBar {
		C.HideTitleBar(a.app)
	}

	// Full Size Content
	if titlebar.FullSizeContent {
		C.FullSizeContent(a.app)
	}

	// Toolbar
	if titlebar.UseToolbar {
		C.UseToolbar(a.app)
	}

	if titlebar.HideToolbarSeparator {
		C.HideToolbarSeparator(a.app)
	}

	if titlebar.TitlebarAppearsTransparent {
		C.TitlebarAppearsTransparent(a.app)
	}

	// Process window Appearance
	if mac.Appearance != "" {
		C.SetAppearance(a.app, a.string2CString(string(mac.Appearance)))
	}

	// Check if the webview should be transparent
	if mac.WebviewIsTransparent {
		C.WebviewIsTransparent(a.app)
	}

	// Check if window should be translucent
	if mac.WindowBackgroundIsTranslucent {
		C.WindowBackgroundIsTranslucent(a.app)
	}

	// Process menu
	//applicationMenu := options.GetApplicationMenu(a.config)
	applicationMenu := a.menuManager.GetApplicationMenuJSON()
	if applicationMenu != "" {
		C.SetApplicationMenu(a.app, a.string2CString(applicationMenu))
	}

	// Process tray
	trays, err := a.menuManager.GetTrayMenus()
	if err != nil {
		return err
	}
	if trays != nil {
		for _, tray := range trays {
			println("Adding tray menu: " + tray)
			//C.AddTray(a.app, a.string2CString(tray))
		}
	}
	os.Exit(1)

	// Process context menus
	contextMenus := options.GetContextMenus(a.config)
	if contextMenus != nil {
		contextMenusJSON, err := processContextMenus(contextMenus)
		if err != nil {
			return err
		}
		C.SetContextMenus(a.app, a.string2CString(contextMenusJSON))
	}

	return nil
}
