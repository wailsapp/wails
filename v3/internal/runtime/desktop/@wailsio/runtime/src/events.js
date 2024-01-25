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
 * @typedef {import("./types").WailsEvent} WailsEvent
 */
import {newRuntimeCallerWithID, objectNames} from "./runtime";

import {EventTypes} from "./event_types";
export const Types = EventTypes;

// Setup
window._wails = window._wails || {};
window._wails.dispatchWailsEvent = dispatchWailsEvent;

const call = newRuntimeCallerWithID(objectNames.Events, '');
const EmitMethod = 0;
const eventListeners = new Map();

class Listener {
    constructor(eventName, callback, maxCallbacks) {
        this.eventName = eventName;
        this.maxCallbacks = maxCallbacks || -1;
        this.Callback = (data) => {
            callback(data);
            if (this.maxCallbacks === -1) return false;
            this.maxCallbacks -= 1;
            return this.maxCallbacks === 0;
        };
    }
}

export class WailsEvent {
    constructor(name, data = null) {
        this.name = name;
        this.data = data;
    }
}

export function setup() {
}

function dispatchWailsEvent(event) {
    let listeners = eventListeners.get(event.name);
    if (listeners) {
        let toRemove = listeners.filter(listener => {
            let remove = listener.Callback(event);
            if (remove) return true;
        });
        if (toRemove.length > 0) {
            listeners = listeners.filter(l => !toRemove.includes(l));
            if (listeners.length === 0) eventListeners.delete(event.name);
            else eventListeners.set(event.name, listeners);
        }
    }
}

/**
 * Register a callback function to be called multiple times for a specific event.
 *
 * @param {string} eventName - The name of the event to register the callback for.
 * @param {function} callback - The callback function to be called when the event is triggered.
 * @param {number} maxCallbacks - The maximum number of times the callback can be called for the event. Once the maximum number is reached, the callback will no longer be called.
 *
 @return {function} - A function that, when called, will unregister the callback from the event.
 */
export function OnMultiple(eventName, callback, maxCallbacks) {
    let listeners = eventListeners.get(eventName) || [];
    const thisListener = new Listener(eventName, callback, maxCallbacks);
    listeners.push(thisListener);
    eventListeners.set(eventName, listeners);
    return () => listenerOff(thisListener);
}

/**
 * Registers a callback function to be executed when the specified event occurs.
 *
 * @param {string} eventName - The name of the event.
 * @param {function} callback - The callback function to be executed. It takes no parameters.
 * @return {function} - A function that, when called, will unregister the callback from the event. */
export function On(eventName, callback) { return OnMultiple(eventName, callback, -1); }

/**
 * Registers a callback function to be executed only once for the specified event.
 *
 * @param {string} eventName - The name of the event.
 * @param {function} callback - The function to be executed when the event occurs.
 * @return {function} - A function that, when called, will unregister the callback from the event.
 */
export function Once(eventName, callback) { return OnMultiple(eventName, callback, 1); }

/**
 * Removes the specified listener from the event listeners collection.
 * If all listeners for the event are removed, the event key is deleted from the collection.
 *
 * @param {Object} listener - The listener to be removed.
 */
function listenerOff(listener) {
    const eventName = listener.eventName;
    let listeners = eventListeners.get(eventName).filter(l => l !== listener);
    if (listeners.length === 0) eventListeners.delete(eventName);
    else eventListeners.set(eventName, listeners);
}


/**
 * Removes event listeners for the specified event names.
 *
 * @param {string} eventName - The name of the event to remove listeners for.
 * @param {...string} additionalEventNames - Additional event names to remove listeners for.
 * @return {undefined}
 */
export function Off(eventName, ...additionalEventNames) {
    let eventsToRemove = [eventName, ...additionalEventNames];
    eventsToRemove.forEach(eventName => eventListeners.delete(eventName));
}
/**
 * Removes all event listeners.
 *
 * @function OffAll
 * @returns {void}
 */
export function OffAll() { eventListeners.clear(); }

/**
 * Emits an event using the given event name.
 *
 * @param {WailsEvent} event - The name of the event to emit.
 * @returns {any} - The result of the emitted event.
 */
export function Emit(event) { return call(EmitMethod, event); }
