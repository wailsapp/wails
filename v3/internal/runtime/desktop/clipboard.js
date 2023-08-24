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

import {newRuntimeCallerWithID, objectNames} from "./runtime";

let call = newRuntimeCallerWithID(objectNames.Clipboard);

let ClipboardSetText = 0;
let ClipboardText = 1;

/**
 * Set the Clipboard text
 */
export function SetText(text) {
    void call(ClipboardSetText, {text});
}

/**
 * Get the Clipboard text
 * @returns {Promise<string>}
 */
export function Text() {
    return call(ClipboardText);
}