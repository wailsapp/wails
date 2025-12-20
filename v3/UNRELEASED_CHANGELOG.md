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
- Add selfupdate service for application self-updating via GitHub/GitLab/Gitea releases (#1178)
  - Supports GitHub, GitLab, and Gitea release sources
  - ECDSA and PGP signature verification for secure updates
  - Themable progress UI with customizable colors and icons
  - Progress events for download tracking
  - Automatic platform detection (OS/arch)
- Add `wails3 tool selfupdate` CLI commands for update signing:
  - `keygen` - Generate ECDSA signing keys
  - `sign` - Sign application binaries
  - `verify` - Verify signed binaries

## Changed
<!-- Changes in existing functionality -->
- Update the documentation for Window `X/Y` options @ruhuang2001

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
