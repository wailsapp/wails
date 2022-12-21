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

/**
 * Set the Size of the window
 *
 * @export
 * @param {string} text
 */
export function ClipboardSetText(text) {
    return Call(":wails:ClipboardSetText", [text]);
}

/**
 * Get the text content of the clipboard
 *
 * @export
 * @return {Promise<{string}>} Text content of the clipboard

 */
export function ClipboardGetText() {
    return Call(":wails:ClipboardGetText");
}