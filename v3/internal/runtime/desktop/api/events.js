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

import {EventTypes} from './event_types';
import {WailsEvent} from '../events';


/**
 * The Events API provides methods to interact with the event system.
 */
export const Events = {
    /**
     * Emit an event
     * @param {string} name
     * @param {any=} data
     * @returns {Promise<void>}
     */
    Emit: (name, data) => wails.Events.Emit(new wails.Events.WailsEvent(name, data)),
    /**
     * Subscribe to an event
     * @param {string} name - name of the event
     * @param {(any) => void} callback - callback to call when the event is emitted
     * @returns {() => void} unsubscribeMethod - method to unsubscribe from the event
     */
    On: (name, callback) => wails.Events.On(name, callback),
    /**
     * Subscribe to an event once
     * @param {string} name - name of the event
     * @param {(any) => void} callback - callback to call when the event is emitted
     * @returns {() => void} unsubscribeMethod - method to unsubscribe from the event
     */
    Once: (name, callback) => wails.Events.Once(name, callback),
    /**
     * Subscribe to an event multiple times
     * @param {string} name - name of the event
     * @param {(any) => void} callback - callback to call when the event is emitted
     * @param {number} count - number of times to call the callback
     * @returns {() => void} unsubscribeMethod - method to unsubscribe from the event
     */
    OnMultiple: (name, callback, count) => wails.Events.OnMultiple(name, callback, count),
    /**
     * Unsubscribe from an event
     * @param {string} name - name of the event to unsubscribe from
     * @param {...string} additionalNames - additional names of events to unsubscribe from
     * @returns {void}
     */
    Off: (name, ...additionalNames) => wails.Events.Off(name, additionalNames),
    /**
     * Unsubscribe all listeners from all events
     * @returns {void}
     */
    OffAll: () => wails.Events.OffAll(),

    Windows: EventTypes.Windows,
    Mac: EventTypes.Mac,
    Common: EventTypes.Common,

};
