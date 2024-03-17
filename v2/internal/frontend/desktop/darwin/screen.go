//go:build darwin
// +build darwin

package darwin

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa -framework WebKit -framework AppKit
#import <Foundation/Foundation.h>
#include <AppKit/AppKit.h>
#include <stdlib.h>

#import "Application.h"
#import "WailsContext.h"

typedef struct Screen {
	int isCurrent;
	int isPrimary;
	int height;
	int width;
	int pHeight;
	int pWidth;
} Screen;


int GetNumScreens(){
	return [[NSScreen screens] count];
}

int screenUniqueID(NSScreen *screen){
	// adapted from https://stackoverflow.com/a/1237490/4188138
    NSDictionary* screenDictionary = [screen deviceDescription];
    NSNumber* screenID = [screenDictionary objectForKey:@"NSScreenNumber"];
    CGDirectDisplayID aID = [screenID unsignedIntValue];
	return aID;
}

Screen GetNthScreen(int nth, void *inctx){
	WailsContext *ctx = (__bridge WailsContext*) inctx;
	NSArray<NSScreen *> *screens = [NSScreen screens];
	NSScreen* nthScreen = [screens objectAtIndex:nth];
	NSScreen* currentScreen = [ctx getCurrentScreen];

	Screen returnScreen;
	returnScreen.isCurrent = (int)(screenUniqueID(currentScreen)==screenUniqueID(nthScreen));
	// TODO properly handle screen mirroring
	// from apple documentation:
	// https://developer.apple.com/documentation/appkit/nsscreen/1388393-screens?language=objc
	// The screen at index 0 in the returned array corresponds to the primary screen of the user’s system. This is the screen that contains the menu bar and whose origin is at the point (0, 0). In the case of mirroring, the first screen is the largest drawable display; if all screens are the same size, it is the screen with the highest pixel depth. This primary screen may not be the same as the one returned by the mainScreen method, which returns the screen with the active window.
	returnScreen.isPrimary = nth==0;
	returnScreen.height = (int) nthScreen.frame.size.height;
	returnScreen.width =  (int) nthScreen.frame.size.width;

	returnScreen.pWidth = 0;
	returnScreen.pHeight = 0;

	// https://stackoverflow.com/questions/13859109/how-to-programmatically-determine-native-pixel-resolution-of-retina-macbook-pro
	CGDirectDisplayID sid = ((NSNumber *)[nthScreen.deviceDescription
    	objectForKey:@"NSScreenNumber"]).unsignedIntegerValue;

	CFArrayRef ms = CGDisplayCopyAllDisplayModes(sid, NULL);
	CFIndex n = CFArrayGetCount(ms);
	for (int i = 0; i < n; i++) {
		CGDisplayModeRef m = (CGDisplayModeRef) CFArrayGetValueAtIndex(ms, i);
		if (CGDisplayModeGetIOFlags(m) & kDisplayModeNativeFlag) {
			// This corresponds with "System Settings" -> General -> About -> Displays
			returnScreen.pWidth = CGDisplayModeGetPixelWidth(m);
			returnScreen.pHeight = CGDisplayModeGetPixelHeight(m);
			break;
		}
	}
	CFRelease(ms);

	if (returnScreen.pWidth == 0 || returnScreen.pHeight == 0) {
		// If there was no native resolution take a best fit approach and use the backing pixel size.
		NSRect pSize = [nthScreen convertRectToBacking:nthScreen.frame];
		returnScreen.pHeight = (int) pSize.size.height;
		returnScreen.pWidth = (int) pSize.size.width;
	}
	return returnScreen;
}

*/
import "C"

import (
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/frontend"
)

func GetAllScreens(wailsContext unsafe.Pointer) ([]frontend.Screen, error) {
	err := error(nil)
	screens := []frontend.Screen{}
	numScreens := int(C.GetNumScreens())
	for screeNum := 0; screeNum < numScreens; screeNum++ {
		screenNumC := C.int(screeNum)
		cScreen := C.GetNthScreen(screenNumC, wailsContext)

		screen := frontend.Screen{
			Height:    int(cScreen.height),
			Width:     int(cScreen.width),
			IsCurrent: cScreen.isCurrent == C.int(1),
			IsPrimary: cScreen.isPrimary == C.int(1),

			Size: frontend.ScreenSize{
				Height: int(cScreen.height),
				Width:  int(cScreen.width),
			},
			PhysicalSize: frontend.ScreenSize{
				Height: int(cScreen.pHeight),
				Width:  int(cScreen.pWidth),
			},
		}
		screens = append(screens, screen)
	}
	return screens, err
}
