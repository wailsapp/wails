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
        case 122: return @"f1";
        case 120: return @"f2";
        case 99: return @"f3";
        case 118: return @"f4";
        case 96: return @"f5";
        case 97: return @"f6";
        case 98: return @"f7";
        case 100: return @"f8";
        case 101: return @"f9";
        case 109: return @"f10";
        case 103: return @"f11";
        case 111: return @"f12";
        case 105: return @"f13";
        case 107: return @"f14";
        case 113: return @"f15";
        case 106: return @"f16";
        case 64: return @"f17";
        case 79: return @"f18";
        case 80: return @"f19";
        case 90: return @"f20";
        // Letter keys
        case 0: return @"a";
        case 11: return @"b";
        case 8: return @"c";
        case 2: return @"d";
        case 14: return @"e";
        case 3: return @"f";
        case 5: return @"g";
        case 4: return @"h";
        case 34: return @"i";
        case 38: return @"j";
        case 40: return @"k";
        case 37: return @"l";
        case 46: return @"m";
        case 45: return @"n";
        case 31: return @"o";
        case 35: return @"p";
        case 12: return @"q";
        case 15: return @"r";
        case 1: return @"s";
        case 17: return @"t";
        case 32: return @"u";
        case 9: return @"v";
        case 13: return @"w";
        case 7: return @"x";
        case 16: return @"y";
        case 6: return @"z";
        // Number keys
        case 29: return @"0";
        case 18: return @"1";
        case 19: return @"2";
        case 20: return @"3";
        case 21: return @"4";
        case 23: return @"5";
        case 22: return @"6";
        case 26: return @"7";
        case 28: return @"8";
        case 25: return @"9";
        // Other special keys
        case 51: return @"delete";
        case 117: return @"forward delete";
        case 123: return @"left";
        case 124: return @"right";
        case 126: return @"up";
        case 125: return @"down";
        case 48: return @"tab";
        case 53: return @"escape";
        case 49: return @"space";
        // Punctuation and other keys (for a standard US layout)
        case 33: return @"[";
        case 30: return @"]";
        case 43: return @",";
        case 27: return @"-";
        case 39: return @"'";
        case 44: return @"/";
        case 47: return @".";
        case 41: return @";";
        case 24: return @"=";
        case 50: return @"`";
        case 42: return @"\\";
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