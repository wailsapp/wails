//go:build darwin

#ifndef _DIALOGS_DELEGATE_H_
#define _DIALOGS_DELEGATE_H_

#import <Cocoa/Cocoa.h>

// Conditionally import UniformTypeIdentifiers based on OS version
#if __MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
#import <UniformTypeIdentifiers/UTType.h>
#endif

// OpenPanel delegate to handle file filtering
@interface OpenPanelDelegate : NSObject <NSOpenSavePanelDelegate>
@property (nonatomic, strong) NSArray *allowedExtensions;
@end

#endif