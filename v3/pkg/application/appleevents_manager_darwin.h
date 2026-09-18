//go:build darwin && !ios && !server

#import <Cocoa/Cocoa.h>

// Registers the shared Wails handler object for the given event class and
// ID with NSAppleEventManager. Must run on the main thread.
void appleEventsRegister(const char *eventClass, const char *eventID);

// Completes a suspended event: fills its reply from replyJSON and resumes
// it. Must run on the main thread.
void appleEventsResume(void *suspensionID, const char *replyJSON);

// Sends an Apple Event to the application with the given bundle identifier
// and waits for the reply. directJSON is the wire form of the direct object
// (see appleevents_manager.go). Returns the reply as JSON, or NULL with
// *errorMessage set. The caller frees both strings.
char *appleEventsSend(const char *bundleID, const char *eventClass, const char *eventID,
                      const char *directJSON, char **errorMessage);
