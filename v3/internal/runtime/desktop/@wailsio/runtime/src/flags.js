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
 * Retrieves the value associated with the specified key from the flag map.
 *
 * @param {string} keyString - The key to retrieve the value for.
 * @return {*} - The value associated with the specified key.
 */
export function GetFlag(keyString) {
    try {
        return window._wails.flags[keyString];
    } catch (e) {
        throw new Error("Unable to retrieve flag '" + keyString + "': " + e);
    }
}
