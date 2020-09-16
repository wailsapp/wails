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
 * Invoke sends the given message to the backend
 *
 * @param {string} message
 */
function Invoke(message) {
	if (window.wailsbridge && window.wailsbridge.websocket) {
		window.wailsbridge.websocket.send(JSON.stringify(message));
	} else {
		console.log('Invoke called with: ' + message + ' but no runtime is available');
	}

	// Also send to listeners
	if (listeners.length > 0) {
		for (var i = 0; i < listeners.length; i++) {
			listeners[i](message);
		}
	}
}

/**
 * Sends a message to the backend based on the given type, payload and callbackID
 *
 * @export
 * @param {string} message
 */

export function SendMessage(message) {
	Invoke(message);
}
