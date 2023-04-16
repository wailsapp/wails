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

/**
 * @typedef {import("./api/types").WailsEvent} WailsEvent
 */

import {newRuntimeCaller} from "./runtime";

let call = newRuntimeCaller("events");

/**
 * The Listener class defines a listener! :-)
 *
 * @class Listener
 */
class Listener {
    /**
     * Creates an instance of Listener.
     * @param {string} eventName
     * @param {function} callback
     * @param {number} maxCallbacks
     * @memberof Listener
     */
    constructor(eventName, callback, maxCallbacks) {
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
 * WailsEvent defines a custom event. It is passed to event listeners.
 *
 * @class WailsEvent
 * @property {string} name - Name of the event
 * @property {any} data - Data associated with the event
 */
export class WailsEvent {
    /**
     * Creates an instance of WailsEvent.
     * @param {string} name - Name of the event
     * @param {any=null} data - Data associated with the event
     * @memberof WailsEvent
     */
    constructor(name, data = null) {
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
 * @param {function(WailsEvent): void} callback
 * @param {number} maxCallbacks
 * @returns {function} A function to cancel the listener
 */
export function OnMultiple(eventName, callback, maxCallbacks) {
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
 * @param {function(WailsEvent): void} callback
 * @returns {function} A function to cancel the listener
 */
export function On(eventName, callback) {
    return OnMultiple(eventName, callback, -1);
}

/**
 * Registers an event listener that will be invoked once then destroyed
 *
 * @export
 * @param {string} eventName
 * @param {function(WailsEvent): void} callback
 @returns {function} A function to cancel the listener
 */
export function Once(eventName, callback) {
    return OnMultiple(eventName, callback, 1);
}

/**
 * listenerOff unregisters a listener previously registered with On
 *
 * @param {Listener} listener
 */
function listenerOff(listener) {
    const eventName = listener.eventName;
    // Remove local listener
    let listeners = eventListeners.get(eventName).filter(l => l !== listener);
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
 * @param {WailsEvent} event
 */
export function dispatchWailsEvent(event) {
    console.log("dispatching event: ", {event});
    let listeners = eventListeners.get(event.name);
    if (listeners) {
        // iterate listeners and call callback. If callback returns true, remove listener
        let toRemove = [];
        listeners.forEach(listener => {
            let remove = listener.Callback(event);
            if (remove) {
                toRemove.push(listener);
            }
        });
        // remove listeners
        if (toRemove.length > 0) {
            listeners = listeners.filter(l => !toRemove.includes(l));
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
export function Off(eventName, ...additionalEventNames) {
    let eventsToRemove = [eventName, ...additionalEventNames];
    eventsToRemove.forEach(eventName => {
        eventListeners.delete(eventName);
    });
}

/**
 * OffAll unregisters all listeners
 * [v3 CHANGE] OffAll only unregisters listeners within the current window
 *
 */
export function OffAll() {
    eventListeners.clear();
}

/**
 * Emit an event
 * @param {WailsEvent} event The event to emit
 */
export function Emit(event) {
    void call("Emit", event);
}