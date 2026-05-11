# Generator Updates

This file tracks significant changes to the WebView2 bindings generator and the generated bindings.

## Current IDL Version

Last updated from: **WebView2 SDK 1.0.2903.40**

IDL file: `scripts/WebView2.1.0.2903.40.idl`

## Version Compatibility Matrix

| Wails Version | Min WebView2 Runtime | Notes |
|---------------|---------------------|-------|
| v3.0.0       | 1.0.1418.22        | Initial v3 release |

## Recent Changes

### [Date] - Version X.Y.Z

**IDL Version**: 1.0.XXXX.XX

**Binding Changes**:
- Added interfaces: [list new interfaces]
- Removed interfaces: [if any]
- Fixed bugs: [list bug fixes]
- Breaking changes: [if any]

**Generator Changes**:
- Improved [feature]
- Fixed [issue]
- Added [capability]

**Impact**:
- [Describe impact on Wails applications]
- [Any migration needed]

## Historical Changes

### Initial Release (v1.0.0)

**IDL Version**: 1.0.1823.32

**Initial Implementation**:
- First working version of the generator
- Basic COM interface support
- String and integer types
- Core WebView2 interfaces

**Known Issues**:
- Memory leaks in string handling
- Incorrect COM reference counting
- Missing error propagation

### v2.0.0 - Complete Rewrite

**IDL Version**: 1.0.2903.40

**Major Changes**:

1. **Memory Management Fixes**
   - Fixed critical memory leak in string handling
   - Proper `CoTaskMemFree` calls for all COM-allocated strings
   - Correct COM reference counting

2. **Error Handling**
   - All HRESULT return codes now checked
   - Proper error propagation to Go layer
   - Better error messages for debugging

3. **Interface Coverage**
   - Added all missing interfaces up to 1.0.2903.40
   - Complete method coverage for all interfaces
   - Correct interface inheritance chain

4. **Code Quality**
   - Fixed incorrect VARIANT handling
   - Added BSTR conversion helpers
   - Improved type safety

**Breaking Changes**:
- Signature changes for methods that previously ignored errors
- Different error handling patterns required
- Updated COM helper APIs

**Impact**:
- Applications no longer leak memory
- Much more stable WebView2 integration
- Full access to modern WebView2 features

## Updating IDL

To update to a new WebView2 SDK version:

1. Download the new IDL from NuGet:
   ```bash
   go run scripts/generator/cmd/webview2gen/main.go --version=1.0.XXXX.XX full
   ```

2. Verify the generator works:
   ```bash
   go run scripts/generator/cmd/webview2gen/main.go verify
   ```

3. Update this file:
   - Update the IDL version
   - Document new interfaces added
   - Note any breaking changes

4. Update the version compatibility matrix in README.md

## Testing Checklist

When updating the generator or bindings:

- [ ] Run `webview2gen verify` - ensure committed files match generated
- [ ] Run unit tests: `go test ./scripts/generator/...`
- [ ] Check formatting: `gofmt -l pkg/webview2`
- [ ] Lint: `golangci-lint run scripts/generator/...`
- [ ] Test on Windows with real WebView2 app
- [ ] Update this file with changes
- [ ] Update README.md if needed
- [ ] Tag new version: `git tag webview2/vX.Y.Z`
