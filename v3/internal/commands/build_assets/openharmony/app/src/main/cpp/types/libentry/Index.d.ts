/*
 * Type declarations for the native Wails bridge module (libentry.so).
 *
 * The C++ layer registers the Go backend's exported functions as a NAPI
 * module; ArkTS imports them by name.
 */

export interface Asset {
  body: ArrayBuffer;
  mime: string;
}

export type Dispatcher = (id: number, name: string, payload: string) => void;

export const open: (callback: Dispatcher) => boolean;
export const result: (id: number, value: string) => boolean;
export const close: () => boolean;
export const start: () => boolean;
export const ready: () => boolean;
export const asset: (path: string) => Asset | null;
export const request: (id: string, payload: string) => boolean;
export const message: (message: string) => void;
export const runOnMainThread: (id: number) => boolean;
