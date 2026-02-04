//go:build linux && !gtk4

package main

func getGTKVersionString() string {
	return "GTK3 (WebKit2GTK 4.1)"
}
