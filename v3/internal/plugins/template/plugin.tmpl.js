// plugin.js
// This file should contain helper functions for the that can be used by the frontend.
// Below are examples of how to use JSDoc to define the Hashes struct and the exported functions.

/**
 * @typedef {Object} Hashes - A collection of hashes.
 * @property {string} md5 - The MD5 hash of a string, represented as a hexadecimal string.
 * @property {string} sha1 - The SHA-1 hash of a string, represented as a hexadecimal string.
 * @property {string} sha256 - The SHA-256 hash of a string, represented as a hexadecimal string.
 */

/**
 * Generate all hashes for a string.
 * @param input {string} - The string to generate hashes for.
 * @returns {Promise<Hashes>}
 */
export function All(input) {
    return wails.Plugin("{{.Name}}", "All", input);
}

/**
 * Generate the MD5 hash for a string.
 * @param input {string} - The string to generate the hash for.
 * @returns {Promise<string>}
 */
export function MD5(input) {
    return wails.Plugin("{{.Name}}", "MD5", input);
}

/**
 * Generate the SHA-1 hash for a string.
 * @param input {string} - The string to generate the hash for.
 * @returns {Promise<string>}
 */
export function SHA1(input) {
    return wails.Plugin("{{.Name}}", "SHA1", input);
}

/**
 * Generate the SHA-256 hash for a string.
 * @param input {string} - The string to generate the hash for.
 * @returns {Promise<string>}
 */
export function SHA256(input) {
    return wails.Plugin("{{.Name}}", "SHA256", input);
}