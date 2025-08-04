/*
 _     __     _ __
| |  / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

import { nanoid } from './nanoid.js';

const runtimeURL = window.location.origin + "/wails/runtime";

// Object Names
export const objectNames = Object.freeze({
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
});
export let clientId = nanoid();

/**
 * Creates a new runtime caller with specified ID.
 *
 * @param object - The object to invoke the method on.
 * @param windowName - The name of the window.
 * @return The new runtime caller function.
 */
export function newRuntimeCaller(object: number, windowName: string = '') {
    return function (method: number, args: any = null, options: RequestInit = {}) {
        return runtimeCallWithID(object, method, windowName, args, options);
    };
}

async function runtimeCallWithID(
    objectID: number, 
    method: number, 
    windowName: string, 
    args: any,
    options: RequestInit = {}
): Promise<Response> {
    let url = new URL(runtimeURL);
    url.searchParams.append("object", objectID.toString());
    url.searchParams.append("method", method.toString());
    if (args) { 
        url.searchParams.append("args", JSON.stringify(args)); 
    }

    let headers: Record<string, string> = {
        ["x-wails-client-id"]: clientId
    }
    if (windowName) {
        headers["x-wails-window-name"] = windowName;
    }

    // Merge headers with provided options
    const requestOptions: RequestInit = {
        ...options,
        headers: {
            ...headers,
            ...(options.headers || {})
        }
    };

    let response = await fetch(url, requestOptions);
    
    // Don't automatically throw on !response.ok - let caller handle it
    // This allows proper error handling with structured error responses
    
    return response; // Return response object for flexible handling
}
