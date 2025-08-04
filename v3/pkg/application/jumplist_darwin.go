//go:build darwin

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
	app        *darwinApp
	categories []JumpListCategory
}

func (app *darwinApp) CreateJumpList() *JumpList {
	return &JumpList{
		app:        app,
		categories: []JumpListCategory{},
	}
}

func (j *JumpList) AddCategory(category JumpListCategory) {
	// Stub implementation for macOS
	j.categories = append(j.categories, category)
}

func (j *JumpList) ClearCategories() {
	// Stub implementation for macOS
	j.categories = []JumpListCategory{}
}

func (j *JumpList) Apply() error {
	// Stub implementation for macOS
	// Jump lists are Windows-specific, so this is a no-op on macOS
	return nil
}