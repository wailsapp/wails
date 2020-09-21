package ffenestri

/*
#cgo darwin CFLAGS: -DFFENESTRI_DARWIN=1
#cgo darwin LDFLAGS: -framework WebKit -lobjc

extern void TitlebarAppearsTransparent(void *);
extern void HideTitle(void *);
*/
import "C"

func (a *Application) processPlatformSettings() {

	// HideTitle
	if a.config.Mac.HideTitle {
		C.HideTitle(a.app)
	}

	// if a.config.Mac.TitlebarAppearsTransparent {
	// 	C.TitlebarAppearsTransparent(a.app)
	// }
}
