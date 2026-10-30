// Measures how WKPreferences.inactiveSchedulingPolicy affects a WKWebView in a
// hidden or miniaturised NSWindow while another window of the app is key.
// The page runs a 100ms setInterval; after a settling delay the probe counts
// the ticks it receives in 10 seconds (100 = unthrottled).
//
// Build: clang -fobjc-arc -framework Cocoa -framework WebKit -o measure measure.m
// Usage: ./measure <default|suspend|throttle|none> <nevershown|orderout|miniaturise|offscreen|visiblekey>
// Environment:
//   WORKER=1   run the timer in a dedicated worker instead of the page
//   BUSY=1     run a 50ms busy loop per tick in a worker; reports iterations/1000
//   RAF=1      count requestAnimationFrame callbacks instead of timer ticks
//   DELAY=<s>  settling delay before counting (default 12, including 2s setup)
#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

static NSString *policyArg, *stateArg;

@interface Probe : NSObject <NSApplicationDelegate, WKScriptMessageHandler>
@property NSWindow *worker, *front;
@property WKWebView *web;
@property NSInteger ticks;
@property BOOL counting;
@end

@implementation Probe
- (void)userContentController:(WKUserContentController *)c didReceiveScriptMessage:(WKScriptMessage *)m {
    if (self.counting) self.ticks += getenv("BUSY") ? [m.body integerValue] / 1000 : 1;
}

- (void)applicationDidFinishLaunching:(NSNotification *)n {
    WKWebViewConfiguration *cfg = [WKWebViewConfiguration new];
    [cfg.userContentController addScriptMessageHandler:self name:@"tick"];
    if (@available(macOS 14.0, *)) {
        if ([policyArg isEqualToString:@"none"]) cfg.preferences.inactiveSchedulingPolicy = WKInactiveSchedulingPolicyNone;
        if ([policyArg isEqualToString:@"throttle"]) cfg.preferences.inactiveSchedulingPolicy = WKInactiveSchedulingPolicyThrottle;
        if ([policyArg isEqualToString:@"suspend"]) cfg.preferences.inactiveSchedulingPolicy = WKInactiveSchedulingPolicySuspend;
    }
    NSRect r = NSMakeRect(100, 100, 400, 300);
    NSWindowStyleMask style = NSWindowStyleMaskTitled | NSWindowStyleMaskClosable | NSWindowStyleMaskMiniaturizable;
    self.worker = [[NSWindow alloc] initWithContentRect:r styleMask:style backing:NSBackingStoreBuffered defer:NO];
    self.worker.releasedWhenClosed = NO;
    self.web = [[WKWebView alloc] initWithFrame:r configuration:cfg];
    self.worker.contentView = self.web;
    NSString *html = getenv("RAF")
        ? @"<script>const f = () => { window.webkit.messageHandlers.tick.postMessage(1); requestAnimationFrame(f); }; requestAnimationFrame(f);</script>"
        : getenv("BUSY")
        ? @"<script>const w = new Worker(URL.createObjectURL(new Blob(['setInterval(() => { const end = performance.now() + 50; let n = 0; while (performance.now() < end) n++; postMessage(n); }, 100)'], {type: 'text/javascript'}))); w.onmessage = (e) => window.webkit.messageHandlers.tick.postMessage(e.data);</script>"
        : getenv("WORKER")
        ? @"<script>const w = new Worker(URL.createObjectURL(new Blob(['setInterval(() => postMessage(1), 100)'], {type: 'text/javascript'}))); w.onmessage = () => window.webkit.messageHandlers.tick.postMessage(1);</script>"
        : @"<script>setInterval(() => window.webkit.messageHandlers.tick.postMessage(1), 100);</script>";
    [self.web loadHTMLString:html baseURL:[NSURL URLWithString:@"https://probe.invalid/"]];

    self.front = [[NSWindow alloc] initWithContentRect:NSMakeRect(600, 100, 300, 200) styleMask:style backing:NSBackingStoreBuffered defer:NO];
    self.front.releasedWhenClosed = NO;
    [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
    [NSApp activateIgnoringOtherApps:YES];

    if (![stateArg isEqualToString:@"nevershown"]) [self.worker makeKeyAndOrderFront:nil];
    dispatch_after(dispatch_time(DISPATCH_TIME_NOW, 2 * NSEC_PER_SEC), dispatch_get_main_queue(), ^{
        if ([stateArg isEqualToString:@"orderout"]) [self.worker orderOut:nil];
        if ([stateArg isEqualToString:@"miniaturise"]) [self.worker miniaturize:nil];
        if ([stateArg isEqualToString:@"offscreen"]) [self.worker setFrameOrigin:NSMakePoint(-10000, -10000)];
        if (![stateArg isEqualToString:@"visiblekey"]) [self.front makeKeyAndOrderFront:nil];
    });
    // Let throttling settle, then count ticks for 10s (100 expected unthrottled).
    dispatch_after(dispatch_time(DISPATCH_TIME_NOW, (int64_t)(getenv("DELAY") ? atoi(getenv("DELAY")) : 12) * NSEC_PER_SEC), dispatch_get_main_queue(), ^{
        self.ticks = 0; self.counting = YES;
    });
    dispatch_after(dispatch_time(DISPATCH_TIME_NOW, (int64_t)((getenv("DELAY") ? atoi(getenv("DELAY")) : 12) + 10) * NSEC_PER_SEC), dispatch_get_main_queue(), ^{
        printf("%-9s %-12s ticks/10s=%ld visible=%d miniaturized=%d\n", policyArg.UTF8String, stateArg.UTF8String,
               (long)self.ticks, self.worker.isVisible, self.worker.isMiniaturized);
        exit(0);
    });
}
@end

int main(int argc, const char **argv) {
    @autoreleasepool {
        policyArg = [NSString stringWithUTF8String:argv[1]];
        stateArg = [NSString stringWithUTF8String:argv[2]];
        NSApplication *app = [NSApplication sharedApplication];
        Probe *p = [Probe new];
        app.delegate = p;
        [app run];
    }
    return 0;
}
