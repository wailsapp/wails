// +build linux

package ffenestri

/*

#cgo linux CFLAGS: -DFFENESTRI_LINUX=1
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0

#include <stdlib.h>
#include "ffenestri.h"


*/
import "C"
import "github.com/wailsapp/wails/v2/internal/features"

func (a *Application) processOSFeatureFlags(features *features.Features) {

	// Process Linux features
	// linux := features.Linux

}
