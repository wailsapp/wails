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
export function newRuntimeCaller(object, windowName = '') {
    return function (method, args = null) {
        return runtimeCallWithID(object, method, windowName, args);
    };
}
async function runtimeCallWithID(objectID, method, windowName, args) {
    var _a, _b;
    let url = new URL(runtimeURL);
    url.searchParams.append("object", objectID.toString());
    url.searchParams.append("method", method.toString());
    if (args) {
        url.searchParams.append("args", JSON.stringify(args));
    }
    let headers = {
        ["x-wails-client-id"]: clientId
    };
    if (windowName) {
        headers["x-wails-window-name"] = windowName;
    }
    let response = await fetch(url, { headers });
    if (!response.ok) {
        throw new Error(await response.text());
    }
    if (((_b = (_a = response.headers.get("Content-Type")) === null || _a === void 0 ? void 0 : _a.indexOf("application/json")) !== null && _b !== void 0 ? _b : -1) !== -1) {
        return response.json();
    }
    else {
        return response.text();
    }
}
