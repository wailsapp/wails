package application

var blankWindowEventContext = &WindowEventContext{}

const (
	droppedFiles         = "droppedFiles"
	dropTargetDetailsKey = "dropTargetDetails"
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

func (c WindowEventContext) setDropTargetDetails(details *DropTargetDetails) {
	if c.data == nil {
		c.data = make(map[string]any)
	}
	if details == nil {
		c.data[dropTargetDetailsKey] = nil
		return
	}
	c.data[dropTargetDetailsKey] = details
}

// DropTargetDetails retrieves information about the drop target element.
func (c WindowEventContext) DropTargetDetails() *DropTargetDetails {
	if c.data == nil {
		c.data = make(map[string]any)
	}
	details, ok := c.data[dropTargetDetailsKey]
	if !ok {
		return nil
	}
	if details == nil {
		return nil
	}
	result, ok := details.(*DropTargetDetails)
	if !ok {
		return nil
	}
	return result
}

func newWindowEventContext() *WindowEventContext {
	return &WindowEventContext{
		data: make(map[string]any),
	}
}
