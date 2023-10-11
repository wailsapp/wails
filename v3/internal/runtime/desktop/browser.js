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

let call = newRuntimeCallerWithID(objectNames.Browser);

let BrowserOpenURL = 0;

/**
 * Open a browser window to the given URL
 * @param {string} url - The URL to open
 */
export function OpenURL(url) {
    void call(BrowserOpenURL, {url});
}
