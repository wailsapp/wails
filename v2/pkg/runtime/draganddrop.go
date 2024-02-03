package runtime

import (
	"context"
	"fmt"
)

// OnFileDrop returns a slice of file path strings when a drop is finished.
func OnFileDrop(ctx context.Context, callback func(x, y int, paths []string)) {
	if callback == nil {
		LogError(ctx, "OnFileDrop called with a nil callback")
		return
	}
	EventsOn(ctx, "wails:file-drop", func(optionalData ...interface{}) {
		if len(optionalData) != 3 {
			callback(0, 0, nil)
		}
		x, ok := optionalData[0].(int)
		if !ok {
			LogError(ctx, fmt.Sprintf("invalid x coordinate in drag and drop: %v", optionalData[0]))
		}
		y, ok := optionalData[1].(int)
		if !ok {
			LogError(ctx, fmt.Sprintf("invalid y coordinate in drag and drop: %v", optionalData[1]))
		}
		paths, ok := optionalData[2].([]string)
		if !ok {
			LogError(ctx, fmt.Sprintf("invalid path data in drag and drop: %v", optionalData[2]))
		}
		callback(x, y, paths)
	})
}

// OnFileDropOff removes the drag and drop listeners and handlers.
func OnFileDropOff(ctx context.Context) {
	EventsOff(ctx, "wails:file-drop")
}
