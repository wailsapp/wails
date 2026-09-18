package application

import (
	"errors"
	"fmt"
	"sort"
	"unsafe"
)

// MacTabOrder selects where a window is inserted relative to the receiving
// window when it joins that window's tab group. Values match
// NSWindowOrderingMode so they can be passed straight to AppKit.
type MacTabOrder int

const (
	// MacTabOrderAbove inserts the new tab after the receiving window
	// (NSWindowAbove).
	MacTabOrderAbove MacTabOrder = 1
	// MacTabOrderBelow inserts the new tab before the receiving window
	// (NSWindowBelow).
	MacTabOrderBelow MacTabOrder = -1
)

func validMacTabOrder(order MacTabOrder) bool {
	return order == MacTabOrderAbove || order == MacTabOrderBelow
}

var (
	// ErrMacWindowTabbingDisallowed means the window was created with
	// Mac.TabbingMode left unset or set to MacWindowTabbingModeDisallowed.
	// Tabbing mode is fixed at creation; create the window with Automatic or
	// Preferred to use the tab APIs.
	ErrMacWindowTabbingDisallowed = errors.New("window tabbing is disallowed: create the window with Mac.TabbingMode set to Automatic or Preferred")
	// ErrMacWindowTabNotCreated means the native window does not exist yet.
	// Tab operations need a live NSWindow, so call them after the window has
	// been shown (for example from a menu handler or after App.Run starts).
	ErrMacWindowTabNotCreated = errors.New("the native window has not been created yet")
	// ErrMacWindowTabTargetRequired means a nil or not-yet-created window was
	// passed to AddTab.
	ErrMacWindowTabTargetRequired = errors.New("a created window to add as a tab is required")
	// ErrMacWindowTabsUnsupported means native window tabs are unavailable on
	// the current platform.
	ErrMacWindowTabsUnsupported = errors.New("macOS window tabs are unavailable on this platform")
)

// MacWindowTabGroup is a live handle for the NSWindowTabGroup a window
// belongs to. It resolves the native group on every call, so it stays valid
// as tabs are added, detached or closed. Obtain one with
// WebviewWindow.TabGroup or NativeWindow.TabGroup. A nil handle is safe to
// use: every method returns its zero value.
//
// Tab groups can contain both WebviewWindow and NativeWindow members.
// Windows returns only the WebviewWindow members because the Window
// interface is WebView-specific; NativeWindows returns the rest and Count
// reports the full native tab count.
type MacWindowTabGroup struct {
	// handle returns the NSWindow of the window that produced this group
	// handle, or nil once that window has been destroyed.
	handle func() unsafe.Pointer
}

func newMacWindowTabGroup(handle func() unsafe.Pointer) *MacWindowTabGroup {
	if handle == nil {
		return nil
	}
	return &MacWindowTabGroup{handle: handle}
}

func (g *MacWindowTabGroup) nsWindow() unsafe.Pointer {
	if g == nil || g.handle == nil {
		return nil
	}
	return g.handle()
}

// Identifier returns the tabbing identifier shared by every window in the
// group, or "" when the group is gone.
func (g *MacWindowTabGroup) Identifier() string {
	nsWindow := g.nsWindow()
	if nsWindow == nil {
		return ""
	}
	return macWindowTabsIdentifier(nsWindow)
}

// Count returns the number of native windows in the group, including
// NativeWindow members and windows not managed by Wails.
func (g *MacWindowTabGroup) Count() int {
	nsWindow := g.nsWindow()
	if nsWindow == nil {
		return 0
	}
	return len(macWindowTabsGroupWindows(nsWindow))
}

// Windows returns the WebviewWindow members of the group in tab order.
// NativeWindow members are skipped; see NativeWindows.
func (g *MacWindowTabGroup) Windows() []Window {
	nsWindow := g.nsWindow()
	if nsWindow == nil {
		return nil
	}
	var result []Window
	for _, pointer := range macWindowTabsGroupWindows(nsWindow) {
		if window := windowForNSWindow(pointer); window != nil {
			result = append(result, window)
		}
	}
	return result
}

// NativeWindows returns the NativeWindow members of the group in tab order.
func (g *MacWindowTabGroup) NativeWindows() []*NativeWindow {
	nsWindow := g.nsWindow()
	if nsWindow == nil {
		return nil
	}
	var result []*NativeWindow
	for _, pointer := range macWindowTabsGroupWindows(nsWindow) {
		if window := nativeWindowForNSWindow(pointer); window != nil {
			result = append(result, window)
		}
	}
	return result
}

// SelectedWindow returns the WebviewWindow whose tab is currently selected,
// or nil when the selected tab is not a WebviewWindow.
func (g *MacWindowTabGroup) SelectedWindow() Window {
	nsWindow := g.nsWindow()
	if nsWindow == nil {
		return nil
	}
	selected := macWindowTabsSelectedWindow(nsWindow)
	if selected == nil {
		return nil
	}
	if window := windowForNSWindow(selected); window != nil {
		return window
	}
	return nil
}

// SelectNext selects the tab after the current one, wrapping around.
func (g *MacWindowTabGroup) SelectNext() {
	if nsWindow := g.nsWindow(); nsWindow != nil {
		macWindowTabsSelectNext(nsWindow)
	}
}

// SelectPrevious selects the tab before the current one, wrapping around.
func (g *MacWindowTabGroup) SelectPrevious() {
	if nsWindow := g.nsWindow(); nsWindow != nil {
		macWindowTabsSelectPrevious(nsWindow)
	}
}

// Select makes the given window's tab the selected tab. Windows that are not
// members of the group, and nil or not-yet-created windows, are ignored.
func (g *MacWindowTabGroup) Select(window Window) {
	nsWindow := g.nsWindow()
	if nsWindow == nil || window == nil {
		return
	}
	target := window.NativeWindow()
	if target == nil {
		return
	}
	macWindowTabsSelect(nsWindow, target)
}

// SelectNative makes the given NativeWindow's tab the selected tab.
func (g *MacWindowTabGroup) SelectNative(window *NativeWindow) {
	nsWindow := g.nsWindow()
	if nsWindow == nil || window == nil {
		return
	}
	target := window.NativeWindow()
	if target == nil {
		return
	}
	macWindowTabsSelect(nsWindow, target)
}

// IsTabBarVisible reports whether the tab bar is showing.
func (g *MacWindowTabGroup) IsTabBarVisible() bool {
	nsWindow := g.nsWindow()
	if nsWindow == nil {
		return false
	}
	return macWindowTabsIsTabBarVisible(nsWindow)
}

// ToggleTabBar shows or hides the tab bar, like View > Show/Hide Tab Bar.
func (g *MacWindowTabGroup) ToggleTabBar() {
	if nsWindow := g.nsWindow(); nsWindow != nil {
		macWindowTabsToggleTabBar(nsWindow)
	}
}

// IsOverviewVisible reports whether the tab overview (the grid of tab
// thumbnails) is showing.
func (g *MacWindowTabGroup) IsOverviewVisible() bool {
	nsWindow := g.nsWindow()
	if nsWindow == nil {
		return false
	}
	return macWindowTabsIsOverviewVisible(nsWindow)
}

// ToggleTabOverview shows or hides the tab overview, like View > Show All
// Tabs.
func (g *MacWindowTabGroup) ToggleTabOverview() {
	if nsWindow := g.nsWindow(); nsWindow != nil {
		macWindowTabsToggleOverview(nsWindow)
	}
}

// tabbingAllowed reports whether the window was created with a tabbing mode
// that lets AppKit place it in a tab group.
func (w *WebviewWindow) tabbingAllowed() bool {
	if w == nil {
		return false
	}
	switch w.options.Mac.TabbingMode {
	case MacWindowTabbingModeAutomatic, MacWindowTabbingModePreferred:
		return true
	default:
		return false
	}
}

// AddTab adds other to this window's tab group, inserting it after
// (MacTabOrderAbove) or before (MacTabOrderBelow) this window's tab. Both
// windows must have been created with a tabbing mode other than Disallowed
// and both native windows must already exist. Off macOS it returns
// ErrMacWindowTabsUnsupported.
func (w *WebviewWindow) AddTab(other Window, ordered MacTabOrder) error {
	if w == nil {
		return ErrMacWindowTabNotCreated
	}
	if !validMacTabOrder(ordered) {
		return fmt.Errorf("unknown macOS tab order %d", ordered)
	}
	if !w.tabbingAllowed() {
		return ErrMacWindowTabbingDisallowed
	}
	if other == nil {
		return ErrMacWindowTabTargetRequired
	}
	if otherWebview, ok := other.(*WebviewWindow); ok {
		if otherWebview == nil {
			return ErrMacWindowTabTargetRequired
		}
		if otherWebview == w {
			return fmt.Errorf("a window cannot be added as a tab of itself")
		}
		if !otherWebview.tabbingAllowed() {
			return fmt.Errorf("%w (window %q)", ErrMacWindowTabbingDisallowed, otherWebview.Name())
		}
	}
	nsWindow := w.NativeWindow()
	if nsWindow == nil {
		return ErrMacWindowTabNotCreated
	}
	target := other.NativeWindow()
	if target == nil {
		return ErrMacWindowTabTargetRequired
	}
	return macWindowTabsAdd(nsWindow, target, ordered)
}

// AddNativeTab adds a NativeWindow to this window's tab group. NativeWindow
// has no tabbing-mode option and uses AppKit's automatic mode, so only this
// window's tabbing mode is checked. Off macOS it returns
// ErrMacWindowTabsUnsupported.
func (w *WebviewWindow) AddNativeTab(other *NativeWindow, ordered MacTabOrder) error {
	if w == nil {
		return ErrMacWindowTabNotCreated
	}
	if !validMacTabOrder(ordered) {
		return fmt.Errorf("unknown macOS tab order %d", ordered)
	}
	if !w.tabbingAllowed() {
		return ErrMacWindowTabbingDisallowed
	}
	if other == nil {
		return ErrMacWindowTabTargetRequired
	}
	nsWindow := w.NativeWindow()
	if nsWindow == nil {
		return ErrMacWindowTabNotCreated
	}
	target := other.NativeWindow()
	if target == nil {
		return ErrMacWindowTabTargetRequired
	}
	return macWindowTabsAdd(nsWindow, target, ordered)
}

// TabGroup returns a handle for the tab group this window belongs to. It
// returns nil off macOS, before the native window exists, and when the
// window cannot be tabbed. A window created with Automatic or Preferred
// tabbing normally belongs to a group even while it is the only tab.
func (w *WebviewWindow) TabGroup() *MacWindowTabGroup {
	if w == nil || !w.tabbingAllowed() {
		return nil
	}
	nsWindow := w.NativeWindow()
	if nsWindow == nil || !macWindowTabsHasGroup(nsWindow) {
		return nil
	}
	return newMacWindowTabGroup(w.NativeWindow)
}

// MoveTabToNewWindow detaches this window from its tab group into its own
// window, like the Window > Move Tab to New Window menu item. It is a no-op
// when the window is not tabbed.
func (w *WebviewWindow) MoveTabToNewWindow() Window {
	if w == nil {
		return w
	}
	if nsWindow := w.NativeWindow(); nsWindow != nil {
		macWindowTabsMoveToNewWindow(nsWindow)
	}
	return w
}

// MergeAllWindows merges every window that can be tabbed into this window's
// tab group, like Window > Merge All Windows.
func (w *WebviewWindow) MergeAllWindows() Window {
	if w == nil {
		return w
	}
	if nsWindow := w.NativeWindow(); nsWindow != nil {
		macWindowTabsMergeAllWindows(nsWindow)
	}
	return w
}

// SetTabTitle sets the title shown on this window's tab. It defaults to the
// window title; pass "" to restore that. It requires the native window to
// exist and is a no-op elsewhere.
func (w *WebviewWindow) SetTabTitle(title string) Window {
	if w == nil {
		return w
	}
	if nsWindow := w.NativeWindow(); nsWindow != nil {
		macWindowTabsSetTitle(nsWindow, title)
	}
	return w
}

// SetTabTooltip sets the tooltip shown when hovering this window's tab. Pass
// "" to remove it. It requires the native window to exist and is a no-op
// elsewhere.
func (w *WebviewWindow) SetTabTooltip(tooltip string) Window {
	if w == nil {
		return w
	}
	if nsWindow := w.NativeWindow(); nsWindow != nil {
		macWindowTabsSetToolTip(nsWindow, tooltip)
	}
	return w
}

// TabGroup returns a handle for the tab group this native window belongs
// to, or nil when it is not part of one. NativeWindow uses AppKit's
// automatic tabbing mode, so it can be added to a WebviewWindow's group with
// WebviewWindow.AddNativeTab.
func (w *NativeWindow) TabGroup() *MacWindowTabGroup {
	if w == nil {
		return nil
	}
	nsWindow := w.NativeWindow()
	if nsWindow == nil || !macWindowTabsHasGroup(nsWindow) {
		return nil
	}
	return newMacWindowTabGroup(w.NativeWindow)
}

// TabGroups returns one handle per distinct tab group that contains at least
// one WebviewWindow, ordered by the lowest window ID in each group. It is
// empty off macOS.
func (wm *WindowManager) TabGroups() []*MacWindowTabGroup {
	if wm == nil || wm.app == nil {
		return nil
	}
	type entry struct {
		group  *MacWindowTabGroup
		lowest uint
	}
	seen := make(map[unsafe.Pointer]*entry)
	for _, window := range wm.GetAll() {
		webview, ok := window.(*WebviewWindow)
		if !ok || webview == nil {
			continue
		}
		nsWindow := webview.NativeWindow()
		if nsWindow == nil || !webview.tabbingAllowed() {
			continue
		}
		pointer := macWindowTabsGroupPointer(nsWindow)
		if pointer == nil {
			continue
		}
		if existing := seen[pointer]; existing != nil {
			if webview.ID() < existing.lowest {
				existing.lowest = webview.ID()
			}
			continue
		}
		seen[pointer] = &entry{group: newMacWindowTabGroup(webview.NativeWindow), lowest: webview.ID()}
	}
	entries := make([]*entry, 0, len(seen))
	for _, item := range seen {
		entries = append(entries, item)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].lowest < entries[j].lowest })
	result := make([]*MacWindowTabGroup, 0, len(entries))
	for _, item := range entries {
		result = append(result, item.group)
	}
	return result
}

// windowForNSWindow finds the WebviewWindow that owns a native window
// pointer, or nil when the pointer belongs to no Wails WebView window.
func windowForNSWindow(pointer unsafe.Pointer) Window {
	if pointer == nil || globalApplication == nil || globalApplication.Window == nil {
		return nil
	}
	for _, window := range globalApplication.Window.GetAll() {
		if window == nil {
			continue
		}
		if window.NativeWindow() == pointer {
			return window
		}
	}
	return nil
}

// nativeWindowForNSWindow finds the NativeWindow that owns a native window
// pointer, or nil when the pointer belongs to no Wails native window.
func nativeWindowForNSWindow(pointer unsafe.Pointer) *NativeWindow {
	if pointer == nil || globalApplication == nil {
		return nil
	}
	globalApplication.nativeWindowsLock.RLock()
	windows := make([]*NativeWindow, 0, len(globalApplication.nativeWindows))
	for _, window := range globalApplication.nativeWindows {
		windows = append(windows, window)
	}
	globalApplication.nativeWindowsLock.RUnlock()
	for _, window := range windows {
		if window != nil && window.NativeWindow() == pointer {
			return window
		}
	}
	return nil
}
