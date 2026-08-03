# Wails Enhancement Proposal (WEP)

## Native macOS Window Chrome: Toolbar, Split Panes, Titlebar Accessories, and Tab Windows

**WEP Number**: (leave blank, assigned on acceptance)
**Status**: Draft
**Author**: Lea Anthony
**Created**: 2026-08-03
**Discussion**: (none yet)
**Implementor**: Lea Anthony (reference implementation drafted with Claude assistance)
**Target**: Wails v3

## Summary

Adds four independent, macOS-only additions to `v3/pkg/application` that expose native Cocoa
window-chrome constructs Wails currently has no path to:

- **`MacToolbar`** — real `NSToolbarItem`s (buttons, a segmented group, a search field, flexible
  space, a sidebar toggle) via `WebviewWindow.SetToolbar`, replacing today's empty-spacer
  `NSToolbar` created by `MacTitleBar.UseToolbar`.
- **`MacSplitWindow`** — a native `NSSplitViewController` layout of independent panes (sidebar,
  content, inspector, ...), each either its own `WKWebView` or a native content kind (text editor,
  PDF, Quick Look preview).
- **Titlebar accessories** — `WebviewWindow.AddTitlebarAccessory`, a small control docked in the
  titlebar itself via `NSTitlebarAccessoryViewController`, independent of the toolbar and of any
  split layout.
- **`MacTabWindow`** — an `NSTabViewController`-backed window switching between whole-page tabs,
  distinct from Wails' existing OS-level window-tabbing (`MacWindowTabbingMode`, which merges
  separate whole windows into one titlebar's tab strip).

Full API surface is in [Detailed Design](#detailed-design) below.

## Motivation

Today `WebviewWindowOptions.Mac` covers window-level chrome (titlebar style, backdrop, corner
radius, appearance) but stops at the boundary of "one window, one webview, filling the content
area." Cocoa offers several well-established constructs for structuring a window's content and
titlebar beyond that boundary — `NSToolbar` with real items, `NSSplitViewController` for
sidebar/content/inspector layouts, `NSTitlebarAccessoryViewController` for lightweight titlebar
controls, `NSTabViewController` for in-window tabs — and apps that want a native macOS feel
(Mail, Notes, Xcode, Finder) lean on all four. Wails currently gives an app author no way to
reach any of them; the only related existing surface (`MacTitleBar.UseToolbar`) creates a
toolbar with no items at all, purely for its effect on titlebar metrics.

This is explicitly **not** the same problem `WebviewPanel` (#4880) solves: a `WebviewPanel`
composites a floating webview into an existing content area at an absolute x/y offset, z-ordered
above other content. A split window's entire content *is* a native split layout — peer panes with
a real draggable divider that AppKit owns the geometry of, not a floating overlay. Reusing
`WebviewPanel` for split panes would leave methods like `SetBounds`/`SetPosition`/`SetZIndex`/
`DockLeft` on a pane whose geometry the split view actually owns, which is why this proposes a
distinct type for panes rather than extending `WebviewPanel`.

## Detailed Design

Two subsystems are independent of everything else (`MacToolbar`, titlebar accessories); two
build directly on `MacSplitWindow`'s "distinct type, no window-level URL/ExecJS" reasoning
(native pane content kinds, `MacTabWindow`). All four follow the codebase's existing id-tagged
native object + narrow `//export` bridge + channel + dispatch pattern already used for native
menu items, so no new plumbing style is introduced.

### Toolbar

```go
type MacToolbar struct {
	Items []MacToolbarItem
}

type MacToolbarItemKind int

const (
	ToolbarButton MacToolbarItemKind = iota
	ToolbarGroup                     // renders as NSToolbarItemGroup
	ToolbarSearchField                // NSSearchToolbarItem
	ToolbarFlexibleSpace
	ToolbarSidebarToggle
)

type MacToolbarItem struct {
	ID         string // NSToolbarItem identifier
	Kind       MacToolbarItemKind
	Label      string
	SymbolName string // SF Symbol name
	Bordered   bool
	Prominent  bool
	TintColor  *RGBA
	BadgeCount int // 0 = no badge

	Items []MacToolbarItem // for Kind == ToolbarGroup

	// TracksSplitPane, for Kind == ToolbarSidebarToggle: informational only
	// in the current implementation. AppKit's own toggleSidebar: responder-
	// chain action already finds and toggles the split view's sidebar item
	// with zero extra wiring for the common single-sidebar case, so this
	// field isn't consumed by anything yet.
	TracksSplitPane string

	OnClick  func(*Context)         // required for ToolbarButton
	OnSearch func(*Context, string) // required for ToolbarSearchField
}

func (w *WebviewWindow) SetToolbar(toolbar *MacToolbar) Window
```

A button-kind item with no `OnClick` (or a search item with no `OnSearch`) fails the whole
`SetToolbar` call rather than rendering as an inert button; the previous toolbar (if any) is left
in place and the error is delivered via `Window.Error`. There is no per-item update method —
changing anything about an existing item (e.g. a badge count) means rebuilding the `MacToolbar`
and calling `SetToolbar` again.

### Split windows

A pane's content isn't always a webview: `NSSplitViewItem` wraps an `NSViewController`, and its
view can be any native content. `Content` is a sealed interface.

```go
type MacSplitOrientation int

const (
	SplitHorizontal MacSplitOrientation = iota // panes left-to-right, vertical divider
	SplitVertical                              // panes top-to-bottom, horizontal divider
)

type MacSplitPaneBehavior int

const (
	PaneBehaviorDefault MacSplitPaneBehavior = iota
	PaneBehaviorSidebar
	PaneBehaviorContentList
	PaneBehaviorInspector
)

// MacPaneContent is a sealed interface: MacWebviewContent is implemented;
// MacTextEditorContent / MacPDFContent / MacQuickLookContent (below) are
// proposed but not yet built -- see Reference Implementation.
type MacPaneContent interface {
	isMacPaneContent()
}

type MacWebviewContent struct {
	URL                string
	WebviewPreferences MacWebviewPreferences // not yet wired per-pane; every pane gets the same defaults
}

func (MacWebviewContent) isMacPaneContent() {}

// MacTextEditorContent: a native NSTextView (inside the usual NSScrollView),
// for plain/rich text, not a code-editor replacement. For syntax
// highlighting or LSP-backed autocomplete, use MacWebviewContent pointed at
// a bundled Monaco/CodeMirror build instead, that already works with no new
// API. OnChange fires with the full current text on each change; debouncing
// or diffing is left to the caller.
type MacTextEditorContent struct {
	Text     string // initial content
	Editable bool
	OnChange func(*Context, string)
}

func (MacTextEditorContent) isMacPaneContent() {}

// MacPDFContent: PDFKit.PDFView. Requires linking PDFKit (or the Quartz
// umbrella framework that includes it), a new cgo LDFLAGS entry.
type MacPDFContent struct {
	FilePath    string // exactly one of FilePath or Data
	Data        []byte
	InitialPage int // 0 = first page
}

func (MacPDFContent) isMacPaneContent() {}

// MacQuickLookContent: QLPreviewView (QuickLookUI, part of the Quartz
// umbrella framework on macOS), a universal file preview for any type
// macOS already knows how to preview, no per-format code. Requires the
// same new LDFLAGS entry as MacPDFContent.
type MacQuickLookContent struct {
	FilePath string // required
}

func (MacQuickLookContent) isMacPaneContent() {}

type MacSplitPaneOptions struct {
	Name     string // unique within the window; used for Pane(name)
	Content  MacPaneContent
	Behavior MacSplitPaneBehavior

	MinThickness               float64
	MaxThickness               float64
	PreferredThicknessFraction float64
	HoldingPriority             float64 // <= 0 leaves AppKit's own default
	Collapsible                  bool
	StartCollapsed                bool
}

type MacSplitWindowOptions struct {
	Title  string
	Width  int
	Height int

	Orientation  MacSplitOrientation
	AutosaveName string // NSSplitView.autosaveName; persists divider positions across launches
	Panes        []MacSplitPaneOptions // at least 1; pane 0 reuses the window's own webview

	TitleBar           MacTitleBar
	Appearance         MacAppearanceType
	Backdrop           MacBackdrop // also covers what an earlier draft of this proposal had as a
	                                // separate LiquidGlass field: Backdrop already has a
	                                // MacBackdropLiquidGlass value, so a second field would have
	                                // been redundant with the existing WebviewWindowOptions shape.
	WindowLevel        MacWindowLevel
	CollectionBehavior MacWindowCollectionBehavior
}

// NewSplitWindow creates a native split-pane window: a real *WebviewWindow
// underneath (so it gets full window chrome -- toolbar, standard titlebar
// behaviour, etc. -- for free), with its content area replaced by an
// NSSplitViewController. Errors if Panes is empty, if two panes share a
// Name, or on non-macOS builds.
func NewSplitWindow(options MacSplitWindowOptions) (*MacSplitWindow, error)

// MacSplitPane exposes only what's common to every content kind: identity
// and collapse state, both properties of the pane's slot in the split view,
// not of its content. Content-specific operations (URL, text, page number)
// live on the typed handle Content() returns, because a pane holding a PDF
// has no URL to set and a pane holding a webview has no page number.
type MacSplitPane struct{ /* unexported */ }

func (p *MacSplitPane) ID() uint
func (p *MacSplitPane) Name() string
func (p *MacSplitPane) IsCollapsed() bool
func (p *MacSplitPane) SetCollapsed(collapsed bool)

// Content returns the pane's typed content handle. Currently only
// *MacWebviewPane is implemented; the other three content kinds' typed
// handles (below) are proposed, not yet built.
func (p *MacSplitPane) Content() any

// SetContent (proposed, not yet built): replaces the pane's content in
// place -- the old native view is torn down and the new one takes over the
// same NSSplitViewItem, so the divider position and the pane's slot in the
// layout are untouched. Any handle previously returned by Content() is
// invalid after this call, fetch it again.
func (p *MacSplitPane) SetContent(content MacPaneContent) error

// MacWebviewPane: the only content kind with a JS runtime. Implemented.
type MacWebviewPane struct{ /* unexported */ }

func (p *MacWebviewPane) SetURL(url string)
func (p *MacWebviewPane) Reload()
func (p *MacWebviewPane) ExecJS(js string)

// MacTextEditorPane / MacPDFPane / MacQuickLookPane (proposed, not yet built):
type MacTextEditorPane struct{ /* unexported */ }

func (p *MacTextEditorPane) Text() string
func (p *MacTextEditorPane) SetText(text string) *MacTextEditorPane
func (p *MacTextEditorPane) SetEditable(editable bool) *MacTextEditorPane

type MacPDFPane struct{ /* unexported */ }

func (p *MacPDFPane) CurrentPage() int
func (p *MacPDFPane) GoToPage(page int) *MacPDFPane
func (p *MacPDFPane) SetDocument(filePath string) *MacPDFPane // swaps the displayed document

type MacQuickLookPane struct{ /* unexported */ }

func (p *MacQuickLookPane) SetFilePath(path string) *MacQuickLookPane

// MacSplitWindow: window chrome only (position, show/hide, appearance).
// No URL/ExecJS/SetURL: there is no single webview to apply them to.
// Implemented.
type MacSplitWindow struct{ /* unexported */ }

func (w *MacSplitWindow) Show() *MacSplitWindow
func (w *MacSplitWindow) Hide() *MacSplitWindow
func (w *MacSplitWindow) Close()
func (w *MacSplitWindow) Focus() *MacSplitWindow
func (w *MacSplitWindow) SetPosition(x, y int) *MacSplitWindow
func (w *MacSplitWindow) SetSize(width, height int) *MacSplitWindow
func (w *MacSplitWindow) Window() *WebviewWindow // the real underlying window, e.g. for SetToolbar
func (w *MacSplitWindow) Panes() []*MacSplitPane
func (w *MacSplitWindow) Pane(name string) (*MacSplitPane, bool)

// AddPane/RemovePane (proposed, not yet built): construction-time panes
// only so far.
func (w *MacSplitWindow) AddPane(options MacSplitPaneOptions, atIndex int) (*MacSplitPane, error)
func (w *MacSplitWindow) RemovePane(pane *MacSplitPane) error
```

### Titlebar accessory

Lighter-weight than the toolbar: a single small control docked in the titlebar strip, for windows
that don't want to adopt a full `NSToolbar`. Follows the same addressable-webview pattern as split
panes rather than inventing a native-control mini-API: an accessory hosts its own small webview,
sized to fixed dimensions, so anything renderable in HTML/CSS/JS works without a second widget
system. Implemented for `WebviewWindow`; applies equally to `MacSplitWindow`/`MacTabWindow` since
both wrap a real `*WebviewWindow` for their chrome.

```go
type MacTitlebarAccessoryPosition int

const (
	AccessoryLeading  MacTitlebarAccessoryPosition = iota // NSTitlebarAccessoryViewController.layoutAttribute = .leading
	AccessoryTrailing                                     // .trailing
)

type MacTitlebarAccessoryOptions struct {
	Name     string // unique per window
	URL      string
	Position MacTitlebarAccessoryPosition
	Width    float64 // defaults to 180 if <= 0
	Height   float64 // defaults to 28 if <= 0; oversized content is clipped
}

func (w *WebviewWindow) AddTitlebarAccessory(options MacTitlebarAccessoryOptions) (*MacTitlebarAccessory, error)
func (w *WebviewWindow) RemoveTitlebarAccessory(accessory *MacTitlebarAccessory) error

type MacTitlebarAccessory struct{ /* unexported */ }

func (a *MacTitlebarAccessory) ID() uint
func (a *MacTitlebarAccessory) Name() string
func (a *MacTitlebarAccessory) SetURL(url string)
func (a *MacTitlebarAccessory) ExecJS(js string)
```

An accessory's webview always has `underPageBackgroundColor` cleared: unlike a split pane, it
always sits directly on the titlebar's own background rather than optionally against a
material-backed `NSVisualEffectView`/`NSGlassEffectView`, so there's no "opaque content" case
the way there is for a `PaneBehaviorDefault` split pane.

### In-window tabs (proposed, not yet built)

Distinct from both native OS window-tabbing (`MacWindow.TabbingMode`, already in Wails, an OS tab
strip merging separate whole windows) and split panes (which show every pane at once):
`NSTabViewController` shows exactly one tab's content at a time, switching swaps it. Same
addressable-webview identity pattern as split panes and titlebar accessories, and the same
"distinct type, no window-level URL/ExecJS" reasoning as `MacSplitWindow`.

```go
type MacTabStyle int

const (
	TabStyleSegmented MacTabStyle = iota // NSTabViewController.tabStyle = .segmentedControlOnTop
	TabStyleToolbar                      // .toolbar (macOS 26+): tabs render inside the unified toolbar itself.
	                                      // Interaction with SetToolbar's own items on the same window needs
	                                      // verifying against real AppKit behavior before this ships; flagged,
	                                      // not assumed to just compose. This is what example 5 is for.
)

type MacTabOptions struct {
	Name       string // unique; used for Tab(name)
	Title      string // tab label
	SymbolName string // SF Symbol shown on the tab button
	URL        string

	WebviewPreferences MacWebviewPreferences
}

type MacTabWindowOptions struct {
	Title  string
	Width  int
	Height int

	TabStyle MacTabStyle
	Tabs     []MacTabOptions // at least 1

	TitleBar           MacTitleBar
	Appearance         MacAppearanceType
	Backdrop           MacBackdrop
	WindowLevel        MacWindowLevel
	CollectionBehavior MacWindowCollectionBehavior
}

func NewTabWindow(options MacTabWindowOptions) (*MacTabWindow, error)

type MacTab struct{ /* unexported */ }

func (t *MacTab) ID() uint
func (t *MacTab) Name() string
func (t *MacTab) SetURL(url string) *MacTab
func (t *MacTab) Reload() *MacTab
func (t *MacTab) ExecJS(js string)
func (t *MacTab) SetTitle(title string) *MacTab // NSTabViewItem.label

// MacTabWindow: window chrome only, same reasoning as MacSplitWindow.
type MacTabWindow struct{ /* unexported */ }

func (w *MacTabWindow) Show() *MacTabWindow
func (w *MacTabWindow) Hide() *MacTabWindow
func (w *MacTabWindow) Close()
func (w *MacTabWindow) Focus() *MacTabWindow
func (w *MacTabWindow) SetPosition(x, y int) *MacTabWindow
func (w *MacTabWindow) SetSize(width, height int) *MacTabWindow
func (w *MacTabWindow) Tabs() []*MacTab
func (w *MacTabWindow) Tab(name string) (*MacTab, bool)
func (w *MacTabWindow) SelectedTab() *MacTab
func (w *MacTabWindow) SelectTab(tab *MacTab) *MacTabWindow
func (w *MacTabWindow) AddTab(options MacTabOptions, atIndex int) (*MacTab, error)
func (w *MacTabWindow) RemoveTab(tab *MacTab) error
```

### Composed example

```go
win, _ := application.NewSplitWindow(application.MacSplitWindowOptions{
	Orientation:  application.SplitHorizontal,
	AutosaveName: "main.split",
	Panes: []application.MacSplitPaneOptions{
		{Name: "sidebar", Behavior: application.PaneBehaviorSidebar,
			Content:      application.MacWebviewContent{URL: "/sidebar.html"},
			MinThickness: 180, MaxThickness: 320, Collapsible: true},
		{Name: "content", Behavior: application.PaneBehaviorContentList,
			Content: application.MacWebviewContent{URL: "/content.html"}},
	},
})
contentPane, _ := win.Pane("content")
content := contentPane.Content().(*application.MacWebviewPane)

win.Window().SetToolbar(&application.MacToolbar{
	Items: []application.MacToolbarItem{
		{ID: "sidebar", Kind: application.ToolbarSidebarToggle},
		{ID: "flex", Kind: application.ToolbarFlexibleSpace},
		{ID: "search", Kind: application.ToolbarSearchField, OnSearch: func(ctx *application.Context, q string) {
			content.SetURL(searchURLFor(q))
		}},
	},
})

win.Show()
```

## Non-Goals

- `WebviewPanel` (#4880) is untouched, a separate feature (see Motivation for why it isn't reused
  here).
- Windows and Linux: this proposal is macOS-only throughout; every new type/function is a no-op
  or returns an error on every other platform.
- Native `NSOutlineView` sidebar lists: would need a data-source bridging API, not just a content
  kind. A sidebar pane today is a webview or (proposed) one of the other content kinds, rendering
  its own list in HTML if it wants one.
- Sheets, popovers.
- Toolbar user-customization (`allowsUserCustomization` and the customization palette).

## Platform Considerations

Every addition in this proposal is macOS-only (`darwin && !ios && !server` build tag). On every
other platform:

- `SetToolbar`, `AddTitlebarAccessory`/`RemoveTitlebarAccessory` are no-ops (methods on the
  existing cross-platform `*WebviewWindow` type).
- `NewSplitWindow`/`NewTabWindow` return an error rather than a usable value.

Minimum OS versions used by specific pieces: `NSToolbarItemGroup` selection APIs and `.bordered`
(10.15), SF Symbol images (11.0), `NSSplitViewItem` sidebar/contentList factories (10.11),
inspector factory (11.0), `WKWebView.underPageBackgroundColor` (12.0, guarded — panes/accessories
simply stay opaque on older systems rather than failing). `TabStyleToolbar` targets macOS 26+
specifically, per its own note above.

## Pros/Cons

**Pros**: closes a real gap between what Cocoa offers and what Wails exposes, without which an
app author's only recourse for native chrome is cgo they write and maintain themselves outside
Wails. Each piece is additive and independently adoptable. All four reuse the existing
id-tagged-object + `//export` bridge + channel + dispatch pattern, so there's no new
architectural style for maintainers to learn. Split panes, titlebar accessories, and tab windows
all reuse the same underlying webview-construction and windowId-keyed asset/IPC routing as a
top-level `WebviewWindow`, so asset serving and JS↔Go IPC work identically everywhere without
touching the core window registry or `messageprocessor.go`.

**Cons**: four largely-independent new subsystems is a meaningfully larger surface to maintain
than a single option struct; macOS-only work always risks drifting further from Windows/Linux
parity; the native pane content kinds pull in two new linked frameworks (PDFKit, QuickLookUI).

## Alternatives Considered

- **Exposing the raw `NSWindow`/`NSView` pointer** for an app author to drive AppKit directly via
  their own cgo. Rejected: defeats the purpose of a Go-first framework, and reintroduces exactly
  the maintenance burden this proposal exists to remove.
- **A single larger `WebviewWindowOptions` field** for all of this (e.g. `Mac.SplitPanes []...`),
  construction-time only. Rejected for `MacSplitWindow`/`MacTabWindow`: both need a distinct
  window-chrome method set (no `URL`/`ExecJS` at the window level, since there's no single
  webview), which doesn't fit cleanly into the existing `WebviewWindow` type. `MacToolbar` and
  titlebar accessories *do* fit as `WebviewWindow` methods, and are exposed that way.
- **Reusing `WebviewPanel` for split panes** — see Motivation.
- **`HoldingPriority *float64`** (pointer, to distinguish "unset" from an explicit `0`) was the
  original draft's shape; changed to a plain `float64` with `<= 0` meaning "unset" during
  implementation, since 0 isn't a meaningful `NSLayoutPriority` value in practice (real values run
  roughly 1–1000) and the pointer indirection wasn't worth it for a case that can't come up.

## Backwards Compatibility

Purely additive: every new type, constant, and method is new surface, and every existing
`WebviewWindowOptions`/`MacWindow` field is untouched. No existing public API changes.

## Security and Privacy

No new data handling. Split panes, titlebar accessories, and tab windows all authenticate and
route through the same `wails://` scheme handler and windowId-keyed IPC path an ordinary
`WebviewWindow` already uses — same-origin and asset-serving behaviour is identical, not weakened
or bypassed for these new webview instances.

## Test Plan

Five example apps under `v3/examples/`, each independently runnable and demonstrating a distinct
slice of the API with no overlap between them — collectively the acceptance criteria for the
feature:

| # | Directory | App | Demonstrates | Status |
|---|---|---|---|---|
| 1 | `mac-toolbar` | Notes-style single-webview app | `MacToolbar`, all item kinds, `Bordered`/`Prominent`/`TintColor`, `BadgeCount` | Built |
| 2 | `mac-split-sidebar` | Reader app: sidebar nav, content, inspector | `NewSplitWindow`, pane behaviors, thickness/collapse options, `AutosaveName`, `ToolbarSidebarToggle` | Built |
| 3 | `mac-split-native-panes` | File browser: sidebar list, native preview pane | `MacTextEditorContent`/`MacPDFContent`/`MacQuickLookContent`, `SetContent`, `AddPane`/`RemovePane` | Not built |
| 4 | `mac-titlebar-accessory` | Plain window with a status accessory | `AddTitlebarAccessory`, both positions, no toolbar/split required | Built |
| 5 | `mac-tabview` | Settings window: General/Appearance/Advanced | `NewTabWindow`, both tab styles, `AddTab`/`RemoveTab`/`SelectTab` | Not built |

Per-example acceptance:

1. **mac-toolbar**: every item is a real `NSToolbarItem`, not a DOM button under the traffic
   lights; each fires its Go callback; the search field submits its query to Go.
2. **mac-split-sidebar**: dragging the divider works and the position survives a relaunch; the
   toolbar toggle collapses/expands the sidebar; sidebar and inspector panes get automatic glass
   without any extra code.
3. **mac-split-native-panes**: text files show a native editable pane and `OnChange` fires while
   typing; PDFs show native rendering with working page navigation; everything else falls back to
   Quick Look with zero per-format code.
4. **mac-titlebar-accessory**: both accessories render and update independently (`SetURL`/`ExecJS`)
   without touching the main window's content.
5. **mac-tabview**: both tab styles render correctly; resolves whether `TabStyleToolbar` conflicts
   with a window's own `SetToolbar` items by construction, not assumption.

Given this environment has no interactive Screen Recording permission (no way to take a real
screenshot), verification for the built examples has relied on: `go build`/`go vet` across
darwin, windows (`CGO_ENABLED=0`), linux+server (`CGO_ENABLED=0`), and a real iOS SDK
cross-compile; running the built binaries and checking the asset-server log shows every expected
distinct pane/accessory URL being requested independently; a temporary Go-side event listener
confirming JS→Go IPC round-trips correctly from a non-top-level webview; and, once, a logged walk
of the native view hierarchy to confirm a claimed fix (material behind sidebar/inspector panes)
had a real native view to act on. None of that is a substitute for actually looking at the
window — the split-sidebar example shipped once already with a real visual defect (panes weren't
actually translucent, plus leftover debug labels in the UI) that no build or log check caught,
found only when it was run and looked at directly.

## Reference Implementation

Drafted on local branch `feat/macos-window-chrome` (based on `master`), not yet pushed or opened
as a PR. Current state:

- **Built and example-verified**: `MacToolbar`, `MacSplitWindow` with `MacWebviewContent` panes,
  titlebar accessories. Examples 1, 2, and 4 above.
- **Not yet built**: native pane content kinds (`MacTextEditorContent`/`MacPDFContent`/
  `MacQuickLookContent`, `MacSplitPane.SetContent`, `AddPane`/`RemovePane`), `MacTabWindow`.
  Examples 3 and 5.

## Maintenance Plan

Follows the existing per-platform file convention already used throughout `v3/pkg/application`:
a cross-platform file with no build tag for types and dispatch, a `_darwin.go`/`.h`/`.m` trio for
the real implementation, and a same-signature no-op/error stub for every other platform. New
`.m`/`.h` files must carry their own `//go:build darwin && !ios && !server` line (not just rely
on the `_darwin` filename suffix) — a real gap found during this implementation: two `.m` files
missing that line were silently pulled into an `ios`-tagged build and failed there, invisible to
the ordinary darwin build.

## Conclusion

Closes a real, specific gap between Cocoa's window-chrome capabilities and what Wails exposes,
using patterns the codebase already has (id-tagged native objects, per-webview windowId-keyed
IPC) rather than introducing new ones. Roughly 60% built and verified as of this draft (toolbar,
split windows with webview panes, titlebar accessories); native pane content kinds and tab
windows remain.
