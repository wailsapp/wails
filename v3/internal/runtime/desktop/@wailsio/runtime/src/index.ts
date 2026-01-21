/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

// Setup
window._wails = window._wails || {};

import "./contextmenu.js";
import "./drag.js";

// Re-export public API
import * as Application from "./application.js";
import * as Browser from "./browser.js";
import * as Call from "./calls.js";
import * as Clipboard from "./clipboard.js";
import * as Create from "./create.js";
import * as Dialogs from "./dialogs.js";
import * as Events from "./events.js";
import * as Flags from "./flags.js";
import * as Screens from "./screens.js";
import * as System from "./system.js";
import * as IOS from "./ios.js";
import Window, { handleDragEnter, handleDragLeave, handleDragOver } from "./window.js";
import * as WML from "./wml.js";

export {
    Application,
    Browser,
    Call,
    Clipboard,
    Dialogs,
    Events,
    Flags,
    Screens,
    System,
    IOS,
    Window,
    WML
};

/**
 * An internal utility consumed by the binding generator.
 *
 * @ignore
 */
export { Create };

export * from "./cancellable.js";

// Export transport interfaces and utilities
export {
    setTransport,
    getTransport,
    type RuntimeTransport,
    objectNames,
    clientId,
} from "./runtime.js";

// Notify backend
window._wails.invoke = System.invoke;

// Register platform handlers (internal API)
// Note: Window is the thisWindow instance (default export from window.ts)
// Binding ensures 'this' correctly refers to the current window instance
window._wails.handlePlatformFileDrop = Window.HandlePlatformFileDrop.bind(Window);

// Linux-specific drag handlers (GTK intercepts DOM drag events)
window._wails.handleDragEnter = handleDragEnter;
window._wails.handleDragLeave = handleDragLeave;
window._wails.handleDragOver = handleDragOver;

System.invoke("wails:runtime:ready");

// Load optional window init scripts (fire and forget, matching current async behavior)
// The backend identifies the window from the x-wails-window-id header
fetch('/wails/init.js')
    .then(r => r.ok && r.status !== 204 ? r.text() : null)
    .then(js => { if (js) eval(js); })
    .catch(() => {}); // Silently ignore errors

fetch('/wails/init.css')
    .then(r => r.ok && r.status !== 204 ? r.text() : null)
    .then(css => {
        if (css) {
            const style = document.createElement('style');
            style.textContent = css;
            document.head.appendChild(style);
        }
    })
    .catch(() => {}); // Silently ignore errors
