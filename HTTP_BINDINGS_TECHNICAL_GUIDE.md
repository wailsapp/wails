# Wails v3 HTTP-Only Bindings Technical Guide

## Executive Summary

This document provides a comprehensive technical analysis of Wails v3's current bindings architecture and proposes a migration path to a fully HTTP-based protocol, eliminating the use of JavaScript `eval()` for returning results from Go method calls.

**Current State**: HTTP Request → Go Handler → JavaScript eval() for responses
**Target State**: HTTP Request → Go Handler → HTTP Response (full round-trip)

---

## Table of Contents

1. [Current Architecture (BEFORE State)](#current-architecture-before-state)
2. [Cancellation Mechanism](#cancellation-mechanism)
3. [Proposed HTTP-Only Architecture (AFTER State)](#proposed-http-only-architecture-after-state)
4. [Migration Implementation Guide](#migration-implementation-guide)
5. [Trade-offs and Considerations](#trade-offs-and-considerations)

---

## Current Architecture (BEFORE State)

### Overview

The current Wails v3-alpha bindings use a **hybrid HTTP + JavaScript eval()** approach:
- **Requests** are sent via HTTP
- **Responses** are returned via JavaScript execution (`eval`)

### Component Flow Diagram

```
┌─────────────┐
│  Frontend   │
│  (Browser)  │
└──────┬──────┘
       │
       │ 1. HTTP POST Request
       │    /wails/runtime?object=0&method=0
       │    Headers: x-wails-client-id, x-wails-window-name
       │    Body: JSON { call-id, methodID/methodName, args }
       │
       ▼
┌──────────────────────────────────────────────────┐
│  AssetServer (assetserver.go)                     │
│  ┌──────────────────────────────────────────┐    │
│  │  MessageProcessor (messageprocessor.go)   │    │
│  │  ServeHTTP() → HandleRuntimeCallWithIDs() │    │
│  └──────────┬───────────────────────────────┘    │
│             │                                      │
│             ▼                                      │
│  ┌──────────────────────────────────────────┐    │
│  │ processCallMethod()                       │    │
│  │ (messageprocessor_call.go:63)             │    │
│  └──────────┬───────────────────────────────┘    │
└─────────────┼──────────────────────────────────────┘
              │
              │ 2. Immediate HTTP 200 OK
              │    (empty response body)
              │
              ▼
       ┌──────────────┐
       │   Frontend   │
       │   Promise    │
       │   awaits...  │
       └──────────────┘

              │ 3. Go goroutine executes
              │    boundMethod.Call(ctx, args)
              │    (bindings.go:261)
              │
              ▼
┌──────────────────────────────────────────┐
│  Bound Go Method Execution                │
│  (User's Go code runs here)               │
└──────────┬───────────────────────────────┘
           │
           │ 4. Result or Error returned
           │
           ▼
┌──────────────────────────────────────────┐
│  callCallback() or callErrorCallback()    │
│  (messageprocessor_call.go:18-34)         │
└──────────┬───────────────────────────────┘
           │
           │ 5. window.ExecJS()
           │    Executes JavaScript via platform API
           │    (webview_window.go:320-343)
           │
           ▼
┌──────────────────────────────────────────┐
│  Platform-Specific JavaScript Injection   │
│  - Windows: ICoreWebView2.ExecuteScript   │
│  - Linux: webkit_web_view_evaluate_javascript │
│  - macOS: evaluateJavaScript             │
└──────────┬───────────────────────────────┘
           │
           │ 6. JavaScript executed in browser context
           │
           ▼
┌──────────────────────────────────────────┐
│  JavaScript Callbacks (calls.ts:72-137)   │
│  _wails.callResultHandler(id, data, isJSON) │
│  OR                                        │
│  _wails.callErrorHandler(id, data, isJSON) │
└──────────┬───────────────────────────────┘
           │
           │ 7. Promise resolved/rejected
           │
           ▼
       ┌──────────────┐
       │   Frontend   │
       │   receives   │
       │   result     │
       └──────────────┘
```

### Key Files and Components

#### Frontend (TypeScript/JavaScript)

**File**: `v3/internal/runtime/desktop/@wailsio/runtime/src/calls.ts`

```typescript
// Lines 177-208: Main call function
export function Call(options: CallOptions): CancellablePromise<any> {
    const id = generateID();

    const result = CancellablePromise.withResolvers<any>();
    callResponses.set(id, { resolve: result.resolve, reject: result.reject });

    // HTTP REQUEST sent here
    const request = call(CallBinding, Object.assign({ "call-id": id }, options));
    let running = false;

    request.then(() => {
        running = true;
    }, (err) => {
        callResponses.delete(id);
        result.reject(err);
    });

    const cancel = () => {
        callResponses.delete(id);
        return cancelCall(CancelMethod, {"call-id": id}).catch((err) => {
            console.error("Error while requesting binding call cancellation:", err);
        });
    };

    result.oncancelled = () => {
        if (running) {
            return cancel();
        } else {
            return request.then(cancel);
        }
    };

    return result.promise;
}

// Lines 72-88: Result handler (called via eval)
function resultHandler(id: string, data: string, isJSON: boolean): void {
    const resolvers = getAndDeleteResponse(id);
    if (!resolvers) {
        return;
    }

    if (!data) {
        resolvers.resolve(undefined);
    } else if (!isJSON) {
        resolvers.resolve(data);
    } else {
        try {
            resolvers.resolve(JSON.parse(data));
        } catch (err: any) {
            resolvers.reject(new TypeError("could not parse result: " + err.message, { cause: err }));
        }
    }
}

// Lines 98-137: Error handler (called via eval)
function errorHandler(id: string, data: string, isJSON: boolean): void {
    const resolvers = getAndDeleteResponse(id);
    if (!resolvers) {
        return;
    }

    if (!isJSON) {
        resolvers.reject(new Error(data));
    } else {
        let error: any;
        try {
            error = JSON.parse(data);
        } catch (err: any) {
            resolvers.reject(new TypeError("could not parse error: " + err.message, { cause: err }));
            return;
        }

        let options: ErrorOptions = {};
        if (error.cause) {
            options.cause = error.cause;
        }

        let exception;
        switch (error.kind) {
            case "ReferenceError":
                exception = new ReferenceError(error.message, options);
                break;
            case "TypeError":
                exception = new TypeError(error.message, options);
                break;
            case "RuntimeError":
                exception = new RuntimeError(error.message, options);
                break;
            default:
                exception = new Error(error.message, options);
                break;
        }

        resolvers.reject(exception);
    }
}
```

**How it's used** (Generated bindings):
```javascript
// v3/examples/binding/assets/bindings/.../greetservice.js
export function Greet(name, ...counts) {
    return $Call.ByID(1411160069, name, counts);
}
```

#### Backend (Go)

**File**: `v3/pkg/application/messageprocessor_call.go`

```go
// Lines 63-195: Main call processing
func (m *MessageProcessor) processCallMethod(method int, rw http.ResponseWriter, r *http.Request, window Window, params QueryParams) {
	args, err := params.Args()
	if err != nil {
		m.httpError(rw, "Invalid binding call:", fmt.Errorf("unable to parse arguments: %w", err))
		return
	}

	callID := args.String("call-id")
	if callID == nil || *callID == "" {
		m.httpError(rw, "Invalid binding call:", errors.New("missing argument 'call-id'"))
		return
	}

	switch method {
	case CallBinding:
		var options CallOptions
		err := params.ToStruct(&options)
		if err != nil {
			m.httpError(rw, "Invalid binding call:", fmt.Errorf("error parsing call options: %w", err))
			return
		}

		ctx, cancel := context.WithCancel(context.WithoutCancel(r.Context()))

		// Schedule cancel in case panics happen before starting the call.
		cancelRequired := true
		defer func() {
			if cancelRequired {
				cancel()
			}
		}()

		ambiguousID := false
		func() {
			m.l.Lock()
			defer m.l.Unlock()

			if m.runningCalls[*callID] != nil {
				ambiguousID = true
			} else {
				m.runningCalls[*callID] = cancel
			}
		}()

		if ambiguousID {
			m.httpError(rw, "Invalid binding call:", fmt.Errorf("ambiguous call id: %s", *callID))
			return
		}

		m.ok(rw) // ⚠️ HTTP RESPONSE SENT IMMEDIATELY WITH STATUS 200

		// Log call
		var methodRef any = options.MethodName
		if options.MethodName == "" {
			methodRef = options.MethodID
		}
		m.Info("Binding call started:", "id", *callID, "method", methodRef)

		// ⚠️ ACTUAL EXECUTION HAPPENS IN GOROUTINE
		go func() {
			defer handlePanic()
			defer func() {
				m.l.Lock()
				defer m.l.Unlock()
				delete(m.runningCalls, *callID)
			}()
			defer cancel()

			var boundMethod *BoundMethod
			if options.MethodName != "" {
				boundMethod = globalApplication.bindings.Get(&options)
				if boundMethod == nil {
					m.callErrorCallback(window, "Binding call failed:", callID, &CallError{
						Kind:    ReferenceError,
						Message: fmt.Sprintf("unknown bound method name '%s'", options.MethodName),
					})
					return
				}
			} else {
				boundMethod = globalApplication.bindings.GetByID(options.MethodID)
				if boundMethod == nil {
					m.callErrorCallback(window, "Binding call failed:", callID, &CallError{
						Kind:    ReferenceError,
						Message: fmt.Sprintf("unknown bound method id %d", options.MethodID),
					})
					return
				}
			}

			// Prepare args for logging. This should never fail since json.Unmarshal succeeded before.
			jsonArgs, _ := json.Marshal(options.Args)
			var jsonResult []byte
			defer func() {
				m.Info("Binding call complete:", "id", *callID, "method", boundMethod, "args", string(jsonArgs), "result", string(jsonResult))
			}()

			// Set the context values for the window
			if window != nil {
				ctx = context.WithValue(ctx, WindowKey, window)
			}

			// ⚠️ ACTUAL GO METHOD EXECUTION
			result, err := boundMethod.Call(ctx, options.Args)
			if cerr := (*CallError)(nil); errors.As(err, &cerr) {
				switch cerr.Kind {
				case ReferenceError, TypeError:
					m.callErrorCallback(window, "Binding call failed:", callID, cerr)
				case RuntimeError:
					m.callErrorCallback(window, "Bound method returned an error:", callID, cerr)
				}
				return
			}

			if result != nil {
				// convert result to json
				jsonResult, err = json.Marshal(result)
				if err != nil {
					m.callErrorCallback(window, "Binding call failed:", callID, &CallError{
						Kind:    TypeError,
						Message: fmt.Sprintf("error marshaling result: %s", err),
					})
					return
				}
			}

			// ⚠️ RESULT DELIVERED VIA JAVASCRIPT EVAL
			m.callCallback(window, callID, string(jsonResult))
		}()

		cancelRequired = false

	default:
		m.httpError(rw, "Invalid binding call:", fmt.Errorf("unknown method: %d", method))
		return
	}
}

// Lines 32-34: Success callback - uses JavaScript eval
func (m *MessageProcessor) callCallback(window Window, callID *string, result string) {
	window.CallResponse(*callID, result)
}

// Lines 18-30: Error callback - uses JavaScript eval
func (m *MessageProcessor) callErrorCallback(window Window, message string, callID *string, err error) {
	m.Error(message, "id", *callID, "error", err)
	if cerr := (*CallError)(nil); errors.As(err, &cerr) {
		if data, jsonErr := json.Marshal(cerr); jsonErr == nil {
			window.CallError(*callID, string(data), true)
			return
		} else {
			m.Error("Unable to convert data to JSON. Please report this to the Wails team!", "id", *callID, "error", jsonErr)
		}
	}

	window.CallError(*callID, err.Error(), false)
}
```

**File**: `v3/pkg/application/webview_window.go`

```go
// Lines 320-343: JavaScript eval execution for responses
func (w *WebviewWindow) CallError(callID string, result string, isJSON bool) {
	if w.impl != nil {
		w.impl.execJS(
			fmt.Sprintf(
				"_wails.callErrorHandler('%s', '%s', %t);",
				callID,
				template.JSEscapeString(result),
				isJSON,
			),
		)
	}
}

func (w *WebviewWindow) CallResponse(callID string, result string) {
	if w.impl != nil {
		w.impl.execJS(
			fmt.Sprintf(
				"_wails.callResultHandler('%s', '%s', true);",
				callID,
				template.JSEscapeString(result),
			),
		)
	}
}
```

### Data Structures

#### Frontend Call Options
```typescript
type CallOptions = {
    methodID: number;
    methodName?: never;
    args: any[];
} | {
    methodID?: never;
    methodName: string;
    args: any[];
};
```

#### Backend Call Options
```go
type CallOptions struct {
	MethodID   uint32            `json:"methodID"`
	MethodName string            `json:"methodName"`
	Args       []json.RawMessage `json:"args"`
}
```

### Why eval() is Currently Required

1. **Asynchronous Goroutine Execution**: The HTTP response is sent immediately (line 112 in `messageprocessor_call.go`), while the actual Go method executes in a goroutine

2. **No HTTP Connection to Return To**: Once the HTTP response completes, there's no connection to send the result back through

3. **Window-Specific Responses**: Results must be delivered to the specific browser window that made the request

4. **Platform Integration**: Each platform (Windows/Linux/macOS) provides APIs to execute JavaScript in the webview context:
   - **Windows**: `ICoreWebView2::ExecuteScript`
   - **Linux**: `webkit_web_view_evaluate_javascript`
   - **macOS**: `evaluateJavaScript:completionHandler:`

---

## Cancellation Mechanism

### Current Implementation

**File**: `v3/pkg/application/messageprocessor_call.go`

```go
// Lines 36-61: Cancel request handler
func (m *MessageProcessor) processCallCancelMethod(method int, rw http.ResponseWriter, r *http.Request, window Window, params QueryParams) {
	args, err := params.Args()
	if err != nil {
		m.httpError(rw, "Invalid binding call:", fmt.Errorf("unable to parse arguments: %w", err))
		return
	}

	callID := args.String("call-id")
	if callID == nil || *callID == "" {
		m.httpError(rw, "Invalid binding call:", errors.New("missing argument 'call-id'"))
		return
	}

	var cancel func()
	func() {
		m.l.Lock()
		defer m.l.Unlock()
		cancel = m.runningCalls[*callID]
	}()

	if cancel != nil {
		cancel()
		m.Info("Binding call canceled:", "id", *callID)
	}
	m.ok(rw)
}
```

### How Cancellation Works

1. **Frontend Initiates Cancel**:
   ```typescript
   // calls.ts:193-197
   const cancel = () => {
       callResponses.delete(id);
       return cancelCall(CancelMethod, {"call-id": id}).catch((err) => {
           console.error("Error while requesting binding call cancellation:", err);
       });
   };
   ```

2. **HTTP Request Sent**: `POST /wails/runtime?object=10&method=0` with `{"call-id": "abc123"}`

3. **MessageProcessor Lookup**: Finds the cancel function stored in `runningCalls` map

4. **Context Cancellation**: Calls the `cancel()` function, which cancels the context

5. **Go Method Receives Cancellation**: The bound Go method can check `ctx.Done()` to detect cancellation

6. **Cleanup**: Cancel function is removed from map in the goroutine's defer

### Cancellation Flow Diagram

```
Frontend                MessageProcessor         Goroutine
   │                           │                      │
   │ cancel() called           │                      │
   ├──────────────────────────►│                      │
   │ HTTP POST                  │                      │
   │ object=10&method=0         │                      │
   │                            │                      │
   │◄───────────────────────────┤                      │
   │ 200 OK                     │                      │
   │                            │                      │
   │                            ├─────────────────────►│
   │                            │ ctx.Done() closed    │
   │                            │                      │
   │                            │                      ├─ Go method
   │                            │                      │  detects cancel
   │                            │                      │  and exits
   │                            │                      │
   │                            │◄─────────────────────┤
   │                            │ defer cleanup        │
   │                            │                      │
```

### Storage Mechanism

**MessageProcessor State**:
```go
type MessageProcessor struct {
	logger *slog.Logger

	runningCalls map[string]context.CancelFunc  // ← Stores cancel functions by call-id
	l            sync.Mutex                      // ← Protects the map
}
```

**Registration** (processCallMethod:96-105):
```go
func() {
	m.l.Lock()
	defer m.l.Unlock()

	if m.runningCalls[*callID] != nil {
		ambiguousID = true
	} else {
		m.runningCalls[*callID] = cancel  // ← Store cancel function
	}
}()
```

**Cleanup** (processCallMethod:124-127):
```go
defer func() {
	m.l.Lock()
	defer m.l.Unlock()
	delete(m.runningCalls, *callID)  // ← Remove after completion
}()
```

---

## Proposed HTTP-Only Architecture (AFTER State)

### Overview

Replace the eval-based response mechanism with **HTTP streaming** or **long-polling**:
- HTTP request opens
- Go method executes
- Result returned via HTTP response body
- HTTP connection closes

### Option 1: HTTP Streaming (Server-Sent Events)

#### Flow Diagram

```
┌─────────────┐
│  Frontend   │
└──────┬──────┘
       │
       │ 1. HTTP POST Request (kept alive)
       │    /wails/runtime?object=0&method=0
       │    Headers: Accept: text/event-stream
       │
       ▼
┌──────────────────────────────────────────────────┐
│  MessageProcessor                                 │
│  ┌──────────────────────────────────────────┐    │
│  │ processCallMethod()                       │    │
│  │ - Does NOT send immediate response        │    │
│  │ - Sets flusher for streaming              │    │
│  └──────────┬───────────────────────────────┘    │
└─────────────┼──────────────────────────────────────┘
              │
              │ 2. Connection held open
              │    Flusher configured
              │
              ▼
       ┌──────────────┐
       │  Goroutine   │
       │  Executes    │
       └──────┬───────┘
              │
              │ 3. boundMethod.Call(ctx, args)
              │
              ▼
┌──────────────────────────────────────────┐
│  Bound Go Method Execution                │
└──────────┬───────────────────────────────┘
           │
           │ 4. Result or Error
           │
           ▼
┌──────────────────────────────────────────┐
│  Write to HTTP Response                   │
│  - Flush SSE event: data: {...}          │
│  - Close connection                       │
└──────────┬───────────────────────────────┘
           │
           │ 5. HTTP Response received
           │
           ▼
       ┌──────────────┐
       │   Frontend   │
       │   receives   │
       │   result     │
       └──────────────┘
```

#### Implementation Details

**Frontend** (`calls.ts`):
```typescript
export function Call(options: CallOptions): CancellablePromise<any> {
    const id = generateID();

    const result = CancellablePromise.withResolvers<any>();

    // Create SSE connection
    const eventSource = new EventSource(
        `/wails/runtime?object=0&method=0&call-id=${id}&...`
    );

    eventSource.onmessage = (event) => {
        const response = JSON.parse(event.data);

        if (response.error) {
            result.reject(createError(response.error));
        } else {
            result.resolve(response.result);
        }

        eventSource.close();
    };

    eventSource.onerror = (err) => {
        result.reject(new Error("Connection failed"));
        eventSource.close();
    };

    result.oncancelled = () => {
        eventSource.close();
        // Send cancel request via separate HTTP call
        return fetch(`/wails/runtime?object=10&method=0`, {
            method: 'POST',
            body: JSON.stringify({ "call-id": id })
        });
    };

    return result.promise;
}
```

**Backend** (`messageprocessor_call.go`):
```go
func (m *MessageProcessor) processCallMethod(method int, rw http.ResponseWriter, r *http.Request, window Window, params QueryParams) {
	args, err := params.Args()
	if err != nil {
		m.httpError(rw, "Invalid binding call:", fmt.Errorf("unable to parse arguments: %w", err))
		return
	}

	callID := args.String("call-id")
	if callID == nil || *callID == "" {
		m.httpError(rw, "Invalid binding call:", errors.New("missing argument 'call-id'"))
		return
	}

	switch method {
	case CallBinding:
		var options CallOptions
		err := params.ToStruct(&options)
		if err != nil {
			m.httpError(rw, "Invalid binding call:", fmt.Errorf("error parsing call options: %w", err))
			return
		}

		// Set up SSE headers
		rw.Header().Set("Content-Type", "text/event-stream")
		rw.Header().Set("Cache-Control", "no-cache")
		rw.Header().Set("Connection", "keep-alive")

		flusher, ok := rw.(http.Flusher)
		if !ok {
			m.httpError(rw, "Streaming not supported", errors.New("cannot stream"))
			return
		}

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		// Register for cancellation
		func() {
			m.l.Lock()
			defer m.l.Unlock()
			m.runningCalls[*callID] = cancel
		}()

		defer func() {
			m.l.Lock()
			defer m.l.Unlock()
			delete(m.runningCalls, *callID)
		}()

		// Execute synchronously (NOT in goroutine)
		var boundMethod *BoundMethod
		if options.MethodName != "" {
			boundMethod = globalApplication.bindings.Get(&options)
			if boundMethod == nil {
				m.writeSSEError(rw, flusher, &CallError{
					Kind:    ReferenceError,
					Message: fmt.Sprintf("unknown bound method name '%s'", options.MethodName),
				})
				return
			}
		} else {
			boundMethod = globalApplication.bindings.GetByID(options.MethodID)
			if boundMethod == nil {
				m.writeSSEError(rw, flusher, &CallError{
					Kind:    ReferenceError,
					Message: fmt.Sprintf("unknown bound method id %d", options.MethodID),
				})
				return
			}
		}

		// Set window context
		if window != nil {
			ctx = context.WithValue(ctx, WindowKey, window)
		}

		// Execute method
		result, err := boundMethod.Call(ctx, options.Args)

		if err != nil {
			m.writeSSEError(rw, flusher, err)
			return
		}

		// Write result as SSE event
		jsonResult, err := json.Marshal(map[string]any{
			"result": result,
		})
		if err != nil {
			m.writeSSEError(rw, flusher, err)
			return
		}

		fmt.Fprintf(rw, "data: %s\n\n", jsonResult)
		flusher.Flush()

	default:
		m.httpError(rw, "Invalid binding call:", fmt.Errorf("unknown method: %d", method))
	}
}

func (m *MessageProcessor) writeSSEError(rw http.ResponseWriter, flusher http.Flusher, err error) {
	var cerr *CallError
	if !errors.As(err, &cerr) {
		cerr = &CallError{
			Kind:    RuntimeError,
			Message: err.Error(),
		}
	}

	jsonErr, _ := json.Marshal(map[string]any{
		"error": cerr,
	})

	fmt.Fprintf(rw, "data: %s\n\n", jsonErr)
	flusher.Flush()
}
```

### Option 2: Synchronous HTTP (Simpler)

Keep the connection open and return the result directly in the HTTP response body.

#### Flow Diagram

```
Frontend                MessageProcessor         Go Method
   │                           │                      │
   │ HTTP POST                 │                      │
   ├──────────────────────────►│                      │
   │                            │                      │
   │                            ├─────────────────────►│
   │                            │ Execute synchronously│
   │                            │                      │
   │                            │◄─────────────────────┤
   │                            │ Return result        │
   │◄───────────────────────────┤                      │
   │ HTTP 200 + JSON body       │                      │
   │                            │                      │
```

#### Implementation

**Frontend** (`calls.ts`):
```typescript
export function Call(options: CallOptions): CancellablePromise<any> {
    const id = generateID();

    const abortController = new AbortController();
    const result = CancellablePromise.withResolvers<any>();

    // Make HTTP request with abort signal
    const requestPromise = fetch(
        `/wails/runtime?object=0&method=0`,
        {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'x-wails-client-id': gn,
                'x-wails-window-name': windowName
            },
            body: JSON.stringify({
                "call-id": id,
                ...options
            }),
            signal: abortController.signal
        }
    );

    requestPromise
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            return response.json();
        })
        .then(data => {
            if (data.error) {
                result.reject(createError(data.error));
            } else {
                result.resolve(data.result);
            }
        })
        .catch(err => {
            if (err.name === 'AbortError') {
                result.reject(new CancelError("Call cancelled"));
            } else {
                result.reject(err);
            }
        });

    result.oncancelled = () => {
        abortController.abort();
        // Optionally send explicit cancel to clean up server-side
        return fetch(`/wails/runtime?object=10&method=0`, {
            method: 'POST',
            body: JSON.stringify({ "call-id": id })
        }).catch(() => {});
    };

    return result.promise;
}
```

**Backend** (`messageprocessor_call.go`):
```go
func (m *MessageProcessor) processCallMethod(method int, rw http.ResponseWriter, r *http.Request, window Window, params QueryParams) {
	args, err := params.Args()
	if err != nil {
		m.httpError(rw, "Invalid binding call:", fmt.Errorf("unable to parse arguments: %w", err))
		return
	}

	callID := args.String("call-id")
	if callID == nil || *callID == "" {
		m.httpError(rw, "Invalid binding call:", errors.New("missing argument 'call-id'"))
		return
	}

	switch method {
	case CallBinding:
		var options CallOptions
		err := params.ToStruct(&options)
		if err != nil {
			m.httpError(rw, "Invalid binding call:", fmt.Errorf("error parsing call options: %w", err))
			return
		}

		// Create cancellable context from request context
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		// Register for explicit cancellation
		func() {
			m.l.Lock()
			defer m.l.Unlock()
			m.runningCalls[*callID] = cancel
		}()

		defer func() {
			m.l.Lock()
			defer m.l.Unlock()
			delete(m.runningCalls, *callID)
		}()

		// Find bound method
		var boundMethod *BoundMethod
		if options.MethodName != "" {
			boundMethod = globalApplication.bindings.Get(&options)
			if boundMethod == nil {
				m.jsonError(rw, &CallError{
					Kind:    ReferenceError,
					Message: fmt.Sprintf("unknown bound method name '%s'", options.MethodName),
				})
				return
			}
		} else {
			boundMethod = globalApplication.bindings.GetByID(options.MethodID)
			if boundMethod == nil {
				m.jsonError(rw, &CallError{
					Kind:    ReferenceError,
					Message: fmt.Sprintf("unknown bound method id %d", options.MethodID),
				})
				return
			}
		}

		// Set window context
		if window != nil {
			ctx = context.WithValue(ctx, WindowKey, window)
		}

		// Execute method SYNCHRONOUSLY (blocking)
		result, err := boundMethod.Call(ctx, options.Args)

		if err != nil {
			m.jsonError(rw, err)
			return
		}

		// Write result to HTTP response
		m.json(rw, map[string]any{
			"result": result,
		})

	default:
		m.httpError(rw, "Invalid binding call:", fmt.Errorf("unknown method: %d", method))
	}
}

func (m *MessageProcessor) jsonError(rw http.ResponseWriter, err error) {
	var cerr *CallError
	if !errors.As(err, &cerr) {
		cerr = &CallError{
			Kind:    RuntimeError,
			Message: err.Error(),
		}
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK) // Still 200 to distinguish from HTTP errors

	json.NewEncoder(rw).Encode(map[string]any{
		"error": cerr,
	})
}
```

---

## Migration Implementation Guide

### Phase 1: Preparation

1. **Add Feature Flag**
   ```go
   // In application options
   type ApplicationOptions struct {
       // ... existing fields

       // UseHTTPOnlyBindings enables full HTTP round-trip for bindings
       // instead of eval-based responses
       UseHTTPOnlyBindings bool
   }
   ```

2. **Update MessageProcessor** to support both modes
   ```go
   type MessageProcessor struct {
       logger            *slog.Logger
       runningCalls      map[string]context.CancelFunc
       l                 sync.Mutex
       useHTTPOnly       bool // ← New field
   }
   ```

### Phase 2: Backend Implementation

1. **Modify `processCallMethod`** to conditionally use HTTP-only mode:
   ```go
   func (m *MessageProcessor) processCallMethod(method int, rw http.ResponseWriter, r *http.Request, window Window, params QueryParams) {
       // ... existing validation code ...

       if m.useHTTPOnly {
           m.processCallMethodHTTP(method, rw, r, window, params, callID, &options)
       } else {
           m.processCallMethodEval(method, rw, r, window, params, callID, &options)
       }
   }
   ```

2. **Implement `processCallMethodHTTP`** (synchronous version shown above)

3. **Keep existing `processCallMethodEval`** for backward compatibility

### Phase 3: Frontend Implementation

1. **Add runtime detection**:
   ```typescript
   // Detect if HTTP-only mode is enabled
   let useHTTPOnly = false;

   async function detectHTTPOnlyMode() {
       try {
           const response = await fetch('/wails/capabilities');
           const caps = await response.json();
           useHTTPOnly = caps.httpOnlyBindings || false;
       } catch (e) {
           useHTTPOnly = false;
       }
   }
   ```

2. **Implement dual-mode Call function**:
   ```typescript
   export function Call(options: CallOptions): CancellablePromise<any> {
       if (useHTTPOnly) {
           return CallHTTP(options);
       } else {
           return CallEval(options);
       }
   }
   ```

3. **Implement `CallHTTP`** (synchronous version shown above)

4. **Keep existing eval-based implementation as `CallEval`**

### Phase 4: Testing

1. **Unit Tests**: Test both modes independently
2. **Integration Tests**: Test cancellation in both modes
3. **Performance Tests**: Compare latency and throughput
4. **Stress Tests**: Many concurrent calls

### Phase 5: Migration Path

1. **v3-beta.1**: Ship with feature flag disabled by default
2. **v3-beta.2**: Enable by default, allow opt-out
3. **v3-rc.1**: Remove eval-based code entirely

---

## Trade-offs and Considerations

### Advantages of HTTP-Only

1. **Security**: Eliminates use of JavaScript eval()
2. **Simplicity**: Standard HTTP request/response pattern
3. **Debugging**: Easier to inspect in browser dev tools
4. **Compatibility**: Works with standard HTTP middleware/proxies
5. **Timeout Handling**: Can use standard HTTP timeouts

### Disadvantages of HTTP-Only

1. **Connection Management**: Must keep HTTP connections open during Go execution
2. **Browser Connection Limits**: HTTP/1.1 limits to 6-8 concurrent connections per domain (Note: HTTP/2 removes this limit)
3. **Latency**: Slightly higher overhead for connection setup/teardown
4. **Memory**: Open connections consume more resources than callbacks

### Concurrency Clarification

**"Synchronous" refers to the HTTP lifecycle, NOT global execution:**

```
Current (eval):                HTTP-Only (sync):
─────────────                  ──────────────────

Request 1 → [HTTP close]       Request 1 → [HTTP open] ──┐
              ↓                                           │
           [goroutine 1]                            [goroutine 1]
                                                           │
Request 2 → [HTTP close]       Request 2 → [HTTP open] ──┼─┐
              ↓                                           │ │
           [goroutine 2]                            [goroutine 2]
                                                           │ │
Request 3 → [HTTP close]       Request 3 → [HTTP open] ──┼─┼─┐
              ↓                                           │ │ │
           [goroutine 3]                            [goroutine 3]
                                                           │ │ │
ALL RUN CONCURRENTLY           ALL RUN CONCURRENTLY       │ │ │
                                                           ↓ ↓ ↓
                                                        [HTTP responses]
```

**Key Points**:
- ✅ Multiple methods CAN run concurrently (Go's http.Server creates a goroutine per request)
- ✅ No serialization or queuing of method calls
- ⚠️ Browser limits concurrent HTTP/1.1 connections to 6-8 (rarely hit in practice)
- ✅ HTTP/2 removes connection limits via multiplexing

### Recommended Approach

**Use Option 2 (Synchronous HTTP)** because:

1. **Simplest Implementation**: Minimal changes to existing code
2. **Best Debugging**: Standard HTTP tools work perfectly
3. **Natural Timeout Handling**: HTTP timeouts work out of the box
4. **Browser Connection Limits**: Unlikely to be hit in practice (6 concurrent bindings calls is very high)
5. **Context Cancellation**: Works perfectly with HTTP request context cancellation

### Browser Connection Limit Mitigation

If connection limits become an issue:

1. **Use HTTP/2**: Multiplexing removes per-connection limits
2. **Connection Pooling**: Reuse connections for multiple calls
3. **Request Queuing**: Frontend queues calls when at limit
4. **Hybrid Approach**: Long-running calls use HTTP, short calls use current eval method

---

## Proof of Concept Code

### Files to Modify

1. **`v3/pkg/application/messageprocessor_call.go`** - Main call processing
2. **`v3/internal/runtime/desktop/@wailsio/runtime/src/calls.ts`** - Frontend call implementation
3. **`v3/pkg/application/application.go`** - Add feature flag
4. **`v3/pkg/application/messageprocessor.go`** - Constructor changes

### Estimated Work

- **Backend**: 2-3 days
  - Refactor processCallMethod: 4 hours
  - Implement HTTP-only version: 4 hours
  - Testing and debugging: 8-12 hours

- **Frontend**: 1-2 days
  - Refactor Call function: 2 hours
  - Implement HTTP-only version: 3 hours
  - Testing and debugging: 4-8 hours

- **Documentation**: 1 day
  - Update API docs
  - Migration guide
  - Examples

**Total Estimated Time**: 4-6 days for full implementation

---

## Appendix: Platform-Specific JavaScript Execution APIs

### Windows (WebView2)

```cpp
// File: webview_window_windows.go (via CGo)
HRESULT ICoreWebView2::ExecuteScript(
    LPCWSTR javaScript,
    ICoreWebView2ExecuteScriptCompletedHandler* handler
);
```

### Linux (WebKitGTK)

```c
// File: linux_purego.go (via purego)
void webkit_web_view_evaluate_javascript(
    WebKitWebView* web_view,
    const gchar* script,
    gssize length,
    const gchar* world_name,
    const gchar* source_uri,
    GCancellable* cancellable,
    GAsyncReadyCallback callback,
    gpointer user_data
);
```

### macOS (WKWebView)

```objc
// File: webview_window_darwin.go (via Objective-C)
- (void)evaluateJavaScript:(NSString *)javaScriptString
         completionHandler:(void (^)(id, NSError *error))completionHandler;
```

---

## Conclusion

Migrating from eval-based responses to HTTP-only bindings is **technically feasible** and **architecturally sound**. The recommended synchronous HTTP approach provides:

- ✅ Elimination of JavaScript eval()
- ✅ Standard HTTP debugging
- ✅ Natural timeout and cancellation handling
- ✅ Minimal architectural changes
- ✅ Backward compatibility during migration

The main consideration is **managing open HTTP connections**, but this is unlikely to be a practical limitation for typical Wails applications.

**Recommended Next Steps**:
1. Create feature branch
2. Implement Option 2 (Synchronous HTTP) with feature flag
3. Test with existing Wails examples
4. Gather community feedback
5. Refine based on real-world usage
6. Make default in beta release

---

## Generated By

**Claude Code** (Anthropic)
**Date**: 2025-10-01
**Branch**: `http-only-bindings` (based on `v3-alpha`)
**Repository**: github.com/wailsapp/wails
