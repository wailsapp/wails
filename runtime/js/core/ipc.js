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

/**
 * Invoke sends the given message to the backend
 *
 * @param {string} message
 */
function Invoke(message) {
	if ( window.wailsbridge ) {
		window.wailsbridge.websocket.send(message);
	} else {
		window.external.invoke(message);
	}
}

/**
 * Sends a message to the backend based on the given type, payload and callbackID
 *
 * @export
 * @param {string} type
 * @param {string} payload
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
