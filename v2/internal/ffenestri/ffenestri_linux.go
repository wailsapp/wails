package ffenestri

/*
#cgo linux CFLAGS: -DFFENESTRI_LINUX=1
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0


#include "ffenestri.h"
#include "ffenestri_linux.h"

*/
import "C"

func (a *Application) processPlatformSettings() error {

	return nil
}
