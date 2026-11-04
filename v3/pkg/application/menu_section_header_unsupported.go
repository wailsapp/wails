//go:build !darwin || ios || server

package application

func menuSectionHeadersSupported() bool { return false }
