package application

var blankWindowEventContext = &WindowEventContext{}

const (
	droppedFiles       = "droppedFiles"
	dropZoneDetailsKey = "dropZoneDetails"
)

type WindowEventContext struct {
	// contains filtered or unexported fields
	data map[string]any
}

func (c WindowEventContext) DroppedFiles() []string {
	if c.data == nil {
		c.data = make(map[string]any)
	}
	files, ok := c.data[droppedFiles]
	if !ok {
		return nil
	}
	result, ok := files.([]string)
	if !ok {
		return nil
	}
	return result
}

func (c WindowEventContext) setDroppedFiles(files []string) {
	if c.data == nil {
		c.data = make(map[string]any)
	}
	c.data[droppedFiles] = files
}

func (c WindowEventContext) setCoordinates(x, y int) {
	if c.data == nil {
		c.data = make(map[string]any)
	}
	c.data["x"] = x
	c.data["y"] = y
}

func (c WindowEventContext) setDropZoneDetails(details *DropZoneDetails) {
	if c.data == nil {
		c.data = make(map[string]any)
	}
	if details == nil {
		c.data[dropZoneDetailsKey] = nil
		return
	}
	c.data[dropZoneDetailsKey] = details
}

// DropZoneDetails retrieves the detailed drop zone information, if available.
func (c WindowEventContext) DropZoneDetails() *DropZoneDetails {
	if c.data == nil {
		c.data = make(map[string]any)
	}
	details, ok := c.data[dropZoneDetailsKey]
	if !ok {
		return nil
	}
	// Explicitly type assert, handle if it's nil (though setDropZoneDetails should handle it)
	if details == nil {
		return nil
	}
	result, ok := details.(*DropZoneDetails)
	if !ok {
		// This case indicates a programming error if data was set incorrectly
		return nil
	}
	return result
}

func newWindowEventContext() *WindowEventContext {
	return &WindowEventContext{
		data: make(map[string]any),
	}
}
