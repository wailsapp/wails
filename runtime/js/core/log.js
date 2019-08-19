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

import { SendMessage } from './ipc';

/**
 * Sends a log message to the backend with the given level + message
 *
 * @param {string} level
 * @param {string} message
 */
function sendLogMessage(level, message) {

	// Log Message
	const payload = {
		level: level,
		message: message,
	};
	SendMessage('log', payload);
}

/**
 * Log the given debug message with the backend
 *
 * @export
 * @param {string} message
 */
export function Debug(message) {
	sendLogMessage('debug', message);
}

/**
 * Log the given info message with the backend
 *
 * @export
 * @param {string} message
 */
export function Info(message) {
	sendLogMessage('info', message);
}

/**
 * Log the given warning message with the backend
 *
 * @export
 * @param {string} message
 */
export function Warning(message) {
	sendLogMessage('warning', message);
}

/**
 * Log the given error message with the backend
 *
 * @export
 * @param {string} message
 */
export function Error(message) {
	sendLogMessage('error', message);
}

/**
 * Log the given fatal message with the backend
 *
 * @export
 * @param {string} message
 */
export function Fatal(message) {
	sendLogMessage('fatal', message);
}

