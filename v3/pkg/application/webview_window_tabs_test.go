package application

import (
	"errors"
	"testing"
)

func TestMacTabOrderValuesMatchNSWindowOrderingMode(t *testing.T) {
	if MacTabOrderAbove != 1 || MacTabOrderBelow != -1 {
		t.Fatalf("unexpected tab order values: above=%d below=%d", MacTabOrderAbove, MacTabOrderBelow)
	}
	if !validMacTabOrder(MacTabOrderAbove) || !validMacTabOrder(MacTabOrderBelow) {
		t.Fatal("documented tab orders were rejected")
	}
	if validMacTabOrder(0) || validMacTabOrder(MacTabOrder(2)) {
		t.Fatal("unknown tab order was accepted")
	}
}

func TestWebviewWindowTabbingAllowedFollowsCreationMode(t *testing.T) {
	cases := map[MacWindowTabbingMode]bool{
		MacWindowTabbingModeDefault:    false,
		MacWindowTabbingModeDisallowed: false,
		MacWindowTabbingModeAutomatic:  true,
		MacWindowTabbingModePreferred:  true,
	}
	for mode, want := range cases {
		window := &WebviewWindow{options: WebviewWindowOptions{Mac: MacWindow{TabbingMode: mode}}}
		if got := window.tabbingAllowed(); got != want {
			t.Fatalf("tabbingAllowed() for mode %d = %v, want %v", mode, got, want)
		}
	}
	var nilWindow *WebviewWindow
	if nilWindow.tabbingAllowed() {
		t.Fatal("nil window reported tabbing as allowed")
	}
}

func TestAddTabRejectsDisallowedWindowsBeforeTouchingNative(t *testing.T) {
	disallowed := &WebviewWindow{options: WebviewWindowOptions{Mac: MacWindow{TabbingMode: MacWindowTabbingModeDisallowed}}}
	preferred := &WebviewWindow{options: WebviewWindowOptions{Mac: MacWindow{TabbingMode: MacWindowTabbingModePreferred}}}
	unset := &WebviewWindow{}

	if err := disallowed.AddTab(preferred, MacTabOrderAbove); !errors.Is(err, ErrMacWindowTabbingDisallowed) {
		t.Fatalf("disallowed receiver error = %v, want ErrMacWindowTabbingDisallowed", err)
	}
	if err := unset.AddTab(preferred, MacTabOrderAbove); !errors.Is(err, ErrMacWindowTabbingDisallowed) {
		t.Fatalf("unset receiver error = %v, want ErrMacWindowTabbingDisallowed", err)
	}
	if err := preferred.AddTab(disallowed, MacTabOrderAbove); !errors.Is(err, ErrMacWindowTabbingDisallowed) {
		t.Fatalf("disallowed target error = %v, want ErrMacWindowTabbingDisallowed", err)
	}
	if err := preferred.AddTab(preferred, MacTabOrderAbove); err == nil {
		t.Fatal("adding a window as a tab of itself was accepted")
	}
	if err := preferred.AddTab(nil, MacTabOrderAbove); !errors.Is(err, ErrMacWindowTabTargetRequired) {
		t.Fatalf("nil target error = %v, want ErrMacWindowTabTargetRequired", err)
	}
	if err := preferred.AddTab(preferred, MacTabOrder(7)); err == nil {
		t.Fatal("unknown tab order was accepted")
	}
	if err := preferred.AddNativeTab(nil, MacTabOrderBelow); !errors.Is(err, ErrMacWindowTabTargetRequired) {
		t.Fatalf("nil native target error = %v, want ErrMacWindowTabTargetRequired", err)
	}
	if err := disallowed.AddNativeTab(&NativeWindow{}, MacTabOrderBelow); !errors.Is(err, ErrMacWindowTabbingDisallowed) {
		t.Fatalf("disallowed native receiver error = %v, want ErrMacWindowTabbingDisallowed", err)
	}
}

func TestAddTabRequiresCreatedWindows(t *testing.T) {
	receiver := &WebviewWindow{options: WebviewWindowOptions{Mac: MacWindow{TabbingMode: MacWindowTabbingModePreferred}}}
	other := &WebviewWindow{options: WebviewWindowOptions{Mac: MacWindow{TabbingMode: MacWindowTabbingModeAutomatic}}}
	if err := receiver.AddTab(other, MacTabOrderAbove); !errors.Is(err, ErrMacWindowTabNotCreated) {
		t.Fatalf("uncreated receiver error = %v, want ErrMacWindowTabNotCreated", err)
	}
	if err := receiver.AddNativeTab(&NativeWindow{}, MacTabOrderAbove); !errors.Is(err, ErrMacWindowTabNotCreated) {
		t.Fatalf("uncreated native receiver error = %v, want ErrMacWindowTabNotCreated", err)
	}
	var nilWindow *WebviewWindow
	if err := nilWindow.AddTab(other, MacTabOrderAbove); !errors.Is(err, ErrMacWindowTabNotCreated) {
		t.Fatalf("nil receiver error = %v, want ErrMacWindowTabNotCreated", err)
	}
}

func TestTabGroupIsNilWithoutNativeWindow(t *testing.T) {
	preferred := &WebviewWindow{options: WebviewWindowOptions{Mac: MacWindow{TabbingMode: MacWindowTabbingModePreferred}}}
	if preferred.TabGroup() != nil {
		t.Fatal("uncreated window returned a tab group")
	}
	disallowed := &WebviewWindow{}
	if disallowed.TabGroup() != nil {
		t.Fatal("disallowed window returned a tab group")
	}
	var nilWindow *WebviewWindow
	if nilWindow.TabGroup() != nil {
		t.Fatal("nil window returned a tab group")
	}
	if (&NativeWindow{}).TabGroup() != nil {
		t.Fatal("uncreated native window returned a tab group")
	}
	var nilNative *NativeWindow
	if nilNative.TabGroup() != nil {
		t.Fatal("nil native window returned a tab group")
	}
}

func TestTabMutatorsAreSafeBeforeCreation(t *testing.T) {
	window := &WebviewWindow{options: WebviewWindowOptions{Mac: MacWindow{TabbingMode: MacWindowTabbingModePreferred}}}
	if window.MoveTabToNewWindow() != window || window.MergeAllWindows() != window ||
		window.SetTabTitle("Tab") != window || window.SetTabTooltip("Tip") != window {
		t.Fatal("tab mutators did not return the receiver")
	}
	var nilWindow *WebviewWindow
	nilWindow.MoveTabToNewWindow()
	nilWindow.MergeAllWindows()
	nilWindow.SetTabTitle("Tab")
	nilWindow.SetTabTooltip("Tip")
}

func TestNilMacWindowTabGroupIsSafe(t *testing.T) {
	var group *MacWindowTabGroup
	if group.Windows() != nil || group.NativeWindows() != nil || group.SelectedWindow() != nil {
		t.Fatal("nil group returned members")
	}
	if group.Count() != 0 || group.Identifier() != "" {
		t.Fatal("nil group returned a count or identifier")
	}
	if group.IsTabBarVisible() || group.IsOverviewVisible() {
		t.Fatal("nil group reported visible chrome")
	}
	group.SelectNext()
	group.SelectPrevious()
	group.Select(nil)
	group.SelectNative(nil)
	group.ToggleTabBar()
	group.ToggleTabOverview()
	if newMacWindowTabGroup(nil) != nil {
		t.Fatal("nil handle produced a group")
	}
}

func TestMacWindowTabGroupWithDeadHandleIsSafe(t *testing.T) {
	group := newMacWindowTabGroup((&WebviewWindow{}).NativeWindow)
	if group == nil {
		t.Fatal("expected a group handle")
	}
	if group.Windows() != nil || group.Count() != 0 || group.Identifier() != "" {
		t.Fatal("dead handle returned members")
	}
	if group.SelectedWindow() != nil || group.IsTabBarVisible() || group.IsOverviewVisible() {
		t.Fatal("dead handle reported state")
	}
	group.SelectNext()
	group.SelectPrevious()
	group.Select(&WebviewWindow{})
	group.ToggleTabBar()
	group.ToggleTabOverview()
}

func TestWindowManagerTabGroupsNilSafe(t *testing.T) {
	var manager *WindowManager
	if manager.TabGroups() != nil {
		t.Fatal("nil manager returned tab groups")
	}
	if (&WindowManager{}).TabGroups() != nil {
		t.Fatal("manager without an app returned tab groups")
	}
}

func TestWindowForNSWindowWithoutApplication(t *testing.T) {
	if windowForNSWindow(nil) != nil || nativeWindowForNSWindow(nil) != nil {
		t.Fatal("nil pointer resolved to a window")
	}
}
