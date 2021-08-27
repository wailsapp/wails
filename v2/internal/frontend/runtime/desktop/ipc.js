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
 * SendMessage sends the given message to the backend
 *
 * @param {string} message
 */
window.WailsInvoke = function (message) {

	// Call Platform specific invoke method
	if (PLATFORM === "windows") {
		window.chrome.webview.postMessage(message);
	} else if (PLATFORM === "darwin") {
		window.blah();
	} else {
		console.error("Unsupported Platform");
	}
};
