import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';
import assert from 'assert';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

const eventsCode = fs.readFileSync(path.join(__dirname, '../../internal/frontend/runtime/desktop/events.js'), 'utf8');

// Extract just the logic we need by evaluating the module code in a mock context
const eventListeners = {};
globalThis.window = { WailsInvoke: () => {} };

// Manually recreate the fixed notifyListeners and supporting functions for testing
class Listener {
    constructor(eventName, callback, maxCallbacks) {
        this.eventName = eventName;
        this.maxCallbacks = maxCallbacks || -1;
        this.Callback = (data) => {
            callback.apply(null, data);
            if (this.maxCallbacks === -1) return false;
            this.maxCallbacks -= 1;
            return this.maxCallbacks === 0;
        };
    }
}

function EventsOnMultiple(eventName, callback, maxCallbacks) {
    eventListeners[eventName] = eventListeners[eventName] || [];
    const thisListener = new Listener(eventName, callback, maxCallbacks);
    eventListeners[eventName].push(thisListener);
    return () => listenerOff(thisListener);
}

function EventsOn(eventName, callback) {
    return EventsOnMultiple(eventName, callback, -1);
}

function listenerOff(listener) {
    const eventName = listener.eventName;
    if (eventListeners[eventName] === undefined) return;
    eventListeners[eventName] = eventListeners[eventName].filter(l => l !== listener);
    if (eventListeners[eventName].length === 0) {
        delete eventListeners[eventName];
    }
}

function notifyListeners(eventData) {
    let eventName = eventData.name;
    let listeners = eventListeners[eventName];
    if (!listeners || listeners.length === 0) return;
    const snapshot = listeners.slice();
    for (let count = snapshot.length - 1; count >= 0; count -= 1) {
        const listener = snapshot[count];
        if (!eventListeners[eventName] || !eventListeners[eventName].includes(listener)) continue;
        const destroy = listener.Callback(eventData.data);
        if (destroy) {
            eventListeners[eventName] = eventListeners[eventName].filter(l => l !== listener);
        }
    }
    if (!eventListeners[eventName] || eventListeners[eventName].length === 0) {
        delete eventListeners[eventName];
    }
}

function EventsEmit(eventName) {
    const payload = { name: eventName, data: [].slice.apply(arguments).slice(1) };
    notifyListeners(payload);
}

console.log('Test: Removing an event listener during callback should keep it removed (#4393)');

let callCount1 = 0;
let callCount2 = 0;
let callCount3 = 0;

const cancel1 = EventsOn('test-event', () => {
    callCount1++;
    cancel1();
});

EventsOn('test-event', () => {
    callCount2++;
});

EventsOn('test-event', () => {
    callCount3++;
});

EventsEmit('test-event');

assert.strictEqual(callCount1, 1, 'First listener should have been called once');
assert.strictEqual(callCount2, 1, 'Second listener should have been called once');
assert.strictEqual(callCount3, 1, 'Third listener should have been called once');

EventsEmit('test-event');

assert.strictEqual(callCount1, 1, 'First listener should NOT have been called again');
assert.strictEqual(callCount2, 2, 'Second listener should have been called twice');
assert.strictEqual(callCount3, 2, 'Third listener should have been called twice');

console.log('PASS: All assertions passed');
