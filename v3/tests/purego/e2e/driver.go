package main

import (
	"fmt"
	"runtime"
	"strings"
)

// driver injects real OS input into this process's windows: XTEST under Xvfb
// on Linux (driver_linux.go), NSEvents posted to the application on macOS
// (driver_darwin.go). Nothing here calls Wails internals: the e2e probe only
// sees what a user would see.
type driver interface {
	// focus brings the window with this title to the front and gives it
	// keyboard focus.
	focus(title string) error
	// key presses and releases one chord, e.g. key("ctrl+shift+k"),
	// key("return"), key("escape"), key("down"), key("f10").
	key(chord string) error
	// typeText types plain text into the focused control.
	typeText(text string) error
	// click and rightClick press a mouse button at a point inside the window
	// with this title, in the window's content coordinates (points, origin
	// top left), as the page reports them with getBoundingClientRect.
	click(title string, x, y int) error
	rightClick(title string, x, y int) error
}

// chord is a parsed key chord: modifiers plus one key name, all lower case.
type chord struct {
	ctrl, shift, alt, cmd bool
	key                   string
}

// parseChord parses "ctrl+shift+k". "cmdorctrl" means cmd on macOS and ctrl
// elsewhere, matching Wails accelerators.
func parseChord(s string) (chord, error) {
	var c chord
	parts := strings.Split(strings.ToLower(s), "+")
	for i, p := range parts {
		p = strings.TrimSpace(p)
		if i == len(parts)-1 {
			if p == "" {
				return c, fmt.Errorf("chord %q has no key", s)
			}
			c.key = p
			break
		}
		switch p {
		case "ctrl", "control":
			c.ctrl = true
		case "shift":
			c.shift = true
		case "alt", "option", "optionoralt":
			c.alt = true
		case "cmd", "command", "super":
			c.cmd = true
		case "cmdorctrl", "commandorcontrol":
			if runtime.GOOS == "darwin" {
				c.cmd = true
			} else {
				c.ctrl = true
			}
		default:
			return c, fmt.Errorf("chord %q: unknown modifier %q", s, p)
		}
	}
	return c, nil
}
