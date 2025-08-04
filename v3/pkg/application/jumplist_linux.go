//go:build linux

package application

type JumpListItemType int

const (
	JumpListItemTypeTask JumpListItemType = iota
	JumpListItemTypeSeparator
)

type JumpListItem struct {
	Type        JumpListItemType
	Title       string
	Description string
	FilePath    string
	Arguments   string
	IconPath    string
	IconIndex   int
}

type JumpListCategory struct {
	Name  string
	Items []JumpListItem
}

type JumpList struct {
	app        *linuxApp
	categories []JumpListCategory
}

func (app *linuxApp) CreateJumpList() *JumpList {
	return &JumpList{
		app:        app,
		categories: []JumpListCategory{},
	}
}

func (j *JumpList) AddCategory(category JumpListCategory) {
	// Stub implementation for Linux
	j.categories = append(j.categories, category)
}

func (j *JumpList) ClearCategories() {
	// Stub implementation for Linux
	j.categories = []JumpListCategory{}
}

func (j *JumpList) Apply() error {
	// Stub implementation for Linux
	// Jump lists are Windows-specific, so this is a no-op on Linux
	return nil
}