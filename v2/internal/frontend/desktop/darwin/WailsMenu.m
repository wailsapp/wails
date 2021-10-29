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
    NSMenuItem *result = [[NSMenuItem new] autorelease];
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
    }
    [self setAutoenablesItems:NO];
    return self;
}

- (void) appendSubmenu :(WailsMenu*)child {
    NSMenuItem *childMenuItem = [[NSMenuItem new] autorelease];
    [childMenuItem setTitle:[child title]];
    [self addItem:childMenuItem];
    [childMenuItem setSubmenu:child];
}

- (void) appendRole :(WailsContext*)ctx :(Role)role {

    switch(role) {
        case AppMenu:
        {
            NSString *appName = [[NSProcessInfo processInfo] processName];
            NSString *cap = [appName capitalizedString];
            WailsMenu *appMenu = [[WailsMenu new] initWithNSTitle:cap];
            id quitTitle = [@"Quit " stringByAppendingString:cap];
            NSMenuItem* quitMenuItem = [self newMenuItem:quitTitle :@selector(Quit) :@"q" :NSEventModifierFlagCommand];
            quitMenuItem.target = ctx;
            if (ctx.aboutTitle != nil) {
                [appMenu addItem:[self newMenuItemWithContext :ctx :[@"About " stringByAppendingString:cap] :@selector(About) :nil :0]];
            }
            [appMenu addItem:quitMenuItem];
            [self appendSubmenu:appMenu];
            break;
        }
        case EditMenu:
        {
            WailsMenu *editMenu = [[WailsMenu new] initWithNSTitle:@"Edit"];
            [editMenu addItem:[self newMenuItem:@"Undo" :@selector(undoActionName) :@"z" :NSEventModifierFlagCommand]];
            [editMenu addItem:[self newMenuItem:@"Redo" :@selector(redoActionName) :@"z" :(NSEventModifierFlagShift | NSEventModifierFlagCommand)]];
            [editMenu addItem:[NSMenuItem separatorItem]];
            [editMenu addItem:[self newMenuItem:@"Cut" :@selector(cut:) :@"x" :NSEventModifierFlagCommand]];
            [editMenu addItem:[self newMenuItem:@"Copy" :@selector(copy:) :@"c" :NSEventModifierFlagCommand]];
            [editMenu addItem:[self newMenuItem:@"Paste" :@selector(paste:) :@"v" :NSEventModifierFlagCommand]];
            [editMenu addItem:[self newMenuItem:@"Paste and Match Style" :@selector(pasteAsRichText:) :@"v" :(NSEventModifierFlagOption | NSEventModifierFlagShift | NSEventModifierFlagCommand)]];
            [editMenu addItem:[self newMenuItem:@"Delete" :@selector(delete:) :[self accel:"backspace"] :0]];
            [editMenu addItem:[self newMenuItem:@"Select All" :@selector(selectAll:) :@"a" :NSEventModifierFlagCommand]];
            [editMenu addItem:[NSMenuItem separatorItem]];
//            NSMenuItem *speechMenuItem = [[NSMenuItem new] autorelease];
//            [speechMenuItem setTitle:@"Speech"];
//            [editMenu addItem:speechMenuItem];
            WailsMenu *speechMenu =  [[WailsMenu new] initWithNSTitle:@"Speech"];
            [speechMenu addItem:[self newMenuItem:@"Start Speaking" :@selector(startSpeaking:) :@""]];
            [speechMenu addItem:[self newMenuItem:@"Stop Speaking" :@selector(stopSpeaking:) :@""]];
            [editMenu appendSubmenu:speechMenu];
            [self appendSubmenu:editMenu];
            
            break;
        }
    }
}

- (void*) AppendMenuItem :(WailsContext*)ctx :(const char*)label :(const char *)shortcutKey :(int)modifiers :(bool)disabled :(bool)checked :(int)menuItemID {
    NSString *nslabel = @"";
    if (label != nil ) {
        nslabel = [NSString stringWithUTF8String:label];
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


- (NSString*) accel :(const char *)key {

    // Guard against no accelerator key
    if( key == NULL ) {
        return [NSString stringWithUTF8String:""];
    }

    if( STREQ(key, "backspace") ) {
        return unicode(0x0008);
    }
    if( STREQ(key, "tab") ) {
        return unicode(0x0009);
    }
    if( STREQ(key, "return") ) {
        return unicode(0x000d);
    }
    if( STREQ(key, "enter") ) {
        return unicode(0x000d);
    }
    if( STREQ(key, "escape") ) {
        return unicode(0x001b);
    }
    if( STREQ(key, "left") ) {
        return unicode(0x001c);
    }
    if( STREQ(key, "right") ) {
        return unicode(0x001d);
    }
    if( STREQ(key, "up") ) {
        return unicode(0x001e);
    }
    if( STREQ(key, "down") ) {
        return unicode(0x001f);
    }
    if( STREQ(key, "space") ) {
        return unicode(0x0020);
    }
    if( STREQ(key, "delete") ) {
        return unicode(0x007f);
    }
    if( STREQ(key, "home") ) {
        return unicode(0x2196);
    }
    if( STREQ(key, "end") ) {
        return unicode(0x2198);
    }
    if( STREQ(key, "page up") ) {
        return unicode(0x21de);
    }
    if( STREQ(key, "page down") ) {
        return unicode(0x21df);
    }
    if( STREQ(key, "f1") ) {
        return unicode(0xf704);
    }
    if( STREQ(key, "f2") ) {
        return unicode(0xf705);
    }
    if( STREQ(key, "f3") ) {
        return unicode(0xf706);
    }
    if( STREQ(key, "f4") ) {
        return unicode(0xf707);
    }
    if( STREQ(key, "f5") ) {
        return unicode(0xf708);
    }
    if( STREQ(key, "f6") ) {
        return unicode(0xf709);
    }
    if( STREQ(key, "f7") ) {
        return unicode(0xf70a);
    }
    if( STREQ(key, "f8") ) {
        return unicode(0xf70b);
    }
    if( STREQ(key, "f9") ) {
        return unicode(0xf70c);
    }
    if( STREQ(key, "f10") ) {
        return unicode(0xf70d);
    }
    if( STREQ(key, "f11") ) {
        return unicode(0xf70e);
    }
    if( STREQ(key, "f12") ) {
        return unicode(0xf70f);
    }
    if( STREQ(key, "f13") ) {
        return unicode(0xf710);
    }
    if( STREQ(key, "f14") ) {
        return unicode(0xf711);
    }
    if( STREQ(key, "f15") ) {
        return unicode(0xf712);
    }
    if( STREQ(key, "f16") ) {
        return unicode(0xf713);
    }
    if( STREQ(key, "f17") ) {
        return unicode(0xf714);
    }
    if( STREQ(key, "f18") ) {
        return unicode(0xf715);
    }
    if( STREQ(key, "f19") ) {
        return unicode(0xf716);
    }
    if( STREQ(key, "f20") ) {
        return unicode(0xf717);
    }
    if( STREQ(key, "f21") ) {
        return unicode(0xf718);
    }
    if( STREQ(key, "f22") ) {
        return unicode(0xf719);
    }
    if( STREQ(key, "f23") ) {
        return unicode(0xf71a);
    }
    if( STREQ(key, "f24") ) {
        return unicode(0xf71b);
    }
    if( STREQ(key, "f25") ) {
        return unicode(0xf71c);
    }
    if( STREQ(key, "f26") ) {
        return unicode(0xf71d);
    }
    if( STREQ(key, "f27") ) {
        return unicode(0xf71e);
    }
    if( STREQ(key, "f28") ) {
        return unicode(0xf71f);
    }
    if( STREQ(key, "f29") ) {
        return unicode(0xf720);
    }
    if( STREQ(key, "f30") ) {
        return unicode(0xf721);
    }
    if( STREQ(key, "f31") ) {
        return unicode(0xf722);
    }
    if( STREQ(key, "f32") ) {
        return unicode(0xf723);
    }
    if( STREQ(key, "f33") ) {
        return unicode(0xf724);
    }
    if( STREQ(key, "f34") ) {
        return unicode(0xf725);
    }
    if( STREQ(key, "f35") ) {
        return unicode(0xf726);
    }
//  if( STREQ(key, "Insert") ) {
//    return unicode(0xf727);
//  }
//  if( STREQ(key, "PrintScreen") ) {
//    return unicode(0xf72e);
//  }
//  if( STREQ(key, "ScrollLock") ) {
//    return unicode(0xf72f);
//  }
    if( STREQ(key, "numLock") ) {
        return unicode(0xf739);
    }

    return [NSString stringWithUTF8String:key];
}


@end


