import { Call, Callback, callbacks } from './calls'
import { expect, describe, it, beforeAll, vi, afterEach } from 'vitest'

beforeAll(() => {
    window.WailsInvoke = vi.fn(() => {})
    window.runtime = {
        LogDebug: vi.fn(),
    }
})

afterEach(() => {
    vi.clearAllMocks()
    Object.keys(callbacks).forEach(key => delete callbacks[key])
})

describe('Callback', () => {
    it('should reject with Error object when binding returns error (issue #4379)', async () => {
        const promise = Call('main/App.SetHello', ['test'], 0)

        const invokeCall = window.WailsInvoke.mock.calls[0][0]
        const payload = JSON.parse(invokeCall.slice(1))
        const callbackID = payload.callbackID

        Callback(JSON.stringify({
            callbackid: callbackID,
            error: "some error message"
        }))

        try {
            await promise
            expect.unreachable('should have rejected')
        } catch (e) {
            expect(e).toBeInstanceOf(Error)
            expect(e.message).toBe('some error message')
        }
    })

    it('should resolve with result when binding succeeds', async () => {
        const promise = Call('main/App.GetValue', [], 0)

        const invokeCall = window.WailsInvoke.mock.calls[0][0]
        const payload = JSON.parse(invokeCall.slice(1))
        const callbackID = payload.callbackID

        Callback(JSON.stringify({
            callbackid: callbackID,
            result: "hello world"
        }))

        const result = await promise
        expect(result).toBe("hello world")
    })

    it('should reject with Error on timeout', async () => {
        const promise = Call('main/App.SlowCall', [], 50)

        try {
            await promise
            expect.unreachable('should have rejected')
        } catch (e) {
            expect(e).toBeInstanceOf(Error)
            expect(e.message).toContain('timed out')
        }
    })

    it('should reject with Error object for ObfuscatedCall errors', async () => {
        const promise = window.ObfuscatedCall('abc123', [], 0)

        const invokeCall = window.WailsInvoke.mock.calls[0][0]
        const payload = JSON.parse(invokeCall.slice(1))
        const callbackID = payload.callbackID

        Callback(JSON.stringify({
            callbackid: callbackID,
            error: "obfuscated error"
        }))

        try {
            await promise
            expect.unreachable('should have rejected')
        } catch (e) {
            expect(e).toBeInstanceOf(Error)
            expect(e.message).toBe('obfuscated error')
        }
    })
})
