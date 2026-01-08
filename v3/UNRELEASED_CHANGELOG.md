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

## Changed
<!-- Changes in existing functionality -->
- **BREAKING:** Rename `EnableDragAndDrop` to `EnableFileDrop` in window options
- **BREAKING:** Rename `DropZoneDetails` to `DropTargetDetails` in event context
- **BREAKING:** Rename `DropZoneDetails()` method to `DropTargetDetails()` on `WindowEventContext`
- **BREAKING:** Remove `WindowDropZoneFilesDropped` event, use `WindowFilesDropped` instead
- **BREAKING:** Change HTML attribute from `data-wails-dropzone` to `data-file-drop-target`
- **BREAKING:** Change CSS hover class from `wails-dropzone-hover` to `file-drop-target-active`
- **BREAKING:** Remove `DragEffect`, `OnEnterEffect`, `OnOverEffect` options from Windows (were part of removed IDropTarget)

## Fixed
<!-- Bug fixes -->
- Fix outdated Manager API references in documentation (31 files updated to use new pattern like `app.Window.New()`, `app.Event.Emit()`, etc.)
- Fix file drag-and-drop on Windows not working at non-100% display scaling
- Fix HTML5 internal drag-and-drop being broken when file drop was enabled on Windows
- Fix file drop coordinates being in wrong pixel space on Windows (physical vs CSS pixels)
- Fix file drag-and-drop on Linux not working reliably with hover effects
- Fix HTML5 internal drag-and-drop being broken when file drop was enabled on Linux

## Deprecated
<!-- Soon-to-be removed features -->

## Removed
<!-- Features removed in this release -->
- Remove native `IDropTarget` implementation on Windows in favor of JavaScript-based approach (matches v2 behavior)

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
