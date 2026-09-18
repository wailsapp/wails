//go:build !darwin || ios || server

package application

func macSidebarApplySnapshot(*MacSidebar)                   {}
func macSidebarApplySelection(*MacSidebar, *MacSidebarItem) {}
