export declare const eventListeners: Map<string, Listener[]>;
export declare class Listener {
    eventName: string;
    callback: (data: any) => void;
    maxCallbacks: number;
    constructor(eventName: string, callback: (data: any) => void, maxCallbacks: number);
    dispatch(data: any): boolean;
}
export declare function listenerOff(listener: Listener): void;
