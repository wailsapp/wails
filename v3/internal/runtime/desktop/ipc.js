/*
 _       __      _ __
| |     / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/
/* jshint esversion: 9 */



let postMessage = null;

(function () {
	// Credit: https://stackoverflow.com/a/2631521
	let _deeptest = function (s) {
		let obj = window[s.shift()];
		while (obj && s.length) obj = obj[s.shift()];
		return obj;
	};
	let windows = _deeptest(["chrome", "webview", "postMessage"]);
	let mac_linux = _deeptest(["webkit", "messageHandlers", "external", "postMessage"]);

	if (!windows && !mac_linux) {
		console.error("Unsupported Platform");
		return;
	}

	if (windows) {
		postMessage = (message) => window.chrome.webview.postMessage(message);
	}
	if (mac_linux) {
		postMessage = (message) => window.webkit.messageHandlers.external.postMessage(message);
	}
})();

export function invoke(message, id) {
	if( id && id !== -1) {
		postMessage("WINDOWID:"+ id + ":" + message);
	} else {
		postMessage(message);
	}
}
