# Unreleased Changes

<!-- 
This file is used to collect changelog entries for the next v3-alpha release.
Add your changes under the appropriate sections below.

Guidelines:
- Follow the "Keep a Changelog" format (https://keepachangelog.com/)
- Write clear, concise descriptions of changes
- Include the impact on users when relevant
- Use present tense ("Add feature" not "Added feature")
- Reference issue/PR numbers when applicable

This file is automatically processed by the nightly release workflow.
After processing, the content will be moved to the main changelog and this file will be reset.
-->

## Added
<!-- New features, capabilities, or enhancements -->
- Add `MessageDialog.WithButton()` method for builder pattern chaining when adding buttons without configuration (#4792)
- Add `MessageDialog.WithDefaultButton()` method for adding a button marked as default (Enter key) with builder pattern chaining (#4810)
- Add `MessageDialog.WithCancelButton()` method for adding a button marked as cancel (Escape key) with builder pattern chaining (#4810)

## Changed
<!-- Changes in existing functionality -->
- **BREAKING**: `MessageDialog.Show()` now returns `(string, error)` - the clicked button's label and an error if the dialog could not be displayed. This enables synchronous dialog workflows and proper error handling (#4792) by @leaanthony
- Switch to goccy/go-json for all runtime JSON processing (method bindings, events, webview requests, notifications, kvstore), improving performance by 21-63% and reducing memory allocations by 40-60%
- Optimize BoundMethod struct layout and cache isVariadic flag to reduce per-call overhead
- Use stack-allocated argument buffer for methods with <=8 arguments to avoid heap allocations
- Optimize result collection in method calls to avoid slice allocation for single return values
- Use sync.Map for MIME type cache to improve concurrent performance
- Use buffer pool for HTTP transport request body reading
- Lazily allocate CloseNotify channel in content type sniffer to reduce per-request allocations
- Remove debug CSS logging from asset server
- Expand MIME type extension map to cover 50+ common web formats (fonts, audio, video, etc.)

## Fixed
<!-- Bug fixes -->
- Fix `IsCancel` button not responding to Escape key on Linux (GTK) (#4810)

## Deprecated
<!-- Soon-to-be removed features -->

## Removed
<!-- Features removed in this release -->
- Remove github.com/wailsapp/mimetype dependency in favor of expanded extension map + stdlib http.DetectContentType, reducing binary size by ~1.2MB
- Remove gopkg.in/ini.v1 dependency by implementing minimal .desktop file parser for Linux file explorer, saving ~45KB
- Remove samber/lo from runtime code by using Go 1.21+ stdlib slices package and minimal internal helpers, saving ~310KB

## Security
<!-- Security-related changes -->

---

### Example Entries:

**Added:**
- Add support for custom window icons in application options
- Add new `SetWindowIcon()` method to runtime API (#1234)

**Changed:**
- Update minimum Go version requirement to 1.21
- Improve error messages for invalid configuration files

**Fixed:**
- Fix memory leak in event system during window close operations (#5678)
- Fix crash when using context menus on Linux with Wayland

**Security:**
- Update dependencies to address CVE-2024-12345 in third-party library
