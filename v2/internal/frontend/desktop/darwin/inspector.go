//go:build darwin && !(dev || debug || devtools)

package darwin

import (
	"unsafe"
)

func showInspector(_ unsafe.Pointer) {
}
