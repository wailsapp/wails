/*
 _       __      _ __    
| |     / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  ) 
|__/|__/\__,_/_/_/____/  
The electron alternative for Go
(c) Lea Anthony 2019-present
*/
/* jshint esversion: 6 */

import {Call} from './calls';

// This is where we bind go method wrappers
window.go = {};

export function SetBindings(bindingsMap) {
	try {
		bindingsMap = JSON.parse(bindingsMap);
	} catch (e) {
		console.error(e);
	}

	// Initialise the bindings map
	window.go = window.go || {};

	// Iterate package names
	Object.keys(bindingsMap).forEach((packageName) => {

		// Create inner map if it doesn't exist
		window.go[packageName] = window.go[packageName] || {};

		// Iterate struct names
		Object.keys(bindingsMap[packageName]).forEach((structName) => {

			// Create inner map if it doesn't exist
			window.go[packageName][structName] = window.go[packageName][structName] || {};

			Object.keys(bindingsMap[packageName][structName]).forEach((methodName) => {

				window.go[packageName][structName][methodName] = function () {

					// No timeout by default
					let timeout = 0;

					// Actual function
					function dynamic() {
						const args = [].slice.call(arguments);
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
