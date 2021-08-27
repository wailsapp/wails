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

import {log} from "./log";
import Overlay from "./Overlay.svelte";
import {hideOverlay, showOverlay} from "./store";

let components = {};

// Sets up the overlay
components.overlay = new Overlay({
    target: document.body,
    anchor: document.querySelector('#wails-spinner'),
});

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
    hideOverlay();
    setupIPCBridge();
    clearInterval(connectTimer);
    websocket.onclose = handleDisconnect;
    websocket.onmessage = handleMessage;
}

// Handles websocket disconnects
function handleDisconnect() {
    log('Disconnected from backend');
    websocket = null;
    showOverlay();
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
    }, 250);
}

function handleMessage(message) {

    if (message.data === "reload") {
        window.runtime.WindowReload();
        return;
    }

    // As a bridge we ignore js and css injections
    switch (message.data[0]) {
        // Notifications
        case 'n':
            window.wails.EventsNotify(message.data.slice(1));
            break;
        case 'c':
            const callbackData = message.data.slice(1);
            window.wails.Callback(callbackData);
            break;
        default:
            log('Unknown message: ' + message.data);
    }
}
