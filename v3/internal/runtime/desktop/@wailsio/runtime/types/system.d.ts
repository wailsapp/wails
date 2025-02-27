export declare function invoke(msg: any): void;
/**
 * Retrieves the system dark mode status.
 *
 * @returns A promise that resolves to a boolean value indicating if the system is in dark mode.
 */
export declare function IsDarkMode(): Promise<boolean>;
/**
 * Fetches the capabilities of the application from the server.
 *
 * @returns A promise that resolves to an object containing the capabilities.
 */
export declare function Capabilities(): Promise<Record<string, any>>;
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
export declare function Environment(): Promise<EnvironmentInfo>;
/**
 * Checks if the current operating system is Windows.
 *
 * @return True if the operating system is Windows, otherwise false.
 */
export declare function IsWindows(): boolean;
/**
 * Checks if the current operating system is Linux.
 *
 * @returns Returns true if the current operating system is Linux, false otherwise.
 */
export declare function IsLinux(): boolean;
/**
 * Checks if the current environment is a macOS operating system.
 *
 * @returns True if the environment is macOS, false otherwise.
 */
export declare function IsMac(): boolean;
/**
 * Checks if the current environment architecture is AMD64.
 *
 * @returns True if the current environment architecture is AMD64, false otherwise.
 */
export declare function IsAMD64(): boolean;
/**
 * Checks if the current architecture is ARM.
 *
 * @returns True if the current architecture is ARM, false otherwise.
 */
export declare function IsARM(): boolean;
/**
 * Checks if the current environment is ARM64 architecture.
 *
 * @returns Returns true if the environment is ARM64 architecture, otherwise returns false.
 */
export declare function IsARM64(): boolean;
/**
 * Reports whether the app is being run in debug mode.
 *
 * @returns True if the app is being run in debug mode.
 */
export declare function IsDebug(): boolean;
