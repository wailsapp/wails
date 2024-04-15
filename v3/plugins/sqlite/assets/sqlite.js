
import {Call} from '/wails/runtime.js';

/**
 * Open a sqlite DB.
 * @param filename {string} - file to open.
 * @returns {Promise<void>}
 */
export function Open(filename) {
    return Call.ByID(147348976, filename);
}

/**
 * Close a sqlite DB.
 * @returns {Promise<void>}
 */
export function Close() {
    return Call.ByID(3998329564);
}

/**
 * Execute a SQL statement.
 * @param statement {string} - SQL statement to execute.
 * @param args {...any} - Arguments to pass to the statement.
 * @returns {Promise<void>}
 */
export function Execute(statement, ...args) {
    return Call.ByID(2804887383, statement, ...args);
}

/**
 * Perform a select query.
 * @param statement {string} - Select SQL statement.
 * @param args {...any} - Arguments to pass to the statement.
 * @returns {Promise<any>}
 */
export function Select(statement, ...args) {
    return Call.ByID(2209315040, statement, ...args);
}
