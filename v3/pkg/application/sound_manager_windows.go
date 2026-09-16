//go:build windows && !server

package application

import (
	"syscall"
	"unsafe"
)

var (
	soundModUser32        = syscall.NewLazyDLL("user32.dll")
	soundModWinmm         = syscall.NewLazyDLL("winmm.dll")
	soundProcMessageBeep  = soundModUser32.NewProc("MessageBeep")
	soundProcPlaySoundW   = soundModWinmm.NewProc("PlaySoundW")
	soundMessageBeepAlert = uintptr(0xFFFFFFFF) // MB_OK is 0; -1 plays the default beep
)

const (
	soundSndAsync     = 0x0001
	soundSndNoDefault = 0x0002
	soundSndAlias     = 0x00010000
	soundSndFilename  = 0x00020000
)

func systemSoundsDir() string {
	return ""
}

func soundBeep() {
	_, _, _ = soundProcMessageBeep.Call(soundMessageBeepAlert)
}

func soundPlayWithFlags(name string, flags uintptr) error {
	ptr, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return err
	}
	ret, _, _ := soundProcPlaySoundW.Call(uintptr(unsafe.Pointer(ptr)), 0, flags)
	if ret == 0 {
		return ErrSoundNotSupported
	}
	return nil
}

func soundPlayNamed(name string) error {
	return soundPlayWithFlags(name, soundSndAsync|soundSndAlias|soundSndNoDefault)
}

func soundPlayFile(path string) error {
	return soundPlayWithFlags(path, soundSndAsync|soundSndFilename|soundSndNoDefault)
}

func soundPlayData([]byte) error {
	return ErrSoundNotSupported
}
