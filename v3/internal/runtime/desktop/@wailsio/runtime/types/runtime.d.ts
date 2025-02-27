export declare const objectNames: Readonly<{
    Call: 0;
    Clipboard: 1;
    Application: 2;
    Events: 3;
    ContextMenu: 4;
    Dialog: 5;
    Window: 6;
    Screens: 7;
    System: 8;
    Browser: 9;
    CancelCall: 10;
}>;
export declare let clientId: string;
/**
 * Creates a new runtime caller with specified ID.
 *
 * @param object - The object to invoke the method on.
 * @param windowName - The name of the window.
 * @return The new runtime caller function.
 */
export declare function newRuntimeCaller(object: number, windowName?: string): (method: number, args?: any) => Promise<any>;
