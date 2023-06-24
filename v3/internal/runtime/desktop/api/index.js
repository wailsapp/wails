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

import * as Clipboard from "./clipboard";
import * as Application from "./application";
import * as Screens from "./screens";
import * as Dialogs from "./dialogs";
import * as Events from "./events";
import * as Window from "./window";

export { Clipboard, Application, Screens, Dialogs, Events, Window };

/**
 * Call a plugin method
 * @param {string} pluginName - name of the plugin
 * @param {string} methodName - name of the method
 * @param {...any} args - arguments to pass to the method
 * @returns {Promise<any>} - promise that resolves with the result
 */
export const Plugin = (pluginName, methodName, ...args) => {
    return wails.Plugin(pluginName, methodName, ...args);
};


