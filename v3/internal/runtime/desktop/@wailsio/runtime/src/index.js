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

import * as Application from "./application";
import * as Browser from "./browser";
import * as Clipboard from "./clipboard";
import * as ContextMenu from "./contextmenu";
import * as Flags from "./flags";
import * as Runtime from "./runtime";
import * as Screens from "./screens";
import * as System from "./system";
import * as Window from "./window";

export { Application, Browser, Clipboard, ContextMenu, Flags, Runtime, Screens, System, Window };

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
