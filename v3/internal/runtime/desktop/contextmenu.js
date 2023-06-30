import {newRuntimeCaller} from "./runtime";

let call = newRuntimeCaller("contextmenu");

function openContextMenu(id, x, y, data) {
    void call("OpenContextMenu", {id, x, y, data});
}

export function setupContextMenus() {
    window.addEventListener('contextmenu', contextMenuHandler);
}

function contextMenuHandler(event) {
    // Check for custom context menu
    let element = event.target;
    let customContextMenu = window.getComputedStyle(element).getPropertyValue("--custom-contextmenu");
    customContextMenu = customContextMenu ? customContextMenu.trim() : "";
    if (customContextMenu) {
        event.preventDefault();
        let customContextMenuData = window.getComputedStyle(element).getPropertyValue("--custom-contextmenu-data");
        openContextMenu(customContextMenu, event.clientX, event.clientY, customContextMenuData);
        return
    }

    processDefaultContextMenu(event);
}


/*
--default-contextmenu: auto; (default) will show the default context menu if contentEditable is true OR text has been selected OR element is input or textarea
--default-contextmenu: show; will always show the default context menu
--default-contextmenu: hide; will always hide the default context menu

This rule is inherited like normal CSS rules, so nesting works as expected
*/
function processDefaultContextMenu(event) {
    // Debug builds always show the menu
    if (DEBUG) {
        return;
    }

    // Process default context menu
    const element = event.target;
    const computedStyle = window.getComputedStyle(element);
    const defaultContextMenuAction = computedStyle.getPropertyValue("--default-contextmenu").trim();
    switch (defaultContextMenuAction) {
        case "show":
            return;
        case "hide":
            event.preventDefault();
            return;
        default:
            // Check if contentEditable is true
            if (element.isContentEditable) {
                return;
            }

            // Check if text has been selected
            const selection = window.getSelection();
            const hasSelection = (selection.toString().length > 0)
            if (hasSelection) {
                for (let i = 0; i < selection.rangeCount; i++) {
                    const range = selection.getRangeAt(i);
                    const rects = range.getClientRects();
                    for (let j = 0; j < rects.length; j++) {
                        const rect = rects[j];
                        if (document.elementFromPoint(rect.left, rect.top) === element) {
                            return;
                        }
                    }
                }
            }
            // Check if tagname is input or textarea
            if (element.tagName === "INPUT" || element.tagName === "TEXTAREA") {
                if (hasSelection || (!element.readOnly && !element.disabled)) {
                    return;
                }
            }

            // hide default context menu
            event.preventDefault();
    }
}
