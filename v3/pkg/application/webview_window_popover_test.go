package application

import (
	"errors"
	"testing"
)

func TestMacPopoverBehaviorValuesMatchNSPopoverBehavior(t *testing.T) {
	if MacPopoverBehaviorApplicationDefined != 0 || MacPopoverBehaviorTransient != 1 || MacPopoverBehaviorSemitransient != 2 {
		t.Fatal("popover behaviors do not match NSPopoverBehavior")
	}
	if !validMacPopoverBehavior(MacPopoverBehaviorSemitransient) || validMacPopoverBehavior(MacPopoverBehavior(3)) || validMacPopoverBehavior(-1) {
		t.Fatal("behavior validation is wrong")
	}
}

func TestMacRectEdgeConvertsToNSRectEdge(t *testing.T) {
	cases := map[MacRectEdge]int{
		MacRectEdgeDefault: 1,
		MacRectEdgeMinX:    0,
		MacRectEdgeMinY:    1,
		MacRectEdgeMaxX:    2,
		MacRectEdgeMaxY:    3,
	}
	for edge, want := range cases {
		if got := edge.nsRectEdge(); got != want {
			t.Fatalf("nsRectEdge(%d) = %d, want %d", edge, got, want)
		}
	}
	if validMacRectEdge(MacRectEdge(5)) || validMacRectEdge(-1) {
		t.Fatal("unknown edge was accepted")
	}
}

func TestMacPopoverOptionDefaults(t *testing.T) {
	options, err := (MacPopoverOptions{}).normalise()
	if err != nil {
		t.Fatalf("normalise: %v", err)
	}
	if options.Width != macPopoverDefaultWidth {
		t.Fatalf("default width = %v", options.Width)
	}
	if want := MacAccessoryDefaultTitlebarHeight + 2*macPopoverContentPadding; options.Height != want {
		t.Fatalf("default height = %v, want %v", options.Height, want)
	}
	if options.Behavior != MacPopoverBehaviorApplicationDefined || options.PreferredEdge != MacRectEdgeDefault || options.DisableAnimation {
		t.Fatalf("defaults changed: %+v", options)
	}

	content := NewMacAccessory(MacAccessoryLayoutBottom).SetHeight(44)
	options, err = (MacPopoverOptions{Content: content, Width: 200}).normalise()
	if err != nil {
		t.Fatalf("normalise with content: %v", err)
	}
	if options.Width != 200 || options.Height != 44+2*macPopoverContentPadding {
		t.Fatalf("content-sized options = %v x %v", options.Width, options.Height)
	}

	for name, invalid := range map[string]MacPopoverOptions{
		"behavior": {Behavior: MacPopoverBehavior(9)},
		"edge":     {PreferredEdge: MacRectEdge(9)},
		"width":    {Width: -1},
		"height":   {Height: -1},
	} {
		if _, err := invalid.normalise(); err == nil {
			t.Fatalf("%s: invalid options were accepted", name)
		}
	}
}

func TestNewMacPopoverIsUsableWithInvalidOptions(t *testing.T) {
	popover := NewMacPopover(MacPopoverOptions{Behavior: MacPopoverBehavior(7), Width: -5, PreferredEdge: MacRectEdge(9)})
	defer popover.Destroy()
	options := popover.Options()
	if options.Behavior != MacPopoverBehaviorApplicationDefined || options.Width != macPopoverDefaultWidth || options.PreferredEdge != MacRectEdgeDefault {
		t.Fatalf("invalid options were not replaced: %+v", options)
	}
	if popover.ID() == 0 || popover.IsDestroyed() || popover.IsShown() {
		t.Fatal("new popover state is wrong")
	}
	if popover.Content() != nil {
		t.Fatal("Content must be nil when none was given")
	}
	if popover.resolveEdge(MacRectEdgeDefault) != MacRectEdgeMinY || popover.resolveEdge(MacRectEdgeMaxY) != MacRectEdgeMaxY {
		t.Fatal("resolveEdge is wrong")
	}
	if err := popover.SetBehavior(MacPopoverBehavior(4)); err == nil {
		t.Fatal("unknown behavior was accepted by SetBehavior")
	}
	if err := popover.SetBehavior(MacPopoverBehaviorTransient); err != nil || popover.Options().Behavior != MacPopoverBehaviorTransient {
		t.Fatalf("SetBehavior = %v", err)
	}
	popover.SetContentSize(0, 10)
	if popover.Options().Width != macPopoverDefaultWidth {
		t.Fatal("SetContentSize accepted a non-positive width")
	}
	popover.SetContentSize(100, 50)
	if o := popover.Options(); o.Width != 100 || o.Height != 50 {
		t.Fatalf("SetContentSize did not store the size: %+v", o)
	}
}

func TestMacPopoverShowErrorPaths(t *testing.T) {
	var nilPopover *MacPopover
	window := &WebviewWindow{options: WebviewWindowOptions{Name: "uncreated"}}
	if err := nilPopover.ShowRelativeTo(Rect{}, window, MacRectEdgeDefault); !errors.Is(err, ErrMacPopoverRequired) {
		t.Fatalf("nil popover = %v", err)
	}
	if err := macToolbarItemShowPopover(&MacToolbarItem{}, nilPopover); !errors.Is(err, ErrMacPopoverRequired) {
		t.Fatalf("nil popover from toolbar = %v", err)
	}
	if err := macSystemTrayShowPopover(&SystemTray{}, nilPopover); !errors.Is(err, ErrMacPopoverRequired) {
		t.Fatalf("nil popover from tray = %v", err)
	}
	nilPopover.Close()
	nilPopover.Destroy()
	if nilPopover.IsShown() || !nilPopover.IsDestroyed() || nilPopover.ID() != 0 || nilPopover.Content() != nil {
		t.Fatal("nil popover reported state")
	}
	nilPopover.OnClose(func() {})()

	popover := NewMacPopover(MacPopoverOptions{Behavior: MacPopoverBehaviorTransient})
	if !macPopoversSupported() {
		defer popover.Destroy()
		if err := popover.ShowRelativeTo(Rect{}, window, MacRectEdgeDefault); !errors.Is(err, ErrMacPopoverUnsupported) {
			t.Fatalf("ShowRelativeTo off macOS = %v", err)
		}
		if err := popover.ShowRelativeToNative(Rect{}, &NativeWindow{}, MacRectEdgeDefault); !errors.Is(err, ErrMacPopoverUnsupported) {
			t.Fatalf("ShowRelativeToNative off macOS = %v", err)
		}
		if err := macToolbarItemShowPopover(&MacToolbarItem{}, popover); !errors.Is(err, ErrMacPopoverUnsupported) {
			t.Fatalf("toolbar show off macOS = %v", err)
		}
		if err := macSystemTrayShowPopover(&SystemTray{}, popover); !errors.Is(err, ErrMacPopoverUnsupported) {
			t.Fatalf("tray show off macOS = %v", err)
		}
		return
	}
	// Destroying before any show never touches the native layer, so the
	// destroyed error path can be exercised without an application.
	popover.Destroy()
	popover.Destroy()
	if !popover.IsDestroyed() {
		t.Fatal("Destroy did not mark the popover")
	}
	if err := popover.ShowRelativeTo(Rect{}, window, MacRectEdgeDefault); !errors.Is(err, ErrMacPopoverDestroyed) {
		t.Fatalf("show after destroy = %v", err)
	}
	if err := macToolbarItemShowPopover(&MacToolbarItem{}, popover); !errors.Is(err, ErrMacPopoverDestroyed) {
		t.Fatalf("toolbar show after destroy = %v", err)
	}
	if err := macSystemTrayShowPopover(&SystemTray{}, popover); !errors.Is(err, ErrMacPopoverDestroyed) {
		t.Fatalf("tray show after destroy = %v", err)
	}
	macPopoverRegistry.RLock()
	_, registered := macPopoverRegistry.popovers[popover.id]
	macPopoverRegistry.RUnlock()
	if registered {
		t.Fatal("destroyed popover is still registered")
	}
}

func TestMacPopoverCloseCallbacks(t *testing.T) {
	popover := NewMacPopover(MacPopoverOptions{})
	var closed int
	remove := popover.OnClose(func() { closed++ })
	popover.OnClose(func() { closed += 10 })
	handleMacPopoverClosed(popover.id)
	if closed != 11 {
		t.Fatalf("callbacks ran %d", closed)
	}
	remove()
	handleMacPopoverClosed(popover.id)
	if closed != 21 {
		t.Fatalf("callbacks after removal ran %d", closed)
	}
	handleMacPopoverClosed(popover.id + 1000) // unknown popovers are ignored
	popover.Destroy()
	handleMacPopoverClosed(popover.id) // destroyed popovers are ignored
	if closed != 21 {
		t.Fatal("callbacks ran after Destroy")
	}
}

func TestMacPopoverImplementsAccessoryHost(t *testing.T) {
	var host macAccessoryWindow = NewMacPopover(MacPopoverOptions{})
	if host.ID() == 0 {
		t.Fatal("popover ID must be non-zero")
	}
	host.Error("no application: %s", "ignored")
}
