//go:build darwin && !ios && !server

#import "appleevents_manager_darwin.h"
#import "application_darwin_delegate.h"
#import <CoreServices/CoreServices.h>

// Delivers an incoming event to Go. When suspensionID is non-NULL the event
// has been suspended and Go resumes it later through appleEventsResume; when
// it is NULL Go runs the handler inline and returns the reply JSON, which the
// caller frees.
extern char *appleEventsDeliver(char *eventJSON, void *suspensionID);

static NSString *wailsAEFourCharString(FourCharCode code) {
    char chars[5] = {
        (char)((code >> 24) & 0xFF), (char)((code >> 16) & 0xFF),
        (char)((code >> 8) & 0xFF), (char)(code & 0xFF), 0,
    };
    return [[[NSString alloc] initWithBytes:chars length:4 encoding:NSMacOSRomanStringEncoding] autorelease];
}

static FourCharCode wailsAEFourCharCode(NSString *code) {
    NSData *data = [code dataUsingEncoding:NSMacOSRomanStringEncoding allowLossyConversion:YES];
    if (data.length != 4) {
        return 0;
    }
    const unsigned char *bytes = data.bytes;
    return ((FourCharCode)bytes[0] << 24) | ((FourCharCode)bytes[1] << 16) |
           ((FourCharCode)bytes[2] << 8) | (FourCharCode)bytes[3];
}

static FourCharCode wailsAEFourCharCodeC(const char *code) {
    return wailsAEFourCharCode([NSString stringWithUTF8String:code ? code : ""]);
}

// Descriptor -> wire dictionary (see the codec comment in appleevents_manager.go).
static NSDictionary *wailsAEDescriptorToWire(NSAppleEventDescriptor *desc) {
    if (desc == nil) {
        return @{@"t": @"null"};
    }
    DescType type = [desc descriptorType];
    switch (type) {
    case typeNull:
        return @{@"t": @"null"};
    case typeTrue:
    case typeFalse:
    case typeBoolean:
        return @{@"t": @"bool", @"v": @([desc booleanValue])};
    case typeSInt16:
    case typeSInt32:
    case typeUInt16:
    case typeUInt32:
    case typeSInt64:
    case typeUInt64: {
        NSAppleEventDescriptor *coerced = [desc coerceToDescriptorType:typeSInt64];
        if (coerced != nil && [[coerced data] length] == sizeof(int64_t)) {
            int64_t value = 0;
            [[coerced data] getBytes:&value length:sizeof(value)];
            return @{@"t": @"int", @"v": @(value)};
        }
        return @{@"t": @"int", @"v": @([desc int32Value])};
    }
    case typeIEEE32BitFloatingPoint:
    case typeIEEE64BitFloatingPoint:
        return @{@"t": @"float", @"v": @([desc doubleValue])};
    case typeFileURL:
    case typeAlias:
    case typeBookmarkData: {
        NSURL *url = [desc fileURLValue];
        if (url == nil) {
            url = [[desc coerceToDescriptorType:typeFileURL] fileURLValue];
        }
        if (url != nil && url.path != nil) {
            return @{@"t": @"file", @"v": url.path};
        }
        break;
    }
    case typeUnicodeText:
    case typeUTF8Text:
    case typeUTF16ExternalRepresentation:
    case typeChar:
    case typeCString:
    case typeStyledText:
    case typeIntlText: {
        NSString *text = [desc stringValue];
        if (text != nil) {
            return @{@"t": @"string", @"v": text};
        }
        break;
    }
    case typeAEList: {
        NSInteger count = [desc numberOfItems];
        NSMutableArray *items = [NSMutableArray arrayWithCapacity:(NSUInteger)MAX(count, 0)];
        for (NSInteger i = 1; i <= count; i++) {
            [items addObject:wailsAEDescriptorToWire([desc descriptorAtIndex:i])];
        }
        return @{@"t": @"list", @"v": items};
    }
    default:
        break;
    }
    NSData *data = [desc data] ?: [NSData data];
    return @{@"t": @"raw", @"type": wailsAEFourCharString(type), @"v": [data base64EncodedStringWithOptions:0]};
}

// Wire dictionary -> descriptor. Returns nil for malformed input.
static NSAppleEventDescriptor *wailsAEWireToDescriptor(id wire) {
    if (![wire isKindOfClass:[NSDictionary class]]) {
        return nil;
    }
    NSString *kind = wire[@"t"];
    id value = wire[@"v"];
    if (![kind isKindOfClass:[NSString class]] || [kind isEqualToString:@"null"]) {
        return [NSAppleEventDescriptor nullDescriptor];
    }
    if ([kind isEqualToString:@"string"] && [value isKindOfClass:[NSString class]]) {
        return [NSAppleEventDescriptor descriptorWithString:value];
    }
    if ([kind isEqualToString:@"file"] && [value isKindOfClass:[NSString class]]) {
        return [NSAppleEventDescriptor descriptorWithFileURL:[NSURL fileURLWithPath:value]];
    }
    if ([kind isEqualToString:@"bool"] && [value isKindOfClass:[NSNumber class]]) {
        return [NSAppleEventDescriptor descriptorWithBoolean:[value boolValue]];
    }
    if ([kind isEqualToString:@"int"] && [value isKindOfClass:[NSNumber class]]) {
        long long number = [value longLongValue];
        if (number >= INT32_MIN && number <= INT32_MAX) {
            return [NSAppleEventDescriptor descriptorWithInt32:(SInt32)number];
        }
        int64_t wide = (int64_t)number;
        return [NSAppleEventDescriptor descriptorWithDescriptorType:typeSInt64 bytes:&wide length:sizeof(wide)];
    }
    if ([kind isEqualToString:@"float"] && [value isKindOfClass:[NSNumber class]]) {
        return [NSAppleEventDescriptor descriptorWithDouble:[value doubleValue]];
    }
    if ([kind isEqualToString:@"list"] && [value isKindOfClass:[NSArray class]]) {
        NSAppleEventDescriptor *list = [NSAppleEventDescriptor listDescriptor];
        for (id item in (NSArray *)value) {
            NSAppleEventDescriptor *itemDesc = wailsAEWireToDescriptor(item);
            if (itemDesc == nil) {
                return nil;
            }
            [list insertDescriptor:itemDesc atIndex:[list numberOfItems] + 1];
        }
        return list;
    }
    if ([kind isEqualToString:@"raw"] && [value isKindOfClass:[NSString class]]) {
        NSString *typeString = wire[@"type"];
        NSData *data = [[[NSData alloc] initWithBase64EncodedString:value options:0] autorelease];
        if (![typeString isKindOfClass:[NSString class]] || data == nil) {
            return nil;
        }
        return [NSAppleEventDescriptor descriptorWithDescriptorType:wailsAEFourCharCode(typeString) data:data];
    }
    return nil;
}

static id wailsAEParseJSON(const char *json) {
    if (json == NULL) {
        return nil;
    }
    NSData *data = [NSData dataWithBytes:json length:strlen(json)];
    return [NSJSONSerialization JSONObjectWithData:data options:0 error:nil];
}

static NSString *wailsAEJSONString(id object) {
    NSData *data = [NSJSONSerialization dataWithJSONObject:object options:0 error:nil];
    if (data == nil) {
        return nil;
    }
    return [[[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding] autorelease];
}

// Serialises an incoming event: its class and ID, direct object and other
// parameters.
static NSString *wailsAEEventToJSON(NSAppleEventDescriptor *event) {
    NSMutableDictionary *root = [NSMutableDictionary dictionary];
    root[@"class"] = wailsAEFourCharString((FourCharCode)[[event attributeDescriptorForKeyword:keyEventClassAttr] typeCodeValue]);
    root[@"id"] = wailsAEFourCharString((FourCharCode)[[event attributeDescriptorForKeyword:keyEventIDAttr] typeCodeValue]);
    NSMutableDictionary *params = [NSMutableDictionary dictionary];
    NSInteger count = [event numberOfItems];
    for (NSInteger i = 1; i <= count; i++) {
        AEKeyword keyword = [event keywordForDescriptorAtIndex:i];
        NSDictionary *wire = wailsAEDescriptorToWire([event descriptorAtIndex:i]);
        if (keyword == keyDirectObject) {
            root[@"direct"] = wire;
        } else {
            params[wailsAEFourCharString(keyword)] = wire;
        }
    }
    root[@"params"] = params;
    return wailsAEJSONString(root);
}

// Fills a reply descriptor from {"result":wire,"errn":n,"errs":"text"}.
// Replies of type typeNull (the sender did not ask for one) are left alone.
static void wailsAEFillReply(NSAppleEventDescriptor *reply, const char *replyJSON) {
    if (reply == nil || [reply descriptorType] == typeNull) {
        return;
    }
    NSDictionary *wire = wailsAEParseJSON(replyJSON);
    if (![wire isKindOfClass:[NSDictionary class]]) {
        return;
    }
    NSAppleEventDescriptor *result = wailsAEWireToDescriptor(wire[@"result"]);
    if (result != nil && [result descriptorType] != typeNull) {
        [reply setParamDescriptor:result forKeyword:keyDirectObject];
    }
    NSNumber *errorNumber = wire[@"errn"];
    if ([errorNumber isKindOfClass:[NSNumber class]] && [errorNumber intValue] != 0) {
        [reply setParamDescriptor:[NSAppleEventDescriptor descriptorWithInt32:[errorNumber intValue]] forKeyword:keyErrorNumber];
        NSString *errorString = wire[@"errs"];
        if ([errorString isKindOfClass:[NSString class]] && errorString.length > 0) {
            [reply setParamDescriptor:[NSAppleEventDescriptor descriptorWithString:errorString] forKeyword:keyErrorString];
        }
    }
}

// WailsAppleEventHandler is the standalone object registered with
// NSAppleEventManager for every class/ID the application handles. It never
// touches the application delegate.
@interface WailsAppleEventHandler : NSObject
+ (instancetype)shared;
- (void)handleAppleEvent:(NSAppleEventDescriptor *)event withReplyEvent:(NSAppleEventDescriptor *)reply;
@end

@implementation WailsAppleEventHandler

+ (instancetype)shared {
    static WailsAppleEventHandler *shared = nil;
    static dispatch_once_t once;
    dispatch_once(&once, ^{
        shared = [[WailsAppleEventHandler alloc] init];
    });
    return shared;
}

- (void)handleAppleEvent:(NSAppleEventDescriptor *)event withReplyEvent:(NSAppleEventDescriptor *)reply {
  @autoreleasepool {
    AEEventClass eventClass = (AEEventClass)[[event attributeDescriptorForKeyword:keyEventClassAttr] typeCodeValue];
    AEEventID eventID = (AEEventID)[[event attributeDescriptorForKeyword:keyEventIDAttr] typeCodeValue];
    if (eventClass == kInternetEventClass && eventID == kAEGetURL) {
        // Keep the built-in custom URL scheme delivery working.
        [CustomProtocolSchemeHandler handleGetURLEvent:event withReplyEvent:reply];
    }

    NSString *json = wailsAEEventToJSON(event);
    if (json == nil) {
        return;
    }
    // Events from other processes arrive on the main thread and are
    // suspended so the Go handler can run on its own goroutine without
    // blocking the UI. Events this process sends to itself with kAEWaitReply
    // are dispatched directly on the sending thread, where suspension has
    // nothing to resume against, so they are answered inline instead.
    NSAppleEventManagerSuspensionID suspensionID = NULL;
    if ([NSThread isMainThread]) {
        suspensionID = [[NSAppleEventManager sharedAppleEventManager] suspendCurrentAppleEvent];
    }
    char *replyJSON = appleEventsDeliver((char *)[json UTF8String], (void *)suspensionID);
    if (suspensionID == NULL && replyJSON != NULL) {
        wailsAEFillReply(reply, replyJSON);
    }
    if (replyJSON != NULL) {
        free(replyJSON);
    }
  }
}

@end

void appleEventsRegister(const char *eventClass, const char *eventID) {
    @autoreleasepool {
        [[NSAppleEventManager sharedAppleEventManager] setEventHandler:[WailsAppleEventHandler shared]
                                                          andSelector:@selector(handleAppleEvent:withReplyEvent:)
                                                        forEventClass:wailsAEFourCharCodeC(eventClass)
                                                           andEventID:wailsAEFourCharCodeC(eventID)];
    }
}

void appleEventsResume(void *suspensionID, const char *replyJSON) {
    if (suspensionID == NULL) {
        return;
    }
    @autoreleasepool {
        NSAppleEventManager *manager = [NSAppleEventManager sharedAppleEventManager];
        NSAppleEventManagerSuspensionID sid = (NSAppleEventManagerSuspensionID)suspensionID;
        wailsAEFillReply([manager replyAppleEventForSuspensionID:sid], replyJSON);
        [manager resumeWithSuspensionID:sid];
    }
}

char *appleEventsSend(const char *bundleID, const char *eventClass, const char *eventID,
                      const char *directJSON, char **errorMessage) {
    @autoreleasepool {
        NSAppleEventDescriptor *target = [NSAppleEventDescriptor descriptorWithBundleIdentifier:[NSString stringWithUTF8String:bundleID]];
        NSAppleEventDescriptor *event = [NSAppleEventDescriptor appleEventWithEventClass:wailsAEFourCharCodeC(eventClass)
                                                                                 eventID:wailsAEFourCharCodeC(eventID)
                                                                        targetDescriptor:target
                                                                                returnID:kAutoGenerateReturnID
                                                                           transactionID:kAnyTransactionID];
        NSAppleEventDescriptor *direct = wailsAEWireToDescriptor(wailsAEParseJSON(directJSON));
        if (direct == nil) {
            *errorMessage = strdup("apple events: could not encode the direct object");
            return NULL;
        }
        if ([direct descriptorType] != typeNull) {
            [event setParamDescriptor:direct forKeyword:keyDirectObject];
        }
        NSError *error = nil;
        NSAppleEventDescriptor *reply = [event sendEventWithOptions:(kAEWaitReply | kAECanInteract | kAECanSwitchLayer)
                                                            timeout:kAEDefaultTimeout
                                                              error:&error];
        if (reply == nil) {
            NSString *message = error.localizedDescription ?: @"apple events: send failed";
            if (error != nil) {
                message = [NSString stringWithFormat:@"%@ (error %ld)", message, (long)error.code];
            }
            *errorMessage = strdup([message UTF8String]);
            return NULL;
        }
        NSMutableDictionary *wire = [NSMutableDictionary dictionary];
        NSAppleEventDescriptor *result = [reply paramDescriptorForKeyword:keyDirectObject];
        if (result != nil) {
            wire[@"result"] = wailsAEDescriptorToWire(result);
        }
        NSAppleEventDescriptor *errorNumber = [reply paramDescriptorForKeyword:keyErrorNumber];
        wire[@"errn"] = @(errorNumber != nil ? [errorNumber int32Value] : 0);
        NSAppleEventDescriptor *errorString = [reply paramDescriptorForKeyword:keyErrorString];
        if (errorString != nil && [errorString stringValue] != nil) {
            wire[@"errs"] = [errorString stringValue];
        }
        NSString *json = wailsAEJSONString(wire);
        if (json == nil) {
            *errorMessage = strdup("apple events: could not decode the reply");
            return NULL;
        }
        return strdup([json UTF8String]);
    }
}
