package runtime

import "context"

// DragAndDropOnMotion returns X and Y coordinates of the mouse position inside the window while it is dragging something.
func DragAndDropOnMotion(ctx context.Context, callback func(x, y int)) {
	EventsOn(ctx, "wails.dnd.motion", func(optionalData ...interface{}) {
		if callback != nil && len(optionalData) == 2 {
			var x, y int
			x = optionalData[0].(int)
			y = optionalData[1].(int)
			callback(x, y)
		}
	})
}

// DragAndDropOnDrop returns a slice of file path strings when a drop is finished.
func DragAndDropOnDrop(ctx context.Context, callback func(paths ...string)) {
	EventsOn(ctx, "wails.dnd.drop", func(optionalData ...interface{}) {
		if callback != nil && len(optionalData) > 0 {
			paths := optionalData[0].([]string)
			callback(paths...)
		}
	})
}
