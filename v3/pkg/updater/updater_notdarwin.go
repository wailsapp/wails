//go:build !darwin

package updater

func bundleTarget(exe string) string { return exe }
