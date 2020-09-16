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
 * Log the given debug message with the backend
 *
 * @export
 * @param {string} message
 */
function Debug(message) {
	window.wails.Log.Debug(message);
}

/**
 * Log the given info message with the backend
 *
 * @export
 * @param {string} message
 */
function Info(message) {
	window.wails.Log.Info(message);
}

/**
 * Log the given warning message with the backend
 *
 * @export
 * @param {string} message
 */
function Warning(message) {
	window.wails.Log.Warning(message);
}

/**
 * Log the given error message with the backend
 *
 * @export
 * @param {string} message
 */
function Error(message) {
	window.wails.Log.Error(message);
}

/**
 * Log the given fatal message with the backend
 *
 * @export
 * @param {string} message
 */
function Fatal(message) {
	window.wails.Log.Fatal(message);
}

module.exports = {
	Debug: Debug,
	Info: Info,
	Warning: Warning,
	Error: Error,
	Fatal: Fatal
};
