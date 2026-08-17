# Unreleased Changes

<!-- 
This file is used to collect changelog entries for the next v3 release.
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
- Add `Options.Linux.ApplicationID` to override the GTK application id, which sandboxed (flatpak) builds have to set to the id their manifest declares (#5972)

## Changed
<!-- Changes in existing functionality -->
- Default `Options.Linux.ProgramName` to `ApplicationID` when only the latter is set, so windows keep matching their `.desktop` file on Wayland, where GTK takes the surface `app_id` from the program name (#5972)

## Fixed
<!-- Bug fixes -->
- Report a Linux `ApplicationID` that GTK would reject and fall back to the derived id, instead of letting `gtk_application_new()` assert and take the process down during startup (#5972)
- Fix Linux application ids derived from an application `Name` that starts with a digit, such as `1example`, producing `org.wails.1example`, which GTK rejects (#5972)

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
