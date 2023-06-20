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
Default: Show default context menu if contentEditable: true OR text has been selected OR --default-contextmenu: show OR tagname is input or textarea
--default-contextmenu: show will always show the context menu
--default-contextmenu: hide will always hide the context menu

Anything nested under a tag with --default-contextmenu: hide will not show the context menu unless it is explicitly set with --default-contextmenu: show
 */
function processDefaultContextMenu(event) {
    // Process default context menu
    let element = event.target;
    let defaultContextMenuAction = window.getComputedStyle(element).getPropertyValue("--default-contextmenu");
    defaultContextMenuAction = defaultContextMenuAction ? defaultContextMenuAction.trim() : "";
    switch(defaultContextMenuAction) {
        case "show":
            return;
        case "hide":
            event.preventDefault();
            return;
        default:
            // Check if contentEditable is true
            let contentEditable = element.getAttribute("contentEditable");
            if (contentEditable && contentEditable.toLowerCase() === "true") {
                return;
            }

            // Check if text has been selected
            let selection = window.getSelection();
            if (selection && selection.toString().length > 0) {
                return;
            }

            // Check if tagname is input or textarea
            let tagName = element.tagName.toLowerCase();
            if (tagName === "input" || tagName === "textarea") {
                return;
            }

            // hide default context menu
            event.preventDefault();
    }
}
