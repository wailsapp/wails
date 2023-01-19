/*
 _       __      _ __
| |     / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/
/* jshint esversion: 6 */

/**
 * WailsInvoke sends the given message to the backend
 *
 * @param {string} message
 */

(function () {
	window.WailsInvoke = (message) => {
		WINDOWS && window.chrome.webview.postMessage(message);
		(DARWIN || LINUX) && window.webkit.messageHandlers.wails.postMessage(message);
	}
})();