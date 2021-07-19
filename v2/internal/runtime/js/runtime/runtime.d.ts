export = wailsapp__runtime;

interface Level {
    TRACE: 1,
    DEBUG: 2,
    INFO: 3,
    WARNING: 4,
    ERROR: 5,
}

declare const wailsapp__runtime: {
    Events: {
        Emit(eventName: string, data?: any): void;
        On(eventName: string, callback: (data?: any) => void): void;
        OnMultiple(eventName: string, callback: (data?: any) => void, maxCallbacks: number): void;
        Once(eventName: string, callback: (data?: any) => void): void;
    };
    // Init(callback: () => void): void;
    Log: {
        Debug(message: string): void;
        Error(message: string): void;
        Fatal(message: string): void;
        Info(message: string): void;
        Warning(message: string): void;
        Level: Level;
    };
};


