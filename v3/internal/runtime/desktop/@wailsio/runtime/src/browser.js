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

const call = newRuntimeCallerWithID(objectNames.Browser, '');
const BrowserOpenURL = 0;

/**
 * Open a browser window to the given URL
 * @param {string} url - The URL to open
 * @returns {Promise<string>}
 */
export function OpenURL(url) {
    return call(BrowserOpenURL, {url});
}
