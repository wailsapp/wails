package runtime

import (
	"context"
	"fmt"
)

// DragAndDropOn returns a slice of file path strings when a drop is finished.
func DragAndDropOn(ctx context.Context, callback func(x, y int, paths []string)) {
	EventsOn(ctx, "wails:dnd:drop", func(optionalData ...interface{}) {

		if len(optionalData) != 3 {
			callback(0, 0, nil)
			return
		}

		if callback != nil && len(optionalData) > 0 {
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
		}
	})
}

// DragAndDropOff removes the drag and drop listeners and handlers.
func DragAndDropOff(ctx context.Context) {
	EventsOff(ctx, "wails:dnd:drop")
}
