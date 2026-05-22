package updater

// Event names emitted by the Updater. Subscribe in Go via app.Event.On(name, ...)
// or in JavaScript via wails.Events.On(name, ...). Payload types are documented
// inline next to each constant.
const (
	// EventCheckStarted fires before a Check round trip. Payload: nil.
	EventCheckStarted = "updater:check-started"

	// EventUpdateAvailable fires when Check returns a newer release. Payload: *Release.
	EventUpdateAvailable = "updater:update-available"

	// EventNoUpdate fires when Check confirms the caller is up to date. Payload: nil.
	EventNoUpdate = "updater:no-update"

	// EventDownloadStarted fires when the Updater begins streaming bytes from
	// a provider. Payload: *Release.
	EventDownloadStarted = "updater:download-started"

	// EventDownloadProgress fires periodically during download (~10/sec). Payload: Progress.
	EventDownloadProgress = "updater:download-progress"

	// EventDownloadComplete fires once all bytes are on disk and the file has
	// been closed, but BEFORE verification. Payload: *Release.
	EventDownloadComplete = "updater:download-complete"

	// EventVerifying fires when the Updater begins verifying the downloaded
	// artifact. Payload: *Release.
	EventVerifying = "updater:verifying"

	// EventInstalling fires when the Updater begins swapping the binary.
	// Payload: *Release.
	EventInstalling = "updater:installing"

	// EventUpdateReady fires when an update is installed and a restart is
	// pending. Payload: *Release.
	EventUpdateReady = "updater:update-ready"

	// EventError fires whenever any stage fails. Payload: ErrorInfo.
	EventError = "updater:error"

	// EventMeta fires once per session before the first state-snapshot
	// replay, carrying host-side context the page can't derive from any
	// Release: the version currently running, and the version the user
	// has marked skipped (or "" if none). Payload: Meta.
	EventMeta = "updater:meta"
)

// Meta is the payload of EventMeta — host-side context the default window
// template uses to render the "from" version in the update pill and the
// "v1.2.3 · This is the latest version" pill in the up-to-date state.
type Meta struct {
	CurrentVersion string `json:"currentVersion"`
	SkippedVersion string `json:"skippedVersion,omitempty"`
}
