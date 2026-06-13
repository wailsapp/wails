// ObjC-only helpers shared between mobile_features_ios.m and
// webview_window_ios.m. NOT included by cgo (.go) files — it exposes UIKit
// types that cgo cannot parse.
#ifndef WAILS_MOBILE_FEATURES_IOS_INTERNAL_H
#define WAILS_MOBILE_FEATURES_IOS_INTERNAL_H

#import <UIKit/UIKit.h>

// Orientation lock and status-bar appearance are driven by global state set
// from Go (ios_set_orientation / ios_set_status_bar) and read back by the
// WailsViewController overrides.
UIInterfaceOrientationMask mfSupportedOrientations(void);
UIStatusBarStyle mfStatusBarStyle(void);
BOOL mfStatusBarHidden(void);

#endif
