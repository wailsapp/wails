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

import {EventsOn} from "./events";

/**
 * postMessageWithAdditionalObjects checks the browser's capability of sending postMessageWithAdditionalObjects
 *
 * @returns {boolean}
 * @constructor
 */
export function CanResolveFilePaths() {
    return window.chrome?.webview?.postMessageWithAdditionalObjects != null;
}

export function ResolveFilePaths(files) {
    // Only for windows webview2 >= 1.0.1774.30
    // https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2webmessagereceivedeventargs2?view=webview2-1.0.1823.32#applies-to
    if (!window.chrome?.webview?.postMessageWithAdditionalObjects) {
        reject(new Error("Unsupported Platform"));
        return;
    }

    chrome.webview.postMessageWithAdditionalObjects(`file:drop:`, files);
}
