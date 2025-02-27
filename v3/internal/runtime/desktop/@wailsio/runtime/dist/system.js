/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/
import { newRuntimeCaller, objectNames } from "./runtime.js";
const call = newRuntimeCaller(objectNames.System);
const SystemIsDarkMode = 0;
const SystemEnvironment = 1;
const _invoke = (function () {
    var _a, _b, _c, _d, _e;
    try {
        if ((_b = (_a = window.chrome) === null || _a === void 0 ? void 0 : _a.webview) === null || _b === void 0 ? void 0 : _b.postMessage) {
            return window.chrome.webview.postMessage.bind(window.chrome.webview);
        }
        else if ((_e = (_d = (_c = window.webkit) === null || _c === void 0 ? void 0 : _c.messageHandlers) === null || _d === void 0 ? void 0 : _d['external']) === null || _e === void 0 ? void 0 : _e.postMessage) {
            return window.webkit.messageHandlers['external'].postMessage.bind(window.webkit.messageHandlers['external']);
        }
    }
    catch (e) { }
    console.warn('\n%c⚠️ Browser Environment Detected %c\n\n%cOnly UI previews are available in the browser. For full functionality, please run the application in desktop mode.\nMore information at: https://v3.wails.io/learn/build/#using-a-browser-for-development\n', 'background: #ffffff; color: #000000; font-weight: bold; padding: 4px 8px; border-radius: 4px; border: 2px solid #000000;', 'background: transparent;', 'color: #ffffff; font-style: italic; font-weight: bold;');
    return null;
})();
export function invoke(msg) {
    _invoke === null || _invoke === void 0 ? void 0 : _invoke(msg);
}
/**
 * Retrieves the system dark mode status.
 *
 * @returns A promise that resolves to a boolean value indicating if the system is in dark mode.
 */
export function IsDarkMode() {
    return call(SystemIsDarkMode);
}
/**
 * Fetches the capabilities of the application from the server.
 *
 * @returns A promise that resolves to an object containing the capabilities.
 */
export async function Capabilities() {
    let response = await fetch("/wails/capabilities");
    if (response.ok) {
        return response.json();
    }
    else {
        throw new Error("could not fetch capabilities: " + response.statusText);
    }
}
/**
 * Retrieves environment details.
 *
 * @returns A promise that resolves to an object containing OS and system architecture.
 */
export function Environment() {
    return call(SystemEnvironment);
}
/**
 * Checks if the current operating system is Windows.
 *
 * @return True if the operating system is Windows, otherwise false.
 */
export function IsWindows() {
    return window._wails.environment.OS === "windows";
}
/**
 * Checks if the current operating system is Linux.
 *
 * @returns Returns true if the current operating system is Linux, false otherwise.
 */
export function IsLinux() {
    return window._wails.environment.OS === "linux";
}
/**
 * Checks if the current environment is a macOS operating system.
 *
 * @returns True if the environment is macOS, false otherwise.
 */
export function IsMac() {
    return window._wails.environment.OS === "darwin";
}
/**
 * Checks if the current environment architecture is AMD64.
 *
 * @returns True if the current environment architecture is AMD64, false otherwise.
 */
export function IsAMD64() {
    return window._wails.environment.Arch === "amd64";
}
/**
 * Checks if the current architecture is ARM.
 *
 * @returns True if the current architecture is ARM, false otherwise.
 */
export function IsARM() {
    return window._wails.environment.Arch === "arm";
}
/**
 * Checks if the current environment is ARM64 architecture.
 *
 * @returns Returns true if the environment is ARM64 architecture, otherwise returns false.
 */
export function IsARM64() {
    return window._wails.environment.Arch === "arm64";
}
/**
 * Reports whether the app is being run in debug mode.
 *
 * @returns True if the app is being run in debug mode.
 */
export function IsDebug() {
    return Boolean(window._wails.environment.Debug);
}
