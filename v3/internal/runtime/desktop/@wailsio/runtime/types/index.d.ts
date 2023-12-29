export namespace Events {
    export { On };
    export { Off };
    export { OnMultiple };
    export { WailsEvent };
    export { OffAll };
    export { Emit };
    export { Once };
}
export namespace Call {
    export { Plugin };
    export { ByID };
    export { ByName };
}
export namespace Dialogs {
    export { Info };
    export { Error };
    export { OpenFile };
    export { Question };
    export { Warning };
    export { SaveFile };
}
import * as Application from "./application";
import * as Browser from "./browser";
import * as Clipboard from "./clipboard";
import * as ContextMenu from "./contextmenu";
import * as Flags from "./flags";
import * as Runtime from "./runtime";
import * as Screens from "./screens";
import * as System from "./system";
import * as Window from "./window";
import { On } from './events';
import { Off } from './events';
import { OnMultiple } from './events';
import { WailsEvent } from './events';
import { OffAll } from './events';
import { Emit } from './events';
import { Once } from './events';
import { Plugin } from "./calls";
import { ByID } from "./calls";
import { ByName } from "./calls";
import { Info } from "./dialogs";
import { Error } from "./dialogs";
import { OpenFile } from "./dialogs";
import { Question } from "./dialogs";
import { Warning } from "./dialogs";
import { SaveFile } from "./dialogs";
export { Application, Browser, Clipboard, ContextMenu, Flags, Runtime, Screens, System, Window };
