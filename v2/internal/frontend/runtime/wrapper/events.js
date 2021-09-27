/*
 _       __      _ __    
| |     / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  ) 
|__/|__/\__,_/_/_/____/  
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

/* jshint esversion: 9 */


/**
 * Registers an event listener that will be invoked `maxCallbacks` times before being destroyed
 *
 * @export
 * @param {string} eventName
 * @param {function} callback
 * @param {number} maxCallbacks
 */
export function EventsOnMultiple(eventName, callback, maxCallbacks) {
	window.runtime.EventsOnMultiple(eventName, callback, maxCallbacks);
}

/**
 * Registers an event listener that will be invoked every time the event is emitted
 *
 * @export
 * @param {string} eventName
 * @param {function} callback
 */
export function EventsOn(eventName, callback) {
	OnMultiple(eventName, callback, -1);
}

/**
 * Registers an event listener that will be invoked once then destroyed
 *
 * @export
 * @param {string} eventName
 * @param {function} callback
 */
export function EventsOnce(eventName, callback) {
	OnMultiple(eventName, callback, 1);
}


/**
 * Emit an event with the given name and data
 *
 * @export
 * @param {string} eventName
 */
export function EventsEmit(eventName) {
	let args = [eventName].slice.call(arguments);
	return window.runtime.EventsEmit.apply(null, args);
}
