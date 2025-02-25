/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

import { invoke, IsWindows } from "./system.js";
import { GetFlag } from "./flags.js";
import { canTrackButtons } from "./utils.js";

// Setup
let canDrag = false;
let dragging = false;

let resizable = false;
let canResize = false;
let resizing = false;
let resizeEdge: string = "";
let defaultCursor = "auto";

let buttons = 0;
const buttonsTracked = canTrackButtons();

window._wails = window._wails || {};
window._wails.setResizable = (value: boolean): void => {
    resizable = value;
    if (!resizable) {
        // Stop resizing if in progress.
        canResize = resizing = false;
        setResize();
    }
};

window.addEventListener('mousedown', onMouseDown, { capture: true });
window.addEventListener('mousemove', onMouseMove, { capture: true });
window.addEventListener('mouseup', onMouseUp, { capture: true });
for (const ev of ['click', 'contextmenu', 'dblclick', 'pointerdown', 'pointerup']) {
    window.addEventListener(ev, suppressEvent, { capture: true });
}

function suppressEvent(event: Event) {
    if (dragging || resizing) {
        // Suppress all button events during dragging & resizing.
        event.stopImmediatePropagation();
        event.stopPropagation();
        event.preventDefault();
    }
}

function onMouseDown(event: MouseEvent): void {
    buttons = buttonsTracked ? event.buttons : (buttons | (1 << event.button));

    if (dragging || resizing) {
        // After dragging or resizing has started, only lifting the primary button can stop it.
        // Do not let any other events through.
        suppressEvent(event);
        return;
    }

    if ((canDrag || canResize) && (buttons & 1) && event.button !== 0) {
        // We were ready before, the primary is pressed and was not released:
        // still ready, but let events bubble through the window.
        return;
    }

    // Reset readiness state.
    canDrag = false;
    canResize = false;

    // Check for resizing readiness.
    if (resizeEdge) {
        if (event.button === 0 && event.detail === 1) {
            // Ready to resize if the primary button was pressed for the first time.
            canResize = true;
            invoke("wails:resize:" + resizeEdge);
        }

        // Do not start drag operations within resize edges.
        return;
    }

    let target: HTMLElement;

    if (event.target instanceof HTMLElement) {
        target = event.target;
    } else if (!(event.target instanceof HTMLElement) && event.target instanceof Node) {
        target = event.target.parentElement ?? document.body;
    } else {
        target = document.body;
    }

    const style = window.getComputedStyle(target);
    const setting = style.getPropertyValue("--wails-draggable").trim();
    if (setting === "drag" && event.button === 0 && event.detail === 1) {
        // Ready to drag if the primary button was pressed for the first time on a draggable element.
        // Ignore clicks on the scrollbar.
        if (
            event.offsetX - parseFloat(style.paddingLeft) < target.clientWidth
            && event.offsetY - parseFloat(style.paddingTop) < target.clientHeight
        ) {
            canDrag = true;
            invoke("wails:drag");
        }
    }
}

function onMouseUp(event: MouseEvent) {
    buttons = buttonsTracked ? event.buttons : (buttons & ~(1 << event.button));

    if (event.button === 0) {
        if (resizing) {
            // Let mouseup event bubble when a drag ends, but not when a resize ends.
            suppressEvent(event);
        }

        // Stop dragging and resizing when the primary button is lifted.
        canDrag = false;
        dragging = false;
        canResize = false;
        resizing = false;
        return;
    }

    // After dragging or resizing has started, only lifting the primary button can stop it.
    suppressEvent(event);
    return;
}

const cursorForEdge = Object.freeze({
    "se-resize": "nwse-resize",
    "sw-resize": "nesw-resize",
    "nw-resize": "nwse-resize",
    "ne-resize": "nesw-resize",
    "w-resize": "ew-resize",
    "n-resize": "ns-resize",
    "s-resize": "ns-resize",
    "e-resize": "ew-resize",
})

function setResize(edge?: keyof typeof cursorForEdge): void {
    if (edge) {
        if (!resizeEdge) { defaultCursor = document.body.style.cursor; }
        document.body.style.cursor = cursorForEdge[edge];
    } else if (!edge && resizeEdge) {
        document.body.style.cursor = defaultCursor;
    }

    resizeEdge = edge || "";
}

function onMouseMove(event: MouseEvent): void {
    if (canResize && resizeEdge) {
        // Start resizing.
        resizing = true;
    } else if (canDrag) {
        // Start dragging.
        dragging = true;
    }

    if (dragging || resizing) {
        // Either drag or resize is ongoing,
        // reset readiness and stop processing.
        canDrag = canResize = false;
        return;
    }

    if (!resizable || !IsWindows()) {
        if (resizeEdge) { setResize(); }
        return;
    }

    const resizeHandleHeight = GetFlag("system.resizeHandleHeight") || 5;
    const resizeHandleWidth = GetFlag("system.resizeHandleWidth") || 5;

    // Extra pixels for the corner areas.
    const cornerExtra = GetFlag("resizeCornerExtra") || 10;

    const rightBorder = (window.outerWidth - event.clientX) < resizeHandleWidth;
    const leftBorder = event.clientX < resizeHandleWidth;
    const topBorder = event.clientY < resizeHandleHeight;
    const bottomBorder = (window.outerHeight - event.clientY) < resizeHandleHeight;

    // Adjust for corner areas.
    const rightCorner = (window.outerWidth - event.clientX) < (resizeHandleWidth + cornerExtra);
    const leftCorner = event.clientX < (resizeHandleWidth + cornerExtra);
    const topCorner = event.clientY < (resizeHandleHeight + cornerExtra);
    const bottomCorner = (window.outerHeight - event.clientY) < (resizeHandleHeight + cornerExtra);

    if (!leftCorner && !topCorner && !bottomCorner && !rightCorner) {
        // Optimisation: out of all corner areas implies out of borders.
        setResize();
    }
    // Detect corners.
    else if (rightCorner && bottomCorner) setResize("se-resize");
    else if (leftCorner && bottomCorner) setResize("sw-resize");
    else if (leftCorner && topCorner) setResize("nw-resize");
    else if (topCorner && rightCorner) setResize("ne-resize");
    // Detect borders.
    else if (leftBorder) setResize("w-resize");
    else if (topBorder) setResize("n-resize");
    else if (bottomBorder) setResize("s-resize");
    else if (rightBorder) setResize("e-resize");
    // Out of border area.
    else setResize();
}
