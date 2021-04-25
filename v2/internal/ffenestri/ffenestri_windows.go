package ffenestri

/*

#cgo windows,amd64 LDFLAGS: -L./windows/x64 -lwebview -lWebView2Loader -lgdi32

#include "ffenestri.h"
#include "ffenestri_windows.h"

*/
import "C"

func (a *Application) processPlatformSettings() error {

	return nil
}
