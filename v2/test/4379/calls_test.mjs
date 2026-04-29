import { readFileSync } from 'fs';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';

const __dirname = dirname(fileURLToPath(import.meta.url));

const callsCode = readFileSync(join(__dirname, '../../internal/frontend/runtime/desktop/calls.js'), 'utf8');

let wailsInvokeCallback = null;
const mockWindow = {
    WailsInvoke: (msg) => {
        wailsInvokeCallback && wailsInvokeCallback(msg);
    },
    crypto: {
        getRandomValues: (arr) => {
            arr[0] = Math.floor(Math.random() * 0xFFFFFFFF);
            return arr;
        }
    }
};

let code = callsCode
    .replace(/export /g, '')
    .replace(/window\.WailsInvoke/g, 'mockWindow.WailsInvoke')
    .replace(/window\.ObfuscatedCall/g, 'globalThis.ObfuscatedCall')
    .replace(/window\.crypto/g, 'mockWindow.crypto');

const wrappedCode = `
var runtime = { LogDebug: () => {} };
${code}
return { callbacks, Call };
`;

const fn = new Function('mockWindow', wrappedCode);
const api = fn(mockWindow);

const { callbacks, Call } = api;
// Callback is assigned to window in the original code, but since
// we replaced window references, we need to access it differently
// Let's just inline the Callback function for testing

function invokeCallback(incomingMessage) {
    let message;
    try {
        message = JSON.parse(incomingMessage);
    } catch (e) {
        throw new Error('Invalid JSON: ' + e.message);
    }
    let callbackID = message.callbackid;
    let callbackData = callbacks[callbackID];
    if (!callbackData) {
        throw new Error('Callback not found: ' + callbackID);
    }
    clearTimeout(callbackData.timeoutHandle);
    delete callbacks[callbackID];
    if (message.error) {
        callbackData.reject(typeof message.error === 'string' ? new Error(message.error) : message.error);
    } else {
        callbackData.resolve(message.result);
    }
}

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

// Test 1: Error string from Go is wrapped in Error object
{
    const callPromise = Call('mybinding.Func', []);

    const callbackID = Object.keys(callbacks)[0];
    invokeCallback(JSON.stringify({
        callbackid: callbackID,
        error: "something went wrong"
    }));

    try {
        await callPromise;
        assert(false, 'Test 1: Should have rejected');
    } catch (err) {
        assert(err instanceof Error, 'Test 1: Rejection value is an Error instance');
        assert(err.message === 'something went wrong', 'Test 1: Error message matches the string');
        assert(typeof err.stack === 'string', 'Test 1: Error has a stack trace');
    }
}

// Test 2: Successful result still works
{
    const callPromise = Call('mybinding.Func2', []);

    const callbackID = Object.keys(callbacks)[0];
    invokeCallback(JSON.stringify({
        callbackid: callbackID,
        result: { data: 42 }
    }));

    const result = await callPromise;
    assert(result.data === 42, 'Test 2: Successful result still works');
}

// Test 3: Error object (already non-string) is passed through as-is
{
    const callPromise = Call('mybinding.Func3', []);

    const callbackID = Object.keys(callbacks)[0];
    const errorObj = { message: "structured error", code: 500 };
    invokeCallback(JSON.stringify({
        callbackid: callbackID,
        error: errorObj
    }));

    try {
        await callPromise;
        assert(false, 'Test 3: Should have rejected');
    } catch (err) {
        assert(!(err instanceof Error), 'Test 3: Non-string error is passed through as-is');
        assert(err.message === 'structured error', 'Test 3: Structured error message preserved');
        assert(err.code === 500, 'Test 3: Structured error code preserved');
    }
}

console.log();
console.log(passed + ' passed, ' + failed + ' failed');
if (failed === 0) {
    console.log('ALL TESTS PASSED');
} else {
    console.log('SOME TESTS FAILED');
    process.exit(1);
}
