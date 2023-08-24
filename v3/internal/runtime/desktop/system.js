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

let call = newRuntimeCallerWithID(objectNames.System);

let SystemIsDarkMode = 0;

/**
 * Determines if the system is currently using dark mode
 * @returns {Promise<boolean>}
 */
export function IsDarkMode() {
    return call(SystemIsDarkMode);
}