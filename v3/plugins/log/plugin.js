// plugin.js
// This file should contain helper functions for the that can be used by the frontend.
// Below are examples of how to use JSDoc to define the Hashes struct and the exported functions.

/**
 * Log at the Trace level.
 * @param input {string} - The message in printf format.
 * @param args {...any} - The arguments for the log message.
 * @returns {Promise<void|Error>}
 */
function Trace(input, ...args) {
    return wails.Plugin("log", "Trace", input, ...args);
}

/**
 * Log at the Debug level.
 * @param input {string} - The message in printf format.
 * @param args {...any} - The arguments for the log message.
 * @returns {Promise<void|Error>}
 */

function Debug(input, ...args) {
    return wails.Plugin("log", "Debug", input, ...args);
}

/**
 * Log at the Info level.
 * @param input {string} - The message in printf format.
 * @param args {...any} - The arguments for the log message.
 * @returns {Promise<void|Error>}
 */
function Info(input, ...args) {
    return wails.Plugin("log", "Info", input, ...args);
}

/**
 * Log at the Warning level.
 * @param input {string} - The message in printf format.
 * @param args {...any} - The arguments for the log message.
 * @returns {Promise<void|Error>}
 */
function Warning(input, ...args) {
    return wails.Plugin("log", "Warning", input, ...args);
}

/**
 * Log at the Error level.
 * @param input {string} - The message in printf format.
 * @param args {...any} - The arguments for the log message.
 * @returns {Promise<void|Error>}
 */
function Error(input, ...args) {
    return wails.Plugin("log", "Error", input, ...args);
}

/**
 * Log at the Fatal level.
 * @param input {string} - The message in printf format.
 * @param args {...any} - The arguments for the log message.
 * @returns {Promise<void|Error>}
 */
function Fatal(input, ...args) {
    return wails.Plugin("log", "Fatal", input, ...args);
}

/**
 * SetLevel sets the logging level
 * @param level {Level} The log level to set
 * @returns {Promise<void>}
 */
function SetLevel(level) {
    return wails.Plugin("log", "SetLevel", level);
}

/**
 * Log Level.
 * @readonly
 * @enum {number}
 */
let Level = {
    Trace: 1,
    Debug: 2,
    Info: 3,
    Warning: 4,
    Error: 5,
    Fatal: 6,
};

window.Logger = {
    Trace,
    Debug,
    Info,
    Warning,
    Error,
    Fatal,
    SetLevel,
    Level,
}
