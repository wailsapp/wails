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

const _invoke = (() => {
    try {
        if(window?.chrome?.webview) {
            return (msg) => window.chrome.webview.postMessage(msg);
        }
        if(window?.webkit?.messageHandlers?.external) {
            return (msg) => window.webkit.messageHandlers.external.postMessage(msg);
        }
    } catch(e) {
        console.warn('\n%c⚠️ Browser Environment Detected %c\n\n%cOnly UI previews are available in the browser. For full functionality, please run the application in desktop mode.\nMore information at: https://v3.wails.io/learn/build/#using-a-browser-for-development\n',
            'background: #ffffff; color: #000000; font-weight: bold; padding: 4px 8px; border-radius: 4px; border: 2px solid #000000;',
            'background: transparent;',
            'color: #ffffff; font-style: italic; font-weight: bold;');
    }
    return null;
})();

export function invoke(msg) {
    if (!_invoke) return;
    return _invoke(msg);
}

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
export function Capabilities() {
    let response = fetch("/wails/capabilities");
    return response.json();
}

/**
 * @typedef {Object} OSInfo
 * @property {string} Branding - The branding of the OS.
 * @property {string} ID - The ID of the OS.
 * @property {string} Name - The name of the OS.
 * @property {string} Version - The version of the OS.
 */

/**
 * @typedef {Object} EnvironmentInfo
 * @property {string} Arch - The architecture of the system.
 * @property {boolean} Debug - True if the application is running in debug mode, otherwise false.
 * @property {string} OS - The operating system in use.
 * @property {OSInfo} OSInfo - Details of the operating system.
 * @property {Object} PlatformInfo - Additional platform information.
 */

/**
 * @function
 * Retrieves environment details.
 * @returns {Promise<EnvironmentInfo>} - A promise that resolves to an object containing OS and system architecture.
 */
export function Environment() {
    return call(environment);
}

/**
 * Checks if the current operating system is Windows.
 *
 * @return {boolean} True if the operating system is Windows, otherwise false.
 */
export function IsWindows() {
    return window._wails.environment.OS === "windows";
}

/**
 * Checks if the current operating system is Linux.
 *
 * @returns {boolean} Returns true if the current operating system is Linux, false otherwise.
 */
export function IsLinux() {
    return window._wails.environment.OS === "linux";
}

/**
 * Checks if the current environment is a macOS operating system.
 *
 * @returns {boolean} True if the environment is macOS, false otherwise.
 */
export function IsMac() {
    return window._wails.environment.OS === "darwin";
}

/**
 * Checks if the current environment architecture is AMD64.
 * @returns {boolean} True if the current environment architecture is AMD64, false otherwise.
 */
export function IsAMD64() {
    return window._wails.environment.Arch === "amd64";
}

/**
 * Checks if the current architecture is ARM.
 *
 * @returns {boolean} True if the current architecture is ARM, false otherwise.
 */
export function IsARM() {
    return window._wails.environment.Arch === "arm";
}

/**
 * Checks if the current environment is ARM64 architecture.
 *
 * @returns {boolean} - Returns true if the environment is ARM64 architecture, otherwise returns false.
 */
export function IsARM64() {
    return window._wails.environment.Arch === "arm64";
}

export function IsDebug() {
    return window._wails.environment.Debug === true;
}

