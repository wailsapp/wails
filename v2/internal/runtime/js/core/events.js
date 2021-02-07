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

import { Error } from './log';
import { SendMessage } from 'ipc';

// Defines a single listener with a maximum number of times to callback

/**
 * The Listener class defines a listener! :-)
 *
 * @class Listener
 */
class Listener {
	/**
	 * Creates an instance of Listener.
	 * @param {function} callback
	 * @param {number} maxCallbacks
	 * @memberof Listener
	 */
	constructor(callback, maxCallbacks) {
		// Default of -1 means infinite
		maxCallbacks = maxCallbacks || -1;
		// Callback invokes the callback with the given data
		// Returns true if this listener should be destroyed
		this.Callback = (data) => {
			callback.apply(null, data);
			// If maxCallbacks is infinite, return false (do not destroy)
			if (maxCallbacks === -1) {
				return false;
			}
			// Decrement maxCallbacks. Return true if now 0, otherwise false
			maxCallbacks -= 1;
			return maxCallbacks === 0;
		};
	}
}

let eventListeners = {};

/**
 * Registers an event listener that will be invoked `maxCallbacks` times before being destroyed
 *
 * @export
 * @param {string} eventName
 * @param {function} callback
 * @param {number} maxCallbacks
 */
export function OnMultiple(eventName, callback, maxCallbacks) {
	eventListeners[eventName] = eventListeners[eventName] || [];
	const thisListener = new Listener(callback, maxCallbacks);
	console.log('Pushing event listener: ' + eventName);
	eventListeners[eventName].push(thisListener);
}

/**
 * Registers an event listener that will be invoked every time the event is emitted
 *
 * @export
 * @param {string} eventName
 * @param {function} callback
 */
export function On(eventName, callback) {
	OnMultiple(eventName, callback);
}

/**
 * Registers listeners for when the system theme changes from light/dark. A bool is
 * sent to the listener, true if it is dark mode.
 *
 * @export
 * @param {function} callback
 */
export function OnThemeChange(callback) {
	On('wails:system:themechange', callback);
}

/**
 * Registers an event listener that will be invoked once then destroyed
 *
 * @export
 * @param {string} eventName
 * @param {function} callback
 */
export function Once(eventName, callback) {
	OnMultiple(eventName, callback, 1);
}

function notifyListeners(eventData) {

	// Get the event name
	let eventName = eventData.name;

	// Check if we have any listeners for this event
	if (eventListeners[eventName]) {

		// Keep a list of listener indexes to destroy
		const newEventListenerList = eventListeners[eventName].slice();

		// Iterate listeners
		for (let count = 0; count < eventListeners[eventName].length; count += 1) {

			// Get next listener
			const listener = eventListeners[eventName][count];

			let data = eventData.data;

			// Do the callback
			const destroy = listener.Callback(data);
			if (destroy) {
				// if the listener indicated to destroy itself, add it to the destroy list
				newEventListenerList.splice(count, 1);
			}
		}

		// Update callbacks with new list of listeners
		eventListeners[eventName] = newEventListenerList;
	}
}

/**
 * Notify informs frontend listeners that an event was emitted with the given data
 *
 * @export
 * @param {string} notifyMessage - encoded notification message

 */
export function Notify(notifyMessage) {

	// Parse the message
	var message;
	try {
		message = JSON.parse(notifyMessage);
	} catch (e) {
		const error = 'Invalid JSON passed to Notify: ' + notifyMessage;
		throw new Error(error);
	}

	notifyListeners(message);
}

/**
 * Emit an event with the given name and data
 *
 * @export
 * @param {string} eventName
 */
export function Emit(eventName) {

	const payload = {
		name: eventName,
		data: [].slice.apply(arguments).slice(1),
	};

	// Notify JS listeners
	notifyListeners(payload);

	// Notify Go listeners
	SendMessage('Ej' + JSON.stringify(payload));

}