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

import {log} from "./log";

let websocket = null;
let connectTimer;

window.onbeforeunload = function () {
    if (websocket) {
        websocket.onclose = function () {
        };
        websocket.close();
        websocket = null;
    }
};

// ...and attempt to connect
connect();


function setupIPCBridge() {
    window.WailsInvoke = (message) => {
        websocket.send(message);
    };
}

// Handles incoming websocket connections
function handleConnect() {
    log('Connected to backend');
    setupIPCBridge();
    clearInterval(connectTimer);
    websocket.onclose = handleDisconnect;
    websocket.onmessage = handleMessage;
}

// Handles websocket disconnects
function handleDisconnect() {
    log('Disconnected from backend');
    websocket = null;
    connect();
}

// Try to connect to the backend every 1s (default value).
function connect() {
    connectTimer = setInterval(function () {
        if (websocket == null) {
            websocket = new WebSocket('ws://' + window.location.hostname + ':34115/wails/ipc');
            websocket.onopen = handleConnect;
            websocket.onerror = function (e) {
                e.stopImmediatePropagation();
                e.stopPropagation();
                e.preventDefault();
                websocket = null;
                return false;
            };
        }
    }, 1000);
}

function handleMessage(message) {

    // ignore
    if (message === "drag") {
        return;
    }

    // As a bridge we ignore js and css injections
    switch (message.data[0]) {
        // Wails library - inject!
        // case 'b':
        // 	message = message.data.slice(1);
        // 	addScript(message);
        // 	log('Loaded Wails Runtime');
        //
        // 	// We need to now send a message to the backend telling it
        // 	// we have loaded (System Start)
        // 	window.wailsInvoke('SS');
        //
        // 	// Now wails runtime is loaded, wails for the ready event
        // 	// and callback to the main app
        // 	// window.wails.Events.On('wails:loaded', function () {
        // 	if (callback) {
        // 		log('Notifying application');
        // 		callback(window.wails);
        // 	}
        // 	// });
        // 	break;
        // Notifications
        case 'n':
            window.wails.EventsNotify(message.data.slice(1));
            break;
        // 	// 	// Binding
        // 	// case 'b':
        // 	// 	const binding = message.data.slice(1);
        // 	// 	//log("Binding: " + binding)
        // 	// 	window.wails._.NewBinding(binding);
        // 	// 	break;
        // 	// 	// Call back
        case 'c':
            const callbackData = message.data.slice(1);
            window.wails.Callback(callbackData);
            break;

        // Application menu
        /*
        case 'M':
            const menuData = message.data.slice(1);
            window.appmenu = JSON.parse(menuData);
            ApplicationMenu.set(window.appmenu);
            console.log("Set application menu");
            break;
            */
        case 'T':
            const serverStartTime = message.data.slice(1);
            const lastServerStartTime = localStorage.getItem('wails:lastServerStartTime');
            if (lastServerStartTime === null) {
                localStorage.setItem('wails:lastServerStartTime', serverStartTime);
            }
            if (serverStartTime !== lastServerStartTime) {
                localStorage.setItem('wails:lastServerStartTime', serverStartTime);
                location.reload(true);
                console.log("Reloaded");
            }
            break;
        // 	// Tray
        // case 'T':
        // 	const trayMessage = message.data.slice(1);
        // 	switch (trayMessage[0]) {
        // 	case 'S':
        // 		// Set tray
        // 		const trayJSON = trayMessage.slice(1);
        // 		let tray = JSON.parse(trayJSON);
        // 		setTray(tray);
        // 		break;
        // 	case 'U':
        // 		// Update label
        // 		const updateTrayLabelJSON = trayMessage.slice(1);
        // 		let trayLabelData = JSON.parse(updateTrayLabelJSON);
        // 		updateTrayLabel(trayLabelData);
        // 		break;
        // 	case 'D':
        // 		// Delete Tray Menu
        // 		const id = trayMessage.slice(1);
        // 		deleteTrayMenu(id);
        // 		break;
        // 	default:
        // 		log('Unknown tray message: ' + message.data);
        // 	}
        // 	break;

        default:
            log('Unknown message: ' + message.data);
    }
}
