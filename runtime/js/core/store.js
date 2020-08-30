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

	if( !window.wails) {
		throw Error('Wails is not initialised');
	}

	let callbacks = [];
	
	this.subscribe = (callback) => {
		callbacks.push(callback);
	};

	this.set = (newdata) => {
		
		data = newdata;

		// Emit the data
		window.wails.Events.Emit('wails:sync:store:updated:'+name, data);
	};

	this.update = (updater) => {
		var newValue = updater(data);
		this.set(newValue);
	};

	// Setup event listener
	window.wails.Events.On('wails:sync:store:updated:'+name, function(result) {

		// Save data
		data = result;

		// Notify listeners
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