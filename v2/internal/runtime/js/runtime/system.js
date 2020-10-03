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

const Events = require('./events');

/**
 * Registers an event listener that will be invoked when the user changes the
 * desktop theme (light mode / dark mode). The callback receives a boolean which
 * indicates if dark mode is enabled.
 *
 * @export
 * @param {function} callback The callback to invoke on theme change
 */
function OnThemeChange(callback) {
	Events.On("wails:system:themechange", callback);
}

/**
 * Checks if dark mode is curently enabled.
 *
 * @export
 * @returns {Promise}
 */
function DarkModeEnabled() {
	return window.wails._.SystemCall("IsDarkMode");
}

module.exports = {
	OnThemeChange: OnThemeChange,
	DarkModeEnabled: DarkModeEnabled,
};