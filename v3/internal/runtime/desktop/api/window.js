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

/**
 * @typedef {import("./wails").Size} Size
 * @typedef {import("./wails").Position} Position
 * @typedef {import("./wails").Screen} Screen
 */


/**
 * The Window API provides methods to interact with the window.
 */
export const Window = {
    /**
     * Center the window.
     * @returns {Promise<void>}
     */
    Center: () => void wails.Window.Center(),
    /**
     * Set the window title.
     * @param {string} title
     * @returns {Promise<void>}
     */
    SetTitle: (title) => void wails.Window.SetTitle(title),

    /**
     * Makes the window fullscreen.
     * @returns {Promise<void>}
     */
    Fullscreen: () => void wails.Window.Fullscreen(),

    /**
     * Unfullscreen the window.
     * @returns {Promise<void>}
     */
    UnFullscreen: () => void wails.Window.UnFullscreen(),

    /**
     * Set the window size.
     * @param {number} width The window width
     * @param {number} height The window height
     */
    SetSize: (width, height) => void wails.Window.SetSize(width, height),

    /**
     * Get the window size.
     * @returns {Promise<Size>} The window size
     */
    Size: () => {
        return wails.Window.Size();
    },

    /**
     * Set the window maximum size.
     * @param {number} width
     * @param {number} height
     * @returns {Promise<void>}
     */
    SetMaxSize: (width, height) => void wails.Window.SetMaxSize(width, height),

    /**
     * Set the window minimum size.
     * @param {number} width
     * @param {number} height
     * @returns {Promise<void>}
     */
    SetMinSize: (width, height) => void wails.Window.SetMinSize(width, height),

    /**
     * Set window to be always on top.
     * @param {boolean} onTop Whether the window should be always on top
     * @returns {Promise<void>}
     */
    SetAlwaysOnTop: (onTop) => void wails.Window.SetAlwaysOnTop(onTop),

    /**
     * Set the window position relative to the current monitor.
     * @param {number} x
     * @param {number} y
     * @returns {Promise<void>}
     */
    SetRelativePosition: (x, y) => void wails.Window.SetRelativePosition(x, y),

    /**
     * Get the window position relative to the current monitor.
     * @returns {Promise<Position>} The window position
     */
    RelativePosition: () => {
        return wails.Window.RelativePosition();
    },

    /**
     * Set the absolute window position.
     * @param {number} x
     * @param {number} y
     * @returns {Promise<void>}
     */
    SetAbsolutePosition: (x, y) => void wails.Window.SetAbsolutePosition(x, y),

    /**
     * Get the absolute window position.
     * @returns {Promise<Position>} The window position
     */
    AbsolutePosition: () => {
        return wails.Window.AbsolutePosition();
    },

    /**
     * Get the screen the window is on.
     * @returns {Promise<Screen>}
     */
    Screen: () => {
        return wails.Window.Screen();
    },

    /**
     * Hide the window
     * @returns {Promise<void>}
     */
    Hide: () => void wails.Window.Hide(),

    /**
     * Maximise the window
     * @returns {Promise<void>}
     */
    Maximise: () => void wails.Window.Maximise(),

    /**
     * Show the window
     * @returns {Promise<void>}
     */
    Show: () => void wails.Window.Show(),

    /**
     * Close the window
     * @returns {Promise<void>}
     */
    Close: () => void wails.Window.Close(),

    /**
     * Toggle the window maximise state
     * @returns {Promise<void>}
     */
    ToggleMaximise: () => void wails.Window.ToggleMaximise(),

    /**
     * Unmaximise the window
     * @returns {Promise<void>}
     */
    UnMaximise: () => void wails.Window.UnMaximise(),

    /**
     * Minimise the window
     * @returns {Promise<void>}
     */
    Minimise: () => void wails.Window.Minimise(),

    /**
     * Unminimise the window
     * @returns {Promise<void>}
     */
    UnMinimise: () => void wails.Window.UnMinimise(),

    /**
     * Set the background colour of the window.
     * @param {number} r - The red value between 0 and 255
     * @param {number} g - The green value between 0 and 255
     * @param {number} b - The blue value between 0 and 255
     * @param {number} a - The alpha value between 0 and 255
     * @returns {Promise<void>}
     */
    SetBackgroundColour: (r, g, b, a) => void wails.Window.SetBackgroundColour(r, g, b, a),
};
