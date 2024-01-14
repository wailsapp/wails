package runtime

import "context"

// HandleDragAndDrop returns a slice of file path strings when a drop is finished.
func HandleDragAndDrop(ctx context.Context, callback func(paths ...string)) {
	EventsOn(ctx, "wails.dnd.drop", func(optionalData ...interface{}) {
		if callback != nil && len(optionalData) > 0 {
			paths := optionalData[0].([]string)
			callback(paths...)
		}
	})
}
