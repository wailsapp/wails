//go:build darwin && purego && !ios && !server

package application

import (
	"testing"

	"github.com/ebitengine/purego/objc"
)

// probeWindow is a window as far as windowDropsEscapeCommand reads one: its
// styleMask and delegate, plus the DisableEscapeExitsFullscreen flag recorded
// for its pointer. It needs no display.
type probeWindow struct {
	obj       id
	styleMask uint
	delegate  id
}

var probeWindows = map[uintptr]*probeWindow{}

func newProbeWindow(t *testing.T) *probeWindow {
	t.Helper()
	loadFrameworks()
	cls := class("WailsTestEscapeProbeWindow")
	if cls.isNil() {
		probe := func(self objc.ID) *probeWindow { return probeWindows[uintptr(self)] }
		cls = registerDelegateClass("WailsTestEscapeProbeWindow", "NSObject", nil, []objc.MethodDef{
			{Cmd: sel_("styleMask"), Fn: func(self objc.ID, cmd objc.SEL) uint { return probe(self).styleMask }},
			{Cmd: sel_("delegate"), Fn: func(self objc.ID, cmd objc.SEL) objc.ID { return objc.ID(probe(self).delegate) }},
		})
	}
	p := &probeWindow{obj: cls.send("new")}
	probeWindows[p.obj.ptr()] = p
	t.Cleanup(func() {
		windowDisableEscape.Delete(p.obj.ptr())
		delete(probeWindows, p.obj.ptr())
		p.obj.send("release")
	})
	return p
}

func newCancellingDelegate(t *testing.T) id {
	t.Helper()
	cls := class("WailsTestCancellingDelegate")
	if cls.isNil() {
		cls = registerDelegateClass("WailsTestCancellingDelegate", "NSObject", nil, []objc.MethodDef{
			{Cmd: sel_("cancelOperation:"), Fn: func(self objc.ID, cmd objc.SEL, sender objc.ID) {}},
		})
	}
	delegate := cls.send("new")
	t.Cleanup(func() { delegate.send("release") })
	return delegate
}

// TestEscapeFullscreenCommandFilterPurego is the purego counterpart of the
// cgo harness's --filter mode (testdata/escape_fullscreen/main.m).
func TestEscapeFullscreenCommandFilterPurego(t *testing.T) {
	const titled = 1 << 0
	window := newProbeWindow(t)
	drops := func(selector string) bool { return windowDropsEscapeCommand(window.obj, sel_(selector)) }

	windowDisableEscape.Store(window.obj.ptr(), true)
	window.styleMask = titled | nsWindowStyleMaskFullScreen
	if !drops("cancelOperation:") {
		t.Error("cancelOperation: is not dropped in fullscreen with the option on")
	}
	// Every other command still reaches NSWindow.
	for _, selector := range []string{"cancel:", "insertNewline:", "complete:"} {
		if drops(selector) {
			t.Errorf("%s is dropped", selector)
		}
	}
	// Only while fullscreen, and only when asked for.
	window.styleMask = titled
	if drops("cancelOperation:") {
		t.Error("cancelOperation: is dropped outside fullscreen")
	}
	window.styleMask = titled | nsWindowStyleMaskFullScreen
	windowDisableEscape.Store(window.obj.ptr(), false)
	if drops("cancelOperation:") {
		t.Error("cancelOperation: is dropped with the option off")
	}
	windowDisableEscape.Delete(window.obj.ptr())
	if drops("cancelOperation:") {
		t.Error("cancelOperation: is dropped for a window that never set the option")
	}
	// A delegate that handles the command still gets it.
	windowDisableEscape.Store(window.obj.ptr(), true)
	window.delegate = class("NSObject").send("new")
	defer window.delegate.send("release")
	if !drops("cancelOperation:") {
		t.Error("cancelOperation: is not dropped with a delegate that does not handle it")
	}
	window.delegate = newCancellingDelegate(t)
	if drops("cancelOperation:") {
		t.Error("cancelOperation: is dropped although the delegate handles it")
	}
}

// TestWindowClassesOverrideDoCommandBySelector checks that both window classes
// implement doCommandBySelector: themselves. NSResponder already answers it,
// so instancesRespondToSelector: would pass without the override; compare the
// IMP with the AppKit superclass's instead.
func TestWindowClassesOverrideDoCommandBySelector(t *testing.T) {
	loadFrameworks()
	for _, tc := range []struct {
		name  string
		cls   id
		super string
	}{
		{"WailsWebviewNSWindow", registerWebviewWindowClass(), "NSWindow"},
		{"WailsWebviewNSPanel", registerWebviewPanelClass(), "NSPanel"},
	} {
		ours := get[uintptr](tc.cls, "instanceMethodForSelector:", sel_("doCommandBySelector:"))
		inherited := get[uintptr](class(tc.super), "instanceMethodForSelector:", sel_("doCommandBySelector:"))
		if ours == inherited {
			t.Errorf("%s inherits doCommandBySelector: from %s", tc.name, tc.super)
		}
	}
}
