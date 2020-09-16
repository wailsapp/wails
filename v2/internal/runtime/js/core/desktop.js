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

// Setup global error handler
window.onerror = function (/*msg, url, lineNo, columnNo, error*/) {
	// window.wails.Log.Error('**** Caught Unhandled Error ****');
	// window.wails.Log.Error('Message: ' + msg);
	// window.wails.Log.Error('URL: ' + url);
	// window.wails.Log.Error('Line No: ' + lineNo);
	// window.wails.Log.Error('Column No: ' + columnNo);
	// window.wails.Log.Error('error: ' + error);
	(function () { window.wails.Log.Error(new Error().stack); })();
};

// Initialise the Runtime
Init();

// Load Bindings if they exist
if (window.wailsbindings) {
	SetBindings(window.wailsbindings);
}

// Emit loaded event. Leaving this for now. It will show any errors if runtime fails to load.
window.wails.Events.Emit('wails:loaded');

