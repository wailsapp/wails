//go:build linux

package single_instance

/*
#cgo CFLAGS:
#cgo LDFLAGS:

*/
import "C"
import "fmt"

func (p *Plugin) activeInstance(pid int) error {
	//	C.activateApplicationWithProcessID(C.int(pid))
	fmt.Println("[linux] activateInstance - not implemented: ", pid)
	return nil
}
