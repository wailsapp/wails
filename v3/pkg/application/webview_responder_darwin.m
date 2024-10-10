#include "webview_responder_darwin.h"
#include "webview_window_darwin.h"

extern void processWindowKeyDownEvent(unsigned int, const char*);

@implementation WebviewResponder
- (WebviewResponder *) initAttachToWindow:(NSWindow *)window {
	self = [super init];

	self.w = window;
	[window setNextResponder:self];

	return self;
}
- (void)keyDown:(NSEvent *)event {
	// TODO: FIX: ctrl+l never reaches this function
	NSUInteger modifierFlags = event.modifierFlags;
    // Create an array to hold the modifier strings
    NSMutableArray *modifierStrings = [NSMutableArray array];
    // Check for modifier flags and add corresponding strings to the array
    if (modifierFlags & NSEventModifierFlagShift) {
        [modifierStrings addObject:@"shift"];
    }
    if (modifierFlags & NSEventModifierFlagControl) {
        [modifierStrings addObject:@"ctrl"];
    }
    if (modifierFlags & NSEventModifierFlagOption) {
        [modifierStrings addObject:@"option"];
    }
    if (modifierFlags & NSEventModifierFlagCommand) {
        [modifierStrings addObject:@"cmd"];
    }
    NSString *keyString = [self keyStringFromEvent:event];
    if (keyString.length > 0) {
        [modifierStrings addObject:keyString];
    }
    // Combine the modifier strings with the key character
    NSString *keyEventString = [modifierStrings componentsJoinedByString:@"+"];
    const char* utf8String = [keyEventString UTF8String];
    WebviewWindowDelegate *delegate = (WebviewWindowDelegate*)self.w.delegate;
    processWindowKeyDownEvent(delegate.windowId, utf8String);
}
- (NSString *)keyStringFromEvent:(NSEvent *)event {
    // Get the pressed key
    switch ([event keyCode]) {
        // Function keys
        case kVK_F1: return @"f1";
        case kVK_F2: return @"f2";
        case kVK_F3: return @"f3";
        case kVK_F4: return @"f4";
        case kVK_F5: return @"f5";
        case kVK_F6: return @"f6";
        case kVK_F7: return @"f7";
        case kVK_F8: return @"f8";
        case kVK_F9: return @"f9";
        case kVK_F10: return @"f10";
        case kVK_F11: return @"f11";
        case kVK_F12: return @"f12";
        case kVK_F13: return @"f13";
        case kVK_F14: return @"f14";
        case kVK_F15: return @"f15";
        case kVK_F16: return @"f16";
        case kVK_F17: return @"f17";
        case kVK_F18: return @"f18";
        case kVK_F19: return @"f19";
        case kVK_F20: return @"f20";
        // Letter keys
        case kVK_ANSI_A: return @"a";
        case kVK_ANSI_B: return @"b";
        case kVK_ANSI_C: return @"c";
        case kVK_ANSI_D: return @"d";
        case kVK_ANSI_E: return @"e";
        case kVK_ANSI_F: return @"f";
        case kVK_ANSI_G: return @"g";
        case kVK_ANSI_H: return @"h";
        case kVK_ANSI_I: return @"i";
        case kVK_ANSI_J: return @"j";
        case kVK_ANSI_K: return @"k";
        case kVK_ANSI_L: return @"l";
        case kVK_ANSI_M: return @"m";
        case kVK_ANSI_N: return @"n";
        case kVK_ANSI_O: return @"o";
        case kVK_ANSI_P: return @"p";
        case kVK_ANSI_Q: return @"q";
        case kVK_ANSI_R: return @"r";
        case kVK_ANSI_S: return @"s";
        case kVK_ANSI_T: return @"t";
        case kVK_ANSI_U: return @"u";
        case kVK_ANSI_V: return @"v";
        case kVK_ANSI_W: return @"w";
        case kVK_ANSI_X: return @"x";
        case kVK_ANSI_Y: return @"y";
        case kVK_ANSI_Z: return @"z";
        // Number keys
        case kVK_ANSI_0: return @"0";
        case kVK_ANSI_1: return @"1";
        case kVK_ANSI_2: return @"2";
        case kVK_ANSI_3: return @"3";
        case kVK_ANSI_4: return @"4";
        case kVK_ANSI_5: return @"5";
        case kVK_ANSI_6: return @"6";
        case kVK_ANSI_7: return @"7";
        case kVK_ANSI_8: return @"8";
        case kVK_ANSI_9: return @"9";
        // Other special keys
        case kVK_Delete: return @"delete";
        case kVK_ForwardDelete: return @"forward delete";
        case kVK_LeftArrow: return @"left";
        case kVK_RightArrow: return @"right";
        case kVK_UpArrow: return @"up";
        case kVK_DownArrow: return @"down";
        case kVK_Tab: return @"tab";
        case kVK_Escape: return @"escape";
        case kVK_Space: return @"space";
        // Punctuation and other keys (for a standard US layout)
        case kVK_ANSI_LeftBracket: return @"[";
        case kVK_ANSI_RightBracket: return @"]";
        case kVK_ANSI_Comma: return @",";
        case kVK_ANSI_Minus: return @"-";
        case kVK_ANSI_Quote: return @"'";
        case kVK_ANSI_Slash: return @"/";
        case kVK_ANSI_Period: return @".";
        case kVK_ANSI_Semicolon: return @";";
        case kVK_ANSI_Equal: return @"=";
        case kVK_ANSI_Grave: return @"`";
        case kVK_ANSI_Backslash: return @"\\";
        default: return [self specialKeyStringFromEvent:event];
    }
}
- (NSString *)specialKeyStringFromEvent:(NSEvent *)event {
    // Check for special keys like escape and tab
    NSString *characters = [event characters];
    if (characters.length == 0) {
        return @"";
    }
	
    if ([characters isEqualToString:@"\r"]) {
        return @"enter";
    }
    if ([characters isEqualToString:@"\b"]) {
        return @"backspace";
    }
    if ([characters isEqualToString:@"\e"]) {
        return @"escape";
    }
    // page down
    if ([characters isEqualToString:@"\x0B"]) {
        return @"page down";
    }
    // page up
    if ([characters isEqualToString:@"\x0E"]) {
        return @"page up";
    }
    // home
    if ([characters isEqualToString:@"\x01"]) {
        return @"home";
    }
    // end
    if ([characters isEqualToString:@"\x04"]) {
        return @"end";
    }
    // clear
    if ([characters isEqualToString:@"\x0C"]) {
        return @"clear";
    }
    // default
	return @"";
}
@end