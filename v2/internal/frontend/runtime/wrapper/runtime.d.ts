export interface Position {
    x: number;
    y: number;
}

export interface Size {
    w: number;
    h: number;
}

export interface runtime {
    EventsEmit(eventName: string, data?: any): void;

    EventsOn(eventName: string, callback: (data?: any) => void): void;

    EventsOnMultiple(eventName: string, callback: (data?: any) => void, maxCallbacks: number): void;

    EventsOnce(eventName: string, callback: (data?: any) => void): void;

    LogTrace(message: string): void;

    LogDebug(message: string): void;

    LogError(message: string): void;

    LogFatal(message: string): void;

    LogInfo(message: string): void;

    LogWarning(message: string): void;

    WindowReload(): void;

    WindowCenter(): void;

    WindowSetTitle(title: string): void;

    WindowFullscreen(): void;

    WindowUnfullscreen(): void;

    WindowSetSize(width: number, height: number): Promise<Size>;

    WindowGetSize(): Promise<Size>;

    WindowSetMaxSize(width: number, height: number): void;

    WindowSetMinSize(width: number, height: number): void;

    WindowSetPosition(x: number, y: number): void;

    WindowGetPosition(): Promise<Position>;

    WindowHide(): void;

    WindowShow(): void;

    WindowMaximise(): void;

    WindowToggleMaximise(): void;

    WindowUnmaximise(): void;

    WindowMinimise(): void;

    WindowUnminimise(): void;

    WindowSetRGBA(R: number, G: number, B: number, A: number): void;

    BrowserOpenURL(url: string): void;

    Quit(): void;
}

declare global {
    interface Window {
        runtime: runtime;
    }
}
