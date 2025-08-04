package application

// JumpListItemType represents the type of jump list item
type JumpListItemType int

const (
	// JumpListItemTypeTask represents a task item
	JumpListItemTypeTask JumpListItemType = iota
	// JumpListItemTypeSeparator represents a separator (Windows only)
	JumpListItemTypeSeparator
)

// JumpListItem represents a single item in a jump list
type JumpListItem struct {
	Type        JumpListItemType
	Title       string
	Description string
	FilePath    string
	Arguments   string
	IconPath    string
	IconIndex   int
}

// JumpListCategory represents a category of items in a jump list
type JumpListCategory struct {
	Name  string
	Items []JumpListItem
}

// JumpList provides an interface for managing application jump lists.
// This is primarily a Windows feature, but the API is designed to be
// cross-platform safe (no-op on non-Windows platforms).
type JumpList struct {
	app        platformApp
	categories []JumpListCategory
}

// AddCategory adds a category to the jump list.
// If the category name is empty, the items will be added as tasks.
func (j *JumpList) AddCategory(category JumpListCategory) {
	j.categories = append(j.categories, category)
}

// ClearCategories removes all categories from the jump list.
func (j *JumpList) ClearCategories() {
	j.categories = []JumpListCategory{}
}

// Apply applies the current jump list configuration.
// On Windows, this updates the taskbar jump list.
// On other platforms, this is a no-op.
func (j *JumpList) Apply() error {
	// Platform-specific implementation
	return j.applyPlatform()
}