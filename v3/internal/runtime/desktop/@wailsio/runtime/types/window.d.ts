/**
 * Gets the specified window.
 *
 * @param {string} windowName - The name of the window to get.
 * @return {Object} - The specified window object.
 */
export function Get(windowName: string): any;
declare const _default: any;
export default _default;
export type Screen = import("./screens").Screen;
export type Position = {
    /**
     * - The X coordinate.
     */
    X: number;
    /**
     * - The Y coordinate.
     */
    Y: number;
};
export type Size = {
    /**
     * - The width.
     */
    X: number;
    /**
     * - The height.
     */
    Y: number;
};
export type Rect = {
    /**
     * - The X coordinate of the top-left corner.
     */
    X: number;
    /**
     * - The Y coordinate of the top-left corner.
     */
    Y: number;
    /**
     * - The width of the rectangle.
     */
    Width: number;
    /**
     * - The height of the rectangle.
     */
    Height: number;
};
/**
 * The rotation of the screen. Can be one of 'Zero', 'Ninety', 'OneEighty', 'TwoSeventy'.
 */
export type Rotation = ('Zero' | 'Ninety' | 'OneEighty' | 'TwoSeventy');
