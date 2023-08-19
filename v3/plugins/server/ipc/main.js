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
import { nanoid } from 'nanoid/non-secure';

let components = {};
let source = null;

function handleCallback(e) {
    const payload = JSON.parse(e.data);
    _wails.callCallback(payload.id,
                        payload.result,
                        true);
}

function handleCallbackError(e) {
    const payload = JSON.parse(e.data);
    _wails.callErrorCallback(payload.id, payload.result);
}

function handleDialog(e) {
    const payload = JSON.parse(e.data);
    _wails.dialogCallback(payload.id,
                          payload.result,
                          true);
}

function handleDialogError(e) {
    const payload = JSON.parse(e.data);
    _wails.dialogErrorCallback(payload.id, payload.result);
}

function handleWailsEvent(e) {
    console.log("WailsEvent: " + e.data)
}

window.addEventListener('DOMContentLoaded', () => {
    components.overlay = new Overlay({
        target: document.body,
        anchor: document.querySelector('#wails-spinner'),
    });
    connect();
});

window.onbeforeunload = function () {
    if (source) {
        source.onclose = function () { };
        source.close();
        source = null;
    }
};

// Handles sse connections
function handleConnect(e) {
    hideOverlay();
    source.onclose = handleDisconnect;

}

// Handles SSE disconnects
// EventSource will attempt to reconnect on it's own
function handleDisconnect(e) {
    if (this.readyState == EventSource.CONNECTING) {
        showOverlay();
    } else {
        console.log(e);
    }
}

function _connect() {
    if (source == null) {
        source = new EventSource("/server/events?clientId="+wails.clientId);
        source.onopen = handleConnect;
        source.onerror = handleDisconnect;
        source.addEventListener('cb', handleCallback);
        source.addEventListener('cberror', handleCallbackError);
        source.addEventListener('dlgcb', handleDialog);
        source.addEventListener('dlgcberror', handleDialogError);
        source.addEventListener('wailsevent', handleWailsEvent);
    }
}

// Try to connect to the backend every .5s
function connect() {
    _connect();
}
