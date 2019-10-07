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
 * Determines if the given identifier is valid Javascript
 *
 * @param {boolean} name
 * @returns
 */
function isValidIdentifier(name) {
	// Don't xss yourself :-)
	try {
		new Function('var ' + name);
		return true;
	} catch (e) {
		return false;
	}
}

/**
 * NewBinding creates a new binding from the given binding name
 *
 * @export
 * @param {string} bindingName
 * @returns
 */
// eslint-disable-next-line max-lines-per-function
export function NewBinding(bindingName) {

	// Get all the sections of the binding
	var bindingSections = [].concat(bindingName.split('.').splice(1));
	var pathToBinding = window.backend;

	// Check if we have a path (IE Struct)
	if (bindingSections.length > 1) {
		// Iterate over binding sections, adding them to the window.backend object
		for (let index = 0; index < bindingSections.length-1; index += 1) {
			const name = bindingSections[index];
			// Is name a valid javascript identifier?
			if (!isValidIdentifier(name)) {
				return new Error(`${name} is not a valid javascript identifier.`);
			}
			if (!pathToBinding[name]) {
				pathToBinding[name] = {};
			}
			pathToBinding = pathToBinding[name];
		}
	}

	// Get the actual function/method call name
	const name = bindingSections.pop();

	// Is name a valid javascript identifier?
	if (!isValidIdentifier(name)) {
		return new Error(`${name} is not a valid javascript identifier.`);
	}

	// Add binding call
	pathToBinding[name] = function () {

		// No timeout by default
		var timeout = 0;

		// Actual function
		function dynamic() {
			var args = [].slice.call(arguments);
			return Call(bindingName, args, timeout);
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
}
