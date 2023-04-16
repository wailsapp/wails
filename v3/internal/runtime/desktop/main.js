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
import {dispatchWailsEvent, Emit, Off, OffAll, On, Once, OnMultiple} from "./events";
import {dialogCallback, dialogErrorCallback, Error, Info, OpenFile, Question, SaveFile, Warning,} from "./dialogs";
import {enableContextMenus} from "./contextmenu";
import {reloadWML} from "./wml";

window.wails = {
    ...newRuntime(null),
};

// Internal wails endpoints
window._wails = {
    dialogCallback,
    dialogErrorCallback,
    dispatchWailsEvent,
    callCallback,
    callErrorCallback,
};

export function newRuntime(windowName) {
    return {
        Clipboard: {
            ...Clipboard
        },
        Application: {
            ...Application,
            GetWindowByName(windowName) {
                return newRuntime(windowName);
            }
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
        Window: newWindow(windowName),
    };
}

if (DEBUG) {
    console.log("Wails v3.0.0 Debug Mode Enabled");
}

enableContextMenus(true);

document.addEventListener("DOMContentLoaded", function(event) {
    reloadWML();
});