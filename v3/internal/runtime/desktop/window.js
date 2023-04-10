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
 * @typedef {import("../api/types").Size} Size
 * @typedef {import("../api/types").Position} Position
 * @typedef {import("../api/types").Screen} Screen
 */

import {newRuntimeCaller} from "./runtime";

export function newWindow(windowName) {
    let call = newRuntimeCaller("window", windowName);
    return {
        // Reload: () => call('WR'),
        // ReloadApp: () => call('WR'),
        // SetSystemDefaultTheme: () => call('WASDT'),
        // SetLightTheme: () => call('WALT'),
        // SetDarkTheme: () => call('WADT'),
        // IsFullscreen: () => call('WIF'),
        // IsMaximized: () => call('WIM'),
        // IsMinimized: () => call('WIMN'),
        // IsWindowed: () => call('WIF'),


        /**
         * Centers the window.
         */
        Center: () => void call('Center'),

        /**
         * Set the window title.
         * @param title
         */
        SetTitle: (title) => void call('SetTitle', {title}),

        /**
         * Makes the window fullscreen.
         */
        Fullscreen: () => void call('Fullscreen'),

        /**
         * Unfullscreen the window.
         */
        UnFullscreen: () => void call('UnFullscreen'),

        /**
         * Set the window size.
         * @param {number} width The window width
         * @param {number} height The window height
         */
        SetSize: (width, height) => call('SetSize', {width,height}),

        /**
         * Get the window size.
         * @returns {Promise<Size>} The window size
         */
        Size: () => { return call('Size'); },

        /**
         * Set the window maximum size.
         * @param {number} width
         * @param {number} height
         */
        SetMaxSize: (width, height) => void call('SetMaxSize', {width,height}),

        /**
         * Set the window minimum size.
         * @param {number} width
         * @param {number} height
         */
        SetMinSize: (width, height) => void call('SetMinSize', {width,height}),

        /**
         * Set window to be always on top.
         * @param {boolean} onTop Whether the window should be always on top
         */
        SetAlwaysOnTop: (onTop) => void call('SetAlwaysOnTop', {alwaysOnTop:onTop}),

        /**
         * Set the window position.
         * @param {number} x
         * @param {number} y
         */
        SetPosition: (x, y) => call('SetPosition', {x,y}),

        /**
         * Get the window position.
         * @returns {Promise<Position>} The window position
         */
        Position: () => { return call('Position'); },

        /**
         * Get the screen the window is on.
         * @returns {Promise<Screen>}
         */
        Screen: () => { return call('Screen'); },

        /**
         * Hide the window
         */
        Hide: () => void call('Hide'),

        /**
         * Maximise the window
         */
        Maximise: () => void call('Maximise'),

        /**
         * Show the window
         */
        Show: () => void call('Show'),

        /**
         * Close the window
         */
        Close: () => void call('Close'),

        /**
         * Toggle the window maximise state
         */
        ToggleMaximise: () => void call('ToggleMaximise'),

        /**
         * Unmaximise the window
         */
        UnMaximise: () => void call('UnMaximise'),

        /**
         * Minimise the window
         */
        Minimise: () => void call('Minimise'),

        /**
         * Unminimise the window
         */
        UnMinimise: () => void call('UnMinimise'),

        /**
         * Set the background colour of the window.
         * @param {number} r - A value between 0 and 255
         * @param {number} g - A value between 0 and 255
         * @param {number} b - A value between 0 and 255
         * @param {number} a - A value between 0 and 255
         */
        SetBackgroundColour: (r, g, b, a) => void call('SetBackgroundColour', {r, g, b, a}),
    };
}
