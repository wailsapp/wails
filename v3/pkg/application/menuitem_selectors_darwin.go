// File: v3/pkg/application/menuitem_selectors_darwin.go

//go:build darwin

package application

// #cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
// #cgo LDFLAGS: -framework Cocoa
// #include <stdlib.h>
// #include "menuitem_darwin.h"
import "C"
import "unsafe"

var roleToSelector = map[Role]string{
    AppMenuRole:              "",  // This is a special case, handled separately
    AboutRole:                "orderFrontStandardAboutPanel:",
    ServicesMenuRole:         "",  // This is a submenu, no direct selector
    HideRole:                 "hide:",
    HideOthersRole:           "hideOtherApplications:",
    ShowAllRole:              "unhideAllApplications:",
    QuitRole:                 "terminate:",
    WindowMenuRole:           "",  // This is a submenu, no direct selector
    MinimizeRole:             "performMiniaturize:",
    ZoomRole:                 "performZoom:",
    BringAllToFrontRole:      "arrangeInFront:",
    CloseWindowRole:          "performClose:",
    EditMenuRole:             "",  // This is a submenu, no direct selector
    UndoRole:                 "undo:",
    RedoRole:                 "redo:",
    CutRole:                  "cut:",
    CopyRole:                 "copy:",
    PasteRole:                "paste:",
    DeleteRole:               "delete:",
    SelectAllRole:            "selectAll:",
    FindMenuRole:             "",  // This is a submenu, no direct selector
    FindRole:                 "performTextFinderAction:",
    FindAndReplaceRole:       "performTextFinderAction:",
    FindNextRole:             "performTextFinderAction:",
    FindPreviousRole:         "performTextFinderAction:",
    UseSelectionForFindRole:  "performTextFinderAction:",
    ViewMenuRole:             "",  // This is a submenu, no direct selector
    ToggleFullScreenRole:     "toggleFullScreen:",
    FileMenuRole:             "",  // This is a submenu, no direct selector
    NewRole:                  "newDocument:",
    OpenRole:                 "openDocument:",
    CloseRole:                "performClose:",
    SaveRole:                 "saveDocument:",
    SaveAsRole:               "saveDocumentAs:",
    RevertRole:               "revertDocumentToSaved:",
    PrintRole:                "printDocument:",
    HelpMenuRole:             "",  // This is a submenu, no direct selector
    HelpRole:                 "showHelp:",
    NoRole:                   "",  // No specific selector for this role
}

func getSelectorForRole(role Role) *C.char {
    if selector, ok := roleToSelector[role]; ok && selector != "" {
        return C.CString(selector)
    }
    return nil
}

//export processMenuItemClick
func processMenuItemClick(menuItemID C.uint) {
    if menuItem := getMenuItemByID(uint(menuItemID)); menuItem != nil {
        menuItem.handleClick()
    }
}