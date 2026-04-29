const eventsJS = `
var eventListeners = {};
var window = { WailsInvoke: () => {} };

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
    const newEventListenerList = eventListeners[eventName]?.slice() || [];
    if (newEventListenerList.length) {
        for (let count = newEventListenerList.length - 1; count >= 0; count -= 1) {
            const listener = newEventListenerList[count];
            let data = eventData.data;
            const destroy = listener.Callback(data);
            if (destroy) {
                newEventListenerList.splice(count, 1);
            }
        }
        const currentListeners = eventListeners[eventName];
        if (currentListeners) {
            const survivingListeners = newEventListenerList.filter(
                l => currentListeners.includes(l)
            );
            if (survivingListeners.length === 0) {
                delete eventListeners[eventName];
            } else {
                eventListeners[eventName] = survivingListeners;
            }
        } else {
            delete eventListeners[eventName];
        }
    }
}
`;

function runTest(name, script, assert) {
    const vm = require('vm');
    const context = vm.createContext({});
    vm.runInContext(eventsJS, context);
    vm.runInContext(script, context);
    const result = vm.runInContext(assert, context);
    if (result !== true) {
        throw new Error(`${name}: assertion failed: ${result}`);
    }
    console.log(`  ✓ ${name}`);
}

let passed = 0;
let failed = 0;

function test(name, script, assert) {
    try {
        runTest(name, script, assert);
        passed++;
    } catch (e) {
        console.error(`  ✗ ${name}: ${e.message}`);
        failed++;
    }
}

console.log('Testing event listener removal during callback handling:');

test(
    'listener removed via cancel() during callback stays removed',
    `
    let removed = false;
    let cancel = EventsOn("test", () => {
        if (!removed) {
            removed = true;
            cancel();
        }
    });
    notifyListeners({name: "test", data: []});
    var count1 = eventListeners["test"] ? eventListeners["test"].length : 0;
    notifyListeners({name: "test", data: []});
    var count2 = eventListeners["test"] ? eventListeners["test"].length : 0;
    `,
    `(function() {
        if (count1 !== 0) return "expected 0 listeners after first emit, got " + count1;
        if (count2 !== 0) return "expected 0 listeners after second emit, got " + count2;
        return true;
    })()`
);

test(
    'other listeners survive when one removes itself',
    `
    let results = [];
    let cancel1 = EventsOn("multi", () => { results.push("a"); cancel1(); });
    EventsOn("multi", () => { results.push("b"); });
    notifyListeners({name: "multi", data: []});
    var countAfterFirst = eventListeners["multi"] ? eventListeners["multi"].length : 0;
    notifyListeners({name: "multi", data: []});
    var countAfterSecond = eventListeners["multi"] ? eventListeners["multi"].length : 0;
    `,
    `(function() {
        if (countAfterFirst !== 1) return "expected 1 listener after first emit, got " + countAfterFirst;
        if (countAfterSecond !== 1) return "expected 1 listener after second emit, got " + countAfterSecond;
        if (JSON.stringify(results) !== '["b","a","b"]') return "unexpected results: " + JSON.stringify(results);
        return true;
    })()`
);

test(
    'all listeners removed via EventsOff during callback',
    `
    let results = [];
    EventsOn("offtest", () => { results.push(1); });
    EventsOn("offtest", () => { results.push(2); delete eventListeners["offtest"]; });
    EventsOn("offtest", () => { results.push(3); });
    var initialCount = eventListeners["offtest"].length;
    notifyListeners({name: "offtest", data: []});
    var finalCount = eventListeners["offtest"] !== undefined ? eventListeners["offtest"].length : 0;
    `,
    `(function() {
        if (initialCount !== 3) return "expected 3 initial listeners, got " + initialCount;
        if (finalCount !== 0) return "expected 0 listeners after off during callback, got " + finalCount;
        return true;
    })()`
);

test(
    'self-removing listener with EventsOnce works correctly',
    `
    let results = [];
    EventsOnMultiple("once", () => { results.push("once"); }, 1);
    EventsOn("once", () => { results.push("persistent"); });
    notifyListeners({name: "once", data: []});
    var count1 = eventListeners["once"].length;
    notifyListeners({name: "once", data: []});
    var count2 = eventListeners["once"].length;
    `,
    `(function() {
        if (count1 !== 1) return "expected 1 listener after first emit, got " + count1;
        if (count2 !== 1) return "expected 1 listener after second emit, got " + count2;
        if (JSON.stringify(results) !== '["persistent","once","persistent"]') return "unexpected results: " + JSON.stringify(results);
        return true;
    })()`
);

console.log(`\nResults: ${passed} passed, ${failed} failed`);
process.exit(failed > 0 ? 1 : 0);
