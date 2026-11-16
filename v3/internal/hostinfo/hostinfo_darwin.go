package hostinfo

import (
	"errors"

	"golang.org/x/sys/unix"
)

func Processors() ([]Processor, error) {
	model, err := unix.Sysctl("machdep.cpu.brand_string")
	if err != nil {
		return nil, err
	}
	cores, err := unix.SysctlUint32("hw.physicalcpu")
	if err != nil {
		return nil, err
	}
	return []Processor{{Model: model, Cores: int(cores)}}, nil
}

func Memory() (uint64, error) {
	return unix.SysctlUint64("hw.memsize")
}

func Graphics() ([]GraphicsCard, error) {
	return nil, errors.New("use the doctor's macOS graphics collector")
}
