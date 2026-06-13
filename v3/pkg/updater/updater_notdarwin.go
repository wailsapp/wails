//go:build !darwin || ios

package updater

func bundleTarget(exe string) string { return exe }
