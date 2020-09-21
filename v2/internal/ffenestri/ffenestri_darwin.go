package ffenestri

/*
#cgo darwin CFLAGS: -DFFENESTRI_DARWIN=1
#cgo darwin LDFLAGS: -framework WebKit -lobjc

extern void TitlebarAppearsTransparent(void *);
extern void HideTitle(void *);
extern void HideTitleBar(void *);
*/
import "C"

func (a *Application) processPlatformSettings() {

	// HideTitle
	if a.config.Mac.HideTitle {
		C.HideTitle(a.app)
	}

	// HideTitleBar
	if a.config.Mac.HideTitleBar {
		C.HideTitleBar(a.app)
	}

	// if a.config.Mac.TitlebarAppearsTransparent {
	// 	C.TitlebarAppearsTransparent(a.app)
	// }
}
