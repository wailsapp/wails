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

// Setup
window._wails = window._wails || {};
window._wails.callResultHandler = resultHandler;
window._wails.callErrorHandler = errorHandler;


const CallBinding = 0;
const call = newRuntimeCallerWithID(objectNames.Call, '');
const cancelCall = newRuntimeCallerWithID(objectNames.CancelCall, '');
let callResponses = new Map();

/**
 * Generates a unique ID using the nanoid library.
 *
 * @return {string} - A unique ID that does not exist in the callResponses set.
 */
function generateID() {
    let result;
    do {
        result = nanoid();
    } while (callResponses.has(result));
    return result;
}

/**
 * Handles the result of a call request.
 *
 * @param {string} id - The id of the request to handle the result for.
 * @param {string} data - The result data of the request.
 * @param {boolean} isJSON - Indicates whether the data is JSON or not.
 *
 * @return {undefined} - This method does not return any value.
 */
function resultHandler(id, data, isJSON) {
    const promiseHandler = getAndDeleteResponse(id);
    if (promiseHandler) {
        promiseHandler.resolve(isJSON ? JSON.parse(data) : data);
    }
}

/**
 * Handles the error from a call request.
 *
 * @param {string} id - The id of the promise handler.
 * @param {string} message - The error message to reject the promise handler with.
 *
 * @return {void}
 */
function errorHandler(id, message) {
    const promiseHandler = getAndDeleteResponse(id);
    if (promiseHandler) {
        promiseHandler.reject(message);
    }
}

/**
 * Retrieves and removes the response associated with the given ID from the callResponses map.
 *
 * @param {any} id - The ID of the response to be retrieved and removed.
 *
 * @returns {any} The response object associated with the given ID.
 */
function getAndDeleteResponse(id) {
    const response = callResponses.get(id);
    callResponses.delete(id);
    return response;
}

/**
 * Executes a call using the provided type and options.
 *
 * @param {string|number} type - The type of call to execute.
 * @param {Object} [options={}] - Additional options for the call.
 * @return {Promise} - A promise that will be resolved or rejected based on the result of the call. It also has a cancel method to cancel a long running request.
 */
function callBinding(type, options = {}) {
    const id = generateID();
    const doCancel = () => { return cancelCall(type, {"call-id": id}) };
    let queuedCancel = false, callRunning = false;
    let p = new Promise((resolve, reject) => {
        options["call-id"] = id;
        callResponses.set(id, { resolve, reject });
        call(type, options).
            then((_) => {
                callRunning = true;
                if (queuedCancel) {
                    return doCancel();
                }
            }).
            catch((error) => {
                reject(error);
                callResponses.delete(id);
            });
    });
    p.cancel = () => {
        if (callRunning) {
            return doCancel();
        } else {
            queuedCancel = true;
        }
    };

    return p;
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
 * @param {string} methodName - The name of the method in the format 'package.struct.method'.
 * @param {...*} args - The arguments to pass to the method.
 * @throws {Error} If the name is not a string or is not in the correct format.
 * @returns {*} The result of the method execution.
 */
export function ByName(methodName, ...args) {
    return callBinding(CallBinding, {
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
