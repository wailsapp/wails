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

import * as events from "./events";

// Setup
let dndEnabled = false;
let initialised = false;
let nextDeactivate = null;
let nextDeactivateTimeout = null;

window._wails = window._wails || {};

window._wails.enableDragAndDrop = function(value) {
    dndEnabled = value;
}

const DROP_TARGET_ACTIVE = "wails-drop-target-active";
const CSS_DROP_PROPERTY = "--wails-drop-target";

/**
 * checkStyleDropTarget checks if the style has the drop target attribute
 * 
 * @param {CSSStyleDeclaration} style 
 * @returns 
 */
function checkStyleDropTarget(style) {
    const cssDropValue = style.getPropertyValue(CSS_DROP_PROPERTY).trim();
    if (cssDropValue) {
        if (cssDropValue === 'drop' ) {
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
    if (!dndEnabled) {
        return;
    }
    e.preventDefault();

    const element = e.target;

    // Trigger debounce function to deactivate drop targets
    if(nextDeactivate) nextDeactivate();

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
    if (!dndEnabled) {
        return;
    }
    e.preventDefault();

    // Find the close drop target element
    if (!e.target || !checkStyleDropTarget(getComputedStyle(e.target))) {
        return null;
    }

    // Trigger debounce function to deactivate drop targets
    if(nextDeactivate) nextDeactivate();
    
    // Use debounce technique to tackle dragleave events on overlapping elements and drop target elements
    nextDeactivate = () => {
        // Deactivate all drop targets, new drop target will be activated on next dragover event
        Array.from(document.getElementsByClassName(DROP_TARGET_ACTIVE)).forEach(el => el.classList.remove(DROP_TARGET_ACTIVE));
        // Reset nextDeactivate
        nextDeactivate = null;
        // Clear timeout
        if (nextDeactivateTimeout) {
            clearTimeout(nextDeactivateTimeout);
            nextDeactivateTimeout = null;
        }
    }

    // Set timeout to deactivate drop targets if not triggered by next drag event
    nextDeactivateTimeout = setTimeout(() => {
        if(nextDeactivate) nextDeactivate();
    }, 10);
}

/**
 * onDrop is called when the drop event is emitted.
 * @param {DragEvent} e 
 * @returns 
 */
function onDrop(e) {
    if (!dndEnabled) {
        return;
    }
    e.preventDefault();

    // Find the close drop target element
    if (!e.target || !checkStyleDropTarget(getComputedStyle(e.target))) {
        return null;
    }

    if (canResolveFilePaths()) {
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
        resolveFilePaths(e.x, e.y, files);
    }

    // Trigger debounce function to deactivate drop targets
    if(nextDeactivate) nextDeactivate();

    // Deactivate all drop targets
    Array.from(document.getElementsByClassName(DROP_TARGET_ACTIVE)).forEach(el => el.classList.remove(DROP_TARGET_ACTIVE));

}

/**
 * postMessageWithAdditionalObjects checks the browser's capability of sending postMessageWithAdditionalObjects
 *
 * @returns {boolean}
 * @constructor
 */
function canResolveFilePaths() {
    return window.chrome?.webview?.postMessageWithAdditionalObjects != null;
}

/**
 * ResolveFilePaths sends drop events to the GO side to resolve file paths on windows.
 *
 * @param {number} x
 * @param {number} y
 * @param {any[]} files
 */
function resolveFilePaths(x, y, files) {
    // Only for windows webview2 >= 1.0.1774.30
    // https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2webmessagereceivedeventargs2?view=webview2-1.0.1823.32#applies-to
    if (window.chrome?.webview?.postMessageWithAdditionalObjects) {
        chrome.webview.postMessageWithAdditionalObjects(`wails:file:drop:${x}:${y}`, files);
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
 */
export function OnFileDrop(callback) {
    if (typeof callback !== "function") {
        console.error("DragAndDropCallback is not a function");
        return;
    }


    if (initialised) {
        return;
    }
    initialised = true;

    window.addEventListener('dragover', onDragOver);
    window.addEventListener('dragleave', onDragLeave);
    window.addEventListener('drop', onDrop);
    events.On("wails:file-drop", callback);
}

/**
 * OnFileDropOff removes the drag and drop listeners and handlers.
 */
export function OnFileDropOff() {
    window.removeEventListener('dragover', onDragOver);
    window.removeEventListener('dragleave', onDragLeave);
    window.removeEventListener('drop', onDrop);
    events.Off("wails:file-drop");
    initialised = false;
}
