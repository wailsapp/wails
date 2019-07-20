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
import { SendMessage } from './ipc';

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

var eventListeners = {};

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
 * Registers an event listener that will be invoked once then destroyed
 *
 * @export
 * @param {string} eventName
 * @param {function} callback
 */
export function Once(eventName, callback) {
	OnMultiple(eventName, callback, 1);
}

/**
 * Notify informs frontend listeners that an event was emitted with the given data
 *
 * @export
 * @param {string} eventName
 * @param {string} data
 */
export function Notify(eventName, data) {

	// Check if we have any listeners for this event
	if (eventListeners[eventName]) {

		// Keep a list of listener indexes to destroy
		const newEventListenerList = eventListeners[eventName].slice();

		// Iterate listeners
		for (let count = 0; count < eventListeners[eventName].length; count += 1) {

			// Get next listener
			const listener = eventListeners[eventName][count];

			// Parse data if we have it
			var parsedData = [];
			if (data) {
				try {
					parsedData = JSON.parse(data);
				} catch (e) {
					Error('Invalid JSON data sent to notify. Event name = ' + eventName);
				}
			}
			// Do the callback
			const destroy = listener.Callback(parsedData);
			if (destroy) {
				// if the listener indicated to destroy itself, add it to the destroy list
				newEventListenerList.splice(count, 1);
			}
		}

		// Update callbacks with new list of listners
		eventListeners[eventName] = newEventListenerList;
	}
}

/**
 * Emit an event with the given name and data
 *
 * @export
 * @param {string} eventName
 */
export function Emit(eventName) {

	// Calculate the data
	var data = JSON.stringify([].slice.apply(arguments).slice(1));

	// Notify backend
	const payload = {
		name: eventName,
		data: data,
	};
	SendMessage('event', payload);
}

// Callbacks for the heartbeat calls
const heartbeatCallbacks = {};

/**
 * Heartbeat emits the event `eventName`, every `timeInMilliseconds` milliseconds until 
 * the event is acknowledged via `Event.Acknowledge`. Once this happens, `callback` is invoked ONCE
 *
 * @export
 * @param {string} eventName
 * @param {number} timeInMilliseconds
 * @param {function} callback
 */
export function Heartbeat(eventName, timeInMilliseconds, callback) {

	// Declare interval variable
	let interval = null;

	// Setup callback
	function dynamicCallback() {
		// Kill interval
		clearInterval(interval);
		// Callback
		callback();
	}

	// Register callback
	heartbeatCallbacks[eventName] = dynamicCallback;

	// Start emitting the event
	interval = setInterval(function () {
		Emit(eventName);
	}, timeInMilliseconds);
}

/**
 * Acknowledges a heartbeat event by name
 *
 * @export
 * @param {string} eventName
 */
export function Acknowledge(eventName) {
	// If we are waiting for acknowledgement for this event type
	if (heartbeatCallbacks[eventName]) {
		// Acknowledge!
		heartbeatCallbacks[eventName]();
	} else {
		throw new Error(`Cannot acknowledge unknown heartbeat '${eventName}'`);
	}
}