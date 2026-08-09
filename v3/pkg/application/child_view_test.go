package application

import "testing"

type testChildView struct {
	attached int
	detached int
}

func (v *testChildView) Attach(Window) { v.attached++ }
func (v *testChildView) Detach()       { v.detached++ }

func TestWebviewWindowChildViewsAreManagedOnce(t *testing.T) {
	window := &WebviewWindow{}
	view := &testChildView{}

	window.AddChildView(view)
	window.AddChildView(view)
	if got := len(window.childViews); got != 1 {
		t.Fatalf("registered views = %d, want 1", got)
	}

	window.RemoveChildView(view)
	if view.detached != 1 {
		t.Fatalf("detach calls = %d, want 1", view.detached)
	}
	if got := len(window.childViews); got != 0 {
		t.Fatalf("registered views after removal = %d, want 0", got)
	}
}

func TestWebviewWindowRemovingUnknownChildViewDoesNotDetach(t *testing.T) {
	window := &WebviewWindow{}
	registered := &testChildView{}
	unknown := &testChildView{}
	window.AddChildView(registered)

	window.RemoveChildView(unknown)
	if unknown.detached != 0 {
		t.Fatalf("unknown view detached %d times", unknown.detached)
	}
}
