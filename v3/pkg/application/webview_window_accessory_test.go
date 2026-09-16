package application

import (
	"errors"
	"testing"
	"unsafe"
)

func TestMacScrollEdgeEffectStyleValues(t *testing.T) {
	if MacScrollEdgeEffectStyleAutomatic != 0 || MacScrollEdgeEffectStyleSoft != 1 ||
		MacScrollEdgeEffectStyleHard != 2 {
		t.Fatalf("unexpected scroll-edge style values: %d, %d, %d",
			MacScrollEdgeEffectStyleAutomatic, MacScrollEdgeEffectStyleSoft,
			MacScrollEdgeEffectStyleHard)
	}
	for _, style := range []MacScrollEdgeEffectStyle{
		MacScrollEdgeEffectStyleAutomatic,
		MacScrollEdgeEffectStyleSoft,
		MacScrollEdgeEffectStyleHard,
	} {
		if !validMacScrollEdgeEffectStyle(style) {
			t.Fatalf("validMacScrollEdgeEffectStyle(%d) = false", style)
		}
	}
	if validMacScrollEdgeEffectStyle(MacScrollEdgeEffectStyle(99)) {
		t.Fatal("unknown scroll-edge style was accepted")
	}
}

func TestWrapMacAccessoryViewControllerRejectsNil(t *testing.T) {
	controller, err := WrapMacAccessoryViewController(nil)
	if controller != nil {
		t.Fatal("nil native controller produced a wrapper")
	}
	if !errors.Is(err, ErrMacAccessoryControllerRequired) {
		t.Fatalf("error = %v, want ErrMacAccessoryControllerRequired", err)
	}
}

func TestMacAccessoryViewControllerRejectsUnknownStyleBeforeNativeCall(t *testing.T) {
	marker := byte(0)
	controller := &MacAccessoryViewController{
		native: unsafe.Pointer(&marker),
		kind:   MacAccessoryViewControllerKindTitlebar,
	}
	if err := controller.SetPreferredScrollEdgeEffectStyle(MacScrollEdgeEffectStyle(99)); err == nil {
		t.Fatal("unknown scroll-edge style was accepted")
	}
}

func TestNilMacAccessoryViewControllerIsSafe(t *testing.T) {
	var controller *MacAccessoryViewController
	if controller.NativeController() != nil {
		t.Fatal("nil wrapper returned a native controller")
	}
	if controller.Kind() != MacAccessoryViewControllerKindUnknown {
		t.Fatal("nil wrapper returned a controller kind")
	}
	if controller.SupportsPreferredScrollEdgeEffectStyle() {
		t.Fatal("nil wrapper reported style support")
	}
	if !errors.Is(controller.SetPreferredScrollEdgeEffectStyle(MacScrollEdgeEffectStyleAutomatic),
		ErrMacAccessoryControllerRequired) {
		t.Fatal("nil wrapper did not report the required-controller error")
	}
	if _, err := controller.PreferredScrollEdgeEffectStyle(); !errors.Is(err, ErrMacAccessoryControllerRequired) {
		t.Fatalf("getter error = %v, want ErrMacAccessoryControllerRequired", err)
	}
}

func TestMacAccessoryLayoutConstants(t *testing.T) {
	if MacAccessoryLayoutLeading != 0 || MacAccessoryLayoutTrailing != 1 ||
		MacAccessoryLayoutBottom != 2 || MacAccessoryLayoutTop != 3 {
		t.Fatalf("unexpected layout values: %d, %d, %d, %d", MacAccessoryLayoutLeading,
			MacAccessoryLayoutTrailing, MacAccessoryLayoutBottom, MacAccessoryLayoutTop)
	}
	for layout, want := range map[MacAccessoryLayout][2]bool{
		MacAccessoryLayoutLeading:  {true, false},
		MacAccessoryLayoutTrailing: {true, false},
		MacAccessoryLayoutBottom:   {true, true},
		MacAccessoryLayoutTop:      {false, true},
	} {
		if layout.titlebar() != want[0] || layout.splitItem() != want[1] {
			t.Fatalf("%s: titlebar=%v splitItem=%v, want %v/%v", layout, layout.titlebar(), layout.splitItem(), want[0], want[1])
		}
		if !validMacAccessoryLayout(layout) {
			t.Fatalf("validMacAccessoryLayout(%s) = false", layout)
		}
	}
	if validMacAccessoryLayout(MacAccessoryLayout(42)) {
		t.Fatal("unknown layout was accepted")
	}
	if NewMacAccessory(MacAccessoryLayout(42)).Layout() != MacAccessoryLayoutBottom {
		t.Fatal("unknown layout did not fall back to Bottom")
	}
	if MacAccessoryDefaultTitlebarHeight != 28 || MacAccessoryDefaultSplitItemHeight != 36 {
		t.Fatal("unexpected default accessory heights")
	}
}

func TestMacAccessoryControlOrdering(t *testing.T) {
	accessory := NewMacAccessory(MacAccessoryLayoutBottom)
	search := accessory.AddSearch("Search")
	segmented := accessory.AddSegmented([]string{"All", "Starred"}, 5)
	button := accessory.AddButton("Compose")
	symbolButton := accessory.AddSymbolButton("square.and.pencil")
	menuButton := accessory.AddMenuButton("Sort", NewMenu())
	label := accessory.AddLabel("Ready")
	space := accessory.AddFlexibleSpace()
	marker := byte(0)
	native := accessory.AddNativeView(unsafe.Pointer(&marker))

	controls := accessory.Controls()
	want := []*MacAccessoryControl{search, segmented, button, symbolButton, menuButton, label, space, native}
	if len(controls) != len(want) {
		t.Fatalf("got %d controls, want %d", len(controls), len(want))
	}
	kinds := []MacAccessoryControlKind{
		MacAccessoryControlSearch, MacAccessoryControlSegmented, MacAccessoryControlButton,
		MacAccessoryControlButton, MacAccessoryControlMenuButton, MacAccessoryControlLabel,
		MacAccessoryControlFlexibleSpace, MacAccessoryControlNativeView,
	}
	var previousID uint64
	for index, control := range controls {
		if control != want[index] {
			t.Fatalf("control %d out of order", index)
		}
		if control.Kind() != kinds[index] {
			t.Fatalf("control %d kind = %d, want %d", index, control.Kind(), kinds[index])
		}
		if control.Accessory() != accessory {
			t.Fatalf("control %d does not point back to its accessory", index)
		}
		if control.internalID <= previousID {
			t.Fatalf("control %d identifier %d is not increasing", index, control.internalID)
		}
		previousID = control.internalID
	}
	if segmented.SelectedSegment() != -1 {
		t.Fatalf("out-of-range initial selection %d was kept", segmented.SelectedSegment())
	}
	segmented.SetSelectedSegment(1)
	if segmented.SelectedSegment() != 1 {
		t.Fatal("SetSelectedSegment did not update the model")
	}
	segmented.SetSelectedSegment(7)
	if segmented.SelectedSegment() != -1 {
		t.Fatal("SetSelectedSegment accepted an out-of-range index")
	}
	if label.SetText("Busy").Text() != "Busy" {
		t.Fatal("SetText did not update the model")
	}
	if symbolButton.snapshot().symbol != "square.and.pencil" || button.snapshot().text != "Compose" {
		t.Fatal("button constructors did not store their content")
	}
	if native.snapshot().nativeView != unsafe.Pointer(&marker) {
		t.Fatal("AddNativeView did not keep the view pointer")
	}
	if accessory.SetHeight(-1).Height() != 0 || accessory.SetHeight(40).Height() != 40 {
		t.Fatal("SetHeight validation failed")
	}
	if !accessory.SetHidden(true).IsHidden() {
		t.Fatal("SetHidden did not update the model")
	}
	if accessory.SetPreferredScrollEdgeEffectStyle(MacScrollEdgeEffectStyleSoft).PreferredScrollEdgeEffectStyle() != MacScrollEdgeEffectStyleSoft {
		t.Fatal("scroll-edge style was not stored")
	}
	if accessory.Controller() != nil || accessory.IsAttached() {
		t.Fatal("a detached accessory reported a controller")
	}
}

func TestMacAccessoryDefaultHeights(t *testing.T) {
	titlebar := NewMacAccessory(MacAccessoryLayoutLeading)
	if titlebar.resolvedHeight() != MacAccessoryDefaultTitlebarHeight {
		t.Fatalf("titlebar default height = %v", titlebar.resolvedHeight())
	}
	pane := NewMacAccessory(MacAccessoryLayoutTop)
	if pane.resolvedHeight() != MacAccessoryDefaultSplitItemHeight {
		t.Fatalf("pane default height = %v", pane.resolvedHeight())
	}
	bottom := NewMacAccessory(MacAccessoryLayoutBottom)
	bottom.target = macAccessoryTarget{pane: &MacSplitPane{}}
	if bottom.resolvedHeight() != MacAccessoryDefaultSplitItemHeight {
		t.Fatalf("bottom pane default height = %v", bottom.resolvedHeight())
	}
	bottom.SetHeight(44)
	if bottom.resolvedHeight() != 44 {
		t.Fatalf("explicit height = %v, want 44", bottom.resolvedHeight())
	}
}

func TestMacAccessoryRemoveUnregistersCallbacks(t *testing.T) {
	accessory := NewMacAccessory(MacAccessoryLayoutTrailing)
	button := accessory.AddButton("One")
	search := accessory.AddSearch("Two")
	registerMacAccessoryControls(accessory)
	if macAccessoryControlByID(button.internalID) != button || macAccessoryControlByID(search.internalID) != search {
		t.Fatal("controls were not registered")
	}

	window := &WebviewWindow{}
	accessory.target = macAccessoryTarget{window: window}
	queueMacAccessory(window, accessory)
	if pending := pendingMacAccessories(window); len(pending) != 1 || pending[0] != accessory {
		t.Fatal("accessory was not queued")
	}

	accessory.Remove()
	if macAccessoryControlByID(button.internalID) != nil || macAccessoryControlByID(search.internalID) != nil {
		t.Fatal("Remove left control callbacks registered")
	}
	if len(pendingMacAccessories(window)) != 0 {
		t.Fatal("Remove left the accessory queued")
	}
	if !accessory.target.empty() || accessory.IsAttached() {
		t.Fatal("Remove did not clear the attach target")
	}
	// Removed accessories keep their controls and can be attached again.
	if len(accessory.Controls()) != 2 {
		t.Fatal("Remove dropped the controls")
	}
	var nilAccessory *MacAccessory
	nilAccessory.Remove()
}

func TestMacAccessoryEventRouting(t *testing.T) {
	accessory := NewMacAccessory(MacAccessoryLayoutBottom)
	clicked := make(chan struct{}, 1)
	button := accessory.AddButton("Go").OnClick(func(*Context) { clicked <- struct{}{} })
	searched := make(chan string, 1)
	search := accessory.AddSearch("Find").OnSearch(func(_ *Context, query string) { searched <- query })
	type selection struct {
		index int
		label string
	}
	selected := make(chan selection, 1)
	segmented := accessory.AddSegmented([]string{"Inbox", "Sent"}, 0).
		OnSelectionChange(func(_ *Context, index int, label string) { selected <- selection{index, label} })
	// A control without a callback must be ignored rather than panic.
	silent := accessory.AddButton("Silent")

	registerMacAccessoryControls(accessory)
	t.Cleanup(func() { unregisterMacAccessoryControls(accessory) })

	handleMacAccessoryEvent(macAccessoryEvent{controlID: button.internalID, kind: macAccessoryEventClick})
	select {
	case <-clicked:
	default:
		t.Fatal("OnClick was not invoked")
	}
	handleMacAccessoryEvent(macAccessoryEvent{controlID: search.internalID, kind: macAccessoryEventSearch, text: "peach"})
	select {
	case query := <-searched:
		if query != "peach" {
			t.Fatalf("OnSearch query = %q", query)
		}
	default:
		t.Fatal("OnSearch was not invoked")
	}
	if search.Text() != "peach" {
		t.Fatal("search event did not update the control text")
	}
	handleMacAccessoryEvent(macAccessoryEvent{controlID: segmented.internalID, kind: macAccessoryEventSelection, index: 1})
	select {
	case got := <-selected:
		if got.index != 1 || got.label != "Sent" {
			t.Fatalf("OnSelectionChange = %+v", got)
		}
	default:
		t.Fatal("OnSelectionChange was not invoked")
	}
	if segmented.SelectedSegment() != 1 {
		t.Fatal("selection event did not update the model")
	}
	handleMacAccessoryEvent(macAccessoryEvent{controlID: silent.internalID, kind: macAccessoryEventClick})
	handleMacAccessoryEvent(macAccessoryEvent{controlID: 1 << 62, kind: macAccessoryEventClick})

	// Unregistered controls stop receiving events.
	unregisterMacAccessoryControls(accessory)
	handleMacAccessoryEvent(macAccessoryEvent{controlID: button.internalID, kind: macAccessoryEventClick})
	select {
	case <-clicked:
		t.Fatal("OnClick fired after the control was unregistered")
	default:
	}
}

func TestMacAccessoryAttachValidation(t *testing.T) {
	window := &WebviewWindow{}
	if err := window.AddTitlebarAccessory(nil); !errors.Is(err, ErrMacAccessoryRequired) {
		t.Fatalf("nil accessory error = %v", err)
	}
	if err := window.AddTitlebarAccessory(NewMacAccessory(MacAccessoryLayoutTop)); !errors.Is(err, ErrMacAccessoryLayoutInvalid) {
		t.Fatalf("top layout in titlebar error = %v", err)
	}
	native := &NativeWindow{}
	if err := native.AddTitlebarAccessory(NewMacAccessory(MacAccessoryLayoutTop)); !errors.Is(err, ErrMacAccessoryLayoutInvalid) {
		t.Fatalf("NativeWindow top layout error = %v", err)
	}

	split := NewMacSplitView()
	pane := split.AddPrimaryContent().MacSplitPane
	if err := pane.AddTopAccessory(nil); !errors.Is(err, ErrMacAccessoryRequired) {
		t.Fatalf("nil pane accessory error = %v", err)
	}
	if err := pane.AddTopAccessory(NewMacAccessory(MacAccessoryLayoutBottom)); !errors.Is(err, ErrMacAccessoryLayoutInvalid) {
		t.Fatalf("bottom layout for AddTopAccessory error = %v", err)
	}
	if err := pane.AddBottomAccessory(NewMacAccessory(MacAccessoryLayoutLeading)); !errors.Is(err, ErrMacAccessoryLayoutInvalid) {
		t.Fatalf("leading layout for AddBottomAccessory error = %v", err)
	}
	if err := (&MacSplitPane{}).AddTopAccessory(NewMacAccessory(MacAccessoryLayoutTop)); !errors.Is(err, ErrMacAccessoryPaneRequired) {
		t.Fatalf("detached pane error = %v", err)
	}
	var nilPane *MacSplitPane
	if err := nilPane.AddBottomAccessory(NewMacAccessory(MacAccessoryLayoutBottom)); !errors.Is(err, ErrMacAccessoryPaneRequired) {
		t.Fatalf("nil pane error = %v", err)
	}

	accessory := NewMacAccessory(MacAccessoryLayoutLeading)
	err := window.AddTitlebarAccessory(accessory)
	paneAccessory := NewMacAccessory(MacAccessoryLayoutTop)
	paneErr := pane.AddTopAccessory(paneAccessory)
	if !macAccessoriesSupported {
		if !errors.Is(err, ErrMacAccessoryUnsupported) || !errors.Is(paneErr, ErrMacAccessoryUnsupported) {
			t.Fatalf("errors off macOS = %v / %v, want ErrMacAccessoryUnsupported", err, paneErr)
		}
		if !accessory.target.empty() {
			t.Fatal("unsupported attach claimed the accessory")
		}
		return
	}
	// On macOS the native window does not exist in this test, so both
	// accessories are queued for the creation path.
	if err != nil || paneErr != nil {
		t.Fatalf("queued attach errors = %v / %v", err, paneErr)
	}
	if pending := pendingMacAccessories(window); len(pending) != 1 || pending[0] != accessory {
		t.Fatal("titlebar accessory was not queued for the window")
	}
	if pending := pendingMacAccessories(pane); len(pending) != 1 || pending[0] != paneAccessory {
		t.Fatal("pane accessory was not queued for the pane")
	}
	if !paneAccessory.target.top || paneAccessory.target.pane != pane {
		t.Fatal("pane accessory target was not recorded")
	}
	if err := window.AddTitlebarAccessory(accessory); !errors.Is(err, ErrMacAccessoryAttached) {
		t.Fatalf("second attach error = %v, want ErrMacAccessoryAttached", err)
	}
	if err := pane.AddTopAccessory(paneAccessory); !errors.Is(err, ErrMacAccessoryAttached) {
		t.Fatalf("second pane attach error = %v, want ErrMacAccessoryAttached", err)
	}
	accessory.Remove()
	paneAccessory.Remove()
	if len(pendingMacAccessories(window)) != 0 || len(pendingMacAccessories(pane)) != 0 {
		t.Fatal("Remove did not drain the pending queues")
	}
	// A removed accessory can be queued again.
	if err := window.AddTitlebarAccessory(accessory); err != nil {
		t.Fatalf("re-attach after Remove failed: %v", err)
	}
	accessory.Remove()
}

func TestMacAccessoryNilSafety(t *testing.T) {
	var accessory *MacAccessory
	accessory.SetHidden(true).SetHeight(10).SetFullScreenMinHeight(1).SetAutomaticallyAdjustsSize(true).
		SetAutomaticallyAppliesContentInsets(false).SetPreferredScrollEdgeEffectStyle(MacScrollEdgeEffectStyleHard)
	if accessory.AddSearch("x") != nil || accessory.AddButton("x") != nil || accessory.AddFlexibleSpace() != nil ||
		accessory.AddLabel("x") != nil || accessory.AddMenuButton("x", nil) != nil ||
		accessory.AddSegmented(nil, 0) != nil || accessory.AddSymbolButton("x") != nil || accessory.AddNativeView(nil) != nil {
		t.Fatal("nil accessory produced a control")
	}
	if accessory.Controls() != nil || accessory.IsHidden() || accessory.Height() != 0 || accessory.Controller() != nil {
		t.Fatal("nil accessory reported state")
	}
	var control *MacAccessoryControl
	control.SetEnabled(false).SetHidden(true).SetTooltip("t").SetWidth(1).SetText("t").SetSymbol("s").
		SetPlaceholder("p").SetIncremental(true).SetSegmentSymbols("a").SetSelectedSegment(0).SetMenu(nil).
		OnClick(nil).OnSearch(nil).OnSelectionChange(nil)
	if control.Kind() != MacAccessoryControlFlexibleSpace || control.Accessory() != nil ||
		control.Text() != "" || control.SelectedSegment() != -1 {
		t.Fatal("nil control reported state")
	}
	// Controls cannot be added to an attached accessory, but without an
	// application there is nothing to log; the control is simply returned.
	attached := NewMacAccessory(MacAccessoryLayoutBottom)
	marker := byte(0)
	attached.native = unsafe.Pointer(&marker)
	if attached.AddButton("late") == nil || len(attached.Controls()) != 0 {
		t.Fatal("late control was installed into the model")
	}
}
