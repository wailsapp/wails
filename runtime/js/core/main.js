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
import { On, OnMultiple, Emit, Notify, Heartbeat, Acknowledge } from './events';
import { NewBinding } from './bindings';
import { Callback } from './calls';
import { AddScript, InjectCSS } from './utils';

// Initialise global if not already
window.wails = window.wails || {};
window.backend = {};

// Setup internal calls
var internal = {
	NewBinding,
	Callback,
	Notify,
	AddScript,
	InjectCSS,
	Init,
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

// Emit loaded event
Emit('wails:loaded');

// Nothing to init in production
export function Init(callback) {
	callback();
}