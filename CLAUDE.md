# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Wails is a Go-based desktop application framework that combines Go backends with web frontends. The project maintains two parallel versions:
- **v2/**: Stable production version (Go 1.21+)
- **v3/**: Alpha development version with next-generation features (Go 1.24+)

## Common Commands

### Root Level Tasks (using Taskfile.yaml)
```bash
task v2:lint                   # Lint v2 codebase (golangci-lint with 3m timeout)
task v2:download               # Update v2 dependencies (go mod tidy)
task v3:install                # Install wails3 CLI tool
task v3:test:examples          # Build all v3 example applications
task v3:precommit              # Run v3 tests and formatting
task format:md                 # Format all markdown files
```

### Direct Go Commands
```bash
# For v2 development
cd v2 && go mod tidy
cd v2 && golangci-lint run ./... --timeout=3m -v

# For v3 development  
cd v3 && go install ./cmd/wails3
cd v3/examples && ../bin/wails3 build
```

## Architecture Overview

### Multi-Version Structure
The codebase maintains two complete implementations:
- `v2/`: Production-ready framework with stable API
- `v3/`: Next-generation rewrite with enhanced capabilities

Each version contains:
- `cmd/`: CLI implementation
- `pkg/`: Public API packages
- `internal/`: Internal implementation details
- `examples/`: Demonstration applications

### Key Components

#### Application Framework
- **Application**: Core app lifecycle, window management, and native OS integration
- **Binding System**: Automatic Go method exposure to JavaScript with TypeScript generation
- **Asset Server**: Web asset serving with development/production modes
- **Menu System**: Cross-platform native menu integration

#### CLI Architecture
The CLI tools (`wails` for v2, `wails3` for v3) provide:
- Project initialization with frontend framework templates
- Development server with live reload
- Production build with single binary output
- Cross-platform packaging and distribution

#### Frontend Integration
- Frontend-agnostic (React, Vue, Svelte, Vanilla JS, etc.)
- Embedded web server for development
- Asset bundling for production
- Bi-directional event system between Go and JavaScript

### Development Patterns

#### Template-Driven Development
- Project templates in `v2/pkg/templates/` and `v3/internal/templates/`
- Generator system for creating boilerplate code
- Example-driven development with 30+ examples in v3

#### Cross-Platform Support
- Native webview integration (WebView2 on Windows, WebKit on macOS/Linux)
- OS-specific capabilities detection and feature flags
- Platform-specific build configurations and installers

#### Testing Strategy
- Example applications serve as integration tests
- CLI testing through example builds
- Manual testing across multiple platforms

## Important Notes

- The v3 branch is in active alpha development - expect breaking changes
- All v3 examples must build successfully (`task v3:test:examples`)
- Linting uses a permissive `.golangci.yml` configuration focused on formatting
- Contributors must follow the community guide at https://wails.io/community-guide
- The project supports 9+ human languages for documentation

## Current Refactor Project

**PRD Reference**: See `message.txt` for the complete Product Requirements Document outlining the App API restructuring project.

**Objective**: Refactor the monolithic `App` struct API into organized manager structs for better discoverability and maintainability. This involves grouping related functionality (Windows, ContextMenus, KeyBindings, Browser, Environment, Dialogs, Events, Menus) into dedicated manager structs while maintaining backward compatibility.