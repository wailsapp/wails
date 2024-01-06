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
import { newRuntimeCallerWithID, objectNames } from "./runtime";
import { nanoid } from 'nanoid/non-secure';

const CallBinding = 0;
const call = newRuntimeCallerWithID(objectNames.Call, '');
let callResponses = new Map();

window._wails = window._wails || {};
window._wails.callResultHandler = resultHandler;
window._wails.callErrorHandler = errorHandler;

function generateID() {
    let result;
    do {
        result = nanoid();
    } while (callResponses.has(result));
    return result;
}

export function resultHandler(id, data, isJSON) {
    const promiseHandler = getAndDeleteResponse(id);
    if (promiseHandler) {
        promiseHandler.resolve(isJSON ? JSON.parse(data) : data);
    }
}

export function errorHandler(id, message) {
    const promiseHandler = getAndDeleteResponse(id);
    if (promiseHandler) {
        promiseHandler.reject(message);
    }
}

function getAndDeleteResponse(id) {
    const response = callResponses.get(id);
    callResponses.delete(id);
    return response;
}

function callBinding(type, options = {}) {
    return new Promise((resolve, reject) => {
        const id = generateID();
        options["call-id"] = id;
        callResponses.set(id, { resolve, reject });
        call(type, options).catch((error) => {
            reject(error);
            callResponses.delete(id);
        });
    });
}

/**
 * Call method.
 *
 * @param {Object} options - The options for the method.
 * @returns {Object} - The result of the call.
 */
export function Call(options) {
    return callBinding(CallBinding, options);
}

/**
 * Executes a method by name.
 *
 * @param {string} name - The name of the method in the format 'package.struct.method'.
 * @param {...*} args - The arguments to pass to the method.
 * @throws {Error} If the name is not a string or is not in the correct format.
 * @returns {*} The result of the method execution.
 */
export function ByName(name, ...args) {
    if (typeof name !== "string" || name.split(".").length !== 3) {
        throw new Error("CallByName requires a string in the format 'package.struct.method'");
    }
    let [packageName, structName, methodName] = name.split(".");
    return callBinding(CallBinding, {
        packageName,
        structName,
        methodName,
        args
    });
}

/**
 * Calls a method by its ID with the specified arguments.
 *
 * @param {number} methodID - The ID of the method to call.
 * @param {...*} args - The arguments to pass to the method.
 * @return {*} - The result of the method call.
 */
export function ByID(methodID, ...args) {
    return callBinding(CallBinding, {
        methodID,
        args
    });
}

/**
 * Calls a method on a plugin.
 *
 * @param {string} pluginName - The name of the plugin.
 * @param {string} methodName - The name of the method to call.
 * @param {...*} args - The arguments to pass to the method.
 * @returns {*} - The result of the method call.
 */
export function Plugin(pluginName, methodName, ...args) {
    return callBinding(CallBinding, {
        packageName: "wails-plugins",
        structName: pluginName,
        methodName,
        args
    });
}
