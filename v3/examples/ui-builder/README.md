# UI Builder

A drag & drop page builder, generated from the built-in `ui-builder` template.
Drag components from the palette onto the artboard, tune them in the inspector,
rearrange them in the layers panel, then save the layout or export it as a
standalone HTML page through native dialogs provided by the Go backend.

![UI Builder](screenshot.png)

## Running

```
cd v3/examples/ui-builder
wails3 dev        # live reload
wails3 build      # production binary in bin/
```

The frontend build output (`frontend/dist`) and the generated bindings
(`frontend/bindings`) are committed, so `go run .` also works straight away.

## Try it

- Drag **Section**, **Row** or **Card** in first: they are containers, and an
  empty one shows a landing pad. Drop leaf components inside.
- Click any component to edit it in the inspector. Double-click grabs its
  container. `⇧⇥` walks up to the parent, `↑`/`↓` move between siblings.
- Switch the artboard between desktop, tablet and phone widths, or toggle the
  light/dark theme, from the top bar. **Preview** hides the builder chrome.
- **Open** one of the sample layouts in [`layouts/`](layouts/):
  `dashboard.uib.json` (a dark analytics page) or `signup-form.uib.json`.
- **Export HTML** writes a self-contained page that uses the same stylesheet
  as the canvas, so it looks exactly like the artboard.
- Press `?` for the full list of keyboard shortcuts.

## How it works

| Path | Purpose |
| --- | --- |
| `main.go` | Creates the Wails app and the builder window |
| `layoutservice.go` | `LayoutService` — save / open layouts and export HTML via native dialogs |
| `frontend/src/components.ts` | The component registry: defaults, inspector fields and renderers |
| `frontend/src/theme.ts` | Stylesheet for the components you build (shared with the export) |
| `frontend/src/store.ts` | Document state, selection and undo/redo history |
| `frontend/src/canvas.ts` | Artboard rendering, selection and drag & drop |
| `frontend/src/palette.ts` · `inspector.ts` · `layers.ts` | The three panels |
| `frontend/src/exporter.ts` | Turns a layout into a standalone HTML page |
| `frontend/public/style.css` | The builder's own chrome (panels, toolbar, indicators) |

Layouts are plain JSON. The builder also autosaves the current design to the
browser's local storage, so closing the app never loses work.

To scaffold your own copy of this app run `wails3 init -t ui-builder`.
