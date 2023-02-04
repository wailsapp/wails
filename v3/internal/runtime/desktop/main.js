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

import {dialogCallback, dialogErrorCallback, Error, Info, OpenFile, Question, SaveFile, Warning,} from "./dialogs";

import * as Clipboard from './clipboard';
import * as Application from './application';

import {newWindow} from "./window";
import {dispatchCustomEvent, Emit, Off, OffAll, On, Once, OnMultiple} from "./events";

// Internal wails endpoints
window.wails = {
    ...newRuntime(-1),
};

window._wails = {
    dialogCallback,
    dialogErrorCallback,
    dispatchCustomEvent,
}


export function newRuntime(id) {
    return {
        Clipboard: {
            ...Clipboard
        },
        Application: {
            ...Application
        },
        Dialog: {
            Info,
            Warning,
            Error,
            Question,
            OpenFile,
            SaveFile,
        },
        Events: {
            Emit,
            On,
            Once,
            OnMultiple,
            Off,
            OffAll,
        },
        Window: newWindow(id),
    }
}

if (DEBUG) {
    console.log("Wails v3.0.0 Debug Mode Enabled");
}

