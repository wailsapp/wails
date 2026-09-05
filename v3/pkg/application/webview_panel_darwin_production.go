//go:build darwin && !ios && !server && production && !devtools

package application

import "unsafe"

func configureEmbeddedPanelDevTools(_ unsafe.Pointer, _ bool) {}
func openEmbeddedPanelDevTools(_ unsafe.Pointer)              {}
