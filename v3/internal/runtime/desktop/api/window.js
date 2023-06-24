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
 * The Window API provides methods to interact with the window.
 */
export const Window = {
    /**
     * Center the window.
     */
    Center: () => void wails.Window.Center(),
    /**
     * Set the window title.
     * @param title
     */
    SetTitle: (title) => void wails.Window.SetTitle(title),

    /**
     * Makes the window fullscreen.
     */
    Fullscreen: () => void wails.Window.Fullscreen(),

    /**
     * Unfullscreen the window.
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
     */
    SetMaxSize: (width, height) => void wails.Window.SetMaxSize(width, height),

    /**
     * Set the window minimum size.
     * @param {number} width
     * @param {number} height
     */
    SetMinSize: (width, height) => void wails.Window.SetMinSize(width, height),

    /**
     * Set window to be always on top.
     * @param {boolean} onTop Whether the window should be always on top
     */
    SetAlwaysOnTop: (onTop) => void wails.Window.SetAlwaysOnTop(onTop),

    /**
     * Set the window position relative to the current monitor.
     * @param {number} x
     * @param {number} y
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
     */
    Hide: () => void wails.Window.Hide(),

    /**
     * Maximise the window
     */
    Maximise: () => void wails.Window.Maximise(),

    /**
     * Show the window
     */
    Show: () => void wails.Window.Show(),

    /**
     * Close the window
     */
    Close: () => void wails.Window.Close(),

    /**
     * Toggle the window maximise state
     */
    ToggleMaximise: () => void wails.Window.ToggleMaximise(),

    /**
     * Unmaximise the window
     */
    UnMaximise: () => void wails.Window.UnMaximise(),

    /**
     * Minimise the window
     */
    Minimise: () => void wails.Window.Minimise(),

    /**
     * Unminimise the window
     */
    UnMinimise: () => void wails.Window.UnMinimise(),

    /**
     * Set the background colour of the window.
     * @param {number} r - The red value between 0 and 255
     * @param {number} g - The green value between 0 and 255
     * @param {number} b - The blue value between 0 and 255
     * @param {number} a - The alpha value between 0 and 255
     */
    SetBackgroundColour: (r, g, b, a) => void wails.Window.SetBackgroundColour(r, g, b, a),
};
