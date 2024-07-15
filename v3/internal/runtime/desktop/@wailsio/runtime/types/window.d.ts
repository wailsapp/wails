export default thisWindow;
export type Screen = import("./screens").Screen;
/**
 * A record describing the position of a window.
 */
export type Position = {
    /**
     * - The horizontal position of the window
     */
    x: number;
    /**
     * - The vertical position of the window
     */
    y: number;
};
/**
 * A record describing the size of a window.
 */
export type Size = {
    /**
     * - The width of the window
     */
    width: number;
    /**
     * - The height of the window
     */
    height: number;
};
/**
 * The window within which the script is running.
 *
 * @type {Window}
 */
declare const thisWindow: Window;
declare class Window {
    /**
     * Initialises a window object with the specified name.
     *
     * @private
     * @param {string} name - The name of the target window.
     */
    private constructor();
    /**
     * Gets the specified window.
     *
     * @public
     * @param {string} name - The name of the window to get.
     * @return {Window} - The corresponding window object.
     */
    public Get(name: string): Window;
    /**
     * Returns the absolute position of the window.
     *
     * @public
     * @return {Promise<Position>} - The current absolute position of the window.
     */
    public Position(): Promise<Position>;
    /**
     * Centers the window on the screen.
     *
     * @public
     * @return {Promise<void>}
     */
    public Center(): Promise<void>;
    /**
     * Closes the window.
     *
     * @public
     * @return {Promise<void>}
     */
    public Close(): Promise<void>;
    /**
     * Disables min/max size constraints.
     *
     * @public
     * @return {Promise<void>}
     */
    public DisableSizeConstraints(): Promise<void>;
    /**
     * Enables min/max size constraints.
     *
     * @public
     * @return {Promise<void>}
     */
    public EnableSizeConstraints(): Promise<void>;
    /**
     * Focuses the window.
     *
     * @public
     * @return {Promise<void>}
     */
    public Focus(): Promise<void>;
    /**
     * Forces the window to reload the page assets.
     *
     * @public
     * @return {Promise<void>}
     */
    public ForceReload(): Promise<void>;
    /**
     * Doc.
     *
     * @public
     * @return {Promise<void>}
     */
    public Fullscreen(): Promise<void>;
    /**
     * Returns the screen that the window is on.
     *
     * @public
     * @return {Promise<Screen>} - The screen the window is currently on
     */
    public GetScreen(): Promise<Screen>;
    /**
     * Returns the current zoom level of the window.
     *
     * @public
     * @return {Promise<number>} - The current zoom level
     */
    public GetZoom(): Promise<number>;
    /**
     * Returns the height of the window.
     *
     * @public
     * @return {Promise<number>} - The current height of the window
     */
    public Height(): Promise<number>;
    /**
     * Hides the window.
     *
     * @public
     * @return {Promise<void>}
     */
    public Hide(): Promise<void>;
    /**
     * Returns true if the window is focused.
     *
     * @public
     * @return {Promise<boolean>} - Whether the window is currently focused
     */
    public IsFocused(): Promise<boolean>;
    /**
     * Returns true if the window is fullscreen.
     *
     * @public
     * @return {Promise<boolean>} - Whether the window is currently fullscreen
     */
    public IsFullscreen(): Promise<boolean>;
    /**
     * Returns true if the window is maximised.
     *
     * @public
     * @return {Promise<boolean>} - Whether the window is currently maximised
     */
    public IsMaximised(): Promise<boolean>;
    /**
     * Returns true if the window is minimised.
     *
     * @public
     * @return {Promise<boolean>} - Whether the window is currently minimised
     */
    public IsMinimised(): Promise<boolean>;
    /**
     * Maximises the window.
     *
     * @public
     * @return {Promise<void>}
     */
    public Maximise(): Promise<void>;
    /**
     * Minimises the window.
     *
     * @public
     * @return {Promise<void>}
     */
    public Minimise(): Promise<void>;
    /**
     * Returns the name of the window.
     *
     * @public
     * @return {Promise<string>} - The name of the window
     */
    public Name(): Promise<string>;
    /**
     * Opens the development tools pane.
     *
     * @public
     * @return {Promise<void>}
     */
    public OpenDevTools(): Promise<void>;
    /**
     * Returns the relative position of the window to the screen.
     *
     * @public
     * @return {Promise<Position>} - The current relative position of the window
     */
    public RelativePosition(): Promise<Position>;
    /**
     * Reloads the page assets.
     *
     * @public
     * @return {Promise<void>}
     */
    public Reload(): Promise<void>;
    /**
     * Returns true if the window is resizable.
     *
     * @public
     * @return {Promise<boolean>} - Whether the window is currently resizable
     */
    public Resizable(): Promise<boolean>;
    /**
     * Restores the window to its previous state if it was previously minimised, maximised or fullscreen.
     *
     * @public
     * @return {Promise<void>}
     */
    public Restore(): Promise<void>;
    /**
     * Sets the absolute position of the window.
     *
     * @public
     * @param {number} x - The desired horizontal absolute position of the window
     * @param {number} y - The desired vertical absolute position of the window
     * @return {Promise<void>}
     */
    public SetPosition(x: number, y: number): Promise<void>;
    /**
     * Sets the window to be always on top.
     *
     * @public
     * @param {boolean} alwaysOnTop - Whether the window should stay on top
     * @return {Promise<void>}
     */
    public SetAlwaysOnTop(alwaysOnTop: boolean): Promise<void>;
    /**
     * Sets the background colour of the window.
     *
     * @public
     * @param {number} r - The desired red component of the window background
     * @param {number} g - The desired green component of the window background
     * @param {number} b - The desired blue component of the window background
     * @param {number} a - The desired alpha component of the window background
     * @return {Promise<void>}
     */
    public SetBackgroundColour(r: number, g: number, b: number, a: number): Promise<void>;
    /**
     * Removes the window frame and title bar.
     *
     * @public
     * @param {boolean} frameless - Whether the window should be frameless
     * @return {Promise<void>}
     */
    public SetFrameless(frameless: boolean): Promise<void>;
    /**
     * Disables the system fullscreen button.
     *
     * @public
     * @param {boolean} enabled - Whether the fullscreen button should be enabled
     * @return {Promise<void>}
     */
    public SetFullscreenButtonEnabled(enabled: boolean): Promise<void>;
    /**
     * Sets the maximum size of the window.
     *
     * @public
     * @param {number} width - The desired maximum width of the window
     * @param {number} height - The desired maximum height of the window
     * @return {Promise<void>}
     */
    public SetMaxSize(width: number, height: number): Promise<void>;
    /**
     * Sets the minimum size of the window.
     *
     * @public
     * @param {number} width - The desired minimum width of the window
     * @param {number} height - The desired minimum height of the window
     * @return {Promise<void>}
     */
    public SetMinSize(width: number, height: number): Promise<void>;
    /**
     * Sets the relative position of the window to the screen.
     *
     * @public
     * @param {number} x - The desired horizontal relative position of the window
     * @param {number} y - The desired vertical relative position of the window
     * @return {Promise<void>}
     */
    public SetRelativePosition(x: number, y: number): Promise<void>;
    /**
     * Sets whether the window is resizable.
     *
     * @public
     * @param {boolean} resizable - Whether the window should be resizable
     * @return {Promise<void>}
     */
    public SetResizable(resizable: boolean): Promise<void>;
    /**
     * Sets the size of the window.
     *
     * @public
     * @param {number} width - The desired width of the window
     * @param {number} height - The desired height of the window
     * @return {Promise<void>}
     */
    public SetSize(width: number, height: number): Promise<void>;
    /**
     * Sets the title of the window.
     *
     * @public
     * @param {string} title - The desired title of the window
     * @return {Promise<void>}
     */
    public SetTitle(title: string): Promise<void>;
    /**
     * Sets the zoom level of the window.
     *
     * @public
     * @param {number} zoom - The desired zoom level
     * @return {Promise<void>}
     */
    public SetZoom(zoom: number): Promise<void>;
    /**
     * Shows the window.
     *
     * @public
     * @return {Promise<void>}
     */
    public Show(): Promise<void>;
    /**
     * Returns the size of the window.
     *
     * @public
     * @return {Promise<Size>} - The current size of the window
     */
    public Size(): Promise<Size>;
    /**
     * Toggles the window between fullscreen and normal.
     *
     * @public
     * @return {Promise<void>}
     */
    public ToggleFullscreen(): Promise<void>;
    /**
     * Toggles the window between maximised and normal.
     *
     * @public
     * @return {Promise<void>}
     */
    public ToggleMaximise(): Promise<void>;
    /**
     * Un-fullscreens the window.
     *
     * @public
     * @return {Promise<void>}
     */
    public UnFullscreen(): Promise<void>;
    /**
     * Un-maximises the window.
     *
     * @public
     * @return {Promise<void>}
     */
    public UnMaximise(): Promise<void>;
    /**
     * Un-minimises the window.
     *
     * @public
     * @return {Promise<void>}
     */
    public UnMinimise(): Promise<void>;
    /**
     * Returns the width of the window.
     *
     * @public
     * @return {Promise<number>} - The current width of the window
     */
    public Width(): Promise<number>;
    /**
     * Zooms the window.
     *
     * @public
     * @return {Promise<void>}
     */
    public Zoom(): Promise<void>;
    /**
     * Increases the zoom level of the webview content.
     *
     * @public
     * @return {Promise<void>}
     */
    public ZoomIn(): Promise<void>;
    /**
     * Decreases the zoom level of the webview content.
     *
     * @public
     * @return {Promise<void>}
     */
    public ZoomOut(): Promise<void>;
    /**
     * Resets the zoom level of the webview content.
     *
     * @public
     * @return {Promise<void>}
     */
    public ZoomReset(): Promise<void>;
}
