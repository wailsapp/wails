import { CancellablePromise } from "./cancellable.js";
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
export declare class RuntimeError extends Error {
    /**
     * Constructs a new RuntimeError instance.
     * @param message - The error message.
     * @param options - Options to be forwarded to the Error constructor.
     */
    constructor(message?: string, options?: ErrorOptions);
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
export declare function Call(options: CallOptions): CancellablePromise<any>;
/**
 * Calls a bound method by name with the specified arguments.
 * See {@link Call} for details.
 *
 * @param methodName - The name of the method in the format 'package.struct.method'.
 * @param args - The arguments to pass to the method.
 * @returns The result of the method call.
 */
export declare function ByName(methodName: string, ...args: any[]): CancellablePromise<any>;
/**
 * Calls a method by its numeric ID with the specified arguments.
 * See {@link Call} for details.
 *
 * @param methodID - The ID of the method to call.
 * @param args - The arguments to pass to the method.
 * @return The result of the method call.
 */
export declare function ByID(methodID: number, ...args: any[]): CancellablePromise<any>;
