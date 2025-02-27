/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/
import { newRuntimeCaller, objectNames } from "./runtime.js";
const call = newRuntimeCaller(objectNames.Screens);
const getAll = 0;
const getPrimary = 1;
const getCurrent = 2;
/**
 * Gets all screens.
 *
 * @returns A promise that resolves to an array of Screen objects.
 */
export function GetAll() {
    return call(getAll);
}
/**
 * Gets the primary screen.
 *
 * @returns A promise that resolves to the primary screen.
 */
export function GetPrimary() {
    return call(getPrimary);
}
/**
 * Gets the current active screen.
 *
 * @returns A promise that resolves with the current active screen.
 */
export function GetCurrent() {
    return call(getCurrent);
}
