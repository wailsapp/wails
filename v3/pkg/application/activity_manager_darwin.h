//go:build darwin && !ios && !server

#ifndef activity_manager_darwin_h
#define activity_manager_darwin_h

#include <stdbool.h>

// wailsActivityPublish creates an NSUserActivity from the UserActivity JSON
// ({"type","title","userInfo","webpageURL","eligibleForHandoff",
// "eligibleForSearch","eligibleForPrediction","keywords"}), stores it under
// id and makes it current. Call on the main thread.
void wailsActivityPublish(unsigned long long id, const char *json);

// wailsActivityUpdate replaces userInfo on the stored activity and marks it
// as needing save. Returns false when id is unknown. Call on the main thread.
bool wailsActivityUpdate(unsigned long long id, const char *userInfoJSON);

// wailsActivityInvalidate invalidates and forgets the stored activity. Call
// on the main thread.
void wailsActivityInvalidate(unsigned long long id);

// wailsActivityCurrentJSON returns the UserActivity JSON of the activity most
// recently made current through wailsActivityPublish, or "null" when none is
// live. Caller frees. Call on the main thread.
char* wailsActivityCurrentJSON(void);

#endif /* activity_manager_darwin_h */
