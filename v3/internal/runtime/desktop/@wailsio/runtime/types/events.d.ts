export function dispatchWailsEvent(event: any): void;
/**
 * Register a callback function to be called multiple times for a specific event.
 *
 * @param {string} eventName - The name of the event to register the callback for.
 * @param {function} callback - The callback function to be called when the event is triggered.
 * @param {number} maxCallbacks - The maximum number of times the callback can be called for the event. Once the maximum number is reached, the callback will no longer be called.
 *
 @return {function} - A function that, when called, will unregister the callback from the event.
 */
export function OnMultiple(eventName: string, callback: Function, maxCallbacks: number): Function;
/**
 * Registers a callback function to be executed when the specified event occurs.
 *
 * @param {string} eventName - The name of the event.
 * @param {function} callback - The callback function to be executed. It takes no parameters.
 * @return {function} - A function that, when called, will unregister the callback from the event. */
export function On(eventName: string, callback: Function): Function;
/**
 * Registers a callback function to be executed only once for the specified event.
 *
 * @param {string} eventName - The name of the event.
 * @param {function} callback - The function to be executed when the event occurs.
 * @return {void@return {function} - A function that, when called, will unregister the callback from the event.
 */
export function Once(eventName: string, callback: Function): void;
/**
 * Removes event listeners for the specified event names.
 *
 * @param {string} eventName - The name of the event to remove listeners for.
 * @param {...string} additionalEventNames - Additional event names to remove listeners for.
 * @return {undefined}
 */
export function Off(eventName: string, ...additionalEventNames: string[]): undefined;
/**
 * Removes all event listeners.
 *
 * @function OffAll
 * @returns {void}
 */
export function OffAll(): void;
/**
 * Emits an event using the given event name.
 *
 * @param {WailsEvent} event - The name of the event to emit.
 * @returns {any} - The result of the emitted event.
 */
export function Emit(event: WailsEvent): any;
export class WailsEvent {
    constructor(name: any, data?: any);
    name: any;
    data: any;
}
