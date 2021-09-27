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
	if (window.wailsbridge) {
		window.wailsbridge.websocket.send(message);
	} else {
		window.external.invoke(message);
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
 * @param {string} type
 * @param {Object} payload
 * @param {string=} callbackID
 */
export function SendMessage(type, payload, callbackID) {
	const message = {
		type,
		callbackID,
		payload
	};

	Invoke(JSON.stringify(message));
}
