import assert from "node:assert/strict";
import test from "node:test";
import environment from "./protocol-test-environment.mjs";

test("protocol environment isolates browsing-context state and preserves native APIs", () => {
    const first = { Event, AbortSignal, fetch };
    const second = {};
    const a = environment.setup(first);
    const b = environment.setup(second);
    assert.equal(first.window.sessionStorage.getItem("missing"), null);
    first.window.sessionStorage.setItem(5, 9);
    assert.equal(first.window.sessionStorage.getItem("5"), "9");
    assert.equal(second.window.sessionStorage.getItem("5"), null);
    first.window.sessionStorage.setItem("5", "updated");
    assert.equal(first.window.sessionStorage.getItem(5), "updated");
    first.window.sessionStorage.clear();
    assert.equal(first.window.sessionStorage.getItem(5), null);
    first.window.name = "reloaded-window";
    assert.equal(second.window.name, "");
    assert.equal(first.window.location.origin, "http://localhost:3000");
    assert.equal(first.Event, Event);
    assert.equal(first.AbortSignal, AbortSignal);
    assert.equal(first.fetch, fetch);
    assert.equal(first.document.createElement, undefined);
    a.teardown();
    b.teardown();
    assert.deepEqual(first, { Event, AbortSignal, fetch });
    assert.deepEqual(second, {});
});

test("teardown restores original global descriptors", () => {
    const original = { browser: "original" };
    const target = {};
    Object.defineProperty(target, "window", { get: () => original, configurable: true });
    const descriptor = Object.getOwnPropertyDescriptor(target, "window");
    const scope = environment.setup(target);
    assert.notEqual(target.window, original);
    scope.teardown();
    assert.equal(target.window, original);
    assert.deepEqual(Object.getOwnPropertyDescriptor(target, "window"), descriptor);
});
