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

// const windows = 0;
// const macos = 1;
// const linux = 2;

window.WailsInvoke = function (message) {

	// Call Platform specific invoke method
	if (PLATFORM === 0) {
		window.chrome.webview.postMessage(message);
	} else if (PLATFORM === 1) {
		window.webkit.messageHandlers.external.postMessage(message);
	} else if (PLATFORM === 2) {
		console.error("Unsupported Platform");
	} else {
		console.error("Unsupported Platform");
	}
};
