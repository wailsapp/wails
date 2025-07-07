//go:build runtimedevtools

package linux

/*
#include "gtk/gtk.h"
#include <webkit2/webkit2.h>

void ShowInspector(void *webview);
*/
import "C"

func (f *Frontend) OpenDevTools() {
	C.ShowInspector(f.mainWindow.webView)
}