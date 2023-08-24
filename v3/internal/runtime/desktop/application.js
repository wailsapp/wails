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

let call = newRuntimeCallerWithID(objectNames.Application);

let methods = {
    Hide: 0,
    Show: 1,
    Quit: 2,
}

/**
 * Hide the application
 */
export function Hide() {
    void call(methods.Hide);
}

/**
 * Show the application
 */
export function Show() {
    void call(methods.Show);
}


/**
 * Quit the application
 */
export function Quit() {
    void call(methods.Quit);
}