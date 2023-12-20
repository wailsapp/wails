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
    Center: () => wails.Window.Center(),
    /**
     * Set the window title.
     * @param {string} title
     * @returns {Promise<void>}
     */
    SetTitle: (title) => wails.Window.SetTitle(title),

    /**
     * Makes the window fullscreen.
     * @returns {Promise<void>}
     */
    Fullscreen: () => wails.Window.Fullscreen(),

    /**
     * Unfullscreen the window.
     * @returns {Promise<void>}
     */
    UnFullscreen: () => wails.Window.UnFullscreen(),

    /**
     * Set the window size.
     * @param {number} width The window width
     * @param {number} height The window height
     * @returns {Promise<void>}
     */
    SetSize: (width, height) => wails.Window.SetSize(width, height),

    /**
     * Get the window size.
     * @returns {Promise<Size>} The window size
     */
    Size: () => wails.Window.Size(),

    /**
     * Set the window maximum size.
     * @param {number} width
     * @param {number} height
     * @returns {Promise<void>}
     */
    SetMaxSize: (width, height) => wails.Window.SetMaxSize(width, height),

    /**
     * Set the window minimum size.
     * @param {number} width
     * @param {number} height
     * @returns {Promise<void>}
     */
    SetMinSize: (width, height) => wails.Window.SetMinSize(width, height),

    /**
     * Set window to be always on top.
     * @param {boolean} onTop Whether the window should be always on top
     * @returns {Promise<void>}
     */
    SetAlwaysOnTop: (onTop) => wails.Window.SetAlwaysOnTop(onTop),

    /**
     * Set the window position relative to the current monitor.
     * @param {number} x
     * @param {number} y
     * @returns {Promise<void>}
     */
    SetRelativePosition: (x, y) => wails.Window.SetRelativePosition(x, y),

    /**
     * Get the window position relative to the current monitor.
     * @returns {Promise<Position>} The window position
     */
    RelativePosition: () => wails.Window.RelativePosition(),

    /**
     * Set the absolute window position.
     * @param {number} x
     * @param {number} y
     * @returns {Promise<void>}
     */
    SetAbsolutePosition: (x, y) => wails.Window.SetAbsolutePosition(x, y),

    /**
     * Get the absolute window position.
     * @returns {Promise<Position>} The window position
     */
    AbsolutePosition: () => wails.Window.AbsolutePosition(),

    /**
     * Get the screen the window is on.
     * @returns {Promise<Screen>}
     */
    Screen: () => wails.Window.Screen(),

    /**
     * Hide the window
     * @returns {Promise<void>}
     */
    Hide: () => wails.Window.Hide(),

    /**
     * Maximise the window
     * @returns {Promise<void>}
     */
    Maximise: () => wails.Window.Maximise(),

    /**
     * Show the window
     * @returns {Promise<void>}
     */
    Show: () => wails.Window.Show(),

    /**
     * Close the window
     * @returns {Promise<void>}
     */
    Close: () => wails.Window.Close(),

    /**
     * Toggle the window maximise state
     * @returns {Promise<void>}
     */
    ToggleMaximise: () => wails.Window.ToggleMaximise(),

    /**
     * Unmaximise the window
     * @returns {Promise<void>}
     */
    UnMaximise: () => wails.Window.UnMaximise(),

    /**
     * Minimise the window
     * @returns {Promise<void>}
     */
    Minimise: () => wails.Window.Minimise(),

    /**
     * Unminimise the window
     * @returns {Promise<void>}
     */
    UnMinimise: () => wails.Window.UnMinimise(),

    /**
     * Set the background colour of the window.
     * @param {number} r - The red value between 0 and 255
     * @param {number} g - The green value between 0 and 255
     * @param {number} b - The blue value between 0 and 255
     * @param {number} a - The alpha value between 0 and 255
     * @returns {Promise<void>}
     */
    SetBackgroundColour: (r, g, b, a) => wails.Window.SetBackgroundColour(r, g, b, a),
};
