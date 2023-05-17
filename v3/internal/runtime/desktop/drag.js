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

let shouldDrag = false;

export function dragTest(e) {
    let val = window.getComputedStyle(e.target).getPropertyValue("--wails-draggable");
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

function onMouseDown(e) {
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
    document.body.style.cursor = window.wails.previousCursor || 'auto';
    shouldDrag = false;
}

function onMouseMove(e) {
    if (shouldDrag) {
        shouldDrag = false;
        let mousePressed = e.buttons !== undefined ? e.buttons : e.which;
        if (mousePressed > 0) {
            window.wails.previousCursor = document.body.style.cursor;
            document.body.style.cursor = 'grab';
            invoke("drag");
            return;
        }
    }
}