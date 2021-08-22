export = wailsapp__runtime;

interface Level {
    TRACE: 1,
    DEBUG: 2,
    INFO: 3,
    WARNING: 4,
    ERROR: 5,
}

declare const wailsapp__runtime: {
    EventsEmit(eventName: string, data?: any): void;
    EventsOn(eventName: string, callback: (data?: any) => void): void;
    EventsOnMultiple(eventName: string, callback: (data?: any) => void, maxCallbacks: number): void;
    EventsOnce(eventName: string, callback: (data?: any) => void): void;
    // Init(callback: () => void): void;
    LogDebug(message: string): void;
    LogError(message: string): void;
    LogFatal(message: string): void;
    LogInfo(message: string): void;
    LogWarning(message: string): void;
    LogLevel: Level;
};


