export interface Position {
    x: number;
    y: number;
}

export interface Size {
    w: number;
    h: number;
}


export function EventsEmit(eventName: string, data?: any): void;

export function EventsOn(eventName: string, callback: (data?: any) => void): void;

export function EventsOnMultiple(eventName: string, callback: (data?: any) => void, maxCallbacks: number): void;

export function EventsOnce(eventName: string, callback: (data?: any) => void): void;

export function EventsOff(eventName: string): void;

export function LogTrace(message: string): void;

export function LogDebug(message: string): void;

export function LogError(message: string): void;

export function LogFatal(message: string): void;

export function LogInfo(message: string): void;

export function LogWarning(message: string): void;

export function WindowReload(): void;

export function WindowSetSystemDefaultTheme(): void;

export function WindowSetLightTheme(): void;

export function WindowSetDarkTheme(): void;

export function WindowCenter(): void;

export function WindowSetTitle(title: string): void;

export function WindowFullscreen(): void;

export function WindowUnfullscreen(): void;

export function WindowSetSize(width: number, height: number): Promise<Size>;

export function WindowGetSize(): Promise<Size>;

export function WindowSetMaxSize(width: number, height: number): void;

export function WindowSetMinSize(width: number, height: number): void;

export function WindowSetPosition(x: number, y: number): void;

export function WindowGetPosition(): Promise<Position>;

export function WindowHide(): void;

export function WindowShow(): void;

export function WindowMaximise(): void;

export function WindowToggleMaximise(): void;

export function WindowUnmaximise(): void;

export function WindowMinimise(): void;

export function WindowUnminimise(): void;

export function WindowSetRGBA(R: number, G: number, B: number, A: number): void;

export function BrowserOpenURL(url: string): void;

export function Quit(): void;
