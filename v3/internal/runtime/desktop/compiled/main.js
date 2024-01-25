/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

import {debugLog} from "../@wailsio/runtime/src/log";

window._wails = window._wails || {};

import * as Application from "../@wailsio/runtime/src/application";
import * as Browser from "../@wailsio/runtime/src/browser";
import * as Clipboard from "../@wailsio/runtime/src/clipboard";
import * as Flags from "../@wailsio/runtime/src/flags";
import * as Screens from "../@wailsio/runtime/src/screens";
import * as System from "../@wailsio/runtime/src/system";
import * as Window from "../@wailsio/runtime/src/window";
import * as WML from '../@wailsio/runtime/src/wml';
import * as Events from "../@wailsio/runtime/src/events";
import * as Dialogs from "../@wailsio/runtime/src/dialogs";
import * as Call from "../@wailsio/runtime/src/calls";
import {invoke} from "../@wailsio/runtime/src/system";

/***
 This technique for proper load detection is taken from HTMX:

 BSD 2-Clause License

 Copyright (c) 2020, Big Sky Software
 All rights reserved.

 Redistribution and use in source and binary forms, with or without
 modification, are permitted provided that the following conditions are met:

 1. Redistributions of source code must retain the above copyright notice, this
 list of conditions and the following disclaimer.

 2. Redistributions in binary form must reproduce the above copyright notice,
 this list of conditions and the following disclaimer in the documentation
 and/or other materials provided with the distribution.

 THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

 ***/

window._wails.invoke=invoke;

window.wails = window.wails || {};
window.wails.Application = Application;
window.wails.Browser = Browser;
window.wails.Call = Call;
window.wails.Clipboard = Clipboard;
window.wails.Dialogs = Dialogs;
window.wails.Events = Events;
window.wails.Flags = Flags;
window.wails.Screens = Screens;
window.wails.System = System;
window.wails.Window = Window;
window.wails.WML = WML;


let isReady = false
document.addEventListener('DOMContentLoaded', function() {
    isReady = true
    window._wails.invoke('wails:runtime:ready');
    if(DEBUG) {
        debugLog("Wails Runtime Loaded");
    }
})

function whenReady(fn) {
    if (isReady || document.readyState === 'complete') {
        fn();
    } else {
        document.addEventListener('DOMContentLoaded', fn);
    }
}

whenReady(() => {
    WML.Reload();
});
