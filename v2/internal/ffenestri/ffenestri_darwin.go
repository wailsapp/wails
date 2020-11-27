package ffenestri

/*
#cgo darwin CFLAGS: -DFFENESTRI_DARWIN=1
#cgo darwin LDFLAGS: -framework WebKit -lobjc

extern void TitlebarAppearsTransparent(void *);
extern void HideTitle(void *);
extern void HideTitleBar(void *);
extern void FullSizeContent(void *);
extern void UseToolbar(void *);
extern void HideToolbarSeparator(void *);
extern void DisableFrame(void *);
extern void SetAppearance(void *, const char *);
extern void WebviewIsTransparent(void *);
extern void SetWindowBackgroundIsTranslucent(void *);
extern void SetMenu(void *, const char *);
*/
import "C"
import (
	"encoding/json"

	"github.com/wailsapp/wails/v2/pkg/menu"
)

func (a *Application) processPlatformSettings() error {

	mac := a.config.Mac
	titlebar := mac.TitleBar

	// HideTitle
	if titlebar.HideTitle {
		C.HideTitle(a.app)
	}

	// HideTitleBar
	if titlebar.HideTitleBar {
		C.HideTitleBar(a.app)
	}

	// Full Size Content
	if titlebar.FullSizeContent {
		C.FullSizeContent(a.app)
	}

	// Toolbar
	if titlebar.UseToolbar {
		C.UseToolbar(a.app)
	}

	if titlebar.HideToolbarSeparator {
		C.HideToolbarSeparator(a.app)
	}

	if titlebar.TitlebarAppearsTransparent {
		C.TitlebarAppearsTransparent(a.app)
	}

	// Process window Appearance
	if mac.Appearance != "" {
		C.SetAppearance(a.app, a.string2CString(string(mac.Appearance)))
	}

	// Check if the webview should be transparent
	if mac.WebviewIsTransparent {
		C.WebviewIsTransparent(a.app)
	}

	// Check if window should be translucent
	if mac.WindowBackgroundIsTranslucent {
		C.SetWindowBackgroundIsTranslucent(a.app)
	}

	// Process menu
	if mac.Menu != nil {

		/*
			As radio groups need to be manually managed on OSX,
			we preprocess the menu to determine the radio groups.
			This is defined as any adjacent menu item of type "RadioType".
			We keep a record of every radio group member we discover by saving
			a list of all members of the group and the number of members
			in the group (this last one is for optimisation at the C layer).

			Example:
			{
				"RadioGroups": [
					{
						"Members": [
							"option-1",
							"option-2",
							"option-3"
						],
						"Length": 3
					}
				]
			}
		*/
		processedMenu := NewProcessedMenu(mac.Menu)
		menuJSON, err := json.Marshal(processedMenu)
		if err != nil {
			return err
		}
		C.SetMenu(a.app, a.string2CString(string(menuJSON)))
	}

	return nil
}

// ProcessedMenu is the original menu with the addition
// of radio groups extracted from the menu data
type ProcessedMenu struct {
	Menu              *menu.Menu
	RadioGroups       []*RadioGroup
	currentRadioGroup []string
}

// RadioGroup holds all the members of the same radio group
type RadioGroup struct {
	Members []string
	Length  int
}

// NewProcessedMenu processed the given menu and returns
// the original menu with the extracted radio groups
func NewProcessedMenu(menu *menu.Menu) *ProcessedMenu {
	result := &ProcessedMenu{
		Menu:              menu,
		RadioGroups:       []*RadioGroup{},
		currentRadioGroup: []string{},
	}

	result.processMenu()

	return result
}

func (p *ProcessedMenu) processMenu() {
	// Loop over top level menus
	for _, item := range p.Menu.Items {
		// Process MenuItem
		p.processMenuItem(item)
	}

	p.finaliseRadioGroup()
}

func (p *ProcessedMenu) processMenuItem(item *menu.MenuItem) {

	switch item.Type {

	// We need to recurse submenus
	case menu.SubmenuType:

		// Finalise any current radio groups as they don't trickle down to submenus
		p.finaliseRadioGroup()

		// Process each submenu item
		for _, subitem := range item.SubMenu {
			p.processMenuItem(subitem)
		}
	case menu.RadioType:
		// Add the item to the radio group
		p.currentRadioGroup = append(p.currentRadioGroup, item.ID)
	default:
		p.finaliseRadioGroup()
	}
}

func (p *ProcessedMenu) finaliseRadioGroup() {

	// If we were processing a radio group, fix up the references
	if len(p.currentRadioGroup) > 0 {

		// Create new radiogroup
		group := &RadioGroup{
			Members: p.currentRadioGroup,
			Length:  len(p.currentRadioGroup),
		}
		p.RadioGroups = append(p.RadioGroups, group)

		// Empty the radio group
		p.currentRadioGroup = []string{}
	}
}
