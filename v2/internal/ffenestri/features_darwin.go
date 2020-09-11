package ffenestri

/*

#cgo darwin CFLAGS: -DFFENESTRI_DARWIN=1
#cgo darwin LDFLAGS: -framework WebKit -lobjc

#include <stdlib.h>
#include "ffenestri.h"


*/
import "C"
import "github.com/wailsapp/wails/v2/internal/features"

func (a *Application) processOSFeatureFlags(features *features.Features) {

}
