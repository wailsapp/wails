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
import {ByID, ByName, Plugin} from "./calls";

import * as Application from "./application";
import * as Browser from "./browser";
import * as Clipboard from "./clipboard";
import * as Flags from "./flags";
import * as Screens from "./screens";
import * as System from "./system";
import * as Window from "./window";
import * as WML from './wml';
import * as Events from "./events";
import * as Dialogs from "./dialogs";
import * as Call from "./calls";
import {setupEventCallbacks} from "./events";

export { Application, Browser, Call, Clipboard, Dialogs, Events, Flags, Screens, System, Window, WML};


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

let isReady = false
document.addEventListener('DOMContentLoaded', function() {
    isReady = true
})

function whenReady(fn) {
    if (isReady || document.readyState === 'complete') {
        fn();
    } else {
        document.addEventListener('DOMContentLoaded', fn);
    }
}

whenReady(() => {
    setupContextMenus();
    setupDrag();
    setupEventCallbacks();
    WML.Reload();
});
