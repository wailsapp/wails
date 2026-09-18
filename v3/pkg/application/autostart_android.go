//go:build android

package application

type unsupportedAutostart struct{}

func newAutostartImpl(_ *App) autostartImpl { return unsupportedAutostart{} }

func (unsupportedAutostart) enable(AutostartOptions) error  { return ErrAutostartNotSupported }
func (unsupportedAutostart) disable() error                 { return ErrAutostartNotSupported }
func (unsupportedAutostart) status() (AutostartStatus, error) {
	return AutostartStatus{}, ErrAutostartNotSupported
}
