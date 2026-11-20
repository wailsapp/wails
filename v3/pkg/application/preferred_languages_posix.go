//go:build server || !(windows || darwin || android)

package application

import "os"

// platformPreferredLanguages reads the POSIX locale environment. On a Linux
// desktop the session exports the user's choice into the environment that
// launched apps inherit. In server builds, on any OS, this is the host's
// setting, not the remote viewer's: an app serving browsers should use the
// request's Accept-Language instead.
func platformPreferredLanguages() []string {
	return posixPreferredLanguages(os.Getenv)
}
