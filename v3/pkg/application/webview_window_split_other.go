//go:build !darwin || ios || server

package application

// The native split view exists only on desktop macOS. These stubs keep the
// public API buildable everywhere; pane handles are inert and the window
// keeps its ordinary primary WebView.

func macSplitPaneApplyMinimumThickness(*MacSplitPane, float64)       {}
func macSplitPaneApplyMaximumThickness(*MacSplitPane, float64)       {}
func macSplitPaneApplyPreferredFraction(*MacSplitPane, float64)      {}
func macSplitPaneApplyHoldingPriority(*MacSplitPane, float64)        {}
func macSplitPaneApplyCollapsible(*MacSplitPane, bool)               {}
func macSplitPaneApplyCanCollapseFromResize(*MacSplitPane, bool)     {}
func macSplitPaneApplyContentLayout(*MacSplitPane, MacContentLayout) {}
func macSplitPaneApplyCollapsed(*MacSplitPane, bool)                 {}
func macSplitPaneApplyToggle(*MacSplitPane)                          {}
