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

let call = newRuntimeCaller("application");

/**
 * Hide the application
 */
export function Hide() {
    void call("Hide");
}

/**
 * Show the application
 */
export function Show() {
    void call("Show");
}


/**
 * Quit the application
 */
export function Quit() {
    void call("Quit");
}