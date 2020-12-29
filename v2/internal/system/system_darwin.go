// +build darwin

package system

import "github.com/wailsapp/wails/v2/internal/system/operatingsystem"

func (i *Info) discover() error {
	var err error
	osinfo, err := operatingsystem.Info()
	if err != nil {
		return err
	}
	i.OS = osinfo
	return nil
}
