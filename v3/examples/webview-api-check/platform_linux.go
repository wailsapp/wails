//go:build linux

package main

import "os/exec"

func getLinuxWebViewInfo() string {
	// Try to get WebKitGTK version from pkg-config
	// For GTK4 builds, this will be webkitgtk-6.0
	// For GTK3 builds, this will be webkit2gtk-4.1
	out, err := exec.Command("pkg-config", "--modversion", "webkitgtk-6.0").Output()
	if err == nil {
		return "WebKitGTK " + string(out[:len(out)-1])
	}
	out, err = exec.Command("pkg-config", "--modversion", "webkit2gtk-4.1").Output()
	if err == nil {
		return "WebKit2GTK " + string(out[:len(out)-1])
	}
	return "WebKitGTK (unknown version)"
}

func getGTKVersionInfo() string {
	// Try GTK4 first
	out, err := exec.Command("pkg-config", "--modversion", "gtk4").Output()
	if err == nil {
		return "GTK " + string(out[:len(out)-1])
	}
	out, err = exec.Command("pkg-config", "--modversion", "gtk+-3.0").Output()
	if err == nil {
		return "GTK " + string(out[:len(out)-1])
	}
	return "GTK (unknown version)"
}
