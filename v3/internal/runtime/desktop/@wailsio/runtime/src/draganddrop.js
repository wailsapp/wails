/*
 _     __     _ __
| |   / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

/* jshint esversion: 9 */

import * as events from "./events";

let dndEnabled = false;
let initialised = false;
let dropTargets = [];
let currentDropTarget = null;
let rafId = null;
let observer = null;

const DROP_TARGET_ACTIVE = "wails-drop-target-active";
const CSS_DROP_PROPERTY = "--wails-drop-target";

window._wails = window._wails || {};

window._wails.enableDragAndDrop = function(value) {
    dndEnabled = value;
};

function initDropTargets() {
    const potentialTargets = document.querySelectorAll(`[style*="${CSS_DROP_PROPERTY}"]`);
    dropTargets = Array.from(potentialTargets)
        .filter(el => getComputedStyle(el).getPropertyValue(CSS_DROP_PROPERTY).trim() === 'drop')
        .map(el => ({
            element: el,
            rect: el.getBoundingClientRect()
        }));
}

function updateDropTargets() {
    dropTargets.forEach(target => {
        target.rect = target.element.getBoundingClientRect();
    });
}

function isPointInRect(x, y, rect) {
    return x >= rect.left && x <= rect.right && y >= rect.top && y <= rect.bottom;
}

function findDropTarget(x, y) {
    for (let i = dropTargets.length - 1; i >= 0; i--) {
        if (isPointInRect(x, y, dropTargets[i].rect)) {
            return dropTargets[i].element;
        }
    }
    return null;
}

function setDropTarget(target) {
    if (currentDropTarget !== target) {
        if (currentDropTarget) {
            currentDropTarget.classList.remove(DROP_TARGET_ACTIVE);
        }
        if (target) {
            console.log("Setting drop target to", target);
            target.classList.add(DROP_TARGET_ACTIVE);
        }
        currentDropTarget = target;
    }
}

function onDrag(e) {
    if (!dndEnabled) return;
    e.preventDefault();

    if (rafId) {
        cancelAnimationFrame(rafId);
    }

    rafId = requestAnimationFrame(() => {
        const target = findDropTarget(e.clientX, e.clientY);
        setDropTarget(target);
        rafId = null;
    });
}

function onDrop(e) {
    if (!dndEnabled) return;
    e.preventDefault();

    const target = findDropTarget(e.clientX, e.clientY);
    if (!target) return;

    if (canResolveFilePaths()) {
        const files = e.dataTransfer.items
            ? Array.from(e.dataTransfer.items).filter(item => item.kind === 'file').map(item => item.getAsFile())
            : Array.from(e.dataTransfer.files);

        resolveFilePaths(e.clientX, e.clientY, files);
    }

    dropTargets.forEach(target => {
        target.element.classList.remove(DROP_TARGET_ACTIVE);
    });
    currentDropTarget = null;
}

function canResolveFilePaths() {
    return window.chrome?.webview?.postMessageWithAdditionalObjects != null;
}

function resolveFilePaths(x, y, files) {
    if (window.chrome?.webview?.postMessageWithAdditionalObjects) {
        chrome.webview.postMessageWithAdditionalObjects(`wails:file:drop:${x}:${y}`, files);
    }
}

function setupMutationObserver() {
    observer = new MutationObserver((mutations) => {
        let needsRecalculation = false;
        for (let mutation of mutations) {
            if (mutation.type === 'childList') {
                needsRecalculation = true;
                break;
            } else if (mutation.type === 'attributes' && mutation.attributeName === 'style') {
                const target = mutation.target;
                const hasDropTarget = target.style.getPropertyValue(CSS_DROP_PROPERTY).trim() === 'drop';
                const wasDropTarget = dropTargets.some(dt => dt.element === target);

                if (hasDropTarget !== wasDropTarget) {
                    needsRecalculation = true;
                    break;
                }
            }
        }
        if (needsRecalculation) {
            recalculateDropTargets();
        }
    });

    observer.observe(document.body, {
        subtree: true,
        attributes: true,
        attributeFilter: ['style'],
        childList: true,
    });
}

function recalculateDropTargets() {
    initDropTargets();
}

export function OnFileDrop(callback) {
    if (typeof callback !== "function") {
        console.error("DragAndDropCallback is not a function");
        return;
    }

    if (initialised) return;
    initialised = true;

    initDropTargets();
    setupMutationObserver();

    window.addEventListener('dragover', onDrag);
    window.addEventListener('drop', onDrop);
    window.addEventListener('resize', updateDropTargets);
    events.On("wails:file-drop", callback);
}

export function OnFileDropOff() {
    window.removeEventListener('dragover', onDrag);
    window.removeEventListener('drop', onDrop);
    window.removeEventListener('resize', updateDropTargets);
    events.Off("wails:file-drop");
    initialised = false;
    dropTargets = [];
    currentDropTarget = null;
    if (rafId) {
        cancelAnimationFrame(rafId);
        rafId = null;
    }
    if (observer) {
        observer.disconnect();
        observer = null;
    }
}
