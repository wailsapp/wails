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


import {Call} from "./calls";

export function WindowReload() {
    window.location.reload();
}

/**
 * Place the window in the center of the screen
 *
 * @export
 */
export function WindowCenter() {
    window.WailsInvoke('Wc');
}

/**
 * Sets the window title
 *
 * @param {string} title
 * @export
 */
export function WindowSetTitle(title) {
    window.WailsInvoke('WT' + title);
}

/**
 * Makes the window go fullscreen
 *
 * @export
 */
export function WindowFullscreen() {
    window.WailsInvoke('WF');
}

/**
 * Reverts the window from fullscreen
 *
 * @export
 */
export function WindowUnFullscreen() {
    window.WailsInvoke('Wf');
}

/**
 * Set the Size of the window
 *
 * @export
 * @param {number} width
 * @param {number} height
 */
export function WindowSetSize(width, height) {
    window.WailsInvoke('Ws:' + width + ':' + height);
}

/**
 * Get the Size of the window
 *
 * @export
 * @return {Promise<{w: number, h: number}>} The size of the window

 */
export function WindowGetSize() {
    return Call(":wails:WindowGetSize");
}

/**
 * Set the maximum size of the window
 *
 * @export
 * @param {number} width
 * @param {number} height
 */
export function WindowSetMaxSize(width, height) {
    window.WailsInvoke('WZ:' + width + ':' + height);
}

/**
 * Set the minimum size of the window
 *
 * @export
 * @param {number} width
 * @param {number} height
 */
export function WindowSetMinSize(width, height) {
    window.WailsInvoke('Wz:' + width + ':' + height);
}

/**
 * Set the Position of the window
 *
 * @export
 * @param {number} x
 * @param {number} y
 */
export function WindowSetPosition(x, y) {
    window.WailsInvoke('Wp:' + x + ':' + y);
}

/**
 * Get the Position of the window
 *
 * @export
 * @return {Promise<{x: number, y: number}>} The position of the window
 */
export function WindowGetPosition() {
    return Call(":wails:WindowGetPos");
}

/**
 * Hide the Window
 *
 * @export
 */
export function WindowHide() {
    window.WailsInvoke('WH');
}

/**
 * Show the Window
 *
 * @export
 */
export function WindowShow() {
    window.WailsInvoke('WS');
}

/**
 * Maximise the Window
 *
 * @export
 */
export function WindowMaximise() {
    window.WailsInvoke('WM');
}

/**
 * Unmaximise the Window
 *
 * @export
 */
export function WindowUnmaximise() {
    window.WailsInvoke('WU');
}

/**
 * Minimise the Window
 *
 * @export
 */
export function WindowMinimise() {
    window.WailsInvoke('Wm');
}

/**
 * Unminimise the Window
 *
 * @export
 */
export function WindowUnminimise() {
    window.WailsInvoke('Wu');
}


/**
 * Sets the background colour of the window
 *
 * @export
 * @param {RGBA} RGBA background colour
 */
export function WindowSetRGBA(RGBA) {
    let rgba = JSON.stringify(RGBA);
    window.WailsInvoke('Wr:' + rgba);
}

