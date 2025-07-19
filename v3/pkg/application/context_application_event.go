package application

import "log"

var blankApplicationEventContext = &ApplicationEventContext{}

const (
	CONTEXT_OPENED_FILES = "openedFiles"
	CONTEXT_FILENAME     = "filename"
	CONTEXT_URL          = "url"
)

// ApplicationEventContext is the context of an application event
type ApplicationEventContext struct {
	// contains filtered or unexported fields
	data map[string]any
}

// OpenedFiles returns the opened files from the event context if it was set
func (c ApplicationEventContext) OpenedFiles() []string {
	files, ok := c.data[CONTEXT_OPENED_FILES]
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
	c.data[CONTEXT_OPENED_FILES] = files
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

// IsDarkMode returns true if the event context has a dark mode
func (c ApplicationEventContext) IsDarkMode() bool {
	return c.getBool("isDarkMode")
}

// HasVisibleWindows returns true if the event context has a visible window
func (c ApplicationEventContext) HasVisibleWindows() bool {
	return c.getBool("hasVisibleWindows")
}

func (c *ApplicationEventContext) setData(data map[string]any) {
	c.data = data
}

func (c *ApplicationEventContext) setOpenedWithFile(filepath string) {
	c.data[CONTEXT_FILENAME] = filepath
}

func (c *ApplicationEventContext) setURL(openedWithURL string) {
	c.data[CONTEXT_URL] = openedWithURL
}

// Filename returns the filename from the event context if it was set
func (c ApplicationEventContext) Filename() string {
	filename, ok := c.data[CONTEXT_FILENAME]
	if !ok {
		return ""
	}
	result, ok := filename.(string)
	if !ok {
		return ""
	}
	return result
}

// URL returns the URL from the event context if it was set
func (c ApplicationEventContext) URL() string {
	url, ok := c.data[CONTEXT_URL]
	if !ok {
		log.Println("URL not found in event context")
		return ""
	}
	result, ok := url.(string)
	if !ok {
		log.Println("URL not a string in event context")
		return ""
	}
	return result
}

func newApplicationEventContext() *ApplicationEventContext {
	return &ApplicationEventContext{
		data: make(map[string]any),
	}
}
