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

export const abilityStarted: () => void;
export const abilityResumed: () => void;
export const abilityPaused: () => void;
export const abilityStopped: () => void;
export const lowMemory: () => void;
export const emitSystemEvent: (name: string, payload: string) => void;
export const emitEvent: (name: string, payload: string) => void;
export const pageStarted: (url: string) => void;
export const pageFinished: (url: string) => void;
export const shutdown: () => void;

export const dialogResult: (id: number, buttonIndex: number) => void;
