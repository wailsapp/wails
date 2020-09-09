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
 * Create a new Store with the given name and optional default value
 *
 * @export
 * @param {string} name
 * @param {*} optionalDefault
 */
function New(name, optionalDefault) {
	return window.wails.Store.New(name, optionalDefault);
}

module.exports = {
	New: New,
};
