/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The lightweight framework for web-like apps
(c) Lea Anthony 2019-present
*/
/* jshint esversion: 9 */

/**
 * Initialises platform specific code
 */

export const System = {
	AppType: 'desktop',
	Platform: () => 'windows',
};

export function SendMessage(message) {
	window.wailsInvoke(message);
}

export function Init() {

	// Setup drag handler
	// Based on code from: https://github.com/patr0nus/DeskGap
	window.addEventListener('mousedown', function (e) {
		let currentElement = e.target;
		while (currentElement != null) {
			if (currentElement.hasAttribute('data-wails-no-drag')) {
				break;
			} else if (currentElement.hasAttribute('data-wails-drag')) {
				window.wailsInvoke('wails-drag');
				break;
			}
			currentElement = currentElement.parentElement;
		}
	});

	// Setup context menu hook
	window.addEventListener('contextmenu', function (e) {
		let currentElement = e.target;
		let contextMenuId;
		while (currentElement != null) {
			contextMenuId = currentElement.dataset['wails-context-menu-id'];
			if (contextMenuId != null) {
				break;
			}
			currentElement = currentElement.parentElement;
		}
		if (contextMenuId != null || window.disableWailsDefaultContextMenu) {
			e.preventDefault();
		}
		if( contextMenuId != null ) {
			let contextData = currentElement.dataset['wails-context-menu-data'];
			let message = {
				id: contextMenuId,
				data: contextData || '',
			};
			window.wailsInvoke('C'+JSON.stringify(message));
		}
	});
}