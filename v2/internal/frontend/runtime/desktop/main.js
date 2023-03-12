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
import * as Clipboard from "./clipboard";


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
    ...Clipboard,
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
        shouldDrag: false,
        deferDragToMouseMove: false,
        cssDragProperty: "--wails-draggable",
        cssDragValue: "drag",
    }
};

// Set the bindings
if (window.wailsbindings) {
    window.wails.SetBindings(window.wailsbindings);
    delete window.wails.SetBindings;
}

// This is evaluated at build time in package.json
// const dev = 0;
// const production = 1;
if (ENV === 1) {
    delete window.wailsbindings;
}

let dragTest = function (e) {
    var val = window.getComputedStyle(e.target).getPropertyValue(window.wails.flags.cssDragProperty);
    if (val) {
      val = val.trim();
    }
    
    if (val !== window.wails.flags.cssDragValue) {
        return false;
    }

    if (e.buttons !== 1) {
        // Do not start dragging if not the primary button has been clicked.
        return false;
    }

    if (e.detail !== 1) {
        // Do not start dragging if more than once has been clicked, e.g. when double clicking
        return false;
    }

    return true;
};

window.wails.setCSSDragProperties = function (property, value) {
    window.wails.flags.cssDragProperty = property;
    window.wails.flags.cssDragValue = value;
}

window.addEventListener('mousedown', (e) => {

    // Check for resizing
    if (window.wails.flags.resizeEdge) {
        window.WailsInvoke("resize:" + window.wails.flags.resizeEdge);
        e.preventDefault();
        return;
    }

    if (dragTest(e)) {
        if (window.wails.flags.disableScrollbarDrag) {
            // This checks for clicks on the scroll bar
            if (e.offsetX > e.target.clientWidth || e.offsetY > e.target.clientHeight) {
                return;
            }
        }
        if (window.wails.flags.deferDragToMouseMove) {
            window.wails.flags.shouldDrag = true;
        } else {
            e.preventDefault()
            window.WailsInvoke("drag");
        }
        return;
    } else {
        window.wails.flags.shouldDrag = false;
    }
});

window.addEventListener('mouseup', () => {
    window.wails.flags.shouldDrag = false;
});

function setResize(cursor) {
    document.documentElement.style.cursor = cursor || window.wails.flags.defaultCursor;
    window.wails.flags.resizeEdge = cursor;
}

window.addEventListener('mousemove', function (e) {
    if (window.wails.flags.shouldDrag) {
        window.wails.flags.shouldDrag = false;
        let mousePressed = e.buttons !== undefined ? e.buttons : e.which;
        if (mousePressed > 0) {
            window.WailsInvoke("drag");
            return;
        }
    }
    if (!window.wails.flags.enableResize) {
        return;
    }
    if (window.wails.flags.defaultCursor == null) {
        window.wails.flags.defaultCursor = document.documentElement.style.cursor;
    }
    if (window.outerWidth - e.clientX < window.wails.flags.borderThickness && window.outerHeight - e.clientY < window.wails.flags.borderThickness) {
        document.documentElement.style.cursor = "se-resize";
    }
    let rightBorder = window.outerWidth - e.clientX < window.wails.flags.borderThickness;
    let leftBorder = e.clientX < window.wails.flags.borderThickness;
    let topBorder = e.clientY < window.wails.flags.borderThickness;
    let bottomBorder = window.outerHeight - e.clientY < window.wails.flags.borderThickness;

    // If we aren't on an edge, but were, reset the cursor to default
    if (!leftBorder && !rightBorder && !topBorder && !bottomBorder && window.wails.flags.resizeEdge !== undefined) {
        setResize();
    } else if (rightBorder && bottomBorder) setResize("se-resize");
    else if (leftBorder && bottomBorder) setResize("sw-resize");
    else if (leftBorder && topBorder) setResize("nw-resize");
    else if (topBorder && rightBorder) setResize("ne-resize");
    else if (leftBorder) setResize("w-resize");
    else if (topBorder) setResize("n-resize");
    else if (bottomBorder) setResize("s-resize");
    else if (rightBorder) setResize("e-resize");

});

// Setup context menu hook
window.addEventListener('contextmenu', function (e) {
    if (window.wails.flags.disableWailsDefaultContextMenu) {
        e.preventDefault();
    }
});

window.WailsInvoke("runtime:ready");