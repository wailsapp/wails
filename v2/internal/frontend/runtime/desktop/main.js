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
import { eventListeners, EventsEmit, EventsNotify, EventsOff, EventsOn, EventsOnce, EventsOnMultiple } from './events';
import { Callback, callbacks } from './calls';
import { SetBindings } from "./bindings";
import * as Window from "./window";
import * as Browser from "./browser";


export function Quit() {
    window.WailsInvoke('Q');
}

// The JS runtime
window.runtime = {
    ...Log,
    ...Window,
    ...Browser,
    EventsOn,
    EventsOnce,
    EventsOnMultiple,
    EventsEmit,
    EventsOff,
    Quit
};

// Internal wails endpoints
window.wails = {
    Callback,
    EventsNotify,
    SetBindings,
    eventListeners,
    callbacks
};

// Set the bindings
window.wails.SetBindings(window.wailsbindings);
delete window.wails.SetBindings;

// This is evaluated at build time in package.json
if (ENV === "production") {
    delete window.wailsbindings;
}

// Setup drag handler
// Based on code from: https://github.com/patr0nus/DeskGap
window.addEventListener('mousedown', (e) => {
    let currentElement = e.target;
    while (currentElement != null) {
        if (currentElement.hasAttribute('data-wails-no-drag')) {
            break;
        } else if (currentElement.hasAttribute('data-wails-drag')) {
            window.WailsInvoke("drag");
            break;
        }
        currentElement = currentElement.parentElement;
    }
});
