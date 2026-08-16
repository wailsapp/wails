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

function decodeBatch(body) {
    const bytes = new Uint8Array(body.buffer, body.byteOffset, body.byteLength);
    const view = new DataView(bytes.buffer, bytes.byteOffset, bytes.byteLength);
    const count = view.getUint32(0);
    const frames = [];
    let offset = 4;
    for (let i = 0; i < count; i++) {
        const length = view.getUint32(offset);
        offset += 4;
        frames.push(Array.from(bytes.subarray(offset, offset + length)));
        offset += length;
    }
    expect(offset).toBe(bytes.byteLength);
    return frames;
}

beforeEach(() => {
    vi.resetModules();
    window.sessionStorage.clear();
    window.name = "";
    window._wails = {};
});

afterEach(() => {
    vi.useRealTimers();
    vi.unstubAllGlobals();
    vi.restoreAllMocks();
    delete window._wails;
});

describe("WailsSocket protocol handling", () => {
    it("can close while its open request is under backpressure", async () => {
        const poll = deferred();
        let openPosts = 0;
        let closePosts = 0;
        const fetch = vi.fn((input, init = {}) => {
            if (String(input).endsWith("/poll")) return poll.promise;
            switch (init.headers?.["x-wails-stream-kind"]) {
                case "1":
                    openPosts++;
                    return Promise.resolve(response(429));
                case "2":
                    closePosts++;
                    return Promise.resolve(response(204));
                default:
                    return Promise.resolve(response(204));
            }
        });
        vi.stubGlobal("fetch", fetch);

        const { WailsSocket } = await import("./stream");
        const socket = new WailsSocket("test");
        await vi.waitFor(() => expect(openPosts).toBeGreaterThanOrEqual(2), { interval: 1, timeout: 100 });

        socket.close();
        await vi.waitFor(() => expect(closePosts).toBe(1), { interval: 1, timeout: 100 });
        await vi.waitFor(() => expect(socket.readyState).toBe(WailsSocket.CLOSED), { interval: 1, timeout: 100 });
    });

    it("can close while its open request is stalled", async () => {
        const poll = deferred();
        let openPosts = 0;
        let openAborts = 0;
        let closePosts = 0;
        const fetch = vi.fn((input, init = {}) => {
            if (String(input).endsWith("/poll")) return poll.promise;
            switch (init.headers?.["x-wails-stream-kind"]) {
                case "1":
                    openPosts++;
                    return new Promise((resolve, reject) => {
                        const onAbort = () => {
                            openAborts++;
                            reject(new DOMException("aborted", "AbortError"));
                        };
                        if (init.signal?.aborted) onAbort();
                        else init.signal?.addEventListener("abort", onAbort, { once: true });
                    });
                case "2":
                    closePosts++;
                    return Promise.resolve(response(204));
                default:
                    return Promise.resolve(response(204));
            }
        });
        vi.stubGlobal("fetch", fetch);

        const { WailsSocket } = await import("./stream");
        const socket = new WailsSocket("test");
        await vi.waitFor(() => expect(openPosts).toBe(1), { interval: 1, timeout: 100 });

        socket.close();
        await vi.waitFor(() => expect(openAborts).toBe(1), { interval: 1, timeout: 100 });
        await vi.waitFor(() => expect(closePosts).toBe(1), { interval: 1, timeout: 100 });
        await vi.waitFor(() => expect(socket.readyState).toBe(WailsSocket.CLOSED), { interval: 1, timeout: 100 });
    });

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

    it("keeps generations monotonic when storage is blocked and the clock moves backwards", async () => {
        const generations = [];
        const fetch = vi.fn((input, init = {}) => {
            if (String(input).endsWith("/poll")) return deferred().promise;
            if (init.headers?.["x-wails-stream-kind"] === "1") {
                generations.push(Number(init.headers["x-wails-stream-generation"]));
            }
            return Promise.resolve(response(204));
        });
        vi.stubGlobal("fetch", fetch);
        vi.spyOn(Storage.prototype, "getItem").mockImplementation(() => {
            throw new DOMException("storage disabled", "SecurityError");
        });
        vi.spyOn(Storage.prototype, "setItem").mockImplementation(() => {
            throw new DOMException("storage disabled", "SecurityError");
        });
        window.name = "application-window";

        vi.stubGlobal("performance", { timeOrigin: 2_000, now: () => 0 });
        let runtime = await import("./stream");
        const first = new runtime.WailsSocket("test");
        await vi.waitFor(() => expect(generations).toHaveLength(1), { interval: 1, timeout: 100 });
        first._closed(1000, "", true);

        window._wails = {};
        vi.resetModules();
        vi.stubGlobal("performance", { timeOrigin: 1_000, now: () => 0 });
        runtime = await import("./stream");
        const second = new runtime.WailsSocket("test");
        await vi.waitFor(() => expect(generations).toHaveLength(2), { interval: 1, timeout: 100 });
        second._closed(1000, "", true);

        expect(generations[1]).toBeGreaterThan(generations[0]);
        expect(window.name.startsWith("application-window")).toBe(true);
    });

    it("keeps generations monotonic when storage becomes blocked between pages", async () => {
        const generations = [];
        const fetch = vi.fn((input, init = {}) => {
            if (String(input).endsWith("/poll")) return deferred().promise;
            if (init.headers?.["x-wails-stream-kind"] === "1") {
                generations.push(Number(init.headers["x-wails-stream-generation"]));
            }
            return Promise.resolve(response(204));
        });
        vi.stubGlobal("fetch", fetch);
        vi.stubGlobal("performance", { timeOrigin: 2_000, now: () => 0 });

        let runtime = await import("./stream");
        const first = new runtime.WailsSocket("test");
        await vi.waitFor(() => expect(generations).toHaveLength(1), { interval: 1, timeout: 100 });
        first._closed(1000, "", true);

        vi.spyOn(Storage.prototype, "getItem").mockImplementation(() => {
            throw new DOMException("storage disabled", "SecurityError");
        });
        vi.spyOn(Storage.prototype, "setItem").mockImplementation(() => {
            throw new DOMException("storage disabled", "SecurityError");
        });
        window._wails = {};
        vi.resetModules();
        vi.stubGlobal("performance", { timeOrigin: 1_000, now: () => 0 });

        runtime = await import("./stream");
        const second = new runtime.WailsSocket("test");
        await vi.waitFor(() => expect(generations).toHaveLength(2), { interval: 1, timeout: 100 });
        second._closed(1000, "", true);

        expect(generations[1]).toBeGreaterThan(generations[0]);
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

    it.each([400, 401, 403, 404, 405, 409, 422, 451])("treats permanent poll response %i as terminal", async (status) => {
        vi.useFakeTimers();
        let pollCalls = 0;
        const fetch = vi.fn((input) => {
            if (!String(input).endsWith("/poll")) return Promise.resolve(response(204));
            pollCalls++;
            return Promise.resolve(response(status));
        });
        vi.stubGlobal("fetch", fetch);

        const { WailsSocket } = await import("./stream");
        const socket = new WailsSocket("test");
        const events = [];
        socket.addEventListener("error", () => events.push("error"));
        socket.addEventListener("close", () => events.push("close"));

        await vi.advanceTimersByTimeAsync(10_000);

        expect(pollCalls).toBe(1);
        expect(events).toEqual(["error", "close"]);
        expect(socket.readyState).toBe(WailsSocket.CLOSED);
    });

    it("treats a retired session as a clean terminal poll response", async () => {
        vi.useFakeTimers();
        let pollCalls = 0;
        const fetch = vi.fn((input) => {
            if (!String(input).endsWith("/poll")) return Promise.resolve(response(204));
            pollCalls++;
            return Promise.resolve(response(410));
        });
        vi.stubGlobal("fetch", fetch);

        const { WailsSocket } = await import("./stream");
        const socket = new WailsSocket("test");
        const events = [];
        socket.addEventListener("error", () => events.push("error"));
        socket.addEventListener("close", () => events.push("close"));
        await vi.advanceTimersByTimeAsync(10_000);

        expect(pollCalls).toBe(1);
        expect(events).toEqual(["close"]);
        expect(socket.readyState).toBe(WailsSocket.CLOSED);
    });

    it("stops after a successful poll response body cannot be read", async () => {
        vi.useFakeTimers();
        let pollCalls = 0;
        const fetch = vi.fn((input) => {
            if (!String(input).endsWith("/poll")) return Promise.resolve(response(204));
            pollCalls++;
            const failed = response(200);
            failed.arrayBuffer = () => Promise.reject(new TypeError("response body interrupted"));
            return Promise.resolve(failed);
        });
        vi.stubGlobal("fetch", fetch);

        const { WailsSocket } = await import("./stream");
        const socket = new WailsSocket("test");
        const events = [];
        socket.addEventListener("error", () => events.push("error"));
        socket.addEventListener("close", () => events.push("close"));
        await vi.advanceTimersByTimeAsync(10_000);

        expect(pollCalls).toBe(1);
        expect(events).toEqual(["error", "close"]);
        expect(socket.readyState).toBe(WailsSocket.CLOSED);
    });

    it("backs off recoverable poll failures and stops after the last connection closes", async () => {
        vi.useFakeTimers();
        const finalPoll = deferred();
        const pollTimes = [];
        let pollCalls = 0;
        let pollAborts = 0;
        const fetch = vi.fn((input, init = {}) => {
            if (!String(input).endsWith("/poll")) return Promise.resolve(response(204));
            pollTimes.push(Date.now());
            pollCalls++;
            if (pollCalls === 1) return Promise.reject(new TypeError("network unavailable"));
            const retryableStatuses = [408, 425, 429, 500, 503, 599];
            if (pollCalls <= retryableStatuses.length + 1) {
                return Promise.resolve(response(retryableStatuses[pollCalls - 2]));
            }
            return new Promise((resolve, reject) => {
                const onAbort = () => {
                    pollAborts++;
                    reject(new DOMException("aborted", "AbortError"));
                };
                if (init.signal?.aborted) onAbort();
                else init.signal?.addEventListener("abort", onAbort, { once: true });
                finalPoll.promise.then(resolve, reject);
            });
        });
        vi.stubGlobal("fetch", fetch);

        const { WailsSocket } = await import("./stream");
        const socket = new WailsSocket("test");
        await vi.advanceTimersByTimeAsync(0);
        expect(pollCalls).toBe(1);
        const delays = [250, 500, 1000, 2000, 4000, 5000, 5000];
        for (let i = 0; i < delays.length; i++) {
            await vi.advanceTimersByTimeAsync(delays[i] - 1);
            expect(pollCalls).toBe(i + 1);
            await vi.advanceTimersByTimeAsync(1);
            expect(pollCalls).toBe(i + 2);
        }
        expect(pollTimes.map((time) => time - pollTimes[0])).toEqual([
            0, 250, 750, 1750, 3750, 7750, 12750, 17750,
        ]);

        socket._closed(1000, "peer closed", true);
        await vi.advanceTimersByTimeAsync(10_000);
        expect(pollCalls).toBe(8);
        expect(pollAborts).toBe(1);
    });

    it("cancels a recoverable poll retry timer after the last connection closes", async () => {
        vi.useFakeTimers();
        let pollCalls = 0;
        const fetch = vi.fn((input) => {
            if (!String(input).endsWith("/poll")) return Promise.resolve(response(204));
            pollCalls++;
            return Promise.reject(new TypeError("network unavailable"));
        });
        vi.stubGlobal("fetch", fetch);

        const { WailsSocket } = await import("./stream");
        const socket = new WailsSocket("test");
        await vi.advanceTimersByTimeAsync(0);
        expect(pollCalls).toBe(1);
        expect(vi.getTimerCount()).toBe(1);

        socket._closed(1000, "peer closed", true);
        await vi.advanceTimersByTimeAsync(30_000);
        expect(pollCalls).toBe(1);
        expect(vi.getTimerCount()).toBe(0);
    });

    it("keeps the shared poll for remaining connections and restarts after wind-down", async () => {
        let pollCalls = 0;
        let pollAborts = 0;
        const fetch = vi.fn((input, init = {}) => {
            if (!String(input).endsWith("/poll")) return Promise.resolve(response(204));
            pollCalls++;
            return new Promise((resolve, reject) => {
                const onAbort = () => {
                    pollAborts++;
                    reject(new DOMException("aborted", "AbortError"));
                };
                if (init.signal?.aborted) onAbort();
                else init.signal?.addEventListener("abort", onAbort, { once: true });
            });
        });
        vi.stubGlobal("fetch", fetch);

        const { WailsSocket } = await import("./stream");
        const first = new WailsSocket("first");
        const second = new WailsSocket("second");
        await vi.waitFor(() => expect(pollCalls).toBe(1), { interval: 1, timeout: 100 });

        first._closed(1000, "peer closed", true);
        await Promise.resolve();
        expect(pollAborts).toBe(0);

        second._closed(1000, "peer closed", true);
        const replacement = new WailsSocket("replacement");
        await vi.waitFor(() => expect(pollAborts).toBe(1), { interval: 1, timeout: 100 });
        await vi.waitFor(() => expect(pollCalls).toBe(2), { interval: 1, timeout: 100 });

        replacement._closed(1000, "peer closed", true);
        await vi.waitFor(() => expect(pollAborts).toBe(2), { interval: 1, timeout: 100 });
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

    it("can close while an outbound data frame is under backpressure", async () => {
        const poll = deferred();
        let dataPosts = 0;
        let closePosts = 0;
        const fetch = vi.fn((input, init = {}) => {
            if (String(input).endsWith("/poll")) return poll.promise;
            switch (init.headers?.["x-wails-stream-kind"]) {
                case "0":
                    dataPosts++;
                    return Promise.resolve(response(429));
                case "2":
                    closePosts++;
                    return Promise.resolve(response(204));
                default:
                    return Promise.resolve(response(204));
            }
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
        await vi.waitFor(() => expect(dataPosts).toBeGreaterThanOrEqual(2), { interval: 1, timeout: 100 });

        socket.close();
        expect(socket.bufferedAmount).toBe(0);
        await vi.waitFor(() => expect(closePosts).toBe(1), { interval: 1, timeout: 100 });
        await vi.waitFor(() => expect(socket.readyState).toBe(WailsSocket.CLOSED), { interval: 1, timeout: 100 });

        expect(socket.bufferedAmount).toBe(0);
        expect(events).toEqual(["close"]);
    });

    it("can close while an outbound data request is stalled", async () => {
        const poll = deferred();
        let dataPosts = 0;
        let closePosts = 0;
        let dataAborts = 0;
        const fetch = vi.fn((input, init = {}) => {
            if (String(input).endsWith("/poll")) return poll.promise;
            switch (init.headers?.["x-wails-stream-kind"]) {
                case "0":
                    dataPosts++;
                    return new Promise((resolve, reject) => {
                        const onAbort = () => {
                            dataAborts++;
                            reject(new DOMException("aborted", "AbortError"));
                        };
                        if (init.signal?.aborted) onAbort();
                        else init.signal?.addEventListener("abort", onAbort, { once: true });
                    });
                case "2":
                    closePosts++;
                    return Promise.resolve(response(204));
                default:
                    return Promise.resolve(response(204));
            }
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
        await vi.waitFor(() => expect(dataPosts).toBe(1), { interval: 1, timeout: 100 });

        socket.close();
        expect(socket.bufferedAmount).toBe(0);
        await vi.waitFor(() => expect(dataAborts).toBe(1), { interval: 1, timeout: 100 });
        await vi.waitFor(() => expect(closePosts).toBe(1), { interval: 1, timeout: 100 });
        await vi.waitFor(() => expect(socket.readyState).toBe(WailsSocket.CLOSED), { interval: 1, timeout: 100 });

        expect(socket.bufferedAmount).toBe(0);
        expect(events).toEqual(["close"]);
    });

    it("can close while an outbound Blob conversion is stalled", async () => {
        const poll = deferred();
        const conversion = deferred();
        let conversionStarted = 0;
        let dataPosts = 0;
        let closePosts = 0;
        const fetch = vi.fn((input, init = {}) => {
            if (String(input).endsWith("/poll")) return poll.promise;
            switch (init.headers?.["x-wails-stream-kind"]) {
                case "0":
                    dataPosts++;
                    return Promise.resolve(response(204));
                case "2":
                    closePosts++;
                    return Promise.resolve(response(204));
                default:
                    return Promise.resolve(response(204));
            }
        });
        vi.stubGlobal("fetch", fetch);

        const { WailsSocket } = await import("./stream");
        const socket = new WailsSocket("test");
        await Promise.resolve();
        socket._opened();

        const blob = new Blob([new Uint8Array(8)]);
        Object.defineProperty(blob, "arrayBuffer", {
            value: () => {
                conversionStarted++;
                return conversion.promise;
            },
        });
        socket.send(blob);
        expect(socket.bufferedAmount).toBe(8);
        await vi.waitFor(() => expect(conversionStarted).toBe(1), { interval: 1, timeout: 100 });
        await vi.waitFor(() => expect(socket._pending).toHaveLength(0), { interval: 1, timeout: 100 });

        socket.close();
        expect(socket.bufferedAmount).toBe(0);
        await vi.waitFor(() => expect(closePosts).toBe(1), { interval: 1, timeout: 100 });
        await vi.waitFor(() => expect(socket.readyState).toBe(WailsSocket.CLOSED), { interval: 1, timeout: 100 });
        expect(dataPosts).toBe(0);
    });

    it("can discard a queued Blob conversion before the flush starts", async () => {
        const poll = deferred();
        const conversion = deferred();
        let closePosts = 0;
        const fetch = vi.fn((input, init = {}) => {
            if (String(input).endsWith("/poll")) return poll.promise;
            if (init.headers?.["x-wails-stream-kind"] === "2") closePosts++;
            return Promise.resolve(response(204));
        });
        vi.stubGlobal("fetch", fetch);

        const { WailsSocket } = await import("./stream");
        const socket = new WailsSocket("test");
        socket._opened();

        const blob = new Blob([new Uint8Array(8)]);
        Object.defineProperty(blob, "arrayBuffer", { value: () => conversion.promise });
        socket.send(blob);
        socket.close();

        await vi.waitFor(() => expect(closePosts).toBe(1), { interval: 1, timeout: 100 });
        await vi.waitFor(() => expect(socket.readyState).toBe(WailsSocket.CLOSED), { interval: 1, timeout: 100 });
        await Promise.resolve();
        expect(socket.bufferedAmount).toBe(0);
    });

    it("accounts for and releases a sustained pending send burst", async () => {
        const poll = deferred();
        let dataStarted = 0;
        const fetch = vi.fn((input, init = {}) => {
            if (String(input).endsWith("/poll")) return poll.promise;
            if (init.headers?.["x-wails-stream-kind"] !== "0") {
                return Promise.resolve(response(204));
            }
            dataStarted++;
            return new Promise((resolve, reject) => {
                const onAbort = () => reject(new DOMException("aborted", "AbortError"));
                if (init.signal?.aborted) onAbort();
                else init.signal?.addEventListener("abort", onAbort, { once: true });
            });
        });
        vi.stubGlobal("fetch", fetch);

        const { WailsSocket } = await import("./stream");
        const socket = new WailsSocket("test");
        await Promise.resolve();
        socket._opened();

        const frameBytes = 1024;
        const frames = 256;
        for (let i = 0; i < frames; i++) socket.send(new Uint8Array(frameBytes));
        expect(socket.bufferedAmount).toBe(frames * frameBytes);
        await vi.waitFor(() => expect(dataStarted).toBe(1), { interval: 1, timeout: 100 });

        socket._closed(1000, "peer closed", true);
        expect(socket.bufferedAmount).toBe(0);
        await Promise.resolve();
        await Promise.resolve();
        expect(socket.bufferedAmount).toBe(0);
        expect(dataStarted).toBe(1);
    });

    it("releases every pending byte after a terminal send failure", async () => {
        const poll = deferred();
        const dataPost = deferred();
        let dataPosts = 0;
        const fetch = vi.fn((input, init = {}) => {
            if (String(input).endsWith("/poll")) return poll.promise;
            if (init.headers?.["x-wails-stream-kind"] !== "0") {
                return Promise.resolve(response(204));
            }
            dataPosts++;
            return dataPost.promise;
        });
        vi.stubGlobal("fetch", fetch);

        const { WailsSocket } = await import("./stream");
        const socket = new WailsSocket("test");
        await Promise.resolve();
        socket._opened();

        const events = [];
        socket.addEventListener("error", () => events.push("error"));
        socket.addEventListener("close", () => events.push("close"));
        for (let i = 0; i < 32; i++) socket.send(new Uint8Array(2048));
        expect(socket.bufferedAmount).toBe(32 * 2048);
        await vi.waitFor(() => expect(dataPosts).toBe(1), { interval: 1, timeout: 100 });

        dataPost.resolve(response(500));
        await vi.waitFor(() => expect(socket.readyState).toBe(WailsSocket.CLOSED), { interval: 1, timeout: 100 });
        expect(socket.bufferedAmount).toBe(0);
        expect(events).toEqual(["error", "close"]);
        expect(dataPosts).toBe(1);
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

    it("snapshots every mutable binary input when send is called", async () => {
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

        // Hold one request open so every tested value remains queued after
        // send() returns and before its bytes are assembled into a batch.
        socket.send(new Uint8Array([0]));
        await vi.waitFor(() => expect(dataBodies).toHaveLength(1), { interval: 1, timeout: 100 });

        const arrayBuffer = new Uint8Array([1, 2]).buffer;
        const typedBacking = new Uint8Array([90, 3, 4, 91]);
        const typedView = new Uint8Array(typedBacking.buffer, 1, 2);
        const dataViewBacking = new Uint8Array([92, 5, 6, 93]);
        const dataView = new DataView(dataViewBacking.buffer, 1, 2);
        const sharedBuffer = new SharedArrayBuffer(2);
        new Uint8Array(sharedBuffer).set([7, 8]);

        socket.send(arrayBuffer);
        socket.send(typedView);
        socket.send(dataView);
        socket.send(sharedBuffer);

        new Uint8Array(arrayBuffer).fill(0xa1);
        typedBacking.fill(0xa2);
        dataViewBacking.fill(0xa3);
        new Uint8Array(sharedBuffer).fill(0xa4);
        firstDataPost.resolve(response(204));

        await vi.waitFor(() => expect(dataBodies).toHaveLength(2), { interval: 1, timeout: 100 });
        expect(decodeBatch(dataBodies[1])).toEqual([
            [1, 2],
            [3, 4],
            [5, 6],
            [7, 8],
        ]);
        await vi.waitFor(() => expect(socket.bufferedAmount).toBe(0), { interval: 1, timeout: 100 });
        socket._closed(1000, "", true);
    });

    it("snapshots each queued send when one mutable buffer is reused", async () => {
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

        socket.send(new Uint8Array([0]));
        await vi.waitFor(() => expect(dataBodies).toHaveLength(1), { interval: 1, timeout: 100 });

        const reused = new Uint8Array(2);
        reused.set([1, 1]);
        socket.send(reused);
        reused.set([2, 2]);
        socket.send(reused);
        reused.set([3, 3]);
        socket.send(reused);
        reused.fill(9);
        firstDataPost.resolve(response(204));

        await vi.waitFor(() => expect(dataBodies).toHaveLength(2), { interval: 1, timeout: 100 });
        expect(decodeBatch(dataBodies[1])).toEqual([
            [1, 1],
            [2, 2],
            [3, 3],
        ]);
        await vi.waitFor(() => expect(socket.bufferedAmount).toBe(0), { interval: 1, timeout: 100 });
        socket._closed(1000, "", true);
    });

    it("reuses one encoder across string sends", async () => {
        const poll = deferred();
        const NativeTextEncoder = globalThis.TextEncoder;
        let constructions = 0;
        vi.stubGlobal("TextEncoder", class {
            encoder = new NativeTextEncoder();

            constructor() {
                constructions++;
            }

            encode(value) {
                return this.encoder.encode(value);
            }
        });
        vi.stubGlobal("fetch", vi.fn((input) =>
            String(input).endsWith("/poll") ? poll.promise : Promise.resolve(response(204))
        ));

        const { WailsSocket } = await import("./stream");
        const socket = new WailsSocket("test");
        socket._opened();
        socket.send("one");
        socket.send("two");
        socket.send("three");

        await vi.waitFor(() => expect(socket.bufferedAmount).toBe(0), { interval: 1, timeout: 100 });
        expect(constructions).toBe(1);
        socket._closed(1000, "", true);
    });

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

describe("JSONStream in desktop mode", () => {
    it("decodes once for onmessage and distinct listeners", async () => {
        const poll = deferred();
        vi.stubGlobal("fetch", vi.fn((input) =>
            String(input).endsWith("/poll") ? poll.promise : Promise.resolve(response(204))
        ));

        const { JSONStream } = await import("./stream");
        const socket = JSONStream("json");
        socket._opened();
        const received = [];
        socket.addEventListener("message", (event) => received.push(event.data));
        socket.addEventListener("message", (event) => received.push(event.data));

        const payload = new TextEncoder().encode('{"ok":true}');
        socket._message(payload.buffer);

        expect(received).toHaveLength(2);
        expect(received[0]).toBe(received[1]);
        socket._closed(1000, "", true);
    });
});

describe("JSONStream in server mode", () => {
    class FakeSocket {
        binaryType = "blob";
        listeners = new Map();

        send() {}

        addEventListener(type, listener) {
            const listeners = this.listeners.get(type) ?? [];
            if (!listeners.includes(listener)) listeners.push(listener);
            this.listeners.set(type, listeners);
        }

        removeEventListener(type, listener) {
            const listeners = this.listeners.get(type) ?? [];
            this.listeners.set(type, listeners.filter((candidate) => candidate !== listener));
        }

        dispatchEvent(event) {
            for (const listener of [...(this.listeners.get(event.type) ?? [])]) {
                if (typeof listener === "function") listener.call(this, event);
                else listener.handleEvent(event);
            }
            return !event.defaultPrevented;
        }
    }

    it("decodes once for onmessage and distinct listeners", async () => {
        const native = new FakeSocket();
        window._wails.streamFactory = () => native;

        const { JSONStream } = await import("./stream");
        const socket = JSONStream("json");
        const received = [];
        socket.onmessage = (event) => received.push(event.data);
        socket.addEventListener("message", (event) => received.push(event.data));

        const payload = new TextEncoder().encode('{"ok":true}');
        native.dispatchEvent(new MessageEvent("message", { data: payload.buffer }));

        expect(received).toHaveLength(2);
        expect(received[0]).toBe(received[1]);
    });

    it("preserves addEventListener duplicate suppression while decoding", async () => {
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

    it("preserves listener removal while decoding", async () => {
        const native = new FakeSocket();
        window._wails.streamFactory = () => native;

        const { JSONStream } = await import("./stream");
        const socket = JSONStream("json");
        const listener = vi.fn();
        socket.addEventListener("message", listener);
        socket.removeEventListener("message", listener);

        const payload = new TextEncoder().encode('{"ok":true}');
        native.dispatchEvent(new MessageEvent("message", { data: payload.buffer }));

        expect(listener).not.toHaveBeenCalled();
    });

    it("emits one error and drops malformed JSON for every listener", async () => {
        const native = new FakeSocket();
        window._wails.streamFactory = () => native;

        const { JSONStream } = await import("./stream");
        const socket = JSONStream("json");
        const first = vi.fn();
        const second = vi.fn();
        const handler = vi.fn();
        const onerror = vi.fn();
        socket.onmessage = handler;
        socket.addEventListener("message", first);
        socket.addEventListener("message", second);
        socket.addEventListener("error", onerror);

        const payload = new TextEncoder().encode("not json");
        native.dispatchEvent(new MessageEvent("message", { data: payload.buffer }));

        expect(handler).not.toHaveBeenCalled();
        expect(first).not.toHaveBeenCalled();
        expect(second).not.toHaveBeenCalled();
        expect(onerror).toHaveBeenCalledTimes(1);
    });
});
