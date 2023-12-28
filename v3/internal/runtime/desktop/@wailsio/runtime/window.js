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

const center = 0;
const setTitle = 1;
const fullscreen = 2;
const unFullscreen = 3;
const setSize = 4;
const size = 5;
const setMaxSize = 6;
const setMinSize = 7;
const setAlwaysOnTop = 8;
const setRelativePosition = 9;
const relativePosition = 10;
const screen = 11;
const hide = 12;
const maximise = 13;
const unMaximise = 14;
const toggleMaximise = 15;
const minimise = 16;
const unMinimise = 17;
const restore = 18;
const show = 19;
const close = 20;
const setBackgroundColour = 21;
const setResizable = 22;
const width = 23;
const height = 24;
const zoomIn = 25;
const zoomOut = 26;
const zoomReset = 27;
const getZoomLevel = 28;
const setZoomLevel = 29;

const thisWindow = newRuntimeCallerWithID(objectNames.Window, '');

function createWindow(call) {
    return {
        Get: (windowName) => createWindow(newRuntimeCallerWithID(objectNames.Window, windowName)),
        Center: () => call(center),
        SetTitle: (title) => call(setTitle, {title}),
        Fullscreen: () => call(fullscreen),
        UnFullscreen: () => call(unFullscreen),
        SetSize: (width, height) => call(setSize, {width, height}),
        Size: () => call(size),
        SetMaxSize: (width, height) => call(setMaxSize, {width, height}),
        SetMinSize: (width, height) => call(setMinSize, {width, height}),
        SetAlwaysOnTop: (onTop) => call(setAlwaysOnTop, {alwaysOnTop: onTop}),
        SetRelativePosition: (x, y) => call(setRelativePosition, {x, y}),
        RelativePosition: () => call(relativePosition),
        Screen: () => call(screen),
        Hide: () => call(hide),
        Maximise: () => call(maximise),
        UnMaximise: () => call(unMaximise),
        ToggleMaximise: () => call(toggleMaximise),
        Minimise: () => call(minimise),
        UnMinimise: () => call(unMinimise),
        Restore: () => call(restore),
        Show: () => call(show),
        Close: () => call(close),
        SetBackgroundColour: (r, g, b, a) => call(setBackgroundColour, {r, g, b, a}),
        SetResizable: (resizable) => call(setResizable, {resizable}),
        Width: () => call(width),
        Height: () => call(height),
        ZoomIn: () => call(zoomIn),
        ZoomOut: () => call(zoomOut),
        ZoomReset: () => call(zoomReset),
        GetZoomLevel: () => call(getZoomLevel),
        SetZoomLevel: (zoomLevel) => call(setZoomLevel, {zoomLevel}),
    };
}

/**
 * Gets the specified window.
 *
 * @param {string} windowName - The name of the window to get.
 * @return {Object} - The specified window object.
 */
export function Get(windowName) {
    return createWindow(newRuntimeCallerWithID(objectNames.Window, windowName));
}

/**
 * Returns a map of all methods in the current window.
 * @returns {Map} - A map of window methods.
 */
export function WindowMethods(targetWindow) {
    // Create a new map to store methods
    let result = new Map();

    // Iterate over all properties of the window object
    for (let method in targetWindow) {
        // Check if the property is indeed a method (function)
        if(typeof targetWindow[method] === 'function') {
            // Add the method to the map
            result.set(method, targetWindow[method]);
        }

    }
    // Return the map of window methods
    return result;
}
export default {
    ...Get('')
}
