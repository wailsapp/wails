/*
 _       __      _ __    
| |     / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  ) 
|__/|__/\__,_/_/_/____/  
The electron alternative for Go
(c) Lea Anthony 2019-present
*/
/* jshint esversion: 9 */

import * as Log from "./log";
import * as Events from './events';
import * as Window from './window';
import * as Browser from './browser';

export function Quit() {
    window.runtime.Quit();
}


export default {
    ...Log,
    ...Events,
    ...Window,
    ...Browser,
    Quit
};