package runtime

import (
	"context"
	"fmt"
)

func processDropData(ctx context.Context, data ...interface{}) (int, int, []string) {
	if len(data) != 3 {
		return 0, 0, nil
	}
	x, ok := data[0].(int)
	if !ok {
		LogError(ctx, fmt.Sprintf("invalid x coordinate in drag and drop: %v", data[0]))
	}
	y, ok := data[1].(int)
	if !ok {
		LogError(ctx, fmt.Sprintf("invalid y coordinate in drag and drop: %v", data[1]))
	}
	paths, ok := data[2].([]string)
	if !ok {
		LogError(ctx, fmt.Sprintf("invalid path data in drag and drop: %v", data[2]))
	}
	return x, y, paths
}

// OnFileDropHover returns a slice of file path strings when a drop is finished.
func OnFileDropHover(ctx context.Context, callback func(x, y int, paths []string)) {
	EventsOn(ctx, "wails:file-drop-hover", func(optionalData ...interface{}) {
		if callback != nil && len(optionalData) == 3 {
			callback(processDropData(ctx, optionalData))
			return
		}
		callback(0, 0, nil)
	})
}

// OnFileDrop returns a slice of file path strings when a drop is finished.
func OnFileDrop(ctx context.Context, callback func(x, y int, paths []string)) {
	EventsOn(ctx, "wails:file-drop", func(optionalData ...interface{}) {
		if callback != nil && len(optionalData) == 3 {
			callback(processDropData(ctx, optionalData))
			return
		}
		callback(0, 0, nil)
	})
}

// OnFileDropCancelled returns a slice of file path strings when a drop is finished.
func OnFileDropCancelled(ctx context.Context, callback func()) {
	EventsOn(ctx, "wails:file-drop-cancelled", func(optionalData ...interface{}) {
		if callback != nil {
			callback()
		}
	})
}

// DragAndDropOff removes the drag and drop listeners and handlers.
func DragAndDropOff(ctx context.Context) {
	EventsOff(ctx, "wails:file-drop")
}
