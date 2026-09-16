//go:build !darwin || ios || server

package application

func appleEventsRegisterNative(*AppleEventsManager, string, string) error {
	return ErrAppleEventsNotSupported
}

func appleEventsSendNative(string, string, string, string) (string, error) {
	return "", ErrAppleEventsNotSupported
}
