import { Stream } from "./stream.js";

/** Options for loading a bounded media file over a registered Wails stream. */
export interface SourceOptions {
    /** Maximum downloaded bytes per source. Defaults to 32 MiB. Must be a positive safe integer. */
    maxBytes?: number;
    /** Cancels this load without clearing a previously loaded source. */
    signal?: AbortSignal;
}

interface SourceState {
    generation: number;
    controller?: AbortController;
    objectURL?: string;
}

// The npm and bundled runtimes can coexist in one page. Store ownership on the
// element under a shared symbol so either copy can replace or clear its source.
const sourceKey = Symbol.for('wails.media.source');
type ManagedElement = HTMLMediaElement & { [sourceKey]?: SourceState };
const defaultMaxBytes = 32 * 1024 * 1024;

function checkAborted(signal: AbortSignal): void {
    if (signal.aborted) {
        throw signal.reason ?? new DOMException('Media load cancelled', 'AbortError');
    }
}

/**
 * Loads a clip from a Go media.NewHandler registered with app.HandleStream,
 * then assigns a blob URL to an audio or video element. name is a file path
 * relative to that handler's filesystem, not a URL or operating-system path.
 *
 * Desktop transfers use bounded Wails stream responses, including on WebView2,
 * without opening a listener. The complete file is retained before assignment;
 * this enables local playback on Linux but is not progressive media playback.
 * maxBytes defaults to 32 MiB, and the Go handler enforces its own limit too.
 * Browser decoding and blob construction can require additional memory.
 *
 * A new call cancels an earlier load for the same element. The existing source
 * remains until a replacement succeeds. Aborted loads reject with AbortError
 * (or the caller's abort reason). Completion means the source was assigned;
 * use media events and play() for decoding and playback results.
 *
 * Call ClearSource before discarding the element to release its blob. Use these
 * helpers consistently for an element; they do not observe direct src changes,
 * source children or DOM removal. This also works with detached Audio objects.
 */
export async function SetSource(element: HTMLMediaElement, streamName: string, name: string, options: SourceOptions = {}): Promise<void> {
    const maxBytes = options.maxBytes ?? defaultMaxBytes;
    if (!Number.isSafeInteger(maxBytes) || maxBytes <= 0) {
        throw new RangeError('maxBytes must be a positive safe integer');
    }
    if (!streamName.trim() || !name.trim()) {
        throw new TypeError('A media stream and file name are required');
    }
    if (options.signal) checkAborted(options.signal);

    const managed = element as ManagedElement;
    const state = managed[sourceKey] ?? { generation: 0 };
    managed[sourceKey] = state;
    state.controller?.abort();
    const generation = ++state.generation;
    const controller = new AbortController();
    state.controller = controller;
    const cancel = () => controller.abort(options.signal?.reason);
    options.signal?.addEventListener('abort', cancel, { once: true });
    try {
        const blob = await receiveBlob(streamName, name, maxBytes, controller.signal);
        checkAborted(controller.signal);
        if (state.generation !== generation) {
            throw new DOMException('Media source replaced', 'AbortError');
        }
        const objectURL = URL.createObjectURL(blob);
        try {
            element.src = objectURL;
        } catch (error) {
            URL.revokeObjectURL(objectURL);
            throw error;
        }
        if (state.objectURL) URL.revokeObjectURL(state.objectURL);
        state.objectURL = objectURL;
    } finally {
        options.signal?.removeEventListener('abort', cancel);
        if (state.generation === generation) state.controller = undefined;
    }
}

/**
 * Cancels an outstanding SetSource load, resets the media element and revokes
 * the blob it owns. Call this on component unmount or before discarding an Audio
 * object. Existing source children, if any, follow normal browser selection rules
 * when the src attribute is removed. Repeated calls have no effect.
 */
export function ClearSource(element: HTMLMediaElement): void {
    const managed = element as ManagedElement;
    const state = managed[sourceKey];
    if (!state) return;
    ++state.generation;
    state.controller?.abort();
    delete managed[sourceKey];
    if (state.objectURL) {
        try {
            element.pause();
            element.removeAttribute('src');
            element.load();
        } finally {
            URL.revokeObjectURL(state.objectURL);
        }
    }
}

function receiveBlob(streamName: string, name: string, maxBytes: number, signal: AbortSignal): Promise<Blob> {
    return new Promise((resolve, reject) => {
        const stream = Stream(streamName);
        stream.binaryType = 'arraybuffer';
        const parts: BlobPart[] = [];
        let expected: number | undefined;
        let contentType = '';
        let received = 0;
        let finished = false;

        const finish = (error?: unknown) => {
            if (finished) return;
            finished = true;
            signal.removeEventListener('abort', cancel);
            stream.onopen = stream.onmessage = stream.onerror = stream.onclose = null;
            stream.close();
            if (error !== undefined) {
                parts.length = 0;
                reject(error);
            } else {
                try {
                    resolve(new Blob(parts, { type: contentType }));
                } catch (cause) {
                    reject(cause);
                } finally {
                    parts.length = 0;
                }
            }
        };
        const cancel = () => finish(signal.reason ?? new DOMException('Media load cancelled', 'AbortError'));
        signal.addEventListener('abort', cancel, { once: true });
        if (signal.aborted) { cancel(); return; }

        stream.onopen = () => {
            if (finished) return;
            try { stream.send(JSON.stringify({ version: 1, name, maxBytes })); }
            catch (error) { finish(error); }
        };
        stream.onerror = () => finish(new Error('Media stream failed'));
        stream.onmessage = event => {
            if (finished) return;
            try {
                if (!(event.data instanceof ArrayBuffer)) throw new Error('Invalid media frame');
                if (expected === undefined) {
                    if (event.data.byteLength > 4096) throw new Error('Invalid media metadata');
                    const meta = JSON.parse(new TextDecoder().decode(event.data));
                    if (meta.version !== 1) throw new Error('Unsupported media protocol');
                    if (meta.error) throw meta.limit ? new RangeError(meta.error) : new Error(meta.error);
                    if (!Number.isSafeInteger(meta.size) || meta.size < 0 || typeof meta.type !== 'string') {
                        throw new Error('Invalid media metadata');
                    }
                    if (meta.size > maxBytes) throw new RangeError(`Media exceeds the ${maxBytes} byte limit`);
                    expected = meta.size;
                    contentType = meta.type;
                } else {
                    received += event.data.byteLength;
                    if (event.data.byteLength === 0 || event.data.byteLength > 64 * 1024 || received > expected || received > maxBytes) {
                        throw new RangeError('Media stream exceeded its declared limits');
                    }
                    parts.push(event.data);
                }
            } catch (error) { finish(error); }
        };
        stream.onclose = event => {
            if (event.code !== 1000 || expected === undefined || received !== expected) {
                finish(new Error('Media stream closed before the complete file arrived'));
            } else {
                finish();
            }
        };
    });
}
