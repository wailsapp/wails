//go:build linux

package start_at_login

func (p *Plugin) init() error {
	// TBD
	return nil
}

func (p *Plugin) StartAtLogin(enabled bool) error {
	panic("not implemented")
}

func (p *Plugin) IsStartAtLogin() (bool, error) {
	panic("not implemented")
}
