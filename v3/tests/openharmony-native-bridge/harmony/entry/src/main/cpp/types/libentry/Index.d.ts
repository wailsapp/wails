export const open: (callback: (id: number, value: number) => void) => boolean;
export const request: (id: number, delay: number) => boolean;
export const close: () => boolean;
