//go:build !darwin || ios || server

package application

// The native content list exists only on desktop macOS. These stubs keep the
// public API buildable everywhere; the model still works and the handles are
// inert, matching the sidebar and inspector.

func macContentListApplySnapshot(*MacContentList)            {}
func macContentListApplyRow(*MacContentListRow)              {}
func macContentListApplySelection(*MacContentList, []uint64) {}
