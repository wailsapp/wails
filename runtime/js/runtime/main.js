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
const Events = require('./events');
const Init = require('./init');
const Store = require('./store');

module.exports = {
	Log: Log,
	Browser: Browser,
	Events: Events,
	Init: Init,
	Store: Store,
};