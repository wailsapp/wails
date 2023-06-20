import {newRuntimeCaller} from "./runtime";

let call = newRuntimeCaller("contextmenu");

function openContextMenu(id, x, y, data) {
    void call("OpenContextMenu", {id, x, y, data});
}

export function enableContextMenus(enabled) {
    if (enabled) {
        window.addEventListener('contextmenu', contextMenuHandler);
    } else {
        window.removeEventListener('contextmenu', contextMenuHandler);
    }
}

function contextMenuHandler(event) {
    let processed = processContextMenu(event.target, event);
    if (!processed) {
        let defaultContextMenuAction = window.getComputedStyle(event.target).getPropertyValue("--default-contextmenu");
        defaultContextMenuAction = defaultContextMenuAction ? defaultContextMenuAction.trim() : "";
        if (defaultContextMenuAction === 'hide') {
            event.preventDefault();
        }
    }
}

function processContextMenu(element, event) {
    let id = element.getAttribute('data-contextmenu');
    if (id) {
        event.preventDefault();
        openContextMenu(id, event.clientX, event.clientY, element.getAttribute('data-contextmenu-data'));
        return true;
    } else {
        let parent = element.parentElement;
        if (parent) {
            processContextMenu(parent, event);
        }
    }
    return false;
}
