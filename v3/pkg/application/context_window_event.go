package application

var blankWindowEventContext = &WindowEventContext{}

const (
	// FilesDropped is the event name for when files are dropped on the window
	droppedFiles = "droppedFiles"
)

type WindowEventContext struct {
	// contains filtered or unexported fields
	data map[string]any
}

func (c WindowEventContext) DroppedFiles() []string {
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
	c.data[droppedFiles] = files
}

func newWindowEventContext() *WindowEventContext {
	return &WindowEventContext{
		data: make(map[string]any),
	}
}
