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

import {EventsOn} from "./events";

const flags = {
    registered: false,
    defaultUseDropTarget: true,
    useDropTarget: true,
    prevElement: null
};

function onDragOver(e) {
    if (!window.wails.flags.enableWailsDragAndDrop) {
        return;
    }
    e.preventDefault();
    e.stopPropagation();

    if (!flags.useDropTarget) {
        return;
    }

    let targetElement = document.elementFromPoint(e.x, e.y);

    if (targetElement === flags.prevElement) {
        return;
    }

    const style = targetElement.style;
    let cssDropValue = null;
    if (Object.keys(style).findIndex(key => style[key] === window.wails.flags.cssDropProperty) < 0) {
        targetElement = targetElement.closest(`[style*='${window.wails.flags.cssDropProperty}']`);
    }

    if (targetElement === null) {
        return;
    }

    cssDropValue = window.getComputedStyle(targetElement).getPropertyValue(window.wails.flags.cssDropProperty);
    if (cssDropValue) {
        cssDropValue = cssDropValue.trim();
    }

    if (cssDropValue === window.wails.flags.cssDropValue) {
        targetElement.classList.add("wails-drop-target-active");
    } else if (flags.prevElement) {
        targetElement.classList.remove("wails-drop-target-active");
        flags.prevElement.classList.remove("wails-drop-target-active");
    }
    flags.prevElement = targetElement;
}

function onDragLeave(e) {
    if (!window.wails.flags.enableWailsDragAndDrop) {
        return;
    }
    e.preventDefault();
    e.stopPropagation();

    if (!flags.useDropTarget) {
        return;
    }

    let targetElement = document.elementFromPoint(e.x, e.y);
    let cssDropValue = window.getComputedStyle(targetElement).getPropertyValue(window.wails.flags.cssDropProperty);
    if (cssDropValue) {
        cssDropValue = cssDropValue.trim();
    }
    if (cssDropValue !== window.wails.flags.cssDropValue && flags.prevElement) {
        targetElement.classList.remove("wails-drop-target-active");
        flags.prevElement.classList.remove("wails-drop-target-active");
    }
}

function onDrop(e) {
    if (!window.wails.flags.enableWailsDragAndDrop) {
        return;
    }
    e.preventDefault();
    e.stopPropagation();

    if (!flags.useDropTarget) {
        return;
    }

    let targetElement = document.elementFromPoint(e.x, e.y);
    let cssDropValue = window.getComputedStyle(targetElement).getPropertyValue(window.wails.flags.cssDropProperty);
    if (cssDropValue) {
        cssDropValue = cssDropValue.trim();
    }
    if (cssDropValue !== window.wails.flags.cssDropValue) {
        if (flags.prevElement) {
            targetElement.classList.remove("wails-drop-target-active");
            flags.prevElement.classList.remove("wails-drop-target-active");
        }
        return;
    }

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

    if (flags.prevElement) {
        flags.prevElement.classList.remove("wails-drop-target-active");
    }
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
 * @param x
 * @param y
 * @param files
 * @constructor
 */
export function ResolveFilePaths(x, y, files) {
    // Only for windows webview2 >= 1.0.1774.30
    // https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2webmessagereceivedeventargs2?view=webview2-1.0.1823.32#applies-to
    if (window.chrome?.webview?.postMessageWithAdditionalObjects) {
        chrome.webview.postMessageWithAdditionalObjects(`file:drop:${x}:${y}`, files);
        return;
    }
    console.warn("unsupported platform");
}

/**
 * Callback for DragAndDropOn returns a slice of file path strings when a drop is finished.
 *
 * @export
 * @callback DragAndDropCallback
 * @param {number} x - x coordinate of the drop
 * @param {number} y - y coordinate of the drop
 * @param {string[]} paths - A list of file paths.
 */

/**
 * DragAndDropOn listens to drag and drop events and calls the callback with the coordinates of the drop and an array of path strings.
 *
 * @export
 * @param {DragAndDropCallback} callback - Callback for DragAndDropOn returns a slice of file path strings when a drop is finished.
 * @param {boolean} [useDropTarget=true] - Only call the callback when the drop finished on an element that has the drop target style. (--wails-drop-target)
 */
export function DragAndDropOn(callback, useDropTarget) {
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
            let targetElement = document.elementFromPoint(x, y);
            if (!targetElement) {
                return;
            }
            let cssDropValue = window.getComputedStyle(targetElement).getPropertyValue(window.wails.flags.cssDropProperty);
            if (cssDropValue) {
                cssDropValue = cssDropValue.trim();
            }
            if (cssDropValue !== window.wails.flags.cssDropValue) {
                return;
            }
            callback(x, y, paths);
        }
    }

    EventsOn("wails:file-drop", cb);
}

/**
 * DragAndDropOff removes the drag and drop listeners and handlers.
 */
export function DragAndDropOff() {
    window.removeEventListener('dragover', onDragOver);
    window.removeEventListener('dragleave', onDragLeave);
    window.removeEventListener('drop', onDrop);
    EventsOff("wails:file-drop");
    flags.registered = false;
}