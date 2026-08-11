import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

function response(status, body = new ArrayBuffer(0)) {
    return {
        status,
        ok: status >= 200 && status < 300,
        arrayBuffer: async () => body,
        text: async () => "",
        headers: new Headers(),
    };
}

function deferred() {
    let resolve;
    let reject;
    const promise = new Promise((r, j) => { resolve = r; reject = j; });
    return { promise, resolve, reject };
}

beforeEach(() => {
    vi.resetModules();
    window.sessionStorage.clear();
    window._wails = {};
});

afterEach(() => {
    vi.unstubAllGlobals();
    delete window._wails;
});

describe("WailsSocket protocol handling", () => {
    it("advances the page generation across a reload", async () => {
        const generations = [];
        const fetch = vi.fn((input, init = {}) => {
            if (String(input).endsWith("/poll")) return deferred().promise;
            if (init.headers?.["x-wails-stream-kind"] === "1") {
                generations.push(Number(init.headers["x-wails-stream-generation"]));
            }
            return Promise.resolve(response(204));
        });
        vi.stubGlobal("fetch", fetch);

        let runtime = await import("./stream");
        const first = new runtime.WailsSocket("test");
        await vi.waitFor(() => expect(generations).toHaveLength(1), { interval: 1, timeout: 100 });
        first._closed(1000, "", true);

        // A real reload replaces window globals and re-evaluates the runtime,
        // while sessionStorage remains attached to the browsing context.
        window._wails = {};
        vi.resetModules();
        runtime = await import("./stream");
        const second = new runtime.WailsSocket("test");
        await vi.waitFor(() => expect(generations).toHaveLength(2), { interval: 1, timeout: 100 });
        second._closed(1000, "", true);

        expect(generations[0]).toBeGreaterThan(0);
        expect(generations[1]).toBeGreaterThan(generations[0]);
    }, 1000);

    it("creates a valid generation when performance.timeOrigin is unavailable", async () => {
        const generations = [];
        const fetch = vi.fn((input, init = {}) => {
            if (String(input).endsWith("/poll")) return deferred().promise;
            if (init.headers?.["x-wails-stream-kind"] === "1") {
                generations.push(Number(init.headers["x-wails-stream-generation"]));
            }
            return Promise.resolve(response(204));
        });
        vi.stubGlobal("fetch", fetch);
        vi.stubGlobal("performance", { now: () => Date.now() });

        const { WailsSocket } = await import("./stream");
        const socket = new WailsSocket("test");
        await Promise.resolve();
        await Promise.resolve();

        expect(generations).toHaveLength(1);
        expect(Number.isSafeInteger(generations[0])).toBe(true);
        expect(generations[0]).toBeGreaterThan(0);
        socket._closed(1000, "", true);
    });

    it("closes every connection on a truncated poll frame", async () => {
        const truncated = new ArrayBuffer(9);
        const view = new DataView(truncated);
        view.setUint8(0, 0x57);
        view.setUint8(1, 0x53);
        view.setUint8(2, 0x31);
        view.setUint8(3, 0x00);
        view.setUint32(5, 1); // Claims one frame, but carries no frame header.

        const fetch = vi.fn(async (input) =>
            String(input).endsWith("/poll") ? response(200, truncated) : response(204)
        );
        vi.stubGlobal("fetch", fetch);

        const { WailsSocket } = await import("./stream");
        const socket = new WailsSocket("test");
        const closed = new Promise((resolve) => socket.addEventListener("close", resolve));

        const event = await closed;
        expect(event.code).toBe(1002);
        expect(socket.readyState).toBe(WailsSocket.CLOSED);
        expect(fetch.mock.calls.filter(([url]) => String(url).endsWith("/poll"))).toHaveLength(1);
    });

    it("clears bufferedAmount when the peer closes during an in-flight send", async () => {
        const poll = deferred();
        const dataPost = deferred();
        const fetch = vi.fn((input, init = {}) => {
            if (String(input).endsWith("/poll")) return poll.promise;
            if (init.headers?.["x-wails-stream-kind"] === "0") return dataPost.promise;
            return Promise.resolve(response(204));
        });
        vi.stubGlobal("fetch", fetch);

        const { WailsSocket } = await import("./stream");
        const socket = new WailsSocket("test");
        await Promise.resolve(); // Let the open POST complete.
        socket._opened();

        socket.send(new Uint8Array(8));
        expect(socket.bufferedAmount).toBe(8);
        await Promise.resolve(); // Let the data POST enter fetch.

        socket._closed(1000, "", true);
        expect(socket.bufferedAmount).toBe(0);

        dataPost.resolve(response(204));
        await Promise.resolve();
        await Promise.resolve();
        expect(socket.bufferedAmount).toBe(0);
    });

    it("does not emit an error after a clean close when an in-flight send rejects", async () => {
        const poll = deferred();
        const dataPost = deferred();
        const fetch = vi.fn((input, init = {}) => {
            if (String(input).endsWith("/poll")) return poll.promise;
            if (init.headers?.["x-wails-stream-kind"] === "0") return dataPost.promise;
            return Promise.resolve(response(204));
        });
        vi.stubGlobal("fetch", fetch);

        const { WailsSocket } = await import("./stream");
        const socket = new WailsSocket("test");
        await Promise.resolve();
        socket._opened();

        const events = [];
        socket.addEventListener("error", () => events.push("error"));
        socket.addEventListener("close", () => events.push("close"));
        socket.send(new Uint8Array(8));
        await vi.waitFor(() => expect(fetch).toHaveBeenCalledTimes(3), { interval: 1, timeout: 100 });

        socket._closed(1000, "", true);
        dataPost.reject(new Error("request ended after the peer closed"));
        await Promise.resolve();
        await Promise.resolve();

        expect(events).toEqual(["close"]);
    });

    it("keeps every accumulated batch within the webview request limit", async () => {
        const poll = deferred();
        const firstDataPost = deferred();
        const dataBodies = [];
        const fetch = vi.fn((input, init = {}) => {
            if (String(input).endsWith("/poll")) return poll.promise;
            if (init.headers?.["x-wails-stream-kind"] !== "0") {
                return Promise.resolve(response(204));
            }
            dataBodies.push(init.body);
            if (dataBodies.length === 1) return firstDataPost.promise;
            return Promise.resolve(response(204));
        });
        vi.stubGlobal("fetch", fetch);

        const { WailsSocket } = await import("./stream");
        const socket = new WailsSocket("test");
        await Promise.resolve();
        socket._opened();

        // Hold one request open so the following frames accumulate into the
        // next flush, which is when batching is useful and must stay bounded.
        socket.send(new Uint8Array(1));
        await vi.waitFor(() => expect(dataBodies).toHaveLength(1), { interval: 1, timeout: 100 });
        for (let i = 0; i < 6; i++) {
            socket.send(new Uint8Array(100 * 1024));
        }
        firstDataPost.resolve(response(204));

        await vi.waitFor(() => expect(dataBodies).toHaveLength(3), { interval: 1, timeout: 100 });
        for (const body of dataBodies) {
            expect(body.byteLength).toBeLessThanOrEqual(512 * 1024);
        }
        socket._closed(1000, "", true);
    }, 1000);

    it("retries the completing chunk unchanged when Go applies backpressure", async () => {
        const poll = deferred();
        const chunks = [];
        let rejectedFinal = false;
        const fetch = vi.fn((input, init = {}) => {
            if (String(input).endsWith("/poll")) return poll.promise;
            if (init.headers?.["x-wails-stream-kind"] !== "0") {
                return Promise.resolve(response(204));
            }
            const index = Number(init.headers["x-wails-stream-chunk-index"]);
            chunks.push({ index, body: new Uint8Array(init.body) });
            if (index === 1 && !rejectedFinal) {
                rejectedFinal = true;
                return Promise.resolve(response(429));
            }
            return Promise.resolve(response(204));
        });
        vi.stubGlobal("fetch", fetch);

        const { WailsSocket } = await import("./stream");
        const socket = new WailsSocket("test");
        await Promise.resolve();
        socket._opened();

        socket.send(new Uint8Array(600 * 1024).fill(0x5a));
        await vi.waitFor(() => expect(chunks).toHaveLength(3), { interval: 1, timeout: 200 });

        expect(chunks.map((chunk) => chunk.index)).toEqual([0, 1, 1]);
        expect(chunks[2].body).toEqual(chunks[1].body);
        await vi.waitFor(() => expect(socket.bufferedAmount).toBe(0), { interval: 1, timeout: 100 });
        socket._closed(1000, "", true);
    }, 1000);
});

describe("JSONStream in server mode", () => {
    it("preserves addEventListener duplicate suppression while decoding", async () => {
        class FakeSocket extends EventTarget {
            binaryType = "blob";
            send() {}
        }

        const native = new FakeSocket();
        window._wails.streamFactory = () => native;

        const { JSONStream } = await import("./stream");
        const socket = JSONStream("json");
        const listener = vi.fn();
        socket.addEventListener("message", listener);
        socket.addEventListener("message", listener);

        const payload = new TextEncoder().encode('{"ok":true}');
        native.dispatchEvent(new MessageEvent("message", { data: payload.buffer }));

        expect(listener).toHaveBeenCalledTimes(1);
        expect(listener.mock.calls[0][0].data).toEqual({ ok: true });
    });
});
