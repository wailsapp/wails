//go:build !linux && !windows && !darwin

package hostinfo

import "errors"

func Processors() ([]Processor, error) {
	return nil, errors.New("processor information unavailable on this platform")
}

func Memory() (uint64, error) {
	return 0, errors.New("physical memory information unavailable on this platform")
}

func Graphics() ([]GraphicsCard, error) {
	return nil, errors.New("graphics information unavailable on this platform")
}
