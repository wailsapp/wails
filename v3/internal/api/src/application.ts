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

let call : (method: string, args?: any) => Promise<void> = newRuntimeCaller("application");

export function Hide() {
    return call("Hide");
}

export function Show() {
    return call("Show");
}

export function Quit() {
    return call("Quit");
}