import assert from 'node:assert/strict';
import {test} from 'node:test';
import {createLayoutScheduler} from './assets/layout.mjs';

test('resizes during a delayed request apply only the latest measurement next', async () => {
    const frames = [];
    const requests = [];
    const errors = [];
    let width = 100;
    let release;
    const blocked = new Promise(resolve => { release = resolve; });
    const schedule = createLayoutScheduler(async () => {
        requests.push(width);
        if (requests.length === 1) await blocked;
    }, error => errors.push(error), frame => frames.push(frame));

    schedule();
    const finished = frames.shift()();
    width = 200;
    schedule();
    width = 300;
    schedule();
    assert.deepEqual(requests, [100]);
    assert.equal(frames.length, 0);
    release();
    await finished;
    assert.deepEqual(requests, [100, 300]);
    assert.deepEqual(errors, []);
});

test('a failed request reports the error and does not stall later resizes', async () => {
    const frames = [];
    const errors = [];
    const failure = new Error('native request failed');
    let calls = 0;
    const schedule = createLayoutScheduler(async () => {
        if (++calls === 1) throw failure;
    }, error => errors.push(error), frame => frames.push(frame));
    schedule();
    schedule();
    assert.equal(frames.length, 1);
    await frames.shift()();
    schedule();
    await frames.shift()();
    assert.equal(calls, 2);
    assert.deepEqual(errors, [failure]);
});
