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
 * Registers an event listener that will be invoked `maxCallbacks` times before being destroyed
 *
 * @export
 * @param {string} eventName
 * @param {function} callback
 * @param {number} maxCallbacks
 */
function OnMultiple(eventName, callback, maxCallbacks) {
	window.wails.Events.OnMultiple(eventName, callback, maxCallbacks);
}

/**
 * Registers an event listener that will be invoked every time the event is emitted
 *
 * @export
 * @param {string} eventName
 * @param {function} callback
 */
function On(eventName, callback) {
	OnMultiple(eventName, callback);
}

/**
 * Registers an event listener that will be invoked once then destroyed
 *
 * @export
 * @param {string} eventName
 * @param {function} callback
 */
function Once(eventName, callback) {
	OnMultiple(eventName, callback, 1);
}


/**
 * Emit an event with the given name and data
 *
 * @export
 * @param {string} eventName
 */
function Emit(eventName) {
	var args = [eventName].slice.call(arguments);
	return window.wails.Events.Emit.apply(null, args);
}

module.exports = {
	OnMultiple: OnMultiple,
	On: On,
	Once: Once,
	Emit: Emit,
};