/*
 _     __     _ __
| |  / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

/* jshint esversion: 9 */
import { nanoid } from 'nanoid/non-secure';

const runtimeURL = window.location.origin + "/wails/runtime";

// Object Names
export const objectNames = {
    Call: 0,
    Clipboard: 1,
    Application: 2,
    Events: 3,
    ContextMenu: 4,
    Dialog: 5,
    Window: 6,
    Screens: 7,
    System: 8,
    Browser: 9,
    CancelCall: 10,
}
export let clientId = nanoid();

/**
 * Creates a runtime caller function that invokes a specified method on a given object within a specified window context.
 *
 * @param {Object} object - The object on which the method is to be invoked.
 * @param {string} windowName - The name of the window context in which the method should be called.
 * @returns {Function} A runtime caller function that takes the method name and optionally arguments and invokes the method within the specified window context.
 */
export function newRuntimeCaller(object, windowName) {
    return function (method, args=null) {
        return runtimeCall(object + "." + method, windowName, args);
    };
}

/**
 * Creates a new runtime caller with specified ID.
 *
 * @param {object} object - The object to invoke the method on.
 * @param {string} windowName - The name of the window.
 * @return {Function} - The new runtime caller function.
 */
export function newRuntimeCallerWithID(object, windowName) {
    return function (method, args=null) {
        return runtimeCallWithID(object, method, windowName, args);
    };
}


function runtimeCall(method, windowName, args) {
    let url = new URL(runtimeURL);
    if( method ) {
        url.searchParams.append("method", method);
    }
    let fetchOptions = {
        headers: {},
    };
    if (windowName) {
        fetchOptions.headers["x-wails-window-name"] = windowName;
    }
    if (args) {
        url.searchParams.append("args", JSON.stringify(args));
    }
    fetchOptions.headers["x-wails-client-id"] = clientId;

    return new Promise((resolve, reject) => {
        fetch(url, fetchOptions)
            .then(response => {
                if (response.ok) {
                    // check content type
                    if (response.headers.get("Content-Type") && response.headers.get("Content-Type").indexOf("application/json") !== -1) {
                        return response.json();
                    } else {
                        return response.text();
                    }
                }
                reject(Error(response.statusText));
            })
            .then(data => resolve(data))
            .catch(error => reject(error));
    });
}

function runtimeCallWithID(objectID, method, windowName, args) {
    let url = new URL(runtimeURL);
    url.searchParams.append("object", objectID);
    url.searchParams.append("method", method);
    let fetchOptions = {
        headers: {},
    };
    if (windowName) {
        fetchOptions.headers["x-wails-window-name"] = windowName;
    }
    if (args) {
        url.searchParams.append("args", JSON.stringify(args));
    }
    fetchOptions.headers["x-wails-client-id"] = clientId;
    return new Promise((resolve, reject) => {
        fetch(url, fetchOptions)
            .then(response => {
                if (response.ok) {
                    // check content type
                    if (response.headers.get("Content-Type") && response.headers.get("Content-Type").indexOf("application/json") !== -1) {
                        return response.json();
                    } else {
                        return response.text();
                    }
                }
                reject(Error(response.statusText));
            })
            .then(data => resolve(data))
            .catch(error => reject(error));
    });
}
