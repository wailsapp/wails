/**
 * Gets the specified window.
 *
 * @param {string} windowName - The name of the window to get.
 * @return {Object} - The specified window object.
 */
export function Get(windowName: string): any;
/**
 * Centers the window on the screen.
 */
export function Center(): void;
/**
 * Sets the title of the window.
 * @param {string} title - The title to set.
 */
export function SetTitle(title: string): void;
/**
 * Sets the window to fullscreen.
 */
export function Fullscreen(): void;
/**
 * Sets the size of the window.
 * @param {number} width - The width of the window.
 * @param {number} height - The height of the window.
 */
export function SetSize(width: number, height: number): void;
/**
 * Gets the size of the window.
 */
export function Size(): any;
/**
 * Sets the maximum size of the window.
 * @param {number} width - The maximum width of the window.
 * @param {number} height - The maximum height of the window.
 */
export function SetMaxSize(width: number, height: number): void;
/**
 * Sets the minimum size of the window.
 * @param {number} width - The minimum width of the window.
 * @param {number} height - The minimum height of the window.
 */
export function SetMinSize(width: number, height: number): void;
/**
 * Sets the window to always be on top.
 * @param {boolean} onTop - Whether the window should always be on top.
 */
export function SetAlwaysOnTop(onTop: boolean): void;
/**
 * Sets the relative position of the window.
 * @param {number} x - The x-coordinate of the window's position.
 * @param {number} y - The y-coordinate of the window's position.
 */
export function SetRelativePosition(x: number, y: number): void;
/**
 * Gets the relative position of the window.
 */
export function RelativePosition(): any;
/**
 * Gets the screen that the window is on.
 */
export function Screen(): any;
export type Screen = import("./screens").Screen;
/**
 * Hides the window.
 */
export function Hide(): void;
/**
 * Maximises the window.
 */
export function Maximise(): void;
/**
 * Un-maximises the window.
 */
export function UnMaximise(): void;
/**
 * Toggles the maximisation of the window.
 */
export function ToggleMaximise(): void;
/**
 * Minimises the window.
 */
export function Minimise(): void;
/**
 * Un-minimises the window.
 */
export function UnMinimise(): void;
/**
 * Restores the window.
 */
export function Restore(): void;
/**
 * Shows the window.
 */
export function Show(): void;
/**
 * Closes the window.
 */
export function Close(): void;
/**
 * Sets the background colour of the window.
 * @param {number} r - The red component of the colour.
 * @param {number} g - The green component of the colour.
 * @param {number} b - The blue component of the colour.
 * @param {number} a - The alpha component of the colour.
 */
export function SetBackgroundColour(r: number, g: number, b: number, a: number): void;
/**
 * Sets whether the window is resizable.
 * @param {boolean} resizable - Whether the window should be resizable.
 */
export function SetResizable(resizable: boolean): void;
/**
 * Gets the width of the window.
 */
export function Width(): any;
/**
 * Gets the height of the window.
 */
export function Height(): any;
/**
 * Zooms in the window.
 */
export function ZoomIn(): void;
/**
 * Zooms out the window.
 */
export function ZoomOut(): void;
/**
 * Resets the zoom of the window.
 */
export function ZoomReset(): void;
/**
 * Gets the zoom level of the window.
 */
export function GetZoomLevel(): any;
/**
 * Sets the zoom level of the window.
 * @param {number} zoomLevel - The zoom level to set.
 */
export function SetZoomLevel(zoomLevel: number): void;
