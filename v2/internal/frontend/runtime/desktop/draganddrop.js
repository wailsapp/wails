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
 * postMessageWithAdditionalObjects checks the browser's capability of sending postMessageWithAdditionalObjects
 *
 * @returns {boolean}
 * @constructor
 */
export function CanResolveFilePaths() {
    return window.chrome?.webview?.postMessageWithAdditionalObjects != null;
}

export function ResolveFilePaths(x, y, files) {
    // Only for windows webview2 >= 1.0.1774.30
    // https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2webmessagereceivedeventargs2?view=webview2-1.0.1823.32#applies-to
    if (window.chrome?.webview?.postMessageWithAdditionalObjects) {
        chrome.webview.postMessageWithAdditionalObjects(`file:drop:${x}:${y}`, files);
        return;
    }
    console.warn("unsupported platform");
}

/**
 * Callback for DragAndDropOnDrop returns a slice of file path strings when a drop is finished.
 *
 * @export
 * @callback HandleDragAndDropCallback
 * @param {Number} x - x coordinate of the drop
 * @param {Number} y - y coordinate of the drop
 * @param {String[]} paths - A list of file paths.
 */

/**
 * HandleDragAndDrop calls a callback with slice of file path strings when a drop is finished.
 *
 * @export
 * @param {HandleDragAndDropCallback} callback
 * @returns {function} - A function to cancel the listener
 */
export function HandleDragAndDrop(callback) {
    return EventsOn("wails.dnd.drop", callback);
}