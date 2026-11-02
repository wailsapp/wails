//go:build linux && purego && !gtk3 && !android && !server

package application

import "testing"

func TestPuregoNativeLibraries(t *testing.T) {
	if linuxLibsErr != nil {
		t.Fatal(linuxLibsErr)
	}
	if gtk_get_major_version() != 4 {
		t.Fatal("wrong GTK ABI loaded")
	}
	if webkit_get_major_version() < 2 {
		t.Fatal("invalid WebKit version")
	}
}
