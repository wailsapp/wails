/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The lightweight framework for web-like apps
(c) Lea Anthony 2019-present
*/
/* jshint esversion: 9 */
import * as Log from './log';
import {eventListeners, EventsEmit, EventsNotify, EventsOff, EventsOn, EventsOnce, EventsOnMultiple} from './events';
import {Callback, callbacks} from './calls';
import {SetBindings} from "./bindings";

// Backend is where the Go struct wrappers get bound to
window.backend = {};

window.runtime = {
    ...Log,
    EventsOn,
    EventsOnce,
    EventsOnMultiple,
    EventsEmit,
    EventsOff,
};

// Initialise global if not already
window.wails = {
    Callback,
    EventsNotify,
    SetBindings,
    eventListeners,
    callbacks
};

window.wails.SetBindings(window.wailsbindings);
delete window.wails.SetBindings;
delete window.wailsbindings;

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
