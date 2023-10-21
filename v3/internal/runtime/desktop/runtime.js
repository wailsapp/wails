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
}
export let clientId = nanoid();

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

export function newRuntimeCaller(object, windowName) {
    return function (method, args=null) {
        return runtimeCall(object + "." + method, windowName, args);
    };
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

export function newRuntimeCallerWithID(object, windowName) {
    return function (method, args=null) {
        return runtimeCallWithID(object, method, windowName, args);
    };
}
