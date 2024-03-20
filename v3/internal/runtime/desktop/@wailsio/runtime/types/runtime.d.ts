/**
 * Creates a runtime caller function that invokes a specified method on a given object within a specified window context.
 *
 * @param {Object} object - The object on which the method is to be invoked.
 * @param {string} windowName - The name of the window context in which the method should be called.
 * @returns {Function} A runtime caller function that takes the method name and optionally arguments and invokes the method within the specified window context.
 */
export function newRuntimeCaller(object: any, windowName: string): Function;
/**
 * Creates a new runtime caller with specified ID.
 *
 * @param {object} object - The object to invoke the method on.
 * @param {string} windowName - The name of the window.
 * @return {Function} - The new runtime caller function.
 */
export function newRuntimeCallerWithID(object: object, windowName: string): Function;
export namespace objectNames {
    let Call: number;
    let Clipboard: number;
    let Application: number;
    let Events: number;
    let ContextMenu: number;
    let Dialog: number;
    let Window: number;
    let Screens: number;
    let System: number;
    let Browser: number;
    let CancelCall: number;
}
export let clientId: string;
