package events

type ApplicationEventType uint
type WindowEventType uint

var Mac = newMacEvents()

type macEvents struct {
	ApplicationDidBecomeActive                              ApplicationEventType
	ApplicationDidChangeBackingProperties                   ApplicationEventType
	ApplicationDidChangeEffectiveAppearance                 ApplicationEventType
	ApplicationDidChangeIcon                                ApplicationEventType
	ApplicationDidChangeOcclusionState                      ApplicationEventType
	ApplicationDidChangeScreenParameters                    ApplicationEventType
	ApplicationDidChangeStatusBarFrame                      ApplicationEventType
	ApplicationDidChangeStatusBarOrientation                ApplicationEventType
	ApplicationDidFinishLaunching                           ApplicationEventType
	ApplicationDidHide                                      ApplicationEventType
	ApplicationDidResignActive                              ApplicationEventType
	ApplicationDidUnhide                                    ApplicationEventType
	ApplicationDidUpdate                                    ApplicationEventType
	ApplicationWillBecomeActive                             ApplicationEventType
	ApplicationWillFinishLaunching                          ApplicationEventType
	ApplicationWillHide                                     ApplicationEventType
	ApplicationWillResignActive                             ApplicationEventType
	ApplicationWillTerminate                                ApplicationEventType
	ApplicationWillUnhide                                   ApplicationEventType
	ApplicationWillUpdate                                   ApplicationEventType
	WindowDidBecomeKey                                      WindowEventType
	WindowDidBecomeMain                                     WindowEventType
	WindowDidBeginSheet                                     WindowEventType
	WindowDidChangeAlpha                                    WindowEventType
	WindowDidChangeBackingLocation                          WindowEventType
	WindowDidChangeBackingProperties                        WindowEventType
	WindowDidChangeCollectionBehavior                       WindowEventType
	WindowDidChangeEffectiveAppearance                      WindowEventType
	WindowDidChangeOcclusionState                           WindowEventType
	WindowDidChangeOrderingMode                             WindowEventType
	WindowDidChangeScreen                                   WindowEventType
	WindowDidChangeScreenParameters                         WindowEventType
	WindowDidChangeScreenProfile                            WindowEventType
	WindowDidChangeScreenSpace                              WindowEventType
	WindowDidChangeScreenSpaceProperties                    WindowEventType
	WindowDidChangeSharingType                              WindowEventType
	WindowDidChangeSpace                                    WindowEventType
	WindowDidChangeSpaceOrderingMode                        WindowEventType
	WindowDidChangeTitle                                    WindowEventType
	WindowDidChangeToolbar                                  WindowEventType
	WindowDidChangeVisibility                               WindowEventType
	WindowDidClose                                          WindowEventType
	WindowDidDeminiaturize                                  WindowEventType
	WindowDidEndSheet                                       WindowEventType
	WindowDidEnterFullScreen                                WindowEventType
	WindowDidEnterVersionBrowser                            WindowEventType
	WindowDidExitFullScreen                                 WindowEventType
	WindowDidExitVersionBrowser                             WindowEventType
	WindowDidExpose                                         WindowEventType
	WindowDidFocus                                          WindowEventType
	WindowDidMiniaturize                                    WindowEventType
	WindowDidMove                                           WindowEventType
	WindowDidOrderOffScreen                                 WindowEventType
	WindowDidOrderOnScreen                                  WindowEventType
	WindowDidResignKey                                      WindowEventType
	WindowDidResignMain                                     WindowEventType
	WindowDidResize                                         WindowEventType
	WindowDidUnfocus                                        WindowEventType
	WindowDidUpdate                                         WindowEventType
	WindowDidUpdateAlpha                                    WindowEventType
	WindowDidUpdateCollectionBehavior                       WindowEventType
	WindowDidUpdateCollectionProperties                     WindowEventType
	WindowDidUpdateShadow                                   WindowEventType
	WindowDidUpdateTitle                                    WindowEventType
	WindowDidUpdateToolbar                                  WindowEventType
	WindowDidUpdateVisibility                               WindowEventType
	WindowWillBecomeKey                                     WindowEventType
	WindowWillBecomeMain                                    WindowEventType
	WindowWillBeginSheet                                    WindowEventType
	WindowWillChangeOrderingMode                            WindowEventType
	WindowWillClose                                         WindowEventType
	WindowWillDeminiaturize                                 WindowEventType
	WindowWillEnterFullScreen                               WindowEventType
	WindowWillEnterVersionBrowser                           WindowEventType
	WindowWillExitFullScreen                                WindowEventType
	WindowWillExitVersionBrowser                            WindowEventType
	WindowWillFocus                                         WindowEventType
	WindowWillMiniaturize                                   WindowEventType
	WindowWillMove                                          WindowEventType
	WindowWillOrderOffScreen                                WindowEventType
	WindowWillOrderOnScreen                                 WindowEventType
	WindowWillResignMain                                    WindowEventType
	WindowWillResize                                        WindowEventType
	WindowWillUnfocus                                       WindowEventType
	WindowWillUpdate                                        WindowEventType
	WindowWillUpdateAlpha                                   WindowEventType
	WindowWillUpdateCollectionBehavior                      WindowEventType
	WindowWillUpdateCollectionProperties                    WindowEventType
	WindowWillUpdateShadow                                  WindowEventType
	WindowWillUpdateTitle                                   WindowEventType
	WindowWillUpdateToolbar                                 WindowEventType
	WindowWillUpdateVisibility                              WindowEventType
	WindowWillUseStandardFrame                              WindowEventType
	MenuWillOpen                                            ApplicationEventType
	MenuDidOpen                                             ApplicationEventType
	MenuDidClose                                            ApplicationEventType
	MenuWillSendAction                                      ApplicationEventType
	MenuDidSendAction                                       ApplicationEventType
	MenuWillHighlightItem                                   ApplicationEventType
	MenuDidHighlightItem                                    ApplicationEventType
	MenuWillDisplayItem                                     ApplicationEventType
	MenuDidDisplayItem                                      ApplicationEventType
	MenuWillAddItem                                         ApplicationEventType
	MenuDidAddItem                                          ApplicationEventType
	MenuWillRemoveItem                                      ApplicationEventType
	MenuDidRemoveItem                                       ApplicationEventType
	MenuWillBeginTracking                                   ApplicationEventType
	MenuDidBeginTracking                                    ApplicationEventType
	MenuWillEndTracking                                     ApplicationEventType
	MenuDidEndTracking                                      ApplicationEventType
	MenuWillUpdate                                          ApplicationEventType
	MenuDidUpdate                                           ApplicationEventType
	MenuWillPopUp                                           ApplicationEventType
	MenuDidPopUp                                            ApplicationEventType
	MenuWillSendActionToItem                                ApplicationEventType
	MenuDidSendActionToItem                                 ApplicationEventType
	WebViewDidStartProvisionalNavigation                    WindowEventType
	WebViewDidReceiveServerRedirectForProvisionalNavigation WindowEventType
	WebViewDidFinishNavigation                              WindowEventType
	WebViewDidCommitNavigation                              WindowEventType
}

func newMacEvents() macEvents {
	return macEvents{
		ApplicationDidBecomeActive:               0,
		ApplicationDidChangeBackingProperties:    1,
		ApplicationDidChangeEffectiveAppearance:  2,
		ApplicationDidChangeIcon:                 3,
		ApplicationDidChangeOcclusionState:       4,
		ApplicationDidChangeScreenParameters:     5,
		ApplicationDidChangeStatusBarFrame:       6,
		ApplicationDidChangeStatusBarOrientation: 7,
		ApplicationDidFinishLaunching:            8,
		ApplicationDidHide:                       9,
		ApplicationDidResignActive:               10,
		ApplicationDidUnhide:                     11,
		ApplicationDidUpdate:                     12,
		ApplicationWillBecomeActive:              13,
		ApplicationWillFinishLaunching:           14,
		ApplicationWillHide:                      15,
		ApplicationWillResignActive:              16,
		ApplicationWillTerminate:                 17,
		ApplicationWillUnhide:                    18,
		ApplicationWillUpdate:                    19,
		WindowDidBecomeKey:                       20,
		WindowDidBecomeMain:                      21,
		WindowDidBeginSheet:                      22,
		WindowDidChangeAlpha:                     23,
		WindowDidChangeBackingLocation:           24,
		WindowDidChangeBackingProperties:         25,
		WindowDidChangeCollectionBehavior:        26,
		WindowDidChangeEffectiveAppearance:       27,
		WindowDidChangeOcclusionState:            28,
		WindowDidChangeOrderingMode:              29,
		WindowDidChangeScreen:                    30,
		WindowDidChangeScreenParameters:          31,
		WindowDidChangeScreenProfile:             32,
		WindowDidChangeScreenSpace:               33,
		WindowDidChangeScreenSpaceProperties:     34,
		WindowDidChangeSharingType:               35,
		WindowDidChangeSpace:                     36,
		WindowDidChangeSpaceOrderingMode:         37,
		WindowDidChangeTitle:                     38,
		WindowDidChangeToolbar:                   39,
		WindowDidChangeVisibility:                40,
		WindowDidClose:                           41,
		WindowDidDeminiaturize:                   42,
		WindowDidEndSheet:                        43,
		WindowDidEnterFullScreen:                 44,
		WindowDidEnterVersionBrowser:             45,
		WindowDidExitFullScreen:                  46,
		WindowDidExitVersionBrowser:              47,
		WindowDidExpose:                          48,
		WindowDidFocus:                           49,
		WindowDidMiniaturize:                     50,
		WindowDidMove:                            51,
		WindowDidOrderOffScreen:                  52,
		WindowDidOrderOnScreen:                   53,
		WindowDidResignKey:                       54,
		WindowDidResignMain:                      55,
		WindowDidResize:                          56,
		WindowDidUnfocus:                         57,
		WindowDidUpdate:                          58,
		WindowDidUpdateAlpha:                     59,
		WindowDidUpdateCollectionBehavior:        60,
		WindowDidUpdateCollectionProperties:      61,
		WindowDidUpdateShadow:                    62,
		WindowDidUpdateTitle:                     63,
		WindowDidUpdateToolbar:                   64,
		WindowDidUpdateVisibility:                65,
		WindowWillBecomeKey:                      66,
		WindowWillBecomeMain:                     67,
		WindowWillBeginSheet:                     68,
		WindowWillChangeOrderingMode:             69,
		WindowWillClose:                          70,
		WindowWillDeminiaturize:                  71,
		WindowWillEnterFullScreen:                72,
		WindowWillEnterVersionBrowser:            73,
		WindowWillExitFullScreen:                 74,
		WindowWillExitVersionBrowser:             75,
		WindowWillFocus:                          76,
		WindowWillMiniaturize:                    77,
		WindowWillMove:                           78,
		WindowWillOrderOffScreen:                 79,
		WindowWillOrderOnScreen:                  80,
		WindowWillResignMain:                     81,
		WindowWillResize:                         82,
		WindowWillUnfocus:                        83,
		WindowWillUpdate:                         84,
		WindowWillUpdateAlpha:                    85,
		WindowWillUpdateCollectionBehavior:       86,
		WindowWillUpdateCollectionProperties:     87,
		WindowWillUpdateShadow:                   88,
		WindowWillUpdateTitle:                    89,
		WindowWillUpdateToolbar:                  90,
		WindowWillUpdateVisibility:               91,
		WindowWillUseStandardFrame:               92,
		MenuWillOpen:                             93,
		MenuDidOpen:                              94,
		MenuDidClose:                             95,
		MenuWillSendAction:                       96,
		MenuDidSendAction:                        97,
		MenuWillHighlightItem:                    98,
		MenuDidHighlightItem:                     99,
		MenuWillDisplayItem:                      100,
		MenuDidDisplayItem:                       101,
		MenuWillAddItem:                          102,
		MenuDidAddItem:                           103,
		MenuWillRemoveItem:                       104,
		MenuDidRemoveItem:                        105,
		MenuWillBeginTracking:                    106,
		MenuDidBeginTracking:                     107,
		MenuWillEndTracking:                      108,
		MenuDidEndTracking:                       109,
		MenuWillUpdate:                           110,
		MenuDidUpdate:                            111,
		MenuWillPopUp:                            112,
		MenuDidPopUp:                             113,
		MenuWillSendActionToItem:                 114,
		MenuDidSendActionToItem:                  115,
		WebViewDidStartProvisionalNavigation:     116,
		WebViewDidReceiveServerRedirectForProvisionalNavigation: 117,
		WebViewDidFinishNavigation:                              118,
		WebViewDidCommitNavigation:                              119,
	}
}
