# HTTP-Only Bindings Implementation Summary

## Quick Answer to Your Question

**Q: "Does synchronous execution mean only 1 method can be called at any time?"**

**A: No!** Multiple methods run concurrently. "Synchronous" refers to each individual HTTP request/response lifecycle, not global serialization.

### Visual Comparison

```
Current (eval-based):              HTTP-Only (synchronous):
════════════════════              ═══════════════════════════

Call 1 → HTTP closes immediately  Call 1 → HTTP stays open ──┐
         ↓                                                    │
      [Go method in goroutine]                        [Go method runs]
                                                              │
Call 2 → HTTP closes immediately  Call 2 → HTTP stays open ──┼──┐
         ↓                                                    │  │
      [Go method in goroutine]                        [Go method runs]
                                                              │  │
                                                              ↓  ↓
      ALL RUN CONCURRENTLY                            [Responses sent]
```

**Both approaches support unlimited concurrent calls.** The difference is:
- **Current**: HTTP closes, result via eval
- **HTTP-Only**: HTTP stays open, result via response body

---

## Complete Documentation

The full technical guide is in **`HTTP_BINDINGS_TECHNICAL_GUIDE.md`** (40KB, 1,268 lines).

### What's Documented

1. **Current Architecture (BEFORE)**
   - Complete flow diagrams
   - Code examples with line numbers
   - Platform-specific JavaScript execution APIs

2. **Cancellation Mechanism**
   - How `runningCalls` map works
   - Context cancellation flow
   - Cleanup procedures

3. **Proposed HTTP-Only Architecture (AFTER)**
   - Two options: SSE streaming vs synchronous HTTP
   - Complete code implementations
   - Flow diagrams

4. **Concurrency Analysis**
   - How Go's http.Server handles requests
   - Browser connection limits (HTTP/1.1 vs HTTP/2)
   - Why it's not a practical limitation

5. **Migration Guide**
   - 5-phase implementation plan
   - Feature flag approach
   - Backward compatibility

6. **Trade-offs**
   - Security benefits (no eval)
   - Debugging advantages
   - Performance considerations

---

## Key Findings

### Current Architecture Issues

1. **Uses JavaScript eval()** to deliver results
2. **Platform-dependent** JavaScript injection
3. **Harder to debug** (results not in network tab)
4. **Security concerns** (eval execution)

### HTTP-Only Benefits

1. ✅ **No eval()** - standard HTTP responses
2. ✅ **Better debugging** - visible in browser DevTools
3. ✅ **Simpler architecture** - standard request/response
4. ✅ **Standard timeouts** - HTTP timeout handling
5. ✅ **Middleware compatible** - works with HTTP proxies

### Only Trade-off

⚠️ **Browser connection limit** (HTTP/1.1 = 6-8 concurrent)
- Rarely an issue (most apps make 1-2 calls at once)
- HTTP/2 removes limit entirely (multiplexing)
- Calls complete quickly (milliseconds)

---

## Recommended Implementation

**Use Synchronous HTTP (Option 2)**

### Backend Changes

```go
// messageprocessor_call.go
func (m *MessageProcessor) processCallMethod(...) {
    // ... validation ...

    // Execute synchronously (NOT in goroutine)
    result, err := boundMethod.Call(ctx, options.Args)

    if err != nil {
        m.jsonError(rw, err)
        return
    }

    // Return result in HTTP response
    m.json(rw, map[string]any{
        "result": result,
    })
}
```

### Frontend Changes

```typescript
// calls.ts
export function Call(options: CallOptions): CancellablePromise<any> {
    const abortController = new AbortController();

    const response = await fetch('/wails/runtime', {
        method: 'POST',
        body: JSON.stringify(options),
        signal: abortController.signal
    });

    const data = await response.json();
    return data.error ? Promise.reject(data.error) : data.result;
}
```

### Cancellation

```typescript
result.oncancelled = () => {
    abortController.abort(); // Browser cancels HTTP request
    // HTTP context receives cancellation automatically
};
```

---

## Implementation Timeline

### Estimated Effort: 4-6 days

1. **Backend** (2-3 days)
   - Refactor `processCallMethod`: 4 hours
   - Implement HTTP-only version: 4 hours
   - Testing and debugging: 8-12 hours

2. **Frontend** (1-2 days)
   - Refactor `Call` function: 2 hours
   - Implement HTTP-only version: 3 hours
   - Testing and debugging: 4-8 hours

3. **Documentation** (1 day)
   - API docs updates
   - Migration guide
   - Examples

---

## Migration Strategy

### Phase 1: Add Feature Flag
```go
type ApplicationOptions struct {
    UseHTTPOnlyBindings bool // Default: false
}
```

### Phase 2: Dual Implementation
- Keep eval-based for compatibility
- Add HTTP-only behind flag
- Both modes tested

### Phase 3: Gradual Rollout
- **v3-beta.1**: Flag disabled by default
- **v3-beta.2**: Flag enabled by default
- **v3-rc.1**: Remove eval-based code

---

## Files Modified

### Backend
- `v3/pkg/application/messageprocessor_call.go` - Main changes
- `v3/pkg/application/messageprocessor.go` - Add flag
- `v3/pkg/application/application.go` - Options

### Frontend
- `v3/internal/runtime/desktop/@wailsio/runtime/src/calls.ts` - Main changes

### Total Lines Changed: ~200-300

---

## Testing Requirements

1. **Unit Tests**
   - Both modes independently
   - Cancellation in both modes
   - Error handling

2. **Integration Tests**
   - Multiple concurrent calls
   - Long-running methods
   - Network failures

3. **Performance Tests**
   - Latency comparison
   - Throughput under load
   - Memory usage

4. **Browser Tests**
   - Chrome, Firefox, Safari
   - Edge cases with connection limits

---

## Branch Information

**Branch**: `http-only-bindings`
**Base**: `v3-alpha`
**Status**: Documentation complete, implementation ready to start

### To Continue Implementation

```bash
git checkout http-only-bindings

# Start with backend
cd v3/pkg/application
# Modify messageprocessor_call.go

# Then frontend
cd v3/internal/runtime/desktop/@wailsio/runtime/src
# Modify calls.ts
```

---

## Questions & Answers

### Q: Will this break existing apps?
**A:** No - feature flag ensures backward compatibility during migration.

### Q: What about WebSockets?
**A:** Separate concern - WebSockets already not supported by AssetServer.

### Q: Performance impact?
**A:** Minimal - HTTP overhead is ~1-2ms, negligible compared to method execution time.

### Q: Why not keep eval?
**A:** Security best practices, better debugging, simpler architecture.

---

## Next Steps

1. Review this documentation
2. Decide on implementation approach
3. Create implementation tasks
4. Begin Phase 1 (feature flag)
5. Implement backend changes
6. Implement frontend changes
7. Test thoroughly
8. Merge to v3-alpha

---

**Generated by**: Claude Code
**Date**: 2025-10-01
**Total Documentation**: 1,400+ lines
**Repository**: github.com/wailsapp/wails
