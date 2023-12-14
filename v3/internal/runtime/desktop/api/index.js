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

export * from "./clipboard";
export * from "./application";
export * from "./screens";
export * from "./dialogs";
export * from "./events";
export * from "./window";

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


