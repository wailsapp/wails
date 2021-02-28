package ffenestri

/*
#cgo darwin CFLAGS: -DFFENESTRI_DARWIN=1
#cgo darwin LDFLAGS: -framework WebKit -framework CoreFoundation -lobjc

#include "ffenestri.h"
#include "ffenestri_darwin.h"

*/
import "C"

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

	// Set activation policy
	C.SetActivationPolicy(a.app, C.int(mac.ActivationPolicy))

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
			C.AddTrayMenu(a.app, a.string2CString(tray))
		}
	}

	// Process context menus
	contextMenus, err := a.menuManager.GetContextMenus()
	if err != nil {
		return err
	}
	if contextMenus != nil {
		for _, contextMenu := range contextMenus {
			C.AddContextMenu(a.app, a.string2CString(contextMenu))
		}
	}

	// Process URL Handlers
	if a.config.Mac.URLHandlers != nil {
		C.HasURLHandlers(a.app)
	}

	return nil
}
