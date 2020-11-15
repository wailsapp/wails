/*
 _       __      _ __
| |     / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The lightweight framework for web-like apps
(c) Lea Anthony 2019-present
*/

/* jshint esversion: 6 */

import { SendMessage } from 'ipc';

/**
 * Place the window in the center of the screen
 *
 * @export
 */
export function Center() {
	SendMessage('Wc');
}

/**
 * Sets the window title
 *
 * @param {string} title
 * @export
 */
export function SetTitle(title) {
	SendMessage('WT' + title);
}

/**
 * Makes the window go fullscreen
 *
 * @export
 */
export function Fullscreen() {
	SendMessage('WF');
}

/**
 * Reverts the window from fullscreen
 *
 * @export
 */
export function UnFullscreen() {
	SendMessage('Wf');
}

/**
 * Set the Size of the window
 *
 * @export
 * @param {number} width
 * @param {number} height
 */
export function SetSize(width, height) {
	SendMessage('Ws:' + width + ':' + height);
}

/**
 * Set the Position of the window
 *
 * @export
 * @param {number} x
 * @param {number} y
 */
export function SetPosition(x, y) {
	SendMessage('Wp:' + x + ':' + y);
}

/**
 * Hide the Window
 *
 * @export
 */
export function Hide() {
	SendMessage('WH');
}

/**
 * Show the Window
 *
 * @export
 */
export function Show() {
	SendMessage('WS');
}

/**
 * Maximise the Window
 *
 * @export
 */
export function Maximise() {
	SendMessage('WM');
}

/**
 * Unmaximise the Window
 *
 * @export
 */
export function Unmaximise() {
	SendMessage('WU');
}

/**
 * Minimise the Window
 *
 * @export
 */
export function Minimise() {
	SendMessage('Wm');
}

/**
 * Unminimise the Window
 *
 * @export
 */
export function Unminimise() {
	SendMessage('Wu');
}

/**
 * Close the Window
 *
 * @export
 */
export function Close() {
	SendMessage('WC');
}
