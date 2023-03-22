// plugin.js
// This file should contain helper functions for the that can be used by the frontend.
// Below are examples of how to use JSDoc to define the Hashes struct and the exported functions.

/**
 * Log at the Trace level.
 * @param input {string} - The message in printf format.
 * @param args {...any} - The arguments for the log message.
 * @returns {Promise<void|Error>}
 */
export function Trace(input, ...args) {
    return wails.Plugin("log", "Trace", input, ...args);
}

/**
 * Log at the Debug level.
 * @param input {string} - The message in printf format.
 * @param args {...any} - The arguments for the log message.
 * @returns {Promise<void|Error>}
 */

export function Debug(input, ...args) {
    return wails.Plugin("log", "Debug", input, ...args);
}

/**
 * Log at the Info level.
 * @param input {string} - The message in printf format.
 * @param args {...any} - The arguments for the log message.
 * @returns {Promise<void|Error>}
 */
export function Info(input, ...args) {
    return wails.Plugin("log", "Info", input, ...args);
}

/**
 * Log at the Warning level.
 * @param input {string} - The message in printf format.
 * @param args {...any} - The arguments for the log message.
 * @returns {Promise<void|Error>}
 */
export function Warning(input, ...args) {
    return wails.Plugin("log", "Warning", input, ...args);
}

/**
 * Log at the Error level.
 * @param input {string} - The message in printf format.
 * @param args {...any} - The arguments for the log message.
 * @returns {Promise<void|Error>}
 */
export function Error(input, ...args) {
    return wails.Plugin("log", "Error", input, ...args);
}

/**
 * Log at the Fatal level.
 * @param input {string} - The message in printf format.
 * @param args {...any} - The arguments for the log message.
 * @returns {Promise<void|Error>}
 */
export function Fatal(input, ...args) {
    return wails.Plugin("log", "Fatal", input, ...args);
}