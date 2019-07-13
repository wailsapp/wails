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

// Sends a log message to the backend with the given
// level + message
function sendLogMessage(level, message) {

	// Log Message
	const payload = {
		level: level,
		message: message,
	}
	SendMessage('log', payload)
}

export function Debug(message) {
	sendLogMessage('debug', message);
}

export function Info(message) {
	sendLogMessage('info', message);
}

export function Warning(message) {
	sendLogMessage('warning', message);
}

export function Error(message) {
	sendLogMessage('error', message);
}

export function Fatal(message) {
	sendLogMessage('fatal', message);
}

