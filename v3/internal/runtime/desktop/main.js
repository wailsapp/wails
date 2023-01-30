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

import {Info, Warning, Error, Question, OpenFile, SaveFile, dialogCallback, dialogErrorCallback, } from "./dialogs";

import * as Clipboard from './clipboard';
import {newWindow} from "./window";

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
    }
}

if (DEBUG) {
    console.log("Wails v3.0.0 Debug Mode Enabled");
}

