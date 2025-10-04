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
    return function (method: number, args: any = null) {
        return runtimeCallWithID(object, method, windowName, args);
    };
}

async function runtimeCallWithID(objectID: number, method: number, windowName: string, args: any): Promise<any> {
    let url = new URL(runtimeURL);
    url.searchParams.append("object", objectID.toString());
    url.searchParams.append("method", method.toString());
    if (args) { url.searchParams.append("args", JSON.stringify(args)); }

    let headers: Record<string, string> = {
        ["x-wails-client-id"]: clientId
    }
    if (windowName) {
        headers["x-wails-window-name"] = windowName;
    }

    let response = await fetch(url, { headers });
    if (!response.ok) {
        throw new Error(await response.text());
    }

    const contentType = response.headers.get("Content-Type") || "";

    // Handle different content types
    if (contentType.indexOf("application/json") !== -1) {
        return response.json();
    } else if (contentType.indexOf("image/") !== -1) {
        // Return image as Blob with data URL
        const blob = await response.blob();
        return new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.onloadend = () => resolve(reader.result);
            reader.onerror = reject;
            reader.readAsDataURL(blob);
        });
    } else if (contentType.indexOf("application/octet-stream") !== -1) {
        // Return binary data as ArrayBuffer
        return response.arrayBuffer();
    } else {
        return response.text();
    }
}
