//go:build ios

#ifndef application_ios_delegate_h
#define application_ios_delegate_h

#import <UIKit/UIKit.h>

@class WailsViewController;

@interface WailsAppDelegate : UIResponder <UIApplicationDelegate>
@property (strong, nonatomic) UIWindow *window;
@property (nonatomic, strong) NSMutableArray<WailsViewController *> *viewControllers;
@end

#endif /* application_ios_delegate_h */