//go:build darwin && !ios && !server

#import "webview_window_split_darwin.h"
#import "webview_window_contentlist_darwin.h"
#import "mac_private_api_darwin.h"
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

// Sidebar model and cells. Nodes mirror the Go MacSidebar tree: sections at
// the root, rows beneath sections, rows, or the root itself.
static NSString* const WailsSidebarDragType = @"io.wails.sidebar.row";

@interface WailsSidebarNode : NSObject
@property unsigned long long nodeID;
@property BOOL section;
@property (copy) NSString* label;
@property (copy) NSString* symbolName;
@property (copy) NSString* tooltip;
@property (copy) NSString* accessorySymbol;
@property (retain) NSColor* tintColor;
@property BOOL disabled;
@property BOOL hidden;
@property BOOL expanded;
@property BOOL editable;
@property NSInteger badge;
@property (assign) WailsSidebarNode* parent;
@property (retain) NSMutableArray<WailsSidebarNode*>* children;
@end

@implementation WailsSidebarNode
- (instancetype)init {
    self = [super init];
    if (self) _children = [[NSMutableArray alloc] init];
    return self;
}

- (void)dealloc {
    [_label release];
    [_symbolName release];
    [_tooltip release];
    [_accessorySymbol release];
    [_tintColor release];
    [_children release];
    [super dealloc];
}
@end

static WailsSidebarNode* sidebarNodeInTree(NSArray<WailsSidebarNode*>* nodes, unsigned long long nodeID) {
    for (WailsSidebarNode* node in nodes) {
        if (node.nodeID == nodeID) return node;
        WailsSidebarNode* found = sidebarNodeInTree(node.children, nodeID);
        if (found != nil) return found;
    }
    return nil;
}

// sidebarNodeHasAncestor reports whether node is ancestor or nests beneath it.
static BOOL sidebarNodeHasAncestor(WailsSidebarNode* node, WailsSidebarNode* ancestor) {
    for (WailsSidebarNode* current = node; current != nil; current = current.parent) {
        if (current == ancestor) return YES;
    }
    return NO;
}

static NSImage* sidebarSymbolImage(NSString* symbolName, NSString* description) {
    NSImage* image = nil;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
    if (@available(macOS 11.0, *)) {
        image = symbolName.length > 0
            ? [NSImage imageWithSystemSymbolName:symbolName accessibilityDescription:description]
            : nil;
    }
#endif
    return image;
}

// One source-list row: leading symbol, label, trailing badge count, and an
// optional trailing accessory symbol. Hidden trailing views leave the layout.
@interface WailsSidebarItemCell : NSTableCellView
@property (retain) NSTextField* badgeLabel;
@property (retain) NSImageView* accessoryView;
@end

@implementation WailsSidebarItemCell
- (void)dealloc {
    [_badgeLabel release];
    [_accessoryView release];
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

@class WailsSidebarViewController;

// The outline subclass routes right-clicks and Return to the controller so
// context menus and inline rename behave like Finder's source list.
@interface WailsSidebarOutlineView : NSOutlineView
@property (assign) WailsSidebarViewController* controller;
@end

// A real AppKit source list. Both the outline and its scroll view remain
// transparent so NSSplitViewItem's semantic sidebar material is visible.
@interface WailsSidebarViewController : NSViewController <NSOutlineViewDataSource, NSOutlineViewDelegate, NSTextFieldDelegate>
@property (retain) NSMutableArray<WailsSidebarNode*>* roots;
@property (retain) NSOutlineView* outlineView;
@property (retain) NSColor* surfaceColor;
@property unsigned long long paneID;
@property unsigned long long selectedItemID;
@property BOOL allowsMultipleSelection;
@property BOOL reorderable;
@property BOOL suppressSelectionCallback;
@property BOOL suppressExpansionCallback;
- (void)reloadContents;
- (void)applyOptions;
- (NSMenu*)contextMenuForEvent:(NSEvent*)event;
- (BOOL)beginEditingSelectedRow;
@end

@implementation WailsSidebarViewController

- (void)dealloc {
    [_roots release];
    [_outlineView release];
    [_surfaceColor release];
    [super dealloc];
}

- (void)loadView {
    NSScrollView* scrollView = [[NSScrollView alloc] initWithFrame:NSMakeRect(0, 0, 240, 600)];
    scrollView.drawsBackground = self.surfaceColor != nil;
    if (self.surfaceColor != nil) scrollView.backgroundColor = self.surfaceColor;
    scrollView.borderType = NSNoBorder;
    scrollView.hasVerticalScroller = YES;
    scrollView.autohidesScrollers = YES;

    WailsSidebarOutlineView* outlineView = [[WailsSidebarOutlineView alloc] initWithFrame:scrollView.bounds];
    outlineView.controller = self;
    outlineView.backgroundColor = [NSColor clearColor];
    outlineView.headerView = nil;
    outlineView.floatsGroupRows = YES;
    outlineView.indentationPerLevel = 12.0;
    outlineView.autoresizesOutlineColumn = NO;
    outlineView.selectionHighlightStyle = NSTableViewSelectionHighlightStyleSourceList;
    outlineView.allowsEmptySelection = YES;
    outlineView.allowsMultipleSelection = self.allowsMultipleSelection;
    outlineView.rowSizeStyle = NSTableViewRowSizeStyleDefault;
    outlineView.verticalMotionCanBeginDrag = YES;
    [outlineView registerForDraggedTypes:@[WailsSidebarDragType]];
    [outlineView setDraggingSourceOperationMask:NSDragOperationMove forLocal:YES];
    [outlineView setDraggingSourceOperationMask:NSDragOperationNone forLocal:NO];
    outlineView.target = self;
    outlineView.doubleAction = @selector(handleDoubleClick:);

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
    self.outlineView = outlineView;
    self.view = scrollView;
    [outlineView release];
    [scrollView release];
}

- (void)applyOptions {
    self.outlineView.allowsMultipleSelection = self.allowsMultipleSelection;
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
    WailsSidebarNode* node = (WailsSidebarNode*)item;
    return node.section || [self visibleNodes:node.children].count > 0;
}

- (BOOL)outlineView:(NSOutlineView*)outlineView isGroupItem:(id)item {
    return ((WailsSidebarNode*)item).section;
}

- (NSTableCellView*)newCellWithIdentifier:(NSUserInterfaceItemIdentifier)identifier section:(BOOL)section {
    NSTextField* textField = [NSTextField labelWithString:@""];
    textField.translatesAutoresizingMaskIntoConstraints = NO;
    textField.lineBreakMode = NSLineBreakByTruncatingTail;
    if (section) {
        NSTableCellView* cell = [[[NSTableCellView alloc] initWithFrame:NSMakeRect(0, 0, 220, 22)] autorelease];
        cell.identifier = identifier;
        textField.font = [NSFont systemFontOfSize:11 weight:NSFontWeightSemibold];
        textField.textColor = [NSColor secondaryLabelColor];
        [cell addSubview:textField];
        [NSLayoutConstraint activateConstraints:@[
            [textField.leadingAnchor constraintEqualToAnchor:cell.leadingAnchor constant:4],
            [textField.trailingAnchor constraintEqualToAnchor:cell.trailingAnchor constant:-4],
            [textField.centerYAnchor constraintEqualToAnchor:cell.centerYAnchor]
        ]];
        cell.textField = textField;
        return cell;
    }

    WailsSidebarItemCell* cell = [[[WailsSidebarItemCell alloc] initWithFrame:NSMakeRect(0, 0, 220, 28)] autorelease];
    cell.identifier = identifier;
    textField.delegate = self;
    [textField setContentHuggingPriority:NSLayoutPriorityDefaultLow
        forOrientation:NSLayoutConstraintOrientationHorizontal];
    [textField setContentCompressionResistancePriority:NSLayoutPriorityDefaultLow
        forOrientation:NSLayoutConstraintOrientationHorizontal];

    NSImageView* imageView = [[[NSImageView alloc] initWithFrame:NSZeroRect] autorelease];
    imageView.translatesAutoresizingMaskIntoConstraints = NO;
    imageView.imageScaling = NSImageScaleProportionallyUpOrDown;
    [NSLayoutConstraint activateConstraints:@[
        [imageView.widthAnchor constraintEqualToConstant:16],
        [imageView.heightAnchor constraintEqualToConstant:16]
    ]];

    NSTextField* badge = [NSTextField labelWithString:@""];
    badge.translatesAutoresizingMaskIntoConstraints = NO;
    badge.font = [NSFont monospacedDigitSystemFontOfSize:[NSFont smallSystemFontSize] weight:NSFontWeightMedium];
    badge.textColor = [NSColor secondaryLabelColor];
    badge.alignment = NSTextAlignmentRight;
    [badge setContentHuggingPriority:NSLayoutPriorityRequired
        forOrientation:NSLayoutConstraintOrientationHorizontal];

    NSImageView* accessory = [[[NSImageView alloc] initWithFrame:NSZeroRect] autorelease];
    accessory.translatesAutoresizingMaskIntoConstraints = NO;
    accessory.imageScaling = NSImageScaleProportionallyDown;
    [NSLayoutConstraint activateConstraints:@[
        [accessory.widthAnchor constraintEqualToConstant:14],
        [accessory.heightAnchor constraintEqualToConstant:14]
    ]];

    NSStackView* stack = [NSStackView stackViewWithViews:@[imageView, textField, badge, accessory]];
    stack.translatesAutoresizingMaskIntoConstraints = NO;
    stack.orientation = NSUserInterfaceLayoutOrientationHorizontal;
    stack.alignment = NSLayoutAttributeCenterY;
    stack.distribution = NSStackViewDistributionFill;
    stack.spacing = 6;
    [stack setCustomSpacing:7 afterView:imageView];
    [cell addSubview:stack];
    [NSLayoutConstraint activateConstraints:@[
        [stack.leadingAnchor constraintEqualToAnchor:cell.leadingAnchor constant:2],
        [stack.trailingAnchor constraintEqualToAnchor:cell.trailingAnchor constant:-6],
        [stack.centerYAnchor constraintEqualToAnchor:cell.centerYAnchor]
    ]];
    cell.imageView = imageView;
    cell.textField = textField;
    cell.badgeLabel = badge;
    cell.accessoryView = accessory;
    return cell;
}

- (NSView*)outlineView:(NSOutlineView*)outlineView viewForTableColumn:(NSTableColumn*)tableColumn item:(id)item {
    WailsSidebarNode* node = (WailsSidebarNode*)item;
    NSUserInterfaceItemIdentifier identifier = node.section ? @"WailsSidebarSection" : @"WailsSidebarItem";
    NSTableCellView* cell = [outlineView makeViewWithIdentifier:identifier owner:self];
    if (cell == nil) cell = [self newCellWithIdentifier:identifier section:node.section];
    cell.textField.stringValue = node.label ?: @"";
    cell.toolTip = node.tooltip.length > 0 ? node.tooltip : nil;
    cell.textField.textColor = node.disabled ? [NSColor disabledControlTextColor] :
        (node.section ? [NSColor secondaryLabelColor] : [NSColor labelColor]);
    if (node.section) return cell;

    NSImage* image = sidebarSymbolImage(node.symbolName, node.label);
    cell.imageView.image = image;
    cell.imageView.hidden = image == nil;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101400
    if (@available(macOS 10.14, *)) {
        cell.imageView.contentTintColor = node.disabled ? nil : node.tintColor;
    }
#endif
    cell.textField.editable = node.editable && !node.disabled;
    cell.textField.selectable = cell.textField.editable;

    WailsSidebarItemCell* itemCell = [cell isKindOfClass:[WailsSidebarItemCell class]] ? (WailsSidebarItemCell*)cell : nil;
    itemCell.badgeLabel.stringValue = node.badge > 0 ? [NSString stringWithFormat:@"%ld", (long)node.badge] : @"";
    itemCell.badgeLabel.hidden = node.badge <= 0;
    itemCell.badgeLabel.textColor = node.disabled ? [NSColor disabledControlTextColor] : [NSColor secondaryLabelColor];
    NSImage* accessoryImage = sidebarSymbolImage(node.accessorySymbol, node.accessorySymbol);
    itemCell.accessoryView.image = accessoryImage;
    itemCell.accessoryView.hidden = accessoryImage == nil;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101400
    if (@available(macOS 10.14, *)) {
        itemCell.accessoryView.contentTintColor = [NSColor secondaryLabelColor];
    }
#endif
    return cell;
}

- (BOOL)outlineView:(NSOutlineView*)outlineView shouldSelectItem:(id)item {
    WailsSidebarNode* node = (WailsSidebarNode*)item;
    return !node.section && !node.disabled;
}

- (void)outlineViewSelectionDidChange:(NSNotification*)notification {
    if (self.suppressSelectionCallback) return;
    NSIndexSet* rows = self.outlineView.selectedRowIndexes;
    NSMutableArray<NSNumber*>* ids = [NSMutableArray arrayWithCapacity:rows.count];
    NSInteger clicked = self.outlineView.clickedRow;
    if (clicked >= 0 && [rows containsIndex:(NSUInteger)clicked]) {
        WailsSidebarNode* node = [self.outlineView itemAtRow:clicked];
        if (node != nil && !node.section && !node.disabled) [ids addObject:@(node.nodeID)];
    }
    [rows enumerateIndexesUsingBlock:^(NSUInteger row, BOOL* stop) {
        if ((NSInteger)row == clicked) return;
        WailsSidebarNode* node = [self.outlineView itemAtRow:(NSInteger)row];
        if (node == nil || node.section || node.disabled) return;
        [ids addObject:@(node.nodeID)];
    }];
    self.selectedItemID = ids.count > 0 ? ids[0].unsignedLongLongValue : 0;
    unsigned long long* buffer = ids.count > 0 ? calloc(ids.count, sizeof(unsigned long long)) : NULL;
    for (NSUInteger index = 0; index < ids.count && buffer != NULL; index++) {
        buffer[index] = ids[index].unsignedLongLongValue;
    }
    processMacSidebarSelectionChanged(self.paneID, buffer, (int)ids.count);
    free(buffer);
}

- (void)expandNodes:(NSArray<WailsSidebarNode*>*)nodes {
    for (WailsSidebarNode* node in nodes) {
        if (node.hidden) continue;
        if (node.section || node.expanded) {
            [self.outlineView expandItem:node];
            [self expandNodes:node.children];
        }
    }
}

- (void)outlineViewItemDidExpand:(NSNotification*)notification {
    WailsSidebarNode* node = notification.userInfo[@"NSObject"];
    if (node == nil || node.section) return;
    // Nested rows keep their own state; restore it now that they are visible.
    BOOL suppressed = self.suppressExpansionCallback;
    self.suppressExpansionCallback = YES;
    [self expandNodes:node.children];
    self.suppressExpansionCallback = suppressed;
    if (suppressed) return;
    node.expanded = YES;
    processMacSidebarItemExpanded(node.nodeID, true);
}

- (void)outlineViewItemDidCollapse:(NSNotification*)notification {
    WailsSidebarNode* node = notification.userInfo[@"NSObject"];
    if (node == nil || node.section || self.suppressExpansionCallback) return;
    node.expanded = NO;
    processMacSidebarItemExpanded(node.nodeID, false);
}

- (void)reloadContents {
    if (self.outlineView == nil) return;
    self.suppressSelectionCallback = YES;
    self.suppressExpansionCallback = YES;
    CGFloat availableWidth = self.outlineView.bounds.size.width;
    if (availableWidth > 0) self.outlineView.outlineTableColumn.width = availableWidth;
    [self.outlineView reloadData];
    [self expandNodes:self.roots];
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
    self.suppressExpansionCallback = NO;
    self.suppressSelectionCallback = NO;
}

// Inline rename.

- (BOOL)beginEditingRow:(NSInteger)row {
    if (row < 0) return NO;
    WailsSidebarNode* node = [self.outlineView itemAtRow:row];
    if (node == nil || node.section || node.disabled || !node.editable) return NO;
    [self.outlineView editColumn:0 row:row withEvent:nil select:YES];
    return YES;
}

- (BOOL)beginEditingSelectedRow {
    if (self.outlineView.numberOfSelectedRows != 1) return NO;
    return [self beginEditingRow:self.outlineView.selectedRow];
}

- (void)handleDoubleClick:(id)sender {
    [self beginEditingRow:self.outlineView.clickedRow];
}

- (void)controlTextDidEndEditing:(NSNotification*)notification {
    NSTextField* field = notification.object;
    if (![field isKindOfClass:[NSTextField class]]) return;
    NSInteger row = [self.outlineView rowForView:field];
    if (row < 0) return;
    WailsSidebarNode* node = [self.outlineView itemAtRow:row];
    if (node == nil || node.section || !node.editable) return;
    NSString* label = field.stringValue ?: @"";
    if ([label isEqualToString:node.label ?: @""]) return;
    node.label = label;
    processMacSidebarItemRenamed(node.nodeID, (char*)label.UTF8String);
}

// Context menus.

- (NSMenu*)contextMenuForEvent:(NSEvent*)event {
    NSPoint point = [self.outlineView convertPoint:event.locationInWindow fromView:nil];
    NSInteger row = [self.outlineView rowAtPoint:point];
    unsigned long long itemID = 0;
    if (row >= 0) {
        WailsSidebarNode* node = [self.outlineView itemAtRow:row];
        if (node != nil && !node.section) itemID = node.nodeID;
    }
    return (NSMenu*)processMacSidebarContextMenu(self.paneID, itemID);
}

// Drag reorder. Only rows from this outline are accepted, and a row keeps its
// nesting level: nested rows reorder beneath their parent, other rows move
// between sections and the root.

- (id<NSPasteboardWriting>)outlineView:(NSOutlineView*)outlineView pasteboardWriterForItem:(id)item {
    WailsSidebarNode* node = (WailsSidebarNode*)item;
    if (!self.reorderable || node == nil || node.section || node.disabled) return nil;
    NSPasteboardItem* pasteboardItem = [[[NSPasteboardItem alloc] init] autorelease];
    [pasteboardItem setString:[NSString stringWithFormat:@"%llu", node.nodeID] forType:WailsSidebarDragType];
    return pasteboardItem;
}

- (WailsSidebarNode*)draggedNodeForInfo:(id<NSDraggingInfo>)info {
    NSString* value = [info.draggingPasteboard stringForType:WailsSidebarDragType];
    if (value == nil) return nil;
    return sidebarNodeInTree(self.roots, strtoull(value.UTF8String, NULL, 10));
}

- (BOOL)canMoveNode:(WailsSidebarNode*)node toParent:(WailsSidebarNode*)parent {
    if (node == nil || node.section) return NO;
    if (parent != nil && sidebarNodeHasAncestor(parent, node)) return NO;
    WailsSidebarNode* current = node.parent;
    if (current != nil && !current.section) return parent == current;
    return parent == nil || parent.section;
}

- (NSDragOperation)outlineView:(NSOutlineView*)outlineView validateDrop:(id<NSDraggingInfo>)info
    proposedItem:(id)item proposedChildIndex:(NSInteger)index {
    if (!self.reorderable || info.draggingSource != outlineView) return NSDragOperationNone;
    WailsSidebarNode* dragged = [self draggedNodeForInfo:info];
    WailsSidebarNode* target = (WailsSidebarNode*)item;
    if (index == NSOutlineViewDropOnItemIndex) {
        if (target == nil || !target.section) return NSDragOperationNone;
        [outlineView setDropItem:target dropChildIndex:(NSInteger)[self visibleNodes:target.children].count];
    }
    if (![self canMoveNode:dragged toParent:target]) return NSDragOperationNone;
    return NSDragOperationMove;
}

- (BOOL)outlineView:(NSOutlineView*)outlineView acceptDrop:(id<NSDraggingInfo>)info
    item:(id)item childIndex:(NSInteger)index {
    if (!self.reorderable || info.draggingSource != outlineView) return NO;
    WailsSidebarNode* dragged = [self draggedNodeForInfo:info];
    WailsSidebarNode* target = (WailsSidebarNode*)item;
    if (![self canMoveNode:dragged toParent:target]) return NO;
    NSArray<WailsSidebarNode*>* all = target == nil ? self.roots : target.children;
    NSArray<WailsSidebarNode*>* visible = [self visibleNodes:all];
    NSInteger fullIndex;
    if (index == NSOutlineViewDropOnItemIndex || index >= (NSInteger)visible.count) {
        fullIndex = (NSInteger)all.count;
    } else {
        fullIndex = (NSInteger)[all indexOfObjectIdenticalTo:visible[(NSUInteger)index]];
    }
    processMacSidebarItemMoved(self.paneID, dragged.nodeID, target == nil ? 0 : target.nodeID, (int)fullIndex);
    return YES;
}

@end

@implementation WailsSidebarOutlineView

- (NSMenu*)menuForEvent:(NSEvent*)event {
    // Let NSTableView record clickedRow and draw the contextual highlight.
    NSMenu* fallback = [super menuForEvent:event];
    NSMenu* menu = [self.controller contextMenuForEvent:event];
    return menu != nil ? menu : fallback;
}

- (void)keyDown:(NSEvent*)event {
    NSString* characters = event.charactersIgnoringModifiers;
    if (characters.length == 1) {
        unichar key = [characters characterAtIndex:0];
        if ((key == NSCarriageReturnCharacter || key == NSEnterCharacter) &&
            [self.controller beginEditingSelectedRow]) {
            return;
        }
    }
    [super keyDown:event];
}

@end

static const void* WailsInspectorControlIDAssociationKey = &WailsInspectorControlIDAssociationKey;
static const void* WailsInspectorStepperFieldIDAssociationKey = &WailsInspectorStepperFieldIDAssociationKey;
static const void* WailsInspectorSectionIDAssociationKey = &WailsInspectorSectionIDAssociationKey;

// Mirrors MacInspectorControlKind.
enum {
    WailsInspectorKindLabel = 0,
    WailsInspectorKindTextField = 1,
    WailsInspectorKindCheckbox = 2,
    WailsInspectorKindPopup = 3,
    WailsInspectorKindSlider = 4,
    WailsInspectorKindStepper = 5,
    WailsInspectorKindSegmented = 6,
    WailsInspectorKindColorWell = 7,
    WailsInspectorKindDatePicker = 8,
    WailsInspectorKindButton = 9,
};

@interface WailsInspectorControlModel : NSObject
@property unsigned long long controlID;
@property int kind;
@property (copy) NSString* label;
@property (copy) NSString* value;
@property BOOL checked;
@property (retain) NSArray<NSString*>* options;
@property NSInteger selectedIndex;
@property double number;
@property double minimum;
@property double maximum;
@property double step;
@property (retain) NSColor* color;
@property double dateSeconds;
@property (copy) NSString* tooltip;
@property BOOL disabled;
@property BOOL hidden;
@end

@implementation WailsInspectorControlModel
- (void)dealloc {
    [_label release];
    [_value release];
    [_options release];
    [_color release];
    [_tooltip release];
    [super dealloc];
}
@end

@interface WailsInspectorSectionModel : NSObject
@property unsigned long long sectionID;
@property (copy) NSString* label;
@property BOOL collapsible;
@property BOOL collapsed;
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

static NSString* inspectorNumberString(double value) {
    return [NSString stringWithFormat:@"%g", value];
}

static NSStepper* inspectorStepperInContainer(NSView* container) {
    for (NSView* view in container.subviews) {
        if ([view isKindOfClass:[NSStepper class]]) return (NSStepper*)view;
    }
    return nil;
}

static NSTextField* inspectorFieldInContainer(NSView* container) {
    for (NSView* view in container.subviews) {
        if ([view isKindOfClass:[NSTextField class]]) return (NSTextField*)view;
    }
    return nil;
}

// A native property inspector. Its scroll view, section headings, and every
// control (labels, text fields, checkboxes, pop-ups, sliders, steppers,
// segmented controls, colour wells, date pickers, and buttons) are AppKit.
// The view remains transparent so NSSplitViewItem's semantic inspector
// surface controls the appearance.
@interface WailsInspectorViewController : NSViewController <NSTextFieldDelegate>
@property (retain) NSMutableArray<WailsInspectorSectionModel*>* sections;
@property (retain) NSMutableDictionary<NSNumber*, WailsInspectorControlModel*>* modelsByID;
@property (retain) NSMutableDictionary<NSNumber*, NSView*>* controlsByID;
@property (retain) NSMutableDictionary<NSNumber*, NSView*>* rowsByID;
@property (retain) NSMutableDictionary<NSNumber*, NSTextField*>* nameLabelsByID;
@property (retain) NSMutableDictionary<NSNumber*, NSNumber*>* sectionIDsByControlID;
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
        _sectionIDsByControlID = [[NSMutableDictionary alloc] init];
    }
    return self;
}

- (void)dealloc {
    [_sections release];
    [_modelsByID release];
    [_controlsByID release];
    [_rowsByID release];
    [_nameLabelsByID release];
    [_sectionIDsByControlID release];
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

- (WailsInspectorSectionModel*)sectionForID:(NSNumber*)sectionID {
    if (sectionID == nil) return nil;
    for (WailsInspectorSectionModel* section in self.sections) {
        if (section.sectionID == sectionID.unsignedLongLongValue) return section;
    }
    return nil;
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

- (NSView*)nativeControlForModel:(WailsInspectorControlModel*)model {
    NSView* result = nil;
    switch (model.kind) {
        case WailsInspectorKindLabel: {
            NSTextField* value = [NSTextField labelWithString:model.value ?: @""];
            value.selectable = YES;
            value.lineBreakMode = NSLineBreakByTruncatingTail;
            result = value;
            break;
        }
        case WailsInspectorKindTextField: {
            NSTextField* field = [NSTextField textFieldWithString:model.value ?: @""];
            field.controlSize = NSControlSizeSmall;
            field.font = [NSFont systemFontOfSize:[NSFont smallSystemFontSize]];
            field.delegate = self;
            objc_setAssociatedObject(field, WailsInspectorControlIDAssociationKey,
                @(model.controlID), OBJC_ASSOCIATION_RETAIN);
            result = field;
            break;
        }
        case WailsInspectorKindCheckbox: {
            NSButton* checkbox = [NSButton checkboxWithTitle:model.label ?: @"" target:self
                action:@selector(handleCheckbox:)];
            checkbox.controlSize = NSControlSizeSmall;
            result = checkbox;
            break;
        }
        case WailsInspectorKindPopup: {
            NSPopUpButton* popup = [[[NSPopUpButton alloc] initWithFrame:NSZeroRect pullsDown:NO] autorelease];
            popup.controlSize = NSControlSizeSmall;
            [popup addItemsWithTitles:model.options ?: @[]];
            popup.target = self;
            popup.action = @selector(handlePopup:);
            result = popup;
            break;
        }
        case WailsInspectorKindSlider: {
            NSSlider* slider = [NSSlider sliderWithValue:model.number minValue:model.minimum
                maxValue:model.maximum target:self action:@selector(handleSlider:)];
            slider.controlSize = NSControlSizeSmall;
            slider.continuous = YES;
            result = slider;
            break;
        }
        case WailsInspectorKindStepper: {
            NSTextField* field = [NSTextField textFieldWithString:inspectorNumberString(model.number)];
            field.controlSize = NSControlSizeSmall;
            field.font = [NSFont monospacedDigitSystemFontOfSize:[NSFont smallSystemFontSize]
                weight:NSFontWeightRegular];
            field.alignment = NSTextAlignmentRight;
            field.delegate = self;
            field.translatesAutoresizingMaskIntoConstraints = NO;
            [field.widthAnchor constraintEqualToConstant:64].active = YES;
            objc_setAssociatedObject(field, WailsInspectorStepperFieldIDAssociationKey,
                @(model.controlID), OBJC_ASSOCIATION_RETAIN);

            NSStepper* stepper = [[[NSStepper alloc] initWithFrame:NSZeroRect] autorelease];
            stepper.controlSize = NSControlSizeSmall;
            stepper.minValue = model.minimum;
            stepper.maxValue = model.maximum;
            stepper.increment = model.step;
            stepper.doubleValue = model.number;
            stepper.valueWraps = NO;
            stepper.autorepeat = YES;
            stepper.target = self;
            stepper.action = @selector(handleStepper:);

            NSStackView* container = [NSStackView stackViewWithViews:@[field, stepper]];
            container.orientation = NSUserInterfaceLayoutOrientationHorizontal;
            container.alignment = NSLayoutAttributeCenterY;
            container.spacing = 2;
            result = container;
            break;
        }
        case WailsInspectorKindSegmented: {
            NSSegmentedControl* segmented = [NSSegmentedControl segmentedControlWithLabels:model.options ?: @[]
                trackingMode:NSSegmentSwitchTrackingSelectOne target:self action:@selector(handleSegmented:)];
            segmented.controlSize = NSControlSizeSmall;
            segmented.segmentDistribution = NSSegmentDistributionFillEqually;
            result = segmented;
            break;
        }
        case WailsInspectorKindColorWell: {
            NSColorWell* well = [[[NSColorWell alloc] initWithFrame:NSZeroRect] autorelease];
            well.translatesAutoresizingMaskIntoConstraints = NO;
            [NSLayoutConstraint activateConstraints:@[
                [well.widthAnchor constraintEqualToConstant:44],
                [well.heightAnchor constraintEqualToConstant:23]
            ]];
            well.target = self;
            well.action = @selector(handleColorWell:);
            well.continuous = YES;
            result = well;
            break;
        }
        case WailsInspectorKindDatePicker: {
            NSDatePicker* picker = [[[NSDatePicker alloc] initWithFrame:NSZeroRect] autorelease];
            picker.controlSize = NSControlSizeSmall;
            picker.font = [NSFont systemFontOfSize:[NSFont smallSystemFontSize]];
            picker.datePickerStyle = NSDatePickerStyleTextFieldAndStepper;
            picker.datePickerMode = NSDatePickerModeSingle;
            picker.datePickerElements = NSDatePickerElementFlagYearMonthDay | NSDatePickerElementFlagHourMinute;
            picker.drawsBackground = NO;
            picker.bezeled = NO;
            picker.bordered = NO;
            picker.target = self;
            picker.action = @selector(handleDatePicker:);
            result = picker;
            break;
        }
        case WailsInspectorKindButton: {
            NSButton* button = [NSButton buttonWithTitle:model.label ?: @"" target:self
                action:@selector(handleButton:)];
            button.controlSize = NSControlSizeSmall;
            button.bezelStyle = NSBezelStyleRounded;
            result = button;
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

- (BOOL)modelOmitsNameLabel:(WailsInspectorControlModel*)model {
    return model.kind == WailsInspectorKindCheckbox || model.kind == WailsInspectorKindButton;
}

- (NSView*)rowForModel:(WailsInspectorControlModel*)model control:(NSView*)control {
    NSStackView* row = [[NSStackView alloc] initWithFrame:NSZeroRect];
    row.translatesAutoresizingMaskIntoConstraints = NO;
    row.orientation = NSUserInterfaceLayoutOrientationHorizontal;
    BOOL textual = model.kind == WailsInspectorKindLabel || model.kind == WailsInspectorKindTextField ||
        model.kind == WailsInspectorKindPopup;
    row.alignment = textual ? NSLayoutAttributeFirstBaseline : NSLayoutAttributeCenterY;
    row.distribution = NSStackViewDistributionFill;
    row.spacing = 8;

    if ([self modelOmitsNameLabel:model]) {
        [row addArrangedSubview:control];
    } else {
        NSTextField* nameLabel = [self propertyNameLabel:model.label];
        [row addArrangedSubview:nameLabel];
        self.nameLabelsByID[@(model.controlID)] = nameLabel;
        [row addArrangedSubview:control];
        BOOL fixedWidth = model.kind == WailsInspectorKindColorWell || model.kind == WailsInspectorKindStepper ||
            model.kind == WailsInspectorKindDatePicker;
        if (!fixedWidth) {
            [control setContentCompressionResistancePriority:NSLayoutPriorityDefaultLow
                forOrientation:NSLayoutConstraintOrientationHorizontal];
            [control.widthAnchor constraintGreaterThanOrEqualToConstant:90].active = YES;
        } else {
            [control setContentHuggingPriority:NSLayoutPriorityRequired
                forOrientation:NSLayoutConstraintOrientationHorizontal];
        }
    }
    row.hidden = model.hidden;
    return [row autorelease];
}

- (NSView*)headingForSection:(WailsInspectorSectionModel*)section {
    NSTextField* heading = [NSTextField labelWithString:section.label ?: @""];
    heading.font = [NSFont systemFontOfSize:13 weight:NSFontWeightSemibold];
    heading.textColor = [NSColor labelColor];
    heading.lineBreakMode = NSLineBreakByTruncatingTail;
    heading.translatesAutoresizingMaskIntoConstraints = NO;
    if (!section.collapsible) return heading;

    NSButton* disclosure = [[[NSButton alloc] initWithFrame:NSZeroRect] autorelease];
    [disclosure setButtonType:NSButtonTypeOnOff];
    disclosure.bezelStyle = NSBezelStyleDisclosure;
    disclosure.title = @"";
    disclosure.controlSize = NSControlSizeSmall;
    disclosure.state = section.collapsed ? NSControlStateValueOff : NSControlStateValueOn;
    disclosure.target = self;
    disclosure.action = @selector(handleDisclosure:);
    disclosure.translatesAutoresizingMaskIntoConstraints = NO;
    objc_setAssociatedObject(disclosure, WailsInspectorSectionIDAssociationKey,
        @(section.sectionID), OBJC_ASSOCIATION_RETAIN);

    NSStackView* header = [NSStackView stackViewWithViews:@[disclosure, heading]];
    header.translatesAutoresizingMaskIntoConstraints = NO;
    header.orientation = NSUserInterfaceLayoutOrientationHorizontal;
    header.alignment = NSLayoutAttributeCenterY;
    header.spacing = 2;
    return header;
}

- (void)applySectionVisibility:(WailsInspectorSectionModel*)section {
    for (WailsInspectorControlModel* model in section.controls) {
        NSView* row = self.rowsByID[@(model.controlID)];
        row.hidden = model.hidden || section.collapsed;
    }
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
    [self.sectionIDsByControlID removeAllObjects];

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

        NSView* heading = [self headingForSection:section];
        [self.stackView addArrangedSubview:heading];
        [heading.widthAnchor constraintEqualToAnchor:self.stackView.widthAnchor].active = YES;
        [self.stackView setCustomSpacing:9 afterView:heading];

        for (WailsInspectorControlModel* model in section.controls) {
            NSView* control = [self nativeControlForModel:model];
            if (control == nil) continue;
            NSView* row = [self rowForModel:model control:control];
            [self.stackView addArrangedSubview:row];
            [row.widthAnchor constraintEqualToAnchor:self.stackView.widthAnchor].active = YES;
            self.controlsByID[@(model.controlID)] = control;
            self.rowsByID[@(model.controlID)] = row;
            self.sectionIDsByControlID[@(model.controlID)] = @(section.sectionID);
            [self applyModel:model];
        }
    }
}

- (void)applyModel:(WailsInspectorControlModel*)model {
    NSView* control = self.controlsByID[@(model.controlID)];
    NSView* row = self.rowsByID[@(model.controlID)];
    if (control == nil) return;

    WailsInspectorSectionModel* section = [self sectionForID:self.sectionIDsByControlID[@(model.controlID)]];
    row.hidden = model.hidden || (section != nil && section.collapsed);
    control.toolTip = model.tooltip.length > 0 ? model.tooltip : nil;
    if ([control isKindOfClass:[NSControl class]]) ((NSControl*)control).enabled = !model.disabled;
    NSTextField* nameLabel = self.nameLabelsByID[@(model.controlID)];
    if (nameLabel != nil) nameLabel.stringValue = model.label ?: @"";

    switch (model.kind) {
        case WailsInspectorKindLabel:
        case WailsInspectorKindTextField:
            ((NSTextField*)control).stringValue = model.value ?: @"";
            break;
        case WailsInspectorKindCheckbox:
            ((NSButton*)control).title = model.label ?: @"";
            ((NSButton*)control).state = model.checked ? NSControlStateValueOn : NSControlStateValueOff;
            break;
        case WailsInspectorKindPopup: {
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
        case WailsInspectorKindSlider: {
            NSSlider* slider = (NSSlider*)control;
            slider.minValue = model.minimum;
            slider.maxValue = model.maximum;
            slider.doubleValue = model.number;
            break;
        }
        case WailsInspectorKindStepper: {
            NSStepper* stepper = inspectorStepperInContainer(control);
            NSTextField* field = inspectorFieldInContainer(control);
            stepper.minValue = model.minimum;
            stepper.maxValue = model.maximum;
            stepper.increment = model.step;
            stepper.doubleValue = model.number;
            stepper.enabled = !model.disabled;
            field.stringValue = inspectorNumberString(model.number);
            field.enabled = !model.disabled;
            break;
        }
        case WailsInspectorKindSegmented: {
            NSSegmentedControl* segmented = (NSSegmentedControl*)control;
            NSArray<NSString*>* options = model.options ?: @[];
            BOOL same = segmented.segmentCount == (NSInteger)options.count;
            for (NSInteger index = 0; same && index < segmented.segmentCount; index++) {
                same = [[segmented labelForSegment:index] isEqualToString:options[(NSUInteger)index]];
            }
            if (!same) {
                segmented.segmentCount = (NSInteger)options.count;
                for (NSInteger index = 0; index < segmented.segmentCount; index++) {
                    [segmented setLabel:options[(NSUInteger)index] forSegment:index];
                }
            }
            if (model.selectedIndex >= 0 && model.selectedIndex < segmented.segmentCount) {
                segmented.selectedSegment = model.selectedIndex;
            } else {
                segmented.selectedSegment = -1;
            }
            break;
        }
        case WailsInspectorKindColorWell:
            if (model.color != nil) ((NSColorWell*)control).color = model.color;
            break;
        case WailsInspectorKindDatePicker:
            ((NSDatePicker*)control).dateValue = [NSDate dateWithTimeIntervalSince1970:model.dateSeconds];
            break;
        case WailsInspectorKindButton:
            ((NSButton*)control).title = model.label ?: @"";
            break;
    }
}

- (void)controlTextDidChange:(NSNotification*)notification {
    NSTextField* field = notification.object;
    NSNumber* controlID = objc_getAssociatedObject(field, WailsInspectorControlIDAssociationKey);
    if (controlID != nil && [field isKindOfClass:[NSTextField class]]) {
        processMacInspectorTextChanged(controlID.unsignedLongLongValue, (char*)field.stringValue.UTF8String);
    }
}

- (void)controlTextDidEndEditing:(NSNotification*)notification {
    NSTextField* field = notification.object;
    NSNumber* controlID = objc_getAssociatedObject(field, WailsInspectorStepperFieldIDAssociationKey);
    if (controlID == nil) return;
    WailsInspectorControlModel* model = self.modelsByID[controlID];
    NSStepper* stepper = inspectorStepperInContainer(self.controlsByID[controlID]);
    if (model == nil || stepper == nil) return;
    double value = MIN(MAX(field.doubleValue, model.minimum), model.maximum);
    model.number = value;
    stepper.doubleValue = value;
    field.stringValue = inspectorNumberString(value);
    processMacInspectorNumberChanged(controlID.unsignedLongLongValue, WailsInspectorKindStepper, value);
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

- (void)handleSlider:(NSSlider*)sender {
    NSNumber* controlID = objc_getAssociatedObject(sender, WailsInspectorControlIDAssociationKey);
    if (controlID == nil) return;
    self.modelsByID[controlID].number = sender.doubleValue;
    processMacInspectorNumberChanged(controlID.unsignedLongLongValue, WailsInspectorKindSlider, sender.doubleValue);
}

- (void)handleStepper:(NSStepper*)sender {
    NSNumber* controlID = objc_getAssociatedObject(sender.superview, WailsInspectorControlIDAssociationKey);
    if (controlID == nil) return;
    WailsInspectorControlModel* model = self.modelsByID[controlID];
    model.number = sender.doubleValue;
    inspectorFieldInContainer(sender.superview).stringValue = inspectorNumberString(sender.doubleValue);
    processMacInspectorNumberChanged(controlID.unsignedLongLongValue, WailsInspectorKindStepper, sender.doubleValue);
}

- (void)handleSegmented:(NSSegmentedControl*)sender {
    NSNumber* controlID = objc_getAssociatedObject(sender, WailsInspectorControlIDAssociationKey);
    if (controlID != nil) {
        processMacInspectorSegmentChanged(controlID.unsignedLongLongValue, (int)sender.selectedSegment);
    }
}

- (void)handleColorWell:(NSColorWell*)sender {
    NSNumber* controlID = objc_getAssociatedObject(sender, WailsInspectorControlIDAssociationKey);
    NSColor* color = [sender.color colorUsingColorSpace:[NSColorSpace sRGBColorSpace]];
    if (controlID == nil || color == nil) return;
    self.modelsByID[controlID].color = sender.color;
    processMacInspectorColorChanged(controlID.unsignedLongLongValue,
        (int)lround(color.redComponent * 255.0), (int)lround(color.greenComponent * 255.0),
        (int)lround(color.blueComponent * 255.0), (int)lround(color.alphaComponent * 255.0));
}

- (void)handleDatePicker:(NSDatePicker*)sender {
    NSNumber* controlID = objc_getAssociatedObject(sender, WailsInspectorControlIDAssociationKey);
    if (controlID == nil) return;
    double seconds = sender.dateValue.timeIntervalSince1970;
    self.modelsByID[controlID].dateSeconds = seconds;
    processMacInspectorDateChanged(controlID.unsignedLongLongValue, seconds);
}

- (void)handleButton:(NSButton*)sender {
    NSNumber* controlID = objc_getAssociatedObject(sender, WailsInspectorControlIDAssociationKey);
    if (controlID != nil) processMacInspectorButtonClicked(controlID.unsignedLongLongValue);
}

- (void)handleDisclosure:(NSButton*)sender {
    NSNumber* sectionID = objc_getAssociatedObject(sender, WailsInspectorSectionIDAssociationKey);
    WailsInspectorSectionModel* section = [self sectionForID:sectionID];
    if (section == nil) return;
    section.collapsed = sender.state == NSControlStateValueOff;
    [self applySectionVisibility:section];
    processMacInspectorSectionCollapsed(section.sectionID, section.collapsed);
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
@property BOOL sidebarAllowsMultipleSelection;
@property BOOL sidebarReorderable;
@property (retain) NSMutableArray<WailsSidebarNode*>* sidebarRoots;
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

// sidebarControllerApplyRecord copies the pane's sidebar options onto a
// freshly created controller before its view loads.
static void sidebarControllerApplyRecord(WailsSidebarViewController* sidebar, WailsSplitPaneRecord* record) {
    sidebar.paneID = record.paneID;
    sidebar.allowsMultipleSelection = record.sidebarAllowsMultipleSelection;
    sidebar.reorderable = record.sidebarReorderable;
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

static void configureInspectorModel(WailsInspectorControlModel* model, WailsInspectorControlSpec spec) {
    if (model == nil) return;
    model.kind = spec.kind;
    model.label = spec.label == NULL ? @"" : [NSString stringWithUTF8String:spec.label];
    model.value = spec.value == NULL ? @"" : [NSString stringWithUTF8String:spec.value];
    model.checked = spec.checked;
    model.options = inspectorOptionsFromJSON(spec.optionsJSON);
    model.selectedIndex = spec.selectedIndex;
    model.number = spec.number;
    model.minimum = spec.minimum;
    model.maximum = spec.maximum;
    model.step = spec.step;
    model.color = [NSColor colorWithSRGBRed:spec.red / 255.0 green:spec.green / 255.0
        blue:spec.blue / 255.0 alpha:spec.alpha / 255.0];
    model.dateSeconds = spec.dateSeconds;
    model.tooltip = spec.tooltip == NULL ? @"" : [NSString stringWithUTF8String:spec.tooltip];
    model.disabled = spec.disabled;
    model.hidden = spec.hidden;
}

static void splitViewItemApplyCanCollapseFromResize(NSSplitViewItem* item, BOOL allowed) {
    if (item != nil && [item respondsToSelector:@selector(setCanCollapseFromWindowResize:)]) {
        [item setValue:[NSNumber numberWithBool:allowed] forKey:@"canCollapseFromWindowResize"];
    }
}

// Content list hooks begin (table implemented in webview_window_contentlist_darwin.m).
// The table controller is created lazily so Go can stage rows before the
// window exists; the record's generic viewController retains it.
NSViewController* splitViewContentListController(void* handlePtr, unsigned long long paneID) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record == nil || record.role != WailsSplitPaneRoleContentList) return nil;
    if (record.viewController == nil) record.viewController = wailsContentListCreateController(paneID);
    return record.viewController;
}

static BOOL splitRecordPrepareContentList(void* handlePtr, WailsSplitPaneRecord* record, NSColor* surfaceColor) {
    NSViewController* controller = splitViewContentListController(handlePtr, record.paneID);
    return controller != nil && wailsContentListPrepareForInstall(controller, surfaceColor);
}

static void splitRecordsReloadContentLists(NSArray<WailsSplitPaneRecord*>* records) {
    for (WailsSplitPaneRecord* record in records) {
        if (record.role == WailsSplitPaneRoleContentList) wailsContentListDidInstall(record.viewController);
    }
}
// Content list hooks end.

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
            record.role != WailsSplitPaneRoleInspector &&
            record.role != WailsSplitPaneRoleContentList) return false;
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
    wailsPrivateSetWebviewTransparent(primaryWebView);

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
            sidebar.selectedItemID = record.selectedSidebarItemID;
            sidebarControllerApplyRecord(sidebar, record);
            sidebar.surfaceColor = primaryBackground;
            (void)sidebar.view;
            if (sidebar.view == nil) {
                [sidebar release];
                return false;
            }
            record.sidebarController = sidebar;
            record.viewController = sidebar;
            [sidebar release];
        } else if (record.role == WailsSplitPaneRoleContentList) {
            if (!splitRecordPrepareContentList(handlePtr, record, primaryBackground)) return false;
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
        } else if (record.role == WailsSplitPaneRoleContentList) {
            item = wailsContentListSplitItem(record.viewController);
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
    splitRecordsReloadContentLists(owner.records);

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
            record.role != WailsSplitPaneRoleInspector &&
            record.role != WailsSplitPaneRoleContentList) return false;
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
            sidebar.selectedItemID = record.selectedSidebarItemID;
            sidebarControllerApplyRecord(sidebar, record);
            sidebar.surfaceColor = primaryBackground;
            (void)sidebar.view;
            if (sidebar.view == nil) {
                [sidebar release];
                return false;
            }
            record.sidebarController = sidebar;
            record.viewController = sidebar;
            [sidebar release];
        } else if (record.role == WailsSplitPaneRoleContentList) {
            if (!splitRecordPrepareContentList(handlePtr, record, primaryBackground)) return false;
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
        } else if (record.role == WailsSplitPaneRoleContentList) {
            item = wailsContentListSplitItem(record.viewController);
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
    splitRecordsReloadContentLists(owner.records);
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

// Pane accessory hook: see webview_window_split_darwin.h.
void* splitViewPaneItem(void* handlePtr, unsigned long long paneID) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    return record == nil ? NULL : (void*)record.item;
}

void splitViewSidebarReset(void* handlePtr, unsigned long long paneID) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record == nil || record.role != WailsSplitPaneRoleSidebar) return;
    [record.sidebarRoots removeAllObjects];
    [record.sidebarController reloadContents];
}

void splitViewSidebarSetOptions(void* handlePtr, unsigned long long paneID,
    bool allowsMultipleSelection, bool reorderable) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record == nil || record.role != WailsSplitPaneRoleSidebar) return;
    record.sidebarAllowsMultipleSelection = allowsMultipleSelection;
    record.sidebarReorderable = reorderable;
    if (record.sidebarController != nil) {
        record.sidebarController.allowsMultipleSelection = allowsMultipleSelection;
        record.sidebarController.reorderable = reorderable;
        [record.sidebarController applyOptions];
    }
}

void splitViewSidebarAddSection(void* handlePtr, unsigned long long paneID,
    unsigned long long sectionID, const char* label) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record == nil || record.role != WailsSplitPaneRoleSidebar) return;
    WailsSidebarNode* node = [[WailsSidebarNode alloc] init];
    node.nodeID = sectionID;
    node.section = YES;
    node.label = label == NULL ? @"" : [NSString stringWithUTF8String:label];
    [record.sidebarRoots addObject:node];
    [node release];
}

void splitViewSidebarAddItem(void* handlePtr, unsigned long long paneID,
    unsigned long long parentID, unsigned long long itemID,
    const char* label, const char* symbolName, const char* tooltip,
    bool disabled, bool hidden, bool expanded, bool editable, int badge,
    const char* accessorySymbol, bool hasTint, int red, int green, int blue, int alpha) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record == nil || record.role != WailsSplitPaneRoleSidebar) return;
    WailsSidebarNode* node = [[WailsSidebarNode alloc] init];
    node.nodeID = itemID;
    node.label = label == NULL ? @"" : [NSString stringWithUTF8String:label];
    node.symbolName = symbolName == NULL ? @"" : [NSString stringWithUTF8String:symbolName];
    node.tooltip = tooltip == NULL ? @"" : [NSString stringWithUTF8String:tooltip];
    node.accessorySymbol = accessorySymbol == NULL ? @"" : [NSString stringWithUTF8String:accessorySymbol];
    node.disabled = disabled;
    node.hidden = hidden;
    node.expanded = expanded;
    node.editable = editable;
    node.badge = badge;
    node.tintColor = hasTint ? [NSColor colorWithSRGBRed:red / 255.0 green:green / 255.0
        blue:blue / 255.0 alpha:alpha / 255.0] : nil;
    WailsSidebarNode* parent = parentID == 0 ? nil : sidebarNodeInTree(record.sidebarRoots, parentID);
    node.parent = parent;
    if (parent != nil) [parent.children addObject:node];
    else [record.sidebarRoots addObject:node];
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
    unsigned long long sectionID, const char* label, bool collapsible, bool collapsed) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record == nil || record.role != WailsSplitPaneRoleInspector) return;
    WailsInspectorSectionModel* section = [[WailsInspectorSectionModel alloc] init];
    section.sectionID = sectionID;
    section.label = label == NULL ? @"" : [NSString stringWithUTF8String:label];
    section.collapsible = collapsible;
    section.collapsed = collapsible && collapsed;
    section.controls = [NSMutableArray array];
    [record.inspectorSections addObject:section];
    [section release];
}

void splitViewInspectorAddControl(void* handlePtr, unsigned long long paneID,
    unsigned long long sectionID, unsigned long long controlID,
    WailsInspectorControlSpec spec) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    WailsInspectorSectionModel* section = inspectorSection(record, sectionID);
    if (record == nil || record.role != WailsSplitPaneRoleInspector || section == nil) return;
    WailsInspectorControlModel* model = [[WailsInspectorControlModel alloc] init];
    model.controlID = controlID;
    configureInspectorModel(model, spec);
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
    unsigned long long controlID, WailsInspectorControlSpec spec) {
    WailsSplitPaneRecord* record = splitPaneRecord(handlePtr, paneID);
    if (record == nil || record.role != WailsSplitPaneRoleInspector) return;
    WailsInspectorControlModel* model = record.inspectorModelsByID[@(controlID)];
    if (model == nil || model.kind != spec.kind) return;
    configureInspectorModel(model, spec);
    [record.inspectorController applyModel:model];
}
