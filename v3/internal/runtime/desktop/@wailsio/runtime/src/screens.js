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
const call = newRuntimeCallerWithID(objectNames.Screens, '');

const getAll = 0;
const getPrimary = 1;
const getCurrent = 2;

/**
 * Gets all screens.
 * @returns {Promise<Screen[]>} A promise that resolves to an array of Screen objects.
 */
export function GetAll() {
    return call(getAll);
}
/**
 * Gets the primary screen.
 * @returns {Promise<Screen>} A promise that resolves to the primary screen.
 */
export function GetPrimary() {
    return call(getPrimary);
}
/**
 * Gets the current active screen.
 *
 * @returns {Promise<Screen>} A promise that resolves with the current active screen.
 */
export function GetCurrent() {
    return call(getCurrent);
}