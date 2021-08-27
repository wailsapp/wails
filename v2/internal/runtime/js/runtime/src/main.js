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

import Overlay from './Overlay.svelte';
import MenuBar from './Menubar.svelte';
import {showOverlay} from "./store";
import {StartWebsocket} from "./websocket";

let components = {};

function setupMenuBar() {
	components.menubar = new MenuBar({
		target: document.body,
	});
}

// Sets up the overlay
function setupOverlay() {
	components.overlay = new Overlay({
		target: document.body,
		anchor: document.querySelector('#wails-bridge'),
	});
}

export function InitBridge(callback) {

	setupMenuBar()

	// Setup the overlay
	setupOverlay();

	// Start by showing the overlay...
	showOverlay();

	// ...and attempt to connect
	StartWebsocket(callback);
}
