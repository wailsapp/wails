/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The lightweight framework for web-like apps
(c) Lea Anthony 2019-present
*/
/* jshint esversion: 6 */

/**
 * Initialises platform specific code
 */

export const System = {
    Platform: "darwin",
    AppType: "desktop"
}

export function SendMessage(message) {
    window.webkit.messageHandlers.external.postMessage(message);
}

export function Init() {

    // Setup drag handler
    // Based on code from: https://github.com/patr0nus/DeskGap
    window.addEventListener('mousedown', function (e) {
        var currentElement = e.target;
        while (currentElement != null) {
            if (currentElement.hasAttribute('data-wails-no-drag')) {
                break;
            } else if (currentElement.hasAttribute('data-wails-drag')) {
                window.webkit.messageHandlers.windowDrag.postMessage(null);
                break;
            }
            currentElement = currentElement.parentElement;
        }
    });
}