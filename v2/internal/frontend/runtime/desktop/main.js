/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The lightweight framework for web-like apps
(c) Lea Anthony 2019-present
*/
/* jshint esversion: 9 */
import * as Log from './log';
import {EventsEmit, EventsNotify, EventsOff, EventsOn, EventsOnce, EventsOnMultiple} from './events';
// import {Callback, SystemCall} from './calls';
// import {AddScript, DisableDefaultContextMenu, InjectCSS} from './utils';
// import {AddIPCListener, SendMessage} from 'ipc';

// Backend is where the Go struct wrappers get bound to
window.backend = {};

window.runtime = {
    ...Log,
    EventsOn,
    EventsOnce,
    EventsOnMultiple,
    EventsEmit,
    EventsOff,
};

// Initialise global if not already
window.wails = {
        //     Callback,
        EventsNotify,
        //     AddScript,
        //     InjectCSS,
        //     DisableDefaultContextMenu,
        //     // Init,
        //     AddIPCListener,
        //     SystemCall,
        //     SendMessage,
};

