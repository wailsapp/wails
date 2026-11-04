// Package credentialstore stores signing credentials in the operating system's vault.
package credentialstore

import "errors"

var ErrNotFound = errors.New("credential not found")
var ErrDismissed = errors.New("credential access dismissed")
