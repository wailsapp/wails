import { readFileSync } from 'fs';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';

const __dirname = dirname(fileURLToPath(import.meta.url));

const eventsCode = readFileSync(join(__dirname, '../../internal/frontend/runtime/desktop/events.js'), 'utf8');

const mockWindow = {
    WailsInvoke: () => {}
};

let code = eventsCode
    .replace(/export /g, '')
    .replace(/window\.WailsInvoke/g, 'mockWindow.WailsInvoke');

// Execute the code and capture all declared names into globalThis
const wrappedCode = `
${code}
return {
    eventListeners,
    Listener,
    EventsOnMultiple,
    EventsOn,
    EventsOnce,
    EventsNotify,
    EventsEmit,
    EventsOff,
    EventsOffAll
};
`;

const fn = new Function('mockWindow', wrappedCode);
const api = fn(mockWindow);

const { eventListeners, EventsOnMultiple, EventsOn, EventsOnce, EventsNotify, EventsEmit, EventsOff, EventsOffAll } = api;

let passed = 0;
let failed = 0;

function assert(condition, message) {
    if (condition) {
        console.log('PASS:', message);
        passed++;
    } else {
        console.log('FAIL:', message);
        failed++;
    }
}

function reset() {
    const keys = Object.keys(eventListeners);
    for (const k of keys) delete eventListeners[k];
}

// Test 1: Removing a listener via cancel function during callback
{
    reset();
    const cancel1 = EventsOn('test1', () => {});
    const cancel2 = EventsOn('test1', () => {
        cancel1();
    });
    EventsEmit('test1');
    const listeners = eventListeners['test1'];
    assert(listeners !== undefined && listeners.length === 1,
        'Test 1: After cancelling listener1 in listener2 callback, only 1 listener remains');
    cancel2();
}

// Test 2: Listener removed via cancel stays removed after multiple emits
{
    reset();
    let callCount = 0;
    const cancel1 = EventsOn('test2', () => { callCount++; });
    const cancel2 = EventsOn('test2', () => {
        cancel1();
    });
    EventsEmit('test2');
    EventsEmit('test2');
    assert(callCount === 1,
        'Test 2: Cancelled listener called exactly once');
}

// Test 3: EventsOnce still works
{
    reset();
    let onceCount = 0;
    EventsOnMultiple('test3', () => { onceCount++; }, 1);
    EventsEmit('test3');
    EventsEmit('test3');
    assert(onceCount === 1, 'Test 3: EventsOnce called exactly once');
    assert(eventListeners['test3'] === undefined,
        'Test 3: Event cleaned up after maxCallbacks');
}

// Test 4: Multiple listeners, mixed removal
{
    reset();
    let count1 = 0, count2 = 0, count3 = 0;
    EventsOnMultiple('test4', () => { count1++; }, 1);
    const cancel2 = EventsOn('test4', () => { count2++; });
    EventsOn('test4', () => {
        count3++;
        cancel2();
    });
    EventsEmit('test4');
    assert(count1 === 1, 'Test 4a: maxCallbacks=1 listener called once');
    assert(count2 === 1, 'Test 4b: cancelled listener called once');
    assert(count3 === 1, 'Test 4c: third listener called once');
    const remaining = eventListeners['test4'];
    assert(remaining !== undefined && remaining.length === 1,
        'Test 4d: Only 1 listener remains');
}

// Test 5: EventsNotify also handles removal correctly
{
    reset();
    let count = 0;
    const cancel = EventsOn('test5', () => { count++; });
    EventsOn('test5', () => {
        cancel();
    });
    EventsNotify(JSON.stringify({ name: 'test5', data: [] }));
    EventsNotify(JSON.stringify({ name: 'test5', data: [] }));
    assert(count === 1,
        'Test 5: Listener cancelled during EventsNotify stays removed');
}

console.log();
console.log(passed + ' passed, ' + failed + ' failed');
if (failed === 0) {
    console.log('ALL TESTS PASSED');
} else {
    console.log('SOME TESTS FAILED');
    process.exit(1);
}
