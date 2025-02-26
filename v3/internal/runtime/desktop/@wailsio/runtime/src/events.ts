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
import { Events as Create } from "./create.js";
import { Types } from "./event_types.js";

// Setup
window._wails = window._wails || {};
window._wails.dispatchWailsEvent = dispatchWailsEvent;

const call = newRuntimeCaller(objectNames.Events);
const EmitMethod = 0;

export * from "./event_types.js";

/**
 * A table of data types for all known events.
 * Will be monkey-patched by the binding generator.
 */
export interface CustomEvents {}

/**
 * Either a known event name or an arbitrary string.
 */
export type WailsEventName<E extends keyof CustomEvents = keyof CustomEvents> = E | (string & {});

/**
 * Union of all known system event names.
 */
type SystemEventName = {
    [K in keyof (typeof Types)]: (typeof Types)[K][keyof ((typeof Types)[K])]
} extends (infer M) ? M[keyof M] : never;

/**
 * The data type associated to a given event.
 */
export type WailsEventData<E extends WailsEventName = WailsEventName> =
    E extends keyof CustomEvents ? CustomEvents[E] : (E extends SystemEventName ? null : any);

/**
 * The type of handlers for a given event.
 */
export type WailsEventCallback<E extends WailsEventName = WailsEventName> = (ev: WailsEvent<E>) => void;

/**
 * Represents a system event or a custom event emitted through wails-provided facilities.
 */
export class WailsEvent<E extends WailsEventName = WailsEventName> {
    /**
     * The name of the event.
     */
    name: E;

    /**
     * Optional data associated with the emitted event.
     */
    data: WailsEventData<E>;

    /**
     * Name of the originating window. Omitted for application events.
     * Will be overridden if set manually.
     */
    sender?: string;

    constructor(name: E, data: WailsEventData<E>);
    constructor(name: E extends keyof CustomEvents ? never : E, data?: WailsEventData<E>)
    constructor(name: E, data?: any) {
        this.name = name;
        this.data = data ?? null;
    }
}

function dispatchWailsEvent(event: any) {
    let listeners = eventListeners.get(event.name);
    if (!listeners) {
        return;
    }

    let wailsEvent = new WailsEvent(
        event.name,
        (event.name in Create) ? Create[event.name](event.data) : event.data
    );
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
export function OnMultiple<E extends WailsEventName = WailsEventName>(eventName: E, callback: WailsEventCallback<E>, maxCallbacks: number) {
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
export function On<E extends WailsEventName = WailsEventName>(eventName: E, callback: WailsEventCallback<E>): () => void {
    return OnMultiple(eventName, callback, -1);
}

/**
 * Registers a callback function to be executed only once for the specified event.
 *
 * @param eventName - The name of the event to register the callback for.
 * @param callback - The callback function to be called when the event is triggered.
 * @returns A function that, when called, will unregister the callback from the event.
 */
export function Once<E extends WailsEventName = WailsEventName>(eventName: E, callback: WailsEventCallback<E>): () => void {
    return OnMultiple(eventName, callback, 1);
}

/**
 * Removes event listeners for the specified event names.
 *
 * @param eventNames - The name of the events to remove listeners for.
 */
export function Off(...eventNames: [WailsEventName, ...WailsEventName[]]): void {
    eventNames.forEach(eventName => eventListeners.delete(eventName));
}

/**
 * Removes all event listeners.
 */
export function OffAll(): void {
    eventListeners.clear();
}

/**
 * Emits the given event.
 *
 * @param event - The name of the event to emit.
 * @returns A promise that will be fulfilled once the event has been emitted.
 */
export function Emit<E extends WailsEventName = WailsEventName>(event: WailsEvent<E>): Promise<boolean> {
    return call(EmitMethod, event);
}
