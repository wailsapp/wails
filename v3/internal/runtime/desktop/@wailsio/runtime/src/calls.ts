/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

import { CancellablePromise, type CancellablePromiseWithResolvers } from "./cancellable.js";
import { newRuntimeCaller, objectNames } from "./runtime.js";
import { nanoid } from "./nanoid.js";

// Setup
window._wails = window._wails || {};

type PromiseResolvers = Omit<CancellablePromiseWithResolvers<any>, "promise" | "oncancelled">

const call = newRuntimeCaller(objectNames.Call);
const cancelCall = newRuntimeCaller(objectNames.CancelCall);
const callResponses = new Map<string, PromiseResolvers>();

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
 * @returns A unique ID that does not exist in the callResponses set.
 */
function generateID(): string {
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
export function Call(options: CallOptions): CancellablePromise<any> {
    const id = generateID();

    const result = CancellablePromise.withResolvers<any>();
    callResponses.set(id, { resolve: result.resolve, reject: result.reject });

    const request = call(CallBinding, Object.assign({ "call-id": id }, options));
    let running = true;

    request.then((res) => {
        running = false;
        callResponses.delete(id);
        result.resolve(res);
    }, (err) => {
        running = false;
        callResponses.delete(id);
        result.reject(err);
    });

    const cancel = () => {
        callResponses.delete(id);
        return cancelCall(CancelMethod, {"call-id": id}).catch((err) => {
            console.error("Error while requesting binding call cancellation:", err);
        });
    };

    result.oncancelled = () => {
        if (running) {
            return cancel();
        } else {
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
