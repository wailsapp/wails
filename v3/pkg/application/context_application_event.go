package application

var blankApplicationEventContext = &ApplicationEventContext{}

const (
	openedFiles = "openedFiles"
	filename    = "filename"
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

func (c ApplicationEventContext) getBool(key string) bool {
	mode, ok := c.data[key]
	if !ok {
		return false
	}
	result, ok := mode.(bool)
	if !ok {
		return false
	}
	return result
}

func (c ApplicationEventContext) IsDarkMode() bool {
	return c.getBool("isDarkMode")
}

func (c ApplicationEventContext) HasVisibleWindows() bool {
	return c.getBool("hasVisibleWindows")
}

func (c ApplicationEventContext) setData(data map[string]any) {
	c.data = data
}

func (c ApplicationEventContext) setOpenedWithFile(filepath string) {
	c.data[filename] = filepath
}

func (c ApplicationEventContext) Filename() string {
	filename, ok := c.data[filename]
	if !ok {
		return ""
	}
	result, ok := filename.(string)
	if !ok {
		return ""
	}
	return result
}

func newApplicationEventContext() *ApplicationEventContext {
	return &ApplicationEventContext{
		data: make(map[string]any),
	}
}
