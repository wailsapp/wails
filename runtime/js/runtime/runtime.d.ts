export = wailsapp__runtime;

declare const wailsapp__runtime: {
    Browser: {
        OpenFile(filename: string): Promise<any>;
        OpenURL(url: string): Promise<any>;
    };
    Events: {
        Acknowledge(eventName: string): void;
        Emit(eventName: string, data?: any): void;
        Heartbeat(eventName: string, timeInMilliseconds: number, callback: (data?: any) => void): void;
        On(eventName: string, callback: (data?: any) => void): void;
        OnMultiple(eventName: string, callback: (data?: any) => void, maxCallbacks: number): void;
        Once(eventName: string, callback: (data?: any) => void): void;
    };
    Init(callback: () => void): void;
    Log: {
        Debug(message: string): void;
        Error(message: string): void;
        Fatal(message: string): void;
        Info(message: string): void;
        Warning(message: string): void;
    };
    Store: {
        New(name: string, optionalDefault?: any): any;
    };
};


