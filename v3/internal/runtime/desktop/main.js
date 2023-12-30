/*
 _     __     _ __
| |  / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/
/* jshint esversion: 9 */


import * as Clipboard from './@wailsio/runtime/clipboard';
import * as Application from './@wailsio/runtime/application';
import * as Screens from './@wailsio/runtime/screens';
import * as System from './@wailsio/runtime/system';
import * as Browser from './@wailsio/runtime/browser';
import * as Window from './@wailsio/runtime/window';
import {Plugin, Call, errorHandler as callErrorHandler, resultHandler as callResultHandler, ByID, ByName} from "./@wailsio/runtime/calls";
import {clientId} from './@wailsio/runtime/runtime';
import {dispatchWailsEvent, Emit, Off, OffAll, On, Once, OnMultiple} from "./@wailsio/runtime/events";
import {dialogResultCallback, dialogErrorCallback, Error, Info, OpenFile, Question, SaveFile, Warning} from "./@wailsio/runtime/dialogs";
import {setupContextMenus} from './@wailsio/runtime/contextmenu';
import {reloadWML} from './@wailsio/runtime/wml';
import {setupDrag, endDrag, setResizable} from './@wailsio/runtime/drag';

window.wails = {
    ...newRuntime(null),
    clientId: clientId,
};

// Internal wails endpoints
window._wails = {
    dialogResultCallback,
    dialogErrorCallback,
    dispatchWailsEvent,
    callErrorHandler,
    callResultHandler,
    endDrag,
    setResizable,
};

export function newRuntime(windowName) {
    return {
        Clipboard: {
            ...Clipboard
        },
        Application: {
            ...Application,
        },
        System,
        Screens,
        Browser,
        Call: {
            Call,
            ByID,
            ByName,
            Plugin,
        },
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
        Window: {
            ...Window.Get('')
        },
    };
}

setupContextMenus();
setupDrag();

document.addEventListener("DOMContentLoaded", function() {
    reloadWML();
});
