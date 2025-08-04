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

// Remove global callback handlers setup - HTTP-only implementation
// DELETE: window._wails = window._wails || {};
// DELETE: window._wails.callResultHandler = resultHandler;
// DELETE: window._wails.callErrorHandler = errorHandler;

// Remove callback response map - no longer needed for HTTP-only
// DELETE: const callResponses = new Map<string, PromiseResolvers>();

const call = newRuntimeCaller(objectNames.Call);
const cancelCall = newRuntimeCaller(objectNames.CancelCall);

const CallBinding = 0;
const CancelMethod = 0;

// Configuration
let bindingTimeout = 5 * 60 * 1000; // 5 minutes default

export function setBindingTimeout(timeout: number) {
    bindingTimeout = timeout;
}

export function getBindingTimeout(): number {
    return bindingTimeout;
}

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
 * Exception class for binding errors with HTTP status codes
 */
export class BindingError extends Error {
    constructor(
        public status: number,
        public kind: string,
        message: string,
        public cause?: any
    ) {
        super(message);
        this.name = "BindingError";
    }
}

/**
 * Exception class that will be thrown in case the bound method returns an error.
 * The value of the {@link RuntimeError#name} property is "RuntimeError".
 * @deprecated Use BindingError for new HTTP-only implementations
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

// Remove all callback handler functions - HTTP-only implementation
// DELETE: function resultHandler(...)
// DELETE: function errorHandler(...)
// DELETE: function getAndDeleteResponse(...)

/**
 * Generate unique ID for calls
 */
function generateID(): string {
    let result;
    do {
        result = nanoid();
    } while (false); // No longer need to check against callback map
    return result;
}

/**
 * Call a bound method - HTTP-only implementation
 *
 * In case of failure, the returned promise will reject with a BindingError
 * containing HTTP status information, error kind, message, and optional cause.
 *
 * @param options - A method call descriptor.
 * @returns The result of the call.
 */
export function Call(options: CallOptions): CancellablePromise<any> {
    const id = generateID();
    const abortController = new AbortController();
    
    const result = CancellablePromise.withResolvers<any>();
    
    // Make HTTP request with timeout and cancellation support
    const requestOptions = {
        signal: abortController.signal,
        timeout: getBindingTimeout()
    };
    
    const request = call(CallBinding, Object.assign({ "call-id": id }, options), requestOptions);
    
    request.then(response => {
        // Handle direct HTTP response
        if (!response.ok) {
            // Parse error response
            return response.json().then(errorData => {
                throw new BindingError(
                    response.status,
                    errorData.kind || 'UnknownError',
                    errorData.error || response.statusText,
                    errorData.cause
                );
            }).catch(parseError => {
                // If error response isn't JSON, create generic error
                throw new BindingError(
                    response.status,
                    'HttpError',
                    `HTTP ${response.status}: ${response.statusText}`
                );
            });
        }
        
        // Success - parse JSON response
        const contentType = response.headers.get('Content-Type');
        if (contentType && contentType.includes('application/json')) {
            return response.json();
        } else {
            return response.text();
        }
    }).then(data => {
        result.resolve(data);
    }).catch(err => {
        if (err.name === 'AbortError') {
            result.reject(new Error('Binding call cancelled'));
        } else {
            result.reject(err);
        }
    });
    
    // Cancellation support
    result.oncancelled = () => {
        abortController.abort();
        // Also try to cancel on backend (best effort)
        return cancelCall(CancelMethod, {"call-id": id}).catch(() => {
            // Ignore cancel request failures
        });
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
