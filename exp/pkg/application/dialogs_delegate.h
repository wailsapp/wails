//go:build darwin

#ifndef _DIALOGS_DELEGATE_H_
#define _DIALOGS_DELEGATE_H_

#import <UniformTypeIdentifiers/UTType.h>
#import <Cocoa/Cocoa.h>

// create an NSOpenPanel delegate to handle the callback
@interface OpenPanelDelegate : NSObject <NSOpenSavePanelDelegate>
@property (nonatomic, strong) NSArray *allowedExtensions;
@end

#endif