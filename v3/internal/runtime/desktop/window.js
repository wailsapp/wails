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

import {newRuntimeCallerWithID, objectNames} from "./runtime";

let WindowCenter = 0;
let WindowSetTitle = 1;
let WindowFullscreen = 2;
let WindowUnFullscreen = 3;
let WindowSetSize = 4;
let WindowSize = 5;
let WindowSetMaxSize = 6;
let WindowSetMinSize = 7;
let WindowSetAlwaysOnTop = 8;
let WindowSetRelativePosition = 9;
let WindowRelativePosition = 10;
let WindowScreen = 11;
let WindowHide = 12;
let WindowMaximise = 13;
let WindowUnMaximise = 14;
let WindowToggleMaximise = 15;
let WindowMinimise = 16;
let WindowUnMinimise = 17;
let WindowRestore = 18;
let WindowShow = 19;
let WindowClose = 20;
let WindowSetBackgroundColour = 21;
let WindowSetResizable = 22;
let WindowWidth = 23;
let WindowHeight = 24;
let WindowZoomIn = 25;
let WindowZoomOut = 26;
let WindowZoomReset = 27;
let WindowGetZoomLevel = 28;
let WindowSetZoomLevel = 29;

export function newWindow(windowName) {
    let call = newRuntimeCallerWithID(objectNames.Window, windowName);
    return {

        /**
         * Centers the window.
         * @returns {Promise<void>}
         */
        Center: () => call(WindowCenter),

        /**
         * Set the window title.
         * @param title
         * @returns {Promise<void>}
         */
        SetTitle: (title) => call(WindowSetTitle, {title}),

        /**
         * Makes the window fullscreen.
         * @returns {Promise<void>}
         */
        Fullscreen: () => call(WindowFullscreen),

        /**
         * Unfullscreen the window.
         * @returns {Promise<void>}
         */
        UnFullscreen: () => call(WindowUnFullscreen),

        /**
         * Set the window size.
         * @param {number} width The window width
         * @param {number} height The window height
         * @returns {Promise<void>}
         */
        SetSize: (width, height) => call(WindowSetSize, {width,height}),

        /**
         * Get the window size.
         * @returns {Promise<Size>} The window size
         */
        Size: () => call(WindowSize),

        /**
         * Set the window maximum size.
         * @param {number} width
         * @param {number} height
         * @returns {Promise<void>}
         */
        SetMaxSize: (width, height) => call(WindowSetMaxSize, {width,height}),

        /**
         * Set the window minimum size.
         * @param {number} width
         * @param {number} height
         * @returns {Promise<void>}
         */
        SetMinSize: (width, height) => call(WindowSetMinSize, {width,height}),

        /**
         * Set window to be always on top.
         * @param {boolean} onTop Whether the window should be always on top
         * @returns {Promise<void>}
         */
        SetAlwaysOnTop: (onTop) => call(WindowSetAlwaysOnTop, {alwaysOnTop:onTop}),

        /**
         * Set the window relative position.
         * @param {number} x
         * @param {number} y
         * @returns {Promise<void>}
         */
        SetRelativePosition: (x, y) => call(WindowSetRelativePosition, {x,y}),

        /**
         * Get the window position.
         * @returns {Promise<Position>} The window position
         */
        RelativePosition: () => call(WindowRelativePosition),

        /**
         * Get the screen the window is on.
         * @returns {Promise<Screen>}
         */
        Screen: () => call(WindowScreen),

        /**
         * Hide the window
         * @returns {Promise<void>}
         */
        Hide: () => call(WindowHide),

        /**
         * Maximise the window
         * @returns {Promise<void>}
         */
        Maximise: () => call(WindowMaximise),

        /**
         * Show the window
         * @returns {Promise<void>}
         */
        Show: () => call(WindowShow),

        /**
         * Close the window
         * @returns {Promise<void>}
         */
        Close: () => call(WindowClose),

        /**
         * Toggle the window maximise state
         * @returns {Promise<void>}
         */
        ToggleMaximise: () => call(WindowToggleMaximise),

        /**
         * Unmaximise the window
         * @returns {Promise<void>}
         */
        UnMaximise: () => call(WindowUnMaximise),

        /**
         * Minimise the window
         * @returns {Promise<void>}
         */
        Minimise: () => call(WindowMinimise),

        /**
         * Unminimise the window
         * @returns {Promise<void>}
         */
        UnMinimise: () => call(WindowUnMinimise),

        /**
         * Restore the window
         * @returns {Promise<void>}
         */
        Restore: () => call(WindowRestore),

        /**
         * Set the background colour of the window.
         * @param {number} r - A value between 0 and 255
         * @param {number} g - A value between 0 and 255
         * @param {number} b - A value between 0 and 255
         * @param {number} a - A value between 0 and 255
         * @returns {Promise<void>}
         */
        SetBackgroundColour: (r, g, b, a) => call(WindowSetBackgroundColour, {r, g, b, a}),

        /**
         * Set whether the window can be resized or not
         * @param {boolean} resizable
         * @returns {Promise<void>}
         */
        SetResizable: (resizable) => call(WindowSetResizable, {resizable}),

        /**
         * Get the window width
         * @returns {Promise<number>}
         */
        Width: () => call(WindowWidth),

        /**
         * Get the window height
         * @returns {Promise<number>}
         */
        Height: () => call(WindowHeight),

        /**
         * Zoom in the window
         * @returns {Promise<void>}
         */
        ZoomIn: () => call(WindowZoomIn),

        /**
         * Zoom out the window
         * @returns {Promise<void>}
         */
        ZoomOut: () => call(WindowZoomOut),

        /**
         * Reset the window zoom
         * @returns {Promise<void>}
         */
        ZoomReset: () => call(WindowZoomReset),

        /**
         * Get the window zoom
         * @returns {Promise<number>}
         */
        GetZoomLevel: () => call(WindowGetZoomLevel),

        /**
         * Set the window zoom level
         * @param {number} zoomLevel
         * @returns {Promise<void>}
         */
        SetZoomLevel: (zoomLevel) => call(WindowSetZoomLevel, {zoomLevel}),
    };
}
