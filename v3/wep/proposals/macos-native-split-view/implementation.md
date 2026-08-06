# Native macOS Split View and Source-List Sidebar — Implementation Guide

This guide implements [`proposal.md`](./proposal.md). The central architectural
rule is simple: the window's existing WKWebView is the primary pane, and the
sidebar is AppKit. Do not create a sidebar WebView, even if its page and WebKit
backing are transparent.

## 1. File layout

Platform-neutral API and state:

```text
v3/pkg/application/webview_window_split.go
v3/pkg/application/webview_window_sidebar.go
v3/pkg/application/webview_window_toolbar.go
```

Desktop macOS bridges:

```text
v3/pkg/application/webview_window_split_darwin.go
v3/pkg/application/webview_window_split_darwin.h
v3/pkg/application/webview_window_split_darwin.m
v3/pkg/application/webview_window_sidebar_darwin.go
v3/pkg/application/webview_window_toolbar_darwin.go
v3/pkg/application/webview_window_toolbar_darwin.h
v3/pkg/application/webview_window_toolbar_darwin.m
```

Portable stubs and tests live beside those files. The Daymark demo keeps only
`assets/editor.html`; `split.go` owns the native source-list model.

## 2. Build the platform-neutral model first

`MacSplitView` stores panes in leading-to-trailing declaration order. Each
internal pane wraps `MacSplitPane`; only the primary pane is exposed as
`MacSplitWebviewPane`. A sidebar pane stores a `*MacSidebar` instead of a URL.

Validation before native mutation must enforce:

- at least two panes;
- exactly one primary pane;
- every sidebar pane has a non-nil native sidebar;
- finite, nonnegative dimensions;
- minimum not greater than maximum;
- preferred fraction within `(0, 1]`; and
- holding priority within `1...1000`.

Attaching freezes pane structure. Live pane state remains mutable.

`MacSidebar` stores ordered root entries. An entry is either a selectable item
or a section containing selectable items. Generate monotonically increasing
internal IDs for sections and items. Do not expose identifier parameters in
the public API.

Snapshot model fields under locks before invoking native code. Never execute a
user callback while holding a split, sidebar, section, or item lock.

## 3. Build the native source-list controller

Implement a private Objective-C `NSViewController` with:

- a borderless `NSScrollView`;
- a single-column `NSOutlineView`;
- `NSTableViewSelectionHighlightStyleSourceList`;
- transparent outline and scroll backgrounds;
- group rows for `MacSidebarSection`; and
- native `NSTableCellView` item rows with optional SF Symbols.

Both AppKit views must remain transparent. The semantic
`NSSplitViewItem.sidebarWithViewController:` supplies the sidebar material.
Adding an opaque layer, CSS background, or WebView above it defeats the native
appearance.

Use the outline view's data source for hierarchy and delegate for:

- group-row classification;
- native cells;
- disabled-row selection rejection; and
- selection-change callbacks.

Suppress callbacks while applying programmatic selection or reloading a
snapshot. Programmatic `SetSelectedItem` changes UI state but must not pretend
the user clicked the row.

Keep the outline column sized to the visible outline width during layout so
labels use the complete sidebar rather than the column's initial width.

## 4. Install the split atomically

Run installation on the AppKit application thread after the primary WebView
exists and before the toolbar is attached.

Order:

1. Validate the Go layout.
2. Create the native split owner and pane records.
3. Copy the complete native sidebar snapshot into its pane record.
4. Construct every pane controller without changing the window.
5. Abort and release the candidate if any controller cannot be created.
6. Retain and remove the existing primary WebView from its old superview.
7. Preserve any existing visual-effect or Liquid Glass backdrop and drag/drop
   overlay.
8. Create `NSSplitViewItem` instances with semantic constructors.
9. Reparent the primary WebView into the primary controller.
10. Apply explicitly configured pane overrides and initial collapse state.
11. Register KVO for `collapsed`.
12. Enable `NSWindowStyleMaskFullSizeContentView`.
13. Replace `window.contentViewController` while preserving the window frame.
14. Apply the resolved primary content layout: either constrain its top edge
    to `NSWindow.contentLayoutGuide` or extend it to its pane's full top edge.
15. Associate the existing primary WebView with its primary pane ID.
16. Reload and expand the native sidebar, then apply initial selection.

There is no auxiliary WebKit configuration path, navigation delegate, scheme
handler, script-message bridge, or sidebar asset request.

## 5. Route native selection without public IDs

Register each live `MacSidebarItem` in a process-local map immediately before
native installation. The native row stores only its generated integer ID.

On `outlineViewSelectionDidChange:`:

1. ignore empty, section, disabled, and programmatic selections;
2. send the internal item ID to Go;
3. resolve the retained item handle;
4. update `MacSidebar.selected`; and
5. invoke the current `OnClick` callback asynchronously.

Item additions after installation register the new handle before reloading the
native snapshot. Teardown removes every entry before native resources are
released.

## 6. Apply live sidebar changes

The initial implementation may rebuild the small source-list snapshot for
`SetLabel`, `SetSymbol`, `SetTooltip`, `SetEnabled`, `SetHidden`, and new rows.
Perform the rebuild on the application thread and restore selection after
`reloadData`.

If profiling later shows large source lists need finer updates, add native
item-diff operations behind the same public handles. Do not change the public
API merely to expose AppKit row identifiers.

## 7. Preserve primary WebView behavior

`NSWindowStyleMaskFullSizeContentView` is needed so the semantic sidebar can
extend beneath the unified titlebar. The primary WebView then follows its
resolved `MacContentLayout` policy. `BelowToolbar` constrains its top edge to
`NSWindow.contentLayoutGuide.topAnchor`. `EdgeToEdge` constrains it to its
primary-pane host's top edge, allowing scrolling content to pass beneath the
floating system toolbar. Leading, trailing, and bottom edges always remain
attached to the pane host.

`Automatic` on the pane inherits `MacWindow.ContentLayout`; automatic at the
window level follows `MacTitleBar.FullSizeContent`. Retain the active Auto
Layout constraints and deactivate them before a live
`MacSplitWebviewPane.SetContentLayout` update. Apply all constraint changes on
the AppKit application thread.

On macOS 26 and newer, the standard `NSToolbar` supplies Liquid Glass and the
system scroll-edge treatment. Do not add an `NSGlassEffectView`, gradient,
CSS fade, or manually sampled blur above the WebView. Explicit soft and hard
edge styles belong to macOS 26.1 titlebar or split-item accessory controllers,
not to the standard toolbar itself.

Expose those accessory-only styles through a type-checked non-owning wrapper
around a native `NSTitlebarAccessoryViewController` or
`NSSplitViewItemAccessoryViewController`. Validate the Objective-C class before
constructing the Go wrapper. Map Automatic, Soft, and Hard directly to
`NSScrollEdgeEffectStyle.automaticStyle`, `.softStyle`, and `.hardStyle` under
an `@available(macOS 26.1, *)` guard. Automatic is a successful no-op on older
systems; explicit styles return a distinguishable availability error. Dispatch
all Objective-C property access synchronously to the AppKit main thread and do
not retain the borrowed controller.

The primary-pane host must paint `NSColor.windowBackgroundColor` across its
full bounds for `MacBackdropNormal`. Wails windows begin transparent at the
native root to support explicit transparent backdrop modes, so leaving a
normal split window clear exposes the desktop in the toolbar region and makes
the semantic sidebar materially more transparent than Finder. Restore an
opaque native window background for `MacBackdropNormal`; preserve the clear
root for transparent, translucent, and Liquid Glass modes.

Keep the sidebar's outline and scroll views transparent. The
`NSSplitViewItem` returned by `sidebarWithViewController:` owns the standard
sidebar material; painting another background inside the source list would
discard native vibrancy and appearance behavior.

Do not simulate the overlap boundary with an HTML gradient or a fixed titlebar
height. Both approaches break when toolbar display mode, accessibility text
size, fullscreen state, or macOS appearance changes. Page typography may still
use normal document margins; those margins are not a replacement for native
window layout.

Existing window methods continue to address the one primary WebView:

- `window.SetURL`
- `window.Reload`
- `window.ExecJS`
- runtime bridge messages
- drag and drop
- transparency/backdrop configuration

The `MacSplitWebviewPane` handle delegates those operations to the window. Its
navigation generation increments synchronously when navigation starts. A
completion event captures the current generation before leaving AppKit; stale
completions are ignored.

## 8. Attach the toolbar after the split

The sidebar tracking separator discovers a sidebar divider only if the split
view is already in the same full-size-content window. Installation order must
therefore be:

```text
window/titlebar preferences
native split controller
native toolbar
show window
```

Use AppKit's standard identifiers for the sidebar toggle and tracking
separator. They return no Wails item handles because their actions and native
instances are AppKit-owned.

Apply `MacToolbar.SetDisplayMode` to the detached candidate and replay its
latest value after committing the toolbar, matching live item setter race
handling.

For a Finder-like demo use:

```go
Mac: application.MacWindow{
    Backdrop:      application.MacBackdropNormal,
    ContentLayout: application.MacContentLayoutEdgeToEdge,
    TitleBar: application.MacTitleBar{
        AppearsTransparent:   false,
        FullSizeContent:      true,
        HideToolbarSeparator: true,
        ToolbarStyle:         application.MacToolbarStyleUnified,
    },
}

toolbar := application.NewMacToolbar().
    SetDisplayMode(application.MacToolbarDisplayModeIconOnly)
```

These are explicit native window and toolbar preferences, not a frameless
simulation.

The scroll-edge effect must remain system-owned and pane-local. Use a normal
opaque window, a nontransparent titlebar, `FullSizeContent`, and edge-to-edge
primary-pane constraints. On macOS 26, WebKit and AppKit then construct the
adaptive fade and blur internally. Do not inspect private WebKit view classes,
create an `NSVisualEffectView` overlay, or imitate the effect in CSS.

Prevent diagonal trackpad gestures from shifting a document with no
horizontal content by setting `overflow-x: hidden` and
`overscroll-behavior-x: none` in that document's CSS. This is a WebKit content
policy, not an `NSScrollView` preference: public macOS WebKit does not expose
the internal scrolling object through an API suitable for Wails.

## 9. Teardown

Before destroying the native window:

1. clear the window's split pointer;
2. mark the native owner uninstalled;
3. remove Go pane and sidebar-item registry entries;
4. mark pane, sidebar, and item handles dead;
5. remove collapse KVO;
6. clear the primary WebView association;
7. release sidebar controllers, records, and the split owner once; and
8. allow normal window destruction to release the preserved primary WebView.

Repeated teardown must be safe. Calls on dead handles must not mutate stored
state, enqueue work, or retain callbacks.

## 10. Tests

Platform-neutral tests should cover:

- split validation, ordering, ownership, and freeze behavior;
- generated sidebar IDs and declaration order;
- root items, sections, item snapshots, and selection;
- native selection callback lookup and callback clearing;
- live item state and dead-handle behavior;
- dimension validation and zero reset semantics;
- collapse callback deduplication;
- primary navigation generation handling;
- toolbar display-mode defaults and invalid values;
- duplicate standard sidebar toolbar items; and
- example asset serving without a sidebar document.

Run focused tests under the race detector. Compile desktop macOS natively and
compile Windows plus Linux/server targets to verify stubs.

## 11. Manual acceptance

On a supported macOS release:

1. Compare the running demo beside Finder.
2. Confirm the sidebar material remains visible behind all rows.
3. Confirm the accessibility hierarchy reports native outline/table rows, not
   web content.
4. Select rows with mouse and arrow keys; the primary editor must update.
5. Add a note; a native row must appear and become selected.
6. Search; native rows must filter without creating a page navigation.
7. Toggle the sidebar from toolbar and View menu.
8. Drag the divider and relaunch; autosaved width must return.
9. Enter focus mode and expand the sidebar manually; focus state must reset.
10. Confirm logs contain an editor asset request and no sidebar asset request.

The implementation is not acceptable if the visual result depends on an
opaque or transparent sidebar WKWebView.

## 12. Common failure modes

- **HTML inside a sidebar split item:** native container, non-native content.
- **Transparent WebView workaround:** still wrong ownership, focus,
  accessibility, and selection semantics.
- **Opaque outline/scroll backgrounds:** hides AppKit sidebar material.
- **Tracking separator attached first:** AppKit cannot discover the divider.
- **User-supplied row IDs:** leaks callback plumbing into application APIs.
- **Callback under a model lock:** deadlocks live setters called by callbacks.
- **Selection callback during reload:** falsely reports programmatic changes as
  user actions.
- **Discarding the old content controller:** loses backdrops and overlays.
