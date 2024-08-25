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

import {EventsOn, EventsOff} from "./events";

const flags = {
    registered: false,
    defaultUseDropTarget: true,
    useDropTarget: true,
    nextDeactivate: null,
    nextDeactivateTimeout: null,
};

const DROP_TARGET_ACTIVE = "wails-drop-target-active";

/**
 * checkStyleDropTarget checks if the style has the drop target attribute
 * 
 * @param {CSSStyleDeclaration} style 
 * @returns 
 */
function checkStyleDropTarget(style) {
    const cssDropValue = style.getPropertyValue(window.wails.flags.cssDropProperty).trim();
    if (cssDropValue) {
        if (cssDropValue === window.wails.flags.cssDropValue) {
            return true;
        }
        // if the element has the drop target attribute, but 
        // the value is not correct, terminate finding process.
        // This can be useful to block some child elements from being drop targets.
        return false;
    }
    return false;
}

/**
 * onDragOver is called when the dragover event is emitted.
 * @param {DragEvent} e 
 * @returns 
 */
function onDragOver(e) {
    if (!window.wails.flags.enableWailsDragAndDrop) {
        return;
    }
    e.dataTransfer.dropEffect = 'copy';
    e.preventDefault();

    if (!flags.useDropTarget) {
        return;
    }

    const element = e.target;

    // Trigger debounce function to deactivate drop targets
    if(flags.nextDeactivate) flags.nextDeactivate();

    // if the element is null or element is not child of drop target element
    if (!element || !checkStyleDropTarget(getComputedStyle(element))) {
        return;
    }

    let currentElement = element;
    while (currentElement) {
        // check if currentElement is drop target element
        if (checkStyleDropTarget(currentElement.style)) {
            currentElement.classList.add(DROP_TARGET_ACTIVE);
        }
        currentElement = currentElement.parentElement;
    }
}

/**
 * onDragLeave is called when the dragleave event is emitted.
 * @param {DragEvent} e 
 * @returns 
 */
function onDragLeave(e) {
    if (!window.wails.flags.enableWailsDragAndDrop) {
        return;
    }
    e.preventDefault();

    if (!flags.useDropTarget) {
        return;
    }

    // Find the close drop target element
    if (!e.target || !checkStyleDropTarget(getComputedStyle(e.target))) {
        return null;
    }

    // Trigger debounce function to deactivate drop targets
    if(flags.nextDeactivate) flags.nextDeactivate();
    
    // Use debounce technique to tacle dragleave events on overlapping elements and drop target elements
    flags.nextDeactivate = () => {
        // Deactivate all drop targets, new drop target will be activated on next dragover event
        Array.from(document.getElementsByClassName(DROP_TARGET_ACTIVE)).forEach(el => el.classList.remove(DROP_TARGET_ACTIVE));
        // Reset nextDeactivate
        flags.nextDeactivate = null;
        // Clear timeout
        if (flags.nextDeactivateTimeout) {
            clearTimeout(flags.nextDeactivateTimeout);
            flags.nextDeactivateTimeout = null;
        }
    }

    // Set timeout to deactivate drop targets if not triggered by next drag event
    flags.nextDeactivateTimeout = setTimeout(() => {
        if(flags.nextDeactivate) flags.nextDeactivate();
    }, 50);
}

/**
 * onDrop is called when the drop event is emitted.
 * @param {DragEvent} e 
 * @returns 
 */
function onDrop(e) {
    if (!window.wails.flags.enableWailsDragAndDrop) {
        return;
    }
    e.preventDefault();

    if (CanResolveFilePaths()) {
        // process files
        let files = [];
        if (e.dataTransfer.items) {
            files = [...e.dataTransfer.items].map((item, i) => {
                if (item.kind === 'file') {
                    return item.getAsFile();
                }
            });
        } else {
            files = [...e.dataTransfer.files];
        }
        window.runtime.ResolveFilePaths(e.x, e.y, files);
    }

    if (!flags.useDropTarget) {
        return;
    }

    // Trigger debounce function to deactivate drop targets
    if(flags.nextDeactivate) flags.nextDeactivate();

    // Deactivate all drop targets
    Array.from(document.getElementsByClassName(DROP_TARGET_ACTIVE)).forEach(el => el.classList.remove(DROP_TARGET_ACTIVE));
}

/**
 * postMessageWithAdditionalObjects checks the browser's capability of sending postMessageWithAdditionalObjects
 *
 * @returns {boolean}
 * @constructor
 */
export function CanResolveFilePaths() {
    return window.chrome?.webview?.postMessageWithAdditionalObjects != null;
}

/**
 * ResolveFilePaths sends drop events to the GO side to resolve file paths on windows.
 *
 * @param {number} x
 * @param {number} y
 * @param {any[]} files
 * @constructor
 */
export function ResolveFilePaths(x, y, files) {
    // Only for windows webview2 >= 1.0.1774.30
    // https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2webmessagereceivedeventargs2?view=webview2-1.0.1823.32#applies-to
    if (window.chrome?.webview?.postMessageWithAdditionalObjects) {
        chrome.webview.postMessageWithAdditionalObjects(`file:drop:${x}:${y}`, files);
    }
}

/**
 * Callback for OnFileDrop returns a slice of file path strings when a drop is finished.
 *
 * @export
 * @callback OnFileDropCallback
 * @param {number} x - x coordinate of the drop
 * @param {number} y - y coordinate of the drop
 * @param {string[]} paths - A list of file paths.
 */

/**
 * OnFileDrop listens to drag and drop events and calls the callback with the coordinates of the drop and an array of path strings.
 *
 * @export
 * @param {OnFileDropCallback} callback - Callback for OnFileDrop returns a slice of file path strings when a drop is finished.
 * @param {boolean} [useDropTarget=true] - Only call the callback when the drop finished on an element that has the drop target style. (--wails-drop-target)
 */
export function OnFileDrop(callback, useDropTarget) {
    if (typeof callback !== "function") {
        console.error("DragAndDropCallback is not a function");
        return;
    }

    if (flags.registered) {
        return;
    }
    flags.registered = true;

    const uDTPT = typeof useDropTarget;
    flags.useDropTarget = uDTPT === "undefined" || uDTPT !== "boolean" ? flags.defaultUseDropTarget : useDropTarget;
    window.addEventListener('dragover', onDragOver);
    window.addEventListener('dragleave', onDragLeave);
    window.addEventListener('drop', onDrop);

    let cb = callback;
    if (flags.useDropTarget) {
        cb = function (x, y, paths) {
            const element = document.elementFromPoint(x, y)
            // if the element is null or element is not child of drop target element, return null
            if (!element || !checkStyleDropTarget(getComputedStyle(element))) {
                return null;
            }
            callback(x, y, paths);
        }
    }

    EventsOn("wails:file-drop", cb);
}

/**
 * OnFileDropOff removes the drag and drop listeners and handlers.
 */
export function OnFileDropOff() {
    window.removeEventListener('dragover', onDragOver);
    window.removeEventListener('dragleave', onDragLeave);
    window.removeEventListener('drop', onDrop);
    EventsOff("wails:file-drop");
    flags.registered = false;
}
