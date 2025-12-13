//go:build ios

#ifndef _events_h
#define _events_h

extern void processApplicationEvent(unsigned int, void* data);
extern void processWindowEvent(unsigned int, unsigned int);

#define EventApplicationDidBecomeActive 1238
#define EventApplicationDidEnterBackground 1239
#define EventApplicationDidFinishLaunching 1240
#define EventApplicationDidReceiveMemoryWarning 1241
#define EventApplicationWillEnterForeground 1242
#define EventApplicationWillResignActive 1243
#define EventApplicationWillTerminate 1244
#define EventWindowDidLoad 1245
#define EventWindowWillAppear 1246
#define EventWindowDidAppear 1247
#define EventWindowWillDisappear 1248
#define EventWindowDidDisappear 1249
#define EventWindowSafeAreaInsetsChanged 1250
#define EventWindowOrientationChanged 1251
#define EventWindowTouchBegan 1252
#define EventWindowTouchMoved 1253
#define EventWindowTouchEnded 1254
#define EventWindowTouchCancelled 1255
#define EventWebViewDidStartNavigation 1256
#define EventWebViewDidFinishNavigation 1257
#define EventWebViewDidFailNavigation 1258
#define EventWebViewDecidePolicyForNavigationAction 1259

#define MAX_EVENTS 1260


#endif