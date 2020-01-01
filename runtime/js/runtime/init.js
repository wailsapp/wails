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
 * Initialises the Wails runtime
 *
 * @param {function} callback
 */
function Init(callback) {
	window.wails._.Init(callback);
}

module.exports = Init;
