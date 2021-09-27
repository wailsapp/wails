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
import * as Log from './log';
import * as Browser from './browser';
import {Acknowledge, Emit, Heartbeat, Notify, On, OnMultiple} from './events';
import {NewBinding} from './bindings';
import {Callback} from './calls';
import {AddScript, InjectCSS, InjectFirebug} from './utils';
import {AddIPCListener} from './ipc';
import * as Store from './store';

// Initialise global if not already
window.wails = window.wails || {};
window.backend = {};

// On webkit2gtk >= 2.32, the external object is not passed to the window context.
// However, IE will throw a strict mode error if window.external is assigned to
// so we need to make sure that line of code isn't reached in IE

// Using !window.external transpiles to `window.external = window.external || ...`
// so we have to use an explicit if statement to prevent webpack from optimizing the code.
if (window.external == undefined) {
	window.external = {
		invoke: function (x) {
			window.webkit.messageHandlers.external.postMessage(x);
		}
	};
}

// Setup internal calls
var internal = {
	NewBinding,
	Callback,
	Notify,
	AddScript,
	InjectCSS,
	Init,
	AddIPCListener,
};

// Setup runtime structure
var runtime = {
	Log,
	Browser,
	Events: {
		On,
		OnMultiple,
		Emit,
		Heartbeat,
		Acknowledge,
	},
	Store,
	_: internal,
};

// Augment global
Object.assign(window.wails, runtime);

// Setup global error handler
window.onerror = function (msg, url, lineNo, columnNo, error) {
	window.wails.Log.Error('**** Caught Unhandled Error ****');
	window.wails.Log.Error('Message: ' + msg);
	window.wails.Log.Error('URL: ' + url);
	window.wails.Log.Error('Line No: ' + lineNo);
	window.wails.Log.Error('Column No: ' + columnNo);
	window.wails.Log.Error('error: ' + error);
};

// Use firebug?
if (window.usefirebug) {
	InjectFirebug();
}

// Emit loaded event
Emit('wails:loaded');

// Nothing to init in production
export function Init(callback) {
	callback();
}
