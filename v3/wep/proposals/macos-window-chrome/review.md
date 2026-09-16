# Native macOS chrome branch review

## Scope

This review covers the branch-specific work from `7b4823680` through
`de507f266`, plus the local generated-template correction that changes the
default macOS backdrop from translucent to normal. The review considers the
public Go API, AppKit implementation, examples, tests, build boundaries, and
the next public AppKit chrome worth exposing.

## Executive assessment

The branch now contains a credible native macOS application structure rather
than a collection of visual effects around a WebView. The strongest parts are:

- one toolbar model shared by WebView and native windows;
- semantic AppKit split items instead of simulated HTML panes;
- an actual `NSOutlineView` source-list sidebar;
- an AppKit inspector with live typed controls;
- a primary WKWebView that can be below-toolbar or edge-to-edge;
- a WebView-free experimental `NativeWindow` and `NSTextView` editor; and
- lazy, multi-representation native sharing.

The branch is not yet merge-ready. Construction-time split configuration is
unusable for windows created after `App.Run`, the `server` build is broken,
and `NativeWindowOptions.Mac` promises substantially more behavior than the
native implementation applies. Those are release blockers; another native
control should not be added before they are resolved.

## Native chrome currently available

### Window and titlebar configuration

Existing `WebviewWindow` support now composes with:

- normal, transparent, translucent, and whole-window Liquid Glass backdrops;
- expanded, preference, unified, and unified-compact toolbar styles;
- standard title visibility/transparency/full-size-content behavior;
- below-toolbar and edge-to-edge primary content layout;
- standard AppKit window tabbing mode; and
- existing window levels, collection behavior, appearance, corners, and
  fullscreen behavior.

The normal opaque backdrop is the correct default for ordinary document
windows. A page may still paint any design it wants. The generated template's
previous translucent choice was redundant because that page paints an opaque,
full-window image and gradients.

### `MacToolbar`

The toolbar API exposes real `NSToolbar` content with generated internal
identifiers and typed live handles:

- buttons with callbacks and SF Symbols;
- `NSToolbarItemGroup` segmented groups;
- `NSSearchToolbarItem` with a compatibility fallback;
- `NSSharingServicePickerToolbarItem` with lazy plain-text, HTML, PDF, PNG,
  JPEG, or custom UTI representations;
- flexible space;
- standard sidebar and inspector toggles;
- standard sidebar and inspector tracking separators;
- icon/label display modes; and
- live label, symbol, tooltip, enabled, hidden, bordered, prominent, tint,
  badge, selection, sharing metadata, and callback updates.

This API follows the right identity model: application code retains an item
handle and does not invent an ID merely so Wails can route callbacks.

### `MacSplitView`, source-list sidebar, and inspector

`MacSplitView` uses `NSSplitViewController` and semantic
`NSSplitViewItem` roles. It supports:

- a native leading source-list sidebar;
- exactly one primary WKWebView for `WebviewWindow`, or one native text editor
  for `NativeWindow`;
- a native trailing inspector;
- persisted divider positions;
- minimum, maximum, preferred-fraction, and holding-priority sizing;
- collapse policy, state, toolbar integration, and callbacks; and
- live below-toolbar/edge-to-edge layout changes for WKWebView content.

The sidebar supports sections, root rows, SF Symbols, selection, enabled and
hidden state, tooltips, callbacks, and live snapshot updates. The inspector
supports labels, text fields, checkboxes, popup buttons, sections, live state,
and typed callbacks.

### Scroll-edge style interop

`MacAccessoryViewController` type-checks an externally supplied
`NSTitlebarAccessoryViewController` or
`NSSplitViewItemAccessoryViewController` pointer and exposes automatic, soft,
and hard scroll-edge styles where macOS supports them. This is currently a
low-level interoperability escape hatch, not a complete Wails accessory API:
the branch does not create or attach either controller type for an ordinary Go
application.

### Experimental native windows

`NativeWindow` creates an `NSWindow` without WKWebView. It can host the same
toolbar and a split layout containing the source-list sidebar, inspector, and
an `NSTextView` editor. `Options.NativeOnly` and the `wails_native` build tag
also provide a useful compile-time experiment that removes WebKit and Web-only
runtime machinery.

This work is valuable evidence for a v4 content/window separation, but the v3
surface is intentionally temporary and materially less complete than
`WebviewWindow`.

## Review findings

### [P1] Split windows cannot be constructed after `App.Run`

Both window managers call `runOrDeferToAppRun` inside `NewWithOptions`. Once
the app is running, that method calls `Run` synchronously before returning the
new window. `SetSplitView` is a separate call and rejects any window whose
native implementation already exists.

Consequences:

- a tray callback cannot create and then configure a native editor window;
- a running app cannot create a new WebView window with a native sidebar;
- `NativeWindow.New()` is unusable because a native window requires a split
  layout but `New()` cannot supply one; and
- `NativeWindow.Run` logs and closes the window before the caller has a chance
  to install content.

The API needs an atomic construction path. The least disruptive v3 repair is
to add `SplitView` and optionally `Toolbar` to the window options, validate and
claim them before scheduling `Run`, and retain the setters for pre-`App.Run`
compatibility. A delayed asynchronous `Run` would only replace a deterministic
failure with a race.

### [P1] `server` builds no longer compile

`webview_requests.go` is guarded only by `!wails_native` and unconditionally
calls `webview.NewRequest`. Under `-tags server`, the asset-server WebView
package intentionally has no platform `NewRequest` implementation.

Reproduced with:

```sh
GOWORK=off CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -tags server ./pkg/application
```

The build fails at `pkg/application/webview_requests.go:23` with
`undefined: webview.NewRequest`. The WebView request implementation needs to
be excluded from `server`, with a server-compatible stub or with the consumer
removed from that build graph.

### [P1] `NativeWindowOptions.Mac` silently ignores most advertised options

The public comment says `NativeWindow` reuses existing macOS window chrome
options. The implementation currently applies only a titlebar subset:

- `AppearsTransparent`;
- `FullSizeContent`;
- `HideTitle`;
- `HideToolbarSeparator`; and
- `ToolbarStyle`.

The native window is always created opaque with
`NSColor.windowBackgroundColor`. `Backdrop` is only reduced to a boolean passed
to split installation; transparent, translucent, and Liquid Glass native
windows are not actually configured. Other ignored fields include shadow,
corners, appearance, window level, collection behavior, tabbing mode, Liquid
Glass configuration, and content-layout policy.

Either implement the documented subset consistently or replace `Mac MacWindow`
with a deliberately smaller `NativeMacWindowOptions` until v4. Silent partial
acceptance is the least reliable contract.

### [P2] Toolbar structure becomes silently stale after attachment

Item setters update an installed toolbar, but every `Add*` method only appends
to the Go model. Adding a button, group member, space, toggle, or separator
after installation does not insert it into the active `NSToolbar`, and no
error tells the caller that the native structure is unchanged.

Choose one explicit contract:

1. freeze toolbar structure when claimed and reject later additions; or
2. expose native insert/remove/move operations and keep the model synchronized.

AppKit already provides identifier-based insertion and removal. Typed handles
can preserve the no-public-ID API.

### [P2] Accessory style API has no Go-native producer

The scroll-edge wrapper accepts an unsafe pointer, but Wails no longer exposes
the titlebar-accessory creation API that existed in the branch's first draft,
and it does not expose split-item accessories. A normal Go application cannot
obtain a supported controller without writing its own Objective-C bridge.

This is a useful escape hatch, but it should not be presented as complete
accessory-controller support.

### [P2] Native window failures are log-only

`NativeWindow.Run()` has no error result. Creation, split installation, and
toolbar attachment failures are logged and the window is closed. This makes
recovery, tests, and tray-driven creation difficult. Atomic option validation
would catch most failures before scheduling native creation; remaining native
failures should be observable through a callback, retained error, or a future
returning construction API.

## Recommended next work

### 0. Stabilize the branch

Before expanding the API:

1. repair `server` builds;
2. provide atomic split/toolbar construction for windows created at runtime;
3. make the native-window option contract truthful;
4. define and enforce toolbar structural mutation semantics; and
5. add tests that create split WebView and native windows after `App.Run`.

### 1. Add first-class titlebar and split-item accessories

This is the strongest next native-chrome feature. AppKit places
`NSTitlebarAccessoryViewController` in the titlebar/toolbar region and gives
it system blur and fullscreen behavior. Modern `NSSplitViewItem` also supports
top- and bottom-aligned `NSSplitViewItemAccessoryViewController` instances,
which are the correct place for a Finder-style search/filter strip above a
sidebar or content list.

The first useful Go surface should provide native controls rather than a tiny
accessory WKWebView:

- search field;
- segmented control;
- button or menu button;
- label/status item; and
- custom-native escape hatch.

The returned typed accessory handle can directly expose hidden state,
automatic content insets, and preferred scroll-edge style. This makes the
existing `MacAccessoryViewController` behavior useful without unsafe pointers.

Apple references:

- [NSTitlebarAccessoryViewController](https://developer.apple.com/documentation/appkit/nstitlebaraccessoryviewcontroller)
- [NSSplitViewItem accessories](https://developer.apple.com/documentation/appkit/nssplitviewitem)
- [NSSplitViewItemAccessoryViewController](https://developer.apple.com/documentation/appkit/nssplitviewitemaccessoryviewcontroller)

### 2. Complete the toolbar model

High-value additive toolbar work:

- `NSMenuToolbarItem`, reusing Wails' existing menu-item model;
- visibility priority and navigational semantics;
- centered and selected item handles;
- live insert/remove/move;
- search recents and search-menu templates; and
- user customization plus autosave.

Customization is the case where stable identifiers become application
semantics rather than internal callback plumbing. Preserve generated callback
IDs, but allow an optional stable persistence key on the toolbar and
customizable items. Do not force IDs on applications that do not enable
customization.

Apple references:

- [NSToolbar](https://developer.apple.com/documentation/appkit/nstoolbar)
- [NSMenuToolbarItem](https://developer.apple.com/documentation/appkit/nsmenutoolbaritem)
- [NSToolbarItem standard identifiers](https://developer.apple.com/documentation/appkit/nstoolbaritem/identifier)

### 3. Add the semantic content-list pane

AppKit has a distinct `contentListWithViewController` split-item role. It is
the missing middle column for Finder, Mail, and document-browser layouts:

```text
source-list sidebar | content list | primary document | inspector
```

This should be a native table/outline API with selection and updates, not a
second WKWebView. Design it separately from `MacSidebar`; content lists need
titles, subtitles, badges/accessories, multi-selection, sorting, and often
context menus.

### 4. Enrich existing native panes

Useful additive sidebar and inspector capabilities:

- sidebar hierarchy beyond one section level;
- badges and trailing accessories;
- context menus;
- drag/drop and reordering;
- inline rename;
- inspector sliders, steppers, segmented controls, color wells, buttons, date
  pickers, and disclosure groups; and
- removal/reordering APIs for sections, rows, and controls.

### 5. Explicit window-tab management

Wails exposes `NSWindow.tabbingMode`, but not tab groups themselves. Additive
methods could group windows, select next/previous tabs, detach a tab, toggle
the tab bar/overview, and inspect the group. This is preferable to inventing
an `NSTabViewController` window abstraction unless the product specifically
needs in-window preference tabs rather than native document-window tabs.

### 6. Adjacent native content, after chrome

PDFKit, Quick Look, sheets, popovers, and `NSDocument` integration are valuable
but are content/document architecture rather than the next window-chrome
primitive. They should follow the window/content model decision rather than
drive it accidentally.

## Verification performed

Passing:

```text
GOWORK=off go test -race ./pkg/application
GOWORK=off go test ./examples/mac-toolbar ./internal/templates
GOWORK=off go vet ./pkg/application ./examples/mac-toolbar ./examples/mac-native-editor
GOWORK=off go build -tags 'production wails_native' ./examples/mac-native-editor
```

A fresh vanilla project was generated from the changed template, bindings were
generated, its frontend was built, and its macOS binary compiled successfully.
The generated `main.go` contains `MacBackdropNormal` while retaining the same
page design and matching `BackgroundColour`.

Failing:

```text
GOWORK=off CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -tags server ./pkg/application
```

The linker warnings in local macOS tests are caused by the machine's macOS 26
SDK objects being linked with the project's macOS 11 deployment target; they
did not fail the builds.

