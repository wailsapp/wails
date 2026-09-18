//go:build darwin && !ios && !server

#import "spotlight_manager_darwin.h"
#import <CoreSpotlight/CoreSpotlight.h>
#import <UniformTypeIdentifiers/UniformTypeIdentifiers.h>

// Delivers the outcome of an index operation to Go. errorMessage is NULL on
// success; the callee copies it before returning.
extern void spotlightOperationResult(unsigned long long requestID, const char *errorMessage);

// Entry point implemented in Go (spotlight_manager.go) for continuation
// activities. userInfoJSON is the activity's userInfo as a JSON object.
extern bool spotlightHandleContinueActivityC(char *activityType, char *userInfoJSON);

static id wailsSpotlightParseJSON(const char *json) {
    if (json == NULL) {
        return nil;
    }
    NSData *data = [NSData dataWithBytes:json length:strlen(json)];
    return [NSJSONSerialization JSONObjectWithData:data options:0 error:nil];
}

static void wailsSpotlightComplete(unsigned long long requestID, NSError *error) {
    // CoreSpotlight calls completion handlers on its own queue; the Go side
    // only forwards to a buffered channel so this never blocks.
    @autoreleasepool {
        if (error != nil) {
            NSString *message = [NSString stringWithFormat:@"%@ (error %ld)", error.localizedDescription, (long)error.code];
            spotlightOperationResult(requestID, [message UTF8String]);
        } else {
            spotlightOperationResult(requestID, NULL);
        }
    }
}

static NSString *wailsSpotlightString(NSDictionary *item, NSString *key) {
    id value = item[key];
    return [value isKindOfClass:[NSString class]] ? value : nil;
}

static CSSearchableItemAttributeSet *wailsSpotlightAttributeSet(NSString *contentType) {
    if (@available(macOS 11.0, *)) {
        UTType *type = [UTType typeWithIdentifier:contentType] ?: UTTypeItem;
        return [[[CSSearchableItemAttributeSet alloc] initWithContentType:type] autorelease];
    }
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"
    return [[[CSSearchableItemAttributeSet alloc] initWithItemContentType:contentType] autorelease];
#pragma clang diagnostic pop
}

bool spotlightIsAvailable(void) {
    return [CSSearchableIndex isIndexingAvailable];
}

void spotlightIndexItems(unsigned long long requestID, const char *itemsJSON) {
  @autoreleasepool {
    NSArray *entries = wailsSpotlightParseJSON(itemsJSON);
    if (![entries isKindOfClass:[NSArray class]]) {
        spotlightOperationResult(requestID, "spotlight: malformed item list");
        return;
    }
    NSMutableArray<CSSearchableItem *> *items = [NSMutableArray arrayWithCapacity:entries.count];
    for (NSDictionary *entry in entries) {
        if (![entry isKindOfClass:[NSDictionary class]]) {
            continue;
        }
        NSString *contentType = wailsSpotlightString(entry, @"contentType") ?: @"public.item";
        CSSearchableItemAttributeSet *attributes = wailsSpotlightAttributeSet(contentType);
        attributes.title = wailsSpotlightString(entry, @"title");
        attributes.contentDescription = wailsSpotlightString(entry, @"description");
        NSArray *keywords = entry[@"keywords"];
        if ([keywords isKindOfClass:[NSArray class]] && keywords.count > 0) {
            attributes.keywords = keywords;
        }
        NSString *thumbnail = wailsSpotlightString(entry, @"thumbnail");
        if (thumbnail.length > 0) {
            attributes.thumbnailData = [[[NSData alloc] initWithBase64EncodedString:thumbnail options:0] autorelease];
        }
        NSString *url = wailsSpotlightString(entry, @"url");
        if (url.length > 0) {
            attributes.contentURL = [NSURL URLWithString:url];
        }
        CSSearchableItem *item = [[[CSSearchableItem alloc] initWithUniqueIdentifier:wailsSpotlightString(entry, @"id")
                                                                    domainIdentifier:wailsSpotlightString(entry, @"domain")
                                                                        attributeSet:attributes] autorelease];
        NSNumber *expires = entry[@"expires"];
        if ([expires isKindOfClass:[NSNumber class]] && [expires doubleValue] > 0) {
            item.expirationDate = [NSDate dateWithTimeIntervalSince1970:[expires doubleValue]];
        }
        [items addObject:item];
    }
    [[CSSearchableIndex defaultSearchableIndex] indexSearchableItems:items completionHandler:^(NSError *error) {
        wailsSpotlightComplete(requestID, error);
    }];
  }
}

void spotlightDeleteItems(unsigned long long requestID, const char *idsJSON) {
  @autoreleasepool {
    NSArray *ids = wailsSpotlightParseJSON(idsJSON);
    if (![ids isKindOfClass:[NSArray class]]) {
        spotlightOperationResult(requestID, "spotlight: malformed identifier list");
        return;
    }
    [[CSSearchableIndex defaultSearchableIndex] deleteSearchableItemsWithIdentifiers:ids completionHandler:^(NSError *error) {
        wailsSpotlightComplete(requestID, error);
    }];
  }
}

void spotlightDeleteDomains(unsigned long long requestID, const char *domainsJSON) {
  @autoreleasepool {
    NSArray *domains = wailsSpotlightParseJSON(domainsJSON);
    if (![domains isKindOfClass:[NSArray class]]) {
        spotlightOperationResult(requestID, "spotlight: malformed domain list");
        return;
    }
    [[CSSearchableIndex defaultSearchableIndex] deleteSearchableItemsWithDomainIdentifiers:domains completionHandler:^(NSError *error) {
        wailsSpotlightComplete(requestID, error);
    }];
  }
}

void spotlightDeleteAllItems(unsigned long long requestID) {
    [[CSSearchableIndex defaultSearchableIndex] deleteAllSearchableItemsWithCompletionHandler:^(NSError *error) {
        wailsSpotlightComplete(requestID, error);
    }];
}

bool wailsSpotlightHandleUserActivity(NSUserActivity *activity) {
  @autoreleasepool {
    if (activity == nil) {
        return false;
    }
    NSString *type = activity.activityType;
    if (![type isEqualToString:CSSearchableItemActionType] && ![type isEqualToString:CSQueryContinuationActionType]) {
        return false;
    }
    // Only string-valued entries matter (the item identifier and the query
    // text); anything else is dropped so the payload always serialises.
    NSMutableDictionary *userInfo = [NSMutableDictionary dictionary];
    [activity.userInfo enumerateKeysAndObjectsUsingBlock:^(id key, id value, BOOL *stop) {
        if ([key isKindOfClass:[NSString class]] && [value isKindOfClass:[NSString class]]) {
            userInfo[key] = value;
        }
    }];
    NSData *data = [NSJSONSerialization dataWithJSONObject:userInfo options:0 error:nil];
    NSString *json = data != nil ? [[[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding] autorelease] : @"{}";
    return spotlightHandleContinueActivityC((char *)[type UTF8String], (char *)[json UTF8String]);
  }
}
