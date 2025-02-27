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
import Window from "./window.js";
import * as WML from "./wml.js";
export { Application, Browser, Call, Clipboard, Dialogs, Events, Flags, Screens, System, Window, WML };
/**
 * An internal utility consumed by the binding generator.
 *
 * @ignore
 * @internal
 */
export { Create };
export * from "./cancellable.js";
// Notify backend
window._wails.invoke = System.invoke;
System.invoke("wails:runtime:ready");
