export = wailsapp__runtime;

interface Store {
    get(): any;
    set(value: any): void;
    subscribe(callback: (newvalue: any) => void): void;
    update(callback: (currentvalue: any) => any): void;
}

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
    System: {
        DarkModeEnabled(): Promise<boolean>;
        OnThemeChange(callback: (darkModeEnabled: boolean) => void): void;
        LogLevel(): Store;
    };
    Store: {
        New(name: string, defaultValue?: any): Store;
    }
};


