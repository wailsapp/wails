/*
 _       __      _ __
| |     / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The lightweight framework for web-like apps
(c) Lea Anthony 2019-present
*/

/* jshint esversion: 6 */

import { SendMessage } from 'ipc';

/**
 * Sends a log message to the backend with the given level + message
 *
 * @param {string} level
 * @param {string} message
 */
function sendLogMessage(level, message) {

	// Log Message format:
	// l[type][message]
	SendMessage('L' + level + message);
}

/**
 * Log the given trace message with the backend
 *
 * @export
 * @param {string} message
 */
export function Trace(message) {
	sendLogMessage('T', message);
}

/**
 * Log the given message with the backend
 *
 * @export
 * @param {string} message
 */
export function Print(message) {
	sendLogMessage('P', message);
}

/**
 * Log the given debug message with the backend
 *
 * @export
 * @param {string} message
 */
export function Debug(message) {
	sendLogMessage('D', message);
}

/**
 * Log the given info message with the backend
 *
 * @export
 * @param {string} message
 */
export function Info(message) {
	sendLogMessage('I', message);
}

/**
 * Log the given warning message with the backend
 *
 * @export
 * @param {string} message
 */
export function Warning(message) {
	sendLogMessage('W', message);
}

/**
 * Log the given error message with the backend
 *
 * @export
 * @param {string} message
 */
export function Error(message) {
	sendLogMessage('E', message);
}

/**
 * Log the given fatal message with the backend
 *
 * @export
 * @param {string} message
 */
export function Fatal(message) {
	sendLogMessage('F', message);
}

/**
 * Sets the Log level to the given log level
 *
 * @export
 * @param {number} loglevel
 */
export function SetLogLevel(loglevel) {
	sendLogMessage('S', loglevel);
}

// Log levels
export const Level = {
	TRACE: 1,
	DEBUG: 2,
	INFO: 3,
	WARNING: 4,
	ERROR: 5,
};
