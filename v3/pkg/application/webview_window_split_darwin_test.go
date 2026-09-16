//go:build darwin && !ios && !server

package application

import (
	"testing"
	"unsafe"
)

func TestMacosWindowSidebarSplitLayoutDetection(t *testing.T) {
	window := &WebviewWindow{}
	impl := &macosWebviewWindow{parent: window}
	if impl.hasSidebarSplitLayout() {
		t.Fatal("a window without a split view must not report a sidebar layout")
	}

	split := NewMacSplitView()
	split.AddSidebar(NewMacSidebar())
	split.AddPrimaryContent()
	window.SetSplitView(split)
	if !impl.hasSidebarSplitLayout() {
		t.Fatal("a pending sidebar split layout should be detected before installation")
	}

}

func TestMacosWindowInspectorSplitLayoutDetection(t *testing.T) {
	window := &WebviewWindow{}
	impl := &macosWebviewWindow{parent: window}
	if impl.hasInspectorSplitLayout() {
		t.Fatal("a window without a split view must not report an inspector layout")
	}

	split := NewMacSplitView()
	split.AddPrimaryContent()
	inspector, _, _, _ := newTestInspector()
	split.AddInspector(inspector)
	window.SetSplitView(split)
	if !impl.hasInspectorSplitLayout() {
		t.Fatal("a pending inspector layout should be detected before installation")
	}
}

func TestTeardownSplitViewInvalidatesHandlesIdempotently(t *testing.T) {
	window := &WebviewWindow{}
	impl := &macosWebviewWindow{parent: window}
	split := NewMacSplitView()
	sidebarModel := NewMacSidebar()
	sidebarModel.AddItem("First")
	sidebar := split.AddSidebar(sidebarModel)
	primary := split.AddPrimaryContent()
	window.SetSplitView(split)
	sidebarInternal := split.paneSnapshot()[0]
	registerMacSplitPane(sidebarInternal)
	registerMacSplitPane(primary)
	impl.activeSplitView = split

	impl.teardownSplitView()

	if !sidebar.isDead() || !primary.isDead() || !sidebarModel.dead {
		t.Fatal("teardown should mark every pane handle dead")
	}
	if macSplitPaneByID(sidebar.internalID) != nil || macSplitPaneByID(primary.internalID) != nil {
		t.Fatal("teardown should remove every registry entry")
	}
	window.splitViewLock.RLock()
	remaining := window.splitView
	window.splitViewLock.RUnlock()
	if remaining != nil {
		t.Fatal("teardown should clear the window's split view reference")
	}

	// Dead handles are safe no-ops and repeated teardown must not panic.
	sidebar.SetMinimumThickness(100).Toggle()
	impl.teardownSplitView()
}

func TestInstallSplitViewLateGuards(t *testing.T) {
	window := &WebviewWindow{}
	impl := &macosWebviewWindow{parent: window}
	split := NewMacSplitView()
	split.AddSidebar(NewMacSidebar())
	split.AddPrimaryContent()
	window.SetSplitView(split)

	// Without a native window there is nothing to install into, and the
	// window's own error must come back rather than a crash in AppKit.
	if err := impl.installSplitViewLate(); err == nil || err == ErrMacSplitViewAlreadyInstalled {
		t.Fatalf("installSplitViewLate without an NSWindow = %v, want a native-window error", err)
	}
	if split.isInstalled() {
		t.Fatal("a failed late installation must not mark the layout installed")
	}

	// An installed layout is never replaced.
	impl.nsWindow = unsafe.Pointer(&struct{}{})
	impl.activeSplitView = split
	if err := impl.installSplitViewLate(); err != ErrMacSplitViewAlreadyInstalled {
		t.Fatalf("installSplitViewLate over an installed layout = %v, want ErrMacSplitViewAlreadyInstalled", err)
	}
}
