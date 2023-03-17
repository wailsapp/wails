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


import * as Clipboard from './clipboard';
import * as Application from './application';
import * as Log from './log';
import * as Screens from './screens';
import {Plugin, Call, callErrorCallback, callCallback} from "./calls";
import {newWindow} from "./window";
import {dispatchCustomEvent, Emit, Off, OffAll, On, Once, OnMultiple} from "./events";
import {dialogCallback, dialogErrorCallback, Error, Info, OpenFile, Question, SaveFile, Warning,} from "./dialogs";
import {enableContextMenus} from "./contextmenu";
import {reloadWML} from "./wml";

window.wails = {
    ...newRuntime(-1),
};

// Internal wails endpoints
window._wails = {
    dialogCallback,
    dialogErrorCallback,
    dispatchCustomEvent,
    callCallback,
    callErrorCallback,
};

export function newRuntime(id) {
    return {
        Clipboard: {
            ...Clipboard
        },
        Application: {
            ...Application
        },
        Log,
        Screens,
        Call,
        Plugin,
        WML: {
            Reload: reloadWML,
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
    };
}

if (DEBUG) {
    console.log("Wails v3.0.0 Debug Mode Enabled");
}

enableContextMenus(true);

document.addEventListener("DOMContentLoaded", function(event) {
    reloadWML();
});