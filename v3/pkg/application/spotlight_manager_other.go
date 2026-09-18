//go:build !darwin || ios || server

package application

func spotlightIsAvailable() bool {
	return false
}

func spotlightIndex(string) error {
	return ErrSpotlightNotSupported
}

func spotlightDelete([]string) error {
	return ErrSpotlightNotSupported
}

func spotlightDeleteDomain(string) error {
	return ErrSpotlightNotSupported
}

func spotlightDeleteAll() error {
	return ErrSpotlightNotSupported
}
