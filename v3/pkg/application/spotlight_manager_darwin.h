//go:build darwin && !ios && !server

#import <Cocoa/Cocoa.h>

// Reports whether CoreSpotlight indexing is available.
bool spotlightIsAvailable(void);

// Index and delete operations. Each completes asynchronously through the Go
// callback spotlightOperationResult with the same requestID. itemsJSON is a
// JSON array of items (see spotlightWireItem in spotlight_manager.go) and
// idsJSON a JSON array of strings.
void spotlightIndexItems(unsigned long long requestID, const char *itemsJSON);
void spotlightDeleteItems(unsigned long long requestID, const char *idsJSON);
void spotlightDeleteDomains(unsigned long long requestID, const char *domainsJSON);
void spotlightDeleteAllItems(unsigned long long requestID);

// Optional native entry point for application:continueUserActivity:.
// Returns true when the activity was a Spotlight item selection or search
// continuation and a Go OnOpen handler consumed it. The delegate does not
// need it: OnOpen registers itself with the ActivityManager, whose
// activityContinue export the delegate already calls. It exists for native
// code that holds an NSUserActivity outside that path:
//
//	if (wailsSpotlightHandleUserActivity(userActivity)) { return YES; }
bool wailsSpotlightHandleUserActivity(NSUserActivity *activity);
