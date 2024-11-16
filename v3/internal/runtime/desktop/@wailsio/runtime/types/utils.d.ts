/**
 * Logs a message to the console with custom formatting.
 * @param {string} message - The message to be logged.
 * @return {void}
 */
export function debugLog(message: string): void;
/**
 * Checks whether the browser supports removing listeners by triggering an AbortSignal
 * (see https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/addEventListener#signal)
 *
 * @return {boolean}
 */
export function canAbortListeners(): boolean;
export function whenReady(callback: any): void;
