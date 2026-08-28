//go:build darwin && !ios && !server

#import "webview_window_split_darwin.h"
#import <objc/runtime.h>
#import <string.h>

#ifndef WAILS_NATIVE_ONLY
static const void* WailsSplitPrimaryPaneIDAssociationKey = &WailsSplitPrimaryPaneIDAssociationKey;
#endif
static void* WailsSplitPaneCollapsedKVOContext = &WailsSplitPaneCollapsedKVOContext;

#ifndef WAILS_NATIVE_ONLY
unsigned long long splitPrimaryPaneIDForWebView(WKWebView* webView) {
    if (webView == nil) return 0;
    NSNumber* value = objc_getAssociatedObject(webView, WailsSplitPrimaryPaneIDAssociationKey);
    return value == nil ? 0 : value.unsignedLongLongValue;
}
#else
unsigned long long splitPrimaryPaneIDForWebView(void* webView) {
    return 0;
}
#endif

@interface WailsSidebarNode : NSObject
@property unsigned long long nodeID;
@property BOOL section;
@property (copy) NSString* label;
@property (copy) NSString* subtitle;
@property (copy) NSString* symbolName;
@property (retain) NSData* imageData;
@property (copy) NSString* tooltip;
@property BOOL disabled;
@property BOOL hidden;
@property (retain) NSMutableArray<WailsSidebarNode*>* children;
@end

@implementation WailsSidebarNode
- (void)dealloc {
    [_label release];
    [_subtitle release];
    [_symbolName release];
    [_imageData release];
    [_tooltip release];
    [_children release];
    [super dealloc];
}
@end

@interface WailsSidebarDetailCell : NSTableCellView
@property (retain) NSTextField* subtitleTextField;
@end

@implementation WailsSidebarDetailCell
- (void)dealloc {
    [_subtitleTextField release];
    [super dealloc];
}
@end

// A fixed account/status item outside the outline view. It uses the same
// detailed presentation as a sidebar row, without scrolling or selection.
@interface WailsSidebarFooterView : NSView
@property (nonatomic, retain) WailsSidebarNode* node;
@property (retain) NSImageView* imageView;
@property (retain) NSTextField* titleTextField;
@property (retain) NSTextField* subtitleTextField;
@end

@implementation WailsSidebarFooterView
- (instancetype)initWithFrame:(NSRect)frame {
    self = [super initWithFrame:frame];
    if (self == nil) return nil;
    NSImageView* image = [[NSImageView alloc] initWithFrame:NSZeroRect];
    image.translatesAutoresizingMaskIntoConstraints = NO;
    image.imageScaling = NSImageScaleProportionallyUpOrDown;
    image.wantsLayer = YES;
    image.layer.cornerRadius = 20;
    image.layer.masksToBounds = YES;
    NSTextField* title = [NSTextField labelWithString:@""];
    title.translatesAutoresizingMaskIntoConstraints = NO;
    title.font = [NSFont systemFontOfSize:13 weight:NSFontWeightSemibold];
    title.lineBreakMode = NSLineBreakByTruncatingTail;
    NSTextField* subtitle = [NSTextField labelWithString:@""];
    subtitle.translatesAutoresizingMaskIntoConstraints = NO;
    subtitle.font = [NSFont systemFontOfSize:12 weight:NSFontWeightRegular];
    subtitle.lineBreakMode = NSLineBreakByTruncatingTail;
    subtitle.textColor = [NSColor secondaryLabelColor];
    [self addSubview:image];
    [self addSubview:title];
    [self addSubview:subtitle];
    [NSLayoutConstraint activateConstraints:@[
        [image.leadingAnchor constraintEqualToAnchor:self.leadingAnchor constant:12],
        [image.centerYAnchor constraintEqualToAnchor:self.centerYAnchor],
        [image.widthAnchor constraintEqualToConstant:40],
        [image.heightAnchor constraintEqualToConstant:40],
        [title.leadingAnchor constraintEqualToAnchor:image.trailingAnchor constant:10],
        [title.trailingAnchor constraintEqualToAnchor:self.trailingAnchor constant:-12],
        [title.centerYAnchor constraintEqualToAnchor:self.centerYAnchor constant:-9],
        [subtitle.leadingAnchor constraintEqualToAnchor:title.leadingAnchor],
        [subtitle.trailingAnchor constraintEqualToAnchor:title.trailingAnchor],
        [subtitle.centerYAnchor constraintEqualToAnchor:self.centerYAnchor constant:10]
    ]];
    self.imageView = image;
    self.titleTextField = title;
    self.subtitleTextField = subtitle;
    [image release];
    return self;
}

- (void)setNode:(WailsSidebarNode*)node {
    if (_node == node) return;
    [_node release];
    _node = [node retain];
    self.hidden = node == nil || node.hidden;
    self.toolTip = node.tooltip.length > 0 ? node.tooltip : nil;
    self.titleTextField.stringValue = node.label ?: @"";
    self.subtitleTextField.stringValue = node.subtitle ?: @"";
    BOOL disabled = node.disabled;
    self.titleTextField.textColor = disabled ? [NSColor disabledControlTextColor] : [NSColor labelColor];
    self.subtitleTextField.textColor = disabled ? [NSColor disabledControlTextColor] : [NSColor secondaryLabelColor];
    NSImage* image = node.imageData.length > 0 ? [[[NSImage alloc] initWithData:node.imageData] autorelease] : nil;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
    if (image == nil) {
        if (@available(macOS 11.0, *)) {
            image = node.symbolName.length > 0
                ? [NSImage imageWithSystemSymbolName:node.symbolName accessibilityDescription:node.label]
                : nil;
        }
    }
#endif
    self.imageView.image = image;
    self.imageView.hidden = image == nil;
}

- (void)mouseUp:(NSEvent*)event {
    if (self.node != nil && !self.node.disabled && NSPointInRect([self convertPoint:event.locationInWindow fromView:nil], self.bounds)) {
        processMacSidebarItemSelected(self.node.nodeID);
    }
}

- (void)dealloc {
    [_node release];
    [_imageView release];
    [_titleTextField release];
    [_subtitleTextField release];
    [super dealloc];
}
@end

@interface WailsPrimaryPaneView : NSView
@property (retain) NSColor* fillColor;
@end

@implementation WailsPrimaryPaneView
- (void)dealloc {
    [_fillColor release];
    [super dealloc];
}

- (BOOL)isOpaque {
    return self.fillColor != nil && self.fillColor.alphaComponent >= 1.0;
}

- (void)drawRect:(NSRect)dirtyRect {
    if (self.fillColor == nil) return;
    [self.fillColor setFill];
    NSRectFill(dirtyRect);
}

- (void)viewDidChangeEffectiveAppearance {
    [super viewDidChangeEffectiveAppearance];
    [self setNeedsDisplay:YES];
}
@end

// A real AppKit source list. Both the outline and its scroll view remain
// transparent so NSSplitViewItem's semantic sidebar material is visible.
@interface WailsSidebarViewController : NSViewController <NSOutlineViewDataSource, NSOutlineViewDelegate>
@property (retain) NSMutableArray<WailsSidebarNode*>* roots;
@property (retain) NSOutlineView* outlineView;
@property (retain) NSScrollView* scrollView;
@property (retain) WailsSidebarNode* footer;
@property (retain) WailsSidebarFooterView* footerView;
@property (retain) NSLayoutConstraint* footerHeightConstraint;
@property unsigned long long selectedItemID;
@property BOOL suppressSelectionCallback;
- (void)fitOutlineToViewport;
- (void)reloadContents;
@end

@implementation WailsSidebarViewController

- (void)dealloc {
    [_roots release];
    [_outlineView release];
    [_scrollView release];
    [_footer release];
    [_footerView release];
    [_footerHeightConstraint release];
    [super dealloc];
}

- (void)loadView {
    NSView* container = [[NSView alloc] initWithFrame:NSMakeRect(0, 0, 240, 600)];
    NSScrollView* scrollView = [[NSScrollView alloc] initWithFrame:NSMakeRect(0, 0, 240, 600)];
    scrollView.translatesAutoresizingMaskIntoConstraints = NO;
    // The semantic NSSplitViewItem supplies the native sidebar material.
    // Keeping the scrolling rows transparent makes them resolve against the
    // same surface as the fixed footer instead of covering it with the
    // primary pane's solid background colour.
    scrollView.drawsBackground = NO;
    scrollView.borderType = NSNoBorder;
    scrollView.hasVerticalScroller = YES;
    scrollView.hasHorizontalScroller = NO;
    scrollView.autohidesScrollers = YES;
    scrollView.horizontalScrollElasticity = NSScrollElasticityNone;

    NSOutlineView* outlineView = [[NSOutlineView alloc] initWithFrame:scrollView.bounds];
    outlineView.backgroundColor = [NSColor clearColor];
    outlineView.headerView = nil;
    outlineView.floatsGroupRows = YES;
    outlineView.indentationPerLevel = 12.0;
    outlineView.autoresizingMask = NSViewWidthSizable;
    outlineView.autoresizesOutlineColumn = YES;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
    if (@available(macOS 11.0, *)) {
        // The table style, rather than the deprecated selection-only style,
        // preserves AppKit's source-list material and subdued selection.
        outlineView.style = NSTableViewStyleSourceList;
    } else {
        outlineView.selectionHighlightStyle = NSTableViewSelectionHighlightStyleSourceList;
    }
#else
    outlineView.selectionHighlightStyle = NSTableViewSelectionHighlightStyleSourceList;
#endif
    // Finder keeps focus in the primary content when its source-list rows are
    // clicked. The selected row therefore uses AppKit's subdued,
    // non-emphasized treatment instead of the bright accent fill.
    outlineView.refusesFirstResponder = YES;
    outlineView.allowsEmptySelection = YES;
    outlineView.allowsMultipleSelection = NO;
    outlineView.rowSizeStyle = NSTableViewRowSizeStyleDefault;

    NSTableColumn* column = [[NSTableColumn alloc] initWithIdentifier:@"sidebar"];
    column.resizingMask = NSTableColumnAutoresizingMask;
    column.width = 240;
    outlineView.columnAutoresizingStyle = NSTableViewLastColumnOnlyAutoresizingStyle;
    [outlineView addTableColumn:column];
    outlineView.outlineTableColumn = column;
    [column release];

    outlineView.dataSource = self;
    outlineView.delegate = self;
    scrollView.documentView = outlineView;
    WailsSidebarFooterView* footer = [[WailsSidebarFooterView alloc] initWithFrame:NSMakeRect(0, 0, 240, 64)];
    footer.translatesAutoresizingMaskIntoConstraints = NO;
    [container addSubview:scrollView];
    [container addSubview:footer];
    NSLayoutConstraint* footerHeight = [footer.heightAnchor constraintEqualToConstant:0];
    [NSLayoutConstraint activateConstraints:@[
        [scrollView.leadingAnchor constraintEqualToAnchor:container.leadingAnchor],
        [scrollView.trailingAnchor constraintEqualToAnchor:container.trailingAnchor],
        [scrollView.topAnchor constraintEqualToAnchor:container.topAnchor],
        [scrollView.bottomAnchor constraintEqualToAnchor:footer.topAnchor],
        [footer.leadingAnchor constraintEqualToAnchor:container.leadingAnchor],
        [footer.trailingAnchor constraintEqualToAnchor:container.trailingAnchor],
        [footer.bottomAnchor constraintEqualToAnchor:container.bottomAnchor],
        footerHeight
    ]];
    self.outlineView = outlineView;
    self.scrollView = scrollView;
    self.footerView = footer;
    self.footerHeightConstraint = footerHeight;
    self.view = container;
    [container release];
    [outlineView release];
    [scrollView release];
    [footer release];
}

- (void)viewDidLayout {
    [super viewDidLayout];
    [self fitOutlineToViewport];
}

- (void)fitOutlineToViewport {
    if (self.outlineView == nil || self.scrollView == nil) return;
    CGFloat width = self.scrollView.contentView.bounds.size.width;
    if (width <= 0) return;

    NSRect frame = self.outlineView.frame;
    frame.origin.x = 0;
    frame.size.width = width;
    self.outlineView.frame = frame;
    self.outlineView.outlineTableColumn.width = width;

    NSPoint origin = self.scrollView.contentView.bounds.origin;
    if (origin.x != 0) {
        origin.x = 0;
        [self.scrollView.contentView scrollToPoint:origin];
        [self.scrollView reflectScrolledClipView:self.scrollView.contentView];
    }
}

- (NSArray<WailsSidebarNode*>*)visibleNodes:(NSArray<WailsSidebarNode*>*)nodes {
    NSMutableArray* visible = [NSMutableArray arrayWithCapacity:nodes.count];
    for (WailsSidebarNode* node in nodes) {
        if (!node.hidden) [visible addObject:node];
    }
    return visible;
}

- (NSInteger)outlineView:(NSOutlineView*)outlineView numberOfChildrenOfItem:(id)item {
    NSArray* nodes = item == nil ? self.roots : ((WailsSidebarNode*)item).children;
    return [self visibleNodes:nodes].count;
}

- (id)outlineView:(NSOutlineView*)outlineView child:(NSInteger)index ofItem:(id)item {
    NSArray* nodes = item == nil ? self.roots : ((WailsSidebarNode*)item).children;
    return [self visibleNodes:nodes][index];
}

- (BOOL)outlineView:(NSOutlineView*)outlineView isItemExpandable:(id)item {
    return ((WailsSidebarNode*)item).section;
}

- (BOOL)outlineView:(NSOutlineView*)outlineView isGroupItem:(id)item {
    return ((WailsSidebarNode*)item).section;
}

- (NSTableCellView*)newCellWithIdentifier:(NSUserInterfaceItemIdentifier)identifier section:(BOOL)section {
    NSTableCellView* cell = [[[NSTableCellView alloc] initWithFrame:NSMakeRect(0, 0, 220, section ? 22 : 28)] autorelease];
    cell.identifier = identifier;

    NSTextField* textField = [NSTextField labelWithString:@""];
    textField.translatesAutoresizingMaskIntoConstraints = NO;
    textField.lineBreakMode = NSLineBreakByTruncatingTail;
    if (section) {
        textField.font = [NSFont systemFontOfSize:11 weight:NSFontWeightSemibold];
        textField.textColor = [NSColor secondaryLabelColor];
        [cell addSubview:textField];
        [NSLayoutConstraint activateConstraints:@[
            [textField.leadingAnchor constraintEqualToAnchor:cell.leadingAnchor constant:4],
            [textField.trailingAnchor constraintEqualToAnchor:cell.trailingAnchor constant:-4],
            [textField.centerYAnchor constraintEqualToAnchor:cell.centerYAnchor]
        ]];
    } else {
        NSImageView* imageView = [[[NSImageView alloc] initWithFrame:NSZeroRect] autorelease];
        imageView.translatesAutoresizingMaskIntoConstraints = NO;
        imageView.imageScaling = NSImageScaleProportionallyUpOrDown;
        [cell addSubview:imageView];
        [cell addSubview:textField];
        [NSLayoutConstraint activateConstraints:@[
            [imageView.leadingAnchor constraintEqualToAnchor:cell.leadingAnchor constant:2],
            [imageView.centerYAnchor constraintEqualToAnchor:cell.centerYAnchor],
            [imageView.widthAnchor constraintEqualToConstant:16],
            [imageView.heightAnchor constraintEqualToConstant:16],
            [textField.leadingAnchor constraintEqualToAnchor:imageView.trailingAnchor constant:7],
            [textField.trailingAnchor constraintEqualToAnchor:cell.trailingAnchor constant:-4],
            [textField.centerYAnchor constraintEqualToAnchor:cell.centerYAnchor]
        ]];
        cell.imageView = imageView;
    }
    cell.textField = textField;
    return cell;
}

- (WailsSidebarDetailCell*)newDetailCellWithIdentifier:(NSUserInterfaceItemIdentifier)identifier {
    WailsSidebarDetailCell* cell = [[[WailsSidebarDetailCell alloc] initWithFrame:NSMakeRect(0, 0, 220, 56)] autorelease];
    cell.identifier = identifier;

    NSImageView* imageView = [[[NSImageView alloc] initWithFrame:NSZeroRect] autorelease];
    imageView.translatesAutoresizingMaskIntoConstraints = NO;
    imageView.imageScaling = NSImageScaleProportionallyUpOrDown;
    imageView.wantsLayer = YES;
    imageView.layer.cornerRadius = 20;
    imageView.layer.masksToBounds = YES;

    NSTextField* title = [NSTextField labelWithString:@""];
    title.translatesAutoresizingMaskIntoConstraints = NO;
    title.font = [NSFont systemFontOfSize:13 weight:NSFontWeightSemibold];
    title.lineBreakMode = NSLineBreakByTruncatingTail;

    NSTextField* subtitle = [NSTextField labelWithString:@""];
    subtitle.translatesAutoresizingMaskIntoConstraints = NO;
    subtitle.font = [NSFont systemFontOfSize:12 weight:NSFontWeightRegular];
    subtitle.lineBreakMode = NSLineBreakByTruncatingTail;
    subtitle.textColor = [NSColor secondaryLabelColor];

    [cell addSubview:imageView];
    [cell addSubview:title];
    [cell addSubview:subtitle];
    [NSLayoutConstraint activateConstraints:@[
        [imageView.leadingAnchor constraintEqualToAnchor:cell.leadingAnchor constant:6],
        [imageView.centerYAnchor constraintEqualToAnchor:cell.centerYAnchor],
        [imageView.widthAnchor constraintEqualToConstant:40],
        [imageView.heightAnchor constraintEqualToConstant:40],
        [title.leadingAnchor constraintEqualToAnchor:imageView.trailingAnchor constant:10],
        [title.trailingAnchor constraintEqualToAnchor:cell.trailingAnchor constant:-6],
        [title.centerYAnchor constraintEqualToAnchor:cell.centerYAnchor constant:-9],
        [subtitle.leadingAnchor constraintEqualToAnchor:title.leadingAnchor],
        [subtitle.trailingAnchor constraintEqualToAnchor:title.trailingAnchor],
        [subtitle.centerYAnchor constraintEqualToAnchor:cell.centerYAnchor constant:10]
    ]];
    cell.imageView = imageView;
    cell.textField = title;
    cell.subtitleTextField = subtitle;
    return cell;
}

- (CGFloat)outlineView:(NSOutlineView*)outlineView heightOfRowByItem:(id)item {
    WailsSidebarNode* node = (WailsSidebarNode*)item;
    if (node.section) return 22;
    return node.subtitle.length > 0 || node.imageData.length > 0 ? 56 : 28;
}

- (NSView*)outlineView:(NSOutlineView*)outlineView viewForTableColumn:(NSTableColumn*)tableColumn item:(id)item {
    WailsSidebarNode* node = (WailsSidebarNode*)item;
    BOOL detail = !node.section && (node.subtitle.length > 0 || node.imageData.length > 0);
    NSUserInterfaceItemIdentifier identifier = node.section ? @"WailsSidebarSection" :
        (detail ? @"WailsSidebarDetailItem" : @"WailsSidebarItem");
    NSTableCellView* cell = [outlineView makeViewWithIdentifier:identifier owner:self];
    if (cell == nil) cell = detail ? [self newDetailCellWithIdentifier:identifier] :
        [self newCellWithIdentifier:identifier section:node.section];
    cell.textField.stringValue = node.label ?: @"";
    cell.toolTip = node.tooltip.length > 0 ? node.tooltip : nil;
    cell.textField.textColor = node.disabled ? [NSColor disabledControlTextColor] :
        (node.section ? [NSColor secondaryLabelColor] : [NSColor labelColor]);
    if (detail) {
        WailsSidebarDetailCell* detailCell = (WailsSidebarDetailCell*)cell;
        detailCell.subtitleTextField.stringValue = node.subtitle ?: @"";
        detailCell.subtitleTextField.textColor = node.disabled ? [NSColor disabledControlTextColor] : [NSColor secondaryLabelColor];
        NSImage* image = node.imageData.length > 0 ? [[[NSImage alloc] initWithData:node.imageData] autorelease] : nil;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
        if (image == nil) {
            if (@available(macOS 11.0, *)) {
                image = node.symbolName.length > 0
                    ? [NSImage imageWithSystemSymbolName:node.symbolName accessibilityDescription:node.label]
                    : nil;
            }
        }
#endif
        detailCell.imageView.image = image;
        detailCell.imageView.hidden = image == nil;
    } else if (!node.section) {
        NSImage* image = nil;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
        if (@available(macOS 11.0, *)) {
            image = node.symbolName.length > 0
                ? [NSImage imageWithSystemSymbolName:node.symbolName accessibilityDescription:node.label]
                : nil;
        }
#endif
        cell.imageView.image = image;
        cell.imageView.hidden = image == nil;
    }
    return cell;
}

- (BOOL)outlineView:(NSOutlineView*)outlineView shouldSelectItem:(id)item {
    WailsSidebarNode* node = (WailsSidebarNode*)item;
    return !node.section && !node.disabled;
}

- (void)outlineViewSelectionDidChange:(NSNotification*)notification {
    if (self.suppressSelectionCallback) return;
    NSInteger row = self.outlineView.selectedRow;
    if (row < 0) {
        self.selectedItemID = 0;
        return;
    }
    WailsSidebarNode* node = [self.outlineView itemAtRow:row];
    if (node == nil || node.section || node.disabled) return;
    self.selectedItemID = node.nodeID;
    processMacSidebarItemSelected(node.nodeID);
}

- (void)reloadContents {
    if (self.outlineView == nil) return;
    self.suppressSelectionCallback = YES;
    [self fitOutlineToViewport];
    [self.outlineView reloadData];
    for (WailsSidebarNode* node in self.roots) {
        if (node.section && !node.hidden) [self.outlineView expandItem:node];
    }
    NSIndexSet* selection = [NSIndexSet indexSet];
    if (self.selectedItemID != 0) {
        for (NSInteger row = 0; row < self.outlineView.numberOfRows; row++) {
            WailsSidebarNode* node = [self.outlineView itemAtRow:row];
            if (!node.section && node.nodeID == self.selectedItemID) {
                selection = [NSIndexSet indexSetWithIndex:row];
                break;
            }
        }
    }
    [self.outlineView selectRowIndexes:selection byExtendingSelection:NO];
    self.suppressSelectionCallback = NO;
    self.footerView.node = self.footer;
    self.footerHeightConstraint.constant = self.footerView.hidden ? 0 : 64;
}

@end

static const void* WailsInspectorControlIDAssociationKey = &WailsInspectorControlIDAssociationKey;

@interface WailsInspectorControlModel : NSObject
@property unsigned long long controlID;
@property int kind;
@property (copy) NSString* label;
@property (copy) NSString* value;
@property BOOL checked;
@property (retain) NSArray<NSString*>* options;
@property NSInteger selectedIndex;
@property (copy) NSString* tooltip;
@property BOOL disabled;
@property BOOL hidden;
@end

@implementation WailsInspectorControlModel
- (void)dealloc {
    [_label release];
    [_value release];
    [_options release];
    [_tooltip release];
    [super dealloc];
}
@end

@interface WailsInspectorSectionModel : NSObject
@property unsigned long long sectionID;
@property (copy) NSString* label;
@property (retain) NSMutableArray<WailsInspectorControlModel*>* controls;
@end

@implementation WailsInspectorSectionModel
- (void)dealloc {
    [_label release];
    [_controls release];
    [super dealloc];
}
@end

@interface WailsInspectorDocumentView : NSView
@end

@implementation WailsInspectorDocumentView
- (BOOL)isFlipped { return YES; }
@end

// A native property inspector. Its scroll view, section headings, labels,
// text fields, checkboxes, and pop-up buttons are all AppKit controls. The
// view remains transparent so NSSplitViewItem's semantic inspector surface
// controls the appearance.
@interface WailsInspectorViewController : NSViewController <NSTextFieldDelegate>
@property (retain) NSMutableArray<WailsInspectorSectionModel*>* sections;
@property (retain) NSMutableDictionary<NSNumber*, WailsInspectorControlModel*>* modelsByID;
@property (retain) NSMutableDictionary<NSNumber*, NSControl*>* controlsByID;
@property (retain) NSMutableDictionary<NSNumber*, NSView*>* rowsByID;
@property (retain) NSMutableDictionary<NSNumber*, NSTextField*>* nameLabelsByID;
@property (retain) NSStackView* stackView;
@property (retain) NSColor* surfaceColor;
- (void)reloadContents;
- (void)applyModel:(WailsInspectorControlModel*)model;
@end

@implementation WailsInspectorViewController

- (instancetype)init {
    self = [super init];
    if (self) {
        _sections = [[NSMutableArray alloc] init];
        _modelsByID = [[NSMutableDictionary alloc] init];
        _controlsByID = [[NSMutableDictionary alloc] init];
        _rowsByID = [[NSMutableDictionary alloc] init];
        _nameLabelsByID = [[NSMutableDictionary alloc] init];
    }
    return self;
}

- (void)dealloc {
    [_sections release];
    [_modelsByID release];
    [_controlsByID release];
    [_rowsByID release];
    [_nameLabelsByID release];
    [_stackView release];
    [_surfaceColor release];
    [super dealloc];
}

- (void)loadView {
    NSScrollView* scrollView = [[NSScrollView alloc] initWithFrame:NSMakeRect(0, 0, 280, 600)];
    scrollView.drawsBackground = self.surfaceColor != nil;
    if (self.surfaceColor != nil) scrollView.backgroundColor = self.surfaceColor;
    scrollView.borderType = NSNoBorder;
    scrollView.hasVerticalScroller = YES;
    scrollView.autohidesScrollers = YES;
    scrollView.automaticallyAdjustsContentInsets = YES;

    WailsInspectorDocumentView* document = [[WailsInspectorDocumentView alloc]
        initWithFrame:NSMakeRect(0, 0, 280, 600)];
    document.translatesAutoresizingMaskIntoConstraints = NO;

    NSStackView* stack = [[NSStackView alloc] initWithFrame:NSZeroRect];
    stack.translatesAutoresizingMaskIntoConstraints = NO;
    stack.orientation = NSUserInterfaceLayoutOrientationVertical;
    stack.alignment = NSLayoutAttributeLeading;
    stack.distribution = NSStackViewDistributionFill;
    stack.spacing = 12;
    [document addSubview:stack];
    scrollView.documentView = document;

    NSClipView* clip = scrollView.contentView;
    [NSLayoutConstraint activateConstraints:@[
        [document.leadingAnchor constraintEqualToAnchor:clip.leadingAnchor],
        [document.trailingAnchor constraintEqualToAnchor:clip.trailingAnchor],
        [document.topAnchor constraintEqualToAnchor:clip.topAnchor],
        [document.widthAnchor constraintEqualToAnchor:clip.widthAnchor],
        [stack.topAnchor constraintEqualToAnchor:document.topAnchor constant:18],
        [stack.leadingAnchor constraintEqualToAnchor:document.leadingAnchor constant:14],
        [stack.trailingAnchor constraintEqualToAnchor:document.trailingAnchor constant:-14],
        [stack.bottomAnchor constraintEqualToAnchor:document.bottomAnchor constant:-18]
    ]];

    self.stackView = stack;
    self.view = scrollView;
    [stack release];
    [document release];
    [scrollView release];
}

- (NSTextField*)propertyNameLabel:(NSString*)name {
    NSTextField* label = [NSTextField labelWithString:name ?: @""];
    label.font = [NSFont systemFontOfSize:[NSFont smallSystemFontSize]];
    label.textColor = [NSColor secondaryLabelColor];
    label.lineBreakMode = NSLineBreakByTruncatingTail;
    label.translatesAutoresizingMaskIntoConstraints = NO;
    [label.widthAnchor constraintEqualToConstant:86].active = YES;
    return label;
}

- (NSControl*)nativeControlForModel:(WailsInspectorControlModel*)model {
    NSControl* result = nil;
    switch (model.kind) {
        case 0: {
            NSTextField* value = [NSTextField labelWithString:model.value ?: @""];
            value.selectable = YES;
            value.lineBreakMode = NSLineBreakByTruncatingTail;
            result = value;
            break;
        }
        case 1: {
            NSTextField* field = [NSTextField textFieldWithString:model.value ?: @""];
            field.controlSize = NSControlSizeSmall;
            field.font = [NSFont systemFontOfSize:[NSFont smallSystemFontSize]];
            field.delegate = self;
            result = field;
            break;
        }
        case 2: {
            NSButton* checkbox = [NSButton checkboxWithTitle:model.label ?: @"" target:self
                action:@selector(handleCheckbox:)];
            checkbox.controlSize = NSControlSizeSmall;
            result = checkbox;
            break;
        }
        case 3: {
            NSPopUpButton* popup = [[NSPopUpButton alloc] initWithFrame:NSZeroRect pullsDown:NO];
            popup.controlSize = NSControlSizeSmall;
            [popup addItemsWithTitles:model.options ?: @[]];
            if (model.selectedIndex >= 0 && model.selectedIndex < popup.numberOfItems) {
                [popup selectItemAtIndex:model.selectedIndex];
            } else {
                [popup selectItem:nil];
            }
            popup.target = self;
            popup.action = @selector(handlePopup:);
            result = [popup autorelease];
            break;
        }
    }
    if (result != nil) {
        result.translatesAutoresizingMaskIntoConstraints = NO;
        objc_setAssociatedObject(result, WailsInspectorControlIDAssociationKey,
            @(model.controlID), OBJC_ASSOCIATION_RETAIN);
    }
    return result;
}

- (NSView*)rowForModel:(WailsInspectorControlModel*)model control:(NSControl*)control {
    NSStackView* row = [[NSStackView alloc] initWithFrame:NSZeroRect];
    row.translatesAutoresizingMaskIntoConstraints = NO;
    row.orientation = NSUserInterfaceLayoutOrientationHorizontal;
    row.alignment = model.kind == 2 ? NSLayoutAttributeCenterY : NSLayoutAttributeFirstBaseline;
    row.distribution = NSStackViewDistributionFill;
    row.spacing = 8;

    if (model.kind == 2) {
        [row addArrangedSubview:control];
    } else {
        NSTextField* nameLabel = [self propertyNameLabel:model.label];
        [row addArrangedSubview:nameLabel];
        self.nameLabelsByID[@(model.controlID)] = nameLabel;
        [row addArrangedSubview:control];
        [control setContentCompressionResistancePriority:NSLayoutPriorityDefaultLow
            forOrientation:NSLayoutConstraintOrientationHorizontal];
        [control.widthAnchor constraintGreaterThanOrEqualToConstant:90].active = YES;
    }
    row.hidden = model.hidden;
    return [row autorelease];
}

- (void)reloadContents {
    (void)self.view;
    NSArray<NSView*>* oldViews = [self.stackView.arrangedSubviews copy];
    for (NSView* view in oldViews) {
        [self.stackView removeArrangedSubview:view];
        [view removeFromSuperview];
    }
    [oldViews release];
    [self.controlsByID removeAllObjects];
    [self.rowsByID removeAllObjects];
    [self.nameLabelsByID removeAllObjects];

    BOOL firstSection = YES;
    for (WailsInspectorSectionModel* section in self.sections) {
        if (!firstSection) {
            NSBox* separator = [[[NSBox alloc] initWithFrame:NSZeroRect] autorelease];
            separator.boxType = NSBoxSeparator;
            separator.translatesAutoresizingMaskIntoConstraints = NO;
            [self.stackView addArrangedSubview:separator];
            [separator.widthAnchor constraintEqualToAnchor:self.stackView.widthAnchor].active = YES;
            [self.stackView setCustomSpacing:15 afterView:separator];
        }
        firstSection = NO;

        NSTextField* heading = [NSTextField labelWithString:section.label ?: @""];
        heading.font = [NSFont systemFontOfSize:13 weight:NSFontWeightSemibold];
        heading.textColor = [NSColor labelColor];
        heading.lineBreakMode = NSLineBreakByTruncatingTail;
        heading.translatesAutoresizingMaskIntoConstraints = NO;
        [self.stackView addArrangedSubview:heading];
        [heading.widthAnchor constraintEqualToAnchor:self.stackView.widthAnchor].active = YES;
        [self.stackView setCustomSpacing:9 afterView:heading];

        for (WailsInspectorControlModel* model in section.controls) {
            NSControl* control = [self nativeControlForModel:model];
            if (control == nil) continue;
            NSView* row = [self rowForModel:model control:control];
            [self.stackView addArrangedSubview:row];
            [row.widthAnchor constraintEqualToAnchor:self.stackView.widthAnchor].active = YES;
            self.controlsByID[@(model.controlID)] = control;
            self.rowsByID[@(model.controlID)] = row;
            [self applyModel:model];
        }
    }
}

- (void)applyModel:(WailsInspectorControlModel*)model {
    NSControl* control = self.controlsByID[@(model.controlID)];
    NSView* row = self.rowsByID[@(model.controlID)];
    if (control == nil) return;

    control.enabled = !model.disabled;
    control.toolTip = model.tooltip.length > 0 ? model.tooltip : nil;
    row.hidden = model.hidden;
    NSTextField* nameLabel = self.nameLabelsByID[@(model.controlID)];
    if (nameLabel != nil) nameLabel.stringValue = model.label ?: @"";
    switch (model.kind) {
        case 0:
        case 1:
            ((NSTextField*)control).stringValue = model.value ?: @"";
            break;
        case 2:
            ((NSButton*)control).title = model.label ?: @"";
            ((NSButton*)control).state = model.checked ? NSControlStateValueOn : NSControlStateValueOff;
            break;
        case 3: {
            NSPopUpButton* popup = (NSPopUpButton*)control;
            if (![popup.itemTitles isEqualToArray:model.options ?: @[]]) {
                [popup removeAllItems];
                [popup addItemsWithTitles:model.options ?: @[]];
            }
            if (model.selectedIndex >= 0 && model.selectedIndex < popup.numberOfItems) {
                [popup selectItemAtIndex:model.selectedIndex];
            } else {
                [popup selectItem:nil];
            }
            break;
        }
    }
}

- (void)controlTextDidChange:(NSNotification*)notification {
    NSTextField* field = notification.object;
    NSNumber* controlID = objc_getAssociatedObject(field, WailsInspectorControlIDAssociationKey);
    if (controlID != nil) {
        processMacInspectorTextChanged(controlID.unsignedLongLongValue, (char*)field.stringValue.UTF8String);
    }
}

- (void)handleCheckbox:(NSButton*)sender {
    NSNumber* controlID = objc_getAssociatedObject(sender, WailsInspectorControlIDAssociationKey);
    if (controlID != nil) {
        processMacInspectorToggleChanged(controlID.unsignedLongLongValue,
            sender.state == NSControlStateValueOn);
    }
}

- (void)handlePopup:(NSPopUpButton*)sender {
    NSNumber* controlID = objc_getAssociatedObject(sender, WailsInspectorControlIDAssociationKey);
    if (controlID != nil) {
        processMacInspectorSelectionChanged(controlID.unsignedLongLongValue, (int)sender.indexOfSelectedItem);
    }
}

@end

@interface WailsTextEditorViewController : NSViewController <NSTextViewDelegate>
@property unsigned long long editorID;
@property (retain) NSScrollView* scrollView;
@property (retain) NSTextView* textView;
@property BOOL suppressChange;
@end

@implementation WailsTextEditorViewController
- (void)loadView {
    NSScrollView* scroll = [[NSScrollView alloc] initWithFrame:NSMakeRect(0, 0, 600, 600)];
    scroll.hasVerticalScroller = YES;
    scroll.hasHorizontalScroller = NO;
    scroll.autohidesScrollers = YES;
    scroll.borderType = NSNoBorder;
    scroll.drawsBackground = YES;

    NSTextView* text = [[NSTextView alloc] initWithFrame:scroll.contentView.bounds];
    text.minSize = NSMakeSize(0, 0);
    text.maxSize = NSMakeSize(CGFLOAT_MAX, CGFLOAT_MAX);
    text.verticallyResizable = YES;
    text.horizontallyResizable = NO;
    text.autoresizingMask = NSViewWidthSizable;
    text.textContainer.widthTracksTextView = YES;
    text.textContainer.containerSize = NSMakeSize(scroll.contentSize.width, CGFLOAT_MAX);
    text.textContainerInset = NSMakeSize(24, 22);
    text.richText = NO;
    text.importsGraphics = NO;
    text.usesFindBar = YES;
    text.allowsUndo = YES;
    text.automaticQuoteSubstitutionEnabled = NO;
    text.automaticDashSubstitutionEnabled = NO;
    text.font = [NSFont userFixedPitchFontOfSize:14.0];
    text.delegate = self;
    scroll.documentView = text;

    self.scrollView = scroll;
    self.textView = text;
    self.view = scroll;
    [text release];
    [scroll release];
}
- (void)textDidChange:(NSNotification*)notification {
    if (!self.suppressChange) processMacTextEditorChanged(self.editorID);
}
- (void)setEditorText:(NSString*)value {
    self.suppressChange = YES;
    self.textView.string = value ?: @"";
    self.suppressChange = NO;
}
- (void)dealloc {
    _textView.delegate = nil;
    [_scrollView release];
    [_textView release];
    [super dealloc];
}
@end

@interface WailsSplitPaneRecord : NSObject
@property unsigned long long paneID;
@property int role;
@property BOOL primary;
@property double minThickness;
@property double maxThickness;
@property double preferredFraction;
@property BOOL hasPreferredFraction;
@property double holdingPriority;
@property BOOL hasHoldingPriority;
@property BOOL collapsible;
@property BOOL hasCollapsible;
@property BOOL canCollapseFromResize;
@property BOOL hasCanCollapseFromResize;
@property BOOL startCollapsed;
@property int contentLayout;
@property unsigned long long selectedSidebarItemID;
@property (retain) NSMutableArray<WailsSidebarNode*>* sidebarRoots;
@property (retain) WailsSidebarNode* sidebarFooter;
@property (retain) WailsSidebarViewController* sidebarController;
@property (retain) NSMutableArray<WailsInspectorSectionModel*>* inspectorSections;
@property (retain) NSMutableDictionary<NSNumber*, WailsInspectorControlModel*>* inspectorModelsByID;
@property (retain) WailsInspectorViewController* inspectorController;
@property (retain) NSViewController* viewController;
@property (retain) NSSplitViewItem* item;
#ifndef WAILS_NATIVE_ONLY
@property (retain) WKWebView* webView;
#endif
@property unsigned long long textEditorID;
@property (copy) NSString* initialText;
@property BOOL textEditorEditable;
@property (retain) WailsTextEditorViewController* textEditorController;
@property BOOL observing;
@property BOOL lastCollapsed;
@end

@implementation WailsSplitPaneRecord
- (void)dealloc {
    [_sidebarRoots release];
    [_sidebarFooter release];
    [_sidebarController release];
    [_inspectorSections release];
    [_inspectorModelsByID release];
    [_inspectorController release];
    [_viewController release];
    [_item release];
#ifndef WAILS_NATIVE_ONLY
    [_webView release];
#endif
    [_initialText release];
    [_textEditorController release];
    [super dealloc];
}

- (void)observeValueForKeyPath:(NSString*)keyPath ofObject:(id)object change:(NSDictionary*)change context:(void*)context {
    if (context != WailsSplitPaneCollapsedKVOContext) {
        [super observeValueForKeyPath:keyPath ofObject:object change:change context:context];
        return;
    }
    BOOL collapsed = [change[NSKeyValueChangeNewKey] boolValue];
    if (collapsed == self.lastCollapsed) return;
    self.lastCollapsed = collapsed;
    processMacSplitPaneCollapsed(self.paneID, collapsed);
}
@end

@interface WailsSplitViewOwner : NSObject
@property (copy) NSString* autosaveName;
@property (retain) NSSplitViewController* controller;
@property (retain) NSMutableArray<WailsSplitPaneRecord*>* records;
@property BOOL installed;
@property BOOL torndown;
@end

@implementation WailsSplitViewOwner
- (void)dealloc {
    [_autosaveName release];
    [_controller release];
    [_records release];
    [super dealloc];
}
@end

static WailsSplitViewOwner* splitViewOwner(void* handlePtr) {
    return (WailsSplitViewOwner*)handlePtr;
}

static WailsSplitPaneRecord* splitPaneRecord(void* handlePtr, unsigned long long paneID) {
    WailsSplitViewOwner* owner = splitViewOwner(handlePtr);
    if (owner == nil || owner.torndown) return nil;
    for (WailsSplitPaneRecord* record in owner.records) {
        if (record.paneID == paneID) return record;
    }
    return nil;
}

void splitViewConfigureTextEditor(void* handlePtr, unsigned long long paneID,
    unsigned long long editorID, const char* text, bool editable) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record == nil || record.role != WailsSplitPaneRolePrimary) return;
    record.textEditorID = editorID;
    NSString* value = text == NULL ? [[NSString alloc] init] : [[NSString alloc] initWithUTF8String:text];
    record.initialText = value;
    [value release];
    record.textEditorEditable = editable;
}

static WailsSidebarNode* sidebarSection(WailsSplitPaneRecord* record, unsigned long long sectionID) {
    for (WailsSidebarNode* node in record.sidebarRoots) {
        if (node.section && node.nodeID == sectionID) return node;
    }
    return nil;
}

static WailsInspectorSectionModel* inspectorSection(WailsSplitPaneRecord* record,
    unsigned long long sectionID) {
    for (WailsInspectorSectionModel* section in record.inspectorSections) {
        if (section.sectionID == sectionID) return section;
    }
    return nil;
}

static NSArray<NSString*>* inspectorOptionsFromJSON(const char* optionsJSON) {
    if (optionsJSON == NULL || strlen(optionsJSON) == 0) return @[];
    NSData* data = [[NSString stringWithUTF8String:optionsJSON] dataUsingEncoding:NSUTF8StringEncoding];
    id decoded = data == nil ? nil : [NSJSONSerialization JSONObjectWithData:data options:0 error:nil];
    if (![decoded isKindOfClass:[NSArray class]]) return @[];
    NSMutableArray<NSString*>* result = [NSMutableArray array];
    for (id value in (NSArray*)decoded) {
        if ([value isKindOfClass:[NSString class]]) [result addObject:value];
    }
    return result;
}

static void configureInspectorModel(WailsInspectorControlModel* model, int kind,
    const char* label, const char* value, bool checked, const char* optionsJSON,
    int selectedIndex, const char* tooltip, bool disabled, bool hidden) {
    if (model == nil) return;
    model.kind = kind;
    model.label = label == NULL ? @"" : [NSString stringWithUTF8String:label];
    model.value = value == NULL ? @"" : [NSString stringWithUTF8String:value];
    model.checked = checked;
    model.options = inspectorOptionsFromJSON(optionsJSON);
    model.selectedIndex = selectedIndex;
    model.tooltip = tooltip == NULL ? @"" : [NSString stringWithUTF8String:tooltip];
    model.disabled = disabled;
    model.hidden = hidden;
}

static void splitViewItemApplyCanCollapseFromResize(NSSplitViewItem* item, BOOL allowed) {
    if (item != nil && [item respondsToSelector:@selector(setCanCollapseFromWindowResize:)]) {
        [item setValue:[NSNumber numberWithBool:allowed] forKey:@"canCollapseFromWindowResize"];
    }
}

void* splitViewCreate(const char* autosaveName) {
    WailsSplitViewOwner* owner = [[WailsSplitViewOwner alloc] init];
    if (owner == nil) return NULL;
    owner.records = [NSMutableArray array];
    if (autosaveName != NULL && strlen(autosaveName) > 0) {
        owner.autosaveName = [NSString stringWithUTF8String:autosaveName];
    }
    return owner;
}

void splitViewAddPane(void* handlePtr, unsigned long long paneID, int role, bool primary,
    double minThickness, double maxThickness,
    double preferredFraction, bool hasPreferredFraction,
    double holdingPriority, bool hasHoldingPriority,
	bool collapsible, bool hasCollapsible,
	bool canCollapseFromResize, bool hasCanCollapseFromResize,
	bool startCollapsed, int contentLayout) {
    WailsSplitViewOwner* owner = splitViewOwner(handlePtr);
    if (owner == nil || owner.installed || owner.torndown) return;
    WailsSplitPaneRecord* record = [[WailsSplitPaneRecord alloc] init];
    record.paneID = paneID;
    record.role = role;
    record.primary = primary;
    record.minThickness = minThickness;
    record.maxThickness = maxThickness;
    record.preferredFraction = preferredFraction;
    record.hasPreferredFraction = hasPreferredFraction;
    record.holdingPriority = holdingPriority;
    record.hasHoldingPriority = hasHoldingPriority;
    record.collapsible = collapsible;
    record.hasCollapsible = hasCollapsible;
    record.canCollapseFromResize = canCollapseFromResize;
    record.hasCanCollapseFromResize = hasCanCollapseFromResize;
	record.startCollapsed = startCollapsed;
	record.contentLayout = contentLayout;
    record.sidebarRoots = [NSMutableArray array];
    record.inspectorSections = [NSMutableArray array];
    record.inspectorModelsByID = [NSMutableDictionary dictionary];
    [owner.records addObject:record];
    [record release];
}

#ifndef WAILS_NATIVE_ONLY
bool splitViewInstall(void* handlePtr, void* nsWindow, bool normalBackdrop) {
    WailsSplitViewOwner* owner = splitViewOwner(handlePtr);
    WebviewWindow* window = (WebviewWindow*)nsWindow;
    if (owner == nil || window == nil || owner.installed || owner.torndown || owner.records.count < 2) return false;
    WKWebView* primaryWebView = window.webView;
    if (primaryWebView == nil) return false;

    WailsSplitPaneRecord* primaryRecord = nil;
    for (WailsSplitPaneRecord* record in owner.records) {
        if (record.primary) primaryRecord = record;
        if (!record.primary && record.role != WailsSplitPaneRoleSidebar &&
            record.role != WailsSplitPaneRoleInspector) return false;
    }
    if (primaryRecord == nil) return false;

    // WebviewWindow starts transparent because that is required by Wails'
    // explicit transparent backdrop modes. The old single WebView concealed
    // that implementation detail. A semantic sidebar does not: its material
    // would blend straight through a clear NSWindow and become much more
    // transparent than a normal AppKit sidebar. Restore an opaque native
    // window surface only for the documented normal-backdrop mode.
    NSColor* primaryBackground = nil;
    if (normalBackdrop) {
        primaryBackground = window.backgroundColor;
        if (primaryBackground == nil || primaryBackground.alphaComponent <= 0.0) {
            primaryBackground = [NSColor windowBackgroundColor];
        }
        window.backgroundColor = primaryBackground;
        window.opaque = YES;
    }

    // A transparent document does not make WKWebView's native backing view
    // transparent. Split windows need the WebView to reveal the AppKit pane
    // surface beneath it, otherwise WebKit, the source list, and the unified
    // titlebar resolve to subtly different system colours.
    [primaryWebView setValue:@NO forKey:@"drawsBackground"];

    // Construct every controller before touching the window hierarchy.
    for (WailsSplitPaneRecord* record in owner.records) {
        if (record.primary) {
            NSViewController* paneController = [[NSViewController alloc] init];
            WailsPrimaryPaneView* container = [[WailsPrimaryPaneView alloc] initWithFrame:NSMakeRect(0, 0, 600, 600)];
            container.fillColor = primaryBackground;
            paneController.view = container;
            [container release];
            record.viewController = paneController;
            [paneController release];
        } else if (record.role == WailsSplitPaneRoleSidebar) {
            WailsSidebarViewController* sidebar = [[WailsSidebarViewController alloc] init];
            sidebar.roots = record.sidebarRoots;
            sidebar.footer = record.sidebarFooter;
            sidebar.selectedItemID = record.selectedSidebarItemID;
            (void)sidebar.view;
            if (sidebar.view == nil) {
                [sidebar release];
                return false;
            }
            record.sidebarController = sidebar;
            record.viewController = sidebar;
            [sidebar release];
        } else {
            WailsInspectorViewController* inspector = [[WailsInspectorViewController alloc] init];
            inspector.sections = record.inspectorSections;
            inspector.modelsByID = record.inspectorModelsByID;
            inspector.surfaceColor = primaryBackground;
            (void)inspector.view;
            if (inspector.view == nil) {
                [inspector release];
                return false;
            }
            record.inspectorController = inspector;
            record.viewController = inspector;
            [inspector release];
        }
    }

    NSSplitViewController* controller = [[NSSplitViewController alloc] init];
    if (controller == nil) return false;
    controller.splitView.vertical = YES;
    if (owner.autosaveName.length > 0) controller.splitView.autosaveName = owner.autosaveName;

    // Create and configure every semantic split item before reparenting the
    // primary WebView. Any failure still leaves the original window intact.
    for (WailsSplitPaneRecord* record in owner.records) {
        NSSplitViewItem* item = nil;
        if (record.role == WailsSplitPaneRoleSidebar) {
            item = [NSSplitViewItem sidebarWithViewController:record.viewController];
        } else if (record.role == WailsSplitPaneRoleInspector) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
            if (@available(macOS 11.0, *)) {
                item = [NSSplitViewItem inspectorWithViewController:record.viewController];
            }
#endif
            if (item == nil) item = [NSSplitViewItem splitViewItemWithViewController:record.viewController];
        } else {
            item = [NSSplitViewItem splitViewItemWithViewController:record.viewController];
        }
        if (item == nil) {
            [controller release];
            return false;
        }
        if (record.minThickness > 0) item.minimumThickness = record.minThickness;
        if (record.maxThickness > 0) item.maximumThickness = record.maxThickness;
        if (record.hasPreferredFraction) item.preferredThicknessFraction = record.preferredFraction;
        if (record.hasHoldingPriority) item.holdingPriority = record.holdingPriority;
        if (record.hasCollapsible) item.canCollapse = record.collapsible;
        if (record.hasCanCollapseFromResize) splitViewItemApplyCanCollapseFromResize(item, record.canCollapseFromResize);
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
        if (@available(macOS 11.0, *)) {
            if (record.role == WailsSplitPaneRoleSidebar) item.allowsFullHeightLayout = YES;
        }
#endif
		[controller addSplitViewItem:item];
		record.item = item;
    }

    [primaryWebView retain];
    [primaryWebView removeFromSuperview];

    NSView* backdropView = nil;
    Class glassEffectViewClass = NSClassFromString(@"NSGlassEffectView");
    for (NSView* subview in window.contentView.subviews) {
        if ([subview isKindOfClass:[NSVisualEffectView class]] ||
            (glassEffectViewClass != nil && [subview isKindOfClass:glassEffectViewClass])) {
            backdropView = subview;
            break;
        }
    }
    if (backdropView != nil) {
        [backdropView retain];
        [backdropView removeFromSuperview];
    }

    NSView* dragView = nil;
    Class dragClass = NSClassFromString(@"WebviewDrag");
    if (dragClass != nil) {
        for (NSView* subview in window.contentView.subviews) {
            if ([subview isKindOfClass:dragClass]) {
                dragView = subview;
                break;
            }
        }
    }
    if (dragView != nil) {
        [dragView retain];
        [dragView removeFromSuperview];
    }

    for (WailsSplitPaneRecord* record in owner.records) {
        if (record.primary) {
            NSView* container = record.viewController.view;
            if (backdropView != nil) {
                backdropView.frame = container.bounds;
                backdropView.autoresizingMask = NSViewWidthSizable | NSViewHeightSizable;
                [container addSubview:backdropView];
                NSView* primaryHost = [backdropView respondsToSelector:@selector(contentView)]
                    ? [backdropView valueForKey:@"contentView"] : nil;
                if (primaryHost != nil) {
                    primaryWebView.frame = primaryHost.bounds;
                    primaryWebView.autoresizingMask = NSViewWidthSizable | NSViewHeightSizable;
                    [primaryHost addSubview:primaryWebView];
                } else {
                    primaryWebView.frame = container.bounds;
                    primaryWebView.autoresizingMask = NSViewWidthSizable | NSViewHeightSizable;
                    [container addSubview:primaryWebView positioned:NSWindowAbove relativeTo:backdropView];
                }
            } else {
                primaryWebView.frame = container.bounds;
                primaryWebView.autoresizingMask = NSViewWidthSizable | NSViewHeightSizable;
                [container addSubview:primaryWebView];
            }
            if (dragView != nil) {
                dragView.frame = container.bounds;
                dragView.autoresizingMask = NSViewWidthSizable | NSViewHeightSizable;
                [container addSubview:dragView];
            }
            record.webView = primaryWebView;
        }

    }

    for (WailsSplitPaneRecord* record in owner.records) {
        if (record.startCollapsed) record.item.collapsed = YES;
        record.lastCollapsed = record.item.collapsed;
        [record.item addObserver:record forKeyPath:@"collapsed"
            options:NSKeyValueObservingOptionNew context:WailsSplitPaneCollapsedKVOContext];
        record.observing = YES;
    }

    window.styleMask |= NSWindowStyleMaskFullSizeContentView;
    NSRect frame = window.frame;
    window.contentViewController = controller;
    windowApplyContentLayout(window, primaryRecord.contentLayout);
    [window setFrame:frame display:YES];
    objc_setAssociatedObject(primaryWebView, WailsSplitPrimaryPaneIDAssociationKey,
        [NSNumber numberWithUnsignedLongLong:primaryRecord.paneID], OBJC_ASSOCIATION_RETAIN);

    for (WailsSplitPaneRecord* record in owner.records) {
        [record.sidebarController reloadContents];
        [record.inspectorController reloadContents];
    }

    [primaryWebView release];
    if (dragView != nil) [dragView release];
    if (backdropView != nil) [backdropView release];
    owner.controller = controller;
    [controller release];
    owner.installed = YES;
    return true;
}
#else
bool splitViewInstall(void* handlePtr, void* nsWindow, bool normalBackdrop) {
    return false;
}
#endif

bool splitViewInstallNative(void* handlePtr, void* nsWindow, bool normalBackdrop) {
    WailsSplitViewOwner* owner = splitViewOwner(handlePtr);
    NSWindow* window = (NSWindow*)nsWindow;
    if (owner == nil || window == nil || owner.installed || owner.torndown || owner.records.count < 2) return false;

    WailsSplitPaneRecord* primaryRecord = nil;
    for (WailsSplitPaneRecord* record in owner.records) {
        if (record.primary) primaryRecord = record;
        if (!record.primary && record.role != WailsSplitPaneRoleSidebar &&
            record.role != WailsSplitPaneRoleInspector) return false;
    }
    if (primaryRecord == nil || primaryRecord.textEditorID == 0) return false;

    NSColor* primaryBackground = nil;
    if (normalBackdrop) {
        primaryBackground = window.backgroundColor;
        if (primaryBackground == nil || primaryBackground.alphaComponent <= 0.0) {
            primaryBackground = [NSColor windowBackgroundColor];
        }
        window.backgroundColor = primaryBackground;
        window.opaque = YES;
    }

    for (WailsSplitPaneRecord* record in owner.records) {
        if (record.primary) {
            WailsTextEditorViewController* editor = [[WailsTextEditorViewController alloc] init];
            editor.editorID = record.textEditorID;
            (void)editor.view;
            if (editor.view == nil) {
                [editor release];
                return false;
            }
            [editor setEditorText:record.initialText ?: @""];
            // The NSTextView is authoritative after installation. Keeping the
            // staging NSString here would retain a second complete document.
            record.initialText = nil;
            editor.textView.editable = record.textEditorEditable;
            record.textEditorController = editor;
            record.viewController = editor;
            [editor release];
        } else if (record.role == WailsSplitPaneRoleSidebar) {
            WailsSidebarViewController* sidebar = [[WailsSidebarViewController alloc] init];
            sidebar.roots = record.sidebarRoots;
            sidebar.footer = record.sidebarFooter;
            sidebar.selectedItemID = record.selectedSidebarItemID;
            (void)sidebar.view;
            if (sidebar.view == nil) {
                [sidebar release];
                return false;
            }
            record.sidebarController = sidebar;
            record.viewController = sidebar;
            [sidebar release];
        } else {
            WailsInspectorViewController* inspector = [[WailsInspectorViewController alloc] init];
            inspector.sections = record.inspectorSections;
            inspector.modelsByID = record.inspectorModelsByID;
            inspector.surfaceColor = primaryBackground;
            (void)inspector.view;
            if (inspector.view == nil) {
                [inspector release];
                return false;
            }
            record.inspectorController = inspector;
            record.viewController = inspector;
            [inspector release];
        }
    }

    NSSplitViewController* controller = [[NSSplitViewController alloc] init];
    if (controller == nil) return false;
    controller.splitView.vertical = YES;
    if (owner.autosaveName.length > 0) controller.splitView.autosaveName = owner.autosaveName;

    for (WailsSplitPaneRecord* record in owner.records) {
        NSSplitViewItem* item = nil;
        if (record.role == WailsSplitPaneRoleSidebar) {
            item = [NSSplitViewItem sidebarWithViewController:record.viewController];
        } else if (record.role == WailsSplitPaneRoleInspector) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
            if (@available(macOS 11.0, *)) {
                item = [NSSplitViewItem inspectorWithViewController:record.viewController];
            }
#endif
            if (item == nil) item = [NSSplitViewItem splitViewItemWithViewController:record.viewController];
        } else {
            item = [NSSplitViewItem splitViewItemWithViewController:record.viewController];
        }
        if (item == nil) {
            [controller release];
            return false;
        }
        if (record.minThickness > 0) item.minimumThickness = record.minThickness;
        if (record.maxThickness > 0) item.maximumThickness = record.maxThickness;
        if (record.hasPreferredFraction) item.preferredThicknessFraction = record.preferredFraction;
        if (record.hasHoldingPriority) item.holdingPriority = record.holdingPriority;
        if (record.hasCollapsible) item.canCollapse = record.collapsible;
        if (record.hasCanCollapseFromResize) splitViewItemApplyCanCollapseFromResize(item, record.canCollapseFromResize);
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
        if (@available(macOS 11.0, *)) {
            if (record.role == WailsSplitPaneRoleSidebar) item.allowsFullHeightLayout = YES;
        }
#endif
		[controller addSplitViewItem:item];
		record.item = item;
    }

    for (WailsSplitPaneRecord* record in owner.records) {
        if (record.startCollapsed) record.item.collapsed = YES;
        record.lastCollapsed = record.item.collapsed;
        [record.item addObserver:record forKeyPath:@"collapsed"
            options:NSKeyValueObservingOptionNew context:WailsSplitPaneCollapsedKVOContext];
        record.observing = YES;
    }

    window.styleMask |= NSWindowStyleMaskFullSizeContentView;
    NSRect frame = window.frame;
    window.contentViewController = controller;
    [window setFrame:frame display:YES];
    for (WailsSplitPaneRecord* record in owner.records) {
        [record.sidebarController reloadContents];
        [record.inspectorController reloadContents];
    }
    owner.controller = controller;
    [controller release];
    owner.installed = YES;
    return true;
}

void splitViewTextEditorSetText(void* handlePtr, unsigned long long paneID, const char* text) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record == nil || record.textEditorController == nil) return;
    NSString* value = text == NULL ? [[NSString alloc] init] : [[NSString alloc] initWithUTF8String:text];
    [record.textEditorController setEditorText:value];
    [value release];
}

char* splitViewTextEditorCopyText(void* handlePtr, unsigned long long paneID) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record == nil || record.textEditorController == nil) return NULL;
    const char* value = record.textEditorController.textView.string.UTF8String;
    return value == NULL ? strdup("") : strdup(value);
}

void splitViewTextEditorSetEditable(void* handlePtr, unsigned long long paneID, bool editable) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record != nil) record.textEditorController.textView.editable = editable;
}

void splitViewTextEditorFocus(void* handlePtr, unsigned long long paneID) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record == nil || record.textEditorController == nil) return;
    [record.textEditorController.view.window makeFirstResponder:record.textEditorController.textView];
}

void splitViewTeardown(void* handlePtr) {
    WailsSplitViewOwner* owner = splitViewOwner(handlePtr);
    if (owner == nil || owner.torndown) return;
    owner.torndown = YES;
    for (WailsSplitPaneRecord* record in owner.records) {
        if (record.observing && record.item != nil) {
            [record.item removeObserver:record forKeyPath:@"collapsed" context:WailsSplitPaneCollapsedKVOContext];
            record.observing = NO;
        }
#ifndef WAILS_NATIVE_ONLY
        if (record.primary && record.webView != nil) {
            objc_setAssociatedObject(record.webView, WailsSplitPrimaryPaneIDAssociationKey, nil, OBJC_ASSOCIATION_RETAIN);
        }
        record.webView = nil;
#endif
        record.item = nil;
        record.viewController = nil;
        record.sidebarController = nil;
        record.inspectorController = nil;
        record.textEditorController = nil;
        [record.sidebarRoots removeAllObjects];
        [record.inspectorSections removeAllObjects];
        [record.inspectorModelsByID removeAllObjects];
    }
    [owner.records removeAllObjects];
    owner.controller = nil;
}

void splitViewRelease(void* handlePtr) {
    [splitViewOwner(handlePtr) release];
}

void splitViewPaneSetMinimumThickness(void* handlePtr, unsigned long long paneID, double value) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record.item != nil) record.item.minimumThickness = value == 0 ? NSSplitViewItemUnspecifiedDimension : value;
}
void splitViewPaneSetMaximumThickness(void* handlePtr, unsigned long long paneID, double value) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record.item != nil) record.item.maximumThickness = value == 0 ? NSSplitViewItemUnspecifiedDimension : value;
}
void splitViewPaneSetPreferredFraction(void* handlePtr, unsigned long long paneID, double value) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record.item != nil) record.item.preferredThicknessFraction = value;
}
void splitViewPaneSetHoldingPriority(void* handlePtr, unsigned long long paneID, double value) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record.item != nil) record.item.holdingPriority = value;
}
void splitViewPaneSetCollapsible(void* handlePtr, unsigned long long paneID, bool collapsible) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record.item != nil) record.item.canCollapse = collapsible;
}
void splitViewPaneSetCanCollapseFromWindowResize(void* handlePtr, unsigned long long paneID, bool allowed) {
    splitViewItemApplyCanCollapseFromResize(splitPaneRecord(handlePtr, paneID).item, allowed);
}
void splitViewPaneSetContentLayout(void* handlePtr, unsigned long long paneID, int layout) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record == nil || !record.primary) return;
    record.contentLayout = layout;
#ifndef WAILS_NATIVE_ONLY
    if (record.webView != nil) windowApplyContentLayout(record.webView.window, layout);
#endif
}
void splitViewPaneSetCollapsed(void* handlePtr, unsigned long long paneID, bool collapsed) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record.item != nil) ((NSSplitViewItem*)record.item.animator).collapsed = collapsed;
}
void splitViewPaneToggleCollapsed(void* handlePtr, unsigned long long paneID) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record.item != nil) ((NSSplitViewItem*)record.item.animator).collapsed = !record.item.collapsed;
}

void splitViewSidebarReset(void* handlePtr, unsigned long long paneID) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record == nil || record.role != WailsSplitPaneRoleSidebar) return;
    [record.sidebarRoots removeAllObjects];
    record.sidebarFooter = nil;
    record.sidebarController.footer = nil;
    [record.sidebarController reloadContents];
}

void splitViewSidebarAddSection(void* handlePtr, unsigned long long paneID,
    unsigned long long sectionID, const char* label) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record == nil || record.role != WailsSplitPaneRoleSidebar) return;
    WailsSidebarNode* node = [[WailsSidebarNode alloc] init];
    node.nodeID = sectionID;
    node.section = YES;
    node.label = label == NULL ? @"" : [NSString stringWithUTF8String:label];
    node.children = [NSMutableArray array];
    [record.sidebarRoots addObject:node];
    [node release];
}

void splitViewSidebarAddItem(void* handlePtr, unsigned long long paneID,
    unsigned long long sectionID, unsigned long long itemID,
    const char* label, const char* subtitle, const char* symbolName,
    const unsigned char* imageData, size_t imageDataLength, const char* tooltip,
    bool disabled, bool hidden) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record == nil || record.role != WailsSplitPaneRoleSidebar) return;
    WailsSidebarNode* node = [[WailsSidebarNode alloc] init];
    node.nodeID = itemID;
    node.label = label == NULL ? @"" : [NSString stringWithUTF8String:label];
    node.subtitle = subtitle == NULL ? @"" : [NSString stringWithUTF8String:subtitle];
    node.symbolName = symbolName == NULL ? @"" : [NSString stringWithUTF8String:symbolName];
    node.imageData = imageData == NULL || imageDataLength == 0 ? nil :
        [NSData dataWithBytes:imageData length:imageDataLength];
    node.tooltip = tooltip == NULL ? @"" : [NSString stringWithUTF8String:tooltip];
    node.disabled = disabled;
    node.hidden = hidden;
    WailsSidebarNode* section = sectionID == 0 ? nil : sidebarSection(record, sectionID);
    if (section != nil) [section.children addObject:node];
    else [record.sidebarRoots addObject:node];
    [node release];
}

void splitViewSidebarSetFooter(void* handlePtr, unsigned long long paneID,
    unsigned long long itemID, const char* label, const char* subtitle, const char* symbolName,
    const unsigned char* imageData, size_t imageDataLength, const char* tooltip,
    bool disabled, bool hidden) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record == nil || record.role != WailsSplitPaneRoleSidebar) return;
    WailsSidebarNode* node = [[WailsSidebarNode alloc] init];
    node.nodeID = itemID;
    node.label = label == NULL ? @"" : [NSString stringWithUTF8String:label];
    node.subtitle = subtitle == NULL ? @"" : [NSString stringWithUTF8String:subtitle];
    node.symbolName = symbolName == NULL ? @"" : [NSString stringWithUTF8String:symbolName];
    node.imageData = imageData == NULL || imageDataLength == 0 ? nil :
        [NSData dataWithBytes:imageData length:imageDataLength];
    node.tooltip = tooltip == NULL ? @"" : [NSString stringWithUTF8String:tooltip];
    node.disabled = disabled;
    node.hidden = hidden;
    record.sidebarFooter = node;
    record.sidebarController.footer = node;
    [node release];
}

void splitViewSidebarSetSelectedItem(void* handlePtr, unsigned long long paneID,
    unsigned long long itemID) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record == nil || record.role != WailsSplitPaneRoleSidebar) return;
    record.selectedSidebarItemID = itemID;
    record.sidebarController.selectedItemID = itemID;
    [record.sidebarController reloadContents];
}

void splitViewInspectorReset(void* handlePtr, unsigned long long paneID) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record == nil || record.role != WailsSplitPaneRoleInspector) return;
    [record.inspectorSections removeAllObjects];
    [record.inspectorModelsByID removeAllObjects];
}

void splitViewInspectorAddSection(void* handlePtr, unsigned long long paneID,
    unsigned long long sectionID, const char* label) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record == nil || record.role != WailsSplitPaneRoleInspector) return;
    WailsInspectorSectionModel* section = [[WailsInspectorSectionModel alloc] init];
    section.sectionID = sectionID;
    section.label = label == NULL ? @"" : [NSString stringWithUTF8String:label];
    section.controls = [NSMutableArray array];
    [record.inspectorSections addObject:section];
    [section release];
}

void splitViewInspectorAddControl(void* handlePtr, unsigned long long paneID,
    unsigned long long sectionID, unsigned long long controlID, int kind,
    const char* label, const char* value, bool checked, const char* optionsJSON,
    int selectedIndex, const char* tooltip, bool disabled, bool hidden) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    WailsInspectorSectionModel* section = inspectorSection(record, sectionID);
    if (record == nil || record.role != WailsSplitPaneRoleInspector || section == nil) return;
    WailsInspectorControlModel* model = [[WailsInspectorControlModel alloc] init];
    model.controlID = controlID;
    configureInspectorModel(model, kind, label, value, checked, optionsJSON,
        selectedIndex, tooltip, disabled, hidden);
    [section.controls addObject:model];
    record.inspectorModelsByID[@(controlID)] = model;
    [model release];
}

void splitViewInspectorReload(void* handlePtr, unsigned long long paneID) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record == nil || record.role != WailsSplitPaneRoleInspector) return;
    [record.inspectorController reloadContents];
}

void splitViewInspectorUpdateControl(void* handlePtr, unsigned long long paneID,
    unsigned long long controlID, int kind, const char* label, const char* value,
    bool checked, const char* optionsJSON, int selectedIndex, const char* tooltip,
    bool disabled, bool hidden) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record == nil || record.role != WailsSplitPaneRoleInspector) return;
    WailsInspectorControlModel* model = record.inspectorModelsByID[@(controlID)];
    if (model == nil || model.kind != kind) return;
    configureInspectorModel(model, kind, label, value, checked, optionsJSON,
        selectedIndex, tooltip, disabled, hidden);
    [record.inspectorController applyModel:model];
}
