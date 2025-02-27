import type { Screen } from "./screens.js";
/**
 * A record describing the position of a window.
 */
interface Position {
    /** The horizontal position of the window. */
    x: number;
    /** The vertical position of the window. */
    y: number;
}
/**
 * A record describing the size of a window.
 */
interface Size {
    /** The width of the window. */
    width: number;
    /** The height of the window. */
    height: number;
}
declare const callerSym: unique symbol;
declare class Window {
    private [callerSym];
    /**
     * Initialises a window object with the specified name.
     *
     * @private
     * @param name - The name of the target window.
     */
    constructor(name?: string);
    /**
     * Gets the specified window.
     *
     * @param name - The name of the window to get.
     * @returns The corresponding window object.
     */
    Get(name: string): Window;
    /**
     * Returns the absolute position of the window.
     *
     * @returns The current absolute position of the window.
     */
    Position(): Promise<Position>;
    /**
     * Centers the window on the screen.
     */
    Center(): Promise<void>;
    /**
     * Closes the window.
     */
    Close(): Promise<void>;
    /**
     * Disables min/max size constraints.
     */
    DisableSizeConstraints(): Promise<void>;
    /**
     * Enables min/max size constraints.
     */
    EnableSizeConstraints(): Promise<void>;
    /**
     * Focuses the window.
     */
    Focus(): Promise<void>;
    /**
     * Forces the window to reload the page assets.
     */
    ForceReload(): Promise<void>;
    /**
     * Switches the window to fullscreen mode.
     */
    Fullscreen(): Promise<void>;
    /**
     * Returns the screen that the window is on.
     *
     * @returns The screen the window is currently on.
     */
    GetScreen(): Promise<Screen>;
    /**
     * Returns the current zoom level of the window.
     *
     * @returns The current zoom level.
     */
    GetZoom(): Promise<number>;
    /**
     * Returns the height of the window.
     *
     * @returns The current height of the window.
     */
    Height(): Promise<number>;
    /**
     * Hides the window.
     */
    Hide(): Promise<void>;
    /**
     * Returns true if the window is focused.
     *
     * @returns Whether the window is currently focused.
     */
    IsFocused(): Promise<boolean>;
    /**
     * Returns true if the window is fullscreen.
     *
     * @returns Whether the window is currently fullscreen.
     */
    IsFullscreen(): Promise<boolean>;
    /**
     * Returns true if the window is maximised.
     *
     * @returns Whether the window is currently maximised.
     */
    IsMaximised(): Promise<boolean>;
    /**
     * Returns true if the window is minimised.
     *
     * @returns Whether the window is currently minimised.
     */
    IsMinimised(): Promise<boolean>;
    /**
     * Maximises the window.
     */
    Maximise(): Promise<void>;
    /**
     * Minimises the window.
     */
    Minimise(): Promise<void>;
    /**
     * Returns the name of the window.
     *
     * @returns The name of the window.
     */
    Name(): Promise<string>;
    /**
     * Opens the development tools pane.
     */
    OpenDevTools(): Promise<void>;
    /**
     * Returns the relative position of the window to the screen.
     *
     * @returns The current relative position of the window.
     */
    RelativePosition(): Promise<Position>;
    /**
     * Reloads the page assets.
     */
    Reload(): Promise<void>;
    /**
     * Returns true if the window is resizable.
     *
     * @returns Whether the window is currently resizable.
     */
    Resizable(): Promise<boolean>;
    /**
     * Restores the window to its previous state if it was previously minimised, maximised or fullscreen.
     */
    Restore(): Promise<void>;
    /**
     * Sets the absolute position of the window.
     *
     * @param x - The desired horizontal absolute position of the window.
     * @param y - The desired vertical absolute position of the window.
     */
    SetPosition(x: number, y: number): Promise<void>;
    /**
     * Sets the window to be always on top.
     *
     * @param alwaysOnTop - Whether the window should stay on top.
     */
    SetAlwaysOnTop(alwaysOnTop: boolean): Promise<void>;
    /**
     * Sets the background colour of the window.
     *
     * @param r - The desired red component of the window background.
     * @param g - The desired green component of the window background.
     * @param b - The desired blue component of the window background.
     * @param a - The desired alpha component of the window background.
     */
    SetBackgroundColour(r: number, g: number, b: number, a: number): Promise<void>;
    /**
     * Removes the window frame and title bar.
     *
     * @param frameless - Whether the window should be frameless.
     */
    SetFrameless(frameless: boolean): Promise<void>;
    /**
     * Disables the system fullscreen button.
     *
     * @param enabled - Whether the fullscreen button should be enabled.
     */
    SetFullscreenButtonEnabled(enabled: boolean): Promise<void>;
    /**
     * Sets the maximum size of the window.
     *
     * @param width - The desired maximum width of the window.
     * @param height - The desired maximum height of the window.
     */
    SetMaxSize(width: number, height: number): Promise<void>;
    /**
     * Sets the minimum size of the window.
     *
     * @param width - The desired minimum width of the window.
     * @param height - The desired minimum height of the window.
     */
    SetMinSize(width: number, height: number): Promise<void>;
    /**
     * Sets the relative position of the window to the screen.
     *
     * @param x - The desired horizontal relative position of the window.
     * @param y - The desired vertical relative position of the window.
     */
    SetRelativePosition(x: number, y: number): Promise<void>;
    /**
     * Sets whether the window is resizable.
     *
     * @param resizable - Whether the window should be resizable.
     */
    SetResizable(resizable: boolean): Promise<void>;
    /**
     * Sets the size of the window.
     *
     * @param width - The desired width of the window.
     * @param height - The desired height of the window.
     */
    SetSize(width: number, height: number): Promise<void>;
    /**
     * Sets the title of the window.
     *
     * @param title - The desired title of the window.
     */
    SetTitle(title: string): Promise<void>;
    /**
     * Sets the zoom level of the window.
     *
     * @param zoom - The desired zoom level.
     */
    SetZoom(zoom: number): Promise<void>;
    /**
     * Shows the window.
     */
    Show(): Promise<void>;
    /**
     * Returns the size of the window.
     *
     * @returns The current size of the window.
     */
    Size(): Promise<Size>;
    /**
     * Toggles the window between fullscreen and normal.
     */
    ToggleFullscreen(): Promise<void>;
    /**
     * Toggles the window between maximised and normal.
     */
    ToggleMaximise(): Promise<void>;
    /**
     * Un-fullscreens the window.
     */
    UnFullscreen(): Promise<void>;
    /**
     * Un-maximises the window.
     */
    UnMaximise(): Promise<void>;
    /**
     * Un-minimises the window.
     */
    UnMinimise(): Promise<void>;
    /**
     * Returns the width of the window.
     *
     * @returns The current width of the window.
     */
    Width(): Promise<number>;
    /**
     * Zooms the window.
     */
    Zoom(): Promise<void>;
    /**
     * Increases the zoom level of the webview content.
     */
    ZoomIn(): Promise<void>;
    /**
     * Decreases the zoom level of the webview content.
     */
    ZoomOut(): Promise<void>;
    /**
     * Resets the zoom level of the webview content.
     */
    ZoomReset(): Promise<void>;
}
/**
 * The window within which the script is running.
 */
declare const thisWindow: Window;
export default thisWindow;
