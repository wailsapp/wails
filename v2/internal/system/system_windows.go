// +build windows

package system

import (
	"github.com/wailsapp/wails/v2/internal/system/operatingsystem"
)

func (i *Info) discover() error {

	var err error
	osinfo, err := operatingsystem.Info()
	if err != nil {
		return err
	}
	i.OS = osinfo

	i.Dependencies = append(i.Dependencies, checkNPM())
	i.Dependencies = append(i.Dependencies, checkUPX())
	i.Dependencies = append(i.Dependencies, checkDocker())

	return nil
}

// IsAppleSilicon returns true if the app is running on Apple Silicon
// Credit: https://www.yellowduck.be/posts/detecting-apple-silicon-via-go/
// NOTE: Not applicable to windows
func IsAppleSilicon() bool {
	return false
}
