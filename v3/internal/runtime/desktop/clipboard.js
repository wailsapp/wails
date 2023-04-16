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

import {newRuntimeCaller} from "./runtime";

let call = newRuntimeCaller("clipboard");

/**
 * Set the Clipboard text
 */
export function SetText(text) {
    void call("SetText", {text});
}

/**
 * Get the Clipboard text
 * @returns {Promise<string>}
 */
export function Text() {
    return call("Text");
}