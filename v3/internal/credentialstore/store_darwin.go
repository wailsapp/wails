package credentialstore

import (
	"errors"

	"github.com/zalando/go-keyring"
)

// Keep the established keychain backend: it sends secrets over stdin and
// decodes credentials stored by earlier Wails releases.
func Set(service, user, value string) error { return keyring.Set(service, user, value) }
func Get(service, user string) (string, error) {
	value, err := keyring.Get(service, user)
	if errors.Is(err, keyring.ErrNotFound) {
		return "", ErrNotFound
	}
	return value, err
}
func Delete(service, user string) error {
	err := keyring.Delete(service, user)
	if errors.Is(err, keyring.ErrNotFound) {
		return ErrNotFound
	}
	return err
}
