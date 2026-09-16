//go:build !darwin || ios || server

package application

import "unsafe"

func macWindowTabsAdd(unsafe.Pointer, unsafe.Pointer, MacTabOrder) error {
	return ErrMacWindowTabsUnsupported
}

func macWindowTabsHasGroup(unsafe.Pointer) bool                 { return false }
func macWindowTabsGroupPointer(unsafe.Pointer) unsafe.Pointer   { return nil }
func macWindowTabsGroupWindows(unsafe.Pointer) []unsafe.Pointer { return nil }
func macWindowTabsSelectedWindow(unsafe.Pointer) unsafe.Pointer { return nil }
func macWindowTabsSelectNext(unsafe.Pointer)                    {}
func macWindowTabsSelectPrevious(unsafe.Pointer)                {}
func macWindowTabsSelect(unsafe.Pointer, unsafe.Pointer)        {}
func macWindowTabsIsTabBarVisible(unsafe.Pointer) bool          { return false }
func macWindowTabsToggleTabBar(unsafe.Pointer)                  {}
func macWindowTabsIsOverviewVisible(unsafe.Pointer) bool        { return false }
func macWindowTabsToggleOverview(unsafe.Pointer)                {}
func macWindowTabsIdentifier(unsafe.Pointer) string             { return "" }
func macWindowTabsMoveToNewWindow(unsafe.Pointer)               {}
func macWindowTabsMergeAllWindows(unsafe.Pointer)               {}
func macWindowTabsSetTitle(unsafe.Pointer, string)              {}
func macWindowTabsSetToolTip(unsafe.Pointer, string)            {}
