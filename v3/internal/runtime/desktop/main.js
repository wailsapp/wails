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


import {Info, Warning, Error, Question, OpenFile, SaveFile, dialogCallback, dialogErrorCallback, } from "./dialogs";

import * as Clipboard from './clipboard';
import {newWindow} from "./window";

// export function Environment() {
//     return Call(":wails:Environment");
// }

// Internal wails endpoints
window.wails = {
    ...newRuntime(-1),
};

window._wails = {
    dialogCallback,
    dialogErrorCallback,
}


export function newRuntime(id) {
    return {
        // Log: newLog(id),
        // Browser: newBrowser(id),
        // Screen: newScreen(id),
        // Events: newEvents(id),
        Clipboard: {
            ...Clipboard
        },
        Dialog: {
            Info,
            Warning,
            Error,
            Question,
            OpenFile,
            SaveFile,
        },
        Window: newWindow(id),
        Application: {
            Show: () => invoke("S"),
            Hide: () => invoke("H"),
            Quit: () => invoke("Q"),
        }
        // GetWindow: function (windowID) {
        //     if (!windowID) {
        //         return this.Window;
        //     }
        //     return newWindow(windowID);
        // }
    }
}

if (DEBUG) {
    console.log("Wails v3.0.0 Debug Mode Enabled");
}

