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

import {newRuntimeCallerWithID, objectNames} from "./runtime";

let call = newRuntimeCallerWithID(objectNames.Screens);

let ScreensGetAll = 0;
let ScreensGetPrimary = 1;
let ScreensGetCurrent = 2;

/**
 * Gets all screens.
 * @returns {Promise<Screen[]>}
 */
export function GetAll() {
    return call(ScreensGetAll);
}

/**
 * Gets the primary screen.
 * @returns {Promise<Screen>}
 */
export function GetPrimary() {
    return call(ScreensGetPrimary);
}

/**
 * Gets the current active screen.
 * @returns {Promise<Screen>}
 */
export function GetCurrent() {
    return call(ScreensGetCurrent);
}
