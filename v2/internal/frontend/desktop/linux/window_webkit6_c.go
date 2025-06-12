//go:build linux && webkit_6
// +build linux,webkit_6

package linux

/*
extern void onActivate();
*/
import "C"

//export onActivate
func onActivate() {
	activateWg.Done()
}
