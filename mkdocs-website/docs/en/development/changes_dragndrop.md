## Drag and Drop

Native drag and drop can be enabled per-window. Simply set the
`EnableDragAndDrop` window config option to `true` and the window will allow
files to be dragged onto it. When this happens, the `events.FilesDropped` event
will be emitted. The filenames can then be retrieved from the
`WindowEvent.Context()` using the `DroppedFiles()` method. This returns a slice
of strings containing the filenames.
