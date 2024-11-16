/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/
/* jshint esversion: 9 */
import {invoke, IsWindows} from "./system";
import {GetFlag} from "./flags";

// Setup
let shouldDrag = false;
let resizable = false;
let resizeEdge = null;
let defaultCursor = "auto";

window._wails = window._wails || {};

window._wails.setResizable = function(value) {
    resizable = value;
};

window._wails.endDrag = function() {
    document.body.style.cursor = 'default';
    shouldDrag = false;
};

window.addEventListener('mousedown', onMouseDown);
window.addEventListener('mousemove', onMouseMove);
window.addEventListener('mouseup', onMouseUp);


function dragTest(e) {
    let val = window.getComputedStyle(e.target).getPropertyValue("--wails-draggable");
    let mousePressed = e.buttons !== undefined ? e.buttons : e.which;
    if (!val || val === "" || val.trim() !== "drag" || mousePressed === 0) {
        return false;
    }
    return e.detail === 1;
}

function onMouseDown(e) {

    // Check for resizing
    if (resizeEdge) {
        invoke("wails:resize:" + resizeEdge);
        e.preventDefault();
        return;
    }

    if (dragTest(e)) {
        // This checks for clicks on the scroll bar
        if (e.offsetX > e.target.clientWidth || e.offsetY > e.target.clientHeight) {
            return;
        }
        shouldDrag = true;
    } else {
        shouldDrag = false;
    }
}

function onMouseUp() {
    shouldDrag = false;
}

function setResize(cursor) {
    document.documentElement.style.cursor = cursor || defaultCursor;
    resizeEdge = cursor;
}

function onMouseMove(e) {
    if (shouldDrag) {
        shouldDrag = false;
        let mousePressed = e.buttons !== undefined ? e.buttons : e.which;
        if (mousePressed > 0) {
            invoke("wails:drag");
            return;
        }
    }
    if (!resizable || !IsWindows()) {
        return;
    }
    if (defaultCursor == null) {
        defaultCursor = document.documentElement.style.cursor;
    }
    let resizeHandleHeight = GetFlag("system.resizeHandleHeight") || 5;
    let resizeHandleWidth = GetFlag("system.resizeHandleWidth") || 5;

    // Extra pixels for the corner areas
    let cornerExtra = GetFlag("resizeCornerExtra") || 10;

    let rightBorder = window.outerWidth - e.clientX < resizeHandleWidth;
    let leftBorder = e.clientX < resizeHandleWidth;
    let topBorder = e.clientY < resizeHandleHeight;
    let bottomBorder = window.outerHeight - e.clientY < resizeHandleHeight;

    // Adjust for corners
    let rightCorner = window.outerWidth - e.clientX < (resizeHandleWidth + cornerExtra);
    let leftCorner = e.clientX < (resizeHandleWidth + cornerExtra);
    let topCorner = e.clientY < (resizeHandleHeight + cornerExtra);
    let bottomCorner = window.outerHeight - e.clientY < (resizeHandleHeight + cornerExtra);

    // If we aren't on an edge, but were, reset the cursor to default
    if (!leftBorder && !rightBorder && !topBorder && !bottomBorder && resizeEdge !== undefined) {
        setResize();
    }
    // Adjusted for corner areas
    else if (rightCorner && bottomCorner) setResize("se-resize");
    else if (leftCorner && bottomCorner) setResize("sw-resize");
    else if (leftCorner && topCorner) setResize("nw-resize");
    else if (topCorner && rightCorner) setResize("ne-resize");
    else if (leftBorder) setResize("w-resize");
    else if (topBorder) setResize("n-resize");
    else if (bottomBorder) setResize("s-resize");
    else if (rightBorder) setResize("e-resize");
}