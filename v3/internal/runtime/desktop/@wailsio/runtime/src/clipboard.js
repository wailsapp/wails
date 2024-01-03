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

const call = newRuntimeCallerWithID(objectNames.Clipboard, '');
const ClipboardSetText = 0;
const ClipboardText = 1;

/**
 * Sets the text to the Clipboard.
 *
 * @param {string} text - The text to be set to the Clipboard.
 * @return {Promise} - A Promise that resolves when the operation is successful.
 */
export function SetText(text) {
    return call(ClipboardSetText, {text});
}

/**
 * Get the Clipboard text
 * @returns {Promise<string>} A promise that resolves with the text from the Clipboard.
 */
export function Text() {
    return call(ClipboardText);
}
