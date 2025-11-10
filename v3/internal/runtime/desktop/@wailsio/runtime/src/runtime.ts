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
    if (args) {
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

/**
 * Helper utilities for custom transport implementations
 */

/**
 * Generates a unique message ID for transport requests.
 * Uses the same nanoid implementation as the Wails runtime.
 *
 * @returns A unique string identifier (21 characters)
 *
 * @example
 * ```typescript
 * import { generateMessageID } from '/wails/runtime.js';
 * const msgID = generateMessageID();
 * ```
 */
export function generateMessageID(): string {
    return nanoid();
}

/**
 * Builds a transport request message in the format expected by Wails.
 * Handles JSON stringification of arguments and proper field naming.
 *
 * @param id - Unique message ID (use generateMessageID())
 * @param objectID - Wails object ID (0=Call, 1=Clipboard, etc. - see objectNames)
 * @param method - Method ID within the object
 * @param windowName - Source window name (optional)
 * @param args - Method arguments (will be JSON stringified)
 * @returns Formatted transport request message
 *
 * @example
 * ```typescript
 * import { buildTransportRequest, generateMessageID, objectNames } from '/wails/runtime.js';
 *
 * const message = buildTransportRequest(
 *     generateMessageID(),
 *     objectNames.Call,
 *     0,
 *     '',
 *     { methodName: 'Greet', args: ['World'] }
 * );
 * ```
 */
export function buildTransportRequest(
    id: string,
    objectID: number,
    method: number,
    windowName: string,
    args: any
): any {
    return {
        id: id,
        type: 'request',
        request: {
            object: objectID,
            method: method,
            args: args ? JSON.stringify(args) : undefined,
            windowName: windowName || undefined,
            clientId: clientId
        }
    };
}

/**
 * Handles Wails callback invocation for binding calls.
 * This abstracts the internal mechanism of how Wails processes method call results.
 *
 * For binding calls (object=0, method=0), Wails expects the result to be delivered
 * via window._wails.callResultHandler() rather than resolving the transport promise directly.
 *
 * @param pending - Pending request object with stored args
 * @param responseData - Decoded response data (string)
 * @param contentType - Response content type
 * @returns true if handled as binding call, false otherwise
 *
 * @example
 * ```typescript
 * import { handleWailsCallback } from '/wails/runtime.js';
 *
 * // In your transport's message handler:
 * if (handleWailsCallback(pending, responseData, response.contentType)) {
 *     pending.resolve(); // Wails callback handled it
 * } else {
 *     pending.resolve(responseData); // Direct resolve for non-binding calls
 * }
 * ```
 */
export function handleWailsCallback(
    pending: any,
    responseData: string,
    contentType: string
): boolean {
    // Only for JSON responses (binding calls)
    if (!responseData || !contentType?.includes('application/json')) {
        return false;
    }

    // Extract call-id from stored request
    const callId = pending.request?.args?.['call-id'];

    // Invoke Wails callback handler if available
    if (callId && window._wails?.callResultHandler) {
        window._wails.callResultHandler(callId, responseData, true);
        return true;
    }

    return false;
}

/**
 * Dispatches an event to the Wails event system.
 * Used by custom transports to deliver server-sent events to frontend listeners.
 *
 * @param event - Event object with name, data, and optional sender
 *
 * @example
 * ```typescript
 * import { dispatchWailsEvent } from '/wails/runtime.js';
 *
 * // In your WebSocket transport's message handler:
 * if (msg.type === 'event') {
 *     dispatchWailsEvent(msg.event);
 * }
 * ```
 */
export function dispatchWailsEvent(event: any): void {
    if (window._wails?.dispatchWailsEvent) {
        window._wails.dispatchWailsEvent(event);
    }
}
