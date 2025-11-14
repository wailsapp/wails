# macOS Issues for Wails v3

This document catalogs all known macOS-related issues, bugs, and limitations for Wails v3, compiled from GitHub issues and codebase analysis.

## Open GitHub Issues

### Critical/High Priority

#### #4650 - [v3][macOS] Window behavior issues (works fine on Windows)
- **Reporter:** Yuelioi
- **Date:** October 17, 2025
- **Labels:** Bug, awaiting feedback
- **Status:** Open
- **Description:** Window operations function correctly on Windows but encounter problems on macOS
- **Link:** https://github.com/wailsapp/wails/issues/4650

#### #4649 - When fullscreen on MacOS, press Esc causing app out of fullscreen
- **Reporter:** Not specified
- **Date:** October 17, 2025
- **Status:** Open
- **Description:** Pressing Escape exits fullscreen mode on macOS unexpectedly
- **Link:** https://github.com/wailsapp/wails/issues/4649

#### #4353 - [V3] macOS window close with pending async Go-bound function call crashes
- **Reporter:** joshhardy
- **Date:** June 13, 2025
- **Labels:** Bug, MacOS, v3
- **Status:** Open
- **Description:** Application crashes when closing a macOS window while async Go-bound function calls are still pending
- **Link:** https://github.com/wailsapp/wails/issues/4353

#### #4389 - [v3]: hide/show window crashes the app
- **Reporter:** kunxl-gg
- **Date:** July 2, 2025
- **Labels:** Bug
- **Status:** Open
- **Description:** Application crashes when hiding/showing windows
- **Link:** https://github.com/wailsapp/wails/issues/4389

#### #4365 - nil pointer in processURLRequest
- **Reporter:** Etesam913
- **Date:** June 21, 2025
- **Labels:** Bug, MacOS, v3
- **Status:** Open
- **Description:** Nil pointer error in processURLRequest function
- **Link:** https://github.com/wailsapp/wails/issues/4365

#### #3743 - [v3] [mac] drag n drop does not work as expected after resizing window
- **Reporter:** Etesam913
- **Date:** September 12, 2024
- **Labels:** Bug
- **Status:** Reopened
- **Assignee:** leaanthony
- **Description:** Drag and drop functionality breaks after window resizing on macOS
- **Link:** https://github.com/wailsapp/wails/issues/3743

### Medium Priority

#### #4290 - [V3] MacOS - Print dialog does not open
- **Reporter:** stavros-k
- **Date:** May 19, 2025
- **Labels:** Bug
- **Status:** Open
- **Description:** Print dialog fails to open on macOS
- **Link:** https://github.com/wailsapp/wails/issues/4290

#### #4582 - [v3] How to enable internationalization for Mac apps
- **Reporter:** ToQuery
- **Date:** September 15, 2025
- **Labels:** Bug, MacOS, v3
- **Status:** Open
- **Description:** Issues with enabling internationalization (i18n) for macOS applications
- **Link:** https://github.com/wailsapp/wails/issues/4582

#### #4578 - [v3] bindings generator Incorrect namespace reference when generating TS code from Go type alias across packages
- **Reporter:** MikeLINGxZ
- **Date:** September 11, 2025
- **Labels:** Bug, MacOS, Typescript, v3
- **Status:** Open
- **Description:** TypeScript bindings generator creates incorrect namespace references for Go type aliases
- **Link:** https://github.com/wailsapp/wails/issues/4578

#### #4567 - [V3] How to configure the custom protocol for v3? I reported an error after configuring it
- **Reporter:** zk3151463
- **Date:** September 8, 2025
- **Labels:** Bug, MacOS, TODO, v3
- **Status:** Open
- **Description:** Custom protocol configuration errors on macOS
- **Link:** https://github.com/wailsapp/wails/issues/4567

#### #4669 - [v3] Shell metacharacters not allowed when trying to open a link
- **Reporter:** Etesam913
- **Date:** October 26, 2025
- **Labels:** Bug
- **Status:** Open
- **Description:** Shell metacharacters cause issues when opening links
- **Link:** https://github.com/wailsapp/wails/issues/4669

### Enhancement Requests

#### #4654 - [Feat] Request to add macOS (darwin) support for Windows-only theme and reload functions
- **Reporter:** Not specified
- **Date:** October 19, 2025
- **Labels:** Enhancement
- **Status:** Open
- **Description:** Request to extend Windows-specific theme and reload functionality to macOS
- **Link:** https://github.com/wailsapp/wails/issues/4654

#### #4329 - [V3] Cannot get the accent color of a macOS application
- **Reporter:** Etesam913
- **Date:** May 31, 2025
- **Labels:** Enhancement
- **Status:** Open
- **Description:** No API to retrieve macOS system accent color
- **Link:** https://github.com/wailsapp/wails/issues/4329

#### #3760 - [v3 feature] Add support for MacOS Panels (PR ready)
- **Reporter:** nixpare
- **Date:** September 20, 2024
- **Labels:** Enhancement, v3
- **Status:** Reopened
- **Description:** Feature request to add support for macOS Panels with PR ready
- **Link:** https://github.com/wailsapp/wails/issues/3760

#### #2413 - Installer for mac using packages
- **Reporter:** atarup
- **Date:** February 23, 2023
- **Labels:** Enhancement, Ready For Testing, TODO, v3
- **Status:** Open
- **Description:** Need better installer support for macOS using packages
- **Link:** https://github.com/wailsapp/wails/issues/2413

#### #1789 - Support Start on login
- **Reporter:** leaanthony
- **Date:** August 24, 2022
- **Labels:** TODO, v3
- **Status:** Reopened
- **Description:** Add support for "Start on login" functionality on macOS
- **Link:** https://github.com/wailsapp/wails/issues/1789

#### #1178 - Support Self-Updating
- **Reporter:** leaanthony
- **Date:** February 23, 2022
- **Labels:** Enhancement, TODO, v3
- **Status:** Open
- **Description:** Add self-updating capability for macOS applications
- **Link:** https://github.com/wailsapp/wails/issues/1178

### Development/Build Issues

#### #4652 - [v3] Can't connect to dev server when using pnpm
- **Reporter:** vegidio
- **Date:** October 17, 2025
- **Labels:** Bug
- **Status:** Open
- **Description:** Development server connection issues when using pnpm
- **Link:** https://github.com/wailsapp/wails/issues/4652

#### #4679 - [v3] window is not defined (Next.js output export)
- **Reporter:** shadiramadan
- **Date:** October 30, 2025
- **Labels:** Bug
- **Status:** Open
- **Description:** "window is not defined" error when using Next.js with export output
- **Link:** https://github.com/wailsapp/wails/issues/4679

#### #4112 - Cannot cross compile for windows on Mac after updating from v2.9.2 to v2.10.1
- **Reporter:** GregorExner
- **Date:** March 4, 2025
- **Labels:** Bug, TODO
- **Status:** Open
- **Description:** Cross-compilation from macOS to Windows broken in recent versions
- **Link:** https://github.com/wailsapp/wails/issues/4112

## Code-Level Issues (from v2 codebase that may affect v3)

### TODOs in Darwin-Specific Code

#### Window Start State
- **File:** `v2/internal/frontend/desktop/darwin/Application.m:44`
- **Issue:** `//TODO: Can you start a mac app minimised?`
- **Description:** Uncertainty about whether macOS applications can start in a minimized state
- **Impact:** Window start state functionality may not be fully implemented

#### Screen Mirroring
- **File:** `v2/internal/frontend/desktop/darwin/screen.go:46`
- **Issue:** `// TODO properly handle screen mirroring`
- **Reference:** https://developer.apple.com/documentation/appkit/nsscreen/1388393-screens?language=objc
- **Description:** Screen mirroring scenarios not properly handled
- **Impact:** Multi-monitor setups with screen mirroring may have issues

### Platform Warnings

#### Private API Usage in Debug Builds
- **File:** `v2/pkg/commands/build/build.go:374-375`
- **Warning:** "This darwin build contains the use of private APIs. This will not pass Apple's AppStore approval process. Please use it only as a test build for testing and debug purposes."
- **Description:** Debug/development builds use private Apple APIs
- **Impact:** Cannot submit debug builds to App Store; production builds must ensure no private API usage

## Summary Statistics

- **Total Open Issues:** 19 macOS-specific v3 issues
- **Critical/High Priority Bugs:** 6
- **Medium Priority Bugs:** 5
- **Enhancement Requests:** 6
- **Development/Build Issues:** 3
- **Code-level TODOs:** 2

## Priority Areas for Investigation

1. **Window Management:** Multiple issues related to window behavior, fullscreen, hide/show, and resize
2. **Async Operations:** Crashes when closing windows with pending async Go-bound function calls
3. **Drag and Drop:** Broken after window resizing
4. **Custom Protocols:** Configuration issues
5. **Development Experience:** pnpm support, Next.js compatibility, cross-compilation
6. **Feature Parity:** Theme/reload functions, accent color API, macOS Panels support

## Recommended Next Steps

1. Prioritize fixing the critical crash issues (#4353, #4389)
2. Investigate window behavior problems (#4650, #4649, #3743)
3. Resolve the screen mirroring TODO before v3 release
4. Address the minimized window start state uncertainty
5. Improve development experience (pnpm, Next.js support)
6. Consider implementing the enhancement requests for better macOS integration

---

*Report generated on: 2025-11-14*
*Branch: claude/find-mac-issues-v3-01LnsF7Giq1vywgkGpKXxZn9*
