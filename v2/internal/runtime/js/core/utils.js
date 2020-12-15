/*
 _       __      _ __    
| |     / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  ) 
|__/|__/\__,_/_/_/____/  
The lightweight framework for web-like apps
(c) Lea Anthony 2019-present
*/
/* jshint esversion: 6 */

import { Emit } from './events';

export function AddScript(js, callbackID) {
	var script = document.createElement('script');
	script.text = js;
	document.body.appendChild(script);
	if (callbackID) {
		Emit(callbackID);
	}
}

// Adapted from webview - thanks zserge!
export function InjectCSS(css) {
	try {
		var elem = document.createElement('style');
		elem.setAttribute('type', 'text/css');
		if (elem.styleSheet) {
			elem.styleSheet.cssText = css;
		} else {
			elem.appendChild(document.createTextNode(css));
		}
		var head = document.head || document.getElementsByTagName('head')[0];
		head.appendChild(elem);
	} catch (e) {
		console.log(e);
	}
}

export function DisableDefaultContextMenu() {
	window.disableWailsDefaultContextMenu = true;
}
