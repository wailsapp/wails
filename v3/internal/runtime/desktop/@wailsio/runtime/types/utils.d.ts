/**
 * Logs a message to the console with custom formatting.
 *
 * @param message - The message to be logged.
 */
export declare function debugLog(message: any): void;
/**
 * Checks whether the webview supports the {@link MouseEvent#buttons} property.
 * Looking at you macOS High Sierra!
 */
export declare function canTrackButtons(): boolean;
/**
 * Checks whether the browser supports removing listeners by triggering an AbortSignal
 * (see https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/addEventListener#signal).
 */
export declare function canAbortListeners(): boolean;
/**
 * Resolves the closest HTMLElement ancestor of an event's target.
 */
export declare function eventTarget(event: Event): HTMLElement;
export declare function whenReady(callback: () => void): void;
