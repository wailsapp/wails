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

	var data;

	// Check we are initialised
	if( !window.wails) {
		throw Error('Wails is not initialised');
	}

	// Store for the callbacks
	let callbacks = [];
	
	// Subscribe to updates by providing a callback
	this.subscribe = (callback) => {
		callbacks.push(callback);
	};

	// sets the store data to the provided `newdata` value
	this.set = (newdata) => {
		
		data = newdata;

		// Emit a notification to back end
		window.wails.Events.Emit('wails:sync:store:updatedbyfrontend:'+name, JSON.stringify(data));

		// Notify callbacks
		callbacks.forEach( function(callback) {
			callback(data);
		});
	};

	// update mutates the value in the store by calling the
	// provided method with the current value. The value returned 
	// by the updater function will be set as the new store value
	this.update = (updater) => {
		var newValue = updater(data);
		this.set(newValue);
	};

	// Setup event callback
	window.wails.Events.On('wails:sync:store:updatedbybackend:'+name, function(result) {

		// Parse data
		result = JSON.parse(result);

		// Todo: Potential preprocessing?

		// Save data
		data = result;

		// Notify callbacks
		callbacks.forEach( function(callback) {
			callback(data);
		});

	});

	// Set to the optional default if set
	if( optionalDefault ) {
		this.set(optionalDefault);
	}

	return this;
}