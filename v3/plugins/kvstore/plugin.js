// plugin.js
// This file should contain helper functions for the that can be used by the frontend.
// Below are examples of how to use JSDoc to define the Hashes struct and the exported functions.

/**
 * Get the value of a key.
 * @param key {string} - The store key.
 * @returns {Promise<any>} - The value of the key.
 */
export function Get(key) {
    return wails.Plugin("kvstore", "Get", key);
}

/**
 * Set the value of a key.
 @param key {string} - The store key.
 @param value {any} - The value to set.
 * @returns {Promise<void>}
 */
export function Set(key, value) {
    return wails.Plugin("kvstore", "Set", key, value);
}


/**
 * Save the database to disk.
 * @returns {Promise<void|Error>}
 */
export function Save() {
    return wails.Plugin("kvstore", "Save");
}
