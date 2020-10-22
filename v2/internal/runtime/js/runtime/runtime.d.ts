export = wailsapp__runtime;

interface Store {
    get(): any;
    set(value: any): void;
    subscribe(callback: (newvalue: any) => void): void;
    update(callback: (currentvalue: any) => any): void;
}

interface Level {
	TRACE: 1,
	DEBUG: 2,
	INFO: 3,
	WARNING: 4,
	ERROR: 5,
};

declare const wailsapp__runtime: {
    Browser: {
        Open(target: string): Promise<any>;
    };
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
    System: {
        DarkModeEnabled(): Promise<boolean>;
        OnThemeChange(callback: (darkModeEnabled: boolean) => void): void;
        LogLevel(): Store;
        Platform: string;
        AppType: string
    };
    Store: {
        New(name: string, defaultValue?: any): Store;
    }
};


