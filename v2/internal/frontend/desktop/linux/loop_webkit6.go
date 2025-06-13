//go:build linux && webkit_6
// +build linux,webkit_6

package linux

/*
#cgo pkg-config: gtk4

#include "gtk/gtk.h"

extern void setMainLoop(GMainLoop *loop);
*/
import "C"

var mainLoop *C.GMainLoop

//export setMainLoop
func setMainLoop(loop *C.GMainLoop) {
	mainLoop = loop
}
