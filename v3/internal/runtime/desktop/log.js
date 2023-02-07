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

let call = newRuntimeCaller("log");

/**
 * Logs a message.
 * @param {message} Message to log
 */
export function Log(message) {
    return call("Log", message);
}
