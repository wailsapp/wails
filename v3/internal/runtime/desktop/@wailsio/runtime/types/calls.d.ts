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
export function Call(options: CallOptions): Promise<any>;
/**
 * Calls a bound method by name with the specified arguments.
 * See {@link Call} for details.
 *
 * @param {string} methodName - The name of the method in the format 'package.struct.method'.
 * @param {any[]} args - The arguments to pass to the method.
 * @returns {Promise<any>} The result of the method call.
 */
export function ByName(methodName: string, ...args: any[]): Promise<any>;
/**
 * Calls a method by its numeric ID with the specified arguments.
 * See {@link Call} for details.
 *
 * @param {number} methodID - The ID of the method to call.
 * @param {any[]} args - The arguments to pass to the method.
 * @return {Promise<any>} - The result of the method call.
 */
export function ByID(methodID: number, ...args: any[]): Promise<any>;
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
    constructor(message: string, ...args: any[]);
}
/**
 * Collects all required information for a binding call.
 */
export type CallOptions = {
    /**
     * - The numeric ID of the bound method to call.
     */
    methodID?: number;
    /**
     * - The fully qualified name of the bound method to call.
     */
    methodName?: string;
    /**
     * - Arguments to be passed into the bound method.
     */
    args: any[];
};
