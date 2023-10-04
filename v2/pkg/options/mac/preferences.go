package mac

import "github.com/leaanthony/u"

var Enabled = u.True
var Disabled = u.False

// Preferences allows to set webkit preferences
type Preferences struct {
	// A Boolean value that indicates whether pressing the tab key changes the focus to links and form controls.
	// Set to false by default.
	TabFocusesLinks u.Bool
	// A Boolean value that indicates whether to allow people to select or otherwise interact with text.
	// Set to true by default.
	TextInteractionEnabled u.Bool
}
