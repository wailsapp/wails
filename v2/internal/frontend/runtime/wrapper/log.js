/*
 _       __      _ __    
| |     / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  ) 
|__/|__/\__,_/_/_/____/  
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

/* jshint esversion: 6 */


/**
 * Log the given trace message with the backend
 *
 * @export
 * @param {string} message
 */
export function LogTrace(message) {
    window.runtime.LogTrace(message);
}

/**
 * Log the given debug message with the backend
 *
 * @export
 * @param {string} message
 */
export function LogDebug(message) {
    window.runtime.LogDebug(message);
}

/**
 * Log the given info message with the backend
 *
 * @export
 * @param {string} message
 */
export function LogInfo(message) {
    window.runtime.LogInfo(message);
}

/**
 * Log the given warning message with the backend
 *
 * @export
 * @param {string} message
 */
export function LogWarning(message) {
    window.runtime.LogWarning(message);
}

/**
 * Log the given error message with the backend
 *
 * @export
 * @param {string} message
 */
export function LogError(message) {
    window.runtime.LogError(message);
}

/**
 * Log the given fatal message with the backend
 *
 * @export
 * @param {string} message
 */
export function LogFatal(message) {
    window.runtime.LogFatal(message);
}
