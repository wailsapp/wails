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
*/
import "C"

func (a *Application) processPlatformSettings() {

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

	// For macs we consider "frameless" to mean a combination of options
	if titlebar.TitlebarAppearsTransparent && titlebar.HideTitle {
		C.DisableFrame(a.app)
	}

	// Process window Appearance
	if mac.Appearance != "" {
		C.SetAppearance(a.app, a.string2CString(string(mac.Appearance)))
	}

	// Check if the webview should be transparent
	if mac.WebviewIsTransparent {
		C.WebviewIsTransparent(a.app)
	}
}
