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
*/
import "C"
import (
	"encoding/json"
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
	if mac.Menu != nil {

		/*
			As radio groups need to be manually managed on OSX,
			we preprocess the menu to determine the radio groups.
			This is defined as any adjacent menu item of type "RadioType".
			We keep a record of every radio group member we discover by saving
			a list of all members of the group and the number of members
			in the group (this last one is for optimisation at the C layer).

			Example:
			{
				"RadioGroups": [
					{
						"Members": [
							"option-1",
							"option-2",
							"option-3"
						],
						"Length": 3
					}
				]
			}
		*/
		processedMenu := NewProcessedMenu(mac.Menu)
		menuJSON, err := json.Marshal(processedMenu)
		if err != nil {
			return err
		}
		C.SetMenu(a.app, a.string2CString(string(menuJSON)))
	}

	return nil
}
