//go:build windows

package edge

type COREWEBVIEW2_PROCESS_FAILED_KIND uint32

const (
	// Indicates that the browser process ended unexpectedly.  The WebView
	// automatically moves to the Closed state.  The app has to recreate a new
	// WebView to recover from this failure.
	COREWEBVIEW2_PROCESS_FAILED_KIND_BROWSER_PROCESS_EXITED = 0

	// Indicates that the main frame's render process ended unexpectedly.  A new
	// render process is created automatically and navigated to an error page.
	// You can use the `Reload` method to try to reload the page that failed.
	COREWEBVIEW2_PROCESS_FAILED_KIND_RENDER_PROCESS_EXITED = 1

	// Indicates that the main frame's render process is unresponsive.
	//
	// Note that this does not seem to work right now.
	// Does not fire for simple long running script case, the only related test
	// SitePerProcessBrowserTest::NoCommitTimeoutForInvisibleWebContents is
	// disabled.
	COREWEBVIEW2_PROCESS_FAILED_KIND_RENDER_PROCESS_UNRESPONSIVE = 2

	// Indicates that a frame-only render process ended unexpectedly. The process
	// exit does not affect the top-level document, only a subset of the
	// subframes within it. The content in these frames is replaced with an error
	// page in the frame.
	COREWEBVIEW2_PROCESS_FAILED_KIND_FRAME_RENDER_PROCESS_EXITED = 3

	// Indicates that a utility process ended unexpectedly.
	COREWEBVIEW2_PROCESS_FAILED_KIND_UTILITY_PROCESS_EXITED = 4

	// Indicates that a sandbox helper process ended unexpectedly.
	COREWEBVIEW2_PROCESS_FAILED_KIND_SANDBOX_HELPER_PROCESS_EXITED = 5

	// Indicates that the GPU process ended unexpectedly.
	COREWEBVIEW2_PROCESS_FAILED_KIND_GPU_PROCESS_EXITED = 6

	// Indicates that a PPAPI plugin process ended unexpectedly.
	COREWEBVIEW2_PROCESS_FAILED_KIND_PPAPI_PLUGIN_PROCESS_EXITED = 7

	// Indicates that a PPAPI plugin broker process ended unexpectedly.
	COREWEBVIEW2_PROCESS_FAILED_KIND_PPAPI_BROKER_PROCESS_EXITED = 8

	// Indicates that a process of unspecified kind ended unexpectedly.
	COREWEBVIEW2_PROCESS_FAILED_KIND_UNKNOWN_PROCESS_EXITED = 9
)
