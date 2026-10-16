/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

/**
 * GoStream — the WebSocket programming model without a listening socket.
 *
 * A WebSocket cannot be spoken over a custom URL scheme, so getting one inside
 * a webview means binding a real TCP port. This speaks the same shape over the
 * asset server the app already serves: Go→JS over one held poll per page, JS→Go
 * over a normal POST.
 *
 * {@link Stream} returns synchronously with `readyState === CONNECTING`, the way
 * `new WebSocket(url)` does, so a connection can be created at module scope:
 *
 * ```js
 * export const Telemetry = Stream("telemetry");   // connects in the background
 * ```
 */

import { nanoid } from "./nanoid.js";
import { hasDOM } from "./environment.js";

// Every piece of state that two runtime instances must agree on lives here, on
// window, not in module scope. The runtime can be evaluated more than once in a
// page — the platform injects a copy and an app may import it too — and a
// per-module registry meant both instances allocated connection id 1, ran
// competing polls against one session, and discarded frames addressed to ids
// only the other instance knew about. The session id alone was not enough.
interface StreamGlobals {
    session: string;
    generation: number;
    nextConnID: number;
    connections: Map<number, WailsSocket>;
    polling: boolean;
    pollAbort?: AbortController;
}

const GENERATION_KEY = "__wails_stream_generation";
const GENERATION_NAME_MARKER = "\u001f__wails_stream_generation__:";

// window.name belongs to the browsing context rather than the Document, so it
// survives a reload even when Web Storage is disabled. Preserve any name the
// application assigned and reserve only a suffix for the Stream generation.
function namedPageGeneration(): number {
    try {
        const marker = window.name.lastIndexOf(GENERATION_NAME_MARKER);
        if (marker < 0) return 0;
        const value = Number(window.name.slice(marker + GENERATION_NAME_MARKER.length));
        return Number.isSafeInteger(value) && value > 0 ? value : 0;
    } catch {
        return 0;
    }
}

function rememberNamedPageGeneration(generation: number): void {
    try {
        const marker = window.name.lastIndexOf(GENERATION_NAME_MARKER);
        const applicationName = marker < 0 ? window.name : window.name.slice(0, marker);
        window.name = applicationName + GENERATION_NAME_MARKER + generation;
    } catch {
        // Some embedded browser policies can deny access to the browsing
        // context name. The time origin remains the last-resort fallback.
    }
}

// Request arrival order is not page order: a request started just before a
// navigation can reach Go after the replacement page's first request. Keep a
// monotonic counter in this top-level browsing context so Go can distinguish
// those page incarnations without trusting server scheduling order.
function nextPageGeneration(): number {
    if (!hasDOM) return 1;
    // Anchor the value to the page's creation time as well as the stored
    // counter. If application code clears sessionStorage, the next reload must
    // still be newer than the generation Go has already retired; restarting
    // from 1 would make every stream request receive 410 until the host exits.
    // Performance.timeOrigin is absent on older WebKit releases that Wails
    // still targets. Date.now keeps those engines on the protocol; the stored
    // counter disambiguates reloads that happen within the same millisecond.
    const origin = typeof performance !== "undefined" && Number.isFinite(performance.timeOrigin)
        ? performance.timeOrigin
        : Date.now();
    const pageTime = Math.max(1, Math.floor(origin * 1000));
    let current = namedPageGeneration();
    try {
        const raw = window.sessionStorage.getItem(GENERATION_KEY);
        const stored = raw === null ? 0 : Number(raw);
        if (Number.isSafeInteger(stored) && stored > current) current = stored;
    } catch {
        // Storage can be disabled by policy. window.name supplies the
        // browsing-context counter in that case.
    }

    const next = Number.isSafeInteger(current) && current >= 0 && current < Number.MAX_SAFE_INTEGER
        ? Math.max(current + 1, pageTime)
        : pageTime;
    try {
        window.sessionStorage.setItem(GENERATION_KEY, String(next));
    } catch {
        // The browsing-context fallback below is sufficient when storage is
        // unavailable.
    }
    rememberNamedPageGeneration(next);
    return next;
}

function streamGlobals(): StreamGlobals {
    const fresh = (): StreamGlobals => ({
        session: nanoid(),
        generation: nextPageGeneration(),
        nextConnID: 1,
        connections: new Map<number, WailsSocket>(),
        polling: false,
    });
    if (!hasDOM) return fresh();
    const w = (window as any)._wails = (window as any)._wails || {};
    return (w.__stream = w.__stream || fresh());
}

const G = streamGlobals();

const HDR_SESSION = "x-wails-stream-session";
const HDR_GENERATION = "x-wails-stream-generation";
const HDR_CONN = "x-wails-stream-conn";
const HDR_KIND = "x-wails-stream-kind";
const HDR_NAME = "x-wails-stream-name";
const HDR_CHUNK = "x-wails-stream-chunk";
const HDR_CHUNK_INDEX = "x-wails-stream-chunk-index";
const HDR_CHUNK_TOTAL = "x-wails-stream-chunk-total";
const HDR_BATCH = "x-wails-stream-batch";

const KIND_DATA = 0;
const KIND_OPEN = 1;
const KIND_CLOSE = 2;
const KIND_ERROR = 3;

// Matches the runtime's own limit: stay under WebView2's ~2MB request body
// buffering limit in WebResourceRequested.
const CHUNK_THRESHOLD = 512 * 1024;
const MAX_BATCH_FRAMES = 4096;

// Backoff floor and ceiling for a failing poll. This is the "reasonable lower
// limit" on polling — an error backoff, not a steady-state interval. In normal
// operation the poll re-issues immediately, because the server hold has already
// provided whatever delay was appropriate.
const BACKOFF_MIN = 250;
const BACKOFF_MAX = 5000;
const textEncoder = new TextEncoder();

function streamURL(endpoint: string): string {
    return window.location.origin + "/wails/stream/" + endpoint;
}


class StreamRequestCancelled extends Error {
    constructor() {
        super("stream request cancelled");
        this.name = "AbortError";
    }
}

function sleep(ms: number, signal?: AbortSignal): Promise<void> {
    if (!signal) return new Promise((resolve) => setTimeout(resolve, ms));
    if (signal.aborted) return Promise.reject(new StreamRequestCancelled());

    return new Promise((resolve, reject) => {
        const onAbort = () => {
            clearTimeout(timer);
            reject(new StreamRequestCancelled());
        };
        const timer = setTimeout(() => {
            signal.removeEventListener("abort", onAbort);
            resolve();
        }, ms);
        signal.addEventListener("abort", onAbort, { once: true });
    });
}

/**
 * A stream connection. Implements the useful subset of the `WebSocket`
 * interface, so the same application code works against a real WebSocket in
 * server builds.
 *
 * One deliberate divergence from the standard: `binaryType` defaults to
 * `"arraybuffer"` rather than `"blob"`, because stream frames are always
 * binary and a Blob would force an extra async hop to read every message.
 * Setting it to `"blob"` behaves as the standard describes.
 */
export class WailsSocket extends EventTarget {
    static readonly CONNECTING = 0;
    static readonly OPEN = 1;
    static readonly CLOSING = 2;
    static readonly CLOSED = 3;

    readonly CONNECTING = 0;
    readonly OPEN = 1;
    readonly CLOSING = 2;
    readonly CLOSED = 3;

    /** The stream name this connection was opened against. */
    readonly name: string;
    readonly url: string;

    /** Always empty: subprotocols and extensions are not negotiated. */
    readonly protocol = "";
    readonly extensions = "";

    binaryType: "arraybuffer" | "blob" = "arraybuffer";
    readyState: number = WailsSocket.CONNECTING;

    // Declared for the type; the constructor replaces each with an accessor
    // that registers a real listener. Calling them directly instead would make
    // any wrapper (see JSONStream) receive both the raw and the decoded event.
    onopen: ((ev: Event) => void) | null = null;
    onmessage: ((ev: MessageEvent) => void) | null = null;
    onclose: ((ev: CloseEvent) => void) | null = null;
    onerror: ((ev: Event) => void) | null = null;

    /**
     * @internal Applied to every inbound payload before dispatch. JSONStream
     * replaces it, so `onmessage` and `addEventListener("message", …)` see the
     * same decoded value — intercepting only `onmessage` made the two disagree.
     */
    _decode: (payload: ArrayBuffer) => unknown = (payload) => payload;

    /** @internal */ readonly _id: number;
    private _buffered = 0;
    // Sends are serialised per connection: concurrent fetch() POSTs do not
    // preserve order, and Go relies on send order being the order it observes.
    private _chain: Promise<void> = Promise.resolve();
    private _pending: Promise<Uint8Array>[] = [];
    private _flushing = false;
    private _openAbort = new AbortController();
    private _sendAbort = new AbortController();

    /** Bytes queued by send() that have not yet reached Go. */
    get bufferedAmount(): number {
        return this._buffered;
    }

    constructor(name: string) {
        super();
        this.name = name;
        this._id = G.nextConnID++;
        this.url = hasDOM ? streamURL("poll") + "#" + encodeURIComponent(name) : "";

        for (const type of ["open", "message", "close", "error"]) {
            defineHandlerProperty(this, type);
        }

        G.connections.set(this._id, this);

        // Fire and forget: the ack arrives as an open frame on the poll, which
        // is what moves readyState to OPEN.
        this._chain = this._chain.then(() =>
            postFrame(this._id, KIND_OPEN, new Uint8Array(0), name, this._openAbort.signal)
        ).catch((err) => {
            this._fail(err);
        });

        startPolling();
    }

    /**
     * Send a frame to Go. Accepts anything a WebSocket does; strings are
     * encoded as UTF-8, since Go receives every frame as a byte slice.
     */
    send(data: string | ArrayBufferLike | ArrayBufferView | Blob): void {
        if (this.readyState === WailsSocket.CONNECTING) {
            // Matches the WebSocket specification.
            throw new DOMException("Still in CONNECTING state.", "InvalidStateError");
        }
        if (this.readyState !== WailsSocket.OPEN) {
            return;
        }

        // Resolve now rather than inside the chain. Mutable binary inputs are
        // copied synchronously so send() snapshots their bytes like a native
        // WebSocket, even when a batch waits behind an in-flight request. A
        // Blob is immutable and is therefore safe to read asynchronously once.
        const immediate = toBytesSync(data);
        const snapshot = immediate
            ? Promise.resolve(immediate)
            : toBytes(data, this._sendAbort.signal);
        // close() may discard this promise before the flush has dequeued it.
        // Attach a rejection observer immediately so cancelling an asynchronous
        // Blob read cannot become an unhandled rejection; Promise.all still
        // observes the original rejection when the flush already owns it.
        void snapshot.catch(() => {});

        // A Blob cannot be converted synchronously, but Blob.size is available
        // now, so its bytes are still counted the moment send() returns. Waiting
        // for arrayBuffer() to resolve let a large Blob slip past code using
        // bufferedAmount for backpressure.
        const blobSize =
            !immediate && typeof Blob !== "undefined" && data instanceof Blob ? data.size : 0;
        if (blobSize) {
            this._buffered += blobSize;
        }

        // Account for the bytes now, not when the chain reaches them. Reporting
        // zero while frames wait behind an in-flight batch would tell an
        // application using bufferedAmount for backpressure that it may keep
        // sending.
        if (immediate) {
            this._buffered += immediate.byteLength;
        }
        this._pending.push(snapshot);

        // Queue rather than post. Whatever accumulates while a request is in
        // flight goes out together in the next one, so the per-request cost —
        // a scheme-handler round trip, about eleven cgo calls on macOS — is
        // divided across the batch. Under light load a batch is one frame and
        // this behaves exactly as before; the batching only appears when the
        // sender is outrunning the transport, which is when it is needed.
        if (this._flushing) return;
        this._flushing = true;

        this._chain = this._chain.then(async () => {
            try {
                while (this._pending.length > 0) {
                    // Bound the work and temporary promise array for one pass.
                    // postBatch applies the wire-size bound after resolution.
                    const queued = this._pending.splice(0, MAX_BATCH_FRAMES);
                    const batch = await Promise.all(queued);

                    // A remote close may arrive while a Blob is being read.
                    // The connection is already gone in that case, so do not
                    // issue a late POST for bytes the peer can no longer take.
                    if (this.readyState === WailsSocket.CLOSED) {
                        break;
                    }

                    const total = batch.reduce((n, b) => n + b.byteLength, 0);
                    try {
                        await postBatch(this._id, batch, this._sendAbort.signal);
                    } finally {
                        // _closed resets the counter immediately. A request
                        // already in flight may finish afterwards, and must
                        // not drive bufferedAmount below zero.
                        this._buffered = Math.max(0, this._buffered - total);
                    }
                }
            } finally {
                this._flushing = false;
            }
        }).catch((err) => {
            this._flushing = false;
            this._fail(err);
        });
    }

    /** Close the connection. Go's handler sees Receive fail and returns. */
    close(code: number = 1000, reason: string = ""): void {
        if (this.readyState === WailsSocket.CLOSING || this.readyState === WailsSocket.CLOSED) {
            return;
        }
        this.readyState = WailsSocket.CLOSING;
        this._pending.length = 0;
        this._buffered = 0;
        this._openAbort.abort();

        // A data frame waiting for receiver capacity must not prevent the
        // reserved close control from ever reaching Go. Cancel only the data
        // retry path; the chain below still orders the close after that path
        // has unwound, preserving normal send-before-close ordering.
        this._sendAbort.abort();
        this._chain = this._chain.then(() =>
            postFrame(this._id, KIND_CLOSE, new Uint8Array(0))
        ).catch(() => {
            // The peer is already gone; the local close below is what matters.
        }).then(() => {
            this._closed(code, reason, true);
        });
    }

    /** @internal Called when the open frame comes back from Go. */
    _opened(): void {
        if (this.readyState !== WailsSocket.CONNECTING) return;
        this.readyState = WailsSocket.OPEN;
        this.dispatchEvent(new Event("open"));
    }

    /** @internal Called for each data frame addressed to this connection. */
    _message(payload: ArrayBuffer): void {
        if (this.readyState !== WailsSocket.OPEN) return;
        let data: unknown;
        try {
            data = this._decode(payload);
        } catch {
            this.dispatchEvent(new Event("error"));
            return;
        }
        if (data === payload && this.binaryType === "blob") {
            data = new Blob([payload]);
        }
        this.dispatchEvent(new MessageEvent("message", { data }));
    }

    /** @internal */
    _fail(err: unknown): void {
        // An in-flight send can reject after a clean close has already arrived.
        // The socket's lifecycle is complete then: emitting error after close
        // reverses the WebSocket event order and reports a false failure.
        if (this.readyState === WailsSocket.CLOSED || this.readyState === WailsSocket.CLOSING) return;
        this.dispatchEvent(new Event("error"));
        // 1006: closed abnormally, no close frame. Same code a real WebSocket
        // reports when the connection drops without a handshake.
        this._closed(1006, err instanceof Error ? err.message : String(err), false);
    }

    /** @internal */
    _closed(code: number, reason: string, wasClean: boolean): void {
        if (this.readyState === WailsSocket.CLOSED) return;
        this.readyState = WailsSocket.CLOSED;
        this._openAbort.abort();
        this._sendAbort.abort();
        G.connections.delete(this._id);
        if (G.connections.size === 0) G.pollAbort?.abort();
        this._pending.length = 0;
        this._buffered = 0;

        const ev = typeof CloseEvent === "function"
            ? new CloseEvent("close", { code, reason, wasClean })
            : Object.assign(new Event("close"), { code, reason, wasClean }) as CloseEvent;
        this.dispatchEvent(ev);
    }
}

/**
 * Connect to a stream declared in Go with `app.HandleStream(name, handler)`.
 *
 * Returns immediately with `readyState === CONNECTING`; listen for `open`, or
 * just call {@link WailsSocket.send} once open. Creating a connection at module
 * scope is supported and is the intended shape for generated bindings.
 */
export function Stream(name: string): WailsSocket | WebSocket {
    // Server builds install a factory returning a real WebSocket, since there
    // is a listener there to upgrade on and nothing to emulate. Both objects
    // present the same interface, so nothing above this line cares which it is.
    if (hasDOM) {
        const factory = (window as any)._wails?.streamFactory;
        if (typeof factory === "function") {
            return factory(name);
        }
    }
    return new WailsSocket(name);
}

// Defines an `on<type>` property that registers and unregisters a real event
// listener, the way the DOM defines them. Assigning replaces the previous
// handler; reading returns it.
function defineHandlerProperty(target: EventTarget, type: string): void {
    let current: ((ev: any) => void) | null = null;
    Object.defineProperty(target, "on" + type, {
        get: () => current,
        set(fn: ((ev: any) => void) | null) {
            if (current) target.removeEventListener(type, current);
            current = typeof fn === "function" ? fn : null;
            if (current) target.addEventListener(type, current);
        },
        configurable: true,
        enumerable: true,
    });
}

/**
 * The object-shaped view of a stream returned by {@link JSONStream}: `send`
 * takes a value to stringify, and `ev.data` on a message is the parsed result.
 */
export interface JSONSocket extends Omit<WailsSocket, "send" | "onmessage" | "addEventListener" | "removeEventListener"> {
    send(value: unknown): void;
    onmessage: ((ev: MessageEvent<any>) => void) | null;
    addEventListener(
        type: "message",
        listener: ((this: JSONSocket, ev: MessageEvent<any>) => any) | EventListenerObject | null,
        options?: boolean | AddEventListenerOptions,
    ): void;
    addEventListener(
        type: string,
        listener: EventListenerOrEventListenerObject | null,
        options?: boolean | AddEventListenerOptions,
    ): void;
    removeEventListener(
        type: "message",
        listener: ((this: JSONSocket, ev: MessageEvent<any>) => any) | EventListenerObject | null,
        options?: boolean | EventListenerOptions,
    ): void;
    removeEventListener(
        type: string,
        listener: EventListenerOrEventListenerObject | null,
        options?: boolean | EventListenerOptions,
    ): void;
}

/**
 * A {@link Stream} that speaks objects instead of bytes.
 *
 * `send(value)` marshals with `JSON.stringify`, and `ev.data` in an `onmessage`
 * handler is the parsed object rather than an `ArrayBuffer`. Pairs with
 * `StreamConn.SendJSON` / `ReceiveJSON` on the Go side.
 *
 * ```js
 * const s = JSONStream("telemetry");
 * s.onmessage = (ev) => console.log(ev.data.temperature);
 * s.onopen    = () => s.send({ subscribe: "sensors" });
 * ```
 *
 * A frame that is not valid JSON raises an `error` event and is dropped, rather
 * than throwing from the poll loop and taking the connection down with it.
 */
export function JSONStream(name: string): JSONSocket {
    const socket = Stream(name);
    const decoder = new TextDecoder();
    const parse = (payload: ArrayBuffer | string) =>
        JSON.parse(typeof payload === "string" ? payload : decoder.decode(payload));

    if (socket instanceof WailsSocket) {
        // Decoding through the hook rather than by wrapping onmessage means
        // addEventListener listeners see the parsed object too.
        socket._decode = (payload) => parse(payload);
    } else {
        // Server builds get a native WebSocket, whose dispatch never consults
        // _decode. Patch both listener styles on this instance so JSONStream
        // means the same thing on either transport, which is the point of the
        // shared API.
        const native = socket as WebSocket;
        native.binaryType = "arraybuffer";

        const wrapped = new WeakMap<object, EventListener>();
        const decodedEvents = new WeakMap<Event, MessageEvent | null>();
        const add = native.addEventListener.bind(native);
        const remove = native.removeEventListener.bind(native);

        const decodeEvent = (ev: Event): MessageEvent | null => {
            if (decodedEvents.has(ev)) return decodedEvents.get(ev) ?? null;
            try {
                const value = parse((ev as MessageEvent).data);
                const decoded = new MessageEvent("message", { data: value });
                decodedEvents.set(ev, decoded);
                return decoded;
            } catch {
                decodedEvents.set(ev, null);
                native.dispatchEvent(new Event("error"));
                return null;
            }
        };

        const decode = (listener: any): EventListener => {
            const existing = wrapped.get(listener as object);
            if (existing) return existing;

            const fn: EventListener = (ev) => {
                const decoded = decodeEvent(ev);
                if (!decoded) return;
                if (typeof listener === "function") listener.call(native, decoded);
                else listener.handleEvent(decoded);
            };
            wrapped.set(listener as object, fn);
            return fn;
        };

        native.addEventListener = ((type: string, listener: any, opts?: any) => {
            add(type, type === "message" && listener ? decode(listener) : listener, opts);
        }) as typeof native.addEventListener;

        native.removeEventListener = ((type: string, listener: any, opts?: any) => {
            const fn = type === "message" && listener ? wrapped.get(listener) : undefined;
            remove(type, fn ?? listener, opts);
        }) as typeof native.removeEventListener;

        let handler: ((ev: MessageEvent) => void) | null = null;
        Object.defineProperty(native, "onmessage", {
            get: () => handler,
            set(fn) {
                if (handler) native.removeEventListener("message", handler as EventListener);
                handler = typeof fn === "function" ? fn : null;
                if (handler) native.addEventListener("message", handler as EventListener);
            },
            configurable: true,
            enumerable: true,
        });
    }

    const send = socket.send.bind(socket);
    (socket as unknown as JSONSocket).send = (value: unknown) => send(JSON.stringify(value));

    return socket as unknown as JSONSocket;
}

// Synchronous conversion for everything but a Blob. Avoids a promise per send
// and, more importantly, lets bufferedAmount be exact the moment send() returns.
function toBytesSync(data: string | ArrayBufferLike | ArrayBufferView | Blob): Uint8Array | null {
    if (typeof data === "string") {
        return textEncoder.encode(data);
    }
    if (typeof Blob !== "undefined" && data instanceof Blob) {
        return null;
    }
    if (ArrayBuffer.isView(data)) {
        return new Uint8Array(data.buffer, data.byteOffset, data.byteLength).slice();
    }
    return new Uint8Array(data as ArrayBufferLike).slice();
}

function abortable<T>(promise: Promise<T>, signal?: AbortSignal): Promise<T> {
    if (!signal) return promise;
    if (signal.aborted) return Promise.reject(new StreamRequestCancelled());

    return new Promise((resolve, reject) => {
        const onAbort = () => reject(new StreamRequestCancelled());
        signal.addEventListener("abort", onAbort, { once: true });
        promise.then(
            (value) => {
                signal.removeEventListener("abort", onAbort);
                resolve(value);
            },
            (err) => {
                signal.removeEventListener("abort", onAbort);
                reject(err);
            },
        );
    });
}

async function toBytes(
    data: string | ArrayBufferLike | ArrayBufferView | Blob,
    signal?: AbortSignal,
): Promise<Uint8Array> {
    if (typeof data === "string") {
        return textEncoder.encode(data);
    }
    if (typeof Blob !== "undefined" && data instanceof Blob) {
        return new Uint8Array(await abortable(data.arrayBuffer(), signal));
    }
    if (ArrayBuffer.isView(data)) {
        return new Uint8Array(data.buffer, data.byteOffset, data.byteLength).slice();
    }
    return new Uint8Array(data as ArrayBufferLike).slice();
}

async function postFrame(connID: number, kind: number, body: Uint8Array, name?: string, signal?: AbortSignal): Promise<void> {
    const headers: Record<string, string> = {
        [HDR_SESSION]: G.session,
        [HDR_GENERATION]: String(G.generation),
        [HDR_CONN]: String(connID),
        [HDR_KIND]: String(kind),
        "Content-Type": "application/octet-stream",
    };
    if (name !== undefined) {
        headers[HDR_NAME] = name;
    }

    if (body.byteLength <= CHUNK_THRESHOLD) {
        await postWithRetry(headers, body, signal);
        return;
    }

    // Split oversized frames the same way the runtime splits oversized calls.
    const chunkID = nanoid();
    const total = Math.ceil(body.byteLength / CHUNK_THRESHOLD);
    for (let i = 0; i < total; i++) {
        const slice = body.subarray(i * CHUNK_THRESHOLD, (i + 1) * CHUNK_THRESHOLD);
        await postWithRetry({
            ...headers,
            [HDR_CHUNK]: chunkID,
            [HDR_CHUNK_INDEX]: String(i),
            [HDR_CHUNK_TOTAL]: String(total),
        }, slice, signal);
    }
}

// A 429 means the Go handler has not taken delivery of what is already queued.
// Retrying the same frame keeps ordering (the send chain has not moved on) and
// leaves the request slot free, which is what stops a backed-up connection from
// starving the window's poll.
async function postWithRetry(headers: Record<string, string>, body: Uint8Array, signal?: AbortSignal): Promise<void> {
    // Never retry with zero delay. The receiver being behind is a condition
    // that takes time to clear, and an immediate retry is a busy loop of
    // fetches — each one a scheme-handler round trip, and on the host a cgo
    // call. Start at 1 ms and ramp.
    let wait = 1;
    for (;;) {
        if (signal?.aborted) throw new StreamRequestCancelled();
        const resp = await fetch(streamURL("send"), {
            method: "POST",
            headers,
            body: body as BodyInit,
            signal,
        });
        if (resp.ok) return;
        if (resp.status !== 429) {
            throw new Error(await resp.text());
        }
        await sleep(wait, signal);
        wait = Math.min(wait * 2, 50);
    }
}

// postBatch sends several data frames for one connection in bounded requests.
// Body: count u32, then count x ( len u32 | payload ). The combined body stays
// within CHUNK_THRESHOLD: several individually small frames must not recreate
// the WebView2 body-size problem that chunking single large frames avoids.
async function postBatch(connID: number, frames: Uint8Array[], signal?: AbortSignal): Promise<void> {
    let start = 0;
    while (start < frames.length) {
        if (frames[start].byteLength > CHUNK_THRESHOLD) {
            await postFrame(connID, KIND_DATA, frames[start], undefined, signal);
            start++;
            continue;
        }

        let end = start;
        let size = 4; // frame count
        while (end < frames.length && end - start < MAX_BATCH_FRAMES) {
            const frame = frames[end];
            if (frame.byteLength > CHUNK_THRESHOLD || size + 4 + frame.byteLength > CHUNK_THRESHOLD) {
                break;
            }
            size += 4 + frame.byteLength;
            end++;
        }

        // A near-threshold frame fits as a normal POST body but not inside a
        // batch once its count and length headers are included.
        if (end === start) {
            await postFrame(connID, KIND_DATA, frames[start], undefined, signal);
            start++;
            continue;
        }

        const group = frames.slice(start, end);
        if (group.length === 1) {
            await postFrame(connID, KIND_DATA, group[0], undefined, signal);
        } else {
            await postBatchGroup(connID, group, signal);
        }
        start = end;
    }
}

async function postBatchGroup(connID: number, frames: Uint8Array[], signal?: AbortSignal): Promise<void> {
    if (frames.length === 0 || frames.length > MAX_BATCH_FRAMES) {
        throw new Error("invalid stream batch size");
    }

    let body = buildBatch(frames);
    let sent = 0;
    let wait = 1;
    for (;;) {
        if (signal?.aborted) throw new StreamRequestCancelled();
        const resp = await fetch(streamURL("send"), {
            method: "POST",
            headers: {
                [HDR_SESSION]: G.session,
                [HDR_GENERATION]: String(G.generation),
                [HDR_CONN]: String(connID),
                [HDR_KIND]: String(KIND_DATA),
                [HDR_BATCH]: String(frames.length - sent),
                "Content-Type": "application/octet-stream",
            },
            body: body as BodyInit,
            signal,
        });
        if (resp.ok) return;
        if (resp.status !== 429) {
            throw new Error(await resp.text());
        }
        // The receiver took a prefix of the batch. Resend only what is left,
        // which keeps ordering without re-delivering anything.
        const accepted = Number(resp.headers.get(HDR_BATCH) ?? 0);
        const remaining = frames.length - sent;
        if (!Number.isInteger(accepted) || accepted < 0 || accepted >= remaining) {
            // Zero is valid (no progress), but an out-of-range acknowledgement
            // is a protocol error. Keep zero on the retry path below.
            if (accepted !== 0) throw new Error("invalid stream batch acknowledgement");
        }
        sent += accepted;
        body = buildBatch(frames.slice(sent));
        await sleep(wait, signal);
        wait = Math.min(wait * 2, 50);
    }
}

function buildBatch(frames: Uint8Array[]): Uint8Array {
    let size = 4;
    for (const f of frames) size += 4 + f.byteLength;
    const out = new Uint8Array(size);
    const view = new DataView(out.buffer);
    view.setUint32(0, frames.length);
    let off = 4;
    for (const f of frames) {
        view.setUint32(off, f.byteLength);
        off += 4;
        out.set(f, off);
        off += f.byteLength;
    }
    return out;
}

function startPolling(): void {
    if (G.polling || !hasDOM) return;
    G.polling = true;
    void pollLoop();
}

/**
 * One poll for the whole page, shared by every connection.
 *
 * There is no polling interval and nothing adaptive here on purpose. The server
 * holds the request until a frame exists, so delivery latency is already ~0 and
 * a client-side interval could only add to it. Frames that arrive while a
 * response is in flight simply accumulate and ride the next one, which makes
 * the round trip itself the batching window — it widens as load rises without
 * anything having to measure it.
 */
async function pollLoop(): Promise<void> {
    let backoff = 0;
    const abort = new AbortController();
    G.pollAbort = abort;

    while (G.connections.size > 0) {
        try {
            const resp = await fetch(streamURL("poll"), {
                method: "GET",
                headers: {
                    [HDR_SESSION]: G.session,
                    [HDR_GENERATION]: String(G.generation),
                },
                cache: "no-store",
                signal: abort.signal,
            });

            if (resp.status === 410) {
                // The session is gone: app shutting down, window closed, or we
                // were reaped as unreachable. Nothing to retry.
                closeAll(1001, "session closed");
                break;
            }
            if (!resp.ok && resp.status !== 204) {
                if (!isRetryablePollStatus(resp.status)) {
                    failAll("stream poll failed: " + resp.status);
                    break;
                }
                throw new Error("retryable stream poll failure: " + resp.status);
            }

            backoff = 0;
            if (resp.status !== 204) {
                let buf: ArrayBuffer;
                try {
                    buf = await resp.arrayBuffer();
                } catch {
                    // The Go queue was consumed when the successful response
                    // was produced. Retrying cannot recover a body that was
                    // interrupted after that point, so surface the loss and
                    // stop instead of silently continuing with missing frames.
                    failAll("stream poll response body failed");
                    break;
                }
                const frames = decodeFrames(buf);
                if (frames === null) {
                    // A framing mismatch means a stale cached runtime against a
                    // newer Go side. Retrying cannot fix that, and letting it
                    // fall into the backoff below would spin silently forever,
                    // so fail every connection with a reason instead.
                    closeAll(1002, "unrecognised stream framing");
                    break;
                }
                deliver(frames);
            }
            // Re-poll immediately, on data and on an expired hold alike.
        } catch (err) {
            if (abort.signal.aborted || G.connections.size === 0) break;
            backoff = backoff ? Math.min(backoff * 2, BACKOFF_MAX) : BACKOFF_MIN;
            try {
                await sleep(backoff, abort.signal);
            } catch {
                break;
            }
        }
    }

    if (G.pollAbort === abort) G.pollAbort = undefined;
    G.polling = false;

    // A connection opened while we were winding down: pick the loop back up.
    if (G.connections.size > 0) {
        startPolling();
    }
}

function isRetryablePollStatus(status: number): boolean {
    return status === 408 || status === 425 || status === 429 || (status >= 500 && status <= 599);
}

interface InFrame {
    connID: number;
    kind: number;
    payload: ArrayBuffer;
}

// Validate the complete response before dispatching any part of it. A short
// platform response must not deliver a valid prefix and then throw RangeError:
// the Go queue has already been drained, so retrying cannot recover those
// frames. Treat the whole response as a protocol error instead.
function decodeFrames(buf: ArrayBuffer): InFrame[] | null {
    if (buf.byteLength < 9) return null;
    const dv = new DataView(buf);
    if (dv.getUint8(0) !== 0x57 || dv.getUint8(1) !== 0x53 ||
        dv.getUint8(2) !== 0x31 || dv.getUint8(3) !== 0x00 ||
        (dv.getUint8(4) & ~1) !== 0) {
        return null;
    }

    const count = dv.getUint32(5);
    const frames: InFrame[] = [];
    let off = 9;

    for (let i = 0; i < count; i++) {
        if (off + 9 > buf.byteLength) return null;
        const connID = dv.getUint32(off); off += 4;
        const kind = dv.getUint8(off); off += 1;
        const len = dv.getUint32(off); off += 4;
        if (kind > KIND_ERROR || len > buf.byteLength - off) return null;
        const payload = buf.slice(off, off + len); off += len;
        frames.push({ connID, kind, payload });
    }

    return off === buf.byteLength ? frames : null;
}

function deliver(frames: InFrame[]): void {
    for (const frame of frames) {
        const connID = frame.connID;
        const kind = frame.kind;
        const payload = frame.payload;
        const conn = G.connections.get(connID);
        if (!conn) continue;

        switch (kind) {
            case KIND_DATA:
                conn._message(payload);
                break;
            case KIND_OPEN:
                conn._opened();
                break;
            case KIND_CLOSE:
                conn._closed(1000, "", true);
                break;
            case KIND_ERROR:
                conn._fail(new Error(new TextDecoder().decode(payload)));
                break;
        }
    }
}

function closeAll(code: number, reason: string): void {
    for (const conn of [...G.connections.values()]) {
        conn._closed(code, reason, false);
    }
}

function failAll(reason: string): void {
    for (const conn of [...G.connections.values()]) {
        conn._fail(new Error(reason));
    }
}
