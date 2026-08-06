//go:build !darwin || ios || server

package application

func macInspectorRegisterControlIfInstalled(*MacInspectorControl) {}
func macInspectorApplySnapshot(*MacInspector)                     {}
func macInspectorApplyControl(*MacInspectorControl)               {}
