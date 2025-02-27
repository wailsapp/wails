/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/
/**
 * Retrieves the value associated with the specified key from the flag map.
 *
 * @param key - The key to retrieve the value for.
 * @return The value associated with the specified key.
 */
export function GetFlag(key) {
    try {
        return window._wails.flags[key];
    }
    catch (e) {
        throw new Error("Unable to retrieve flag '" + key + "': " + e, { cause: e });
    }
}
