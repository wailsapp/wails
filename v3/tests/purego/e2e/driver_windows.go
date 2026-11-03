//go:build windows

package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// The Windows driver injects input with SendInput, as a keyboard and mouse
// do, and finds windows by title. It needs the interactive desktop: a
// process in a service or SSH session (session 0) has no windows to drive.
// The page's position in the window and the display scale come from
// calibrate (main.go), so the driver works in client-area pixels.

var (
	user32 = windows.NewLazySystemDLL("user32.dll")

	procFindWindowW         = user32.NewProc("FindWindowW")
	procSetForegroundWindow = user32.NewProc("SetForegroundWindow")
	procGetForegroundWindow = user32.NewProc("GetForegroundWindow")
	procShowWindow          = user32.NewProc("ShowWindow")
	procIsIconic            = user32.NewProc("IsIconic")
	procClientToScreen      = user32.NewProc("ClientToScreen")
	procSetCursorPos        = user32.NewProc("SetCursorPos")
	procSendInput           = user32.NewProc("SendInput")
	procVkKeyScanW          = user32.NewProc("VkKeyScanW")
	procAllowSetForeground  = user32.NewProc("AllowSetForegroundWindow")
)

const (
	inputMouse    = 0
	inputKeyboard = 1

	keyeventfExtendedKey = 0x0001
	keyeventfKeyUp       = 0x0002
	keyeventfUnicode     = 0x0004

	mouseeventfLeftDown  = 0x0002
	mouseeventfLeftUp    = 0x0004
	mouseeventfRightDown = 0x0008
	mouseeventfRightUp   = 0x0010

	swRestore = 9
	asfwAny   = ^uintptr(0) // ASFW_ANY

	vkShift   = 0x10
	vkControl = 0x11
	vkMenu    = 0x12 // Alt
	vkLWin    = 0x5B

	// eventDelay separates input events, a realistic gap between keystrokes.
	eventDelay = 60 * time.Millisecond
)

// input is the Win32 INPUT structure: a type and a union whose largest
// member is MOUSEINPUT. Declaring the union as that member gives it the
// right size and alignment on 32- and 64-bit Windows (28 and 40 bytes in
// all), which SendInput checks.
type input struct {
	typ uint32
	u   mouseInput
}

type keybdInput struct {
	vk, scan    uint16
	flags, time uint32
	extra       uintptr
}

type mouseInput struct {
	dx, dy            int32
	data, flags, time uint32
	extra             uintptr
}

func keyInput(vk, scan uint16, flags uint32) input {
	in := input{typ: inputKeyboard}
	*(*keybdInput)(unsafe.Pointer(&in.u)) = keybdInput{vk: vk, scan: scan, flags: flags}
	return in
}

func buttonInput(flags uint32) input {
	in := input{typ: inputMouse}
	*(*mouseInput)(unsafe.Pointer(&in.u)) = mouseInput{flags: flags}
	return in
}

func sendInput(in input) error {
	n, _, err := procSendInput.Call(1, uintptr(unsafe.Pointer(&in)), unsafe.Sizeof(in))
	time.Sleep(eventDelay)
	if n != 1 {
		return fmt.Errorf("SendInput: %w (is this the interactive desktop?)", err)
	}
	return nil
}

// extendedVK are the keys that live on the extended (non-numpad) block.
var extendedVK = map[uint16]bool{
	0x21: true, 0x22: true, 0x23: true, 0x24: true, // PageUp, PageDown, End, Home
	0x25: true, 0x26: true, 0x27: true, 0x28: true, // arrows
	0x2D: true, 0x2E: true, // Insert, Delete
	vkLWin: true,
}

func keyDown(vk uint16) error {
	var flags uint32
	if extendedVK[vk] {
		flags |= keyeventfExtendedKey
	}
	return sendInput(keyInput(vk, 0, flags))
}

func keyUp(vk uint16) error {
	flags := uint32(keyeventfKeyUp)
	if extendedVK[vk] {
		flags |= keyeventfExtendedKey
	}
	return sendInput(keyInput(vk, 0, flags))
}

// layoutShift is the shift state VkKeyScanW reports a character needs, in
// its result's high byte.
type layoutShift struct{ shift, ctrl, alt bool }

// vkFor returns the virtual-key code for a chord key name, from the
// layout-independent table or, for a printable character, from the active
// keyboard layout, with the modifiers the layout needs for that character
// ("/" is Shift+7 on a German keyboard).
func vkFor(key string) (uint16, layoutShift, error) {
	if vk, ok := windowsVirtualKey(key); ok {
		return vk, layoutShift{}, nil
	}
	if r := []rune(key); len(r) == 1 {
		res, _, _ := procVkKeyScanW.Call(uintptr(r[0]))
		if int16(res) != -1 {
			state := res >> 8
			return uint16(res & 0xFF), layoutShift{shift: state&1 != 0, ctrl: state&2 != 0, alt: state&4 != 0}, nil
		}
	}
	return 0, layoutShift{}, fmt.Errorf("no virtual key for %q on this keyboard layout", key)
}

type winDriver struct {
	mappings map[string]mapping
}

func newDriver() (driver, error) {
	// Let this process (and the dialogs it opens) take the foreground.
	procAllowSetForeground.Call(asfwAny)
	return &winDriver{mappings: map[string]mapping{}}, nil
}

func find(title string) (windows.HWND, error) {
	t, err := windows.UTF16PtrFromString(title)
	if err != nil {
		return 0, err
	}
	h, _, _ := procFindWindowW.Call(0, uintptr(unsafe.Pointer(t)))
	if h == 0 {
		return 0, fmt.Errorf("no window titled %q", title)
	}
	return windows.HWND(h), nil
}

// focus brings the window to the foreground. Windows gives the foreground
// only to the process that has the user's last input: an Alt press first
// makes that this process, as the documented SetForegroundWindow rules allow.
func (d *winDriver) focus(title string) error {
	h, err := find(title)
	if err != nil {
		return err
	}
	if iconic, _, _ := procIsIconic.Call(uintptr(h)); iconic != 0 {
		procShowWindow.Call(uintptr(h), swRestore)
	}
	for deadline := time.Now().Add(3 * time.Second); time.Now().Before(deadline); {
		if fg, _, _ := procGetForegroundWindow.Call(); windows.HWND(fg) == h {
			time.Sleep(eventDelay)
			return nil
		}
		_ = keyDown(vkMenu)
		_ = keyUp(vkMenu)
		procSetForegroundWindow.Call(uintptr(h))
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("could not bring %q to the foreground", title)
}

func (d *winDriver) key(chordText string) error {
	c, err := parseChord(chordText)
	if err != nil {
		return err
	}
	vk, needs, err := vkFor(c.key)
	if err != nil {
		return err
	}
	var mods []uint16
	for _, m := range []struct {
		on bool
		vk uint16
	}{{c.ctrl || needs.ctrl, vkControl}, {c.shift || needs.shift, vkShift}, {c.alt || needs.alt, vkMenu}, {c.cmd, vkLWin}} {
		if m.on {
			mods = append(mods, m.vk)
		}
	}
	var errs []error
	for _, m := range mods {
		errs = append(errs, keyDown(m))
	}
	errs = append(errs, keyDown(vk), keyUp(vk))
	for i := len(mods) - 1; i >= 0; i-- {
		errs = append(errs, keyUp(mods[i]))
	}
	return errors.Join(errs...)
}

// typeText types each character as a Unicode key event, as an input method
// does, so it does not depend on the keyboard layout.
func (d *winDriver) typeText(text string) error {
	for _, u := range windows.StringToUTF16(text) {
		if u == 0 {
			break
		}
		if u == '\n' {
			if err := d.key("return"); err != nil {
				return err
			}
			continue
		}
		if err := sendInput(keyInput(0, u, keyeventfUnicode)); err != nil {
			return err
		}
		if err := sendInput(keyInput(0, u, keyeventfUnicode|keyeventfKeyUp)); err != nil {
			return err
		}
	}
	return nil
}

// moveRaw moves the pointer to client-area pixel (x, y) of the window.
func (d *winDriver) moveRaw(title string, x, y int) error {
	h, err := find(title)
	if err != nil {
		return err
	}
	pt := struct{ x, y int32 }{int32(x), int32(y)}
	if ok, _, err := procClientToScreen.Call(uintptr(h), uintptr(unsafe.Pointer(&pt))); ok == 0 {
		return fmt.Errorf("ClientToScreen: %w", err)
	}
	if ok, _, err := procSetCursorPos.Call(uintptr(pt.x), uintptr(pt.y)); ok == 0 {
		return fmt.Errorf("SetCursorPos: %w", err)
	}
	time.Sleep(eventDelay)
	return nil
}

func (d *winDriver) setMapping(title string, m mapping) { d.mappings[title] = m }

func (d *winDriver) click(title string, x, y int) error {
	return d.mouse(title, x, y, mouseeventfLeftDown, mouseeventfLeftUp)
}

func (d *winDriver) rightClick(title string, x, y int) error {
	return d.mouse(title, x, y, mouseeventfRightDown, mouseeventfRightUp)
}

// mouse focuses the window, moves the pointer to page point (x, y), and
// presses and releases a button there. An uncalibrated window is taken to
// have the page at the client origin, at 100% scale.
func (d *winDriver) mouse(title string, x, y int, down, up uint32) error {
	if err := d.focus(title); err != nil {
		return err
	}
	m, ok := d.mappings[title]
	if !ok {
		m = mapping{scale: 1}
	}
	rx, ry := m.raw(x, y)
	if err := d.moveRaw(title, int(rx), int(ry)); err != nil {
		return err
	}
	return errors.Join(sendInput(buttonInput(down)), sendInput(buttonInput(up)))
}

// chooseFile types the path into the file name box of the common file
// dialog, which has the keyboard focus when the dialog opens, and accepts it.
func (d *winDriver) chooseFile(title, path string) error {
	if err := waitFocus(d, title); err != nil {
		return err
	}
	return steps(
		// Alt+N moves to the file name box, in case focus is elsewhere.
		func() error { return d.key("alt+n") },
		func() error { return d.typeText(strings.ReplaceAll(path, "/", `\`)) },
		func() error { return d.key("return") },
	)
}
