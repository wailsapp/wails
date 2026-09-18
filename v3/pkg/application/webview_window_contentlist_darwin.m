//go:build darwin && !ios && !server

#import "webview_window_contentlist_darwin.h"
#import <string.h>
#import <stdlib.h>

// Content-list model. Rows and columns mirror the Go MacContentList; the
// controller keeps a filtered visibleRows array for the table.

@interface WailsContentListRowModel : NSObject
@property unsigned long long rowID;
@property (copy) NSString* title;
@property (copy) NSString* subtitle;
@property (copy) NSString* detail;
@property (copy) NSString* symbol;
@property (copy) NSString* tooltip;
@property int badge;
@property (retain) NSArray<NSString*>* cells;
@property BOOL disabled;
@property BOOL hidden;
@end

@implementation WailsContentListRowModel
- (void)dealloc {
    [_title release];
    [_subtitle release];
    [_detail release];
    [_symbol release];
    [_tooltip release];
    [_cells release];
    [super dealloc];
}
@end

@interface WailsContentListColumnModel : NSObject
@property (copy) NSString* title;
@property double width;
@property double minWidth;
@property double maxWidth;
@property BOOL sortable;
@property int alignment;
@end

@implementation WailsContentListColumnModel
- (void)dealloc {
    [_title release];
    [super dealloc];
}
@end

static NSString* contentListString(const char* value) {
    return value == NULL ? @"" : ([NSString stringWithUTF8String:value] ?: @"");
}

static NSArray<NSString*>* contentListCellsFromJSON(const char* json) {
    if (json == NULL || strlen(json) == 0) return @[];
    NSData* data = [NSData dataWithBytes:json length:strlen(json)];
    id parsed = [NSJSONSerialization JSONObjectWithData:data options:0 error:NULL];
    if (![parsed isKindOfClass:[NSArray class]]) return @[];
    NSMutableArray<NSString*>* result = [NSMutableArray arrayWithCapacity:[parsed count]];
    for (id value in (NSArray*)parsed) {
        [result addObject:[value isKindOfClass:[NSString class]] ? value : @""];
    }
    return result;
}

static void contentListConfigureRow(WailsContentListRowModel* row, WailsContentListRowSpec spec) {
    row.title = contentListString(spec.title);
    row.subtitle = contentListString(spec.subtitle);
    row.detail = contentListString(spec.detail);
    row.symbol = contentListString(spec.symbol);
    row.tooltip = contentListString(spec.tooltip);
    row.badge = spec.badge;
    row.cells = contentListCellsFromJSON(spec.cellsJSON);
    row.disabled = spec.disabled;
    row.hidden = spec.hidden;
}

static NSImage* contentListSymbolImage(NSString* symbolName) {
    NSImage* image = nil;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
    if (@available(macOS 11.0, *)) {
        image = symbolName.length > 0
            ? [NSImage imageWithSystemSymbolName:symbolName accessibilityDescription:symbolName]
            : nil;
    }
#endif
    return image;
}

static NSTextAlignment contentListTextAlignment(int alignment) {
    switch (alignment) {
    case 1: return NSTextAlignmentCenter;
    case 2: return NSTextAlignmentRight;
    default: return NSTextAlignmentNatural;
    }
}

static NSColor* contentListAccentColor(void) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101400
    if (@available(macOS 10.14, *)) return [NSColor controlAccentColor];
#endif
    return [NSColor systemBlueColor];
}

// WailsContentListBadgeView draws a count capsule like Mail's unread badge.
// Drawing resolves colours each time so appearance changes are honoured.
@interface WailsContentListBadgeView : NSView
@property (nonatomic) int count;
@property (nonatomic) BOOL emphasized;
@end

@implementation WailsContentListBadgeView

- (NSDictionary*)textAttributes {
    NSColor* textColor = self.emphasized ? contentListAccentColor() : [NSColor whiteColor];
    return @{
        NSFontAttributeName: [NSFont monospacedDigitSystemFontOfSize:11 weight:NSFontWeightSemibold],
        NSForegroundColorAttributeName: textColor
    };
}

- (NSSize)intrinsicContentSize {
    if (self.count <= 0) return NSMakeSize(0, 16);
    NSString* text = [NSString stringWithFormat:@"%d", self.count];
    NSSize textSize = [text sizeWithAttributes:[self textAttributes]];
    return NSMakeSize(MAX(18, ceil(textSize.width) + 12), 16);
}

- (void)setCount:(int)count {
    _count = count;
    [self invalidateIntrinsicContentSize];
    self.needsDisplay = YES;
}

- (void)setEmphasized:(BOOL)emphasized {
    _emphasized = emphasized;
    self.needsDisplay = YES;
}

- (void)viewDidChangeEffectiveAppearance {
    [super viewDidChangeEffectiveAppearance];
    self.needsDisplay = YES;
}

- (void)drawRect:(NSRect)dirtyRect {
    if (self.count <= 0) return;
    NSRect bounds = self.bounds;
    NSBezierPath* capsule = [NSBezierPath bezierPathWithRoundedRect:bounds
        xRadius:bounds.size.height / 2 yRadius:bounds.size.height / 2];
    [(self.emphasized ? [NSColor whiteColor] : contentListAccentColor()) setFill];
    [capsule fill];
    NSString* text = [NSString stringWithFormat:@"%d", self.count];
    NSDictionary* attributes = [self textAttributes];
    NSSize textSize = [text sizeWithAttributes:attributes];
    NSPoint origin = NSMakePoint(NSMidX(bounds) - textSize.width / 2, NSMidY(bounds) - textSize.height / 2);
    [text drawAtPoint:origin withAttributes:attributes];
}

@end

// One rich row: leading symbol, title with trailing detail text, and a
// subtitle line with a trailing badge. Hidden trailing views leave the
// layout.
@interface WailsContentListRichCell : NSTableCellView
@property (retain) NSTextField* subtitleField;
@property (retain) NSTextField* detailField;
@property (retain) WailsContentListBadgeView* badgeView;
@property BOOL disabledRow;
@end

@implementation WailsContentListRichCell

- (void)dealloc {
    [_subtitleField release];
    [_detailField release];
    [_badgeView release];
    [super dealloc];
}

- (void)applyColors {
    BOOL emphasized = self.backgroundStyle == NSBackgroundStyleEmphasized;
    NSColor* secondary = emphasized
        ? [[NSColor alternateSelectedControlTextColor] colorWithAlphaComponent:0.85]
        : [NSColor secondaryLabelColor];
    if (self.disabledRow) {
        secondary = [NSColor disabledControlTextColor];
        self.textField.textColor = [NSColor disabledControlTextColor];
    } else if (!emphasized) {
        self.textField.textColor = [NSColor labelColor];
    }
    self.subtitleField.textColor = secondary;
    self.detailField.textColor = secondary;
    self.badgeView.emphasized = emphasized;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101400
    if (@available(macOS 10.14, *)) {
        self.imageView.contentTintColor = self.disabledRow ? [NSColor disabledControlTextColor]
            : (emphasized ? [NSColor alternateSelectedControlTextColor] : contentListAccentColor());
    }
#endif
}

- (void)setBackgroundStyle:(NSBackgroundStyle)backgroundStyle {
    [super setBackgroundStyle:backgroundStyle];
    [self applyColors];
}

@end

@class WailsContentListViewController;

// The table subclass routes right-clicks and Return to the controller.
@interface WailsContentListTableView : NSTableView
@property (assign) WailsContentListViewController* controller;
@end

@interface WailsContentListViewController : NSViewController <NSTableViewDataSource, NSTableViewDelegate>
@property unsigned long long paneID;
@property (retain) NSMutableArray<WailsContentListRowModel*>* rows;
@property (retain) NSArray<WailsContentListRowModel*>* visibleRows;
@property (retain) NSMutableArray<WailsContentListColumnModel*>* columns;
@property (retain) NSMutableArray<NSNumber*>* selectedIDs;
@property (retain) NSTableView* tableView;
@property (retain) NSScrollView* scrollView;
@property (retain) NSTableHeaderView* headerView;
@property (retain) NSTextField* emptyLabel;
@property (retain) NSColor* surfaceColor;
@property (copy) NSString* emptyText;
@property BOOL allowsMultipleSelection;
@property BOOL sortable;
@property int sortColumn;
@property BOOL sortAscending;
@property double rowHeight;
@property int style;
@property BOOL alternatingRows;
@property BOOL headerVisible;
@property BOOL columnsDirty;
@property BOOL suppressCallbacks;
- (void)reloadContents;
- (void)applySelection;
- (void)reloadRowWithID:(unsigned long long)rowID hiddenChanged:(BOOL)hiddenChanged;
- (NSMenu*)contextMenuForEvent:(NSEvent*)event;
- (BOOL)activateSelectedRow;
@end

@implementation WailsContentListViewController

- (instancetype)init {
    self = [super init];
    if (self != nil) {
        _rows = [[NSMutableArray alloc] init];
        _columns = [[NSMutableArray alloc] init];
        _selectedIDs = [[NSMutableArray alloc] init];
        _visibleRows = [[NSArray alloc] init];
        _sortColumn = -1;
        _columnsDirty = YES;
    }
    return self;
}

- (void)dealloc {
    [_rows release];
    [_visibleRows release];
    [_columns release];
    [_selectedIDs release];
    [_tableView release];
    [_scrollView release];
    [_headerView release];
    [_emptyLabel release];
    [_surfaceColor release];
    [_emptyText release];
    [super dealloc];
}

- (void)loadView {
    NSView* container = [[NSView alloc] initWithFrame:NSMakeRect(0, 0, 320, 600)];

    NSScrollView* scrollView = [[NSScrollView alloc] initWithFrame:container.bounds];
    scrollView.autoresizingMask = NSViewWidthSizable | NSViewHeightSizable;
    scrollView.drawsBackground = self.surfaceColor != nil;
    if (self.surfaceColor != nil) scrollView.backgroundColor = self.surfaceColor;
    scrollView.borderType = NSNoBorder;
    scrollView.hasVerticalScroller = YES;
    scrollView.autohidesScrollers = YES;

    WailsContentListTableView* tableView = [[WailsContentListTableView alloc] initWithFrame:scrollView.bounds];
    tableView.controller = self;
    tableView.backgroundColor = [NSColor clearColor];
    tableView.allowsEmptySelection = YES;
    tableView.allowsColumnReordering = NO;
    tableView.allowsColumnSelection = NO;
    tableView.allowsMultipleSelection = self.allowsMultipleSelection;
    tableView.rowSizeStyle = NSTableViewRowSizeStyleCustom;
    tableView.target = self;
    tableView.doubleAction = @selector(handleDoubleClick:);
    tableView.dataSource = self;
    tableView.delegate = self;
    self.headerView = tableView.headerView;

    scrollView.documentView = tableView;
    [container addSubview:scrollView];

    NSTextField* emptyLabel = [NSTextField wrappingLabelWithString:@""];
    emptyLabel.translatesAutoresizingMaskIntoConstraints = NO;
    emptyLabel.alignment = NSTextAlignmentCenter;
    emptyLabel.textColor = [NSColor secondaryLabelColor];
    emptyLabel.font = [NSFont systemFontOfSize:[NSFont systemFontSize]];
    emptyLabel.hidden = YES;
    [container addSubview:emptyLabel];
    [NSLayoutConstraint activateConstraints:@[
        [emptyLabel.centerXAnchor constraintEqualToAnchor:container.centerXAnchor],
        [emptyLabel.centerYAnchor constraintEqualToAnchor:container.centerYAnchor],
        [emptyLabel.leadingAnchor constraintGreaterThanOrEqualToAnchor:container.leadingAnchor constant:16],
        [emptyLabel.trailingAnchor constraintLessThanOrEqualToAnchor:container.trailingAnchor constant:-16]
    ]];

    self.tableView = tableView;
    self.scrollView = scrollView;
    self.emptyLabel = emptyLabel;
    self.view = container;
    [tableView release];
    [scrollView release];
    [container release];
}

- (BOOL)richMode {
    return self.columns.count == 0;
}

- (void)rebuildColumns {
    NSTableView* tableView = self.tableView;
    for (NSTableColumn* column in [[tableView.tableColumns copy] autorelease]) {
        [tableView removeTableColumn:column];
    }
    if ([self richMode]) {
        NSTableColumn* column = [[NSTableColumn alloc] initWithIdentifier:@"WailsContentListRich"];
        column.title = @"";
        column.resizingMask = NSTableColumnAutoresizingMask;
        [tableView addTableColumn:column];
        [column release];
        tableView.columnAutoresizingStyle = NSTableViewLastColumnOnlyAutoresizingStyle;
    } else {
        NSUInteger index = 0;
        for (WailsContentListColumnModel* model in self.columns) {
            NSString* identifier = [NSString stringWithFormat:@"WailsContentListColumn%lu", (unsigned long)index];
            NSTableColumn* column = [[NSTableColumn alloc] initWithIdentifier:identifier];
            column.title = model.title ?: @"";
            column.resizingMask = NSTableColumnAutoresizingMask | NSTableColumnUserResizingMask;
            if (model.minWidth > 0) column.minWidth = model.minWidth;
            if (model.maxWidth > 0) column.maxWidth = model.maxWidth;
            if (model.width > 0) column.width = model.width;
            column.headerCell.alignment = contentListTextAlignment(model.alignment);
            if (model.sortable) {
                column.sortDescriptorPrototype = [NSSortDescriptor
                    sortDescriptorWithKey:[NSString stringWithFormat:@"%lu", (unsigned long)index] ascending:YES];
            }
            [tableView addTableColumn:column];
            [column release];
            index++;
        }
        tableView.columnAutoresizingStyle = NSTableViewUniformColumnAutoresizingStyle;
    }
    self.columnsDirty = NO;
}

- (void)applyOptions {
    NSTableView* tableView = self.tableView;
    if (tableView == nil) return;
    tableView.allowsMultipleSelection = self.allowsMultipleSelection;
    tableView.usesAlternatingRowBackgroundColors = self.alternatingRows;
    tableView.rowHeight = self.rowHeight > 0 ? self.rowHeight : ([self richMode] ? 46 : 24);
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
    if (@available(macOS 11.0, *)) {
        NSTableViewStyle style = NSTableViewStyleAutomatic;
        switch (self.style) {
        case 1: style = NSTableViewStyleInset; break;
        case 2: style = NSTableViewStyleSourceList; break;
        case 3: style = NSTableViewStylePlain; break;
        case 4: style = NSTableViewStyleFullWidth; break;
        default: break;
        }
        tableView.style = style;
    }
#endif
    tableView.headerView = self.headerVisible ? self.headerView : nil;

    NSArray<NSSortDescriptor*>* descriptors = @[];
    if (self.sortable && self.sortColumn >= 0 && (NSUInteger)self.sortColumn < self.columns.count &&
        self.columns[(NSUInteger)self.sortColumn].sortable) {
        descriptors = @[[NSSortDescriptor sortDescriptorWithKey:[NSString stringWithFormat:@"%d", self.sortColumn]
            ascending:self.sortAscending]];
    }
    if (![tableView.sortDescriptors isEqualToArray:descriptors]) tableView.sortDescriptors = descriptors;
}

- (void)refreshVisibleRows {
    NSMutableArray<WailsContentListRowModel*>* visible = [NSMutableArray arrayWithCapacity:self.rows.count];
    for (WailsContentListRowModel* row in self.rows) {
        if (!row.hidden) [visible addObject:row];
    }
    self.visibleRows = visible;
}

- (NSInteger)visibleIndexOfRowID:(unsigned long long)rowID {
    NSInteger index = 0;
    for (WailsContentListRowModel* row in self.visibleRows) {
        if (row.rowID == rowID) return index;
        index++;
    }
    return -1;
}

- (void)reloadContents {
    if (self.tableView == nil) return;
    self.suppressCallbacks = YES;
    if (self.columnsDirty) [self rebuildColumns];
    [self applyOptions];
    [self refreshVisibleRows];
    if ([self richMode]) {
        CGFloat availableWidth = self.tableView.bounds.size.width;
        if (availableWidth > 0) self.tableView.tableColumns.firstObject.width = availableWidth;
    }
    [self.tableView reloadData];
    [self restoreSelectionScrolling:NO];
    self.emptyLabel.stringValue = self.emptyText ?: @"";
    self.emptyLabel.hidden = self.visibleRows.count > 0 || self.emptyText.length == 0;
    self.suppressCallbacks = NO;
}

// restoreSelectionScrolling selects the rows in selectedIDs that are still
// visible and prunes the rest.
- (void)restoreSelectionScrolling:(BOOL)scroll {
    NSMutableIndexSet* indexes = [NSMutableIndexSet indexSet];
    NSMutableArray<NSNumber*>* kept = [NSMutableArray arrayWithCapacity:self.selectedIDs.count];
    for (NSNumber* rowID in self.selectedIDs) {
        NSInteger index = [self visibleIndexOfRowID:rowID.unsignedLongLongValue];
        if (index < 0) continue;
        [indexes addIndex:(NSUInteger)index];
        [kept addObject:rowID];
        if (!self.allowsMultipleSelection) break;
    }
    [self.selectedIDs setArray:kept];
    [self.tableView selectRowIndexes:indexes byExtendingSelection:NO];
    if (scroll && indexes.count > 0) [self.tableView scrollRowToVisible:(NSInteger)indexes.firstIndex];
}

- (void)applySelection {
    if (self.tableView == nil) return;
    self.suppressCallbacks = YES;
    [self restoreSelectionScrolling:YES];
    self.suppressCallbacks = NO;
}

- (void)reloadRowWithID:(unsigned long long)rowID hiddenChanged:(BOOL)hiddenChanged {
    if (self.tableView == nil) return;
    if (hiddenChanged) {
        [self reloadContents];
        return;
    }
    NSInteger index = [self visibleIndexOfRowID:rowID];
    if (index < 0) return;
    NSIndexSet* columns = [NSIndexSet indexSetWithIndexesInRange:NSMakeRange(0, (NSUInteger)self.tableView.numberOfColumns)];
    [self.tableView reloadDataForRowIndexes:[NSIndexSet indexSetWithIndex:(NSUInteger)index] columnIndexes:columns];
}

// Data source and delegate.

- (NSInteger)numberOfRowsInTableView:(NSTableView*)tableView {
    return (NSInteger)self.visibleRows.count;
}

- (WailsContentListRowModel*)rowAtVisibleIndex:(NSInteger)index {
    if (index < 0 || (NSUInteger)index >= self.visibleRows.count) return nil;
    return self.visibleRows[(NSUInteger)index];
}

- (WailsContentListRichCell*)newRichCell {
    WailsContentListRichCell* cell = [[[WailsContentListRichCell alloc] initWithFrame:NSMakeRect(0, 0, 300, 46)] autorelease];
    cell.identifier = @"WailsContentListRichCell";

    NSImageView* imageView = [[[NSImageView alloc] initWithFrame:NSZeroRect] autorelease];
    imageView.translatesAutoresizingMaskIntoConstraints = NO;
    imageView.imageScaling = NSImageScaleProportionallyUpOrDown;
    [NSLayoutConstraint activateConstraints:@[
        [imageView.widthAnchor constraintEqualToConstant:20],
        [imageView.heightAnchor constraintEqualToConstant:20]
    ]];

    NSTextField* title = [NSTextField labelWithString:@""];
    title.translatesAutoresizingMaskIntoConstraints = NO;
    title.font = [NSFont systemFontOfSize:13 weight:NSFontWeightSemibold];
    title.lineBreakMode = NSLineBreakByTruncatingTail;
    [title setContentHuggingPriority:NSLayoutPriorityDefaultLow forOrientation:NSLayoutConstraintOrientationHorizontal];
    [title setContentCompressionResistancePriority:NSLayoutPriorityDefaultLow forOrientation:NSLayoutConstraintOrientationHorizontal];

    NSTextField* detail = [NSTextField labelWithString:@""];
    detail.translatesAutoresizingMaskIntoConstraints = NO;
    detail.font = [NSFont systemFontOfSize:11];
    detail.alignment = NSTextAlignmentRight;
    [detail setContentHuggingPriority:NSLayoutPriorityRequired forOrientation:NSLayoutConstraintOrientationHorizontal];
    [detail setContentCompressionResistancePriority:NSLayoutPriorityRequired forOrientation:NSLayoutConstraintOrientationHorizontal];

    NSTextField* subtitle = [NSTextField labelWithString:@""];
    subtitle.translatesAutoresizingMaskIntoConstraints = NO;
    subtitle.font = [NSFont systemFontOfSize:11];
    subtitle.lineBreakMode = NSLineBreakByTruncatingTail;
    [subtitle setContentHuggingPriority:NSLayoutPriorityDefaultLow forOrientation:NSLayoutConstraintOrientationHorizontal];
    [subtitle setContentCompressionResistancePriority:NSLayoutPriorityDefaultLow forOrientation:NSLayoutConstraintOrientationHorizontal];

    WailsContentListBadgeView* badge = [[[WailsContentListBadgeView alloc] initWithFrame:NSZeroRect] autorelease];
    badge.translatesAutoresizingMaskIntoConstraints = NO;
    [badge setContentHuggingPriority:NSLayoutPriorityRequired forOrientation:NSLayoutConstraintOrientationHorizontal];
    [badge setContentCompressionResistancePriority:NSLayoutPriorityRequired forOrientation:NSLayoutConstraintOrientationHorizontal];

    NSStackView* topRow = [NSStackView stackViewWithViews:@[title, detail]];
    topRow.orientation = NSUserInterfaceLayoutOrientationHorizontal;
    topRow.alignment = NSLayoutAttributeFirstBaseline;
    topRow.distribution = NSStackViewDistributionFill;
    topRow.spacing = 8;

    NSStackView* bottomRow = [NSStackView stackViewWithViews:@[subtitle, badge]];
    bottomRow.orientation = NSUserInterfaceLayoutOrientationHorizontal;
    bottomRow.alignment = NSLayoutAttributeCenterY;
    bottomRow.distribution = NSStackViewDistributionFill;
    bottomRow.spacing = 8;

    NSStackView* textColumn = [NSStackView stackViewWithViews:@[topRow, bottomRow]];
    textColumn.orientation = NSUserInterfaceLayoutOrientationVertical;
    textColumn.alignment = NSLayoutAttributeWidth;
    textColumn.distribution = NSStackViewDistributionFill;
    textColumn.spacing = 2;
    [textColumn setContentHuggingPriority:NSLayoutPriorityDefaultLow forOrientation:NSLayoutConstraintOrientationHorizontal];

    NSStackView* stack = [NSStackView stackViewWithViews:@[imageView, textColumn]];
    stack.translatesAutoresizingMaskIntoConstraints = NO;
    stack.orientation = NSUserInterfaceLayoutOrientationHorizontal;
    stack.alignment = NSLayoutAttributeCenterY;
    stack.distribution = NSStackViewDistributionFill;
    stack.spacing = 8;
    [cell addSubview:stack];
    [NSLayoutConstraint activateConstraints:@[
        [stack.leadingAnchor constraintEqualToAnchor:cell.leadingAnchor constant:4],
        [stack.trailingAnchor constraintEqualToAnchor:cell.trailingAnchor constant:-8],
        [stack.centerYAnchor constraintEqualToAnchor:cell.centerYAnchor]
    ]];
    cell.imageView = imageView;
    cell.textField = title;
    cell.subtitleField = subtitle;
    cell.detailField = detail;
    cell.badgeView = badge;
    return cell;
}

- (NSTableCellView*)newTableCellWithIdentifier:(NSUserInterfaceItemIdentifier)identifier alignment:(NSTextAlignment)alignment {
    NSTableCellView* cell = [[[NSTableCellView alloc] initWithFrame:NSMakeRect(0, 0, 120, 24)] autorelease];
    cell.identifier = identifier;
    NSTextField* textField = [NSTextField labelWithString:@""];
    textField.translatesAutoresizingMaskIntoConstraints = NO;
    textField.lineBreakMode = NSLineBreakByTruncatingTail;
    textField.alignment = alignment;
    textField.font = [NSFont systemFontOfSize:[NSFont systemFontSize]];
    [cell addSubview:textField];
    [NSLayoutConstraint activateConstraints:@[
        [textField.leadingAnchor constraintEqualToAnchor:cell.leadingAnchor constant:2],
        [textField.trailingAnchor constraintEqualToAnchor:cell.trailingAnchor constant:-2],
        [textField.centerYAnchor constraintEqualToAnchor:cell.centerYAnchor]
    ]];
    cell.textField = textField;
    return cell;
}

- (NSView*)tableView:(NSTableView*)tableView viewForTableColumn:(NSTableColumn*)tableColumn row:(NSInteger)rowIndex {
    WailsContentListRowModel* row = [self rowAtVisibleIndex:rowIndex];
    if (row == nil) return nil;
    if ([self richMode]) {
        WailsContentListRichCell* cell = [tableView makeViewWithIdentifier:@"WailsContentListRichCell" owner:self];
        if (cell == nil) cell = [self newRichCell];
        cell.disabledRow = row.disabled;
        cell.textField.stringValue = row.title ?: @"";
        cell.subtitleField.stringValue = row.subtitle ?: @"";
        cell.subtitleField.hidden = row.subtitle.length == 0;
        cell.detailField.stringValue = row.detail ?: @"";
        cell.detailField.hidden = row.detail.length == 0;
        cell.badgeView.count = row.badge;
        cell.badgeView.hidden = row.badge <= 0;
        NSImage* image = contentListSymbolImage(row.symbol);
        cell.imageView.image = image;
        cell.imageView.hidden = image == nil;
        cell.toolTip = row.tooltip.length > 0 ? row.tooltip : nil;
        [cell applyColors];
        return cell;
    }
    NSUInteger columnIndex = [tableView.tableColumns indexOfObjectIdenticalTo:tableColumn];
    NSTextAlignment alignment = columnIndex < self.columns.count
        ? contentListTextAlignment(self.columns[columnIndex].alignment) : NSTextAlignmentNatural;
    NSUserInterfaceItemIdentifier identifier = [NSString stringWithFormat:@"WailsContentListCell%ld", (long)alignment];
    NSTableCellView* cell = [tableView makeViewWithIdentifier:identifier owner:self];
    if (cell == nil) cell = [self newTableCellWithIdentifier:identifier alignment:alignment];
    NSString* value = columnIndex < row.cells.count ? row.cells[columnIndex] : @"";
    if (row.cells.count == 0 && columnIndex == 0) value = row.title ?: @"";
    cell.textField.stringValue = value ?: @"";
    cell.textField.textColor = row.disabled ? [NSColor disabledControlTextColor] : [NSColor labelColor];
    cell.toolTip = row.tooltip.length > 0 ? row.tooltip : nil;
    return cell;
}

- (BOOL)tableView:(NSTableView*)tableView shouldSelectRow:(NSInteger)rowIndex {
    WailsContentListRowModel* row = [self rowAtVisibleIndex:rowIndex];
    return row != nil && !row.disabled;
}

- (void)tableViewSelectionDidChange:(NSNotification*)notification {
    if (self.suppressCallbacks) return;
    NSIndexSet* indexes = self.tableView.selectedRowIndexes;
    NSMutableArray<NSNumber*>* ids = [NSMutableArray arrayWithCapacity:indexes.count];
    [indexes enumerateIndexesUsingBlock:^(NSUInteger index, BOOL* stop) {
        WailsContentListRowModel* row = [self rowAtVisibleIndex:(NSInteger)index];
        if (row != nil && !row.disabled) [ids addObject:@(row.rowID)];
    }];
    [self.selectedIDs setArray:ids];
    unsigned long long* buffer = ids.count > 0 ? calloc(ids.count, sizeof(unsigned long long)) : NULL;
    for (NSUInteger index = 0; index < ids.count && buffer != NULL; index++) {
        buffer[index] = ids[index].unsignedLongLongValue;
    }
    processMacContentListSelectionChanged(self.paneID, buffer, (int)ids.count);
    free(buffer);
}

- (void)tableView:(NSTableView*)tableView sortDescriptorsDidChange:(NSArray<NSSortDescriptor*>*)oldDescriptors {
    if (self.suppressCallbacks || !self.sortable) return;
    NSSortDescriptor* descriptor = tableView.sortDescriptors.firstObject;
    if (descriptor == nil || descriptor.key.length == 0) return;
    int column = descriptor.key.intValue;
    if (column < 0 || (NSUInteger)column >= self.columns.count) return;
    self.sortColumn = column;
    self.sortAscending = descriptor.ascending;
    processMacContentListSortChanged(self.paneID, column, descriptor.ascending);
}

// Activation: double-click or Return on a single selected row.

- (BOOL)activateRow:(NSInteger)index {
    WailsContentListRowModel* row = [self rowAtVisibleIndex:index];
    if (row == nil || row.disabled) return NO;
    processMacContentListRowActivated(self.paneID, row.rowID);
    return YES;
}

- (void)handleDoubleClick:(id)sender {
    [self activateRow:self.tableView.clickedRow];
}

- (BOOL)activateSelectedRow {
    if (self.tableView.numberOfSelectedRows != 1) return NO;
    return [self activateRow:self.tableView.selectedRow];
}

// Context menus.

- (NSMenu*)contextMenuForEvent:(NSEvent*)event {
    NSPoint point = [self.tableView convertPoint:event.locationInWindow fromView:nil];
    NSInteger index = [self.tableView rowAtPoint:point];
    WailsContentListRowModel* row = [self rowAtVisibleIndex:index];
    return (NSMenu*)processMacContentListContextMenu(self.paneID, row == nil ? 0 : row.rowID);
}

@end

@implementation WailsContentListTableView

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
            [self.controller activateSelectedRow]) {
            return;
        }
    }
    [super keyDown:event];
}

@end

// Hooks for webview_window_split_darwin.m.

NSViewController* wailsContentListCreateController(unsigned long long paneID) {
    WailsContentListViewController* controller = [[[WailsContentListViewController alloc] init] autorelease];
    controller.paneID = paneID;
    return controller;
}

static WailsContentListViewController* contentListController(NSViewController* controller) {
    return [controller isKindOfClass:[WailsContentListViewController class]]
        ? (WailsContentListViewController*)controller : nil;
}

bool wailsContentListPrepareForInstall(NSViewController* viewController, NSColor* surfaceColor) {
    WailsContentListViewController* controller = contentListController(viewController);
    if (controller == nil) return false;
    controller.surfaceColor = surfaceColor;
    (void)controller.view;
    return controller.view != nil;
}

NSSplitViewItem* wailsContentListSplitItem(NSViewController* viewController) {
    if (viewController == nil) return nil;
    NSSplitViewItem* item = nil;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
    if (@available(macOS 11.0, *)) {
        item = [NSSplitViewItem contentListWithViewController:viewController];
    }
#endif
    if (item == nil) item = [NSSplitViewItem splitViewItemWithViewController:viewController];
    return item;
}

void wailsContentListDidInstall(NSViewController* viewController) {
    [contentListController(viewController) reloadContents];
}

// C entry points used by Go.

static WailsContentListViewController* contentListForPane(void* handlePtr, unsigned long long paneID) {
    return contentListController(splitViewContentListController(handlePtr, paneID));
}

void splitViewContentListReset(void* handlePtr, unsigned long long paneID) {
    WailsContentListViewController* controller = contentListForPane(handlePtr, paneID);
    if (controller == nil) return;
    [controller.rows removeAllObjects];
    [controller.columns removeAllObjects];
    controller.columnsDirty = YES;
}

void splitViewContentListSetOptions(void* handlePtr, unsigned long long paneID, WailsContentListOptions options) {
    WailsContentListViewController* controller = contentListForPane(handlePtr, paneID);
    if (controller == nil) return;
    controller.allowsMultipleSelection = options.allowsMultipleSelection;
    controller.sortable = options.sortable;
    controller.sortColumn = options.sortColumn;
    controller.sortAscending = options.sortAscending;
    controller.rowHeight = options.rowHeight;
    controller.style = options.style;
    controller.alternatingRows = options.alternatingRows;
    controller.emptyText = contentListString(options.emptyText);
    controller.headerVisible = options.headerVisible;
}

void splitViewContentListAddColumn(void* handlePtr, unsigned long long paneID,
    const char* title, double width, double minWidth, double maxWidth, bool sortable, int alignment) {
    WailsContentListViewController* controller = contentListForPane(handlePtr, paneID);
    if (controller == nil) return;
    WailsContentListColumnModel* column = [[WailsContentListColumnModel alloc] init];
    column.title = contentListString(title);
    column.width = width;
    column.minWidth = minWidth;
    column.maxWidth = maxWidth;
    column.sortable = sortable;
    column.alignment = alignment;
    [controller.columns addObject:column];
    [column release];
    controller.columnsDirty = YES;
}

void splitViewContentListAddRow(void* handlePtr, unsigned long long paneID,
    unsigned long long rowID, WailsContentListRowSpec spec) {
    WailsContentListViewController* controller = contentListForPane(handlePtr, paneID);
    if (controller == nil) return;
    WailsContentListRowModel* row = [[WailsContentListRowModel alloc] init];
    row.rowID = rowID;
    contentListConfigureRow(row, spec);
    [controller.rows addObject:row];
    [row release];
}

void splitViewContentListSetSelection(void* handlePtr, unsigned long long paneID,
    unsigned long long* rowIDs, int count) {
    WailsContentListViewController* controller = contentListForPane(handlePtr, paneID);
    if (controller == nil) return;
    NSMutableArray<NSNumber*>* ids = [NSMutableArray arrayWithCapacity:(NSUInteger)MAX(count, 0)];
    for (int index = 0; index < count && rowIDs != NULL; index++) {
        [ids addObject:@(rowIDs[index])];
    }
    [controller.selectedIDs setArray:ids];
    [controller applySelection];
}

void splitViewContentListReload(void* handlePtr, unsigned long long paneID) {
    [contentListForPane(handlePtr, paneID) reloadContents];
}

void splitViewContentListUpdateRow(void* handlePtr, unsigned long long paneID,
    unsigned long long rowID, WailsContentListRowSpec spec) {
    WailsContentListViewController* controller = contentListForPane(handlePtr, paneID);
    if (controller == nil) return;
    for (WailsContentListRowModel* row in controller.rows) {
        if (row.rowID != rowID) continue;
        BOOL wasHidden = row.hidden;
        contentListConfigureRow(row, spec);
        [controller reloadRowWithID:rowID hiddenChanged:wasHidden != row.hidden];
        return;
    }
}
