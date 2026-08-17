//go:build darwin && !ios && !server && wails_native

package application

/*
#cgo CFLAGS: -DWAILS_NATIVE_ONLY=1
*/
import "C"
