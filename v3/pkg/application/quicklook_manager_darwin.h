//go:build darwin && !ios && !server

#import <Cocoa/Cocoa.h>

// Opens the shared Quick Look panel on the files in pathsJSON (a JSON array
// of absolute paths). Returns false if the panel could not be shown. Must
// run on the main thread.
bool quickLookPreview(const char *pathsJSON);

// Closes the shared Quick Look panel if it exists. Must run on the main thread.
void quickLookClosePreview(void);

// Reports whether the shared Quick Look panel exists and is visible. Must
// run on the main thread.
bool quickLookIsPreviewOpen(void);

// Reports whether QLThumbnailGenerator is available (macOS 10.15+).
bool quickLookThumbnailAvailable(void);

// Starts thumbnail generation; the result arrives through the Go callback
// quickLookThumbnailResult with the same requestID.
void quickLookThumbnail(unsigned long long requestID, const char *path, int width, int height, double scale, bool iconMode);
