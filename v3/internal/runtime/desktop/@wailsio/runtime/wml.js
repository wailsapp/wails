
import {Emit, WailsEvent} from "./events";
import {Question} from "./dialogs";
import {WindowMethods, Get} from "./window";
import {OpenURL} from "./browser";

/**
 * Sends an event with the given name and optional data.
 *
 * @param {string} eventName - The name of the event to send.
 * @param {any} [data=null] - Optional data to send along with the event.
 *
 * @return {void}
 */
function sendEvent(eventName, data=null) {
    let event = new WailsEvent(eventName, data);
    Emit(event);
}

/**
 * Adds event listeners to elements with `wml-event` attribute.
 *
 * @return {void}
 */
function addWMLEventListeners() {
    const elements = document.querySelectorAll('[wml-event]');
    elements.forEach(function (element) {
        const eventType = element.getAttribute('wml-event');
        const confirm = element.getAttribute('wml-confirm');
        const trigger = element.getAttribute('wml-trigger') || "click";

        let callback = function () {
            if (confirm) {
                Question({Title: "Confirm", Message:confirm, Detached: false, Buttons:[{Label:"Yes"},{Label:"No", IsDefault:true}]}).then(function (result) {
                    if (result !== "No") {
                        sendEvent(eventType);
                    }
                });
                return;
            }
            sendEvent(eventType);
        };

        // Remove existing listeners
        element.removeEventListener(trigger, callback);

        // Add new listener
        element.addEventListener(trigger, callback);
    });
}


/**
 * Calls a method on a specified window.
 * @param {string} windowName - The name of the window to call the method on.
 * @param {string} method - The name of the method to call.
 */
function callWindowMethod(windowName, method) {
    let targetWindow = Get(windowName);
    let methodMap = WindowMethods(targetWindow);
    if (!methodMap.has(method)) {
        console.log("Window method " + method + " not found");
    }
    try {
        methodMap.get(method)();
    } catch (e) {
        console.error("Error calling window method '" + method + "': " + e);
    }
}

/**
 * Adds window listeners for elements with the 'wml-window' attribute.
 * Removes any existing listeners before adding new ones.
 *
 * @return {void}
 */
function addWMLWindowListeners() {
    const elements = document.querySelectorAll('[wml-window]');
    elements.forEach(function (element) {
        const windowMethod = element.getAttribute('wml-window');
        const confirm = element.getAttribute('wml-confirm');
        const trigger = element.getAttribute('wml-trigger') || 'click';
        const targetWindow = element.getAttribute('wml-target-window') || '';

        let callback = function () {
            if (confirm) {
                Question({Title: "Confirm", Message:confirm, Buttons:[{Label:"Yes"},{Label:"No", IsDefault:true}]}).then(function (result) {
                    if (result !== "No") {
                        callWindowMethod(targetWindow, windowMethod);
                    }
                });
                return;
            }
            callWindowMethod(targetWindow, windowMethod);
        };

        // Remove existing listeners
        element.removeEventListener(trigger, callback);

        // Add new listener
        element.addEventListener(trigger, callback);
    });
}

/**
 * Adds a listener to elements with the 'wml-openurl' attribute.
 * When the specified trigger event is fired on any of these elements,
 * the listener will open the URL specified by the 'wml-openurl' attribute.
 * If a 'wml-confirm' attribute is provided, a confirmation dialog will be displayed,
 * and the URL will only be opened if the user confirms.
 *
 * @return {void}
 */
function addWMLOpenBrowserListener() {
    const elements = document.querySelectorAll('[wml-openurl]');
    elements.forEach(function (element) {
        const url = element.getAttribute('wml-openurl');
        const confirm = element.getAttribute('wml-confirm');
        const trigger = element.getAttribute('wml-trigger') || "click";

        let callback = function () {
            if (confirm) {
                Question({Title: "Confirm", Message:confirm, Buttons:[{Label:"Yes"},{Label:"No", IsDefault:true}]}).then(function (result) {
                    if (result !== "No") {
                        void OpenURL(url);
                    }
                });
                return;
            }
            void OpenURL(url);
        };

        // Remove existing listeners
        element.removeEventListener(trigger, callback);

        // Add new listener
        element.addEventListener(trigger, callback);
    });
}

/**
 * Reloads the WML page by adding necessary event listeners and browser listeners.
 *
 * @return {void}
 */
export function reloadWML() {
    console.log("Reloading WML");
    addWMLEventListeners();
    addWMLWindowListeners();
    addWMLOpenBrowserListener();
}
