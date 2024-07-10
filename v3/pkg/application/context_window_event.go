package application

var blankWindowEventContext = &WindowEventContext{}

const (
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

func (c WindowEventContext) setLocation(x int, y int) {
	c.data["x"] = x
	c.data["y"] = y
}

func (c WindowEventContext) Location() (int, int) {
	x, ok := c.data["x"].(int)
	if !ok {
		return 0, 0
	}
	y, ok := c.data["y"].(int)
	if !ok {
		return 0, 0
	}
	return x, y
}

func newWindowEventContext() *WindowEventContext {
	return &WindowEventContext{
		data: make(map[string]any),
	}
}
