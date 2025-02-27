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
import { canTrackButtons, eventTarget } from "./utils.js";
// Setup
let canDrag = false;
let dragging = false;
let resizable = false;
let canResize = false;
let resizing = false;
let resizeEdge = "";
let defaultCursor = "auto";
let buttons = 0;
const buttonsTracked = canTrackButtons();
window._wails = window._wails || {};
window._wails.setResizable = (value) => {
    resizable = value;
    if (!resizable) {
        // Stop resizing if in progress.
        canResize = resizing = false;
        setResize();
    }
};
window.addEventListener('mousedown', update, { capture: true });
window.addEventListener('mousemove', update, { capture: true });
window.addEventListener('mouseup', update, { capture: true });
for (const ev of ['click', 'contextmenu', 'dblclick']) {
    window.addEventListener(ev, suppressEvent, { capture: true });
}
function suppressEvent(event) {
    // Suppress click events while resizing or dragging.
    if (dragging || resizing) {
        event.stopImmediatePropagation();
        event.stopPropagation();
        event.preventDefault();
    }
}
// Use constants to avoid comparing strings multiple times.
const MouseDown = 0;
const MouseUp = 1;
const MouseMove = 2;
function update(event) {
    // Windows suppresses mouse events at the end of dragging or resizing,
    // so we need to be smart and synthesize button events.
    let eventType, eventButtons = event.buttons;
    switch (event.type) {
        case 'mousedown':
            eventType = MouseDown;
            if (!buttonsTracked) {
                eventButtons = buttons | (1 << event.button);
            }
            break;
        case 'mouseup':
            eventType = MouseUp;
            if (!buttonsTracked) {
                eventButtons = buttons & ~(1 << event.button);
            }
            break;
        default:
            eventType = MouseMove;
            if (!buttonsTracked) {
                eventButtons = buttons;
            }
            break;
    }
    let released = buttons & ~eventButtons;
    let pressed = eventButtons & ~buttons;
    buttons = eventButtons;
    // Synthesize a release-press sequence if we detect a press of an already pressed button.
    if (eventType === MouseDown && !(pressed & event.button)) {
        released |= (1 << event.button);
        pressed |= (1 << event.button);
    }
    // Suppress all button events during dragging and resizing,
    // unless this is a mouseup event that is ending a drag action.
    if (eventType !== MouseMove // Fast path for mousemove
        && resizing
        || (dragging
            && (eventType === MouseDown
                || event.button !== 0))) {
        event.stopImmediatePropagation();
        event.stopPropagation();
        event.preventDefault();
    }
    // Handle releases
    if (released & 1) {
        primaryUp(event);
    }
    // Handle presses
    if (pressed & 1) {
        primaryDown(event);
    }
    // Handle mousemove
    if (eventType === MouseMove) {
        onMouseMove(event);
    }
    ;
}
function primaryDown(event) {
    // Reset readiness state.
    canDrag = false;
    canResize = false;
    // Ignore repeated clicks on macOS and Linux.
    if (!IsWindows()) {
        if (event.type === 'mousedown' && event.button === 0 && event.detail !== 1) {
            return;
        }
    }
    if (resizeEdge) {
        // Ready to resize if the primary button was pressed for the first time.
        canResize = true;
        // Do not start drag operations when on resize edges.
        return;
    }
    // Retrieve target element
    const target = eventTarget(event);
    // Ready to drag if the primary button was pressed for the first time on a draggable element.
    // Ignore clicks on the scrollbar.
    const style = window.getComputedStyle(target);
    canDrag = (style.getPropertyValue("--wails-draggable").trim() === "drag"
        && (event.offsetX - parseFloat(style.paddingLeft) < target.clientWidth
            && event.offsetY - parseFloat(style.paddingTop) < target.clientHeight));
}
function primaryUp(event) {
    // Stop dragging and resizing.
    canDrag = false;
    dragging = false;
    canResize = false;
    resizing = false;
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
});
function setResize(edge) {
    if (edge) {
        if (!resizeEdge) {
            defaultCursor = document.body.style.cursor;
        }
        document.body.style.cursor = cursorForEdge[edge];
    }
    else if (!edge && resizeEdge) {
        document.body.style.cursor = defaultCursor;
    }
    resizeEdge = edge || "";
}
function onMouseMove(event) {
    if (canResize && resizeEdge) {
        // Start resizing.
        resizing = true;
        invoke("wails:resize:" + resizeEdge);
    }
    else if (canDrag) {
        // Start dragging.
        dragging = true;
        invoke("wails:drag");
    }
    if (dragging || resizing) {
        // Either drag or resize is ongoing,
        // reset readiness and stop processing.
        canDrag = canResize = false;
        return;
    }
    if (!resizable || !IsWindows()) {
        if (resizeEdge) {
            setResize();
        }
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
    else if (rightCorner && bottomCorner)
        setResize("se-resize");
    else if (leftCorner && bottomCorner)
        setResize("sw-resize");
    else if (leftCorner && topCorner)
        setResize("nw-resize");
    else if (topCorner && rightCorner)
        setResize("ne-resize");
    // Detect borders.
    else if (leftBorder)
        setResize("w-resize");
    else if (topBorder)
        setResize("n-resize");
    else if (bottomBorder)
        setResize("s-resize");
    else if (rightBorder)
        setResize("e-resize");
    // Out of border area.
    else
        setResize();
}
