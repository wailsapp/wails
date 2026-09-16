//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>

static void hapticsPerformPattern(int kind) {
	NSHapticFeedbackPattern pattern = NSHapticFeedbackPatternGeneric;
	switch (kind) {
		case 1: pattern = NSHapticFeedbackPatternAlignment; break;
		case 2: pattern = NSHapticFeedbackPatternLevelChange; break;
		default: break;
	}
	[[NSHapticFeedbackManager defaultPerformer] performFeedbackPattern:pattern
	                                                   performanceTime:NSHapticFeedbackPerformanceTimeNow];
}
*/
import "C"

func hapticsPerform(kind HapticKind) {
	InvokeAsync(func() {
		C.hapticsPerformPattern(C.int(kind))
	})
}

func hapticsSupported() bool {
	return true
}
