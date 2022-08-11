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
import * as Log from './log';
import {eventListeners, EventsEmit, EventsNotify, EventsOff, EventsOn, EventsOnce, EventsOnMultiple} from './events';
import {Call, Callback, callbacks} from './calls';
import {SetBindings} from "./bindings";
import * as Window from "./window";
import * as Screen from "./screen";
import * as Browser from "./browser";


export function Quit() {
    window.WailsInvoke('Q');
}

export function Show() {
    window.WailsInvoke('S');
}

export function Hide() {
    window.WailsInvoke('H');
}

export function Environment() {
    return Call(":wails:Environment");
}

// The JS runtime
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
    Environment,
    Show,
    Hide,
    Quit
};

// Internal wails endpoints
window.wails = {
    Callback,
    EventsNotify,
    SetBindings,
    eventListeners,
    callbacks,
    flags: {
        disableScrollbarDrag: false,
        disableWailsDefaultContextMenu: false,
        enableResize: false,
        defaultCursor: null,
        borderThickness: 6,
        shouldDrag: false
    }
};

// Set the bindings
window.wails.SetBindings(window.wailsbindings);
delete window.wails.SetBindings;

// This is evaluated at build time in package.json
// const dev = 0;
// const production = 1;
if (ENV === 0) {
    delete window.wailsbindings;
}

window.addEventListener('mouseup', () => {
    window.wails.flags.shouldDrag = false;
});

function setResize(cursor) {
    document.body.style.cursor = cursor || window.wails.flags.defaultCursor;
    window.wails.flags.resizeEdge = cursor;
}

window.addEventListener('mousedown', function (e) {
   if (e.target.hasAttribute('data-wails-drag') && e.buttons === 1) {
    e.preventDefault();
    window.WailsInvoke("drag");
   }
});

// Setup context menu hook
window.addEventListener('contextmenu', function (e) {
    if (window.wails.flags.disableWailsDefaultContextMenu) {
        e.preventDefault();
    }
});
