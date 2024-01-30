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
window._wails = window._wails || {};
window._wails.setResizable = setResizable;
window._wails.endDrag = endDrag;
window.addEventListener('mousedown', onMouseDown);
window.addEventListener('mousemove', onMouseMove);
window.addEventListener('mouseup', onMouseUp);


let shouldDrag = false;
let resizeEdge = null;
let resizable = false;
let defaultCursor = "auto";

function dragTest(e) {
    let val = window.getComputedStyle(e.target).getPropertyValue("--webkit-app-region");
    if (!val || val === "" || val.trim() !== "drag" || e.buttons !== 1) {
        return false;
    }
    return e.detail === 1;
}

function setResizable(value) {
    resizable = value;
}

function endDrag() {
    document.body.style.cursor = 'default';
    shouldDrag = false;
}

function testResize() {
    if( resizeEdge ) {
        invoke(`resize:${resizeEdge}`);
        return true
    }
    return false;
}

function onMouseDown(e) {
    if(IsWindows() && testResize() || dragTest(e)) {
        shouldDrag = !!isValidDrag(e);
    }
}

function isValidDrag(e) {
    // Ignore drag on scrollbars
    return !(e.offsetX > e.target.clientWidth || e.offsetY > e.target.clientHeight);
}

function onMouseUp(e) {
    let mousePressed = e.buttons !== undefined ? e.buttons : e.which;
    if (mousePressed > 0) {
        endDrag();
    }
}

function setResize(cursor = defaultCursor) {
    document.documentElement.style.cursor = cursor;
    resizeEdge = cursor;
}

function onMouseMove(e) {
    shouldDrag = checkDrag(e);
    if (IsWindows() && resizable) {
        handleResize(e);
    }
}

function checkDrag(e) {
    let mousePressed = e.buttons !== undefined ? e.buttons : e.which;
    if(shouldDrag && mousePressed > 0) {
        invoke("drag");
        return false;
    }
    return shouldDrag;
}

function handleResize(e) {
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
