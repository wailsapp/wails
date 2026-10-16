//go:build windows

package application

import (
	"fmt"

	"github.com/wailsapp/wails/v3/pkg/w32"
)

// NewIconFromResource loads an icon from an embedded Windows resource. It is
// available in both desktop and server builds so that services which only need
// the icon helper (e.g. notifications) compile under the `server` tag, matching
// the behaviour on macOS and Linux.
func NewIconFromResource(instance w32.HINSTANCE, resId uint16) (w32.HICON, error) {
	var err error
	var result w32.HICON
	if result = w32.LoadIconWithResourceID(instance, resId); result == 0 {
		err = fmt.Errorf("cannot load icon from resource with id %v", resId)
	}
	return result, err
}
