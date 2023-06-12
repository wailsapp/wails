package events

type ApplicationEventType uint
type WindowEventType uint

const (
	FilesDropped WindowEventType = iota
)

var Common = newCommonEvents()

type commonEvents struct {
	ApplicationStarted ApplicationEventType
	WindowMaximise     WindowEventType
	WindowUnMaximise   WindowEventType
	WindowFullscreen   WindowEventType
	WindowUnFullscreen WindowEventType
	WindowRestore      WindowEventType
	WindowMinimise     WindowEventType
	WindowUnMinimise   WindowEventType
	WindowClose        WindowEventType
	WindowZoom         WindowEventType
	WindowZoomIn       WindowEventType
	WindowZoomOut      WindowEventType
	WindowZoomReset    WindowEventType
	WindowFocus        WindowEventType
	WindowShow         WindowEventType
	WindowHide         WindowEventType
	WindowDPIChanged   WindowEventType
}

func newCommonEvents() commonEvents {
	return commonEvents{
		ApplicationStarted: 1154,
		WindowMaximise:     1155,
		WindowUnMaximise:   1156,
		WindowFullscreen:   1157,
		WindowUnFullscreen: 1158,
		WindowRestore:      1159,
		WindowMinimise:     1160,
		WindowUnMinimise:   1161,
		WindowClose:        1162,
		WindowZoom:         1163,
		WindowZoomIn:       1164,
		WindowZoomOut:      1165,
		WindowZoomReset:    1166,
		WindowFocus:        1167,
		WindowShow:         1168,
		WindowHide:         1169,
		WindowDPIChanged:   1170,
	}
}

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
	WindowShouldClose                                       WindowEventType
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
	WindowFileDraggingEntered                               WindowEventType
	WindowFileDraggingPerformed                             WindowEventType
	WindowFileDraggingExited                                WindowEventType
}

func newMacEvents() macEvents {
	return macEvents{
		ApplicationDidBecomeActive:               1024,
		ApplicationDidChangeBackingProperties:    1025,
		ApplicationDidChangeEffectiveAppearance:  1026,
		ApplicationDidChangeIcon:                 1027,
		ApplicationDidChangeOcclusionState:       1028,
		ApplicationDidChangeScreenParameters:     1029,
		ApplicationDidChangeStatusBarFrame:       1030,
		ApplicationDidChangeStatusBarOrientation: 1031,
		ApplicationDidFinishLaunching:            1032,
		ApplicationDidHide:                       1033,
		ApplicationDidResignActive:               1034,
		ApplicationDidUnhide:                     1035,
		ApplicationDidUpdate:                     1036,
		ApplicationWillBecomeActive:              1037,
		ApplicationWillFinishLaunching:           1038,
		ApplicationWillHide:                      1039,
		ApplicationWillResignActive:              1040,
		ApplicationWillTerminate:                 1041,
		ApplicationWillUnhide:                    1042,
		ApplicationWillUpdate:                    1043,
		WindowDidBecomeKey:                       1044,
		WindowDidBecomeMain:                      1045,
		WindowDidBeginSheet:                      1046,
		WindowDidChangeAlpha:                     1047,
		WindowDidChangeBackingLocation:           1048,
		WindowDidChangeBackingProperties:         1049,
		WindowDidChangeCollectionBehavior:        1050,
		WindowDidChangeEffectiveAppearance:       1051,
		WindowDidChangeOcclusionState:            1052,
		WindowDidChangeOrderingMode:              1053,
		WindowDidChangeScreen:                    1054,
		WindowDidChangeScreenParameters:          1055,
		WindowDidChangeScreenProfile:             1056,
		WindowDidChangeScreenSpace:               1057,
		WindowDidChangeScreenSpaceProperties:     1058,
		WindowDidChangeSharingType:               1059,
		WindowDidChangeSpace:                     1060,
		WindowDidChangeSpaceOrderingMode:         1061,
		WindowDidChangeTitle:                     1062,
		WindowDidChangeToolbar:                   1063,
		WindowDidChangeVisibility:                1064,
		WindowDidDeminiaturize:                   1065,
		WindowDidEndSheet:                        1066,
		WindowDidEnterFullScreen:                 1067,
		WindowDidEnterVersionBrowser:             1068,
		WindowDidExitFullScreen:                  1069,
		WindowDidExitVersionBrowser:              1070,
		WindowDidExpose:                          1071,
		WindowDidFocus:                           1072,
		WindowDidMiniaturize:                     1073,
		WindowDidMove:                            1074,
		WindowDidOrderOffScreen:                  1075,
		WindowDidOrderOnScreen:                   1076,
		WindowDidResignKey:                       1077,
		WindowDidResignMain:                      1078,
		WindowDidResize:                          1079,
		WindowDidUnfocus:                         1080,
		WindowDidUpdate:                          1081,
		WindowDidUpdateAlpha:                     1082,
		WindowDidUpdateCollectionBehavior:        1083,
		WindowDidUpdateCollectionProperties:      1084,
		WindowDidUpdateShadow:                    1085,
		WindowDidUpdateTitle:                     1086,
		WindowDidUpdateToolbar:                   1087,
		WindowDidUpdateVisibility:                1088,
		WindowShouldClose:                        1089,
		WindowWillBecomeKey:                      1090,
		WindowWillBecomeMain:                     1091,
		WindowWillBeginSheet:                     1092,
		WindowWillChangeOrderingMode:             1093,
		WindowWillClose:                          1094,
		WindowWillDeminiaturize:                  1095,
		WindowWillEnterFullScreen:                1096,
		WindowWillEnterVersionBrowser:            1097,
		WindowWillExitFullScreen:                 1098,
		WindowWillExitVersionBrowser:             1099,
		WindowWillFocus:                          1100,
		WindowWillMiniaturize:                    1101,
		WindowWillMove:                           1102,
		WindowWillOrderOffScreen:                 1103,
		WindowWillOrderOnScreen:                  1104,
		WindowWillResignMain:                     1105,
		WindowWillResize:                         1106,
		WindowWillUnfocus:                        1107,
		WindowWillUpdate:                         1108,
		WindowWillUpdateAlpha:                    1109,
		WindowWillUpdateCollectionBehavior:       1110,
		WindowWillUpdateCollectionProperties:     1111,
		WindowWillUpdateShadow:                   1112,
		WindowWillUpdateTitle:                    1113,
		WindowWillUpdateToolbar:                  1114,
		WindowWillUpdateVisibility:               1115,
		WindowWillUseStandardFrame:               1116,
		MenuWillOpen:                             1117,
		MenuDidOpen:                              1118,
		MenuDidClose:                             1119,
		MenuWillSendAction:                       1120,
		MenuDidSendAction:                        1121,
		MenuWillHighlightItem:                    1122,
		MenuDidHighlightItem:                     1123,
		MenuWillDisplayItem:                      1124,
		MenuDidDisplayItem:                       1125,
		MenuWillAddItem:                          1126,
		MenuDidAddItem:                           1127,
		MenuWillRemoveItem:                       1128,
		MenuDidRemoveItem:                        1129,
		MenuWillBeginTracking:                    1130,
		MenuDidBeginTracking:                     1131,
		MenuWillEndTracking:                      1132,
		MenuDidEndTracking:                       1133,
		MenuWillUpdate:                           1134,
		MenuDidUpdate:                            1135,
		MenuWillPopUp:                            1136,
		MenuDidPopUp:                             1137,
		MenuWillSendActionToItem:                 1138,
		MenuDidSendActionToItem:                  1139,
		WebViewDidStartProvisionalNavigation:     1140,
		WebViewDidReceiveServerRedirectForProvisionalNavigation: 1141,
		WebViewDidFinishNavigation:                              1142,
		WebViewDidCommitNavigation:                              1143,
		WindowFileDraggingEntered:                               1144,
		WindowFileDraggingPerformed:                             1145,
		WindowFileDraggingExited:                                1146,
	}
}

var Windows = newWindowsEvents()

type windowsEvents struct {
	SystemThemeChanged         ApplicationEventType
	APMPowerStatusChange       ApplicationEventType
	APMSuspend                 ApplicationEventType
	APMResumeAutomatic         ApplicationEventType
	APMResumeSuspend           ApplicationEventType
	APMPowerSettingChange      ApplicationEventType
	WebViewNavigationCompleted WindowEventType
}

func newWindowsEvents() windowsEvents {
	return windowsEvents{
		SystemThemeChanged:         1147,
		APMPowerStatusChange:       1148,
		APMSuspend:                 1149,
		APMResumeAutomatic:         1150,
		APMResumeSuspend:           1151,
		APMPowerSettingChange:      1152,
		WebViewNavigationCompleted: 1153,
	}
}
