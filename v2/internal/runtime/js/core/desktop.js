/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The lightweight framework for web-like apps
(c) Lea Anthony 2019-present
*/
/* jshint esversion: 6 */
import { SetBindings } from './bindings';
import { Init } from './main';
import {RaiseError} from '../desktop/darwin';

// Setup global error handler
window.onerror = function (msg, url, lineNo, columnNo, error) {
	const errorMessage = {
		message: msg,
		url: url,
		line: lineNo,
		column: columnNo,
		error: JSON.stringify(error),
		stack: function() { return JSON.stringify(new Error().stack); }(),
	};
	RaiseError(errorMessage);
};

// Initialise the Runtime
Init();

// Load Bindings if they exist
if (window.wailsbindings) {
	SetBindings(window.wailsbindings);
}

// Emit loaded event. Leaving this for now. It will show any errors if runtime fails to load.
window.wails.Events.Emit('wails:loaded');

