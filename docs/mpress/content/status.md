---
title: "Project status"
description: "Wails v3 beta compatibility, security support, and upgrade guidance"
slug: "status"
sourcePath: "status.md"
---

## Current Status: Beta

Check out the [Changelog](/changelog/) for the latest status.

Our goal is a stable v3.0 release. Wails v2 remains the current stable release and continues to receive fixes. Test beta releases with your application before deployment.

## Beta Compatibility Promise

The v3 Beta contract covers desktop applications:

| Platform | Supported targets | Requirements and notes |
| --- | --- | --- |
| Windows | amd64 and arm64 | WebView2 runtime |
| macOS | Intel and Apple Silicon | The macOS and WebKit versions documented in the installation guide |
| Linux | amd64 and arm64 | GTK4 + WebKitGTK 6.0 by default; GTK3 + WebKit2GTK 4.1 remains a `-tags gtk3` legacy option through v3.0.x and is removed in v3.1 |

All targets require Go 1.25 or later for development. Android and iOS support is experimental and does not block the desktop Beta. Beta APIs aim for stability, but prerelease defects and explicitly announced changes may still be corrected before v3.0.0.

## How You Can Contribute

- Test the latest beta release and report reproducible bugs
- Contribute to documentation and examples
- Participate in discussions and provide feedback on draft WEPs
- Submit pull requests for bug fixes, documentation, or accepted WEPs

We welcome contributions from the community. If you'd like to help with these objectives, join the community discussions. Proposals for new functionality belong in a WEP PR, not a feature-request issue.

## Feedback and Updates

Report reproducible problems as issues; propose new functionality through a [WEP (Wails Enhancement Proposal)](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md) PR.

## Using the beta

Pin exact CLI, Go module and frontend runtime versions rather than tracking `latest`. For an existing alpha project, follow the [alpha-to-beta upgrade guide](/migration/alpha-to-beta/).

The [security policy](https://github.com/wailsapp/wails/blob/master/SECURITY.md) lists v3 beta releases as supported and alpha releases as unsupported. Report vulnerabilities through [private vulnerability reporting](https://github.com/wailsapp/wails/security/advisories/new), not public issues.

## Tracked work

- [Open bugs labelled v3](https://github.com/wailsapp/wails/issues?q=is%3Aissue+is%3Aopen+label%3ABug+label%3Av3)
- [Open v3 issues labelled P0 or P1](https://github.com/wailsapp/wails/issues?q=is%3Aissue+is%3Aopen+label%3Av3+label%3AP0%2CP1)
- [Release milestones](https://github.com/wailsapp/wails/milestones)

These live queries depend on issue labels; they are not a complete list or a promise about a release date or scope. Read the issues to assess their impact on your project.
