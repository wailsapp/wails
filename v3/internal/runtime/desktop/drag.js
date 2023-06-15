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

import {invoke} from "./invoke";
import {GetFlag} from "./flags";

let shouldDrag = false;

export function dragTest(e) {
    // if (window.wails.Capabilities['HasNativeDrag'] === true) {
    //     return false;
    // }

    let val = window.getComputedStyle(e.target).getPropertyValue("app-region");
    if (val) {
        val = val.trim();
    }

    if (val !== "drag") {
        return false;
    }

    // Only process the primary button
    if (e.buttons !== 1) {
        return false;
    }

    return e.detail === 1;
}

export function setupDrag() {
    window.addEventListener('mousedown', onMouseDown);
    window.addEventListener('mousemove', onMouseMove);
    window.addEventListener('mouseup', onMouseUp);
}

let resizeEdge = null;

function testResize(e) {
    if( resizeEdge ) {
        invoke("resize:" + resizeEdge);
        return true
    }
    return false;
}

function onMouseDown(e) {

    // Check for resizing on Windows
    if( WINDOWS ) {
        if (testResize()) {
            return;
        }
    }
    if (dragTest(e)) {
        // Ignore drag on scrollbars
        if (e.offsetX > e.target.clientWidth || e.offsetY > e.target.clientHeight) {
            return;
        }
        shouldDrag = true;
    } else {
        shouldDrag = false;
    }
}

function onMouseUp(e) {
    let mousePressed = e.buttons !== undefined ? e.buttons : e.which;
    if (mousePressed > 0) {
        endDrag();
    }
}

export function endDrag() {
    document.body.style.cursor = 'default';
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
            invoke("drag");
        }
        return;
    }

    if (WINDOWS) {
        handleResize(e);
    }
}

let defaultCursor = "auto";

function handleResize(e) {
    let resizeHandleHeight = GetFlag("system.resizeHandleHeight") || 5;
    let resizeHandleWidth = GetFlag("system.resizeHandleWidth") || 5;

    // Extra pixels for the corner areas
    let cornerExtra = GetFlag("resizeCornerExtra") || 3;

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
