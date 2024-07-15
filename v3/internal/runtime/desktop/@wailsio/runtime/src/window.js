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

// Import screen jsdoc definition from ./screens.js
/**
 * @typedef {import("./screens").Screen} Screen
 */


/**
 * A record describing the position of a window.
 *
 * @typedef {Object} Position
 * @property {number} x - The horizontal position of the window
 * @property {number} y - The vertical position of the window
 */


/**
 * A record describing the size of a window.
 *
 * @typedef {Object} Size
 * @property {number} width - The width of the window
 * @property {number} height - The height of the window
 */


import {newRuntimeCallerWithID, objectNames} from "./runtime";

const PositionMethod                    = 0;
const CenterMethod                      = 1;
const CloseMethod                       = 2;
const DisableSizeConstraintsMethod      = 3;
const EnableSizeConstraintsMethod       = 4;
const FocusMethod                       = 5;
const ForceReloadMethod                 = 6;
const FullscreenMethod                  = 7;
const GetScreenMethod                   = 8;
const GetZoomMethod                     = 9;
const HeightMethod                      = 10;
const HideMethod                        = 11;
const IsFocusedMethod                   = 12;
const IsFullscreenMethod                = 13;
const IsMaximisedMethod                 = 14;
const IsMinimisedMethod                 = 15;
const MaximiseMethod                    = 16;
const MinimiseMethod                    = 17;
const NameMethod                        = 18;
const OpenDevToolsMethod                = 19;
const RelativePositionMethod            = 20;
const ReloadMethod                      = 21;
const ResizableMethod                   = 22;
const RestoreMethod                     = 23;
const SetPositionMethod                 = 24;
const SetAlwaysOnTopMethod              = 25;
const SetBackgroundColourMethod         = 26;
const SetFramelessMethod                = 27;
const SetFullscreenButtonEnabledMethod  = 28;
const SetMaxSizeMethod                  = 29;
const SetMinSizeMethod                  = 30;
const SetRelativePositionMethod         = 31;
const SetResizableMethod                = 32;
const SetSizeMethod                     = 33;
const SetTitleMethod                    = 34;
const SetZoomMethod                     = 35;
const ShowMethod                        = 36;
const SizeMethod                        = 37;
const ToggleFullscreenMethod            = 38;
const ToggleMaximiseMethod              = 39;
const UnFullscreenMethod                = 40;
const UnMaximiseMethod                  = 41;
const UnMinimiseMethod                  = 42;
const WidthMethod                       = 43;
const ZoomMethod                        = 44;
const ZoomInMethod                      = 45;
const ZoomOutMethod                     = 46;
const ZoomResetMethod                   = 47;

/**
 * @type {symbol}
 */
const caller = Symbol();

class Window {
    /**
     * Initialises a window object with the specified name.
     *
     * @private
     * @param {string} name - The name of the target window.
     */
    constructor(name = '') {
        /**
         * @private
         * @name {@link caller}
         * @type {(...args: any[]) => any}
         */
        this[caller] = newRuntimeCallerWithID(objectNames.Window, name)

        // bind instance method to make them easily usable in event handlers
        for (const method of Object.getOwnPropertyNames(Window.prototype)) {
            if (
                method !== "constructor"
                && typeof this[method] === "function"
            ) {
                this[method] = this[method].bind(this);
            }
        }
    }

    /**
     * Gets the specified window.
     *
     * @public
     * @param {string} name - The name of the window to get.
     * @return {Window} - The corresponding window object.
     */
    Get(name) {
        return new Window(name);
    }

    /**
     * Returns the absolute position of the window.
     *
     * @public
     * @return {Promise<Position>} - The current absolute position of the window.
     */
    Position() {
        return this[caller](PositionMethod);
    }

    /**
     * Centers the window on the screen.
     *
     * @public
     * @return {Promise<void>}
     */
    Center() {
        return this[caller](CenterMethod);
    }

    /**
     * Closes the window.
     *
     * @public
     * @return {Promise<void>}
     */
    Close() {
        return this[caller](CloseMethod);
    }

    /**
     * Disables min/max size constraints.
     *
     * @public
     * @return {Promise<void>}
     */
    DisableSizeConstraints() {
        return this[caller](DisableSizeConstraintsMethod);
    }

    /**
     * Enables min/max size constraints.
     *
     * @public
     * @return {Promise<void>}
     */
    EnableSizeConstraints() {
        return this[caller](EnableSizeConstraintsMethod);
    }

    /**
     * Focuses the window.
     *
     * @public
     * @return {Promise<void>}
     */
    Focus() {
        return this[caller](FocusMethod);
    }

    /**
     * Forces the window to reload the page assets.
     *
     * @public
     * @return {Promise<void>}
     */
    ForceReload() {
        return this[caller](ForceReloadMethod);
    }

    /**
     * Doc.
     *
     * @public
     * @return {Promise<void>}
     */
    Fullscreen() {
        return this[caller](FullscreenMethod);
    }

    /**
     * Returns the screen that the window is on.
     *
     * @public
     * @return {Promise<Screen>} - The screen the window is currently on
     */
    GetScreen() {
        return this[caller](GetScreenMethod);
    }

    /**
     * Returns the current zoom level of the window.
     *
     * @public
     * @return {Promise<number>} - The current zoom level
     */
    GetZoom() {
        return this[caller](GetZoomMethod);
    }

    /**
     * Returns the height of the window.
     *
     * @public
     * @return {Promise<number>} - The current height of the window
     */
    Height() {
        return this[caller](HeightMethod);
    }

    /**
     * Hides the window.
     *
     * @public
     * @return {Promise<void>}
     */
    Hide() {
        return this[caller](HideMethod);
    }

    /**
     * Returns true if the window is focused.
     *
     * @public
     * @return {Promise<boolean>} - Whether the window is currently focused
     */
    IsFocused() {
        return this[caller](IsFocusedMethod);
    }

    /**
     * Returns true if the window is fullscreen.
     *
     * @public
     * @return {Promise<boolean>} - Whether the window is currently fullscreen
     */
    IsFullscreen() {
        return this[caller](IsFullscreenMethod);
    }

    /**
     * Returns true if the window is maximised.
     *
     * @public
     * @return {Promise<boolean>} - Whether the window is currently maximised
     */
    IsMaximised() {
        return this[caller](IsMaximisedMethod);
    }

    /**
     * Returns true if the window is minimised.
     *
     * @public
     * @return {Promise<boolean>} - Whether the window is currently minimised
     */
    IsMinimised() {
        return this[caller](IsMinimisedMethod);
    }

    /**
     * Maximises the window.
     *
     * @public
     * @return {Promise<void>}
     */
    Maximise() {
        return this[caller](MaximiseMethod);
    }

    /**
     * Minimises the window.
     *
     * @public
     * @return {Promise<void>}
     */
    Minimise() {
        return this[caller](MinimiseMethod);
    }

    /**
     * Returns the name of the window.
     *
     * @public
     * @return {Promise<string>} - The name of the window
     */
    Name() {
        return this[caller](NameMethod);
    }

    /**
     * Opens the development tools pane.
     *
     * @public
     * @return {Promise<void>}
     */
    OpenDevTools() {
        return this[caller](OpenDevToolsMethod);
    }

    /**
     * Returns the relative position of the window to the screen.
     *
     * @public
     * @return {Promise<Position>} - The current relative position of the window
     */
    RelativePosition() {
        return this[caller](RelativePositionMethod);
    }

    /**
     * Reloads the page assets.
     *
     * @public
     * @return {Promise<void>}
     */
    Reload() {
        return this[caller](ReloadMethod);
    }

    /**
     * Returns true if the window is resizable.
     *
     * @public
     * @return {Promise<boolean>} - Whether the window is currently resizable
     */
    Resizable() {
        return this[caller](ResizableMethod);
    }

    /**
     * Restores the window to its previous state if it was previously minimised, maximised or fullscreen.
     *
     * @public
     * @return {Promise<void>}
     */
    Restore() {
        return this[caller](RestoreMethod);
    }

    /**
     * Sets the absolute position of the window.
     *
     * @public
     * @param {number} x - The desired horizontal absolute position of the window
     * @param {number} y - The desired vertical absolute position of the window
     * @return {Promise<void>}
     */
    SetPosition(x, y) {
        return this[caller](SetPositionMethod, { x, y });
    }

    /**
     * Sets the window to be always on top.
     *
     * @public
     * @param {boolean} alwaysOnTop - Whether the window should stay on top
     * @return {Promise<void>}
     */
    SetAlwaysOnTop(alwaysOnTop) {
        return this[caller](SetAlwaysOnTopMethod, { alwaysOnTop });
    }

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
    SetBackgroundColour(r, g, b, a) {
        return this[caller](SetBackgroundColourMethod, { r, g, b, a });
    }

    /**
     * Removes the window frame and title bar.
     *
     * @public
     * @param {boolean} frameless - Whether the window should be frameless
     * @return {Promise<void>}
     */
    SetFrameless(frameless) {
        return this[caller](SetFramelessMethod, { frameless });
    }

    /**
     * Disables the system fullscreen button.
     *
     * @public
     * @param {boolean} enabled - Whether the fullscreen button should be enabled
     * @return {Promise<void>}
     */
    SetFullscreenButtonEnabled(enabled) {
        return this[caller](SetFullscreenButtonEnabledMethod, { enabled });
    }

    /**
     * Sets the maximum size of the window.
     *
     * @public
     * @param {number} width - The desired maximum width of the window
     * @param {number} height - The desired maximum height of the window
     * @return {Promise<void>}
     */
    SetMaxSize(width, height) {
        return this[caller](SetMaxSizeMethod, { width, height });
    }

    /**
     * Sets the minimum size of the window.
     *
     * @public
     * @param {number} width - The desired minimum width of the window
     * @param {number} height - The desired minimum height of the window
     * @return {Promise<void>}
     */
    SetMinSize(width, height) {
        return this[caller](SetMinSizeMethod, { width, height });
    }

    /**
     * Sets the relative position of the window to the screen.
     *
     * @public
     * @param {number} x - The desired horizontal relative position of the window
     * @param {number} y - The desired vertical relative position of the window
     * @return {Promise<void>}
     */
    SetRelativePosition(x, y) {
        return this[caller](SetRelativePositionMethod, { x, y });
    }

    /**
     * Sets whether the window is resizable.
     *
     * @public
     * @param {boolean} resizable - Whether the window should be resizable
     * @return {Promise<void>}
     */
    SetResizable(resizable) {
        return this[caller](SetResizableMethod, { resizable });
    }

    /**
     * Sets the size of the window.
     *
     * @public
     * @param {number} width - The desired width of the window
     * @param {number} height - The desired height of the window
     * @return {Promise<void>}
     */
    SetSize(width, height) {
        return this[caller](SetSizeMethod, { width, height });
    }

    /**
     * Sets the title of the window.
     *
     * @public
     * @param {string} title - The desired title of the window
     * @return {Promise<void>}
     */
    SetTitle(title) {
        return this[caller](SetTitleMethod, { title });
    }

    /**
     * Sets the zoom level of the window.
     *
     * @public
     * @param {number} zoom - The desired zoom level
     * @return {Promise<void>}
     */
    SetZoom(zoom) {
        return this[caller](SetZoomMethod, { zoom });
    }

    /**
     * Shows the window.
     *
     * @public
     * @return {Promise<void>}
     */
    Show() {
        return this[caller](ShowMethod);
    }

    /**
     * Returns the size of the window.
     *
     * @public
     * @return {Promise<Size>} - The current size of the window
     */
    Size() {
        return this[caller](SizeMethod);
    }

    /**
     * Toggles the window between fullscreen and normal.
     *
     * @public
     * @return {Promise<void>}
     */
    ToggleFullscreen() {
        return this[caller](ToggleFullscreenMethod);
    }

    /**
     * Toggles the window between maximised and normal.
     *
     * @public
     * @return {Promise<void>}
     */
    ToggleMaximise() {
        return this[caller](ToggleMaximiseMethod);
    }

    /**
     * Un-fullscreens the window.
     *
     * @public
     * @return {Promise<void>}
     */
    UnFullscreen() {
        return this[caller](UnFullscreenMethod);
    }

    /**
     * Un-maximises the window.
     *
     * @public
     * @return {Promise<void>}
     */
    UnMaximise() {
        return this[caller](UnMaximiseMethod);
    }

    /**
     * Un-minimises the window.
     *
     * @public
     * @return {Promise<void>}
     */
    UnMinimise() {
        return this[caller](UnMinimiseMethod);
    }

    /**
     * Returns the width of the window.
     *
     * @public
     * @return {Promise<number>} - The current width of the window
     */
    Width() {
        return this[caller](WidthMethod);
    }

    /**
     * Zooms the window.
     *
     * @public
     * @return {Promise<void>}
     */
    Zoom() {
        return this[caller](ZoomMethod);
    }

    /**
     * Increases the zoom level of the webview content.
     *
     * @public
     * @return {Promise<void>}
     */
    ZoomIn() {
        return this[caller](ZoomInMethod);
    }

    /**
     * Decreases the zoom level of the webview content.
     *
     * @public
     * @return {Promise<void>}
     */
    ZoomOut() {
        return this[caller](ZoomOutMethod);
    }

    /**
     * Resets the zoom level of the webview content.
     *
     * @public
     * @return {Promise<void>}
     */
    ZoomReset() {
        return this[caller](ZoomResetMethod);
    }
}

/**
 * The window within which the script is running.
 *
 * @type {Window}
 */
const thisWindow = new Window('');

export default thisWindow;
