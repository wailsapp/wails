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

import {EventsOn} from "./events";

/**
 * Callback for DragAndDropOnMotion returns X and Y coordinates of the mouse position inside the window while it is dragging something.
 *
 * @export
 * @callback DragAndDropOnMotionCallback
 * @param {int} x - coordinate of mouse position inside the window.
 * @param {int} y - coordinate of mouse position inside the window.
 */

/**
 * DragAndDropOnMotion calls a callback with X and Y coordinates of the mouse position inside the window while it is dragging something.
 *
 * @export
 * @param {DragAndDropOnMotionCallback} callback
 * @returns {function} - A function to cancel the listener
 */
export function DragAndDropOnMotion(callback) {
    return EventsOn("wails.dnd.motion", callback);
}

/**
 * Callback for DragAndDropOnDrop returns a slice of file path strings when a drop is finished.
 *
 * @export
 * @callback DragAndDropOnDropCallback
 * @param {String[]} paths - A list of file paths.
 */

/**
 * DragAndDropOnDrop calls a callback with slice of file path strings when a drop is finished.
 *
 * @export
 * @param {DragAndDropOnDropCallback} callback
 * @returns {function} - A function to cancel the listener
 */
export function DragAndDropOnDrop(callback) {
    return EventsOn("wails.dnd.drop", callback);
}