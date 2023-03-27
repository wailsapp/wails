/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

/* jshint esversion: 9 */

import {newRuntimeCaller} from "./runtime";

let call : (method: string, args?: any) => Promise<any> = newRuntimeCaller("clipboard");

// SetText sets the clipboard text
export function SetText(text: string) : Promise<void> {
    return call("SetText", {text});
}

// Text returns the clipboard text
export function Text(): Promise<string> {
    return call("Text");
}