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

import { nanoid } from 'nanoid/non-secure';

let call = newRuntimeCaller("call");

let callResponses = new Map();

function generateID() {
    let result;
    do {
        result = nanoid();
    } while (callResponses.has(result));
    return result;
}

export function callCallback(id, data, isJSON) {
    let p = callResponses.get(id);
    if (p) {
        if (isJSON) {
            p.resolve(JSON.parse(data));
        } else {
            p.resolve(data);
        }
        callResponses.delete(id);
    }
}

export function callErrorCallback(id, message) {
    let p = callResponses.get(id);
    if (p) {
        p.reject(message);
        callResponses.delete(id);
    }
}

function callBinding(type, options) {
    return new Promise((resolve, reject) => {
        let id = generateID();
        options = options || {};
        options["call-id"] = id;
        callResponses.set(id, {resolve, reject});
        call(type, options).catch((error) => {
            reject(error);
            callResponses.delete(id);
        });
    });
}

export function Call(options) {
    return callBinding("Call", options);
}

/**
 * Call a plugin method
 * @param {string} pluginName - name of the plugin
 * @param {string} methodName - name of the method
 * @param {...any} args - arguments to pass to the method
 * @returns {Promise<any>} - promise that resolves with the result
 */
export function Plugin(pluginName, methodName, ...args) {
    return callBinding("Call", {
        packageName: "wails-plugins",
        structName: pluginName,
        methodName: methodName,
        args: args,
    });
}