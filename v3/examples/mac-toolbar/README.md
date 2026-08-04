# Native macOS toolbar

This example is a normal titled Wails window with a real AppKit `NSToolbar`.
Nothing in the toolbar is HTML.

It demonstrates:

- a persistent `NSToolbarItemGroup` for Write/Preview mode;
- an `NSSearchToolbarItem` with a Go callback (and the native search-field
  fallback used on macOS 10.13–10.15);
- an `NSSharingServicePickerToolbarItem` backed by a lazy, multi-format Go
  provider (with a native picker fallback on macOS 10.13–10.14);
- native New, Details, Save, and Focus controls;
- generated internal item identifiers, with no IDs in application code;
- live item handles that update labels, symbols, tooltips, presentation,
  badges, selection, enabled state, and visibility;
- a working notes editor with a note list, write/preview modes, search, native
  sharing, persistent local save state, details, focus mode, and keyboard save.

## Run the demo

Run it on macOS from this directory:

```sh
GOWORK=off go run .
```

The Share toolbar item opens the system share popover. Daymark offers PDF, HTML,
and plain text. Choosing Copy, Mail, Messages, Notes, or another installed
sharing service causes AppKit to request the representation or representations
that service understands.

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
4. When AppKit requests HTML or text, `ShareData` copies the latest note under
   the read lock and renders from that snapshot.

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

See [`main.go`](./main.go) for the complete Go provider and toolbar setup, and
[`assets/index.html`](./assets/index.html) for the editor-to-provider state
bridge.
