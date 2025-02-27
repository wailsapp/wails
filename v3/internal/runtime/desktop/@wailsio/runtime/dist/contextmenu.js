/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/
import { newRuntimeCaller, objectNames } from "./runtime.js";
import { IsDebug } from "./system.js";
import { eventTarget } from "./utils";
// setup
window.addEventListener('contextmenu', contextMenuHandler);
const call = newRuntimeCaller(objectNames.ContextMenu);
const ContextMenuOpen = 0;
function openContextMenu(id, x, y, data) {
    void call(ContextMenuOpen, { id, x, y, data });
}
function contextMenuHandler(event) {
    const target = eventTarget(event);
    // Check for custom context menu
    const customContextMenu = window.getComputedStyle(target).getPropertyValue("--custom-contextmenu").trim();
    if (customContextMenu) {
        event.preventDefault();
        const data = window.getComputedStyle(target).getPropertyValue("--custom-contextmenu-data");
        openContextMenu(customContextMenu, event.clientX, event.clientY, data);
    }
    else {
        processDefaultContextMenu(event, target);
    }
}
/*
--default-contextmenu: auto; (default) will show the default context menu if contentEditable is true OR text has been selected OR element is input or textarea
--default-contextmenu: show; will always show the default context menu
--default-contextmenu: hide; will always hide the default context menu

This rule is inherited like normal CSS rules, so nesting works as expected
*/
function processDefaultContextMenu(event, target) {
    // Debug builds always show the menu
    if (IsDebug()) {
        return;
    }
    // Process default context menu
    switch (window.getComputedStyle(target).getPropertyValue("--default-contextmenu").trim()) {
        case 'show':
            return;
        case 'hide':
            event.preventDefault();
            return;
    }
    // Check if contentEditable is true
    if (target.isContentEditable) {
        return;
    }
    // Check if text has been selected
    const selection = window.getSelection();
    const hasSelection = selection && selection.toString().length > 0;
    if (hasSelection) {
        for (let i = 0; i < selection.rangeCount; i++) {
            const range = selection.getRangeAt(i);
            const rects = range.getClientRects();
            for (let j = 0; j < rects.length; j++) {
                const rect = rects[j];
                if (document.elementFromPoint(rect.left, rect.top) === target) {
                    return;
                }
            }
        }
    }
    // Check if tag is input or textarea.
    if (target instanceof HTMLInputElement || target instanceof HTMLTextAreaElement) {
        if (hasSelection || (!target.readOnly && !target.disabled)) {
            return;
        }
    }
    // hide default context menu
    event.preventDefault();
}
