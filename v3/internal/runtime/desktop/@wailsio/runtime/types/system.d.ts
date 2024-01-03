/**
 * @function
 * Retrieves the system dark mode status.
 * @returns {Promise<boolean>} - A promise that resolves to a boolean value indicating if the system is in dark mode.
 */
export function IsDarkMode(): Promise<boolean>;
/**
 * Fetches the capabilities of the application from the server.
 *
 * @async
 * @function Capabilities
 * @returns {Promise<Object>} A promise that resolves to an object containing the capabilities.
 */
export function Capabilities(): Promise<any>;
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
export function Environment(): Promise<EnvironmentInfo>;
/**
 * Checks if the current operating system is Windows.
 *
 * @return {boolean} True if the operating system is Windows, otherwise false.
 */
export function IsWindows(): boolean;
/**
 * Checks if the current operating system is Linux.
 *
 * @returns {boolean} Returns true if the current operating system is Linux, false otherwise.
 */
export function IsLinux(): boolean;
/**
 * Checks if the current environment is a macOS operating system.
 *
 * @returns {boolean} True if the environment is macOS, false otherwise.
 */
export function IsMac(): boolean;
/**
 * Checks if the current environment architecture is AMD64.
 * @returns {boolean} True if the current environment architecture is AMD64, false otherwise.
 */
export function IsAMD64(): boolean;
/**
 * Checks if the current architecture is ARM.
 *
 * @returns {boolean} True if the current architecture is ARM, false otherwise.
 */
export function IsARM(): boolean;
/**
 * Checks if the current environment is ARM64 architecture.
 *
 * @returns {boolean} - Returns true if the environment is ARM64 architecture, otherwise returns false.
 */
export function IsARM64(): boolean;
export function IsDebug(): boolean;
export let invoke: any;
export type EnvironmentInfo = {
    /**
     * - The operating system in use.
     */
    OS: string;
    /**
     * - The architecture of the system.
     */
    Arch: string;
};
