//go:build darwin && !ios && !server

#include "Cocoa/Cocoa.h"
#include "menuitem_darwin.h"
#include "systemtray_darwin.h"

extern void systrayClickCallback(long, int);
extern int systrayPreClickCallback(long, int);

// StatusItemController.m
@implementation StatusItemController

- (void)statusItemClicked:(id)sender {
	NSEvent *event = [NSApp currentEvent];
	systrayClickCallback(self.id, event.type);
}

- (void)menuDidClose:(NSMenu *)menu {
	// Remove the menu from the status item so future clicks invoke the
	// action handler instead of re-showing the menu.
	self.statusItem.menu = nil;
	menu.delegate = nil;
}

@end

@implementation WailsStatusItemGestureObserver

// Recognizers see the event before the button does, so assigning
// statusItem.menu here makes the button enter native menu tracking for this
// very click — proper highlight, and the app is not activated.
- (void)handlePreClick:(NSEvent *)event {
	StatusItemController *controller = self.controller;
	if (controller == nil) {
		return;
	}
	int action = systrayPreClickCallback((long)controller.id, (int)event.type);
	if (action == 1 && controller.cachedMenu != nil) {
		controller.cachedMenu.delegate = controller;
		controller.statusItem.menu = controller.cachedMenu;
	}
}

- (void)mouseDown:(NSEvent *)event {
	[super mouseDown:event];
	[self handlePreClick:event];
}

- (void)rightMouseDown:(NSEvent *)event {
	[super rightMouseDown:event];
	[self handlePreClick:event];
}

- (void)mouseUp:(NSEvent *)event {
	[super mouseUp:event];
	self.state = NSGestureRecognizerStateFailed;
}

- (void)rightMouseUp:(NSEvent *)event {
	[super rightMouseUp:event];
	self.state = NSGestureRecognizerStateFailed;
}

@end

// Create a new system tray
void* systemTrayNew(long id) {
	StatusItemController *controller = [[StatusItemController alloc] init];
	controller.id = id;
	NSStatusItem *statusItem = [[[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength] retain];
	controller.statusItem = statusItem;
	[statusItem setTarget:controller];
	[statusItem setAction:@selector(statusItemClicked:)];
	NSButton *button = statusItem.button;
	[button sendActionOn:(NSEventMaskLeftMouseDown|NSEventMaskRightMouseDown)];

	if (@available(macOS 27.0, *)) {
		// macOS 27: observe the button's mouse-downs through the
		// gesture-recognizer system so the pre-click callback can attach a
		// menu before the button starts native tracking; NSEvent monitors
		// miss Sidecar/touch input there (TN3212).
		WailsStatusItemGestureObserver *observer = [[[WailsStatusItemGestureObserver alloc] initWithTarget:nil action:NULL] autorelease];
		observer.controller = controller;
		observer.delaysPrimaryMouseButtonEvents = NO;
		[button addGestureRecognizer:observer];
		controller.gestureObserver = observer;
	} else {
		// Install a local event monitor that fires BEFORE the button processes
		// the mouse-down.  When the pre-click callback says "show menu", we
		// temporarily set statusItem.menu so the button enters native menu
		// tracking — this gives proper highlight and does not activate the app.
		controller.eventMonitor = [NSEvent addLocalMonitorForEventsMatchingMask:
			(NSEventMaskLeftMouseDown|NSEventMaskRightMouseDown)
			handler:^NSEvent *(NSEvent *event) {
				if (event.window != button.window) return event;

				int action = systrayPreClickCallback((long)controller.id, (int)event.type);
				if (action == 1 && controller.cachedMenu != nil) {
					controller.cachedMenu.delegate = controller;
					statusItem.menu = controller.cachedMenu;
				}
				return event;
			}];
	}

	return (void*)statusItem;
}

void systemTraySetLabel(void* nsStatusItem, char *label) {
	if( label == NULL ) {
		return;
	}
	// Set the label on the main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		NSStatusItem *statusItem = (NSStatusItem *)nsStatusItem;
		statusItem.button.title = [NSString stringWithUTF8String:label];
		free(label);
	});
}

void systemTraySetANSILabel(void* nsStatusItem, void* label) {
    if( label == NULL ) {
        return;
    }

    NSMutableAttributedString* attributedString = (NSMutableAttributedString*) label;

    // Set the label on the main thread.
    dispatch_async(dispatch_get_main_queue(), ^{
        // setAttributedTitle: copies the string, so drop the owning
        // reference we accumulated while building it.
        NSStatusItem *statusItem = (NSStatusItem *)nsStatusItem;
        [statusItem setAttributedTitle:attributedString];
        [attributedString release];
    });
}

void* appendAttributedString(void *currentString, char *title, char *FG, char *BG) {

    NSMutableAttributedString* newString = createAttributedString(title, FG, BG);
    if( currentString != NULL ) {
        NSMutableAttributedString* current = (NSMutableAttributedString*)currentString;
        [current appendAttributedString:newString];
        // appendAttributedString: copies the run; drop the +1 from
        // createAttributedString
        [newString release];
        newString = current;
    }

    return (void*)newString;
}

void* createAttributedString(char *title, char *FG, char *BG) {

    NSMutableDictionary *dictionary = [NSMutableDictionary new];

    // RGBA
    if(FG != NULL && strlen(FG) > 0) {
        unsigned short r, g, b, a;

        // white by default
        r = g = b = a = 255;
        int count = sscanf(FG, "#%02hx%02hx%02hx%02hx", &r, &g, &b, &a);
        if (count > 0) {
            NSColor *colour = [NSColor colorWithCalibratedRed:(CGFloat)r / 255.0
                                                         green:(CGFloat)g / 255.0
                                                          blue:(CGFloat)b / 255.0
                                                         alpha:(CGFloat)a / 255.0];
            dictionary[NSForegroundColorAttributeName] = colour;

        }
    }

    // Calculate BG colour
    if(BG != NULL && strlen(BG) > 0) {
            unsigned short r, g, b, a;

            // white by default
            r = g = b = a = 255;
            int count = sscanf(BG, "#%02hx%02hx%02hx%02hx", &r, &g, &b, &a);
            if (count > 0) {
                NSColor *colour = [NSColor colorWithCalibratedRed:(CGFloat)r / 255.0
                                                             green:(CGFloat)g / 255.0
                                                              blue:(CGFloat)b / 255.0
                                                             alpha:(CGFloat)a / 255.0];
                dictionary[NSBackgroundColorAttributeName] = colour;
            }
    }
    NSMutableAttributedString *attributedString = [[NSMutableAttributedString alloc] initWithString:[NSString stringWithUTF8String:title] attributes:dictionary];
    // The attributed string keeps its own copy of the attributes
    [dictionary release];
    return (void*)attributedString;
}

// Create an nsimage from a byte array
NSImage* imageFromBytes(const unsigned char *bytes, int length) {
	NSData *data = [NSData dataWithBytes:bytes length:length];
	NSImage *image = [[NSImage alloc] initWithData:data];
	return image;
}

// Set the icon on the system tray
void systemTraySetIcon(void* nsStatusItem, void* nsImage, int position, bool isTemplate) {
	// Set the icon on the main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		NSStatusItem *statusItem = (NSStatusItem *)nsStatusItem;
		NSImage *image = (NSImage *)nsImage;

		NSStatusBar *statusBar = [NSStatusBar systemStatusBar];
		CGFloat thickness = [statusBar thickness];
		[image setSize:NSMakeSize(thickness, thickness)];
		if( isTemplate ) {
			[image setTemplate:YES];
		}
		statusItem.button.image = [image autorelease];
		statusItem.button.imagePosition = position;
	});
}

// Destroy system tray
void systemTrayDestroy(void* nsStatusItem) {
	// Remove the status item from the status bar and its associated menu
	dispatch_async(dispatch_get_main_queue(), ^{
		NSStatusItem *statusItem = (NSStatusItem *)nsStatusItem;
		StatusItemController *controller = (StatusItemController *)[statusItem target];
		if (controller.eventMonitor) {
			[NSEvent removeMonitor:controller.eventMonitor];
			controller.eventMonitor = nil;
		}
		if (controller.gestureObserver) {
			[controller.gestureObserver.view removeGestureRecognizer:controller.gestureObserver];
			controller.gestureObserver = nil;
		}
		[[NSStatusBar systemStatusBar] removeStatusItem:statusItem];
		[controller release];
		[statusItem release];
	});
}

// showMenu is used for programmatic OpenMenu() calls.  Click-triggered
// menus are handled by the pre-click hook installed in systemTrayNew.
void showMenu(void* nsStatusItem, void *nsMenu) {
	dispatch_async(dispatch_get_main_queue(), ^{
		NSStatusItem *statusItem = (NSStatusItem *)nsStatusItem;
		NSMenu *menu = (NSMenu *)nsMenu;
		StatusItemController *controller = (StatusItemController *)[statusItem target];

		// Temporarily assign the menu for native tracking.
		menu.delegate = controller;
		statusItem.menu = menu;

		if (@available(macOS 27.0, *)) {
			// With a menu assigned, a click on the button enters native menu
			// tracking (highlights the button, blocks until dismissed)
			// without activating the app. performClick: goes through the
			// button's standard click path; macOS 27 no longer honours
			// synthesised mouseDown: NSEvents (TN3212).
			[statusItem.button performClick:nil];
		} else {
			// Synthesize a mouse-down at the button centre to trigger native
			// menu tracking (highlights the button, blocks until dismissed).
			NSRect frame = [statusItem.button convertRect:statusItem.button.bounds toView:nil];
			NSPoint loc = NSMakePoint(NSMidX(frame), NSMidY(frame));
			NSEvent *event = [NSEvent mouseEventWithType:NSEventTypeLeftMouseDown
			                                    location:loc
			                               modifierFlags:0
			                                   timestamp:[[NSProcessInfo processInfo] systemUptime]
			                                windowNumber:statusItem.button.window.windowNumber
			                                     context:nil
			                                 eventNumber:0
			                                  clickCount:1
			                                    pressure:1.0];
			[statusItem.button mouseDown:event];
		}

		// Menu dismissed — restore custom click handling.
		statusItem.menu = nil;
		menu.delegate = nil;
	});
}

void systemTraySetCachedMenu(void* nsStatusItem, void *nsMenu) {
	NSStatusItem *statusItem = (NSStatusItem *)nsStatusItem;
	StatusItemController *controller = (StatusItemController *)[statusItem target];
	controller.cachedMenu = (NSMenu *)nsMenu;
}

void systemTrayGetBounds(void* nsStatusItem, NSRect *rect, void **outScreen) {
    NSStatusItem *statusItem = (NSStatusItem *)nsStatusItem;
    NSStatusBarButton *button = statusItem.button;
    
    // Get mouse location and find the screen it's on
    NSPoint mouseLocation = [NSEvent mouseLocation];
    NSScreen *screen = nil;
    NSArray *screens = [NSScreen screens];
    
    for (NSScreen *candidate in screens) {
        NSRect frame = [candidate frame];
        if (NSPointInRect(mouseLocation, frame)) {
            screen = candidate;
            break;
        }
    }
    if (!screen) {
        screen = [NSScreen mainScreen];
    }
    
    // Get button frame in screen coordinates
    NSRect buttonFrame = button.frame;
    NSRect buttonFrameScreen = [button.window convertRectToScreen:buttonFrame];
    
    *rect = buttonFrameScreen;
    *outScreen = (void*)screen;
}

NSRect NSScreen_frame(void* screen) {
    return [(NSScreen*)screen frame];
}

int statusBarHeight() {
    NSMenu *mainMenu = [NSApp mainMenu];
    CGFloat menuBarHeight = [mainMenu menuBarHeight];
    return (int)menuBarHeight;
}

void systemTrayPositionWindow(void* nsStatusItem, void* nsWindow, int offset) {
    // Get the status item's button
    NSStatusBarButton *button = [(NSStatusItem*)nsStatusItem button];

    // Get the frame in screen coordinates
    NSRect frame = [button.window convertRectToScreen:button.frame];

    // Get the screen that contains the status item
    NSScreen *screen = [button.window screen];
    if (screen == nil) {
        screen = [NSScreen mainScreen];
    }

    // Get screen's backing scale factor (DPI)
    CGFloat scaleFactor = [screen backingScaleFactor];

    // Get the window's frame
    NSRect windowFrame = [(NSWindow*)nsWindow frame];

    // Calculate the horizontal position (centered under the status item)
    CGFloat windowX = frame.origin.x + (frame.size.width - windowFrame.size.width) / 2;

    // If the window would go off the right edge of the screen, adjust it
    if (windowX + windowFrame.size.width > screen.frame.origin.x + screen.frame.size.width) {
        windowX = screen.frame.origin.x + screen.frame.size.width - windowFrame.size.width;
    }
    // If the window would go off the left edge of the screen, adjust it
    if (windowX < screen.frame.origin.x) {
        windowX = screen.frame.origin.x;
    }

    // Get screen metrics
    NSRect screenFrame = [screen frame];
    NSRect visibleFrame = [screen visibleFrame];

    // Calculate the vertical position
    CGFloat scaledOffset = offset * scaleFactor;
    CGFloat windowY = visibleFrame.origin.y + visibleFrame.size.height - windowFrame.size.height - scaledOffset;

    // Set the window's frame
    windowFrame.origin.x = windowX;
    windowFrame.origin.y = windowY;
    [(NSWindow*)nsWindow setFrame:windowFrame display:YES animate:NO];

    // Set window level to popup menu level so it appears above other windows
    [(NSWindow*)nsWindow setLevel:NSPopUpMenuWindowLevel];

    // Bring window to front
    [(NSWindow*)nsWindow orderFrontRegardless];
}
