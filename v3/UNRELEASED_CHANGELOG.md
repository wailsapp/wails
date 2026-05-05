# Unreleased Changes

<!-- 
This file is used to collect changelog entries for the next v3 alpha release.
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
- Add `Sound` (default/silent/named), `Attachments`, `ThreadID`, `InterruptionLevel`, and `Schedule` fields to `NotificationOptions` in the notifications service. Each gracefully degrades on platforms with limited support; see the package godoc for the per-platform support matrix.
- Add `NotificationService.UpdateNotification(options)` to update an in-flight notification by ID. macOS auto-deduplicates by identifier; Linux uses the D-Bus replaces_id parameter; Windows currently redelivers as a new notification (true replace requires upstream wintoast tag/group support).

## Changed
<!-- Changes in existing functionality -->
- Internal: rewrite the Windows notifications backend to build toast XML directly via the `wintoast` subpackage, dropping the high-level `git.sr.ht/~jackmordaunt/go-toast/v2` builder. Behaviour is preserved; this unlocks first-class custom sounds, hero/inline images, threading, scheduling, scenarios, and tag-based update on Windows.
- Internal: macOS notifications C bridge now accepts a single JSON options blob instead of individual parameters, simplifying future field additions.
- Documentation: clarify that the Windows notifications backend supports reply fields (only Linux does not).

## Fixed
<!-- Bug fixes -->

## Deprecated
<!-- Soon-to-be removed features -->

## Removed
<!-- Features removed in this release -->

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
