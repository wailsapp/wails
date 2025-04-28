//go:build linux
// +build linux

package linux

func (f *Frontend) IsNotificationAvailable() bool {
	return true
}

func (f *Frontend) RequestNotificationAuthorization() (bool, error) {
	return true, nil
}

func (f *Frontend) CheckNotificationAuthorization() (bool, error) {
	return true, nil
}
