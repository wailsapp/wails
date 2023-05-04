//go:build darwin
//
//  WailsMenu.m
//  test
//
//  Created by Lea Anthony on 25/10/21.
//

#import <Foundation/Foundation.h>
#import "WailsMenu.h"
#import "WailsMenuItem.h"
#import "Role.h"

@implementation WailsMenu

- (NSMenuItem*) newMenuItem :(NSString*)title :(SEL)selector :(NSString*)key :(NSEventModifierFlags)flags {
    NSMenuItem *result = [[[NSMenuItem alloc] initWithTitle:title action:selector keyEquivalent:key] autorelease];
    [result setKeyEquivalentModifierMask:flags];
    return result;
}

- (NSMenuItem*) newMenuItemWithContext :(WailsContext*)ctx :(NSString*)title :(SEL)selector :(NSString*)key :(NSEventModifierFlags)flags {
    NSMenuItem *result = [NSMenuItem new];
    if ( title != nil ) {
        [result setTitle:title];
    }
    if (selector != nil) {
        [result setAction:selector];
    }
    if (key) {
        [result setKeyEquivalent:key];
    }
    if( flags != 0 ) {
        [result setKeyEquivalentModifierMask:flags];
    }
    result.target = ctx;
    return result;
}

- (NSMenuItem*) newMenuItem :(NSString*)title :(SEL)selector :(NSString*)key  {
    return [self newMenuItem :title :selector :key :0];
}

- (WailsMenu*) initWithNSTitle:(NSString *)title {
    if( title != nil ) {
        [super initWithTitle:title];
    } else {
        [self init];
    }
    [self setAutoenablesItems:NO];
    return self;
}

- (void) appendSubmenu :(WailsMenu*)child {
    NSMenuItem *childMenuItem = [[NSMenuItem new] autorelease];
    [childMenuItem setTitle:child.title];
    [self addItem:childMenuItem];
    [childMenuItem setSubmenu:child];
}

- (void) appendRole :(WailsContext*)ctx :(Role)role {

    switch(role) {
        case AppMenu:
        {
            NSString *appName = [NSRunningApplication currentApplication].localizedName;
            if( appName == nil ) {
                appName = [[NSProcessInfo processInfo] processName];
            }
            WailsMenu *appMenu = [[[WailsMenu new] initWithNSTitle:appName] autorelease];
            
            if (ctx.aboutTitle != nil) {
                [appMenu addItem:[self newMenuItemWithContext :ctx :[@"About " stringByAppendingString:appName] :@selector(About) :nil :0]];
                [appMenu addItem:[NSMenuItem separatorItem]];
            }

            [appMenu addItem:[self newMenuItem:[@"Hide " stringByAppendingString:appName] :@selector(hide:) :@"h" :NSEventModifierFlagCommand]];
            [appMenu addItem:[self newMenuItem:@"Hide Others" :@selector(hideOtherApplications:) :@"h" :(NSEventModifierFlagOption | NSEventModifierFlagCommand)]];
            [appMenu addItem:[self newMenuItem:@"Show All" :@selector(unhideAllApplications:) :@""]];
            [appMenu addItem:[NSMenuItem separatorItem]];

            id quitTitle = [@"Quit " stringByAppendingString:appName];
            NSMenuItem* quitMenuItem = [self newMenuItem:quitTitle :@selector(Quit) :@"q" :NSEventModifierFlagCommand];
            quitMenuItem.target = ctx;
            [appMenu addItem:quitMenuItem];
            [self appendSubmenu:appMenu];
            break;
        }
        case EditMenu:
        {
            WailsMenu *editMenu = [[[WailsMenu new] initWithNSTitle:@"Edit"] autorelease];
            [editMenu addItem:[self newMenuItem:@"Undo" :@selector(undo:) :@"z" :NSEventModifierFlagCommand]];
            [editMenu addItem:[self newMenuItem:@"Redo" :@selector(redo:) :@"z" :(NSEventModifierFlagShift | NSEventModifierFlagCommand)]];
            [editMenu addItem:[NSMenuItem separatorItem]];
            [editMenu addItem:[self newMenuItem:@"Cut" :@selector(cut:) :@"x" :NSEventModifierFlagCommand]];
            [editMenu addItem:[self newMenuItem:@"Copy" :@selector(copy:) :@"c" :NSEventModifierFlagCommand]];
            [editMenu addItem:[self newMenuItem:@"Paste" :@selector(paste:) :@"v" :NSEventModifierFlagCommand]];
            [editMenu addItem:[self newMenuItem:@"Paste and Match Style" :@selector(pasteAsRichText:) :@"v" :(NSEventModifierFlagOption | NSEventModifierFlagShift | NSEventModifierFlagCommand)]];
            [editMenu addItem:[self newMenuItem:@"Delete" :@selector(delete:) :[self accel:@"backspace"] :0]];
            [editMenu addItem:[self newMenuItem:@"Select All" :@selector(selectAll:) :@"a" :NSEventModifierFlagCommand]];
            [editMenu addItem:[NSMenuItem separatorItem]];
//            NSMenuItem *speechMenuItem = [[NSMenuItem new] autorelease];
//            [speechMenuItem setTitle:@"Speech"];
//            [editMenu addItem:speechMenuItem];
            WailsMenu *speechMenu =  [[[WailsMenu new] initWithNSTitle:@"Speech"] autorelease];
            [speechMenu addItem:[self newMenuItem:@"Start Speaking" :@selector(startSpeaking:) :@""]];
            [speechMenu addItem:[self newMenuItem:@"Stop Speaking" :@selector(stopSpeaking:) :@""]];
            [editMenu appendSubmenu:speechMenu];
            [self appendSubmenu:editMenu];
            
            break;
        }
        case WindowMenu:
        {
            WailsMenu *windowMenu = [[[WailsMenu new] initWithNSTitle:@"Window"] autorelease];
            [windowMenu addItem:[self newMenuItem:@"Minimize" :@selector(performMiniaturize:) :@"m" :NSEventModifierFlagCommand]];
            [windowMenu addItem:[self newMenuItem:@"Zoom" :@selector(performZoom:) :@""]];
            [windowMenu addItem:[NSMenuItem separatorItem]];
            [windowMenu addItem:[self newMenuItem:@"Full Screen" :@selector(enterFullScreenMode:) :@"f" :(NSEventModifierFlagControl | NSEventModifierFlagCommand)]];
            [self appendSubmenu:windowMenu];
            
            break;
        }
    }
}

- (void*) AppendMenuItem :(WailsContext*)ctx :(NSString*)label :(NSString *)shortcutKey :(int)modifiers :(bool)disabled :(bool)checked :(int)menuItemID {
    
    NSString *nslabel = @"";
    if (label != nil ) {
        nslabel = label;
    }
    WailsMenuItem *menuItem = [WailsMenuItem new];
    
    // Label
    menuItem.title = nslabel;
    
    // Process callback
    menuItem.menuItemID = menuItemID;
    menuItem.action = @selector(handleClick);
    menuItem.target = menuItem;
    
    // Shortcut
    if (shortcutKey != nil) {
        [menuItem setKeyEquivalent:[self accel:shortcutKey]];
        [menuItem setKeyEquivalentModifierMask:modifiers];
    }
        
    // Enabled/Disabled
    [menuItem setEnabled:!disabled];
    
    // Checked
    [menuItem setState:(checked ? NSControlStateValueOn : NSControlStateValueOff)];  
    
    [self addItem:menuItem];
    return menuItem;
}

- (void) AppendSeparator {
    [self addItem:[NSMenuItem separatorItem]];
}


- (NSString*) accel :(NSString*)key {

    // Guard against no accelerator key
    if( key == NULL ) {
        return @"";
    }

    if( [key isEqualToString:@"backspace"] ) {
        return unicode(0x0008);
    }
    if( [key isEqualToString:@"tab"] ) {
        return unicode(0x0009);
    }
    if( [key isEqualToString:@"return"] ) {
        return unicode(0x000d);
    }
    if( [key isEqualToString:@"enter"] ) {
        return unicode(0x000d);
    }
    if( [key isEqualToString:@"escape"] ) {
        return unicode(0x001b);
    }
    if( [key isEqualToString:@"left"] ) {
        return unicode(0x001c);
    }
    if( [key isEqualToString:@"right"] ) {
        return unicode(0x001d);
    }
    if( [key isEqualToString:@"up"] ) {
        return unicode(0x001e);
    }
    if( [key isEqualToString:@"down"] ) {
        return unicode(0x001f);
    }
    if( [key isEqualToString:@"space"] ) {
        return unicode(0x0020);
    }
    if( [key isEqualToString:@"delete"] ) {
        return unicode(0x007f);
    }
    if( [key isEqualToString:@"home"] ) {
        return unicode(0x2196);
    }
    if( [key isEqualToString:@"end"] ) {
        return unicode(0x2198);
    }
    if( [key isEqualToString:@"page up"] ) {
        return unicode(0x21de);
    }
    if( [key isEqualToString:@"page down"] ) {
        return unicode(0x21df);
    }
    if( [key isEqualToString:@"f1"] ) {
        return unicode(0xf704);
    }
    if( [key isEqualToString:@"f2"] ) {
        return unicode(0xf705);
    }
    if( [key isEqualToString:@"f3"] ) {
        return unicode(0xf706);
    }
    if( [key isEqualToString:@"f4"] ) {
        return unicode(0xf707);
    }
    if( [key isEqualToString:@"f5"] ) {
        return unicode(0xf708);
    }
    if( [key isEqualToString:@"f6"] ) {
        return unicode(0xf709);
    }
    if( [key isEqualToString:@"f7"] ) {
        return unicode(0xf70a);
    }
    if( [key isEqualToString:@"f8"] ) {
        return unicode(0xf70b);
    }
    if( [key isEqualToString:@"f9"] ) {
        return unicode(0xf70c);
    }
    if( [key isEqualToString:@"f10"] ) {
        return unicode(0xf70d);
    }
    if( [key isEqualToString:@"f11"] ) {
        return unicode(0xf70e);
    }
    if( [key isEqualToString:@"f12"] ) {
        return unicode(0xf70f);
    }
    if( [key isEqualToString:@"f13"] ) {
        return unicode(0xf710);
    }
    if( [key isEqualToString:@"f14"] ) {
        return unicode(0xf711);
    }
    if( [key isEqualToString:@"f15"] ) {
        return unicode(0xf712);
    }
    if( [key isEqualToString:@"f16"] ) {
        return unicode(0xf713);
    }
    if( [key isEqualToString:@"f17"] ) {
        return unicode(0xf714);
    }
    if( [key isEqualToString:@"f18"] ) {
        return unicode(0xf715);
    }
    if( [key isEqualToString:@"f19"] ) {
        return unicode(0xf716);
    }
    if( [key isEqualToString:@"f20"] ) {
        return unicode(0xf717);
    }
    if( [key isEqualToString:@"f21"] ) {
        return unicode(0xf718);
    }
    if( [key isEqualToString:@"f22"] ) {
        return unicode(0xf719);
    }
    if( [key isEqualToString:@"f23"] ) {
        return unicode(0xf71a);
    }
    if( [key isEqualToString:@"f24"] ) {
        return unicode(0xf71b);
    }
    if( [key isEqualToString:@"f25"] ) {
        return unicode(0xf71c);
    }
    if( [key isEqualToString:@"f26"] ) {
        return unicode(0xf71d);
    }
    if( [key isEqualToString:@"f27"] ) {
        return unicode(0xf71e);
    }
    if( [key isEqualToString:@"f28"] ) {
        return unicode(0xf71f);
    }
    if( [key isEqualToString:@"f29"] ) {
        return unicode(0xf720);
    }
    if( [key isEqualToString:@"f30"] ) {
        return unicode(0xf721);
    }
    if( [key isEqualToString:@"f31"] ) {
        return unicode(0xf722);
    }
    if( [key isEqualToString:@"f32"] ) {
        return unicode(0xf723);
    }
    if( [key isEqualToString:@"f33"] ) {
        return unicode(0xf724);
    }
    if( [key isEqualToString:@"f34"] ) {
        return unicode(0xf725);
    }
    if( [key isEqualToString:@"f35"] ) {
        return unicode(0xf726);
    }
//  if( [key isEqualToString:@"Insert"] ) {
//    return unicode(0xf727);
//  }
//  if( [key isEqualToString:@"PrintScreen"] ) {
//    return unicode(0xf72e);
//  }
//  if( [key isEqualToString:@"ScrollLock"] ) {
//    return unicode(0xf72f);
//  }
    if( [key isEqualToString:@"numLock"] ) {
        return unicode(0xf739);
    }

    return key;
}


@end


