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

import { Call } from './calls';

window.backend = {};

/**
 
Map of this format:

{ 
	packageName: { 
		structName: { 
			methodName: { 
				name: "", 
				inputs: [
					{
						type: <type>
					}
				],
				outputs: [
					{
						type: <type>
					}
				] 
			} 
		} 
	} 
}
 */

export function SetBindings(bindingsMap) {
	try {
		bindingsMap = JSON.parse(bindingsMap);
	} catch (e) {
		console.error(e);
	}

	// Initialise the backend map
	window.backend = window.backend || {};

	// Iterate package names
	Object.keys(bindingsMap).forEach((packageName) => {

		// Create inner map if it doesn't exist
		window.backend[packageName] = window.backend[packageName] || {};

		// Iterate struct names
		Object.keys(bindingsMap[packageName]).forEach((structName) => {

			// Create inner map if it doesn't exist
			window.backend[packageName][structName] = window.backend[packageName][structName] || {};

			Object.keys(bindingsMap[packageName][structName]).forEach((methodName) => {

				window.backend[packageName][structName][methodName] = function () {

					// No timeout by default
					let timeout = 0;

					// Actual function
					function dynamic() {
						var args = [].slice.call(arguments);
						return Call([packageName, structName, methodName].join('.'), args, timeout);
					}

					// Allow setting timeout to function
					dynamic.setTimeout = function (newTimeout) {
						timeout = newTimeout;
					};

					// Allow getting timeout to function
					dynamic.getTimeout = function () {
						return timeout;
					};

					return dynamic;
				}();
			});
		});
	});
}
