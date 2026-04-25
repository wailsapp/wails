//go:build linux
// +build linux

package main

import (
	"fmt"
	"os"
)

// This is a manual test program for issue #4673.
// Run on Linux: go run v2/test/4673/test_html_drag_drop.go
//
// Expected behavior:
// 1. Window opens with two HTML elements: a draggable div and a drop zone
// 2. Drag the "DRAG ME" box over the drop zone
// 3. The cursor should show the "copy" indicator (not "not allowed")
// 4. The drop zone should highlight and log "dragover" events
// 5. Dropping on the zone should log "drop" event
//
// The fix ensures that GTK's default drag destination is unset on the webview,
// allowing WebKit to handle HTML5 drag-and-drop events natively.
// When Wails file drops are enabled, the drag destination is re-enabled
// specifically for text/uri-list targets (file drops).

func main() {
	fmt.Println("Test for issue #4673: HTML5 drag-and-drop should work on Linux")
	fmt.Println("")
	fmt.Println("This test must be run on Linux with WebKitGTK.")
	fmt.Println("")
	fmt.Println("Manual verification steps:")
	fmt.Println("1. Create a new wails project: wails init -n dndtest -t vanilla")
	fmt.Println("2. Add a draggable element and a drop zone in index.html:")
	fmt.Println("   <div draggable=\"true\">DRAG ME</div>")
	fmt.Println("   <div id=\"dz\" style=\"border:2px dotted gray;width:300px;height:300px;\">Drop here</div>")
	fmt.Println("3. Add dragover listener: document.getElementById('dz').addEventListener('dragover', (e) => { console.log('dragover'); e.preventDefault(); })")
	fmt.Println("4. Run: wails dev")
	fmt.Println("5. Drag the element over the drop zone - dragover events should fire")
	fmt.Println("6. Without this fix, dragover events would NOT fire on Linux")
	os.Exit(0)
}
