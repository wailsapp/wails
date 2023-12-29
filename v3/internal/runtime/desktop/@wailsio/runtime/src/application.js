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

import { newRuntimeCallerWithID, objectNames } from "./runtime";
const call = newRuntimeCallerWithID(objectNames.Application, '');

const HideMethod = 0;
const ShowMethod = 1;
const QuitMethod = 2;

/**
 * Hides a certain method by calling the HideMethod function.
 *
 * @return {Promise<void>}
 *
 */
export function Hide() {
    return call(HideMethod);
}

/**
 * Calls the ShowMethod and returns the result.
 *
 * @return {Promise<void>}
 */
export function Show() {
    return call(ShowMethod);
}

/**
 * Calls the QuitMethod to terminate the program.
 *
 * @return {Promise<void>}
 */
export function Quit() {
    return call(QuitMethod);
}
