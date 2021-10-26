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
	// Credit: https://stackoverflow.com/a/2631521
	let _deeptest = function (s) {
		var obj = window[s.shift()];
		while (obj && s.length) obj = obj[s.shift()];
		return obj;
	};
	let windows = _deeptest(["chrome", "webview", "postMessage"]);
	let mac = _deeptest(["webkit", "messageHandlers", "external", "postMessage"]);

	if (!windows && !mac) {
		console.error("Unsupported Platform");
		return;
	}

	if (windows) {
		window.WailsInvoke = (message) => window.chrome.webview.postMessage(message);
	}
	if (mac) {
		window.WailsInvoke = (message) => window.webkit.messageHandlers.external.postMessage(message);
	}
})();