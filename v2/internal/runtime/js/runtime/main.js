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

const Log = require('./log');
const Browser = require('./browser');
const Dialog = require('./dialog');
const Events = require('./events');
const Init = require('./init');
const System = require('./system');
const Store = require('./store');
const Window = require('./window');
const Tray = require('./tray');

module.exports = {
	Browser: Browser,
	Dialog: Dialog,
	Events: Events,
	ready: Init.ready,
	Log: Log,
	System: System,
	Store: Store,
	Window: Window,
	Tray: Tray,
};