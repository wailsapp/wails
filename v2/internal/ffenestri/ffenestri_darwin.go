package ffenestri

/*
#cgo darwin CFLAGS: -DFFENESTRI_DARWIN=1
#cgo darwin LDFLAGS: -framework WebKit -lobjc

extern void TitlebarAppearsTransparent(void *);
extern void HideTitle(void *);
extern void HideTitleBar(void *);
extern void FullSizeContent(void *);
extern void UseToolbar(void *);
extern void HideToolbarSeparator(void *);
extern void DisableFrame(void *);
extern void SetAppearance(void *, const char *);
extern void WebviewIsTransparent(void *);
extern void SetWindowBackgroundIsTranslucent(void *);
extern void SetMenu(void *, const char *);
extern void SetTray(void *, const char *);
extern void SetContextMenus(void *, const char *);
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/options"
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
		C.SetWindowBackgroundIsTranslucent(a.app)
	}

	// Process menu
	applicationMenu := options.GetApplicationMenu(a.config)
	if applicationMenu != nil {

		/*
			As radio groups need to be manually managed on OSX,
			we preprocess the menu to determine the radio groups.
			This is defined as any adjacent menu item of type "RadioType".
			We keep a record of every radio group member we discover by saving
			a list of all members of the group and the number of members
			in the group (this last one is for optimisation at the C layer).
		*/
		processedMenu := NewProcessedMenu(applicationMenu)
		applicationMenuJSON, err := json.Marshal(processedMenu)
		if err != nil {
			return err
		}
		C.SetMenu(a.app, a.string2CString(string(applicationMenuJSON)))
	}

	// Process tray
	tray := options.GetTrayMenu(a.config)
	if tray != nil {

		/*
			As radio groups need to be manually managed on OSX,
			we preprocess the menu to determine the radio groups.
			This is defined as any adjacent menu item of type "RadioType".
			We keep a record of every radio group member we discover by saving
			a list of all members of the group and the number of members
			in the group (this last one is for optimisation at the C layer).
		*/
		processedMenu := NewProcessedMenu(tray)
		trayMenuJSON, err := json.Marshal(processedMenu)
		if err != nil {
			return err
		}
		C.SetTray(a.app, a.string2CString(string(trayMenuJSON)))
	}

	// Process context menus
	contextMenus := options.GetContextMenus(a.config)
	if contextMenus != nil {
		contextMenusJSON, err := json.Marshal(contextMenus)
		fmt.Printf("\n\nCONTEXT MENUS:\n %+v\n\n", string(contextMenusJSON))
		if err != nil {
			return err
		}
		C.SetContextMenus(a.app, a.string2CString(string(contextMenusJSON)))
	}

	return nil
}
