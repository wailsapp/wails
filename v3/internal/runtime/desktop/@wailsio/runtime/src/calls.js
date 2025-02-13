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
import { nanoid } from './nanoid.js';

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
        if (!data) {
            promiseHandler.resolve();
        } else if (!isJSON) {
            promiseHandler.resolve(data);
        } else {
            try {
                promiseHandler.resolve(JSON.parse(data));
            } catch (err) {
                promiseHandler.reject(new TypeError("could not parse result: " + err.message, { cause: err }));
            }
        }
    }
}

/**
 * Handles the error from a call request.
 *
 * @param {string} id - The id of the promise handler.
 * @param {string} data - The error data to reject the promise handler with.
 * @param {boolean} isJSON - Indicates whether the data is JSON or not.
 *
 * @return {void}
 */
function errorHandler(id, data, isJSON) {
    const promiseHandler = getAndDeleteResponse(id);
    if (promiseHandler) {
        if (!isJSON) {
            promiseHandler.reject(new Error(data));
        } else {
            let error;
            try {
                error = JSON.parse(data);
            } catch (err) {
                promiseHandler.reject(new TypeError("could not parse error: " + err.message, { cause: err }));
                return;
            }

            let options = {};
            if (error.cause) {
                options.cause = error.cause;
            }

            let exception;
            switch (error.kind) {
                case "ReferenceError":
                    exception = new ReferenceError(error.message, options);
                    break;
                case "TypeError":
                    exception = new TypeError(error.message, options);
                    break;
                case "RuntimeError":
                    exception = new RuntimeError(error.message, options);
                    break;
                default:
                    exception = new Error(error.message, options);
                    break;
            }

            promiseHandler.reject(exception);
        }
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
 * Collects all required information for a binding call.
 *
 * @typedef {Object} CallOptions
 * @property {number} [methodID] - The numeric ID of the bound method to call.
 * @property {string} [methodName] - The fully qualified name of the bound method to call.
 * @property {any[]} args - Arguments to be passed into the bound method.
 */

/**
 * Exception class that will be thrown in case the bound method returns an error.
 * The value of the {@link RuntimeError#name} property is "RuntimeError".
 */
export class RuntimeError extends Error {
    /**
     * Constructs a new RuntimeError instance.
     *
     * @param {string} message - The error message.
     * @param {any[]} args - Optional arguments for the Error constructor.
     */
    constructor(message, ...args) {
        super(message, ...args);
        this.name = "RuntimeError";
    }
}

/**
 * Call a bound method according to the given call options.
 *
 * In case of failure, the returned promise will reject with an exception
 * among ReferenceError (unknown method), TypeError (wrong argument count or type),
 * {@link RuntimeError} (method returned an error), or other (network or internal errors).
 * The exception might have a "cause" field with the value returned
 * by the application- or service-level error marshaling functions.
 *
 * @param {CallOptions} options - A method call descriptor.
 * @returns {Promise<any>} - The result of the call.
 */
export function Call(options) {
    const id = generateID();
    const doCancel = () => { return cancelCall(type, {"call-id": id}) };
    let queuedCancel = false, callRunning = false;
    let p = new Promise((resolve, reject) => {
        options["call-id"] = id;
        callResponses.set(id, { resolve, reject });
        call(CallBinding, options).then((_) => {
            callRunning = true;
            if (queuedCancel) {
                return doCancel();
            }
        }).catch((error) => {
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
 * Calls a bound method by name with the specified arguments.
 * See {@link Call} for details.
 *
 * @param {string} methodName - The name of the method in the format 'package.struct.method'.
 * @param {any[]} args - The arguments to pass to the method.
 * @returns {Promise<any>} The result of the method call.
 */
export function ByName(methodName, ...args) {
    return Call({
        methodName,
        args
    });
}

/**
 * Calls a method by its numeric ID with the specified arguments.
 * See {@link Call} for details.
 *
 * @param {number} methodID - The ID of the method to call.
 * @param {any[]} args - The arguments to pass to the method.
 * @return {Promise<any>} - The result of the method call.
 */
export function ByID(methodID, ...args) {
    return Call({
        methodID,
        args
    });
}
