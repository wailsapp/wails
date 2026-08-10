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

// One session per page load, like the runtime's clientId. A reload gets a new
// one, which is what makes a reload look like a closed socket to Go. nanoid is
// backed by crypto.getRandomValues, which matters in server mode where there is
// no window-id header to bind the session against.
// Stored on window rather than in module scope. The runtime can be instantiated
// more than once in a page — the platform injects a copy and the app may import
// it as well — and a per-module id gave each instance its own session. Since a
// new session id for a window supersedes the previous one, the two instances
// then tore down each other's connections in turn, which killed live streams
// mid-run. One id per page is the invariant that matters.
const sessionID: string = hasDOM
    ? ((window as any)._wails = (window as any)._wails || {},
       (window as any)._wails.__streamSession ||= nanoid())
    : nanoid();

const HDR_SESSION = "x-wails-stream-session";
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

// Backoff floor and ceiling for a failing poll. This is the "reasonable lower
// limit" on polling — an error backoff, not a steady-state interval. In normal
// operation the poll re-issues immediately, because the server hold has already
// provided whatever delay was appropriate.
const BACKOFF_MIN = 250;
const BACKOFF_MAX = 5000;

function streamURL(endpoint: string): string {
    return window.location.origin + "/wails/stream/" + endpoint;
}

let nextConnID = 1;
const connections = new Map<number, WailsSocket>();
let pollRunning = false;

function sleep(ms: number): Promise<void> {
    return new Promise((resolve) => setTimeout(resolve, ms));
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

    /** Bytes queued by send() that have not yet reached Go. */
    get bufferedAmount(): number {
        return this._buffered;
    }

    constructor(name: string) {
        super();
        this.name = name;
        this._id = nextConnID++;
        this.url = hasDOM ? streamURL("poll") + "#" + encodeURIComponent(name) : "";

        for (const type of ["open", "message", "close", "error"]) {
            defineHandlerProperty(this, type);
        }

        connections.set(this._id, this);

        // Fire and forget: the ack arrives as an open frame on the poll, which
        // is what moves readyState to OPEN.
        this._chain = this._chain.then(() =>
            postFrame(this._id, KIND_OPEN, new Uint8Array(0), name)
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

        // Resolved now rather than inside the chain, so a Blob is read once and
        // a typed array is captured at its current identity. The bytes are not
        // copied — see toBytes — so a caller must not mutate a buffer it has
        // handed to send().
        const snapshot = toBytes(data);

        // Queue rather than post. Whatever accumulates while a request is in
        // flight goes out together in the next one, so the per-request cost —
        // a scheme-handler round trip, about eleven cgo calls on macOS — is
        // divided across the batch. Under light load a batch is one frame and
        // this behaves exactly as before; the batching only appears when the
        // sender is outrunning the transport, which is when it is needed.
        this._pending.push(snapshot);
        if (this._flushing) return;
        this._flushing = true;

        this._chain = this._chain.then(async () => {
            try {
                while (this._pending.length > 0) {
                    const batch = await Promise.all(this._pending.splice(0));
                    const total = batch.reduce((n, b) => n + b.byteLength, 0);
                    this._buffered += total;
                    try {
                        await postBatch(this._id, batch);
                    } finally {
                        this._buffered -= total;
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
        this.dispatchEvent(new Event("error"));
        // 1006: closed abnormally, no close frame. Same code a real WebSocket
        // reports when the connection drops without a handshake.
        this._closed(1006, err instanceof Error ? err.message : String(err), false);
    }

    /** @internal */
    _closed(code: number, reason: string, wasClean: boolean): void {
        if (this.readyState === WailsSocket.CLOSED) return;
        this.readyState = WailsSocket.CLOSED;
        connections.delete(this._id);

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
export interface JSONSocket extends Omit<WailsSocket, "send" | "onmessage"> {
    send(value: unknown): void;
    onmessage: ((ev: MessageEvent<any>) => void) | null;
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
    const socket = Stream(name) as WailsSocket;
    const decoder = new TextDecoder();

    // Decoding here rather than by wrapping onmessage means addEventListener
    // listeners see the parsed object too.
    socket._decode = (payload) => JSON.parse(decoder.decode(payload));

    const send = socket.send.bind(socket);
    (socket as unknown as JSONSocket).send = (value: unknown) => send(JSON.stringify(value));

    return socket as unknown as JSONSocket;
}

async function toBytes(data: string | ArrayBufferLike | ArrayBufferView | Blob): Promise<Uint8Array> {
    if (typeof data === "string") {
        return new TextEncoder().encode(data);
    }
    if (typeof Blob !== "undefined" && data instanceof Blob) {
        return new Uint8Array(await data.arrayBuffer());
    }
    // A view, not a copy. Copying every frame costs a full memcpy per send,
    // which at several hundred MB/s is enough allocation churn to dominate the
    // transport and destabilise the page. The ownership rule is documented
    // instead: do not mutate a buffer you have passed to send() until the
    // send has been issued.
    if (ArrayBuffer.isView(data)) {
        return new Uint8Array(data.buffer, data.byteOffset, data.byteLength);
    }
    return new Uint8Array(data as ArrayBufferLike);
}

async function postFrame(connID: number, kind: number, body: Uint8Array, name?: string): Promise<void> {
    const headers: Record<string, string> = {
        [HDR_SESSION]: sessionID,
        [HDR_CONN]: String(connID),
        [HDR_KIND]: String(kind),
        "Content-Type": "application/octet-stream",
    };
    if (name !== undefined) {
        headers[HDR_NAME] = name;
    }

    if (body.byteLength <= CHUNK_THRESHOLD) {
        await postWithRetry(headers, body);
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
        }, slice);
    }
}

// A 429 means the Go handler has not taken delivery of what is already queued.
// Retrying the same frame keeps ordering (the send chain has not moved on) and
// leaves the request slot free, which is what stops a backed-up connection from
// starving the window's poll.
async function postWithRetry(headers: Record<string, string>, body: Uint8Array): Promise<void> {
    // Never retry with zero delay. The receiver being behind is a condition
    // that takes time to clear, and an immediate retry is a busy loop of
    // fetches — each one a scheme-handler round trip, and on the host a cgo
    // call. Start at 1 ms and ramp.
    let wait = 1;
    for (;;) {
        const resp = await fetch(streamURL("send"), { method: "POST", headers, body: body as BodyInit });
        if (resp.ok) return;
        if (resp.status !== 429) {
            throw new Error(await resp.text());
        }
        await sleep(wait);
        wait = Math.min(wait * 2, 50);
    }
}

// postBatch sends several data frames for one connection in a single request.
// Body: count u32, then count x ( len u32 | payload ). A frame past the chunk
// threshold is sent on its own, since it has to be split anyway.
async function postBatch(connID: number, frames: Uint8Array[]): Promise<void> {
    if (frames.length === 1 || frames.some((f) => f.byteLength > CHUNK_THRESHOLD)) {
        for (const f of frames) {
            await postFrame(connID, KIND_DATA, f);
        }
        return;
    }

    let body = buildBatch(frames);
    let sent = 0;
    let wait = 1;
    for (;;) {
        const resp = await fetch(streamURL("send"), {
            method: "POST",
            headers: {
                [HDR_SESSION]: sessionID,
                [HDR_CONN]: String(connID),
                [HDR_KIND]: String(KIND_DATA),
                [HDR_BATCH]: String(frames.length - sent),
                "Content-Type": "application/octet-stream",
            },
            body: body as BodyInit,
        });
        if (resp.ok) return;
        if (resp.status !== 429) {
            throw new Error(await resp.text());
        }
        // The receiver took a prefix of the batch. Resend only what is left,
        // which keeps ordering without re-delivering anything.
        const accepted = Number(resp.headers.get(HDR_BATCH) ?? 0);
        sent += accepted;
        body = buildBatch(frames.slice(sent));
        await sleep(wait);
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
    if (pollRunning || !hasDOM) return;
    pollRunning = true;
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

    while (connections.size > 0) {
        try {
            const resp = await fetch(streamURL("poll"), {
                method: "GET",
                headers: { [HDR_SESSION]: sessionID },
                cache: "no-store",
            });

            if (resp.status === 410) {
                // The session is gone: app shutting down, window closed, or we
                // were reaped as unreachable. Nothing to retry.
                closeAll(1001, "session closed");
                break;
            }
            if (!resp.ok && resp.status !== 204) {
                throw new Error("stream poll failed: " + resp.status);
            }

            backoff = 0;
            if (resp.status !== 204) {
                const buf = await resp.arrayBuffer();
                if (!validFraming(buf)) {
                    // A framing mismatch means a stale cached runtime against a
                    // newer Go side. Retrying cannot fix that, and letting it
                    // fall into the backoff below would spin silently forever,
                    // so fail every connection with a reason instead.
                    closeAll(1002, "unrecognised stream framing");
                    break;
                }
                deliver(buf);
            }
            // Re-poll immediately, on data and on an expired hold alike.
        } catch (err) {
            backoff = backoff ? Math.min(backoff * 2, BACKOFF_MAX) : BACKOFF_MIN;
            await sleep(backoff);
        }
    }

    pollRunning = false;

    // A connection opened while we were winding down: pick the loop back up.
    if (connections.size > 0) {
        startPolling();
    }
}

function validFraming(buf: ArrayBuffer): boolean {
    if (buf.byteLength < 9) return false;
    const dv = new DataView(buf);
    return dv.getUint8(0) === 0x57 && dv.getUint8(1) === 0x53 &&
           dv.getUint8(2) === 0x31 && dv.getUint8(3) === 0x00;
}

function deliver(buf: ArrayBuffer): void {
    const dv = new DataView(buf);
    const count = dv.getUint32(5);
    let off = 9;

    for (let i = 0; i < count; i++) {
        const connID = dv.getUint32(off); off += 4;
        const kind = dv.getUint8(off); off += 1;
        const len = dv.getUint32(off); off += 4;
        const payload = buf.slice(off, off + len); off += len;

        const conn = connections.get(connID);
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
    for (const conn of [...connections.values()]) {
        conn._closed(code, reason, false);
    }
}
