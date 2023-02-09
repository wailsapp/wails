import {newRuntimeCaller} from "./runtime";

let call = newRuntimeCaller("contextmenu");

function openContextMenu(id, x, y, data) {
    return call("OpenContextMenu", {id, x, y, data});
}

function enableContextMenus(enabled) {
    if (enabled) {
        window.addEventListener('contextmenu', contextMenuHandler);
    } else {
        window.removeEventListener('contextmenu', contextMenuHandler);
    }
}

function contextMenuHandler(e) {
    let element = e.target;
    let contextMenuId = element.getAttribute("data-contextmenu-id");
    if (contextMenuId) {
        let contextMenuData = element.getAttribute("data-contextmenu-data");
        console.log({contextMenuId, contextMenuData, x: e.clientX, y: e.clientY});
        e.preventDefault();
        return openContextMenu(contextMenuId, e.clientX, e.clientY, contextMenuData);
    }
}

enableContextMenus(true);