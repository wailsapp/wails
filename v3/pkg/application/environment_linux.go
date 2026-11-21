//go:build linux && !android

package application

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func detectCompositor() string {
	if os.Getenv("HYPRLAND_INSTANCE_SIGNATURE") != "" {
		return "hyprland"
	}
	if os.Getenv("SWAYSOCK") != "" {
		return "sway"
	}
	if os.Getenv("I3SOCK") != "" {
		return "i3"
	}
	if desktop := os.Getenv("XDG_CURRENT_DESKTOP"); desktop != "" {
		return strings.ToLower(desktop)
	}
	return "unknown"
}

func detectFocusFollowsMouse() bool {
	compositor := detectCompositor()
	switch compositor {
	case "hyprland", "sway", "i3":
		return true
	}
	return false
}

func isWayland() bool {
	return os.Getenv("XDG_SESSION_TYPE") == "wayland" ||
		os.Getenv("WAYLAND_DISPLAY") != ""
}

func isTilingWM() bool {
	switch detectCompositor() {
	case "hyprland", "sway", "i3":
		return true
	}
	return false
}

func getCursorPositionFromCompositor() (x, y int, ok bool) {
	switch detectCompositor() {
	case "hyprland":
		out, err := hyprlandIPC("cursorpos")
		if err != nil {
			return 0, 0, false
		}
		return parseCursorPos(strings.TrimSpace(out))
	}
	// Sway's IPC reports no pointer position (GET_SEATS has none), so on
	// sway, as elsewhere, the caller falls back.
	return 0, 0, false
}

func parseCursorPos(s string) (x, y int, ok bool) {
	parts := strings.Split(s, ", ")
	if len(parts) != 2 {
		return 0, 0, false
	}
	var err error
	x, err = strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, false
	}
	y, err = strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, false
	}
	return x, y, true
}

func hyprlandIPC(command string) (string, error) {
	sig := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	if sig == "" {
		return "", fmt.Errorf("HYPRLAND_INSTANCE_SIGNATURE not set")
	}

	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = fmt.Sprintf("/run/user/%d", os.Getuid())
	}

	socketPath := filepath.Join(runtimeDir, "hypr", sig, ".socket.sock")
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(command))
	if err != nil {
		return "", err
	}

	var result strings.Builder
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if n > 0 {
			result.Write(buf[:n])
		}
		if err != nil {
			break
		}
		if n < len(buf) {
			break
		}
	}

	return result.String(), nil
}
