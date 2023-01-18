/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/
/* jshint esversion: 9 */

import "./ipc.js";
import {Callback, callbacks} from './calls';
import {EventsNotify, eventListeners} from "./events";
import {SetBindings} from "./bindings";

import * as Window from "./window";
import * as Screen from "./screen";
import * as Browser from "./browser";
import * as Log from './log';

let windowID = -1;


export function Quit() {
    window.WailsInvoke('Q');
}

export function Show() {
    window.WailsInvoke('S');
}

export function Hide() {
    window.WailsInvoke('H');
}

// export function Environment() {
//     return Call(":wails:Environment");
// }

// Internal wails endpoints
window.wails = {
    Callback,
    callbacks,
    EventsNotify,
    eventListeners,
    SetBindings,
    window: {
        ID: () => {
            return windowID
        },
    }
};

window.runtime = {
    ...Log,
    ...Window,
    ...Browser,
    ...Screen,
    EventsOn,
    EventsOnce,
    EventsOnMultiple,
    EventsEmit,
    EventsOff,
    // Environment,
    Show,
    Hide,
    Quit,
}

// Process the expected runtime config from the backend
if( window.wails_config ) {
    windowID = window.wails_config.windowID;
    window.wails_config = null;
}

if (DEBUG) {
    console.log("Wails v3.0.0 Debug Mode Enabled");
}

