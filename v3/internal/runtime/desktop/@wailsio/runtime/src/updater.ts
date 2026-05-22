/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

/**
 * Updater event name constants.
 *
 * Use these instead of hard-coding string literals when subscribing to
 * updater events from JavaScript:
 *
 *     import { Events, Updater } from "@wailsio/runtime";
 *
 *     Events.On(Updater.Events.UpdateAvailable, (e) => {
 *         console.log("update found:", e.data.version);
 *     });
 *
 *     Events.On(Updater.Events.DownloadProgress, (e) => {
 *         const p = e.data;
 *         console.log(`${p.written} / ${p.total} bytes`);
 *     });
 *
 * Mirrors the Go-side constants in `pkg/updater/events.go` and the
 * user-action constants in `pkg/updater/window_lifecycle.go`. Any
 * changes here must stay in sync with those files — there's an
 * integration test that asserts the strings match.
 */
export const Events = Object.freeze({
    /** A Check round-trip is starting. Payload: null. */
    CheckStarted: "wails:updater:CheckStarted",
    /** Check found a newer release. Payload: Release. */
    UpdateAvailable: "wails:updater:UpdateAvailable",
    /** Check confirmed the caller is up to date. Payload: null. */
    NoUpdate: "wails:updater:NoUpdate",
    /** Download is starting. Payload: Release. */
    DownloadStarted: "wails:updater:DownloadStarted",
    /** Periodic progress tick during download (~10 Hz). Payload: Progress. */
    DownloadProgress: "wails:updater:DownloadProgress",
    /** All bytes are on disk, but verification has not yet started. Payload: Release. */
    DownloadComplete: "wails:updater:DownloadComplete",
    /** Signature / digest verification has started. Payload: Release. */
    Verifying: "wails:updater:Verifying",
    /** The Updater is swapping the binary into place. Payload: Release. */
    Installing: "wails:updater:Installing",
    /** Update is staged and a restart is pending. Payload: Release. */
    UpdateReady: "wails:updater:UpdateReady",
    /** Something failed. Payload: ErrorInfo { stage, message, provider }. */
    Error: "wails:updater:Error",
    /** Host-side context delivered once per session. Payload: Meta { currentVersion, skippedVersion }. */
    Meta: "wails:updater:Meta",

    /** Sub-namespace: user-action events that the UI emits BACK to the host. */
    User: Object.freeze({
        /** User clicked Install on an Available update. */
        Install: "wails:updater:user:Install",
        /** User clicked Restart & Apply on a Ready update. */
        Restart: "wails:updater:user:Restart",
        /** User clicked Skip This Version. */
        Skip: "wails:updater:user:Skip",
        /** User clicked Remind Me Later. */
        Remind: "wails:updater:user:Remind",
        /** User clicked Close / Cancel. */
        Cancel: "wails:updater:user:Cancel",
    }),

    /** Sub-namespace: framework-internal events the UI emits to coordinate
     *  with the host. Most app code can ignore these. */
    Window: Object.freeze({
        /** The window finished loading and asks the host to replay current state. */
        Ready: "wails:updater:window:Ready",
    }),
});
