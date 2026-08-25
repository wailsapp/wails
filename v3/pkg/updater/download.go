package updater

import (
	"context"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"time"
)

// download streams the artifact for rel from p into a temp file, computing
// the verification digest while the bytes flow past. Returns the temp file
// path and the temp directory enclosing it; the caller is responsible for
// verifying then renaming or removing the directory (RemoveAll on dir tears
// the file down with it).
//
// On any error every artifact created by this function (file + enclosing
// directory) is removed before returning.
//
// Progress is emitted via the application event bus at ~10/sec.
func (u *Updater) download(ctx context.Context, p Provider, rel *Release) (path, dir string, err error) {
	u.transition(StateDownloading)
	u.host.Emit(EventDownloadStarted, rel)

	dir, err = os.MkdirTemp("", "wails-update-*")
	if err != nil {
		return "", "", fmt.Errorf("updater: temp dir: %w", err)
	}
	tmpPath := filepath.Join(dir, ".artifact")

	// Track success — on every error path we tear the directory down.
	success := false
	defer func() {
		if !success {
			_ = os.RemoveAll(dir)
		}
	}()

	f, err := os.Create(tmpPath)
	if err != nil {
		return "", "", fmt.Errorf("updater: create temp: %w", err)
	}
	closed := false
	defer func() {
		if !closed {
			_ = f.Close()
		}
	}()

	// Set up a hasher in parallel only when the release will need it.
	var hasher hash.Hash
	if rel.Verification != nil {
		algo := rel.Verification.DigestAlgo
		// Ed25519ph requires SHA-512 regardless of what DigestAlgo says.
		if rel.Verification.SignatureAlgo == "ed25519ph" {
			algo = "sha512"
		}
		hasher, err = digestHasher(algo)
		if err != nil {
			return "", "", err
		}
	}

	var dst io.Writer = f
	if hasher != nil {
		dst = io.MultiWriter(f, hasher)
	}

	prov := p.Name()
	emit := u.host.Emit
	tracker := newProgressTracker(emit, prov)

	progressFn := func(written, total int64) {
		tracker.tick(written, total)
	}

	wrappedDst := &countingWriter{w: dst, onWrite: tracker.add}

	if err := p.Download(ctx, rel, wrappedDst, progressFn); err != nil {
		_ = f.Close()
		closed = true
		return "", "", err
	}

	if err := f.Sync(); err != nil {
		_ = f.Close()
		closed = true
		return "", "", fmt.Errorf("updater: fsync: %w", err)
	}
	if err := f.Close(); err != nil {
		closed = true
		return "", "", fmt.Errorf("updater: close: %w", err)
	}
	closed = true

	// Final progress tick so subscribers see written == total.
	tracker.flush()

	u.host.Emit(EventDownloadComplete, rel)

	// Stash the digest for verify() so we don't re-read the file.
	if hasher != nil {
		u.mu.Lock()
		u.lastDigest = hasher.Sum(nil)
		u.mu.Unlock()
	}
	success = true
	return tmpPath, dir, nil
}

// verify runs the configured verification rules against the digest computed
// during download. It does NOT re-read the file from disk — the streaming
// hash is the authoritative computation.
func (u *Updater) verify(_ string, rel *Release) error {
	if rel.Verification == nil {
		return nil
	}
	u.mu.RLock()
	digest := u.lastDigest
	cfgKey := u.cfg.PublicKey
	u.mu.RUnlock()
	if digest == nil {
		return errors.New("updater: no digest available; download path bug")
	}
	return runVerification(digest, rel.Verification, cfgKey)
}

// progressTracker debounces progress callbacks from the provider and
// computes a smoothed bytes/sec rate for the EventDownloadProgress payload.
type progressTracker struct {
	emit     func(string, ...any) bool
	provider string

	written int64
	total   int64

	start    time.Time
	lastEmit time.Time
}

func newProgressTracker(emit func(string, ...any) bool, provider string) *progressTracker {
	now := time.Now()
	return &progressTracker{emit: emit, provider: provider, start: now, lastEmit: now}
}

// add is called by the countingWriter every Write to update the running total.
func (p *progressTracker) add(n int) {
	p.written += int64(n)
}

// tick is called by the provider's onProgress callback. The caller may pass
// a total; if so, we let it set our total even on the first tick. Throttled
// emits keep the event bus calm during large downloads.
func (p *progressTracker) tick(written, total int64) {
	if written > p.written {
		p.written = written
	}
	if total > 0 {
		p.total = total
	}
	if time.Since(p.lastEmit) < 100*time.Millisecond {
		return
	}
	p.emitNow()
}

func (p *progressTracker) flush() {
	p.emitNow()
}

func (p *progressTracker) emitNow() {
	p.lastEmit = time.Now()
	elapsed := p.lastEmit.Sub(p.start).Seconds()
	var rate float64
	if elapsed > 0 {
		rate = float64(p.written) / elapsed
	}
	p.emit(EventDownloadProgress, Progress{
		Written:  p.written,
		Total:    p.total,
		Rate:     rate,
		Provider: p.provider,
	})
}

// countingWriter wraps an io.Writer and invokes onWrite with the byte count
// for each write that passes through.
type countingWriter struct {
	w       io.Writer
	onWrite func(int)
}

func (c *countingWriter) Write(p []byte) (int, error) {
	n, err := c.w.Write(p)
	if n > 0 && c.onWrite != nil {
		c.onWrite(n)
	}
	return n, err
}
