/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

import * as Runtime from "../@wailsio/runtime/src";

// NOTE: the following methods MUST be imported explicitly because of how esbuild injection works
import {Enable as EnableWML} from "../@wailsio/runtime/src/wml";
import {debugLog} from "../@wailsio/runtime/src/utils";

window.wails = Runtime;
EnableWML();

if (DEBUG) {
    debugLog("Wails Runtime Loaded")
}
