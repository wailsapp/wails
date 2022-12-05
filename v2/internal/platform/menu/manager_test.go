//go:build windows

package menu_test

import (
	"github.com/stretchr/testify/require"
	platformMenu "github.com/wailsapp/wails/v2/internal/platform/menu"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"testing"
)

func TestManager_ProcessClick_Checkbox(t *testing.T) {

	checkbox := menu.Label("Checkbox").SetChecked(false)
	menu1 := &menu.Menu{
		Items: []*menu.MenuItem{
			checkbox,
		},
	}
	menu2 := &menu.Menu{
		Items: []*menu.MenuItem{
			checkbox,
		},
	}
	menuWithNoCheckbox := &menu.Menu{
		Items: []*menu.MenuItem{
			menu.Label("No Checkbox"),
		},
	}
	clicked := false

	tests := []struct {
		name                string
		inputs              []*menu.Menu
		startState          bool
		expectedState       bool
		expectedMenuUpdates map[*menu.Menu][]*menu.MenuItem
		click               func(*menu.CallbackData)
	}{
		{
			name:   "should callback menu checkbox state when clicked (false -> true)",
			inputs: []*menu.Menu{menu1},
			expectedMenuUpdates: map[*menu.Menu][]*menu.MenuItem{
				menu1: {checkbox},
			},
			startState:    false,
			expectedState: true,
		},
		{
			name:          "should callback multiple menus when checkbox state when clicked (false -> true)",
			inputs:        []*menu.Menu{menu1, menu2},
			startState:    false,
			expectedState: true,
			expectedMenuUpdates: map[*menu.Menu][]*menu.MenuItem{
				menu1: {checkbox},
				menu2: {checkbox},
			},
		},
		{
			name:          "should callback only for the menus that the checkbox is in (false -> true)",
			inputs:        []*menu.Menu{menu1, menuWithNoCheckbox},
			startState:    false,
			expectedState: true,
			expectedMenuUpdates: map[*menu.Menu][]*menu.MenuItem{
				menu1: {checkbox},
			},
		},
		{
			name:   "should callback menu checkbox state when clicked (true->false)",
			inputs: []*menu.Menu{menu1},
			expectedMenuUpdates: map[*menu.Menu][]*menu.MenuItem{
				menu1: {checkbox},
			},
			startState:    true,
			expectedState: false,
		},
		{
			name:          "should callback multiple menus when checkbox state when clicked (true->false)",
			inputs:        []*menu.Menu{menu1, menu2},
			startState:    true,
			expectedState: false,
			expectedMenuUpdates: map[*menu.Menu][]*menu.MenuItem{
				menu1: {checkbox},
				menu2: {checkbox},
			},
		},
		{
			name:          "should callback only for the menus that the checkbox is in (true->false)",
			inputs:        []*menu.Menu{menu1, menuWithNoCheckbox},
			startState:    true,
			expectedState: false,
			expectedMenuUpdates: map[*menu.Menu][]*menu.MenuItem{
				menu1: {checkbox},
			},
		},
		{
			name:                "should callback no menus if checkbox not in them",
			inputs:              []*menu.Menu{menuWithNoCheckbox},
			startState:          false,
			expectedState:       false,
			expectedMenuUpdates: nil,
		},
		{
			name:          "should call Click on the checkbox",
			inputs:        []*menu.Menu{menu1, menu2},
			startState:    false,
			expectedState: true,
			expectedMenuUpdates: map[*menu.Menu][]*menu.MenuItem{
				menu1: {checkbox},
				menu2: {checkbox},
			},
			click: func(data *menu.CallbackData) {
				clicked = true
			},
		},
	}
	for _, tt := range tests {

		menusUpdated := map[*menu.Menu][]*menu.MenuItem{}
		clicked = false

		var checkMenuItemStateInMenu func(menu *menu.Menu)

		checkMenuItemStateInMenu = func(menu *menu.Menu) {
			for _, item := range menusUpdated[menu] {
				if item == checkbox {
					require.Equal(t, tt.expectedState, item.Checked)
				}
				if item.SubMenu != nil {
					checkMenuItemStateInMenu(item.SubMenu)
				}
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			m := platformMenu.NewManager()
			checkbox.SetChecked(tt.startState)
			checkbox.Click = tt.click
			for _, thisMenu := range tt.inputs {
				thisMenu := thisMenu
				m.AddMenu(thisMenu, func(menuItem *menu.MenuItem) {
					menusUpdated[thisMenu] = append(menusUpdated[thisMenu], menuItem)
				})
			}
			m.ProcessClick(checkbox)

			// Check the item has the correct state in all the menus
			for thisMenu := range menusUpdated {
				require.EqualValues(t, tt.expectedMenuUpdates[thisMenu], menusUpdated[thisMenu])
			}

			if tt.click != nil {
				require.Equal(t, true, clicked)
			}
		})
	}
}

func TestManager_ProcessClick_RadioGroups(t *testing.T) {

	radio1 := menu.Radio("Radio1", false, nil, nil)
	radio2 := menu.Radio("Radio2", false, nil, nil)
	radio3 := menu.Radio("Radio3", false, nil, nil)
	radio4 := menu.Radio("Radio4", false, nil, nil)
	radio5 := menu.Radio("Radio5", false, nil, nil)
	radio6 := menu.Radio("Radio6", false, nil, nil)

	radioGroupOne := &menu.Menu{
		Items: []*menu.MenuItem{
			radio1,
			radio2,
			radio3,
		},
	}

	radioGroupTwo := &menu.Menu{
		Items: []*menu.MenuItem{
			radio4,
			radio5,
			radio6,
		},
	}

	radioGroupThree := &menu.Menu{
		Items: []*menu.MenuItem{
			radio1,
			radio2,
			radio3,
		},
	}

	clicked := false

	tests := []struct {
		name                string
		inputs              []*menu.Menu
		startState          map[*menu.MenuItem]bool
		selected            *menu.MenuItem
		expectedMenuUpdates map[*menu.Menu][]*menu.MenuItem
		click               func(*menu.CallbackData)
		expectedState       map[*menu.MenuItem]bool
	}{
		{
			name:   "should only set the clicked radio item",
			inputs: []*menu.Menu{radioGroupOne},
			expectedMenuUpdates: map[*menu.Menu][]*menu.MenuItem{
				radioGroupOne: {radio1, radio2, radio3},
			},
			startState: map[*menu.MenuItem]bool{
				radio1: true,
				radio2: false,
				radio3: false,
			},
			selected: radio2,
			expectedState: map[*menu.MenuItem]bool{
				radio1: false,
				radio2: true,
				radio3: false,
			},
		},
		{
			name:   "should not affect other radio groups or menus",
			inputs: []*menu.Menu{radioGroupOne, radioGroupTwo},
			expectedMenuUpdates: map[*menu.Menu][]*menu.MenuItem{
				radioGroupOne: {radio1, radio2, radio3},
			},
			startState: map[*menu.MenuItem]bool{
				radio1: true,
				radio2: false,
				radio3: false,
				radio4: true,
				radio5: false,
				radio6: false,
			},
			selected: radio2,
			expectedState: map[*menu.MenuItem]bool{
				radio1: false,
				radio2: true,
				radio3: false,
				radio4: true,
				radio5: false,
				radio6: false,
			},
		},
		{
			name:   "menus with the same radio group should be updated",
			inputs: []*menu.Menu{radioGroupOne, radioGroupThree},
			expectedMenuUpdates: map[*menu.Menu][]*menu.MenuItem{
				radioGroupOne:   {radio1, radio2, radio3},
				radioGroupThree: {radio1, radio2, radio3},
			},
			startState: map[*menu.MenuItem]bool{
				radio1: true,
				radio2: false,
				radio3: false,
			},
			selected: radio2,
			expectedState: map[*menu.MenuItem]bool{
				radio1: false,
				radio2: true,
				radio3: false,
			},
		},
	}
	for _, tt := range tests {

		menusUpdated := map[*menu.Menu][]*menu.MenuItem{}
		clicked = false

		t.Run(tt.name, func(t *testing.T) {
			m := platformMenu.NewManager()

			for item, value := range tt.startState {
				item.SetChecked(value)
			}

			tt.selected.Click = tt.click
			for _, thisMenu := range tt.inputs {
				thisMenu := thisMenu
				m.AddMenu(thisMenu, func(menuItem *menu.MenuItem) {
					menusUpdated[thisMenu] = append(menusUpdated[thisMenu], menuItem)
				})
			}
			m.ProcessClick(tt.selected)
			require.Equal(t, tt.expectedMenuUpdates, menusUpdated)

			// Check the items have the correct state in all the menus
			for item, expectedValue := range tt.expectedState {
				require.Equal(t, expectedValue, item.Checked)
			}

			if tt.click != nil {
				require.Equal(t, true, clicked)
			}
		})
	}
}
