//go:build darwin && !ios && !server

#import "dialogs_mac_extras_darwin.h"
#import <objc/runtime.h>
#include <math.h>

#if __MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
#import <UniformTypeIdentifiers/UTType.h>
#endif

static char const WailsAlertDelegateKey = 0;
static char const WailsSavePanelFormatsKey = 0;

// Splits a joined C string and frees it. Returns an empty array for NULL.
static NSArray<NSString *> *wailsSplitAndFree(char *joined, NSString *separator) {
	if (joined == NULL) {
		return @[];
	}
	NSString *string = [NSString stringWithUTF8String:joined];
	free(joined);
	if (string == nil || [string length] == 0) {
		return @[];
	}
	return [string componentsSeparatedByString:separator];
}

static NSString *wailsStringAndFree(char *value) {
	if (value == NULL) {
		return nil;
	}
	NSString *string = [NSString stringWithUTF8String:value];
	free(value);
	return string;
}

#pragma mark - NSAlert

@interface WailsAlertDelegate : NSObject <NSAlertDelegate>
@property (nonatomic, assign) int helpCallbackID;
@end

@implementation WailsAlertDelegate
- (BOOL)alertShowHelp:(NSAlert *)alert {
	dialogHelpCallback(self.helpCallbackID);
	return YES;
}
@end

void wailsAlertApplyExtras(void *alertPtr, bool showsSuppression, char *suppressionLabel, bool showsHelp, int helpCallbackID, void *accessoryView) {
	NSAlert *alert = (__bridge NSAlert *)alertPtr;
	NSString *label = wailsStringAndFree(suppressionLabel);
	if (showsSuppression) {
		alert.showsSuppressionButton = YES;
		if (label != nil && [label length] > 0) {
			alert.suppressionButton.title = label;
		}
	}
	if (showsHelp) {
		alert.showsHelp = YES;
		WailsAlertDelegate *delegate = [[WailsAlertDelegate alloc] init];
		delegate.helpCallbackID = helpCallbackID;
		alert.delegate = delegate;
		// The delegate property does not retain; tie the delegate's lifetime
		// to the alert.
		objc_setAssociatedObject(alert, &WailsAlertDelegateKey, delegate, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
		[delegate release];
	}
	if (accessoryView != NULL) {
		alert.accessoryView = (__bridge NSView *)accessoryView;
	}
}

bool wailsAlertSuppressionState(void *alertPtr) {
	NSAlert *alert = (__bridge NSAlert *)alertPtr;
	if (alert == nil || !alert.showsSuppressionButton) {
		return false;
	}
	return alert.suppressionButton.state == NSControlStateValueOn;
}

#pragma mark - Content types

NSArray<NSString *> *wailsPanelApplyContentTypes(void *panelPtr, char *contentTypes) {
	NSSavePanel *panel = (__bridge NSSavePanel *)panelPtr;
	NSArray<NSString *> *identifiers = wailsSplitAndFree(contentTypes, @";");
	if ([identifiers count] == 0) {
		return identifiers;
	}
#if __MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
	if (@available(macOS 11, *)) {
		NSMutableArray *types = [NSMutableArray array];
		for (NSString *identifier in identifiers) {
			UTType *type = [UTType typeWithIdentifier:identifier];
			if (type != nil) {
				[types addObject:type];
			}
		}
		if ([types count] > 0) {
			// Keep any extension-derived types installed by AddFilter.
			NSArray *existing = panel.allowedContentTypes;
			if ([existing count] > 0) {
				[types addObjectsFromArray:existing];
			}
			panel.allowedContentTypes = types;
		}
		return identifiers;
	}
#endif
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"
	// Before macOS 11 allowedFileTypes accepts identifiers as well as
	// extensions.
	NSMutableArray *legacy = [NSMutableArray arrayWithArray:identifiers];
	if ([panel.allowedFileTypes count] > 0) {
		[legacy addObjectsFromArray:panel.allowedFileTypes];
	}
	panel.allowedFileTypes = legacy;
#pragma clang diagnostic pop
	return identifiers;
}

#pragma mark - Save panel formats

@interface WailsSavePanelFormats : NSObject
// The panel owns the accessory view, which owns this object.
@property (nonatomic, assign) NSSavePanel *panel;
@property (nonatomic, assign) unsigned int dialogID;
@property (nonatomic, retain) NSArray<NSString *> *extensions;
@property (nonatomic, retain) NSArray<NSString *> *identifiers;
@property (nonatomic, retain) NSPopUpButton *popup;
- (void)formatChanged:(id)sender;
- (void)applyIndex:(NSInteger)index notify:(BOOL)notify;
@end

@implementation WailsSavePanelFormats
- (void)dealloc {
	[_extensions release];
	[_identifiers release];
	[_popup release];
	[super dealloc];
}

- (void)formatChanged:(id)sender {
	[self applyIndex:self.popup.indexOfSelectedItem notify:YES];
}

- (void)applyIndex:(NSInteger)index notify:(BOOL)notify {
	if (index < 0 || index >= (NSInteger)[self.extensions count]) {
		return;
	}
	NSString *extension = self.extensions[index];
	NSString *identifier = self.identifiers[index];
	NSSavePanel *panel = self.panel;

#if __MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
	if (@available(macOS 11, *)) {
		UTType *type = nil;
		if ([identifier length] > 0) {
			type = [UTType typeWithIdentifier:identifier];
		}
		if (type == nil && [extension length] > 0) {
			type = [UTType typeWithFilenameExtension:extension];
		}
		panel.allowedContentTypes = type != nil ? @[type] : @[];
	} else
#endif
	{
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"
		NSString *legacy = [identifier length] > 0 ? identifier : extension;
		panel.allowedFileTypes = [legacy length] > 0 ? @[legacy] : nil;
#pragma clang diagnostic pop
	}

	if ([extension length] > 0) {
		NSString *name = panel.nameFieldStringValue;
		if ([name length] > 0) {
			NSString *stem = [name stringByDeletingPathExtension];
			if ([stem length] == 0) {
				stem = name;
			}
			NSString *renamed = [stem stringByAppendingPathExtension:extension];
			if (renamed != nil) {
				panel.nameFieldStringValue = renamed;
			}
		}
	}

	if (notify) {
		saveFileDialogFormatCallback(self.dialogID, (int)index);
	}
}
@end

static void wailsSavePanelInstallFormats(NSSavePanel *panel, unsigned int dialogID, NSArray<NSString *> *labels, NSArray<NSString *> *extensions, NSArray<NSString *> *identifiers, int selected) {
	NSUInteger count = [labels count];
	if (count == 0 || [extensions count] != count || [identifiers count] != count) {
		return;
	}
	if (selected < 0 || (NSUInteger)selected >= count) {
		selected = 0;
	}

	WailsSavePanelFormats *formats = [[WailsSavePanelFormats alloc] init];
	formats.panel = panel;
	formats.dialogID = dialogID;
	formats.extensions = extensions;
	formats.identifiers = identifiers;

	NSTextField *label = [NSTextField labelWithString:@"Format:"];
	[label sizeToFit];

	NSPopUpButton *popup = [[[NSPopUpButton alloc] initWithFrame:NSMakeRect(0, 0, 200, 26) pullsDown:NO] autorelease];
	for (NSString *title in labels) {
		// addItemWithTitle: drops duplicate titles; add menu items directly.
		NSMenuItem *item = [[NSMenuItem alloc] initWithTitle:title action:NULL keyEquivalent:@""];
		[[popup menu] addItem:item];
		[item release];
	}
	[popup selectItemAtIndex:selected];
	[popup sizeToFit];
	popup.target = formats;
	popup.action = @selector(formatChanged:);
	formats.popup = popup;

	const CGFloat padding = 8;
	const CGFloat height = 34;
	CGFloat labelWidth = NSWidth(label.frame);
	CGFloat labelHeight = NSHeight(label.frame);
	CGFloat popupWidth = MAX(NSWidth(popup.frame), 160);
	CGFloat popupHeight = NSHeight(popup.frame);

	NSView *container = [[[NSView alloc] initWithFrame:NSMakeRect(0, 0, padding + labelWidth + padding + popupWidth + padding, height)] autorelease];
	label.frame = NSMakeRect(padding, floor((height - labelHeight) / 2), labelWidth, labelHeight);
	popup.frame = NSMakeRect(padding + labelWidth + padding, floor((height - popupHeight) / 2), popupWidth, popupHeight);
	[container addSubview:label];
	[container addSubview:popup];

	objc_setAssociatedObject(container, &WailsSavePanelFormatsKey, formats, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
	[formats release];

	panel.accessoryView = container;
	[formats applyIndex:selected notify:NO];
}

void wailsSavePanelApplyExtras(void *panelPtr, unsigned int dialogID, WailsSavePanelExtras *extras) {
	NSSavePanel *panel = (__bridge NSSavePanel *)panelPtr;
	if (extras == NULL) {
		return;
	}
	if (extras->contentTypes != NULL) {
		wailsPanelApplyContentTypes(panelPtr, extras->contentTypes);
	}
	NSString *nameFieldLabel = wailsStringAndFree(extras->nameFieldLabel);
	if (nameFieldLabel != nil) {
		panel.nameFieldLabel = nameFieldLabel;
	}
	NSArray<NSString *> *tags = wailsSplitAndFree(extras->tags, @";");
	if ([tags count] > 0) {
		panel.showsTagField = YES;
		panel.tagNames = tags;
	}
	NSString *separator = @"\x1f";
	NSArray<NSString *> *labels = wailsSplitAndFree(extras->formatLabels, separator);
	NSArray<NSString *> *extensions = wailsSplitAndFree(extras->formatExtensions, separator);
	NSArray<NSString *> *identifiers = wailsSplitAndFree(extras->formatIdentifiers, separator);
	if ([labels count] > 0) {
		wailsSavePanelInstallFormats(panel, dialogID, labels, extensions, identifiers, extras->selectedFormat);
	}
}

#pragma mark - Text prompt

static void wailsFinishPrompt(int callbackID, NSAlert *alert, NSTextField *field, NSModalResponse response) {
	bool ok = response == NSAlertFirstButtonReturn;
	const char *value = NULL;
	if (ok) {
		value = [[field stringValue] UTF8String];
	}
	dialogPromptCallback(callbackID, (char *)value, ok);
	[alert release];
}

void wailsShowPromptDialog(int callbackID, char *title, char *message, char *placeholder, char *defaultValue, bool secure, char *okLabel, char *cancelLabel, void *window) {
	NSAlert *alert = [[NSAlert alloc] init];
	alert.alertStyle = NSAlertStyleInformational;

	NSString *titleString = wailsStringAndFree(title);
	if (titleString != nil) {
		alert.messageText = titleString;
	}
	NSString *messageString = wailsStringAndFree(message);
	if (messageString != nil) {
		alert.informativeText = messageString;
	}

	NSRect fieldFrame = NSMakeRect(0, 0, 300, 24);
	NSTextField *field = secure ? [[NSSecureTextField alloc] initWithFrame:fieldFrame] : [[NSTextField alloc] initWithFrame:fieldFrame];
	NSString *placeholderString = wailsStringAndFree(placeholder);
	if (placeholderString != nil) {
		field.placeholderString = placeholderString;
	}
	NSString *defaultString = wailsStringAndFree(defaultValue);
	if (defaultString != nil) {
		field.stringValue = defaultString;
	}
	alert.accessoryView = field;
	// The alert retains its accessory view.
	[field release];

	NSString *okString = wailsStringAndFree(okLabel);
	NSButton *ok = [alert addButtonWithTitle:okString != nil ? okString : @"OK"];
	ok.keyEquivalent = @"\r";
	NSString *cancelString = wailsStringAndFree(cancelLabel);
	NSButton *cancel = [alert addButtonWithTitle:cancelString != nil ? cancelString : @"Cancel"];
	cancel.keyEquivalent = @"\033";

	// Focus the field so typing starts immediately.
	[[alert window] setInitialFirstResponder:field];

	if (window == NULL) {
		NSModalResponse response = [alert runModal];
		wailsFinishPrompt(callbackID, alert, field, response);
		return;
	}
	[alert beginSheetModalForWindow:(__bridge NSWindow *)window completionHandler:^(NSModalResponse response) {
		wailsFinishPrompt(callbackID, alert, field, response);
	}];
}

#pragma mark - Colour panel

static void wailsColorComponents(NSColor *color, int *red, int *green, int *blue, int *alpha) {
	*red = 0;
	*green = 0;
	*blue = 0;
	*alpha = 255;
	if (color == nil) {
		return;
	}
	NSColor *rgb = [color colorUsingColorSpace:[NSColorSpace sRGBColorSpace]];
	if (rgb == nil) {
		rgb = [color colorUsingColorSpace:[NSColorSpace deviceRGBColorSpace]];
	}
	if (rgb == nil) {
		return;
	}
	CGFloat r = 0, g = 0, b = 0, a = 1;
	[rgb getRed:&r green:&g blue:&b alpha:&a];
	*red = (int)lround(MIN(MAX(r, 0), 1) * 255.0);
	*green = (int)lround(MIN(MAX(g, 0), 1) * 255.0);
	*blue = (int)lround(MIN(MAX(b, 0), 1) * 255.0);
	*alpha = (int)lround(MIN(MAX(a, 0), 1) * 255.0);
}

@interface WailsColorPanelSession : NSObject
@property (nonatomic, assign) int callbackID;
@property (nonatomic, assign) BOOL changed;
@property (nonatomic, assign) BOOL live;
- (void)colorChanged:(id)sender;
- (void)panelWillClose:(NSNotification *)notification;
@end

static WailsColorPanelSession *wailsActiveColorSession = nil;

@implementation WailsColorPanelSession
- (void)colorChanged:(id)sender {
	self.changed = YES;
	if (self.live) {
		int red, green, blue, alpha;
		wailsColorComponents([[NSColorPanel sharedColorPanel] color], &red, &green, &blue, &alpha);
		dialogColorCallback(self.callbackID, red, green, blue, alpha, false, true);
	}
}

- (void)panelWillClose:(NSNotification *)notification {
	NSColorPanel *panel = [NSColorPanel sharedColorPanel];
	[[NSNotificationCenter defaultCenter] removeObserver:self];
	// NSColorPanel exposes no target getter; only one session is active
	// at a time, so the session owns the target.
	[panel setTarget:nil];
	[panel setAction:NULL];
	int red, green, blue, alpha;
	wailsColorComponents([panel color], &red, &green, &blue, &alpha);
	if (wailsActiveColorSession == self) {
		wailsActiveColorSession = nil;
	}
	dialogColorCallback(self.callbackID, red, green, blue, alpha, true, self.changed);
	[self autorelease];
}
@end

bool wailsShowColorPanel(int callbackID, int red, int green, int blue, int alpha, bool showsAlpha, char *title, bool live) {
	NSString *titleString = wailsStringAndFree(title);
	if (wailsActiveColorSession != nil) {
		return false;
	}
	WailsColorPanelSession *session = [[WailsColorPanelSession alloc] init];
	session.callbackID = callbackID;
	session.live = live;
	wailsActiveColorSession = session;

	NSColorPanel *panel = [NSColorPanel sharedColorPanel];
	panel.showsAlpha = showsAlpha;
	panel.continuous = YES;
	// Set the colour before wiring the action so the initial value is not
	// reported as a change.
	panel.color = [NSColor colorWithSRGBRed:red / 255.0 green:green / 255.0 blue:blue / 255.0 alpha:alpha / 255.0];
	if (titleString != nil) {
		panel.title = titleString;
	}
	panel.target = session;
	panel.action = @selector(colorChanged:);
	[[NSNotificationCenter defaultCenter] addObserver:session selector:@selector(panelWillClose:) name:NSWindowWillCloseNotification object:panel];
	[panel makeKeyAndOrderFront:nil];
	return true;
}

#pragma mark - Font panel

static void wailsReportFont(int callbackID, NSFont *font, bool closed, bool changed) {
	NSString *family = font.familyName;
	NSString *face = [font.fontDescriptor objectForKey:NSFontFaceAttribute];
	NSString *postScriptName = font.fontName;
	const char *familyC = family != nil ? [family UTF8String] : "";
	const char *faceC = face != nil ? [face UTF8String] : "";
	const char *postScriptC = postScriptName != nil ? [postScriptName UTF8String] : "";
	dialogFontCallback(callbackID, (char *)familyC, (char *)faceC, (char *)postScriptC, (double)font.pointSize, closed, changed);
}

@interface WailsFontPanelSession : NSObject
@property (nonatomic, assign) int callbackID;
@property (nonatomic, retain) NSFont *font;
@property (nonatomic, assign) BOOL changed;
@property (nonatomic, assign) BOOL live;
- (void)changeFont:(id)sender;
- (void)panelWillClose:(NSNotification *)notification;
@end

static WailsFontPanelSession *wailsActiveFontSession = nil;

@implementation WailsFontPanelSession
- (void)dealloc {
	[_font release];
	[super dealloc];
}

- (void)changeFont:(id)sender {
	NSFontManager *manager = [sender isKindOfClass:[NSFontManager class]] ? sender : [NSFontManager sharedFontManager];
	NSFont *converted = [manager convertFont:self.font];
	if (converted != nil) {
		self.font = converted;
	}
	self.changed = YES;
	if (self.live) {
		wailsReportFont(self.callbackID, self.font, false, true);
	}
}

- (void)panelWillClose:(NSNotification *)notification {
	[[NSNotificationCenter defaultCenter] removeObserver:self];
	NSFontManager *manager = [NSFontManager sharedFontManager];
	if (manager.target == self) {
		manager.target = nil;
	}
	if (wailsActiveFontSession == self) {
		wailsActiveFontSession = nil;
	}
	wailsReportFont(self.callbackID, self.font, true, self.changed);
	[self autorelease];
}
@end

bool wailsShowFontPanel(int callbackID, char *family, double size, bool live) {
	NSString *familyString = wailsStringAndFree(family);
	if (wailsActiveFontSession != nil) {
		return false;
	}
	if (size <= 0) {
		size = [NSFont systemFontSize];
	}
	NSFontManager *manager = [NSFontManager sharedFontManager];
	NSFont *font = nil;
	if (familyString != nil && [familyString length] > 0) {
		font = [NSFont fontWithName:familyString size:size];
		if (font == nil) {
			font = [manager fontWithFamily:familyString traits:0 weight:5 size:size];
		}
	}
	if (font == nil) {
		font = [NSFont systemFontOfSize:size];
	}

	WailsFontPanelSession *session = [[WailsFontPanelSession alloc] init];
	session.callbackID = callbackID;
	session.live = live;
	session.font = font;
	wailsActiveFontSession = session;

	manager.target = session;
	[manager setSelectedFont:font isMultiple:NO];
	NSFontPanel *panel = [NSFontPanel sharedFontPanel];
	[[NSNotificationCenter defaultCenter] addObserver:session selector:@selector(panelWillClose:) name:NSWindowWillCloseNotification object:panel];
	[panel makeKeyAndOrderFront:nil];
	return true;
}
