// File: v3/pkg/application/menuitem_selectors_darwin.go

//go:build darwin

package application

import "C"

var roleToSelector = map[Role]string{
	//AppMenu:             "", // This is a special case, handled separately
	About: "orderFrontStandardAboutPanel:",
	//ServicesMenu:        "", // This is a submenu, no direct selector
	Hide:       "hide:",
	HideOthers: "hideOtherApplications:",
	ShowAll:    "unhideAllApplications:",
	Quit:       "terminate:",
	//WindowMenu:          "", // This is a submenu, no direct selector
	Minimize:        "performMiniaturize:",
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

func getSelectorForRole(role Role) *C.char {
	if selector, ok := roleToSelector[role]; ok && selector != "" {
		return C.CString(selector)
	}
	return nil
}
