export = wailsapp__runtime;

declare const wailsapp__runtime: {
    Browser: {
        OpenFile(filename: string): Promise<any>;
        OpenURL(url: string): Promise<any>;
    };
    Events: {
        Acknowledge(eventName: string): void; 
        Emit(eventName: string): void;
        Heartbeat(eventName: string, timeInMilliseconds: number, callback: () => void): void;
        On(eventName: string, callback: () => void): void;
        OnMultiple(eventName: string, callback: () => void, maxCallbacks: number): void;
        Once(eventName: string, callback: () => void): void;
    };
    Init(callback: () => void): void;
    Log: {
        Debug(message: string): void;
        Error(message: string): void;
        Fatal(message: string): void;
        Info(message: string): void;
        Warning(message: string): void;
    };
};


