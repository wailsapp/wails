package application

type Context struct {
	// contains filtered or unexported fields
	data map[string]any
}

func newContext() *Context {
	return &Context{
		data: make(map[string]any),
	}
}

const (
	clickedMenuItem   string = "clickedMenuItem"
	menuItemIsChecked string = "menuItemIsChecked"
	contextMenuData   string = "contextMenuData"
)

func (c *Context) ClickedMenuItem() *MenuItem {
	result, exists := c.data[clickedMenuItem]
	if !exists {
		return nil
	}
	return result.(*MenuItem)
}

func (c *Context) IsChecked() bool {
	result, exists := c.data[menuItemIsChecked]
	if !exists {
		return false
	}
	return result.(bool)
}
func (c *Context) ContextMenuData() string {
	result := c.data[contextMenuData]
	if result == nil {
		return ""
	}
	str, ok := result.(string)
	if !ok {
		return ""
	}
	return str
}

func (c *Context) withClickedMenuItem(menuItem *MenuItem) *Context {
	c.data[clickedMenuItem] = menuItem
	return c
}

func (c *Context) withChecked(checked bool) {
	c.data[menuItemIsChecked] = checked
}

func (c *Context) withContextMenuData(data *ContextMenuData) *Context {
	if data == nil {
		return c
	}
	c.data[contextMenuData] = data.Data
	return c
}
