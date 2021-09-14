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
// const linux = 1;
// const macos = 2;

window.WailsInvoke = function (message) {

	// Call Platform specific invoke method
	if (PLATFORM === 0) {
		window.chrome.webview.postMessage(message);
	} else if (PLATFORM === 1) {
		window.blah();
	} else if (PLATFORM === 2) {
		console.error("Unsupported Platform");
	}
};
