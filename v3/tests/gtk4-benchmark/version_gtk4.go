//go:build linux && !gtk3

package main

func getGTKVersionString() string {
	return "GTK4 (WebKitGTK 6.0)"
}
