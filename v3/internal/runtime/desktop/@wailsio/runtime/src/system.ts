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
const SystemCapabilities = 2;
const ApplicationFilesDroppedWithContext = 100; // New method ID for enriched drop event

const _invoke = (function () {
    try {
        if ((window as any).chrome?.webview?.postMessage) {
            return (window as any).chrome.webview.postMessage.bind((window as any).chrome.webview);
        } else if ((window as any).webkit?.messageHandlers?.['external']?.postMessage) {
            return (window as any).webkit.messageHandlers['external'].postMessage.bind((window as any).webkit.messageHandlers['external']);
        }
    } catch(e) {}

    console.warn('\n%c⚠️ Browser Environment Detected %c\n\n%cOnly UI previews are available in the browser. For full functionality, please run the application in desktop mode.\nMore information at: https://v3.wails.io/learn/build/#using-a-browser-for-development\n',
        'background: #ffffff; color: #000000; font-weight: bold; padding: 4px 8px; border-radius: 4px; border: 2px solid #000000;',
        'background: transparent;',
        'color: #ffffff; font-style: italic; font-weight: bold;');
    return null;
})();

export function invoke(msg: any): void {
    _invoke?.(msg);
}

/**
 * Retrieves the system dark mode status.
 *
 * @returns A promise that resolves to a boolean value indicating if the system is in dark mode.
 */
export function IsDarkMode(): Promise<boolean> {
    return call(SystemIsDarkMode);
}

/**
 * Fetches the capabilities of the application from the server.
 *
 * @returns A promise that resolves to an object containing the capabilities.
 */
export async function Capabilities(): Promise<Record<string, any>> {
    return call(SystemCapabilities);
}

export interface OSInfo {
    /** The branding of the OS. */
    Branding: string;
    /** The ID of the OS. */
    ID: string;
    /** The name of the OS. */
    Name: string;
    /** The version of the OS. */
    Version: string;
}

export interface EnvironmentInfo {
    /** The architecture of the system. */
    Arch: string;
    /** True if the application is running in debug mode, otherwise false. */
    Debug: boolean;
    /** The operating system in use. */
    OS: string;
    /** Details of the operating system. */
    OSInfo: OSInfo;
    /** Additional platform information. */
    PlatformInfo: Record<string, any>;
}

/**
 * Retrieves environment details.
 *
 * @returns A promise that resolves to an object containing OS and system architecture.
 */
export function Environment(): Promise<EnvironmentInfo> {
    return call(SystemEnvironment);
}

/**
 * Checks if the current operating system is Windows.
 *
 * @return True if the operating system is Windows, otherwise false.
 */
export function IsWindows(): boolean {
    return window._wails.environment.OS === "windows";
}

/**
 * Checks if the current operating system is Linux.
 *
 * @returns Returns true if the current operating system is Linux, false otherwise.
 */
export function IsLinux(): boolean {
    return window._wails.environment.OS === "linux";
}

/**
 * Checks if the current environment is a macOS operating system.
 *
 * @returns True if the environment is macOS, false otherwise.
 */
export function IsMac(): boolean {
    return window._wails.environment.OS === "darwin";
}

/**
 * Checks if the current environment architecture is AMD64.
 *
 * @returns True if the current environment architecture is AMD64, false otherwise.
 */
export function IsAMD64(): boolean {
    return window._wails.environment.Arch === "amd64";
}

/**
 * Checks if the current architecture is ARM.
 *
 * @returns True if the current architecture is ARM, false otherwise.
 */
export function IsARM(): boolean {
    return window._wails.environment.Arch === "arm";
}

/**
 * Checks if the current environment is ARM64 architecture.
 *
 * @returns Returns true if the environment is ARM64 architecture, otherwise returns false.
 */
export function IsARM64(): boolean {
    return window._wails.environment.Arch === "arm64";
}

/**
 * Reports whether the app is being run in debug mode.
 *
 * @returns True if the app is being run in debug mode.
 */
export function IsDebug(): boolean {
    return Boolean(window._wails.environment.Debug);
}

/**
 * Handles file drops originating from platform-specific code (e.g., macOS native drag-and-drop).
 * Gathers information about the drop target element and sends it back to the Go backend.
 *
 * @param filenames - An array of file paths (strings) that were dropped.
 * @param x - The x-coordinate of the drop event.
 * @param y - The y-coordinate of the drop event.
 */
export function HandlePlatformFileDrop(filenames: string[], x: number, y: number): void {
    const element = document.elementFromPoint(x, y);
    const elementId = element ? element.id : '';
    const classList = element ? Array.from(element.classList) : [];

    const payload = {
        filenames,
        x,
        y,
        elementId,
        classList,
    };

    call(ApplicationFilesDroppedWithContext, payload)
        .then(() => {
            // Optional: Log success or handle if needed
            console.log("Platform file drop processed and sent to Go.");
        })
        .catch(err => {
            // Optional: Log error
            console.error("Error sending platform file drop to Go:", err);
        });
}

