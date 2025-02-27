export { Types } from "./event_types.js";
/**
 * The type of handlers for a given event.
 */
export type Callback = (ev: WailsEvent) => void;
/**
 * Represents a system event or a custom event emitted through wails-provided facilities.
 */
export declare class WailsEvent {
    /**
     * The name of the event.
     */
    name: string;
    /**
     * Optional data associated with the emitted event.
     */
    data: any;
    /**
     * Name of the originating window. Omitted for application events.
     * Will be overridden if set manually.
     */
    sender?: string;
    constructor(name: string, data?: any);
}
/**
 * Register a callback function to be called multiple times for a specific event.
 *
 * @param eventName - The name of the event to register the callback for.
 * @param callback - The callback function to be called when the event is triggered.
 * @param maxCallbacks - The maximum number of times the callback can be called for the event. Once the maximum number is reached, the callback will no longer be called.
 * @returns A function that, when called, will unregister the callback from the event.
 */
export declare function OnMultiple(eventName: string, callback: Callback, maxCallbacks: number): () => void;
/**
 * Registers a callback function to be executed when the specified event occurs.
 *
 * @param eventName - The name of the event to register the callback for.
 * @param callback - The callback function to be called when the event is triggered.
 * @returns A function that, when called, will unregister the callback from the event.
 */
export declare function On(eventName: string, callback: Callback): () => void;
/**
 * Registers a callback function to be executed only once for the specified event.
 *
 * @param eventName - The name of the event to register the callback for.
 * @param callback - The callback function to be called when the event is triggered.
 * @returns A function that, when called, will unregister the callback from the event.
 */
export declare function Once(eventName: string, callback: Callback): () => void;
/**
 * Removes event listeners for the specified event names.
 *
 * @param eventNames - The name of the events to remove listeners for.
 */
export declare function Off(...eventNames: [string, ...string[]]): void;
/**
 * Removes all event listeners.
 */
export declare function OffAll(): void;
/**
 * Emits the given event.
 *
 * @param event - The name of the event to emit.
 * @returns A promise that will be fulfilled once the event has been emitted.
 */
export declare function Emit(event: WailsEvent): Promise<void>;
