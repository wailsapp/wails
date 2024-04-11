// plugin.js
// This file should contain helper functions for the that can be used by the frontend.
// Below are examples of how to use JSDoc to define the Hashes struct and the exported functions.

import { Call } from '/wails/runtime.js';

/**
 * Get the value of a key.
 * @param key {string} - The store key.
 * @returns {Promise<any>} - The value of the key.
 */
export function Get(key) {
    return Call.ByID(3322496224, key);
}

/**
 * Set the value of a key.
 @param key {string} - The store key.
 @param value {any} - The value to set.
 * @returns {Promise<void>}
 */
export function Set(key, value) {
    return Call.ByID(1207638860, key, value);
}


/**
 * Save the database to disk.
 * @returns {Promise<void|Error>}
 */
export function Save() {
    return Call.ByID(1377075201);
}

/**
 * Delete a key from the store.
 * @param key {string} - The key to delete.
 * @returns {Promise<void>}
 */
export function Delete(key) {
    return Call.ByID(737249231, key);
}