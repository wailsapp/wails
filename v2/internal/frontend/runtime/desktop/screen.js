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


import {Call} from "./calls";


/**
 * Gets the all screens. Call this anew each time you want to refresh data from the underlying windowing system.
 * @export
 * @typedef {import('../wrapper/runtime').Screen} Screen
 * @return {Promise<{Screen[]}>} The screens
 */
export function ScreenGetAll() {
    return Call(":wails:ScreenGetAll");
}
