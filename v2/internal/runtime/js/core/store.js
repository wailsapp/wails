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
 * Creates a new sync store with the given name and optional default value
 *
 * @export
 * @param {string} name
 * @param {*} optionalDefault
 */
export function New(name, optionalDefault) {

	// This is the store state
	var state;

	// Check we are initialised
	if( !window.wails) {
		throw Error('Wails is not initialised');
	}

	// Store for the callbacks
	let callbacks = [];
	
	// Subscribe to updates by providing a callback
	let subscribe = function(callback) {
		callbacks.push(callback);
	};

	// sets the store state to the provided `newstate` value
	let set = function(newstate) {
		
		state = newstate;

		// Emit a notification to back end
		window.wails.Events.Emit('wails:sync:store:updatedbyfrontend:'+name, JSON.stringify(state));

		// Notify callbacks
		callbacks.forEach( function(callback) {
			callback(state);
		});
	};

	// update mutates the value in the store by calling the
	// provided method with the current value. The value returned 
	// by the updater function will be set as the new store value
	let update = function(updater) {
		var newValue = updater(state);
		this.set(newValue);
	};

	// get returns the current value of the store
	let get = function() {
		return state;
	};

	// Setup event callback
	window.wails.Events.On('wails:sync:store:updatedbybackend:'+name, function(jsonEncodedState) {

		// Parse state
		let newState = JSON.parse(jsonEncodedState);

		// Todo: Potential preprocessing?

		// Save state
		state = newState;

		// Notify callbacks
		callbacks.forEach( function(callback) {
			callback(state);
		});

	});

	// Set to the optional default if set
	if( optionalDefault ) {
		this.set(optionalDefault);
	}

	// Trigger an update to the store
	window.wails.Events.Emit('wails:sync:store:resync:'+name);

	return {
		subscribe,
		get,
		set,
		update,
	};
}