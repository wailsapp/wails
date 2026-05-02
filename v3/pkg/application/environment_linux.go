//go:build linux

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
	case "sway":
		out, err := swayIPC("get_seats")
		if err != nil {
			return 0, 0, false
		}
		return parseSwayCursor(out)
	}
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

func parseSwayCursor(json string) (x, y int, ok bool) {
	cursorIdx := strings.Index(json, `"cursor"`)
	if cursorIdx == -1 {
		return 0, 0, false
	}

	xIdx := strings.Index(json[cursorIdx:], `"x"`)
	if xIdx == -1 {
		return 0, 0, false
	}
	xStart := cursorIdx + xIdx + 4
	xEnd := strings.IndexAny(json[xStart:], ",}")
	if xEnd == -1 {
		return 0, 0, false
	}
	x, _ = strconv.Atoi(strings.TrimSpace(json[xStart : xStart+xEnd]))

	yIdx := strings.Index(json[cursorIdx:], `"y"`)
	if yIdx == -1 {
		return 0, 0, false
	}
	yStart := cursorIdx + yIdx + 4
	yEnd := strings.IndexAny(json[yStart:], ",}")
	if yEnd == -1 {
		return 0, 0, false
	}
	y, _ = strconv.Atoi(strings.TrimSpace(json[yStart : yStart+yEnd]))

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

func swayIPC(command string) (string, error) {
	socketPath := os.Getenv("SWAYSOCK")
	if socketPath == "" {
		return "", fmt.Errorf("SWAYSOCK not set")
	}

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	msgType := uint32(0)
	if strings.HasPrefix(command, "get_") {
		switch command {
		case "get_seats":
			msgType = 5
		}
	} else {
		msgType = 0
	}

	payload := []byte(command)
	header := make([]byte, 14)
	copy(header[0:6], "i3-ipc")
	header[6] = byte(len(payload))
	header[7] = byte(len(payload) >> 8)
	header[8] = byte(len(payload) >> 16)
	header[9] = byte(len(payload) >> 24)
	header[10] = byte(msgType)
	header[11] = byte(msgType >> 8)
	header[12] = byte(msgType >> 16)
	header[13] = byte(msgType >> 24)

	conn.Write(header)
	conn.Write(payload)

	respHeader := make([]byte, 14)
	_, err = conn.Read(respHeader)
	if err != nil {
		return "", err
	}

	respLen := uint32(respHeader[6]) | uint32(respHeader[7])<<8 | uint32(respHeader[8])<<16 | uint32(respHeader[9])<<24
	respBody := make([]byte, respLen)
	_, err = conn.Read(respBody)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}
