/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

/* jshint esversion: 9 */

import {newRuntimeCallerWithID, objectNames} from "./runtime";
let call = newRuntimeCallerWithID(objectNames.System, '');
const systemIsDarkMode = 0;
const environment = 1;

/**
 * @function
 * Retrieves the system dark mode status.
 * @returns {Promise<boolean>} - A promise that resolves to a boolean value indicating if the system is in dark mode.
 */
export function IsDarkMode() {
    return call(systemIsDarkMode);
}


/**
 * Fetches the capabilities of the application from the server.
 *
 * @async
 * @function Capabilities
 * @returns {Promise<Object>} A promise that resolves to an object containing the capabilities.
 */
export async function Capabilities() {
    let response = fetch("/wails/capabilities");
    return response.json();
}

/**
 * @typedef {object} EnvironmentInfo
 * @property {string} OS - The operating system in use.
 * @property {string} Arch - The architecture of the system.
 */

/**
 * @function
 * Retrieves environment details.
 * @returns {Promise<EnvironmentInfo>} - A promise that resolves to an object containing OS and system architecture.
 */
export function Environment() {
    return call(environment);
}

export let invoke = null;
let environmentCache = null;

Environment()
    .then(result => {
        environmentCache = result;
        invoke = IsWindows() ? window.chrome.webview.postMessage : window.webkit.messageHandlers.external.postMessage;
    })
    .catch(error => {
        console.error(`Error getting Environment: ${error}`);
    });

/**
 * Checks if the current operating system is Windows.
 *
 * @return {boolean} True if the operating system is Windows, otherwise false.
 */
export function IsWindows() {
    return environmentCache.OS === "windows";
}

/**
 * Checks if the current operating system is Linux.
 *
 * @returns {boolean} Returns true if the current operating system is Linux, false otherwise.
 */
export function IsLinux() {
    return environmentCache.OS === "linux";
}

/**
 * Checks if the current environment is a macOS operating system.
 *
 * @returns {boolean} True if the environment is macOS, false otherwise.
 */
export function IsMac() {
    return environmentCache.OS === "darwin";
}

/**
 * Checks if the current environment architecture is AMD64.
 * @returns {boolean} True if the current environment architecture is AMD64, false otherwise.
 */
export function IsAMD64() {
    return environmentCache.Arch === "amd64";
}

/**
 * Checks if the current architecture is ARM.
 *
 * @returns {boolean} True if the current architecture is ARM, false otherwise.
 */
export function IsARM() {
    return environmentCache.Arch === "arm";
}

/**
 * Checks if the current environment is ARM64 architecture.
 *
 * @returns {boolean} - Returns true if the environment is ARM64 architecture, otherwise returns false.
 */
export function IsARM64() {
    return environmentCache.Arch === "arm64";
}

export function IsDebug() {
    return environmentCache.Debug === true;
}