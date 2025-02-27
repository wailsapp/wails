/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/
import { CancellablePromise } from "./cancellable.js";
import { newRuntimeCaller, objectNames } from "./runtime.js";
import { nanoid } from "./nanoid.js";
// Setup
window._wails = window._wails || {};
window._wails.callResultHandler = resultHandler;
window._wails.callErrorHandler = errorHandler;
const call = newRuntimeCaller(objectNames.Call);
const cancelCall = newRuntimeCaller(objectNames.CancelCall);
const callResponses = new Map();
const CallBinding = 0;
const CancelMethod = 0;
/**
 * Exception class that will be thrown in case the bound method returns an error.
 * The value of the {@link RuntimeError#name} property is "RuntimeError".
 */
export class RuntimeError extends Error {
    /**
     * Constructs a new RuntimeError instance.
     * @param message - The error message.
     * @param options - Options to be forwarded to the Error constructor.
     */
    constructor(message, options) {
        super(message, options);
        this.name = "RuntimeError";
    }
}
/**
 * Handles the result of a call request.
 *
 * @param id - The id of the request to handle the result for.
 * @param data - The result data of the request.
 * @param isJSON - Indicates whether the data is JSON or not.
 */
function resultHandler(id, data, isJSON) {
    const resolvers = getAndDeleteResponse(id);
    if (!resolvers) {
        return;
    }
    if (!data) {
        resolvers.resolve(undefined);
    }
    else if (!isJSON) {
        resolvers.resolve(data);
    }
    else {
        try {
            resolvers.resolve(JSON.parse(data));
        }
        catch (err) {
            resolvers.reject(new TypeError("could not parse result: " + err.message, { cause: err }));
        }
    }
}
/**
 * Handles the error from a call request.
 *
 * @param id - The id of the promise handler.
 * @param data - The error data to reject the promise handler with.
 * @param isJSON - Indicates whether the data is JSON or not.
 */
function errorHandler(id, data, isJSON) {
    const resolvers = getAndDeleteResponse(id);
    if (!resolvers) {
        return;
    }
    if (!isJSON) {
        resolvers.reject(new Error(data));
    }
    else {
        let error;
        try {
            error = JSON.parse(data);
        }
        catch (err) {
            resolvers.reject(new TypeError("could not parse error: " + err.message, { cause: err }));
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
        resolvers.reject(exception);
    }
}
/**
 * Retrieves and removes the response associated with the given ID from the callResponses map.
 *
 * @param id - The ID of the response to be retrieved and removed.
 * @returns The response object associated with the given ID, if any.
 */
function getAndDeleteResponse(id) {
    const response = callResponses.get(id);
    callResponses.delete(id);
    return response;
}
/**
 * Generates a unique ID using the nanoid library.
 *
 * @returns A unique ID that does not exist in the callResponses set.
 */
function generateID() {
    let result;
    do {
        result = nanoid();
    } while (callResponses.has(result));
    return result;
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
 * @param options - A method call descriptor.
 * @returns The result of the call.
 */
export function Call(options) {
    const id = generateID();
    const result = CancellablePromise.withResolvers();
    callResponses.set(id, { resolve: result.resolve, reject: result.reject });
    const request = call(CallBinding, Object.assign({ "call-id": id }, options));
    let running = false;
    request.then(() => {
        running = true;
    }, (err) => {
        callResponses.delete(id);
        result.reject(err);
    });
    const cancel = () => {
        callResponses.delete(id);
        return cancelCall(CancelMethod, { "call-id": id }).catch((err) => {
            console.error("Error while requesting binding call cancellation:", err);
        });
    };
    result.oncancelled = () => {
        if (running) {
            return cancel();
        }
        else {
            return request.then(cancel);
        }
    };
    return result.promise;
}
/**
 * Calls a bound method by name with the specified arguments.
 * See {@link Call} for details.
 *
 * @param methodName - The name of the method in the format 'package.struct.method'.
 * @param args - The arguments to pass to the method.
 * @returns The result of the method call.
 */
export function ByName(methodName, ...args) {
    return Call({ methodName, args });
}
/**
 * Calls a method by its numeric ID with the specified arguments.
 * See {@link Call} for details.
 *
 * @param methodID - The ID of the method to call.
 * @param args - The arguments to pass to the method.
 * @return The result of the method call.
 */
export function ByID(methodID, ...args) {
    return Call({ methodID, args });
}
