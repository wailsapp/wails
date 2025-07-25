/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

import { newRuntimeCaller, objectNames } from "./runtime.js";
import { eventListeners, Listener, listenerOff } from "./listener.js";

// Setup
window._wails = window._wails || {};
window._wails.dispatchWailsEvent = dispatchWailsEvent;

const call = newRuntimeCaller(objectNames.Events);
const EmitMethod = 0;

export { Types } from "./event_types.js";

/**
 * The type of handlers for a given event.
 */
export type Callback = (ev: WailsEvent) => void;

/**
 * Represents a system event or a custom event emitted through wails-provided facilities.
 */
export class WailsEvent {
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

    constructor(name: string, data: any = null) {
        this.name = name;
        this.data = data;
    }
}

function dispatchWailsEvent(event: any) {
    let listeners = eventListeners.get(event.name);
    if (!listeners) {
        return;
    }

    let wailsEvent = new WailsEvent(event.name, event.data);
    if ('sender' in event) {
        wailsEvent.sender = event.sender;
    }

    listeners = listeners.filter(listener => !listener.dispatch(wailsEvent));
    if (listeners.length === 0) {
        eventListeners.delete(event.name);
    } else {
        eventListeners.set(event.name, listeners);
    }
}

/**
 * Register a callback function to be called multiple times for a specific event.
 *
 * @param eventName - The name of the event to register the callback for.
 * @param callback - The callback function to be called when the event is triggered.
 * @param maxCallbacks - The maximum number of times the callback can be called for the event. Once the maximum number is reached, the callback will no longer be called.
 * @returns A function that, when called, will unregister the callback from the event.
 */
export function OnMultiple(eventName: string, callback: Callback, maxCallbacks: number) {
    let listeners = eventListeners.get(eventName) || [];
    const thisListener = new Listener(eventName, callback, maxCallbacks);
    listeners.push(thisListener);
    eventListeners.set(eventName, listeners);
    return () => listenerOff(thisListener);
}

/**
 * Registers a callback function to be executed when the specified event occurs.
 *
 * @param eventName - The name of the event to register the callback for.
 * @param callback - The callback function to be called when the event is triggered.
 * @returns A function that, when called, will unregister the callback from the event.
 */
export function On(eventName: string, callback: Callback): () => void {
    return OnMultiple(eventName, callback, -1);
}

/**
 * Registers a callback function to be executed only once for the specified event.
 *
 * @param eventName - The name of the event to register the callback for.
 * @param callback - The callback function to be called when the event is triggered.
 * @returns A function that, when called, will unregister the callback from the event.
 */
export function Once(eventName: string, callback: Callback): () => void {
    return OnMultiple(eventName, callback, 1);
}

/**
 * Removes event listeners for the specified event names.
 *
 * @param eventNames - The name of the events to remove listeners for.
 */
export function Off(...eventNames: [string, ...string[]]): void {
    eventNames.forEach(eventName => eventListeners.delete(eventName));
}

/**
 * Removes all event listeners.
 */
export function OffAll(): void {
    eventListeners.clear();
}

/**
 * Emits an event using the name and data.
 *
 * @returns A promise that will be fulfilled once the event has been emitted.
 * @param name - the name of the event to emit.
 * @param data - the data to be sent with the event.
 */
export function Emit(name: string, data?: any): Promise<void> {
    let eventName: string;
    let eventData: any;

    if (typeof name === 'object' && name !== null && 'name' in name && 'data' in name) {
        // If name is an object with a name property, use it directly
        eventName = name['name'];
        eventData = name['data'];
    } else {
        // Otherwise use the standard parameters
        eventName = name as string;
        eventData = data;
    }

    return call(EmitMethod, { name: eventName, data: eventData });
}

