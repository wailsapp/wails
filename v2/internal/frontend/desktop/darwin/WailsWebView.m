#import "WailsWebView.h"
#import "message.h"


@implementation WailsWebView
@synthesize disableWebViewDragAndDrop;
@synthesize enableDragAndDrop;

- (BOOL)prepareForDragOperation:(id<NSDraggingInfo>)sender
{
  if ( !enableDragAndDrop ) {
    return [super prepareForDragOperation: sender];
  }

  if ( disableWebViewDragAndDrop ) {
    return YES;
  }

  return [super prepareForDragOperation: sender];
}

- (BOOL)performDragOperation:(id <NSDraggingInfo>)sender
{
  if ( !enableDragAndDrop ) {
    return [super performDragOperation: sender];
  }

  NSPasteboard *pboard = [sender draggingPasteboard];

  // if no types, then we'll just let the WKWebView handle the drag-n-drop as normal
  NSArray<NSPasteboardType> * types = [pboard types];
  if( !types )
    return [super performDragOperation: sender];

  // getting all NSURL types
  NSArray<Class> *url_class = @[[NSURL class]];
  NSDictionary *options = @{};
  NSArray<NSURL*> *files = [pboard readObjectsForClasses:url_class options:options];

  // collecting all file paths
  NSMutableArray *files_strs = [[NSMutableArray alloc] init];
  for (NSURL *url in files)
  {
    const char *fs_path = [url fileSystemRepresentation];  //Will be UTF-8 encoded
    NSString *fs_path_str = [[NSString alloc] initWithCString:fs_path encoding:NSUTF8StringEncoding];
    [files_strs addObject:fs_path_str];
//     NSLog( @"performDragOperation: file path: %s", fs_path );
  }

  NSString *joined=[files_strs componentsJoinedByString:@"\n"];
  
  // Release the array of file paths
  [files_strs release];

  int	dragXLocation = [sender draggingLocation].x - [self frame].origin.x;
  int	dragYLocation = [self frame].size.height - [sender draggingLocation].y; // Y coordinate is inverted, so we need to subtract from the height

//   NSLog( @"draggingUpdated: X coord: %d", dragXLocation );
//   NSLog( @"draggingUpdated: Y coord: %d", dragYLocation );

  NSString *message = [NSString stringWithFormat:@"DD:%d:%d:%@", dragXLocation, dragYLocation, joined];

  const char* res = message.UTF8String;

  processMessage(res);

  if ( disableWebViewDragAndDrop ) {
    return YES;
  }

  return [super performDragOperation: sender];
}

- (NSDragOperation)draggingUpdated:(id <NSDraggingInfo>)sender {
  if ( !enableDragAndDrop ) {
    return [super draggingUpdated: sender];
  }

  NSPasteboard *pboard = [sender draggingPasteboard];

  // if no types, then we'll just let the WKWebView handle the drag-n-drop as normal
  NSArray<NSPasteboardType> * types = [pboard types];
  if( !types ) {
    return [super draggingUpdated: sender];
  }

  if ( disableWebViewDragAndDrop ) {
    // we should call supper as otherwise events will not pass
    [super draggingUpdated: sender];

    // pass NSDragOperationGeneric = 4 to show regular hover for drag and drop. As we want to ignore webkit behaviours that depends on webpage
    return 4;
  }

  return [super draggingUpdated: sender];
}

- (NSDragOperation)draggingEntered:(id <NSDraggingInfo>)sender {
  if ( !enableDragAndDrop ) {
    return [super draggingEntered: sender];
  }

  NSPasteboard *pboard = [sender draggingPasteboard];

  // if no types, then we'll just let the WKWebView handle the drag-n-drop as normal
  NSArray<NSPasteboardType> * types = [pboard types];
  if( !types ) {
    return [super draggingEntered: sender];
  }

  if ( disableWebViewDragAndDrop ) {
    // we should call supper as otherwise events will not pass
    [super draggingEntered: sender];

    // pass NSDragOperationGeneric = 4 to show regular hover for drag and drop. As we want to ignore webkit behaviours that depends on webpage
    return 4;
  }

  return [super draggingEntered: sender];
}

@end
