/*
 _       __      _ __    
| |     / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  ) 
|__/|__/\__,_/_/_/____/  
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

/* jshint esversion: 9 */

/**
 * Reloads the Window
 *
 * @export
 */
export function WindowReload() {
	window.runtime.WindowReload();
}

/**
 * Place the window in the center of the screen
 *
 * @export
 */
export function WindowCenter() {
	window.runtime.WindowCenter();
}

/**
 * Sets the window title
 *
 * @param {string} title
 * @export
 */
export function WindowSetTitle(title) {
	window.runtime.WindowSetTitle(title);
}

/**
 * Makes the window go fullscreen
 *
 * @export
 */
export function WindowFullscreen() {
	window.runtime.WindowFullscreen();
}

/**
 * Reverts the window from fullscreen
 *
 * @export
 */
export function WindowUnfullscreen() {
	window.runtime.WindowUnfullscreen();
}

/**
 * Get the Size of the window
 *
 * @export
 * @return {Promise<{w: number, h: number}>} The size of the window

 */
export function WindowGetSize() {
	window.runtime.WindowGetSize();
}


/**
 * Set the Size of the window
 *
 * @export
 * @param {number} width
 * @param {number} height
 */
export function WindowSetSize(width, height) {
	window.runtime.WindowSetSize(width, height);
}

/**
 * Set the maximum size of the window
 *
 * @export
 * @param {number} width
 * @param {number} height
 */
export function WindowSetMaxSize(width, height) {
	window.runtime.WindowSetMaxSize(width, height);
}

/**
 * Set the minimum size of the window
 *
 * @export
 * @param {number} width
 * @param {number} height
 */
export function WindowSetMinSize(width, height) {
	window.runtime.WindowSetMinSize(width, height);
}

/**
 * Set the Position of the window
 *
 * @export
 * @param {number} x
 * @param {number} y
 */
export function WindowSetPosition(x, y) {
	window.runtime.WindowSetPosition(x, y);
}

/**
 * Get the Position of the window
 *
 * @export
 * @return {Promise<{x: number, y: number}>} The position of the window
 */
export function WindowGetPosition() {
	window.runtime.WindowGetPosition();
}

/**
 * Hide the Window
 *
 * @export
 */
export function WindowHide() {
	window.runtime.WindowHide();
}

/**
 * Show the Window
 *
 * @export
 */
export function WindowShow() {
	window.runtime.WindowShow();
}

/**
 * Maximise the Window
 *
 * @export
 */
export function WindowMaximise() {
	window.runtime.WindowMaximise();
}

/**
 * Toggle the Maximise of the Window
 *
 * @export
 */
export function WindowToggleMaximise() {
	window.runtime.WindowToggleMaximise();
}

/**
 * Unmaximise the Window
 *
 * @export
 */
export function WindowUnmaximise() {
	window.runtime.WindowUnmaximise();
}

/**
 * Minimise the Window
 *
 * @export
 */
export function WindowMinimise() {
	window.runtime.WindowMinimise();
}

/**
 * Unminimise the Window
 *
 * @export
 */
export function WindowUnminimise() {
	window.runtime.WindowUnminimise();
}

/**
 * Sets the background colour of the window
 *
 * @export
 * @param {RGBA} RGBA background colour
 */
export function WindowSetRGBA(RGBA) {
	window.runtime.WindowSetRGBA(RGBA);
}
