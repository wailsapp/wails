package application

import "testing"

func nativeEditorSplitForTest() (*MacSplitView, *MacTextEditor) {
	sidebar := NewMacSidebar()
	sidebar.AddSection("Files").AddItem("One")
	editor := NewMacTextEditor().SetText("hello")
	split := NewMacSplitView()
	split.AddSidebar(sidebar)
	split.AddTextEditor(editor)
	return split, editor
}

func TestNativeWindowClaimsNativeEditorSplit(t *testing.T) {
	split, _ := nativeEditorSplitForTest()
	first := newNativeWindow(NativeWindowOptions{Name: "first"})
	second := newNativeWindow(NativeWindowOptions{Name: "second"})

	if err := first.SetSplitView(split); err != nil {
		t.Fatalf("first SetSplitView: %v", err)
	}
	if split.ownerWindow() != first {
		t.Fatal("split owner was not the native window")
	}
	if err := second.SetSplitView(split); err == nil {
		t.Fatal("second window unexpectedly claimed an owned split")
	}
	if err := first.SetSplitView(nil); err != nil {
		t.Fatalf("release split: %v", err)
	}
	if err := second.SetSplitView(split); err != nil {
		t.Fatalf("second SetSplitView after release: %v", err)
	}
}

func TestMacTextEditorOfflineState(t *testing.T) {
	editor := NewMacTextEditor()
	if got := editor.Text(); got != "" {
		t.Fatalf("initial text = %q", got)
	}
	editor.SetText("native text").SetEditable(false)
	if got := editor.Text(); got != "native text" {
		t.Fatalf("Text = %q", got)
	}
	editor.lock.RLock()
	editable := editor.editable
	editor.lock.RUnlock()
	if editable {
		t.Fatal("editor remained editable")
	}
}

func TestMacTextEditorOnlyClearsAcceptedTextVersion(t *testing.T) {
	editor := NewMacTextEditor().SetText("first")
	_, _, _, firstVersion := editor.snapshot()
	editor.SetText("second")
	_, _, _, secondVersion := editor.snapshot()

	editor.clearCachedText(firstVersion)
	if got := editor.Text(); got != "second" {
		t.Fatalf("stale native acceptance cleared newer text: %q", got)
	}
	editor.clearCachedText(secondVersion)
	editor.lock.RLock()
	cached := editor.text
	editor.lock.RUnlock()
	if cached != "" {
		t.Fatalf("accepted text remained cached: %q", cached)
	}
}

func TestNativeAndWebviewWindowsShareToolbarType(t *testing.T) {
	toolbar := NewMacToolbar()
	toolbar.AddButton("Save").OnClick(func(*Context) {})
	native := newNativeWindow(NativeWindowOptions{Name: "native"})
	webview := NewWindow(WebviewWindowOptions{Name: "webview"})

	if err := native.SetToolbar(toolbar); err != nil {
		t.Fatalf("native SetToolbar: %v", err)
	}
	if _, err := claimMacToolbar(toolbar, webview); err == nil {
		t.Fatal("webview unexpectedly claimed native window's toolbar")
	}
	if err := native.SetToolbar(nil); err != nil {
		t.Fatalf("native toolbar release: %v", err)
	}
	if _, err := claimMacToolbar(toolbar, webview); err != nil {
		t.Fatalf("webview could not claim released toolbar: %v", err)
	}
	releaseMacToolbarOwnership(toolbar, webview)
}
