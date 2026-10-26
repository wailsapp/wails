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
- Link the 404 page to matching Wails v2 documentation pages on v2.wails.io in [PR](https://github.com/wailsapp/wails/pull/6366) by @taliesin-ai
- Move the v2 documentation site to v2.wails.io in [PR](https://github.com/wailsapp/wails/pull/6365) by @taliesin-ai
- Serve v3 documentation from wails.io and v2 documentation from v2.wails.io in [PR](https://github.com/wailsapp/wails/pull/6364) by @taliesin-ai
- Update v3.0.0 release notes and correct the security advisory claim in [PR](https://github.com/wailsapp/wails/pull/6363) by @taliesin-ai

## Fixed
<!-- Bug fixes -->
- Serve `/wails/*` assets on Android with their real MIME type instead of a hardcoded `application/json`, so ES modules such as `/wails/runtime.js` are no longer rejected by the WebView and the runtime initialises correctly (#6014)

## Deprecated
<!-- Soon-to-be removed features -->

## Removed
<!-- Features removed in this release -->

## Security
<!-- Security-related changes -->

---
