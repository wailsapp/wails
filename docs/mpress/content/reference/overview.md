---
title: "API Reference"
description: "Look up Wails types, methods, options and CLI commands"
slug: "reference/overview"
sourcePath: "reference/overview.md"
---

## About This Reference

Use this section to look up Wails types, methods, options and commands while working on an application. If you need a guided lesson, start with [Tutorials](/tutorials/overview/). For steps to complete a task, use [How-to Guides](/quick-start/next-steps/#features). For design background, use [Core Concepts](/concepts/architecture/).

**Organisation:**

- [Application](/reference/application/) - Core application APIs
- [Window](/reference/window/) - Window creation and management; [window options](/features/windows/options/) lists configuration fields
- [Menu](/reference/menu/) - Application, context, and system tray APIs; [menu item reference](/features/menus/reference/) covers item types and properties
- [Events](/reference/events/) - Event system and built-in events
- [Dialogs](/reference/dialogs/) - File and message dialogs
- [Frontend Runtime](/reference/frontend-runtime/) - Frontend runtime APIs
- [CLI command overview](/reference/cli/) - Commands grouped by purpose; [flags and examples](/guides/cli/) gives individual command details
- [Update manifest protocol](/reference/update-manifest/) - Update metadata and verification fields
- [Mobile API](/guides/mobile/mobile-api/) - Experimental mobile APIs

## API Conventions

@details{title="Go API Conventions - For developers new to Go"}
### Naming

- <strong>Types</strong>: PascalCase (e.g., `WebviewWindow`)
- <strong>Methods</strong>: PascalCase (e.g., `SetTitle()`)
- <strong>Options</strong>: PascalCase structs (e.g., `WindowOptions`)
- <strong>Constants</strong>: PascalCase (e.g., `WindowStartStateMaximised`)

#### Error Handling

Most methods that can fail return `error` as the last return value. `app.Run()` blocks until the application exits and returns any startup error:

```go
if err := app.Run(); err != nil {
    log.Fatal(err)
}
```

Window construction does not return an error — `app.Window.New()` returns `*WebviewWindow` directly.

#### Context

Service lifecycle methods receive a `context.Context`:

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // ctx is cancelled when the application is shutting down.
    return nil
}
```

The application's lifetime context is available via `app.Context()`. There is no `RunWithContext` — call `app.Run()`.

#### Options Pattern

Configuration uses option structs:

```go
app := application.New(application.Options{
    Name: "My App",
    Description: "A demo application",
    Services: []application.Service{
        application.NewService(&MyService{}),
    },
})
```

@end

### JavaScript API Conventions

#### Naming

- **Functions**: camelCase (e.g., `setTitle()`)
- **Constants**: SCREAMING<em>SNAKE</em>CASE (e.g., `WINDOW_EVENT_FOCUS`)

#### Async by Default

All Go method calls return Promises:

```javascript
// Async/await (recommended)
const result = await MyService.DoSomething()

// Promise chain
MyService.DoSomething()
    .then(result => console.log(result))
    .catch(error => console.error(error))
```

#### Error Handling

Go errors become JavaScript exceptions:

```javascript
try {
    await MyService.MightFail()
} catch (error) {
    console.error('Go error:', error)
}
```

#### Type Safety

TypeScript definitions are auto-generated:

```typescript
// Fully typed
import { Greet } from './bindings/GreetService'

const message: string = await Greet("World")
```

## Package Structure

```
github.com/wailsapp/wails/v3/pkg/
├── application/          # Core application package
│   ├── application.go    # App type
│   ├── webview_window.go # Window management
│   ├── menu.go           # Menu types
│   ├── event_manager.go  # Event system
│   └── dialogs.go        # Dialog APIs
├── events/               # Event constants
└── services/             # Built-in services
    ├── dock/             # macOS dock (includes badge support)
    ├── fileserver/       # File-server service
    ├── kvstore/          # Key/value store
    ├── log/              # Structured logging service
    ├── notifications/    # Notifications service
    └── sqlite/           # SQLite service
```

## Import Paths

### Go

```go
import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)
```

### JavaScript

```javascript
// Auto-generated bindings
import { MyMethod } from './bindings/MyService'

// Runtime APIs
import { Events, Window } from '@wailsio/runtime'
```

## Type Reference

### Common Types

@tabs{sync-key="lang"}
[Go]
```go
// Application
type App struct { /* ... */ }
type Options struct { /* ... */ }

// Window
type WebviewWindow struct { /* ... */ } // implements the Window interface
type WebviewWindowOptions struct { /* ... */ }

// Menu
type Menu struct { /* ... */ }
type MenuItem struct { /* ... */ }

// Events — there is no generic Event type; events are typed by source.
type ApplicationEvent struct { /* ... */ }
type WindowEvent struct { /* ... */ }
type CustomEvent struct { /* ... */ }
type EventListener struct { /* ... */ }

// Dialogs
type OpenFileDialogOptions struct { /* ... */ }
type SaveFileDialogOptions struct { /* ... */ }
```

[TypeScript]
```typescript
// Window runtime
interface WindowOptions {
    title?: string
    width?: number
    height?: number
    // ...
}

// Events
type EventCallback = (data: any) => void

// Bindings (auto-generated)
export function MyMethod(arg: string): Promise<string>
```

@end

## Platform Differences

Some APIs behave differently on different platforms:

| Feature | Windows | macOS | Linux |
| --- | --- | --- | --- |
| **Application Menu** | Window menu bar | Global menu bar | Window menu bar |
| **System Tray** | Notification area | Menu bar | System tray |
| **Dock** | N/A | ✅ Available | N/A |
| **File dialogs** | Native | Native | Native (GTK) |
| **Transparency** | ✅ Full | Requires [`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background) | ⚠️ Limited |

Platform-specific behaviour is documented in each API section.

## Versioning

Wails v3 follows semantic versioning:

- **Major** (v4.0.0): Breaking changes
- **Minor** (v3.1.0): New features, backwards-compatible
- **Patch** (v3.0.1): Bug fixes, backwards-compatible

**Wails v3.0.0:** Stable

## Deprecation Policy

When APIs are deprecated:

1. **Marked in docs** with deprecation notice
2. **Alternative provided** with migration guide
3. **Maintained for 1 major version** before removal
4. **Compiler warnings** (where possible)

## API Stability

### Stable APIs ✅

These APIs are stable and safe for production use:

- Core application APIs
- Window management
- Menu system
- Event system
- File dialogs
- Service bindings

### Experimental APIs ⚠️

Only APIs explicitly documented as experimental are outside the stable desktop compatibility contract:

- Android support
- iOS support
- Other APIs explicitly marked experimental

Experimental APIs are marked in documentation.

## Getting Help

### API Questions

1. **Check this reference** - Complete API documentation
2. **Check examples** - [GitHub examples](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples)
3. **Search Discord** - [Discord server](https://discord.gg/JDdSxwjhGf)
4. **Ask the community** - Discord #help channel

### Reporting API Issues

Found a bug or inconsistency?

1. **Check existing issues** - [GitHub issues](https://github.com/wailsapp/wails/issues)
2. **Create detailed report** - Include code, error, platform
3. **Provide reproduction** - Minimal example that demonstrates issue

## Related Documentation

- [Tutorials](/tutorials/overview/) - Learn by building real applications
- [Guides](/guides/architecture/) - Task-oriented guides for common scenarios
- [Features](/features/windows/basics/) - Feature-by-feature documentation
- [Examples](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples) - Working code examples on GitHub

---

**Browse the API:** Use the navigation on the left to explore specific APIs.
