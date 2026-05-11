# WebView2 Bindings Generator Architecture

This document describes the architecture and design decisions behind the WebView2 bindings generator.

## Overview

The WebView2 bindings generator converts Microsoft's IDL (Interface Definition Language) files into Go code that can interact with the WebView2 COM interfaces. The generator is designed to be:

- **Maintainable**: Easy to update when new WebView2 versions are released
- **Correct**: Properly handles COM semantics and memory management
- **Efficient**: Generates idiomatic Go code with minimal overhead

## Design Decisions

### Programmatic Emitter vs Templates

**Decision**: We use a programmatic emitter (Go code generating Go code) rather than pure templates.

**Rationale**:

1. **Type Safety**: Go code can leverage the compiler to catch errors in the generator itself
2. **Flexibility**: Complex transformations (e.g., type conversions, COM reference counting) are easier in code
3. **Maintainability**: IDE support and testing are better for Go code than template code
4. **Extensibility**: Adding special cases or custom logic doesn't require mixing code and templates

The generator still uses templates for simple, repetitive code (e.g., enum declarations), but the main emission logic is in Go.

### File-Per-Interface Organization

**Decision**: Each COM interface gets its own `.go` file.

**Rationale**:

1. **Build Performance**: Parallel compilation of many small files is faster than one huge file
2. **Code Navigation**: Easier to find specific interfaces
3. **Incremental Changes**: Changes to one interface don't require recompiling others
4. **Git Blame**: Easier to track changes to specific interfaces

Example: `ICoreWebView2_2.go` contains the `ICoreWebView2_2` interface and its vtable.

## Parser Grammar Overview

The IDL parser is built using the `participle` library and handles the following grammar elements:

### Lexer Rules

```
Comment    : (?:#|//)[^\n]*\n?
String     : "(\\"|[^"])*"
UUID       : [0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}
Hex        : 0x[a-fA-F0-9]+
Int        : [0-9]+
Ident      : [a-zA-Z]\w*
Number     : (?:@Int\.)?@Int
```

### Grammar Structure

```go
type IDL struct {
    Imports   []*Import
    Libraries []*Library
}

type Library struct {
    Header       *LibraryHeader
    Name         string
    Declarations []*Declaration
}

type Declaration struct {
    InterfaceForewardDecl string
    Enum                  *EnumDeclaration
    Struct                *StructDeclaration
    Interface             *InterfaceDeclaration
    CppQuote              string
}
```

## Type System

The generator maps IDL types to Go types as follows:

### Basic Types

| IDL Type | Go Type | Notes |
|-----------|----------|-------|
| `int` | `int32` | 32-bit integer |
| `long` | `int32` | 32-bit integer |
| `unsigned long` | `uint32` | 32-bit unsigned |
| `HRESULT` | `uint32` | COM error code |
| `BSTR` | `*uint16` | String pointer, converted via helpers |
| `VARIANT` | `uintptr` | Simplified (note: this is a known limitation) |
| `GUID` | `*GUID` | Custom struct for UUID |
| `interface *` | `*InterfaceStruct` | COM interface pointer |

### Interface Types

Each COM interface is represented as:

```go
type ICoreWebView2_2 struct {
    IUnknown
}

type ICoreWebView2_2Vtbl struct {
    ICoreWebView2_1Vtbl
    SomeMethod ComProc
}
```

This inheritance model allows interface versioning where `ICoreWebView2_2` extends `ICoreWebView2_1`.

### Helper Functions

The generator creates helper functions for common operations:

- **String Conversion**: `UTF16PtrFromString` / `UTF16PtrToString`
- **Memory Management**: `CoTaskMemFree` for freeing COM-allocated memory
- **GUID Handling**: `NewGUID` / `IsEqualGUID` for GUID operations
- **COM Callbacks**: `NewComProc` for creating Go functions that can be called from COM

## Capability System Design

The capability system is implicit in the generated code:

### Interface Versioning

WebView2 uses interface versioning to add new features:

```
ICoreWebView2 (base)
  ↓ extends
ICoreWebView2_2 (adds methods)
  ↓ extends
ICoreWebView2_3 (adds more methods)
```

### Query Interface Pattern

To use a newer interface:

```go
// Start with base interface
var webView2 *ICoreWebView2
// ...

// Try to query for v2
var webView2V2 *ICoreWebView2_2
hr := webView2.Vtbl.QueryInterface(..., &IID_ICoreWebView2_2, unsafe.Pointer(&webView2V2))
if hr == S_OK {
    // Use v2 methods
    webView2V2.SomeNewMethod()
} else {
    // Fallback for older runtimes
}
```

The generator provides all interfaces, but applications must:

1. Check the runtime version
2. Query for capability-specific interfaces
3. Handle failures gracefully

## Code Generation Process

### Step 1: Parse IDL

```go
idl, err := Parser.ParseBytes("", idlData)
```

The parser creates an AST representing the entire IDL file.

### Step 2: Process

```go
err = idl.Process()
```

During processing:

- Forward interface declarations are collected
- Enums are analyzed for type information
- Interfaces are linked to their parent versions
- Methods are analyzed for parameter types and return values

### Step 3: Generate

```go
generatedFiles, err := idl.Generate()
```

The emitter creates Go code:

1. **com.go**: Core COM types and helpers (from template)
2. **Enum files**: One per enum type
3. **Struct files**: One per struct type
4. **Interface files**: One per interface version

### Step 4: Write

```go
for _, file := range generatedFiles {
    fileName := "../pkg/webview2/" + file.FileName
    os.WriteFile(fileName, file.Content.Bytes(), 0644)
}
```

## Memory Management

### COM Reference Counting

All COM interfaces inherit from `IUnknown`, which has:

```go
AddRef() uintptr
Release() uintptr
```

**Rules**:

- When you receive an interface pointer, call `AddRef` if you're storing it
- Call `Release` when you're done with an interface
- The generator provides wrapper methods for common patterns

### String Memory

COM uses `BSTR` (Basic String) for strings:

```go
// COM allocates: caller must free
bstr := someComMethod()
defer CoTaskMemFree(unsafe.Pointer(bstr))
str := UTF16PtrToString(bstr)
```

The generator creates helper functions to handle this safely.

## Testing

### Unit Tests

The generator has unit tests in `scripts/generator/` covering:

- Enum parsing
- Interface method signatures
- Struct field types
- Type mapping correctness

### Integration Tests

The `webview2gen verify` command:

1. Regenerates all bindings from IDL
2. Compares with committed files
3. Fails if any differences found

This ensures the generator works correctly and the committed files are up-to-date.

## Future Improvements

### Known Limitations

1. **VARIANT Handling**: Currently simplified as `uintptr`. Full VARIANT support is complex and may be added later.
2. **Error Wrapping**: COM error codes could be wrapped in Go errors with more context.
3. **Async Patterns**: COM async patterns could use Go channels and goroutines more idiomatically.

### Planned Enhancements

1. **Generics**: Use Go generics for collection types (e.g., `IStringCollection` → `[]string`)
2. **Context Support**: Add `context.Context` for cancellation of long-running COM operations
3. **Performance**: Optimize hot paths with generated code (e.g., bulk operations)
