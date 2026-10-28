# Unreleased Changes

<!--
AUTO-GENERATED FILE: do not edit in ordinary contribution PRs.

The Auto Changelog (v3) workflow generates entries after PRs merge.
Describe the user-visible change in the PR title and description instead.
Release tooling moves staged entries into docs/mpress/content/changelog.md
and resets this file. Do not delete or reset it by hand.

For changes to this automation or explicitly requested release-note corrections,
follow the changelog ownership policy in AGENTS.md and preserve existing entries.
-->

## Added
<!-- New features, capabilities, or enhancements -->

## Changed
<!-- Changes in existing functionality -->
- Improve asset server performance by formatting content lengths without fmt in [PR](https://github.com/wailsapp/wails/pull/6393) by @taliesin-ai
- Reduce memory copying when reading desktop file explorer metadata in [PR](https://github.com/wailsapp/wails/pull/6391) by @taliesin-ai

## Fixed
<!-- Bug fixes -->
- Release removed ring queue references in the asset server to avoid memory leaks in [PR](https://github.com/wailsapp/wails/pull/6397) by @taliesin-ai
- Release consumed content-sniffer buffers in the asset server to reduce memory use in [PR](https://github.com/wailsapp/wails/pull/6396) by @taliesin-ai
- Recognize mixed-case asset file extensions when serving assets in [PR](https://github.com/wailsapp/wails/pull/6395) by @taliesin-ai
- Correctly detect content type in asset server without retaining filenames in [PR](https://github.com/wailsapp/wails/pull/6394) by @taliesin-ai
- Fall back to a default when a desktop entry has no command in [PR](https://github.com/wailsapp/wails/pull/6392) by @taliesin-ai
- Stop the liquid-glass example from showing a dialog before the app starts on non-macOS in [PR](https://github.com/wailsapp/wails/pull/6389) by @taliesin-ai
- Keep window size constraints when fullscreen starts from a constrained size on Windows in [PR](https://github.com/wailsapp/wails/pull/6387) by @taliesin-ai
- Fix dialogs called before Run blocking forever instead of returning in [PR](https://github.com/wailsapp/wails/pull/6390) by @taliesin-ai

## Deprecated
<!-- Soon-to-be removed features -->

## Removed
<!-- Features removed in this release -->

## Security
<!-- Security-related changes -->

---
