const vm = require('vm');

const callsJS = `
var callbacks = {};
var window = { WailsInvoke: () => {} };
var clearTimeout = function() {};

function cryptoRandom() {
    var array = new Uint32Array(1);
    return window.crypto.getRandomValues(array)[0];
}

var randomFunc = cryptoRandom;

function Call(name, args, timeout) {
    if (timeout == null) timeout = 0;
    return new Promise(function(resolve, reject) {
        var callbackID;
        do { callbackID = name + '-' + randomFunc(); } while (callbacks[callbackID]);
        var timeoutHandle;
        if (timeout > 0) {
            timeoutHandle = setTimeout(function() {
                reject(Error('Call to ' + name + ' timed out'));
            }, timeout);
        }
        callbacks[callbackID] = { timeoutHandle, reject, resolve };
        try {
            const payload = { name, args, callbackID };
            window.WailsInvoke('C' + JSON.stringify(payload));
        } catch (e) { console.error(e); }
    });
}

function Callback(incomingMessage) {
    let message;
    try { message = JSON.parse(incomingMessage); }
    catch (e) { throw new Error('Invalid JSON: ' + incomingMessage); }
    let callbackID = message.callbackid;
    let callbackData = callbacks[callbackID];
    if (!callbackData) throw new Error('Callback not found: ' + callbackID);
    clearTimeout(callbackData.timeoutHandle);
    delete callbacks[callbackID];
    if (message.error) {
        callbackData.reject(typeof message.error === 'string' ? new Error(message.error) : message.error);
    } else {
        callbackData.resolve(message.result);
    }
}

// Simulate Call + Callback for testing
function testCall(name, errorMessage, resultValue) {
    return new Promise(function(resolve, reject) {
        var callbackID = name + '-test1';
        callbacks[callbackID] = {
            timeoutHandle: null,
            reject: function(err) { resolve({ rejected: true, error: err }); },
            resolve: function(val) { resolve({ resolved: true, value: val }); }
        };
        if (errorMessage) {
            Callback(JSON.stringify({ callbackid: callbackID, error: errorMessage }));
        } else {
            Callback(JSON.stringify({ callbackid: callbackID, result: resultValue }));
        }
    });
}
`;

async function runTests() {
    const context = vm.createContext({});
    vm.runInContext(callsJS + `
        async function testStringError() {
            var result = await testCall("method1", "something went wrong", null);
            if (!result.rejected) throw new Error("expected rejection");
            if (typeof result.error !== 'object' || !(result.error instanceof Error)) {
                throw new Error("expected Error object, got " + typeof result.error);
            }
            if (result.error.message !== "something went wrong") {
                throw new Error("expected message 'something went wrong', got '" + result.error.message + "'");
            }
            return "pass: string error wrapped in Error object";
        }

        async function testObjectError() {
            var customErr = { code: 42, message: "custom" };
            var result = await testCall("method2", customErr, null);
            if (!result.rejected) throw new Error("expected rejection");
            if (typeof result.error !== 'object') {
                throw new Error("expected object error, got " + typeof result.error);
            }
            if (result.error.code !== 42) {
                throw new Error("expected code 42, got " + result.error.code);
            }
            return "pass: object error passed through unchanged";
        }

        async function testSuccess() {
            var result = await testCall("method3", null, 42);
            if (!result.resolved) throw new Error("expected resolution");
            if (result.value !== 42) throw new Error("expected 42, got " + result.value);
            return "pass: success case works";
        }
    `, context);

    const tests = ['testStringError', 'testObjectError', 'testSuccess'];
    let passed = 0, failed = 0;
    for (const name of tests) {
        try {
            const msg = await vm.runInContext(name + '()', context);
            console.log('  ✓ ' + msg);
            passed++;
        } catch (e) {
            console.error('  ✗ ' + name + ': ' + e.message);
            failed++;
        }
    }
    console.log(`\nResults: ${passed} passed, ${failed} failed`);
    process.exit(failed > 0 ? 1 : 0);
}

runTests();
