interface runtime {
    EventsEmit(eventName: string, data?: any): void;

    EventsOn(eventName: string, callback: (data?: any) => void): void;

    EventsOnMultiple(eventName: string, callback: (data?: any) => void, maxCallbacks: number): void;

    EventsOnce(eventName: string, callback: (data?: any) => void): void;

    LogDebug(message: string): void;

    LogError(message: string): void;

    LogFatal(message: string): void;

    LogInfo(message: string): void;

    LogWarning(message: string): void;

    WindowReload(): void;
}

declare global {
    interface Window {
        runtime: runtime;
    }
}

export {};
