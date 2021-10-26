/*
 _       __      _ __    
| |     / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  ) 
|__/|__/\__,_/_/_/____/  
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

/* jshint esversion: 6 */


/**
 * Log the given message with the backend
 *
 * @export
 * @param {string} message
 */
function Print(message) {
	window.wails.Log.Print(message);
}

/**
 * Log the given trace message with the backend
 *
 * @export
 * @param {string} message
 */
function Trace(message) {
	window.wails.Log.Trace(message);
}

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
 * @param {string} message
 */
function Fatal(message) {
	window.wails.Log.Fatal(message);
}


/**
 * Sets the Log level to the given log level
 *
 * @param {number} loglevel
 */
function SetLogLevel(loglevel) {
	window.wails.Log.SetLogLevel(loglevel);
}

// Log levels
const Level = {
	TRACE: 1,
	DEBUG: 2,
	INFO: 3,
	WARNING: 4,
	ERROR: 5,
};


module.exports = {
	Print: Print,
	Trace: Trace,
	Debug: Debug,
	Info: Info,
	Warning: Warning,
	Error: Error,
	Fatal: Fatal,
	SetLogLevel: SetLogLevel,
	Level: Level,
};
