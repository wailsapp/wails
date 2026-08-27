package application

import (
	"math"
	"strings"
	"sync"
	"testing"
)

func newTestSidebarSplit() (*MacSplitView, *MacSidebar, *MacSplitPane, *MacSplitWebviewPane) {
	sidebar := NewMacSidebar()
	sidebar.AddSection("Notes").AddItem("First").OnClick(func(*Context) {})
	split := NewMacSplitView()
	sidePane := split.AddSidebar(sidebar)
	primary := split.AddPrimaryContent()
	return split, sidebar, sidePane, primary
}

func TestMacSplitViewValidation(t *testing.T) {
	valid, _, _, _ := newTestSidebarSplit()
	if err := validateMacSplitView(valid); err != nil {
		t.Fatalf("valid native sidebar layout rejected: %v", err)
	}

	tests := []struct {
		name string
		make func() *MacSplitView
		want string
	}{
		{"fewer than two panes", func() *MacSplitView {
			s := NewMacSplitView()
			s.AddPrimaryContent()
			return s
		}, "at least two panes"},
		{"no primary", func() *MacSplitView {
			s := NewMacSplitView()
			s.AddSidebar(NewMacSidebar())
			s.AddSidebar(NewMacSidebar())
			return s
		}, "exactly one primary"},
		{"two primaries", func() *MacSplitView {
			s := NewMacSplitView()
			s.AddPrimaryContent()
			s.AddPrimaryContent()
			return s
		}, "exactly one primary"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateMacSplitView(test.make())
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v, want mention of %q", err, test.want)
			}
		})
	}
}

func TestMacSplitViewContainsNativeSidebarAndPrimaryWebViewOnly(t *testing.T) {
	split, sidebar, sidePane, primary := newTestSidebarSplit()
	panes := split.paneSnapshot()
	if len(panes) != 2 || panes[0].MacSplitPane != sidePane || panes[1] != primary {
		t.Fatal("panes should preserve native-sidebar then primary declaration order")
	}
	if panes[0].sidebar != sidebar || panes[0].primary {
		t.Fatal("leading pane should own the AppKit sidebar model, not a WebView")
	}
	if !primary.primary || primary.sidebar != nil {
		t.Fatal("only the primary pane should expose WebView operations")
	}
	if split.AddSidebar(sidebar) != nil {
		t.Fatal("one native sidebar cannot be attached twice")
	}
}

func TestMacSplitViewContainsSemanticInspector(t *testing.T) {
	split, _, _, primary := newTestSidebarSplit()
	inspector, _, _, _ := newTestInspector()
	inspectPane := split.AddInspector(inspector)
	if inspectPane == nil {
		t.Fatal("native inspector should attach as a trailing pane")
	}
	panes := split.paneSnapshot()
	if len(panes) != 3 || panes[1] != primary || panes[2].MacSplitPane != inspectPane {
		t.Fatal("pane declaration order should be sidebar, primary, inspector")
	}
	if panes[2].role != macSplitPaneInspector || panes[2].inspector != inspector || panes[2].primary {
		t.Fatal("trailing pane should own the native inspector model, not a WebView")
	}
	if err := validateMacSplitView(split); err != nil {
		t.Fatalf("valid three-pane native layout rejected: %v", err)
	}
	if split.AddInspector(inspector) != nil {
		t.Fatal("one native inspector cannot be attached twice")
	}
}

func TestMacSplitViewOwnershipAndFreeze(t *testing.T) {
	first := &WebviewWindow{}
	second := &WebviewWindow{}
	split, _, _, _ := newTestSidebarSplit()
	first.SetSplitView(split)
	if split.ownerWindow() != first {
		t.Fatal("first window should own the pending split")
	}
	if err := claimMacSplitView(split, second); err == nil {
		t.Fatal("second window must not claim an attached split")
	}
	if pane := split.AddSidebar(NewMacSidebar()); pane != nil {
		t.Fatal("structural changes must be rejected after attachment")
	}
	first.SetSplitView(nil)
	second.SetSplitView(split)
	if split.ownerWindow() != second {
		t.Fatal("cleared pending split should be reusable")
	}
}

func TestMacSplitPaneSettersAndValidation(t *testing.T) {
	_, _, pane, _ := newTestSidebarSplit()
	pane.SetMinimumThickness(210).
		SetMaximumThickness(340).
		SetPreferredThicknessFraction(.25).
		SetHoldingPriority(260).
		SetCollapsible(true).
		SetCanCollapseFromWindowResize(false).
		SetCollapsed(true)

	pane.lock.RLock()
	if pane.minimumThickness != 210 || pane.maximumThickness != 340 ||
		pane.preferredThickness != .25 || pane.holdingPriority != 260 ||
		!pane.collapsible || pane.canCollapseFromResize || !pane.collapsed {
		pane.lock.RUnlock()
		t.Fatal("pane setters did not update the pending native description")
	}
	pane.lock.RUnlock()

	pane.SetMinimumThickness(math.NaN()).SetMaximumThickness(-1).
		SetPreferredThicknessFraction(2).SetHoldingPriority(0)
	pane.lock.RLock()
	defer pane.lock.RUnlock()
	if pane.minimumThickness != 210 || pane.maximumThickness != 340 ||
		pane.preferredThickness != .25 || pane.holdingPriority != 260 {
		t.Fatal("invalid live values must not corrupt the pane description")
	}
}

func TestMacSplitPrimaryContentLayout(t *testing.T) {
	_, _, _, primary := newTestSidebarSplit()
	if got := primary.ContentLayout(); got != MacContentLayoutAutomatic {
		t.Fatalf("default content layout = %d, want automatic", got)
	}
	primary.SetContentLayout(MacContentLayoutEdgeToEdge)
	if got := primary.ContentLayout(); got != MacContentLayoutEdgeToEdge {
		t.Fatalf("content layout = %d, want edge-to-edge", got)
	}
	primary.SetContentLayout(MacContentLayout(99))
	if got := primary.ContentLayout(); got != MacContentLayoutEdgeToEdge {
		t.Fatalf("invalid content layout changed state to %d", got)
	}
}

func TestMacSplitPrimaryScrollEdgeEffectStyle(t *testing.T) {
	primary := NewMacSplitView().AddPrimaryContent()
	if got := primary.ScrollEdgeEffectStyle(); got != MacScrollEdgeEffectStyleAutomatic {
		t.Fatalf("default scroll-edge style = %d, want automatic", got)
	}
	primary.SetScrollEdgeEffectStyle(MacScrollEdgeEffectStyleSoft)
	if got := primary.ScrollEdgeEffectStyle(); got != MacScrollEdgeEffectStyleSoft {
		t.Fatalf("scroll-edge style = %d, want soft", got)
	}
	primary.SetScrollEdgeEffectStyle(MacScrollEdgeEffectStyle(99))
	if got := primary.ScrollEdgeEffectStyle(); got != MacScrollEdgeEffectStyleSoft {
		t.Fatalf("invalid style changed preference to %d", got)
	}
}

func TestMacSplitPaneCollapseCallbackDeduplicates(t *testing.T) {
	split, _, pane, _ := newTestSidebarSplit()
	internal := split.paneSnapshot()[0]
	registerMacSplitPane(internal)
	t.Cleanup(func() { unregisterMacSplitPane(internal.internalID) })
	count := 0
	pane.OnCollapsedChange(func(_ *Context, collapsed bool) {
		if !collapsed {
			t.Error("expected collapsed state")
		}
		count++
	})
	handleMacSplitPaneCollapsed(internal.internalID, true)
	handleMacSplitPaneCollapsed(internal.internalID, true)
	if count != 1 || !pane.IsCollapsed() {
		t.Fatalf("collapse callback count = %d, collapsed = %v", count, pane.IsCollapsed())
	}
}

func TestMacSplitPrimaryNavigationGeneration(t *testing.T) {
	_, _, _, primary := newTestSidebarSplit()
	registerMacSplitPane(primary)
	t.Cleanup(func() { unregisterMacSplitPane(primary.internalID) })
	loaded := 0
	primary.OnLoaded(func(*Context) { loaded++ })
	first, _ := newMacSplitPaneLoadEvent(primary.internalID)
	handleMacSplitPaneNavigationStarted(primary.internalID)
	handleMacSplitPaneLoaded(first)
	if loaded != 0 {
		t.Fatal("stale completion must not finish a newer primary navigation")
	}
	current, _ := newMacSplitPaneLoadEvent(primary.internalID)
	handleMacSplitPaneLoaded(current)
	handleMacSplitPaneLoaded(current)
	if loaded != 1 {
		t.Fatalf("current navigation callback count = %d, want 1", loaded)
	}
}

func TestMacSplitDeadHandlesAreSafe(t *testing.T) {
	split, sidebar, pane, primary := newTestSidebarSplit()
	item := sidebar.snapshot().entries[0].section.items[0]
	internal := split.paneSnapshot()[0]
	internal.markDead()
	pane.SetMinimumThickness(100).Toggle().OnCollapsedChange(func(*Context, bool) {})
	primary.markDead()
	primary.SetURL("/late").Reload().ExecJS("1+1").OnLoaded(func(*Context) {})
	if !pane.isDead() || !primary.isDead() || !sidebar.dead {
		t.Fatal("teardown should invalidate every native and primary handle")
	}
	if macSidebarItemRegistry[item.internalID] != nil {
		t.Fatal("dead sidebar should unregister its item callbacks")
	}
}

func TestMacSplitViewPrimaryPaneDelegatesPendingState(t *testing.T) {
	window := &WebviewWindow{}
	split, _, _, primary := newTestSidebarSplit()
	primary.SetURL("/editor.html").ExecJS("window.pending = true")
	window.SetSplitView(split)
	primary.lock.RLock()
	defer primary.lock.RUnlock()
	if primary.url != "" || len(primary.pendingJS) != 0 {
		t.Fatal("attachment should transfer pending primary state to the window")
	}
}

func TestMacSplitViewConcurrentConfiguration(t *testing.T) {
	_, _, pane, _ := newTestSidebarSplit()
	var wait sync.WaitGroup
	for index := 0; index < 20; index++ {
		wait.Add(1)
		go func(value int) {
			defer wait.Done()
			pane.SetMinimumThickness(float64(100 + value)).SetCollapsible(value%2 == 0)
		}(index)
	}
	wait.Wait()
}
