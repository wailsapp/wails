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

import "./contextmenu";
import "./drag";

// Re-export public API
import * as Application from "./application";
import * as Browser from "./browser";
import * as Call from "./calls";
import * as Clipboard from "./clipboard";
import * as Create from "./create";
import * as Dialogs from "./dialogs";
import * as Events from "./events";
import * as Flags from "./flags";
import * as Screens from "./screens";
import * as System from "./system";
import Window from "./window";
import * as WML from "./wml";

export {
    Application,
    Browser,
    Call,
    Clipboard,
    Create,
    Dialogs,
    Events,
    Flags,
    Screens,
    System,
    Window,
    WML
};

// Notify backend
window._wails.invoke = System.invoke;
System.invoke("wails:runtime:ready");
