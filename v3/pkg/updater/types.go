package updater

import (
	"context"
	"io"
	"time"
)

// State is the high-level lifecycle phase the Updater is currently in.
type State string

const (
	StateUnconfigured State = "unconfigured"
	StateIdle         State = "idle"
	StateChecking     State = "checking"
	StateUpToDate     State = "up-to-date"
	StateAvailable    State = "available"
	StateDownloading  State = "downloading"
	StateVerifying    State = "verifying"
	StateInstalling   State = "installing"
	StateReady        State = "ready"
	StateError        State = "error"
)

// Stage describes which phase of the update flow produced an error.
type Stage string

const (
	StageCheck    Stage = "check"
	StageDownload Stage = "download"
	StageVerify   Stage = "verify"
	StageInstall  Stage = "install"
)

// Release is what a Provider returns from Check when there is something newer.
type Release struct {
	Version      string         `json:"version"`
	Channel      string         `json:"channel,omitempty"`
	Name         string         `json:"name,omitempty"`
	Notes        string         `json:"notes,omitempty"`
	PublishedAt  time.Time      `json:"publishedAt,omitempty"`
	Artifact     Artifact       `json:"artifact"`
	Verification *Verification  `json:"verification,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`

	// Provider is set by the Updater after Check; provider implementations
	// should leave it empty. Used to route a follow-up Download back to the
	// same source.
	Provider string `json:"provider,omitempty"`
}

// Artifact describes the file to download for a Release on the running platform.
type Artifact struct {
	Filename string `json:"filename"`
	Filetype string `json:"filetype,omitempty"`
	Size     int64  `json:"size,omitempty"`
	Platform string `json:"platform,omitempty"`
	Arch     string `json:"arch,omitempty"`
}

// Verification carries everything the Updater needs to authenticate the
// downloaded bytes. A Provider populates this from the release source; the
// Updater verifies against it using its configured trust root.
//
// Either Digest+DigestAlgo, Signature+SignatureAlgo, or both may be present.
// When both are present, both are checked; either failing fails the update.
//
// Signature verification always uses Config.PublicKey as the trust root. The
// release source has no say in which key authenticates it — that is the entire
// point of pinning a key out-of-band at build time. Releases that ship a
// Signature without a configured Config.PublicKey fail closed.
type Verification struct {
	DigestAlgo    string `json:"digestAlgo,omitempty"`    // "sha256", "sha512"
	Digest        []byte `json:"digest,omitempty"`        // raw digest bytes
	SignatureAlgo string `json:"signatureAlgo,omitempty"` // "ed25519", "ed25519ph", "ecdsa-p256"
	Signature     []byte `json:"signature,omitempty"`     // raw signature bytes
}

// Progress is the payload of EventDownloadProgress.
type Progress struct {
	Written  int64   `json:"written"`
	Total    int64   `json:"total"`
	Rate     float64 `json:"rate"` // bytes/sec smoothed over the last ~1s
	Provider string  `json:"provider,omitempty"`
}

// ErrorInfo is the payload of EventError.
type ErrorInfo struct {
	Stage    Stage  `json:"stage"`
	Message  string `json:"message"`
	Provider string `json:"provider,omitempty"`
}

// CheckRequest carries platform context to a Provider's Check method. The
// Updater fills in defaults from runtime.GOOS / runtime.GOARCH when the user
// does not supply them via Config.
type CheckRequest struct {
	CurrentVersion string
	Platform       string
	Arch           string
}

// Provider abstracts an update source. Implementations are typically tiny:
// resolve the next release for the running platform, then stream the bytes.
// Everything else (verification, atomic write, swap, restart, events) is the
// Updater's job.
type Provider interface {
	// Name identifies the provider in logs and in the `provider` field of
	// progress / error event payloads. Should be stable and short
	// (e.g. "github", "keygen.sh", "appcast").
	Name() string

	// Check returns the available upgrade for the running platform. Return
	// (nil, nil) when the source confirms the caller is already up to date.
	// Errors here put the Updater into the fallback chain.
	Check(ctx context.Context, req CheckRequest) (*Release, error)

	// Download streams the artifact bytes for r to dst. Providers should
	// invoke onProgress periodically; it is never nil.
	Download(ctx context.Context, r *Release, dst io.Writer, onProgress func(written, total int64)) error
}
