package commands

import (
	"context"
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

var manifestReadConsoleInput = windows.NewLazySystemDLL("kernel32.dll").NewProc("ReadConsoleInputW")

// INPUT_RECORD/KEY_EVENT_RECORD follow the Win32 console ABI. Reading input
// events after a bounded wait avoids a blocking line read during cancellation.
type manifestConsoleEvent struct {
	Type                                    uint16
	_                                       uint16
	KeyDown                                 int32
	Repeat, VirtualKey, ScanCode, Character uint16
	Control                                 uint32
}

func readManifestDevInput(ctx context.Context, input *os.File, send func(string) bool) error {
	var line []rune
	handle := windows.Handle(input.Fd())
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		result, err := windows.WaitForSingleObject(handle, 100)
		if err != nil {
			return err
		}
		if result == uint32(windows.WAIT_TIMEOUT) {
			continue
		}
		if result != windows.WAIT_OBJECT_0 {
			return fmt.Errorf("console wait result %d", result)
		}
		var event manifestConsoleEvent
		var count uint32
		ok, _, err := manifestReadConsoleInput.Call(uintptr(handle), uintptr(unsafe.Pointer(&event)), 1, uintptr(unsafe.Pointer(&count)))
		if ok == 0 {
			return err
		}
		if count == 0 || event.Type != 1 || event.KeyDown == 0 || event.Character == 0 {
			continue
		}
		switch event.Character {
		case '\r':
			fmt.Fprint(os.Stdout, "\n")
			if !send(string(line)) {
				return ctx.Err()
			}
			line = nil
		case '\b':
			if len(line) > 0 {
				line = line[:len(line)-1]
				fmt.Fprint(os.Stdout, "\b \b")
			}
		default:
			if event.Character >= 32 && len(line) < 4096 {
				line = append(line, rune(event.Character))
				fmt.Fprintf(os.Stdout, "%c", event.Character)
			}
		}
	}
}
