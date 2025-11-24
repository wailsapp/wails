/**
 * @jest-environment jsdom
 */

import { Call, ByName, ByID, BindingError, setBindingTimeout, getBindingTimeout } from '../calls';

// Mock fetch globally
global.fetch = jest.fn();

// Mock AbortController if not available
if (!global.AbortController) {
    global.AbortController = class AbortController {
        signal = { aborted: false };
        abort() {
            this.signal.aborted = true;
        }
    } as any;
}

describe('HTTP-Only Bindings', () => {
    beforeEach(() => {
        (fetch as jest.Mock).mockClear();
        // Reset timeout to default
        setBindingTimeout(5 * 60 * 1000);
    });

    describe('Call function', () => {
        it('should make HTTP request for successful binding call', async () => {
            const mockResponse = {
                ok: true,
                status: 200,
                headers: new Map([['Content-Type', 'application/json']]),
                json: jest.fn().mockResolvedValue({ result: 'success' })
            };
            (fetch as jest.Mock).mockResolvedValue(mockResponse);

            const result = await Call({
                methodName: 'testMethod',
                args: ['arg1', 'arg2']
            });

            expect(fetch).toHaveBeenCalledWith(
                expect.stringContaining('/wails/runtime'),
                expect.objectContaining({
                    headers: expect.objectContaining({
                        'x-wails-client-id': expect.any(String)
                    })
                })
            );

            expect(result).toEqual({ result: 'success' });
        });

        it('should handle method not found error (404)', async () => {
            const mockErrorResponse = {
                ok: false,
                status: 404,
                json: jest.fn().mockResolvedValue({
                    error: 'Method not found',
                    kind: 'ReferenceError'
                })
            };
            (fetch as jest.Mock).mockResolvedValue(mockErrorResponse);

            await expect(Call({
                methodName: 'nonexistentMethod',
                args: []
            })).rejects.toThrow(BindingError);

            try {
                await Call({ methodName: 'nonexistentMethod', args: [] });
            } catch (error) {
                expect(error).toBeInstanceOf(BindingError);
                expect((error as BindingError).status).toBe(404);
                expect((error as BindingError).kind).toBe('ReferenceError');
                expect((error as BindingError).message).toBe('Method not found');
            }
        });

        it('should handle runtime error (500)', async () => {
            const mockErrorResponse = {
                ok: false,
                status: 500,
                json: jest.fn().mockResolvedValue({
                    error: 'Internal server error',
                    kind: 'RuntimeError'
                })
            };
            (fetch as jest.Mock).mockResolvedValue(mockErrorResponse);

            await expect(Call({
                methodName: 'errorMethod',
                args: []
            })).rejects.toThrow(BindingError);

            try {
                await Call({ methodName: 'errorMethod', args: [] });
            } catch (error) {
                expect(error).toBeInstanceOf(BindingError);
                expect((error as BindingError).status).toBe(500);
                expect((error as BindingError).kind).toBe('RuntimeError');
            }
        });

        it('should handle timeout error (408)', async () => {
            const mockErrorResponse = {
                ok: false,
                status: 408,
                json: jest.fn().mockResolvedValue({
                    error: 'Request timeout',
                    kind: 'TimeoutError'
                })
            };
            (fetch as jest.Mock).mockResolvedValue(mockErrorResponse);

            await expect(Call({
                methodName: 'slowMethod',
                args: []
            })).rejects.toThrow(BindingError);

            try {
                await Call({ methodName: 'slowMethod', args: [] });
            } catch (error) {
                expect(error).toBeInstanceOf(BindingError);
                expect((error as BindingError).status).toBe(408);
                expect((error as BindingError).kind).toBe('TimeoutError');
            }
        });

        it('should handle network errors', async () => {
            (fetch as jest.Mock).mockRejectedValue(new Error('Network error'));

            await expect(Call({
                methodName: 'testMethod',
                args: []
            })).rejects.toThrow('Network error');
        });

        it('should handle non-JSON error responses', async () => {
            const mockErrorResponse = {
                ok: false,
                status: 500,
                statusText: 'Internal Server Error',
                json: jest.fn().mockRejectedValue(new Error('Not JSON'))
            };
            (fetch as jest.Mock).mockResolvedValue(mockErrorResponse);

            await expect(Call({
                methodName: 'testMethod',
                args: []
            })).rejects.toThrow(BindingError);

            try {
                await Call({ methodName: 'testMethod', args: [] });
            } catch (error) {
                expect(error).toBeInstanceOf(BindingError);
                expect((error as BindingError).status).toBe(500);
                expect((error as BindingError).kind).toBe('HttpError');
                expect((error as BindingError).message).toBe('HTTP 500: Internal Server Error');
            }
        });
    });

    describe('ByName function', () => {
        it('should call binding by method name', async () => {
            const mockResponse = {
                ok: true,
                status: 200,
                headers: new Map([['Content-Type', 'application/json']]),
                json: jest.fn().mockResolvedValue('method result')
            };
            (fetch as jest.Mock).mockResolvedValue(mockResponse);

            const result = await ByName('testMethod', 'arg1', 'arg2');

            expect(fetch).toHaveBeenCalledWith(
                expect.stringContaining('methodName'),
                expect.any(Object)
            );
            expect(result).toBe('method result');
        });
    });

    describe('ByID function', () => {
        it('should call binding by method ID', async () => {
            const mockResponse = {
                ok: true,
                status: 200,
                headers: new Map([['Content-Type', 'application/json']]),
                json: jest.fn().mockResolvedValue({ id: 42 })
            };
            (fetch as jest.Mock).mockResolvedValue(mockResponse);

            const result = await ByID(123, 'arg1', 'arg2');

            expect(fetch).toHaveBeenCalledWith(
                expect.stringContaining('methodID'),
                expect.any(Object)
            );
            expect(result).toEqual({ id: 42 });
        });
    });

    describe('Timeout configuration', () => {
        it('should set and get binding timeout', () => {
            const timeout = 10 * 60 * 1000; // 10 minutes
            setBindingTimeout(timeout);
            expect(getBindingTimeout()).toBe(timeout);
        });

        it('should use custom timeout in requests', async () => {
            const customTimeout = 30 * 1000; // 30 seconds
            setBindingTimeout(customTimeout);

            const mockResponse = {
                ok: true,
                status: 200,
                headers: new Map([['Content-Type', 'application/json']]),
                json: jest.fn().mockResolvedValue('success')
            };
            (fetch as jest.Mock).mockResolvedValue(mockResponse);

            await Call({ methodName: 'testMethod', args: [] });

            // Verify that timeout was used in the request options
            expect(fetch).toHaveBeenCalledWith(
                expect.any(String),
                expect.objectContaining({
                    timeout: customTimeout
                })
            );
        });
    });

    describe('Cancellation', () => {
        it('should support cancellation with AbortController', async () => {
            // Mock a slow response
            const mockResponse = new Promise((resolve) => {
                setTimeout(() => {
                    resolve({
                        ok: true,
                        status: 200,
                        headers: new Map([['Content-Type', 'application/json']]),
                        json: jest.fn().mockResolvedValue('success')
                    });
                }, 1000);
            });
            (fetch as jest.Mock).mockReturnValue(mockResponse);

            const promise = Call({ methodName: 'slowMethod', args: [] });

            // Cancel the promise
            setTimeout(() => {
                promise.cancel();
            }, 100);

            await expect(promise).rejects.toThrow('Binding call cancelled');
        });

        it('should handle AbortError from fetch', async () => {
            const abortError = new Error('The operation was aborted');
            abortError.name = 'AbortError';
            (fetch as jest.Mock).mockRejectedValue(abortError);

            const promise = Call({ methodName: 'testMethod', args: [] });

            await expect(promise).rejects.toThrow('Binding call cancelled');
        });
    });

    describe('Large data handling', () => {
        it('should handle large JSON responses', async () => {
            // Create a large response object
            const largeData = {};
            for (let i = 0; i < 10000; i++) {
                largeData[`key_${i}`] = `value_${i}`;
            }

            const mockResponse = {
                ok: true,
                status: 200,
                headers: new Map([['Content-Type', 'application/json']]),
                json: jest.fn().mockResolvedValue(largeData)
            };
            (fetch as jest.Mock).mockResolvedValue(mockResponse);

            const result = await Call({ methodName: 'largeDataMethod', args: [] });

            expect(Object.keys(result)).toHaveLength(10000);
            expect(result['key_0']).toBe('value_0');
            expect(result['key_9999']).toBe('value_9999');
        });

        it('should handle text responses', async () => {
            const mockResponse = {
                ok: true,
                status: 200,
                headers: new Map([['Content-Type', 'text/plain']]),
                json: jest.fn().mockRejectedValue(new Error('Not JSON')),
                text: jest.fn().mockResolvedValue('Plain text response')
            };
            (fetch as jest.Mock).mockResolvedValue(mockResponse);

            const result = await Call({ methodName: 'textMethod', args: [] });

            expect(result).toBe('Plain text response');
        });
    });

    describe('BindingError class', () => {
        it('should create BindingError with all properties', () => {
            const error = new BindingError(404, 'ReferenceError', 'Method not found', { details: 'extra' });

            expect(error.name).toBe('BindingError');
            expect(error.status).toBe(404);
            expect(error.kind).toBe('ReferenceError');
            expect(error.message).toBe('Method not found');
            expect(error.cause).toEqual({ details: 'extra' });
            expect(error).toBeInstanceOf(Error);
        });

        it('should work without cause parameter', () => {
            const error = new BindingError(500, 'RuntimeError', 'Something went wrong');

            expect(error.name).toBe('BindingError');
            expect(error.status).toBe(500);
            expect(error.kind).toBe('RuntimeError');
            expect(error.message).toBe('Something went wrong');
            expect(error.cause).toBeUndefined();
        });
    });

    describe('Request parameters', () => {
        it('should include call-id in request', async () => {
            const mockResponse = {
                ok: true,
                status: 200,
                headers: new Map([['Content-Type', 'application/json']]),
                json: jest.fn().mockResolvedValue('success')
            };
            (fetch as jest.Mock).mockResolvedValue(mockResponse);

            await Call({ methodName: 'testMethod', args: [] });

            const fetchCall = (fetch as jest.Mock).mock.calls[0];
            const url = new URL(fetchCall[0]);
            const args = JSON.parse(url.searchParams.get('args') || '{}');

            expect(args['call-id']).toBeTruthy();
            expect(typeof args['call-id']).toBe('string');
        });

        it('should include method name in request', async () => {
            const mockResponse = {
                ok: true,
                status: 200,
                headers: new Map([['Content-Type', 'application/json']]),
                json: jest.fn().mockResolvedValue('success')
            };
            (fetch as jest.Mock).mockResolvedValue(mockResponse);

            await Call({ methodName: 'testMethod', args: ['arg1'] });

            const fetchCall = (fetch as jest.Mock).mock.calls[0];
            const url = new URL(fetchCall[0]);
            const args = JSON.parse(url.searchParams.get('args') || '{}');

            expect(args.methodName).toBe('testMethod');
            expect(args.args).toEqual(['arg1']);
        });

        it('should include method ID in request', async () => {
            const mockResponse = {
                ok: true,
                status: 200,
                headers: new Map([['Content-Type', 'application/json']]),
                json: jest.fn().mockResolvedValue('success')
            };
            (fetch as jest.Mock).mockResolvedValue(mockResponse);

            await Call({ methodID: 42, args: ['arg1'] });

            const fetchCall = (fetch as jest.Mock).mock.calls[0];
            const url = new URL(fetchCall[0]);
            const args = JSON.parse(url.searchParams.get('args') || '{}');

            expect(args.methodID).toBe(42);
            expect(args.args).toEqual(['arg1']);
        });
    });
});