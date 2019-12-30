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
 * Opens the given URL in the system browser
 *
 * @export
 * @param {string} url
 * @returns
 */
function OpenURL(url) {
	return window.wails.Browser.OpenURL(url);
}

/**
 * Opens the given filename using the system's default file handler
 *
 * @export
 * @param {sting} filename
 * @returns
 */
function OpenFile(filename) {
	return window.wails.Browser.OpenFile(filename);
}

module.exports = {
	OpenURL: OpenURL,
	OpenFile: OpenFile
};