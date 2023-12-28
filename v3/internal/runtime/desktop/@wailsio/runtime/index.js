/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

import {setupContextMenus} from "./contextmenu";
import {setupDrag} from "./drag";
import {reloadWML} from "./wml";
import {Emit, Off, OffAll, On, Once, OnMultiple, WailsEvent} from './events';
import {ByID, ByName, Plugin} from "./calls";
import {Error, Info, OpenFile, Question, SaveFile, Warning} from "./dialogs";

export * as Application from "./application";
export * as Browser from "./browser";
export * as Clipboard from "./clipboard";
export * as ContextMenu from "./contextmenu";
export * as Flags from "./flags";
export * as Runtime from "./runtime";
export * as Screens from "./screens";
export * as System from "./system";
export * as Window from "./window";

export const Events = {
    On,
    Off,
    OnMultiple,
    WailsEvent,
    OffAll,
    Emit,
    Once

}

export const Call = {
    Plugin,
    ByID,
    ByName
}

export const Dialogs = {
    Info,
    Error,
    OpenFile, Question, Warning, SaveFile
}

setupContextMenus();
setupDrag();

document.addEventListener("DOMContentLoaded", function () {
    reloadWML();
});
