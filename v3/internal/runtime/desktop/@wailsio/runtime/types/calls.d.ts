/**
 * Call method.
 *
 * @param {Object} options - The options for the method.
 * @returns {Object} - The result of the call.
 */
export function Call(options: any): any;
/**
 * Executes a method by name.
 *
 * @param {string} name - The name of the method in the format 'package.struct.method'.
 * @param {...*} args - The arguments to pass to the method.
 * @throws {Error} If the name is not a string or is not in the correct format.
 * @returns {*} The result of the method execution.
 */
export function ByName(name: string, ...args: any[]): any;
/**
 * Calls a method by its ID with the specified arguments.
 *
 * @param {number} methodID - The ID of the method to call.
 * @param {...*} args - The arguments to pass to the method.
 * @return {*} - The result of the method call.
 */
export function ByID(methodID: number, ...args: any[]): any;
/**
 * Calls a method on a plugin.
 *
 * @param {string} pluginName - The name of the plugin.
 * @param {string} methodName - The name of the method to call.
 * @param {...*} args - The arguments to pass to the method.
 * @returns {*} - The result of the method call.
 */
export function Plugin(pluginName: string, methodName: string, ...args: any[]): any;
