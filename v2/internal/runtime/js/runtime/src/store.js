
import { writable } from 'svelte/store';
import {log} from "./log";

/** Overlay */
export const overlayVisible = writable(false);

export function showOverlay() {
    overlayVisible.set(true);
}
export function hideOverlay() {
    overlayVisible.set(false);
}

/** Menubar **/
export const menuVisible = writable(false);

export function showMenuBar() {
    menuVisible.set(true);
}
export function hideMenuBar() {
    menuVisible.set(false);
}

/** Trays **/

export const trays = writable([]);
export function setTray(tray) {
    trays.update((current) => {
        // Remove existing if it exists, else add
        const index = current.findIndex(item => item.ID === tray.ID);
        if ( index === -1 ) {
            current.push(tray);
        } else {
            current[index] = tray;
        }
        return current;
    })
}
export function updateTrayLabel(tray) {
    trays.update((current) => {
        // Remove existing if it exists, else add
        const index = current.findIndex(item => item.ID === tray.ID);
        if ( index === -1 ) {
            return log("ERROR: Attempted to update tray index ", tray.ID, "but it doesn't exist")
        }
        current[index].Label = tray.Label;
        return current;
    })
}

export function deleteTrayMenu(id) {
    trays.update((current) => {
        // Remove existing if it exists, else add
        const index = current.findIndex(item => item.ID === id);
        if ( index === -1 ) {
            return log("ERROR: Attempted to delete tray index ", id, "but it doesn't exist")
        }
        current.splice(index, 1);
        return current;
    })
}

export let selectedMenu = writable(null);