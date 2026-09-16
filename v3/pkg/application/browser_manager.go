package application

import (
	"errors"
	"strings"

	"github.com/wailsapp/wails/v3/internal/browser"
)

// BrowserManager manages browser-related operations
type BrowserManager struct {
	app *App
}

// newBrowserManager creates a new BrowserManager instance
func newBrowserManager(app *App) *BrowserManager {
	return &BrowserManager{
		app: app,
	}
}

// OpenURL opens a URL in the default browser
func (bm *BrowserManager) OpenURL(url string) error {
	return browser.OpenURL(url)
}

// OpenFile opens a file in the default browser
func (bm *BrowserManager) OpenFile(path string) error {
	return browser.OpenFile(path)
}

// AppInfo identifies an installed application.
type AppInfo struct {
	// Name is the user-visible application name ("TextEdit").
	Name string `json:"name"`
	// BundleID is the bundle identifier ("com.apple.TextEdit"); empty on
	// platforms without bundles.
	BundleID string `json:"bundleID"`
	// Path is the application bundle or executable path.
	Path string `json:"path"`
}

// ErrOpenWithUnsupported is returned by OpenWith when the platform cannot
// open a file with a specific application.
var ErrOpenWithUnsupported = errors.New("browser: opening with a specific application is not supported on this platform")

// ErrApplicationNotFound is returned when an application cannot be resolved
// from a bundle identifier or path.
var ErrApplicationNotFound = errors.New("browser: application not found")

// ErrApplicationNotRunning is returned by ActivateApplication when no
// running process matches the bundle identifier.
var ErrApplicationNotRunning = errors.New("browser: application is not running")

// OpenWith opens path with a specific application instead of the default
// handler. app is a bundle identifier ("com.apple.TextEdit") or the path
// to an application bundle ("/System/Applications/TextEdit.app").
//
// macOS: NSWorkspace openURLs:withApplicationAtURL:configuration:
// completionHandler:; bundle identifiers are resolved through
// URLForApplicationWithBundleIdentifier:. OpenWith returns once the
// application has accepted the request.
//
// Other platforms: app must name an executable on PATH or be an executable
// path; it is started with path as its only argument. Otherwise
// ErrOpenWithUnsupported.
func (bm *BrowserManager) OpenWith(path string, app string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("browser: path is required")
	}
	if strings.TrimSpace(app) == "" {
		return errors.New("browser: application is required")
	}
	return platformOpenWith(path, app)
}

// ApplicationsForFile lists the installed applications able to open path,
// with the default handler first.
//
// macOS: NSWorkspace URLsForApplicationsToOpenURL: (macOS 12+), with
// LSCopyApplicationURLsForURL on older systems.
//
// Other platforms: empty.
func (bm *BrowserManager) ApplicationsForFile(path string) []AppInfo {
	if strings.TrimSpace(path) == "" {
		return []AppInfo{}
	}
	apps := platformApplicationsForFile(path)
	if apps == nil {
		apps = []AppInfo{}
	}
	return apps
}

// ActivateApplication brings a running application to the front by bundle
// identifier.
//
// macOS: NSRunningApplication activateFromApplication:options: (macOS 14+,
// where activation is cooperative and only the active application can hand
// focus to another) or activateWithOptions: on older systems, with all
// windows. ErrApplicationNotRunning when nothing matches.
//
// Other platforms: ErrApplicationNotRunning.
func (bm *BrowserManager) ActivateApplication(bundleID string) error {
	if strings.TrimSpace(bundleID) == "" {
		return errors.New("browser: bundle identifier is required")
	}
	return platformActivateApplication(bundleID)
}
