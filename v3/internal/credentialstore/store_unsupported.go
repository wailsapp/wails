//go:build !darwin && !windows && !linux && !freebsd && !openbsd && !netbsd && !dragonfly

package credentialstore

import (
	"fmt"
	"runtime"
)

func Set(service, user, value string) error    { return unsupported() }
func Get(service, user string) (string, error) { return "", unsupported() }
func Delete(service, user string) error        { return unsupported() }

func unsupported() error {
	return fmt.Errorf("credential storage is unavailable on %s", runtime.GOOS)
}
