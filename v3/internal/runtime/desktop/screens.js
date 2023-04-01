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

/**
 * @typedef {import("./api/types").Screen} Screen
 */

import {newRuntimeCaller} from "./runtime";

let call = newRuntimeCaller("screens");

/**
 * Gets all screens.
 * @returns {Promise<Screen[]>}
 */
export function GetAll() {
    return call("GetAll");
}

/**
 * Gets the primary screen.
 * @returns {Promise<Screen>}
 */
export function GetPrimary() {
    return call("GetPrimary");
}

/**
 * Gets the current active screen.
 * @returns {Promise<Screen>}
 * @constructor
 */
export function GetCurrent() {
    return call("GetCurrent");
}