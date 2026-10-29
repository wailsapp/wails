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

## Fixed
<!-- Bug fixes -->
- Fix the application-level `MarshalError` option never being applied: `Bindings.Add` fell back from `ServiceOptions.MarshalError` straight to the default marshaler, so an error returned by a service method reached the frontend as `json.Marshal` of the error value (`{}` for an error type without exported fields) unless every service set its own `MarshalError`; the fallback chain is now service → application → default, as documented on both options by @APshenkin

## Deprecated
<!-- Soon-to-be removed features -->

## Removed
<!-- Features removed in this release -->

## Security
<!-- Security-related changes -->

---
