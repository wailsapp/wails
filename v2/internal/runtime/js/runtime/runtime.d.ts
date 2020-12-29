export = wailsapp__runtime;

interface Store {
    get(): any;
    set(value: any): void;
    subscribe(callback: (newvalue: any) => void): void;
    update(callback: (currentvalue: any) => any): void;
}

interface MacTitleBar {
    TitleBarAppearsTransparent: boolean; //  NSWindow.titleBarAppearsTransparent
    HideTitle: boolean; // NSWindow.hideTitle
    HideTitleBar: boolean; // NSWindow.hideTitleBar
    FullSizeContent: boolean; // Makes the webview portion of the window the full size of the window, even over the titlebar
    UseToolbar: boolean; // Set true to add a blank toolbar to the window (makes the title bar larger)
    HideToolbarSeparator: boolean; // Set true to remove the separator between the toolbar and the main content area
}

interface MacAppConfig {
    TitleBar: MacTitleBar;
}
interface LinuxAppConfig {
}
interface WindowsAppConfig {
}

interface AppConfig {
    Title: string; // Application Title
    Width: number; // Window Width
    Height: number; // Window Height
    DisableResize: boolean; // True if resize is disabled
    Fullscreen: boolean; // App started in fullscreen
    MinWidth: number; // Window Minimum Width
    MinHeight: number; // Window Minimum Height
    MaxWidth: number; // Window Maximum Width
    MaxHeight: number; // Window Maximum Height
    StartHidden: boolean; // Start with window hidden
    DevTools: boolean; // Enables the window devtools
    RBGA: number; // The initial window colour. Convert to hex then it'll mean 0xRRGGBBAA
    Mac?: MacAppConfig; // - Configuration when running on Mac
    Linux?: LinuxAppConfig; // - Configuration when running on Linux
    Windows?: WindowsAppConfig; // - Configuration when running on Windows
    Appearance: string; // The default application appearance. Use the values listed here: https://developer.apple.com/documentation/appkit/nsappearance?language=objc
    WebviewIsTransparent: number; // Makes the background of the webview content transparent. Use this with the Alpha part of the window colour to make parts of your application transparent.
    WindowBackgroundIsTranslucent: number; // Makes the transparent parts of the application window translucent. Example: https://en.wikipedia.org/wiki/MacOS_Big_Sur#/media/File:MacOS_Big_Sur_-_Safari_Extensions_category_in_App_Store.jpg
    LogLevel: number; // The initial log level (lower is more verbose)
}
interface Level { 
	TRACE: 1,
	DEBUG: 2,
	INFO: 3,
	WARNING: 4,
	ERROR: 5,
}

interface OpenDialogOptions {
	DefaultDirectory:           string;
	DefaultFilename:            string;
	Title:                      string;
	Filters:                    string;
	AllowFiles:                 boolean;
	AllowDirectories:           boolean;
	AllowMultiple:              boolean;
	ShowHiddenFiles:            boolean;
	CanCreateDirectories:       boolean;
	ResolvesAliases:            boolean;
	TreatPackagesAsDirectories: boolean;
}

interface SaveDialogOptions {
    DefaultDirectory:           string;
    DefaultFilename:            string;
    Title:                      string;
    Filters:                    string;
    ShowHiddenFiles:            boolean;
    CanCreateDirectories:       boolean;
    TreatPackagesAsDirectories: boolean;
}

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
        Platform(): string;
        AppType(): string;
        AppConfig(): AppConfig;
    };
    Store: {
        New(name: string, defaultValue?: any): Store;
    };
    Dialog: {
        Open(options: OpenDialogOptions): Promise<Array<string>>;
        Save(options: SaveDialogOptions): Promise<Array<string>>;
    };
    Tray: {
        SetIcon(trayIconID: string): void;
    }
};


