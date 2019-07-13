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

import { Debug } from './log';
import { SendMessage } from './ipc';

var callbacks = {};

// AwesomeRandom
function cryptoRandom() {
	var array = new Uint32Array(1);
	return window.crypto.getRandomValues(array)[0];
}

// LOLRandom
function basicRandom() {
	return Math.random() * 9007199254740991;
}

// Pick one based on browser capability
var randomFunc;
if (window.crypto) {
	randomFunc = cryptoRandom;
} else {
	randomFunc = basicRandom;
}


// Call sends a message to the backend to call the binding with the
// given data. A promise is returned and will be completed when the
// backend responds. This will be resolved when the call was successful
// or rejected if an error is passed back.
// There is a timeout mechanism. If the call doesn't respond in the given
// time (in milliseconds) then the promise is rejected.

export function Call(bindingName, data, timeout) {

	// Timeout infinite by default
	if (timeout == null || timeout == undefined) {
		timeout = 0;
	}

	// Create a promise
	return new Promise(function (resolve, reject) {

		// Create a unique callbackID
		var callbackID;
		do {
			callbackID = bindingName + '-' + randomFunc();
		} while (callbacks[callbackID]);

		// Set timeout
		if (timeout > 0) {
			var timeoutHandle = setTimeout(function () {
				reject(Error('Call to ' + bindingName + ' timed out. Request ID: ' + callbackID));
			}, timeout);
		}

		// Store callback
		callbacks[callbackID] = {
			timeoutHandle: timeoutHandle,
			reject: reject,
			resolve: resolve
		};

		try {
			const payload = {
				bindingName: bindingName,
				data: JSON.stringify(data),
			};

			// Make the call
			SendMessage('call', payload, callbackID);
		} catch (e) {
			// eslint-disable-next-line
			console.error(e);
		}
	});
}


// Called by the backend to return data to a previously called
// binding invocation
export function Callback(incomingMessage) {

	// Decode the message - Credit: https://stackoverflow.com/a/13865680
	incomingMessage = decodeURIComponent(incomingMessage.replace(/\s+/g, '').replace(/[0-9a-f]{2}/g, '%$&'));

	// Parse the message
	var message;
	try {
		message = JSON.parse(incomingMessage);
	} catch (e) {
		const error = `Invalid JSON passed to callback: ${e.message}. Message: ${incomingMessage}`;
		Debug(error);
		throw new Error(error);
	}
	var callbackID = message.callbackid;
	var callbackData = callbacks[callbackID];
	if (!callbackData) {
		const error = `Callback '${callbackID}' not registed!!!`;
		console.error(error); // eslint-disable-line
		throw new Error(error);
	}
	clearTimeout(callbackData.timeoutHandle);

	delete callbacks[callbackID];

	if (message.error) {
		return callbackData.reject(message.error);
	}
	return callbackData.resolve(message.data);
}

// systemCall is used to call wails methods from the frontend
export function SystemCall(method, data) {
	return Call('.wails.' + method, data);
}