/*
 _     __     _ __
| |  / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

import { nanoid } from "./nanoid.js";

const runtimeURL = window.location.origin + "/wails/runtime";

// Re-export nanoid for custom transport implementations
export { nanoid };

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
 * RuntimeTransport defines the interface for custom IPC transport implementations.
 * Implement this interface to use WebSockets, custom protocols, or any other
 * transport mechanism instead of the default HTTP fetch.
 */
export interface RuntimeTransport {
    /**
     * Send a runtime call and return the response.
     *
     * @param objectID - The Wails object ID (0=Call, 1=Clipboard, etc.)
     * @param method - The method ID to call
     * @param windowName - Optional window name
     * @param args - Arguments to pass (will be JSON stringified if present)
     * @returns Promise that resolves with the response data
     */
    call(objectID: number, method: number, windowName: string, args: any): Promise<any>;
}

/**
 * Custom transport implementation (can be set by user)
 */
let customTransport: RuntimeTransport | null = null;

/**
 * Set a custom transport for all Wails runtime calls.
 * This allows you to replace the default HTTP fetch transport with
 * WebSockets, custom protocols, or any other mechanism.
 *
 * @param transport - Your custom transport implementation
 *
 * @example
 * ```typescript
 * import { setTransport } from '/wails/runtime.js';
 *
 * const wsTransport = {
 *   call: async (objectID, method, windowName, args) => {
 *     // Your WebSocket implementation
 *   }
 * };
 *
 * setTransport(wsTransport);
 * ```
 */
export function setTransport(transport: RuntimeTransport | null): void {
    customTransport = transport;
}

/**
 * Get the current transport (useful for extending/wrapping)
 */
export function getTransport(): RuntimeTransport | null {
    return customTransport;
}

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
    // Use custom transport if available
    if (customTransport) {
        return customTransport.call(objectID, method, windowName, args);
    }

    // Default HTTP fetch transport
    let url = new URL(runtimeURL);

    let body: { object: number; method: number, args?: any } = {
      object: objectID,
      method
    }
    if (args !== null && args !== undefined) {
      body.args = args;
    }

    let headers: Record<string, string> = {
        ["x-wails-client-id"]: clientId,
        ["Content-Type"]: "application/json"
    }
    if (windowName) {
        headers["x-wails-window-name"] = windowName;
    }

    let response = await fetch(url, {
      method: 'POST',
      headers,
      body: JSON.stringify(body)
    });
    if (!response.ok) {
        throw new Error(await response.text());
    }

    if ((response.headers.get("Content-Type")?.indexOf("application/json") ?? -1) !== -1) {
        return response.json();
    } else {
        return response.text();
    }
}
