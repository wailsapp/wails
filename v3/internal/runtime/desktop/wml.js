
import {Emit, WailsEvent} from "./events";
import {Question} from "./dialogs";

function sendEvent(eventName, data=null) {
    let event = new WailsEvent(eventName, data);
    Emit(event);
}

function addWMLEventListeners() {
    const elements = document.querySelectorAll('[data-wml-event]');
    elements.forEach(function (element) {
        const eventType = element.getAttribute('data-wml-event');
        const confirm = element.getAttribute('data-wml-confirm');
        const trigger = element.getAttribute('data-wml-trigger') || "click";

        let callback = function () {
            if (confirm) {
                Question({Title: "Confirm", Message:confirm, Buttons:[{Label:"Yes"},{Label:"No", IsDefault:true}]}).then(function (result) {
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

function callWindowMethod(method) {
    if (wails.Window[method] === undefined) {
        console.log("Window method " + method + " not found");
    }
    wails.Window[method]();
}

function addWMLWindowListeners() {
    const elements = document.querySelectorAll('[data-wml-window]');
    elements.forEach(function (element) {
        const windowMethod = element.getAttribute('data-wml-window');
        const confirm = element.getAttribute('data-wml-confirm');
        const trigger = element.getAttribute('data-wml-trigger') || "click";

        let callback = function () {
            if (confirm) {
                Question({Title: "Confirm", Message:confirm, Buttons:[{Label:"Yes"},{Label:"No", IsDefault:true}]}).then(function (result) {
                    if (result !== "No") {
                        callWindowMethod(windowMethod);
                    }
                });
                return;
            }
            callWindowMethod(windowMethod);
        };

        // Remove existing listeners
        element.removeEventListener(trigger, callback);

        // Add new listener
        element.addEventListener(trigger, callback);
    });
}

export function reloadWML() {
    addWMLEventListeners();
    addWMLWindowListeners();
}
