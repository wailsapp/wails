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

import {invoke} from "./ipc.js";
import {Callback, callbacks} from './calls';
import {EventsNotify, eventListeners} from "./events";
import {SetBindings} from "./bindings";

import {newWindow} from "./window";

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
};


export function newRuntime(id) {
    return {
        // Log: newLog(id),
        // Browser: newBrowser(id),
        // Screen: newScreen(id),
        // Events: newEvents(id),
        Window: newWindow(id),
        Show: () => invoke("S"),
        Hide: () => invoke("H"),
        Quit: () => invoke("Q"),
        // GetWindow: function (windowID) {
        //     if (!windowID) {
        //         return this.Window;
        //     }
        //     return newWindow(windowID);
        // }
    }
}

window.runtime = newRuntime(-1);

if (DEBUG) {
    console.log("Wails v3.0.0 Debug Mode Enabled");
}

