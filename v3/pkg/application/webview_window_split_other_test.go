//go:build !darwin || ios || server

package application

import "testing"

func TestUnsupportedMacNativeSidebarRemainsDescriptiveAndSafe(t *testing.T) {
	split, sidebar, pane, primary := newTestSidebarSplit()
	window := &WebviewWindow{}
	window.SetSplitView(split)
	pane.Toggle()
	primary.SetURL("/editor.html").ExecJS("window.test = true")
	if sidebar == nil || !pane.IsCollapsed() {
		t.Fatal("portable API should retain safe descriptive state")
	}
}
