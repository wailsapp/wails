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

/**
 * Place the window in the center of the screen
 *
 * @export
 */
export function Center() {
	window.wails.Window.Center();
}

/**
 * Set the Size of the window
 *
 * @export
 * @param {number} width
 * @param {number} height
 */
export function SetSize(width, height) {
	window.wails.Window.SetSize(width, height);
}

/**
 * Set the Position of the window
 *
 * @export
 * @param {number} x
 * @param {number} y
 */
export function SetPosition(x, y) {
	window.wails.Window.SetPosition(x, y);
}

/**
 * Hide the Window
 *
 * @export
 */
export function Hide() {
	window.wails.Window.Hide();
}

/**
 * Show the Window
 *
 * @export
 */
export function Show() {
	window.wails.Window.Show();
}

/**
 * Maximise the Window
 *
 * @export
 */
export function Maximise() {
	window.wails.Window.Maximise()
}

/**
 * Unmaximise the Window
 *
 * @export
 */
export function Unmaximise() {
	window.wails.Window.Unmaximise()
}

/**
 * Minimise the Window
 *
 * @export
 */
export function Minimise() {
	window.wails.Window.Minimise();
}

/**
 * Unminimise the Window
 *
 * @export
 */
export function Unminimise() {
	window.wails.Window.Unminimise();
}

/**
 * Close the Window
 *
 * @export
 */
export function Close() {
	window.wails.Window.Close();
}
