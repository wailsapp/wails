//go:build darwin && !ios && !server

#ifndef browser_manager_darwin_h
#define browser_manager_darwin_h

// wailsWorkspaceOpenWith opens path with the application named by app (a
// bundle identifier or a .app path). The outcome is reported through
// wailsCompletionCallback(id, message) with a NULL message on success. Call
// on the main thread.
void wailsWorkspaceOpenWith(unsigned long long id, const char *path, const char *app);

// wailsWorkspaceApplicationsForFile returns a JSON array of
// {"name","bundleID","path"} for the applications able to open path, the
// default handler first. Caller frees.
char* wailsWorkspaceApplicationsForFile(const char *path);

// wailsWorkspaceActivateApplication activates the running application with
// the bundle identifier. Returns a strdup'd error message the caller frees,
// or NULL on success. Call on the main thread.
char* wailsWorkspaceActivateApplication(const char *bundleID);

#endif /* browser_manager_darwin_h */
