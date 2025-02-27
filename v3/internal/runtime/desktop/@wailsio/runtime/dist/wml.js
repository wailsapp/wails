/*
 _     __     _ __
| |  / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/
import { OpenURL } from "./browser.js";
import { Question } from "./dialogs.js";
import { Emit, WailsEvent } from "./events.js";
import { canAbortListeners, whenReady } from "./utils.js";
import Window from "./window.js";
/**
 * Sends an event with the given name and optional data.
 *
 * @param eventName - - The name of the event to send.
 * @param [data=null] - - Optional data to send along with the event.
 */
function sendEvent(eventName, data = null) {
    Emit(new WailsEvent(eventName, data));
}
/**
 * Calls a method on a specified window.
 *
 * @param windowName - The name of the window to call the method on.
 * @param methodName - The name of the method to call.
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
    }
    catch (e) {
        console.error(`Error calling window method '${methodName}': `, e);
    }
}
/**
 * Responds to a triggering event by running appropriate WML actions for the current target.
 */
function onWMLTriggered(ev) {
    const element = ev.currentTarget;
    function runEffect(choice = "Yes") {
        if (choice !== "Yes")
            return;
        const eventType = element.getAttribute('wml-event') || element.getAttribute('data-wml-event');
        const targetWindow = element.getAttribute('wml-target-window') || element.getAttribute('data-wml-target-window') || "";
        const windowMethod = element.getAttribute('wml-window') || element.getAttribute('data-wml-window');
        const url = element.getAttribute('wml-openurl') || element.getAttribute('data-wml-openurl');
        if (eventType !== null)
            sendEvent(eventType);
        if (windowMethod !== null)
            callWindowMethod(targetWindow, windowMethod);
        if (url !== null)
            void OpenURL(url);
    }
    const confirm = element.getAttribute('wml-confirm') || element.getAttribute('data-wml-confirm');
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
    }
    else {
        runEffect();
    }
}
// Private field names.
const controllerSym = Symbol("controller");
const triggerMapSym = Symbol("triggerMap");
const elementCountSym = Symbol("elementCount");
/**
 * AbortControllerRegistry does not actually remember active event listeners: instead
 * it ties them to an AbortSignal and uses an AbortController to remove them all at once.
 */
class AbortControllerRegistry {
    constructor() {
        this[controllerSym] = new AbortController();
    }
    /**
     * Returns an options object for addEventListener that ties the listener
     * to the AbortSignal from the current AbortController.
     *
     * @param element - An HTML element
     * @param triggers - The list of active WML trigger events for the specified elements
     */
    set(element, triggers) {
        return { signal: this[controllerSym].signal };
    }
    /**
     * Removes all registered event listeners and resets the registry.
     */
    reset() {
        this[controllerSym].abort();
        this[controllerSym] = new AbortController();
    }
}
/**
 * WeakMapRegistry maps active trigger events to each DOM element through a WeakMap.
 * This ensures that the mapping remains private to this module, while still allowing garbage
 * collection of the involved elements.
 */
class WeakMapRegistry {
    constructor() {
        this[triggerMapSym] = new WeakMap();
        this[elementCountSym] = 0;
    }
    /**
     * Sets active triggers for the specified element.
     *
     * @param element - An HTML element
     * @param triggers - The list of active WML trigger events for the specified element
     */
    set(element, triggers) {
        if (!this[triggerMapSym].has(element)) {
            this[elementCountSym]++;
        }
        this[triggerMapSym].set(element, triggers);
        return {};
    }
    /**
     * Removes all registered event listeners.
     */
    reset() {
        if (this[elementCountSym] <= 0)
            return;
        for (const element of document.body.querySelectorAll('*')) {
            if (this[elementCountSym] <= 0)
                break;
            const triggers = this[triggerMapSym].get(element);
            if (triggers != null) {
                this[elementCountSym]--;
            }
            for (const trigger of triggers || [])
                element.removeEventListener(trigger, onWMLTriggered);
        }
        this[triggerMapSym] = new WeakMap();
        this[elementCountSym] = 0;
    }
}
const triggerRegistry = canAbortListeners() ? new AbortControllerRegistry() : new WeakMapRegistry();
/**
 * Adds event listeners to the specified element.
 */
function addWMLListeners(element) {
    const triggerRegExp = /\S+/g;
    const triggerAttr = (element.getAttribute('wml-trigger') || element.getAttribute('data-wml-trigger') || "click");
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
 */
export function Enable() {
    whenReady(Reload);
}
/**
 * Reloads the WML page by adding necessary event listeners and browser listeners.
 */
export function Reload() {
    triggerRegistry.reset();
    document.body.querySelectorAll('[wml-event], [wml-window], [wml-openurl], [data-wml-event], [data-wml-window], [data-wml-openurl]').forEach(addWMLListeners);
}
