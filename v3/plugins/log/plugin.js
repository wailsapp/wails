// plugin.js
// This file should contain helper functions for the that can be used by the frontend.
// Below are examples of how to use JSDoc to define the Hashes struct and the exported functions.

/**
 * Log at the Debug level.
 * @param input {string} - The message in printf format.
 * @param args {...any} - The arguments for the log message.
 * @returns {Promise<void|Error>}
 */

function Debug(input, ...args) {
    return wails.CallByID(4111675027, input, ...args);
}

/**
 * Log at the Info level.
 * @param input {string} - The message in printf format.
 * @param args {...any} - The arguments for the log message.
 * @returns {Promise<void|Error>}
 */
function Info(input, ...args) {
    return wails.CallByID(2391172776, input, ...args);
}

/**
 * Log at the Warning level.
 * @param input {string} - The message in printf format.
 * @param args {...any} - The arguments for the log message.
 * @returns {Promise<void|Error>}
 */
function Warning(input, ...args) {
    return wails.CallByID(2762394760, input, ...args);
}

/**
 * Log at the Error level.
 * @param input {string} - The message in printf format.
 * @param args {...any} - The arguments for the log message.
 * @returns {Promise<void|Error>}
 */
function Error(input, ...args) {
    return wails.CallByID(878590242, input, ...args);
}
