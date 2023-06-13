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

// defined in the Taskfile
export let invoke = function(input) {
    if(WINDOWS) {
        chrome.webview.postMessage(input);
    } else {
        webkit.messageHandlers.external.postMessage(input);
    }
}
