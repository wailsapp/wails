//go:build darwin && purego && !ios && !server

package application

// roleToSelector maps a menu Role to the Objective-C selector (by name) that
// AppKit dispatches up the responder chain for that role. Roles that are
// submenus (AppMenu, EditMenu, ...) or that have no built-in action map to the
// empty string and are handled elsewhere.
//
// This mirrors the cgo build's menuitem_selectors_darwin.go, but returns the
// selector as a plain Go string instead of a *C.char since the purego backend
// drives the runtime directly.
var roleToSelector = map[Role]string{
	//AppMenu:             "", // This is a special case, handled separately
	About: "orderFrontStandardAboutPanel:",
	//ServicesMenu:        "", // This is a submenu, no direct selector
	Hide:       "hide:",
	HideOthers: "hideOtherApplications:",
	ShowAll:    "unhideAllApplications:",
	Quit:       "terminate:",
	//WindowMenu:          "", // This is a submenu, no direct selector
	Minimise:        "performMiniaturize:",
	Zoom:            "performZoom:",
	BringAllToFront: "arrangeInFront:",
	CloseWindow:     "performClose:",
	//EditMenu:            "", // This is a submenu, no direct selector
	Undo:      "undo:",
	Redo:      "redo:",
	Cut:       "cut:",
	Copy:      "copy:",
	Paste:     "paste:",
	Delete:    "delete:",
	SelectAll: "selectAll:",
	//FindMenu:            "", // This is a submenu, no direct selector
	Find:           "performTextFinderAction:",
	FindAndReplace: "performTextFinderAction:",
	FindNext:       "performTextFinderAction:",
	FindPrevious:   "performTextFinderAction:",
	//ViewMenu:            "", // This is a submenu, no direct selector
	ToggleFullscreen: "toggleFullScreen:",
	//FileMenu:            "", // This is a submenu, no direct selector
	NewFile:       "newDocument:",
	Open:          "openDocument:",
	Save:          "saveDocument:",
	SaveAs:        "saveDocumentAs:",
	StartSpeaking: "startSpeaking:",
	StopSpeaking:  "stopSpeaking:",
	Revert:        "revertDocumentToSaved:",
	Print:         "printDocument:",
	PageLayout:    "runPageLayout:",
	//HelpMenu:            "", // This is a submenu, no direct selector
	Help: "showHelp:",
	//No:                  "", // No specific selector for this role
}

// getSelectorForRole returns the selector name for a role, or "" if the role
// has no built-in AppKit action (in which case the item uses a custom
// target/action pair driving processMenuItemClick).
func getSelectorForRole(role Role) string {
	if selector, ok := roleToSelector[role]; ok {
		return selector
	}
	return ""
}
