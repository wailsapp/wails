import {newRuntimeCaller} from "./runtime";

let call = newRuntimeCaller("contextmenu");

function openContextMenu(id, x, y, data) {
    return call("OpenContextMenu", {id, x, y, data});
}

export function enableContextMenus(enabled) {
    if (enabled) {
        window.addEventListener('contextmenu', contextMenuHandler);
    } else {
        window.removeEventListener('contextmenu', contextMenuHandler);
    }
}

function contextMenuHandler(event) {
    processContextMenu(event.target, event);
}

function processContextMenu(element, event) {
    let id = element.getAttribute('data-contextmenu');
    if (id) {
        event.preventDefault();
        openContextMenu(id, event.clientX, event.clientY, element.getAttribute('data-contextmenu-data'));
    } else {
        let parent = element.parentElement;
        if (parent) {
            processContextMenu(parent, event);
        }
    }
}
