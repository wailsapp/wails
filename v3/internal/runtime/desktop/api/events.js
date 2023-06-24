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
 * The Events API provides methods to interact with the event system.
 */
export const Events = {
    /**
     * Emit an event
     * @param {string} name
     * @param {any=} data
     */
    Emit: (name, data) => {
        return wails.Events.Emit(name, data);
    },
    /**
     * Subscribe to an event
     * @param {string} name - name of the event
     * @param {(any) => void} callback - callback to call when the event is emitted
     @returns {function()} unsubscribeMethod - method to unsubscribe from the event
     */
    On: (name, callback) => {
        return wails.Events.On(name, callback);
    },
    /**
     * Subscribe to an event once
     * @param {string} name - name of the event
     * @param {(any) => void} callback - callback to call when the event is emitted
     * @returns {function()} unsubscribeMethod - method to unsubscribe from the event
     */
    Once: (name, callback) => {
        return wails.Events.Once(name, callback);
    },
    /**
     * Subscribe to an event multiple times
     * @param {string} name - name of the event
     * @param {(any) => void} callback - callback to call when the event is emitted
     * @param {number} count - number of times to call the callback
     * @returns {Promise<void>} unsubscribeMethod - method to unsubscribe from the event
     */
    OnMultiple: (name, callback, count) => {
        return wails.Events.OnMultiple(name, callback, count);
    },
    /**
     * Unsubscribe from an event
     * @param {string} name - name of the event to unsubscribe from
     * @param {...string} additionalNames - additional names of events to unsubscribe from
     */
    Off: (name, ...additionalNames) => {
        wails.Events.Off(name, additionalNames);
    },
    /**
     * Unsubscribe all listeners from all events
     */
    OffAll: () => {
        wails.Events.OffAll();
    },
};
