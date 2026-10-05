# Unreleased Changes

<!-- 
This file is used to collect changelog entries for the next v3 release.
Add your changes under the appropriate sections below.

Guidelines:
- Follow the "Keep a Changelog" format (https://keepachangelog.com/)
- Write clear, concise descriptions of changes
- Include the impact on users when relevant
- Use present tense ("Add feature" not "Added feature").
- Reference issue/PR numbers when applicable.

This file is automatically processed by the nightly release workflow.
After processing, the content will be moved to the main changelog and this file will be reset.
-->

## Added
<!-- New features, capabilities, or enhancements -->

## Changed
<!-- Changes in existing functionality -->

## Fixed
<!-- Bug fixes -->
- Fix right-clicking a Linux system tray icon running the click handler as well as opening the menu, so an attached window toggled on every right-click (#6018)
- Keep Windows WebView2 recovery configuration failures non-fatal, release failed controllers, and retry initialization within the recovery budget. Fixes #6167.
- Keep WebView2 callback handlers alive until native release on Windows in [PR](https://github.com/wailsapp/wails/pull/6184) by @taliesin-ai

## Deprecated
<!-- Soon-to-be removed features -->

## Removed
<!-- Features removed in the release -->

## Security
<!-- Security-related changes -->
