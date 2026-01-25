# File Drop Example

This example demonstrates how to handle files being dragged from the operating system (Finder, Explorer, file managers) into a Wails application.

Dropped files are automatically categorised by type and displayed in separate buckets: documents, images, or other files.

## How it works

1. Enable file drops in window options:
   ```go
   EnableFileDrop: true
   ```

2. Mark elements as drop targets in HTML:
   ```html
   <div data-file-drop-target>Drop files here</div>
   ```

3. Listen for the `WindowFilesDropped` event:
   ```go
   win.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
       files := event.Context().DroppedFiles()
       details := event.Context().DropTargetDetails()
       // Handle the dropped files
   })
   ```

4. Optionally forward to frontend:
   ```go
   application.Get().Event.Emit("files-dropped", map[string]any{
       "files":   files,
       "details": details,
   })
   ```

## Drop Target Details

When files are dropped, you can get information about the drop location:

- `ElementID` - The ID of the element that received the drop
- `ClassList` - CSS classes on the drop target
- `X`, `Y` - Coordinates of the drop within the element

## Styling

When files are dragged over a valid drop target, Wails adds the `file-drop-target-active` class:

```css
.file-drop-target-active {
    border-color: #4a9eff;
    background: rgba(74, 158, 255, 0.1);
}
```

## Running the example

```bash
go run main.go
```

Then drag files from your desktop or file manager into the drop zone.

## HTML5 Drag and Drop API

This example also includes a demonstration for dragging elements *within* your application via the HTML5 Drag and Drop API.

Scroll down to the `Internal Drag and Drop` section within the launched application to interact with the demo.

## Status

| Platform | Status  |
|----------|---------|
| Mac      | Working |
| Windows  | Working |
| Linux    | Working |
