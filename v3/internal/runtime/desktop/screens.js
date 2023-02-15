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

let call = newRuntimeCaller("screens");

export function GetAll() {
    return call("GetAll");
}

export function GetPrimary() {
    return call("GetPrimary");
}

export function GetCurrent() {
    return call("GetCurrent");
}