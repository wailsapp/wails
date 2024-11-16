/*
 _     __     _ __
| |  / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

import {OpenURL} from "./browser";
import {Question} from "./dialogs";
import {Emit, WailsEvent} from "./events";
import {canAbortListeners, whenReady} from "./utils";
import Window from "./window";

/**
 * Sends an event with the given name and optional data.
 *
 * @param {string} eventName - The name of the event to send.
 * @param {any} [data=null] - Optional data to send along with the event.
 *
 * @return {void}
 */
function sendEvent(eventName, data=null) {
    Emit(new WailsEvent(eventName, data));
}

/**
 * Calls a method on a specified window.
 * @param {string} windowName - The name of the window to call the method on.
 * @param {string} methodName - The name of the method to call.
 */
function callWindowMethod(windowName, methodName) {
    const targetWindow = Window.Get(windowName);
    const method = targetWindow[methodName];

    if (typeof method !== "function") {
        console.error(`Window method '${methodName}' not found`);
        return;
    }

    try {
        method.call(targetWindow);
    } catch (e) {
        console.error(`Error calling window method '${methodName}': `, e);
    }
}

/**
 * Responds to a triggering event by running appropriate WML actions for the current target
 *
 * @param {Event} ev
 * @return {void}
 */
function onWMLTriggered(ev) {
    const element = ev.currentTarget;

    function runEffect(choice = "Yes") {
        if (choice !== "Yes")
            return;

        const eventType = element.getAttribute('wml-event');
        const targetWindow = element.getAttribute('wml-target-window') || "";
        const windowMethod = element.getAttribute('wml-window');
        const url = element.getAttribute('wml-openurl');

        if (eventType !== null)
            sendEvent(eventType);
        if (windowMethod !== null)
            callWindowMethod(targetWindow, windowMethod);
        if (url !== null)
            void OpenURL(url);
    }

    const confirm = element.getAttribute('wml-confirm');

    if (confirm) {
        Question({
            Title: "Confirm",
            Message: confirm,
            Detached: false,
            Buttons: [
                { Label: "Yes" },
                { Label: "No", IsDefault: true }
            ]
        }).then(runEffect);
    } else {
        runEffect();
    }
}

/**
 * @type {symbol}
 */
const controller = Symbol();

/**
 * AbortControllerRegistry does not actually remember active event listeners: instead
 * it ties them to an AbortSignal and uses an AbortController to remove them all at once.
 */
class AbortControllerRegistry {
    constructor() {
        /**
         * Stores the AbortController that can be used to remove all currently active listeners.
         *
         * @private
         * @name {@link controller}
         * @member {AbortController}
         */
        this[controller] = new AbortController();
    }

    /**
     * Returns an options object for addEventListener that ties the listener
     * to the AbortSignal from the current AbortController.
     *
     * @param {HTMLElement} element An HTML element
     * @param {string[]} triggers The list of active WML trigger events for the specified elements
     * @returns {AddEventListenerOptions}
     */
    set(element, triggers) {
        return { signal: this[controller].signal };
    }

    /**
     * Removes all registered event listeners.
     *
     * @returns {void}
     */
    reset() {
        this[controller].abort();
        this[controller] = new AbortController();
    }
}

/**
 * @type {symbol}
 */
const triggerMap = Symbol();

/**
 * @type {symbol}
 */
const elementCount = Symbol();

/**
 * WeakMapRegistry maps active trigger events to each DOM element through a WeakMap.
 * This ensures that the mapping remains private to this module, while still allowing garbage
 * collection of the involved elements.
 */
class WeakMapRegistry {
    constructor() {
        /**
         * Stores the current element-to-trigger mapping.
         *
         * @private
         * @name {@link triggerMap}
         * @member {WeakMap<HTMLElement, string[]>}
         */
        this[triggerMap] = new WeakMap();

        /**
         * Counts the number of elements with active WML triggers.
         *
         * @private
         * @name {@link elementCount}
         * @member {number}
         */
        this[elementCount] = 0;
    }

    /**
     * Sets the active triggers for the specified element.
     *
     * @param {HTMLElement} element An HTML element
     * @param {string[]} triggers The list of active WML trigger events for the specified element
     * @returns {AddEventListenerOptions}
     */
    set(element, triggers) {
        this[elementCount] += !this[triggerMap].has(element);
        this[triggerMap].set(element, triggers);
        return {};
    }

    /**
     * Removes all registered event listeners.
     *
     * @returns {void}
     */
    reset() {
        if (this[elementCount] <= 0)
            return;

        for (const element of document.body.querySelectorAll('*')) {
            if (this[elementCount] <= 0)
                break;

            const triggers = this[triggerMap].get(element);
            this[elementCount] -= (typeof triggers !== "undefined");

            for (const trigger of triggers || [])
                element.removeEventListener(trigger, onWMLTriggered);
        }

        this[triggerMap] = new WeakMap();
        this[elementCount] = 0;
    }
}

const triggerRegistry = canAbortListeners() ? new AbortControllerRegistry() : new WeakMapRegistry();

/**
 * Adds event listeners to the specified element.
 *
 * @param {HTMLElement} element
 * @return {void}
 */
function addWMLListeners(element) {
    const triggerRegExp = /\S+/g;
    const triggerAttr = (element.getAttribute('wml-trigger') || "click");
    const triggers = [];

    let match;
    while ((match = triggerRegExp.exec(triggerAttr)) !== null)
        triggers.push(match[0]);

    const options = triggerRegistry.set(element, triggers);
    for (const trigger of triggers)
        element.addEventListener(trigger, onWMLTriggered, options);
}

/**
 * Schedules an automatic reload of WML to be performed as soon as the document is fully loaded.
 *
 * @return {void}
 */
export function Enable() {
    whenReady(Reload);
}

/**
 * Reloads the WML page by adding necessary event listeners and browser listeners.
 *
 * @return {void}
 */
export function Reload() {
    triggerRegistry.reset();
    document.body.querySelectorAll('[wml-event], [wml-window], [wml-openurl]').forEach(addWMLListeners);
}
