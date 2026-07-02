//go:build windows

package application

// WebView2 process-failure recovery (#5701 / #5705).
//
// On mixed-DPI / dual-GPU systems, moving a window between monitors can kill a
// WebView2 helper process. Field logs show the sequence: a GPU-process exit,
// followed up to ~25 seconds later by the controller "wedging" — every COM call
// (ExecuteScript, SetBounds, Focus, …) permanently returning
// ERROR_INVALID_STATE ("The group or resource is not in the correct state…").
// Once wedged, no in-place call can revive the controller, because the calls
// themselves fail: the only recovery is a new controller.
//
// This file implements a staged recovery ladder driven by two signals:
//
//  1. The ProcessFailed event (processFailed): every event schedules a
//     kind-appropriate action — re-navigate for a dead renderer, delayed
//     health probes for a dead GPU process (the wedge can lag the event),
//     and an unconditional controller rebuild for a dead browser process.
//  2. A low-frequency health watchdog (startWebviewHealthWatchdog): a cheap
//     controller probe that catches any wedge whose ProcessFailed event was
//     missed or never fired.
//
// The rebuild (rebuildWebview) swaps in a fresh edge.Chromium and re-runs
// setupChromium on the existing HWND, which re-embeds, re-applies settings and
// re-navigates to the start URL. The frontend reloads from scratch — the same
// as a manual refresh — which is the accepted cost of reviving a dead UI.
//
// All functions in this file run on the main/UI thread unless noted otherwise.

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/wailsapp/wails/v3/internal/assetserver"
	"github.com/wailsapp/wails/webview2/pkg/edge"
)

// webviewRecoveryAction selects how recoverWebview responds to a scheduled
// recovery check.
type webviewRecoveryAction int

const (
	// webviewRecoveryProbe checks controller health. Used for GPU/utility
	// process exits, where WebView2 usually restarts the helper process on
	// its own and no action is needed. An unhealthy probe is confirmed with a
	// second strike before the disruptive rebuild, because controller calls
	// can fail TRANSIENTLY while the browser process reconfigures during a
	// mixed-DPI drag (#5544) — only a persistent failure is a wedge.
	webviewRecoveryProbe webviewRecoveryAction = iota

	// webviewRecoveryProbeConfirm is the second strike of a failed probe:
	// still unhealthy 2s later → rebuild; recovered → done.
	webviewRecoveryProbeConfirm

	// webviewRecoveryRenavigate re-navigates to the start URL when the
	// controller is healthy (dead renderer, live browser process); an
	// unhealthy controller goes through the confirm strike like a probe.
	webviewRecoveryRenavigate

	// webviewRecoveryRenavigateConfirm is the second strike of a renavigate
	// that found the controller unhealthy: still unhealthy → rebuild;
	// recovered → perform the deferred re-navigation.
	webviewRecoveryRenavigateConfirm

	// webviewRecoveryRebuild rebuilds the controller unconditionally. Used for
	// a dead browser process, where the controller is documented to be
	// unrecoverable (COREWEBVIEW2_PROCESS_FAILED_KIND_BROWSER_PROCESS_EXITED)
	// — no probe involved, so no transient-failure ambiguity.
	webviewRecoveryRebuild
)

// webviewRecoveryConfirmDelay spaces the two strikes of an unhealthy probe.
// Long enough for a transient busy state to clear, short enough that a real
// wedge is still recovered within a few seconds of its first detection.
const webviewRecoveryConfirmDelay = 2 * time.Second

// webviewRuntimeVersion is the WebView2 runtime version string, recorded by
// setupChromium and included in ProcessFailed diagnostics so field logs can
// correlate process deaths with runtime rollouts. Main-thread only.
var webviewRuntimeVersion string

// webviewRebuildActiveGlobal serializes rebuilds ACROSS windows. Embed pumps a
// nested thread-wide message loop (GetMessage with hwnd=0), which reentrantly
// dispatches other windows' queued InvokeAsync callbacks — in the field
// scenario a GPU-process death wedges every window at once, and without this
// flag window B's rebuild would start nested inside window A's Embed pump,
// compounding worst-case main-thread block time. While set, other windows'
// recovery checks reschedule themselves and run after the active rebuild
// finishes. Main-thread only, so no lock is needed.
var webviewRebuildActiveGlobal bool

const (
	// webviewHealthProbeInterval is the watchdog period. Two consecutive
	// failures trigger a rebuild, so worst-case detection latency for a wedge
	// with no ProcessFailed event is ~2 intervals.
	webviewHealthProbeInterval = 15 * time.Second

	// webviewRebuildWindow / webviewRebuildMaxAttempts cap rebuild frequency:
	// at most webviewRebuildMaxAttempts rebuilds per window within
	// webviewRebuildWindow. The cap breaks pathological rebuild loops (e.g. a
	// broken WebView2 runtime that wedges immediately after every rebuild)
	// while still allowing occasional legitimate recoveries — entries expire,
	// so a dock/undock event hours later gets a fresh budget.
	webviewRebuildWindow      = 5 * time.Minute
	webviewRebuildMaxAttempts = 3
)

// processFailed is the WebView2 ProcessFailedCallback and the entry point of
// the recovery ladder. It logs the failure (throttled) and schedules the
// kind-appropriate recovery. Scheduling is NOT throttled: every event matters
// — field logs show a gpu-process-exited followed ~25s later by the actual
// wedge, and a time-only throttle would have swallowed the second event.
// Runs on the main thread (WebView2 delivers events on the UI thread).
func (w *windowsWebviewWindow) processFailed(_ *edge.ICoreWebView2, args *edge.ICoreWebView2ProcessFailedEventArgs) {
	var kind edge.COREWEBVIEW2_PROCESS_FAILED_KIND = edge.COREWEBVIEW2_PROCESS_FAILED_KIND_UNKNOWN_PROCESS_EXITED
	if args != nil {
		if k, err := args.GetProcessFailedKind(); err == nil {
			kind = k
		}
	}
	kindName := processFailedKindName(kind)

	// Throttle the stack-trace log PER KIND: RENDER_PROCESS_UNRESPONSIVE
	// re-fires every few seconds while hung, but a different kind (e.g. the
	// browser process dying after the GPU process — observed in the field with
	// a ~25s gap) must never be swallowed, and kinds alternating within the
	// window must not defeat the throttle by resetting each other.
	now := time.Now()
	if last, seen := w.processFailedLogAt[kindName]; !seen || now.Sub(last) >= 30*time.Second {
		if w.processFailedLogAt == nil {
			w.processFailedLogAt = make(map[string]time.Time, 2)
		}
		w.processFailedLogAt[kindName] = now
		globalApplication.error("WebView2 process failed (kind=%s, window=%d, runtime=%s) — scheduling recovery (#5701)\n%s",
			kindName, w.parent.id, webviewRuntimeVersion, debug.Stack())
	}

	switch kind {
	case edge.COREWEBVIEW2_PROCESS_FAILED_KIND_BROWSER_PROCESS_EXITED:
		// The controller is documented unrecoverable: rebuild. The short delay
		// moves the rebuild (which pumps messages in Embed) out of this COM
		// event handler.
		w.scheduleWebviewRecovery(200*time.Millisecond, kindName, webviewRecoveryRebuild)
	case edge.COREWEBVIEW2_PROCESS_FAILED_KIND_RENDER_PROCESS_EXITED:
		// WebView2 replaces the dead renderer with an error page; re-navigate
		// to restore the app (equivalent to the documented Reload recovery).
		w.scheduleWebviewRecovery(200*time.Millisecond, kindName, webviewRecoveryRenavigate)
	case edge.COREWEBVIEW2_PROCESS_FAILED_KIND_RENDER_PROCESS_UNRESPONSIVE:
		// Log only: the renderer may recover on its own, and auto-reloading a
		// merely-slow page would destroy frontend state. If it escalates to
		// RENDER_PROCESS_EXITED the case above handles it.
	case edge.COREWEBVIEW2_PROCESS_FAILED_KIND_GPU_PROCESS_EXITED:
		// WebView2 restarts the GPU process automatically and usually
		// recovers. But on dual-GPU monitor transitions the controller can
		// wedge long after the event (observed ~25s), so probe on a widening
		// schedule; the watchdog backstops anything later still.
		w.scheduleWebviewRecovery(1*time.Second, kindName, webviewRecoveryProbe)
		w.scheduleWebviewRecovery(5*time.Second, kindName, webviewRecoveryProbe)
		w.scheduleWebviewRecovery(30*time.Second, kindName, webviewRecoveryProbe)
	default:
		// Frame/utility/sandbox/unknown process exits normally self-heal; a
		// single probe confirms the controller survived.
		w.scheduleWebviewRecovery(2*time.Second, kindName, webviewRecoveryProbe)
	}
}

// processFailedKindName gives the COREWEBVIEW2_PROCESS_FAILED_KIND enum a
// human-readable label for the diagnostic log.
func processFailedKindName(k edge.COREWEBVIEW2_PROCESS_FAILED_KIND) string {
	switch k {
	case edge.COREWEBVIEW2_PROCESS_FAILED_KIND_BROWSER_PROCESS_EXITED:
		return "browser-process-exited"
	case edge.COREWEBVIEW2_PROCESS_FAILED_KIND_RENDER_PROCESS_EXITED:
		return "render-process-exited"
	case edge.COREWEBVIEW2_PROCESS_FAILED_KIND_RENDER_PROCESS_UNRESPONSIVE:
		return "render-process-unresponsive"
	case edge.COREWEBVIEW2_PROCESS_FAILED_KIND_FRAME_RENDER_PROCESS_EXITED:
		return "frame-render-process-exited"
	case edge.COREWEBVIEW2_PROCESS_FAILED_KIND_UTILITY_PROCESS_EXITED:
		return "utility-process-exited"
	case edge.COREWEBVIEW2_PROCESS_FAILED_KIND_GPU_PROCESS_EXITED:
		return "gpu-process-exited"
	default:
		return fmt.Sprintf("kind-%d", uint32(k))
	}
}

// scheduleWebviewRecovery runs recoverWebview on the main thread after delay.
// The timer goroutine hop is required even though the callers are already on
// the main thread: dispatchOnMainThread executes inline when on-thread, and
// the rebuild must not run inside a WebView2 COM event handler (Embed pumps a
// nested message loop). Safe to call from any thread.
func (w *windowsWebviewWindow) scheduleWebviewRecovery(delay time.Duration, reason string, action webviewRecoveryAction) {
	time.AfterFunc(delay, func() {
		InvokeAsync(func() {
			w.recoverWebview(reason, action)
		})
	})
}

// recoverWebview executes one scheduled recovery check. Main thread only.
func (w *windowsWebviewWindow) recoverWebview(reason string, action webviewRecoveryAction) {
	if w.webviewRebuildInProgress || w.hwnd == 0 || w.parent.isDestroyed() || globalApplication.performingShutdown {
		return
	}
	if webviewRebuildActiveGlobal {
		// Another window is mid-rebuild and its Embed pump dispatched us —
		// retry after it finishes so rebuilds never nest across windows.
		w.scheduleWebviewRecovery(webviewRecoveryConfirmDelay, reason, action)
		return
	}
	healthy := w.webviewControllerHealthy()
	switch action {
	case webviewRecoveryRebuild:
		w.rebuildWebview(reason)

	case webviewRecoveryProbe, webviewRecoveryRenavigate:
		if !healthy {
			// First strike. Could be a transient busy state (#5544), so
			// confirm before rebuilding; a renavigate keeps its intent
			// through the confirm.
			confirm := webviewRecoveryProbeConfirm
			if action == webviewRecoveryRenavigate {
				confirm = webviewRecoveryRenavigateConfirm
			}
			globalApplication.warning("WebView2 controller probe failed after %s — confirming in %s before rebuild (#5701)", reason, webviewRecoveryConfirmDelay)
			w.scheduleWebviewRecovery(webviewRecoveryConfirmDelay, reason, confirm)
			return
		}
		if action == webviewRecoveryRenavigate {
			globalApplication.warning("WebView2 renderer died (%s) — re-navigating to the start page (#5701)", reason)
			w.renavigateWebview()
		}

	case webviewRecoveryProbeConfirm, webviewRecoveryRenavigateConfirm:
		if !healthy {
			// Second strike: the failure persisted across the confirm delay,
			// so this is a real wedge, not a transient.
			w.rebuildWebview(reason)
			return
		}
		if action == webviewRecoveryRenavigateConfirm {
			globalApplication.warning("WebView2 renderer died (%s) — re-navigating to the start page (#5701)", reason)
			w.renavigateWebview()
		}
	}
}

// webviewControllerHealthy reports whether the WebView2 controller still
// accepts COM calls. NotifyParentWindowPositionChanged is the probe: it is
// cheap, side-effect-free on a healthy controller (WM_MOVE issues it
// routinely), and returns ERROR_INVALID_STATE once the controller has wedged.
// A controller that is still initialising reports healthy — that is start-up,
// not a wedge. Main thread only.
func (w *windowsWebviewWindow) webviewControllerHealthy() bool {
	if w.chromium == nil || !w.chromium.IsReady() {
		return true
	}
	return w.chromium.NotifyParentWindowPositionChanged() == nil
}

// renavigateWebview points the existing (healthy) controller back at the
// application content, mirroring the navigation tail of setupChromium. Used
// when only the render process died and WebView2 replaced the page with its
// internal error page.
func (w *windowsWebviewWindow) renavigateWebview() {
	// IsReady guarantees chromium.webview is non-nil (Navigate dereferences
	// it without a guard). A not-yet-initialised webview cannot need a
	// re-navigation — setupChromium navigates it when it comes up.
	if w.chromium == nil || !w.chromium.IsReady() {
		return
	}
	if w.parent.options.HTML != "" {
		w.chromium.NavigateToString(w.parent.options.HTML)
		return
	}
	startURL, err := assetserver.GetStartURL(w.parent.options.URL)
	if err != nil {
		globalApplication.error("WebView2 recovery: cannot resolve start URL: %s", err)
		return
	}
	w.chromium.Navigate(startURL)
}

// rebuildWebview replaces a dead WebView2 controller with a fresh one on the
// same HWND. setupChromium re-applies every setting from the window options
// against the new chromium (callbacks, permissions, settings, bounds) and
// re-navigates to the start URL; Embed pumps messages until the new controller
// is live, so the method returns with a working webview. Rate-capped to break
// rebuild loops. Main thread only.
func (w *windowsWebviewWindow) rebuildWebview(reason string) {
	if w.webviewRebuildInProgress || webviewRebuildActiveGlobal {
		// Direct callers (the watchdog) retry on their own cadence; scheduled
		// recovery paths reschedule in recoverWebview before reaching here.
		return
	}

	// Drop expired attempts, then enforce the cap.
	now := time.Now()
	recent := w.webviewRebuildTimes[:0]
	for _, t := range w.webviewRebuildTimes {
		if now.Sub(t) < webviewRebuildWindow {
			recent = append(recent, t)
		}
	}
	w.webviewRebuildTimes = recent
	if len(w.webviewRebuildTimes) >= webviewRebuildMaxAttempts {
		if now.Sub(w.lastRebuildSuppressedLogAt) >= time.Minute {
			w.lastRebuildSuppressedLogAt = now
			globalApplication.error("WebView2 rebuild suppressed: %d attempts within %s (%s) — the runtime appears unable to host a webview; an app restart is required (#5701)",
				webviewRebuildMaxAttempts, webviewRebuildWindow, reason)
		}
		return
	}
	w.webviewRebuildTimes = append(w.webviewRebuildTimes, now)

	w.webviewRebuildInProgress = true
	webviewRebuildActiveGlobal = true
	defer func() {
		w.webviewRebuildInProgress = false
		webviewRebuildActiveGlobal = false
	}()

	globalApplication.warning("WebView2 controller wedged (%s) — rebuilding controller, attempt %d (#5701)", reason, len(w.webviewRebuildTimes))

	// Detach the old chromium. Its browser process is usually gone, so no
	// further events arrive; nil the ProcessFailedCallback anyway so any
	// stragglers dispatched during teardown cannot schedule recovery against
	// the new controller, and mark it shutting down to silence its internal
	// paths. Hide() covers the alive-but-wedged case: if the old browser
	// process still owns a child HWND, hiding it stops it overlapping the new
	// controller's window (best-effort — fails with ERROR_INVALID_STATE when
	// the process is already dead, which is fine).
	if old := w.chromium; old != nil {
		old.ProcessFailedCallback = nil
		old.ShuttingDown()
		if old.GetController() != nil { // Hide dereferences the controller
			_ = old.Hide()
		}
	}

	w.chromium = edge.NewChromium()
	if globalApplication.options.ErrorHandler != nil {
		w.chromium.SetErrorCallback(globalApplication.options.ErrorHandler)
	}
	w.setupChromium()

	w.webviewHealthProbeFailures = 0
	// Warning level deliberately: recovery SUCCESS must reach hosts whose log
	// bridge forwards only Warn+ (the failure that triggered it is error-level,
	// so a missing success line would read as an unrecovered wedge).
	globalApplication.warning("WebView2 controller rebuilt after %s (window=%d) (#5701)", reason, w.parent.id)
}

// startWebviewHealthWatchdog begins a low-frequency controller health probe.
// It is the backstop for wedges that never deliver a usable ProcessFailed
// event: two consecutive probe failures (≥1 interval apart) trigger a rebuild.
// The chain stops rescheduling once the window is destroyed or the app shuts
// down. Called once from run() after the initial setupChromium.
func (w *windowsWebviewWindow) startWebviewHealthWatchdog() {
	var tick func()
	tick = func() {
		time.AfterFunc(webviewHealthProbeInterval, func() {
			InvokeAsync(func() {
				if w.hwnd == 0 || w.parent.isDestroyed() || globalApplication.performingShutdown {
					return
				}
				if w.webviewRebuildInProgress || w.webviewControllerHealthy() {
					w.webviewHealthProbeFailures = 0
				} else {
					w.webviewHealthProbeFailures++
					if w.webviewHealthProbeFailures == 1 {
						globalApplication.warning("WebView2 health probe failed — will rebuild the controller if the next probe fails (#5701)")
					} else {
						w.webviewHealthProbeFailures = 0
						w.rebuildWebview("health watchdog: controller stopped accepting calls")
					}
				}
				tick()
			})
		})
	}
	tick()
}
