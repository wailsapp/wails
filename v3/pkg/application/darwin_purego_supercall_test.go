//go:build darwin && purego && !ios && !server

package application

import (
	"testing"

	"github.com/ebitengine/purego/objc"
)

// TestSuperCallFromSubclassInstance guards sendSuperFrom. A method our class
// registers, run on an instance of a SUBCLASS of it (the notch panel is a
// subclass of the panel class; KVO isa-swizzles instances the same way), must
// call the real superclass. purego's objc.ID.SendSuper looked "super" up from
// the instance's class, so it called our own method again until the stack
// overflowed (sendEvent: on the macOS VM).
func TestSuperCallFromSubclassInstance(t *testing.T) {
	registerWebviewPanelClass()
	sub := registerDelegateClass("WailsTestPanelSubclass", "WailsWebviewNSPanel", nil, []objc.MethodDef{
		{Cmd: sel_("wailsTestNoop"), Fn: func(self objc.ID, cmd objc.SEL) {}},
	})
	const borderless, buffered = 0, 2
	// AppKit wants windows created and messaged on the main thread. A
	// regression is a stack overflow, which crashes the test binary.
	runOnMain(func() {
		panel := sub.send("alloc").send("initWithContentRect:styleMask:backing:defer:",
			CGRect{Size: CGSize{Width: 10, Height: 10}}, uint(borderless), uint(buffered), true)
		// Inherited from the panel class, which calls [super cancelOperation:].
		panel.send("cancelOperation:", id(0))
		panel.send("release")
	})
}
