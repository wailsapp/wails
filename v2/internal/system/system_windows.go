// +build windows

package system

import "github.com/wailsapp/wails/v2/internal/system/operatingsystem"

func (i *Info) discover() error {

	var err error
	osinfo, err := operatingsystem.Info()
	if err != nil {
		return err
	}
	i.OS = osinfo

	// dll := syscall.MustLoadDLL("kernel32.dll")
	// p := dll.MustFindProc("GetVersion")
	// v, _, _ := p.Call()
	// fmt.Printf("Windows version %d.%d (Build %d)\n", byte(v), uint8(v>>8), uint16(v>>16))
	return nil
}
