import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

beforeEach(() => {
    vi.resetModules();
});

afterEach(() => {
    delete window.wails;
    delete window._wailsAndroidCallback;
    vi.unstubAllGlobals();
});

async function caller(transport, response) {
    if (transport === "Android") {
        window.wails = {
            invokeAsync: vi.fn((id) => {
                window._wailsAndroidCallback(id, JSON.stringify(response), null);
            }),
        };
    } else {
        vi.stubGlobal("fetch", vi.fn(async () => new Response(JSON.stringify(response), {
            status: 422,
            headers: { "Content-Type": "application/json" },
        })));
    }
    const runtime = await import("./runtime");
    return { call: runtime.newRuntimeCaller(0), RuntimeError: runtime.RuntimeError };
}

describe.each(["HTTP", "Android"])("%s call errors", (transport) => {
    it.each(["ReferenceError", "TypeError", "RuntimeError", "UnknownError"])("restores %s and nested cause", async (kind) => {
        const cause = { detail: { message: "try again", values: [1, true, null] }, retry: false };
        const error = { kind, message: "service failed", cause };
        const { call, RuntimeError } = await caller(transport, transport === "Android" ? { ok: false, error } : error);
        const constructors = { ReferenceError, TypeError, RuntimeError };
        const caught = await call(0).catch((err) => err);
        expect(caught.constructor).toBe(constructors[kind] ?? Error);
        expect(caught.message).toBe(error.message);
        expect(caught.cause).toEqual(cause);
    });

    it.each([null, false, 0, "text", [1, { nested: true }]].map((cause) => [cause]))("preserves arbitrary JSON cause %j", async (cause) => {
        const error = { kind: "RuntimeError", message: "service failed", cause };
        const { call, RuntimeError } = await caller(transport, transport === "Android" ? { ok: false, error } : error);
        const caught = await call(0).catch((err) => err);
        expect(caught).toBeInstanceOf(RuntimeError);
        expect(caught.cause).toEqual(cause);
    });

    it("restores errors without a cause", async () => {
        const error = { kind: "TypeError", message: "invalid arguments" };
        const { call } = await caller(transport, transport === "Android" ? { ok: false, error } : error);
        const caught = await call(0).catch((err) => err);
        expect(caught).toBeInstanceOf(TypeError);
        expect(caught.message).toBe(error.message);
        expect(caught.cause).toBeUndefined();
    });
});

describe("Android response envelopes", () => {
    it.each(["legacy failure", "", null, undefined])("accepts legacy error %j", async (error) => {
        const { call } = await caller("Android", { ok: false, error });
        const caught = await call(0).catch((err) => err);
        expect(caught.constructor).toBe(Error);
        expect(caught.message).toBe(error ?? "unknown runtime call error");
        expect(caught.cause).toBeUndefined();
    });

    it.each([
        [{ ok: true, data: { value: 42 } }, { value: 42 }],
        [{ ok: true, data: null }, null],
        [{ ok: true, text: "hello" }, "hello"],
        [{ ok: true, text: "" }, ""],
    ])("resolves successful envelope %j", async (envelope, expected) => {
        const { call } = await caller("Android", envelope);
        await expect(call(0)).resolves.toEqual(expected);
    });

    it("preserves Java bridge errors", async () => {
        const { call } = await caller("Android", {});
        window.wails.invokeAsync.mockImplementation((id) => {
            window._wailsAndroidCallback(id, null, "bridge failed");
        });
        await expect(call(0)).rejects.toThrow("bridge failed");
    });
});
