# Native macOS toolbar and split view

This example is a normal titled Wails window with a real AppKit `NSToolbar`
and a native three-pane `NSSplitViewController` layout. Nothing in the toolbar,
sidebar, or inspector is HTML. The dividers, pane roles, source-list rows,
property controls, keyboard navigation, collapse animation, divider-position
persistence, and toolbar alignment all come from AppKit. The editor is the
window's only WKWebView.

It demonstrates:

- a native `MacSidebar` rendered by `NSOutlineView` inside an AppKit sidebar
  split item, with generated internal row identifiers and Go callbacks;
- the window's existing and only WKWebView as primary content
  (`MacSplitView.AddPrimaryContent`);
- a semantic trailing AppKit inspector (`MacSplitView.AddInspector`) containing
  native labels, a text field, checkbox, and pop-up button;
- the standard AppKit sidebar toggle (`AddSidebarToggle`) and the sidebar
  tracking separator (`AddSidebarTrackingSeparator`) that keeps the leading
  toolbar section aligned above the moving divider;
- divider-position persistence through `SetAutosaveName`;
- native collapse observation (`OnCollapsedChange`) keeping focus mode, the
  View menu, and the frontend synchronized regardless of whether the change
  came from the toolbar, the menu, a divider gesture, or Go;
- a persistent `NSToolbarItemGroup` for Write/Preview mode;
- an `NSSearchToolbarItem` with a Go callback (and the native search-field
  fallback used on macOS 10.13–10.15);
- an `NSSharingServicePickerToolbarItem` backed by a lazy, multi-format Go
  provider (with a native picker fallback on macOS 10.13–10.14);
- native New, Save, and Focus controls;
- the standard AppKit inspector toggle and tracking separator on macOS 14+,
  with a native functional toggle fallback on older supported releases;
- generated internal item identifiers, with no IDs in application code;
- a configurable native toolbar display mode, using icon-only presentation in
  this Finder-style unified titlebar;
- live handles for application-owned items that update labels, symbols,
  tooltips, presentation, badges, selection, enabled state, and visibility
  (the two AppKit-owned sidebar items intentionally return no handle);
- a working notes editor with a native-sidebar note list, write/preview modes,
  search, native sharing, persistent local save state, editable native
  metadata, focus mode, keyboard save, and a live choice between edge-to-edge
  and below-toolbar content layout.

## Run the demo

Run it on macOS from this directory:

```sh
GOWORK=off go run .
```

The toolbar's leading section (sidebar toggle and New) stays above the native
sidebar as you drag the divider. Search and the document actions live in the
primary-content section, so the expanding native search field cannot outgrow
the sidebar. The final separator and toggle track the inspector's leading
edge. The Share PDF toolbar item opens the system share popover with only a
PDF representation. Choose Mail there and the generated Daymark PDF is
attached to a new message. The provider also contains HTML and plain-text
renderers to demonstrate a multi-format configuration, but this toolbar item
deliberately advertises only PDF so the result is unambiguous.

Use **View → Content Layout** to switch the primary WebView live between the
two native arrangements. In **Edge to Edge**, scroll the note upward: the
standard macOS 26 toolbar floats above the WebView and AppKit supplies the
scroll-edge treatment. **Below Toolbar** constrains the WebView to the
unobscured window content guide for applications that do not want overlap.

## Example layout

- `main.go` creates the application and window, attaches the split view (which
  must be configured before the window is shown), then the toolbar and menu.
- `split.go` builds the native source-list model, owns the note navigation
  state, assembles all three split panes, and wires native selection and
  collapse callbacks.
- `inspector.go` builds the native property model and connects its generated
  control callbacks to the current note.
- `toolbar.go` constructs the native controls, registers their callbacks, and
  keeps toolbar and focus state in sync with editor and sidebar events.
- `menu.go` adds keyboard and menu equivalents, including View → Toggle
  Sidebar driving the native pane.
- `share.go` implements the lazy share provider, the demo's PDF-only
  configuration, and its optional plain-text and HTML renderers.
- `pdf.go` contains the demo's independent PDF renderer.
- `assets/editor.html` is only the primary WebView writing and reading views.
  It contains no sidebar or inspector markup.
- `assets/index.html` is the embedded asset-server entry point and redirects a
  root navigation to the editor. Wails requires this file even when a window
  starts at a more specific embedded page.

The native source list, inspector, and editor communicate through Go and Wails
events (`sidebar:note-selected`, `editor:note-updated`,
`editor:inspector-state`, and `editor:ready`). There is no auxiliary page,
transparent-WebView trick, DOM inspector overlay, or HTML recreation of native
chrome.

## Content layout and Liquid Glass

Liquid Glass is not a custom Wails toolbar backdrop. A standard `NSToolbar`
built against the current macOS SDK adopts the system toolbar appearance on
supported releases, including system grouping and adaptive glass. The content
layout determines whether scrolling content can pass underneath that toolbar.

### Finder-style window recipe

Finder does not make the entire application window transparent. It combines a
normal opaque AppKit window, a unified `NSToolbar`, and content that extends
beneath the toolbar. Configure those concerns independently:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Mac: application.MacWindow{
        Backdrop:      application.MacBackdropNormal,
        ContentLayout: application.MacContentLayoutEdgeToEdge,
        TitleBar: application.MacTitleBar{
            AppearsTransparent:   false,
            FullSizeContent:      true,
            HideToolbarSeparator: true,
            ToolbarStyle:         application.MacToolbarStyleUnified,
        },
    },
})
```

Each option has a separate native responsibility:

| Option | Native result |
| --- | --- |
| `MacBackdropNormal` | Restores an opaque `NSWindow` using the configured background colour or `NSColor.windowBackgroundColor`. |
| `AppearsTransparent: false` | Keeps the standard titlebar backing. AppKit requires a nontransparent titlebar when automatically adjusting an overlapping scroll view. |
| `FullSizeContent: true` | Adds `NSWindowStyleMaskFullSizeContentView`. |
| `MacContentLayoutEdgeToEdge` | Constrains the primary WebView to the pane's full bounds instead of `NSWindow.contentLayoutGuide`. |
| `MacToolbarStyleUnified` | Uses the native unified titlebar and toolbar arrangement. |
| `HideToolbarSeparator: true` | Removes the legacy horizontal separator; AppKit's scroll-edge treatment supplies visual separation on macOS 26+. |

Do not use `MacBackdropTransparent` or `MacBackdropLiquidGlass` for this
design. Finder uses an ordinary opaque window; AppKit supplies the toolbar and
sidebar materials at their semantic locations.

The Liquid Glass behind toolbar controls and the scroll-edge fade are related
but distinct AppKit features. On macOS 26, `NSToolbar` supplies the control
glass automatically. Use `MacToolbar.AddGroup` for controls that belong on one
glass capsule; AppKit also separates unlike controls such as search fields and
segmented selectors. The fade or harder backing over moving content is owned
by the native scroll view beneath the floating controls.

On macOS 26+, `WKWebView` participates in AppKit's genuine scroll-edge
handling. With the native window structure above, the system adapts the
effect's shape to toolbar controls and progressively fades and blurs WebKit
content. There is no HTML gradient, simulated `NSVisualEffectView`, or private
WebKit integration.

AppKit does not expose a soft/hard style setter for an ordinary `NSToolbar`.
Those preferences are public only on real titlebar and split-item accessory
controllers in macOS 26.1+. Consequently Wails does not present a toolbar
style switch that AppKit cannot faithfully implement. Choose
`MacContentLayoutBelowToolbar` when overlap—and therefore the effect—is not
wanted.

Native integrations that own an `NSTitlebarAccessoryViewController` or
`NSSplitViewItemAccessoryViewController` can wrap its borrowed pointer and use
the style API directly:

```go
accessory, err := application.WrapMacAccessoryViewController(nativeController)
if err != nil {
    return err
}

if accessory.SupportsPreferredScrollEdgeEffectStyle() {
    err = accessory.SetPreferredScrollEdgeEffectStyle(
        application.MacScrollEdgeEffectStyleSoft,
    )
}
```

`MacScrollEdgeEffectStyleAutomatic`, `Soft`, and `Hard` map one-to-one to
`NSScrollEdgeEffectStyle`. The setter affects only the wrapped accessory—not
the window toolbar—and never exposes a synthetic blur-strength value. On
macOS before 26.1, `Automatic` is a successful no-op while explicit `Soft` or
`Hard` returns `ErrMacScrollEdgeEffectStyleUnavailable`. The wrapper is
non-owning, so the native window or split item must keep its controller alive.

The demo prevents accidental sideways rubber-banding using WebKit's supported
web-content policy:

```css
html, body {
    overflow-x: hidden;
    overscroll-behavior-x: none;
}
```

This keeps vertical scrolling native while preventing a diagonal trackpad
gesture from shifting a document that has no horizontal scrollable content.

### Content layout policy

A native split layout can override the window policy for its primary WebView:

```go
primary := split.AddPrimaryContent()
primary.SetContentLayout(application.MacContentLayoutEdgeToEdge)
```

The available policies are:

| Value | Behavior |
| --- | --- |
| `MacContentLayoutAutomatic` | Follows `MacWindow.ContentLayout`; at window level, follows `TitleBar.FullSizeContent`. |
| `MacContentLayoutBelowToolbar` | Constrains the WebView to `NSWindow.contentLayoutGuide`, so it never passes beneath the toolbar. |
| `MacContentLayoutEdgeToEdge` | Extends the WebView through the full content area beneath the titlebar and toolbar. |

`SetContentLayout` works before or after the split is installed. Live changes
replace the native Auto Layout constraints on the AppKit application thread.
Unsupported platforms retain the setting as inert descriptive state.

On macOS 26, standard toolbar items receive Liquid Glass automatically and
AppKit groups compatible controls. `MacToolbarItem.SetProminent` and
`SetTintColor` expose the supported item-level emphasis APIs. Do not use
`MacBackdropLiquidGlass` merely to obtain a glass toolbar: that option is a
custom whole-window backdrop for intentionally translucent application
content and is separate from `NSToolbar` presentation.

## Native inspector API

Build the control model first, retain handles that need updates, and append the
inspector after primary content so it occupies the trailing edge:

```go
inspector := application.NewMacInspector()

document := inspector.AddSection("Document")
title := document.AddTextField("Title", "A note")
category := document.AddPopup("Category", []string{"Personal", "Work"}, 0)
pinned := document.AddCheckbox("Pinned", false)

statistics := inspector.AddSection("Statistics")
words := statistics.AddLabel("Words", "0")

title.OnTextChange(func(_ *application.Context, value string) {
    model.Rename(value)
})
category.OnSelectionChange(func(_ *application.Context, index int, value string) {
    model.SetCategory(value)
})
pinned.OnToggle(func(_ *application.Context, checked bool) {
    model.SetPinned(checked)
})

inspectPane := split.AddInspector(inspector)
inspectPane.SetMinimumThickness(240).
    SetMaximumThickness(360).
    SetCollapsible(true)

toolbar.AddInspectorTrackingSeparator()
toolbar.AddInspectorToggle()
```

Wails generates section, control, and callback identifiers internally. The
returned `MacInspectorControl` is the identity and live state handle. Use
`SetValue`, `SetChecked`, `SetOptions`, `SetSelectedIndex`, `SetTooltip`,
`SetEnabled`, and `SetHidden` to update AppKit controls after installation.

Programmatic setters never invoke callbacks. User interaction updates the Go
handle before `OnTextChange`, `OnToggle`, or `OnSelectionChange` runs. This
prevents model-to-view updates from feeding back into application logic while
still making the latest user value immediately readable from the handle.

On macOS 11+, the pane is created with AppKit's semantic
`inspectorWithViewController:` role. macOS 10.13–10.15 use the same native
controls and collapse API in a regular trailing split item. The standard
Inspector toolbar identifiers are available on macOS 14+; older systems use a
native toggle fallback and omit only the unavailable tracking separator.

## Why sharing uses a provider

Wails does not guess which DOM element, editor model, URL, or rendered view an
application intends to share. The application owns that decision.

A `MacShareProvider` advertises the formats it can produce and returns bytes
only when Cocoa requests one:

```go
type MacShareProvider interface {
    ShareRepresentations() []application.MacShareRepresentation
    ShareData(application.MacShareRequest) ([]byte, error)
}
```

This maps naturally to AppKit's `NSItemProvider`:

1. `ShareRepresentations` declares the available Uniform Type Identifiers
   (UTIs), such as plain text, HTML, PDF, PNG, or an application-specific UTI.
2. The user opens the native share popover and selects a service.
3. AppKit requests one or more representations that service can consume.
4. Wails calls `ShareData` lazily on a background callback.
5. The provider returns complete bytes or an error.

An expensive PDF or image render therefore does not happen unless a receiving
service actually requests it. AppKit may request more than one advertised
format, so providers must be safe to call repeatedly and concurrently.

### Choosing a representation

The standard macOS share popover chooses a sharing service, not a file format.
When a provider advertises several representations, AppKit and the selected
service decide which compatible representation to load. An application cannot
rely on the user being offered a PDF/HTML/text picker.

When an action promises a particular format, advertise only that format for
that share item:

```go
share := toolbar.AddShare("Share PDF").
    SetProvider(application.MacShareProviderFunc{
        Available: []application.MacShareRepresentation{
            {ContentType: application.MacShareTypePDF},
        },
        Load: func(request application.MacShareRequest) ([]byte, error) {
            return renderPDF(currentDocument, request.SuggestedName)
        },
    })
```

Daymark uses this PDF-only configuration. Selecting Mail therefore causes
Wails to render PDF bytes and Mail receives the result as an attachment.

## Minimal example

`MacShareProviderFunc` is the shortest way to share simple data:

```go
toolbar := application.NewMacToolbar()

share := toolbar.AddShare("Share").
    SetProvider(application.MacShareProviderFunc{
        Available: []application.MacShareRepresentation{
            {ContentType: application.MacShareTypePlainText},
        },
        Load: func(request application.MacShareRequest) ([]byte, error) {
            return []byte("Hello from Wails"), nil
        },
    }).
    SetSubject("A Wails note").
    SetSuggestedName("Wails Note")

share.OnShared(func(_ *application.Context, service string) {
    log.Printf("shared with %s", service)
})
share.OnShareError(func(_ *application.Context, service string, err error) {
    log.Printf("%s failed: %v", service, err)
})

window.SetToolbar(toolbar)
```

No application-supplied item ID is necessary. Keep the returned `share` handle
to update metadata, replace the provider, enable or disable the item, and
register callbacks.

## Multiple formats and lazy PDF generation

A stateful provider can offer several representations of one logical item:

```go
type reportShareProvider struct {
    mu     sync.RWMutex
    report Report
}

func (p *reportShareProvider) ShareRepresentations() []application.MacShareRepresentation {
    return []application.MacShareRepresentation{
        {ContentType: application.MacShareTypePlainText},
        {ContentType: application.MacShareTypeHTML},
        {ContentType: application.MacShareTypePDF},
    }
}

func (p *reportShareProvider) ShareData(request application.MacShareRequest) ([]byte, error) {
    p.mu.RLock()
    snapshot := p.report.Clone()
    p.mu.RUnlock()

    switch request.ContentType {
    case application.MacShareTypePlainText:
        return []byte(snapshot.PlainText()), nil
    case application.MacShareTypeHTML:
        return []byte(snapshot.HTML()), nil
    case application.MacShareTypePDF:
        // renderPDF is application code. It runs only if AppKit asks for PDF.
        return renderPDF(snapshot, request.SuggestedName)
    default:
        return nil, fmt.Errorf("unsupported share type %q", request.ContentType)
    }
}
```

Install it exactly like the functional provider:

```go
provider := &reportShareProvider{report: currentReport}

share := toolbar.AddShare("Share").
    SetProvider(provider).
    SetSubject("Quarterly report").
    SetSuggestedName("Quarterly Report")
```

The returned byte slice must match the requested UTI:

| Constant | UTI | Expected bytes |
| --- | --- | --- |
| `MacShareTypePlainText` | `public.utf8-plain-text` | UTF-8 plain text |
| `MacShareTypeHTML` | `public.html` | UTF-8 HTML document or fragment |
| `MacShareTypePDF` | `com.adobe.pdf` | A complete PDF document |
| `MacShareTypePNG` | `public.png` | PNG-encoded image data |
| `MacShareTypeJPEG` | `public.jpeg` | JPEG-encoded image data |

Custom `MacShareContentType` values are supported. The receiving application
must understand that UTI and the provider must return bytes in its declared
format.

## Sharing live WebView state

`ShareData` should not synchronously query the DOM or wait for JavaScript while
Cocoa is waiting for data. Keep a thread-safe Go snapshot of the shareable
model instead.

Daymark follows this pattern and renders a styled, paginated PDF when Cocoa
requests `com.adobe.pdf`:

1. The frontend emits `editor:share-content` whenever the current note changes.
2. Go decodes the note into `daymarkShareProvider`.
3. The provider protects its current note with `sync.RWMutex`.
4. When AppKit requests the advertised PDF, `ShareData` copies the latest note
   under the read lock and renders from that snapshot.

The frontend side is intentionally model-oriented:

```js
Events.Emit('editor:share-content', JSON.stringify({
  title: title.value,
  subtitle: subtitle.value,
  body: writing.innerText
}));
```

The Go side updates provider state without recreating the toolbar:

```go
app.Event.On("editor:share-content", func(event *application.CustomEvent) {
    raw, ok := event.Data.(string)
    if !ok {
        return
    }
    var note sharePayload
    if err := json.Unmarshal([]byte(raw), &note); err != nil {
        return
    }
    provider.Update(note)
})
```

This cleanly separates selection from representation. An application can share
the editor model as text, render semantic HTML, create a PDF, export an image,
or offer all of them. Wails never assumes that the visually rendered WebView is
the intended payload.

## Subject and suggested name

`SetSubject` configures services that support a message subject, such as Mail.
It does not modify the provider bytes.

`SetSuggestedName` names the shared logical item. It applies to all advertised
representations because one `NSItemProvider` represents one item; AppKit and
the receiving service use the requested UTI to determine the file type. The
name is available to `ShareData` through `MacShareRequest.SuggestedName`.

Both values can be changed after the toolbar is attached:

```go
share.SetSubject(currentDocument.Title)
share.SetSuggestedName(currentDocument.DisplayName)
```

## Provider lifecycle and replacement

Calling `SetProvider` performs the following operations:

- `ShareRepresentations` is called immediately;
- empty content types are discarded;
- duplicate content types are de-duplicated;
- the representation list is copied;
- the native item is enabled when at least one valid representation exists.

Passing `nil`, or a provider with no valid representations, clears and disables
the Share item:

```go
share.SetProvider(nil)
```

Replacing the provider while a share operation is open is safe. Each native
`NSItemProvider` retains an internal snapshot of the Go provider and advertised
formats with which it started. Internal registration IDs are generated and
released by Wails; they are not part of the public API and applications never
need to manage them.

## Concurrency and errors

Treat `ShareData` like a concurrent renderer:

- AppKit may request multiple formats for one service.
- Requests may overlap or be repeated.
- Protect mutable provider state with a mutex or immutable snapshots.
- Do not assume which service will request which format.
- Return an error for an unsupported or failed representation.
- Do not panic; Wails catches a provider panic and converts it to a sharing
  error, but ordinary errors provide better diagnostics.

`OnShared` runs after a service reports success. `OnShareError` runs when a
selected service or provider fails. Callbacks are asynchronous; protect any
shared application state they access. Live toolbar setters are safe to call
from these callbacks.

The general `SetEnabled` state and provider availability are combined. A Share
item is enabled only when the application enables it and a valid provider is
installed.

## Platform behavior

- macOS 10.15 and newer use `NSSharingServicePickerToolbarItem` directly.
- macOS 10.13–10.14 use an `NSToolbarItem` that opens the native
  `NSSharingServicePicker`.
- `SetSuggestedName` maps to `NSItemProvider.suggestedName` on macOS 10.14 and
  newer.
- On non-macOS platforms the native toolbar bridge is a no-op, allowing shared
  application code to compile without AppKit.

Other toolbar presentation features have their own system requirements: SF
Symbols need macOS 11, bordered items need macOS 10.15, and prominent style,
tint, and badges need macOS 26. The API safely applies only the capabilities
available on the current system.

## Complete implementation

See [`main.go`](./main.go) for the window, split-view, and toolbar setup,
[`split.go`](./split.go) for the native sidebar layout, and
[`assets/editor.html`](./assets/editor.html) for the editor-to-provider state
bridge.
