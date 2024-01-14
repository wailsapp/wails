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
import * as DragAndDrop from "./draganddrop";
import * as ContextMenu from "./contextmenu";


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
    ...DragAndDrop,
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
        disableDefaultContextMenu: false,
        enableResize: false,
        defaultCursor: null,
        borderThickness: 6,
        shouldDrag: false,
        deferDragToMouseMove: true,
        cssDragProperty: "--wails-draggable",
        cssDragValue: "drag",
        cssDropProperty: "--wails-drop-target",
        cssDropValue: "drop",
        enableWailsDragAndDrop: false,
        wailsDropPreviousElement: null,
    }
};

// Set the bindings
if (window.wailsbindings) {
    window.wails.SetBindings(window.wailsbindings);
    delete window.wails.SetBindings;
}

// (bool) This is evaluated at build time in package.json
if (!DEBUG) {
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

window.wails.setCSSDropProperties = function (property, value) {
    window.wails.flags.cssDropProperty = property;
    window.wails.flags.cssDropValue = value;
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
    // always show the contextmenu in debug & dev
    if (DEBUG) return;

    if (window.wails.flags.disableDefaultContextMenu) {
        e.preventDefault();
    } else {
        ContextMenu.processDefaultContextMenu(e);
    }
});

window.addEventListener('dragover', function (e) {
    if (!window.wails.flags.enableWailsDragAndDrop) {
        return;
    }
    e.preventDefault();
    let targetElement = document.elementFromPoint(e.x, e.y);
    if (targetElement === window.wails.flags.wailsDropPreviousElement) {
        return;
    }
    const style = targetElement.style;
    let cssDropValue = null;
    if (Object.keys(style).findIndex(key => style[key] === window.wails.flags.cssDropProperty) < 0) {
        targetElement = targetElement.closest(`[style*='${window.wails.flags.cssDropProperty}']`);
    }
    if (targetElement == null) {
        return;
    }
    cssDropValue = window.getComputedStyle(targetElement).getPropertyValue(window.wails.flags.cssDropProperty);
    if (cssDropValue) {
        cssDropValue = cssDropValue.trim();
    }

    if (cssDropValue === window.wails.flags.cssDropValue) {
        targetElement.classList.add("wails-drop-target-active");
    } else if (window.wails.flags.wailsDropPreviousElement) {
        window.wails.flags.wailsDropPreviousElement.classList.remove("wails-drop-target-active");
    }
    window.wails.flags.wailsDropPreviousElement = targetElement;
})

window.addEventListener('dragleave', function (e) {
    if (!window.wails.flags.enableWailsDragAndDrop) {
        return;
    }
    e.preventDefault();

    let targetElement = document.elementFromPoint(e.x, e.y);
    let cssDropValue = window.getComputedStyle(targetElement).getPropertyValue(window.wails.flags.cssDropProperty);
    if (cssDropValue) {
        cssDropValue = cssDropValue.trim();
    }
    if (cssDropValue !== window.wails.flags.cssDropValue && window.wails.flags.wailsDropPreviousElement) {
        window.wails.flags.wailsDropPreviousElement.classList.remove("wails-drop-target-active");
    }
});

window.addEventListener('drop', function (e) {
    if (!window.wails.flags.enableWailsDragAndDrop) {
        return;
    }
    e.preventDefault();
    let targetElement = document.elementFromPoint(e.x, e.y);
    let cssDropValue = window.getComputedStyle(targetElement).getPropertyValue(window.wails.flags.cssDropProperty);
    if (cssDropValue) {
        cssDropValue = cssDropValue.trim();
    }
    if (cssDropValue !== window.wails.flags.cssDropValue) {
        return;
    }
    // process files
    let files = [];
    if (e.dataTransfer.items) {
        files = [...e.dataTransfer.items].map((item, i) => {
            if (item.kind === 'file') {
                const file = item.getAsFile();
                return file;
            }
        });
    } else {
        files = [...e.dataTransfer.files];
    }

    window.runtime.ResolveFilePaths(files);
    if(window.wails.flags.wailsDropPreviousElement) {
        window.wails.flags.wailsDropPreviousElement.classList.remove("wails-drop-target-active");
    }
});

window.WailsInvoke("runtime:ready");