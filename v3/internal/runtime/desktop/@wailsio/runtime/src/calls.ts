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

const call = newRuntimeCaller(objectNames.Call);
const cancelCall = newRuntimeCaller(objectNames.CancelCall);

const CallBinding = 0;
const CancelMethod = 0

/**
 * Holds all required information for a binding call.
 * May provide either a method ID or a method name, but not both.
 */
export type CallOptions = {
    /** The numeric ID of the bound method to call. */
    methodID: number;
    /** The fully qualified name of the bound method to call. */
    methodName?: never;
    /** Arguments to be passed into the bound method. */
    args: any[];
} | {
    /** The numeric ID of the bound method to call. */
    methodID?: never;
    /** The fully qualified name of the bound method to call. */
    methodName: string;
    /** Arguments to be passed into the bound method. */
    args: any[];
};

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
    constructor(message?: string, options?: ErrorOptions) {
        super(message, options);
        this.name = "RuntimeError";
    }
}

/**
 * Generates a unique ID using the nanoid library.
 *
 * @returns A unique ID.
 */
function generateID(): string {
    return nanoid();
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
export function Call(options: CallOptions): CancellablePromise<any> {
    const id = generateID();
    console.log('[Call] ===== HTTP-ONLY BINDINGS VERSION =====');
    console.log('[Call] Starting call with options:', options);

    const result = CancellablePromise.withResolvers<any>();

    // Make HTTP request that waits for response
    const request = call(CallBinding, Object.assign({ "call-id": id }, options));

    request.then((data) => {
        // data is already parsed JSON from runtimeCallWithID
        console.log('[Call] Received data:', data);
        if (data.error) {
            // Handle error response
            const error = data.error;
            console.log('[Call] Processing error:', error);
            let options: ErrorOptions = {};
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
            result.reject(exception);
        } else {
            // Handle success response
            console.log('[Call] Resolving with result:', data.result);
            result.resolve(data.result);
        }
    }, (err) => {
        // Handle HTTP/network errors
        console.log('[Call] HTTP/network error:', err);
        result.reject(err);
    });

    const cancel = () => {
        return cancelCall(CancelMethod, {"call-id": id}).catch((err) => {
            console.error("Error while requesting binding call cancellation:", err);
        });
    };

    result.oncancelled = () => {
        return request.then(cancel, cancel);
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
export function ByName(methodName: string, ...args: any[]): CancellablePromise<any> {
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
export function ByID(methodID: number, ...args: any[]): CancellablePromise<any> {
    return Call({ methodID, args });
}
