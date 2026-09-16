//go:build (!darwin && !windows) || ios || server

package application

func systemSoundsDir() string {
	return ""
}

func soundBeep() {}

func soundPlayNamed(string) error {
	return ErrSoundNotSupported
}

func soundPlayFile(string) error {
	return ErrSoundNotSupported
}

func soundPlayData([]byte) error {
	return ErrSoundNotSupported
}
