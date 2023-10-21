package application

var blankApplicationEventContext = &ApplicationEventContext{}

const (
	openedFiles = "openedFiles"
)

type ApplicationEventContext struct {
	// contains filtered or unexported fields
	data map[string]any
}

func (c ApplicationEventContext) OpenedFiles() []string {
	files, ok := c.data[openedFiles]
	if !ok {
		return nil
	}
	result, ok := files.([]string)
	if !ok {
		return nil
	}
	return result
}

func (c ApplicationEventContext) setOpenedFiles(files []string) {
	c.data[openedFiles] = files
}

func (c ApplicationEventContext) setIsDarkMode(mode bool) {
	c.data["isDarkMode"] = mode
}

func (c ApplicationEventContext) IsDarkMode() bool {
	mode, ok := c.data["isDarkMode"]
	if !ok {
		return false
	}
	result, ok := mode.(bool)
	if !ok {
		return false
	}
	return result
}

func newApplicationEventContext() *ApplicationEventContext {
	return &ApplicationEventContext{
		data: make(map[string]any),
	}
}
