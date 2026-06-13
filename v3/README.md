# Wails v3 Beta

Welcome to Wails v3! This is the beta release of the next major version of Wails.

## What Beta means

The v3 API in `v3/pkg/application` and related packages is stable. We will not make breaking changes without a deprecation cycle and clear migration guidance.

## Getting Started

Full documentation is available at [v3.wails.io](https://v3.wails.io).

To get started with a new project:
```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
wails3 init -n myapp
```

## Beta scope — what's deferred

The following areas are still under active development and are **not** covered by the v3 stability guarantee:

- **iOS / Android**: Not production-ready in v3 beta; deferred to a future release.
- **GTK4 Linux backend**: Now the default on Linux but still in active stabilisation. Some rendering edge cases may remain.
- **Some niche AppKit APIs**: A small number of macOS-specific window APIs are still being finalised.

The items listed above may have breaking changes without a deprecation cycle.

## Giving feedback

- **Bug reports**: [github.com/wailsapp/wails/issues](https://github.com/wailsapp/wails/issues)
- **Discussions**: [github.com/wailsapp/wails/discussions](https://github.com/wailsapp/wails/discussions)
- **Discord**: [wails.io/discord](https://wails.io/discord)
