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
const listeners = [];

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
	if (PLATFORM === "windows") {
		window.chrome.webview.postMessage(message);
	} else if (PLATFORM === "darwin") {
		window.blah();
	} else {
		console.error("Unsupported Platform");
	}

	// Also send to listeners
	if (listeners.length > 0) {
		for (let i = 0; i < listeners.length; i++) {
			listeners[i](message);
		}
	}
}
