/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

/* jshint esversion: 9 */

import {newRuntimeCaller} from "./runtime";

let call = newRuntimeCaller("events");

/**
 * The Listener class defines a listener! :-)
 *
 * @class Listener
 */
class Listener {
    eventName: string;
    private maxCallbacks: number;
    Callback: (data?:any) => boolean;
    /**
     * Creates an instance of Listener.
     * @param {string} eventName
     * @param {function} callback
     * @param {number} maxCallbacks
     * @memberof Listener
     */
    constructor(eventName: string, callback: (data:any) => void, maxCallbacks?: number) {
        this.eventName = eventName;
        // Default of -1 means infinite
        this.maxCallbacks = maxCallbacks || -1;
        // Callback invokes the callback with the given data
        // Returns true if this listener should be destroyed
        this.Callback = (data) => {
            callback(data);
            // If maxCallbacks is infinite, return false (do not destroy)
            if (this.maxCallbacks === -1) {
                return false;
            }
            // Decrement maxCallbacks. Return true if now 0, otherwise false
            this.maxCallbacks -= 1;
            return this.maxCallbacks === 0;
        };
    }
}


/**
 * CustomEvent defines a custom event. It is passed to event listeners.
 *
 * @class CustomEvent
 */
export class CustomEvent {
    name: string;
    data: any;
    /**
     * Creates an instance of CustomEvent.
     * @param {string} name - Name of the event
     * @param {any} data - Data associated with the event
     * @memberof CustomEvent
     */
    constructor(name: string, data:any) {
        this.name = name;
        this.data = data;
    }
}

export const eventListeners = new Map();

/**
 * Registers an event listener that will be invoked `maxCallbacks` times before being destroyed
 *
 * @export
 * @param {string} eventName
 * @param {function(CustomEvent): void} callback
 * @param {number} maxCallbacks
 * @returns {function} A function to cancel the listener
 */
export function OnMultiple(eventName: string, callback: (data:CustomEvent) => void, maxCallbacks: number) {
    let listeners = eventListeners.get(eventName) || [];
    const thisListener = new Listener(eventName, callback, maxCallbacks);
    listeners.push(thisListener);
    eventListeners.set(eventName, listeners);
    return () => listenerOff(thisListener);
}

/**
 * Registers an event listener that will be invoked every time the event is emitted
 *
 * @export
 * @param {string} eventName
 * @param {function(CustomEvent): void} callback
 * @returns {function} A function to cancel the listener
 */
export function On(eventName: string, callback: (data:CustomEvent) => void) {
    return OnMultiple(eventName, callback, -1);
}

/**
 * Registers an event listener that will be invoked once then destroyed
 *
 * @export
 * @param {string} eventName
 * @param {function(CustomEvent): void} callback
 * @returns {function} A function to cancel the listener
 */
export function Once(eventName: string, callback: (data:CustomEvent) => void) {
    return OnMultiple(eventName, callback, 1);
}

/**
 * listenerOff unregisters a listener previously registered with On
 *
 * @param {Listener} listener
 */
function listenerOff(listener: Listener) {
    const eventName = listener.eventName;
    // Remove local listener
    let listeners = eventListeners.get(eventName).filter((l:Listener) => l !== listener);
    if (listeners.length === 0) {
        eventListeners.delete(eventName);
    } else {
        eventListeners.set(eventName, listeners);
    }
}

/**
 * dispatches an event to all listeners
 *
 * @export
 * @param {CustomEvent} event
 */
export function dispatchCustomEvent(event: CustomEvent) {
    console.log("dispatching event: ", {event});
    let listeners = eventListeners.get(event.name);
    if (listeners) {
        // iterate listeners and call callback. If callback returns true, remove listener
        let toRemove: any[] = [];
        listeners.forEach((listener: Listener) => {
            let remove = listener.Callback(event)
            if (remove) {
                toRemove.push(listener);
            }
        });
        // remove listeners
        if (toRemove.length > 0) {
            listeners = listeners.filter((l:Listener) => !toRemove.includes(l));
            if (listeners.length === 0) {
                eventListeners.delete(event.name);
            } else {
                eventListeners.set(event.name, listeners);
            }
        }
    }
}

/**
 * Off unregisters a listener previously registered with On,
 * optionally multiple listeners can be unregistered via `additionalEventNames`
 *
 [v3 CHANGE] Off only unregisters listeners within the current window
 *
 * @param {string} eventName
 * @param  {...string} additionalEventNames
 */
export function Off(eventName: string, ...additionalEventNames: string[]) {
    let eventsToRemove = [eventName, ...additionalEventNames];
    eventsToRemove.forEach(eventName => {
        eventListeners.delete(eventName);
    })
}

/**
 * OffAll unregisters all listeners
 * [v3 CHANGE] OffAll only unregisters listeners within the current window
 *
 */
export function OffAll() {
    eventListeners.clear();
}

/*
   Emit emits an event to all listeners
 */
export function Emit(event: CustomEvent) {
    return call("Emit", event);
}