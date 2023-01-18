import {writable} from 'svelte/store';

/** Overlay */
export const overlayVisible = writable(false);

export function showOverlay() {
    overlayVisible.set(true);
}

export function hideOverlay() {
    overlayVisible.set(false);
}
