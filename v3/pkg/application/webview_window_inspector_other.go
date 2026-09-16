//go:build !darwin || ios || server

package application

func macInspectorRegisterControlIfInstalled(*MacInspectorControl) {}
func macInspectorRegisterSectionIfInstalled(*MacInspectorSection) {}
func macInspectorApplySnapshot(*MacInspector)                     {}
func macInspectorApplyControl(*MacInspectorControl)               {}
