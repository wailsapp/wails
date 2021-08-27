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

import * as Platform from 'platform';

// IPC Listeners
var listeners = [];

/**
 * Adds a listener to IPC messages
 * @param {function} callback
 */
export function AddIPCListener(callback) {
	listeners.push(callback);
}

/**
 * SendMessage sends the given message to the backend
 *
 * @param {string} message
 */
export function SendMessage(message) {

	// Call Platform specific invoke method
	Platform.SendMessage(message);

	// Also send to listeners
	if (listeners.length > 0) {
		for (var i = 0; i < listeners.length; i++) {
			listeners[i](message);
		}
	}
}
