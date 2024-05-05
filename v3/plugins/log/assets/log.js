// plugin.js
// This file should contain helper functions for the that can be used by the frontend.
// Below are examples of how to use JSDoc to define the Hashes struct and the exported functions.

import {Call} from '/wails/runtime.js';

/**
 * Log at the Debug level.
 * @param input {string} - The message in printf format.
 * @param args {...any} - The arguments for the log message.
 * @returns {Promise<void|Error>}
 */

export function Debug(input, ...args) {
    return Call.ByID(4111675027, input, ...args);
}

/**
 * Log at the Info level.
 * @param input {string} - The message in printf format.
 * @param args {...any} - The arguments for the log message.
 * @returns {Promise<void|Error>}
 */
export function Info(input, ...args) {
    return Call.ByID(2391172776, input, ...args);
}

/**
 * Log at the Warning level.
 * @param input {string} - The message in printf format.
 * @param args {...any} - The arguments for the log message.
 * @returns {Promise<void|Error>}
 */
export function Warning(input, ...args) {
    return Call.ByID(2762394760, input, ...args);
}

/**
 * Log at the Error level.
 * @param input {string} - The message in printf format.
 * @param args {...any} - The arguments for the log message.
 * @returns {Promise<void|Error>}
 */
export function Error(input, ...args) {
    return Call.ByID(878590242, input, ...args);
}

export const LevelDebug  = -4
export const LevelInfo   = 0
export const LevelWarn   = 4
export const LevelError  = 8


/**
 * Set Log level
 * @param level {LogLevel} - The log level to set.
 * @returns {Promise<void|Error>}
 */
export function SetLogLevel(level) {
    return Call.ByID(2758810652, level);
}
