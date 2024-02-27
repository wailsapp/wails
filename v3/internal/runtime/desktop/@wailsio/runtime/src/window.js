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

const thisWindow = Get('');

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
 * Centers the window on the screen.
 */
export function Center() {
    thisWindow.Center();
}

/**
 * Sets the title of the window.
 * @param {string} title - The title to set.
 */
export function SetTitle(title) {
    thisWindow.SetTitle(title);
}

/**
 * Sets the window to fullscreen.
 */
export function Fullscreen() {
    thisWindow.Fullscreen();
}

/**
 * Restores the previous window dimensions and position prior to fullscreen.
 */
export function UnFullscreen() {
    thisWindow.UnFullscreen();
}

/**
 * Sets the size of the window.
 * @param {number} width - The width of the window.
 * @param {number} height - The height of the window.
 */
export function SetSize(width, height) {
    thisWindow.SetSize(width, height);
}

/**
 * Gets the size of the window.
 */
export function Size() {
    return thisWindow.Size();
}

/**
 * Sets the maximum size of the window.
 * @param {number} width - The maximum width of the window.
 * @param {number} height - The maximum height of the window.
 */
export function SetMaxSize(width, height) {
    thisWindow.SetMaxSize(width, height);
}

/**
 * Sets the minimum size of the window.
 * @param {number} width - The minimum width of the window.
 * @param {number} height - The minimum height of the window.
 */
export function SetMinSize(width, height) {
    thisWindow.SetMinSize(width, height);
}

/**
 * Sets the window to always be on top.
 * @param {boolean} onTop - Whether the window should always be on top.
 */
export function SetAlwaysOnTop(onTop) {
    thisWindow.SetAlwaysOnTop(onTop);
}

/**
 * Sets the relative position of the window.
 * @param {number} x - The x-coordinate of the window's position.
 * @param {number} y - The y-coordinate of the window's position.
 */
export function SetRelativePosition(x, y) {
    thisWindow.SetRelativePosition(x, y);
}

/**
 * Gets the relative position of the window.
 */
export function RelativePosition() {
    return thisWindow.RelativePosition();
}

/**
 * Gets the screen that the window is on.
 */
export function Screen() {
    return thisWindow.Screen();
}

/**
 * Hides the window.
 */
export function Hide() {
    thisWindow.Hide();
}

/**
 * Maximises the window.
 */
export function Maximise() {
    thisWindow.Maximise();
}

/**
 * Un-maximises the window.
 */
export function UnMaximise() {
    thisWindow.UnMaximise();
}

/**
 * Toggles the maximisation of the window.
 */
export function ToggleMaximise() {
    thisWindow.ToggleMaximise();
}

/**
 * Minimises the window.
 */
export function Minimise() {
    thisWindow.Minimise();
}

/**
 * Un-minimises the window.
 */
export function UnMinimise() {
    thisWindow.UnMinimise();
}

/**
 * Restores the window.
 */
export function Restore() {
    thisWindow.Restore();
}

/**
 * Shows the window.
 */
export function Show() {
    thisWindow.Show();
}

/**
 * Closes the window.
 */
export function Close() {
    thisWindow.Close();
}

/**
 * Sets the background colour of the window.
 * @param {number} r - The red component of the colour.
 * @param {number} g - The green component of the colour.
 * @param {number} b - The blue component of the colour.
 * @param {number} a - The alpha component of the colour.
 */
export function SetBackgroundColour(r, g, b, a) {
    thisWindow.SetBackgroundColour(r, g, b, a);
}

/**
 * Sets whether the window is resizable.
 * @param {boolean} resizable - Whether the window should be resizable.
 */
export function SetResizable(resizable) {
    thisWindow.SetResizable(resizable);
}

/**
 * Gets the width of the window.
 */
export function Width() {
    return thisWindow.Width();
}

/**
 * Gets the height of the window.
 */
export function Height() {
    return thisWindow.Height();
}

/**
 * Zooms in the window.
 */
export function ZoomIn() {
    thisWindow.ZoomIn();
}

/**
 * Zooms out the window.
 */
export function ZoomOut() {
    thisWindow.ZoomOut();
}

/**
 * Resets the zoom of the window.
 */
export function ZoomReset() {
    thisWindow.ZoomReset();
}

/**
 * Gets the zoom level of the window.
 */
export function GetZoomLevel() {
    return thisWindow.GetZoomLevel();
}

/**
 * Sets the zoom level of the window.
 * @param {number} zoomLevel - The zoom level to set.
 */
export function SetZoomLevel(zoomLevel) {
    thisWindow.SetZoomLevel(zoomLevel);
}
