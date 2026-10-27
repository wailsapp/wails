---
title: "Roadmap"
description: "Wails v3 project status, planned features, and how to contribute"
slug: "status"
sourcePath: "status.md"
---

## Wails v3.0.0 — Stable

Check out the [Changelog](/changelog/) for the latest status.

This documentation covers Wails v3.0.0. After v3.0.0, fixes and improvements ship in weekly releases, with dated nightly previews in between and hotfixes for urgent problems. See the [release policy](/contributing/release-policy/).

## Desktop Compatibility

The v3.0.0 compatibility contract covers desktop applications:

| Platform | Supported targets | Requirements and notes |
| --- | --- | --- |
| Windows | amd64 and arm64 | WebView2 runtime |
| macOS | Intel and Apple Silicon | The macOS and WebKit versions documented in the installation guide |
| Linux | amd64 and arm64 | GTK4 + WebKitGTK 6.0 by default; GTK3 + WebKit2GTK 4.1 remains a `-tags gtk3` legacy option through v3.0.x and is removed in v3.1 |

All targets require Go 1.26 or later for development. Android and iOS support remains experimental and is outside the stable desktop compatibility contract.

## How You Can Contribute

- Test Wails v3.0.0 and report reproducible bugs
- Contribute to documentation and examples
- Participate in discussions and provide feedback on draft WEPs
- Submit pull requests for bug fixes, documentation, or accepted WEPs

We welcome contributions from the community. If you'd like to help with these objectives, join the community discussions. Proposals for new functionality belong in a WEP PR, not a feature-request issue.

## Feedback and Updates

This roadmap is subject to change based on community feedback and project priorities. We'll update it regularly to reflect progress and changes in direction. Report reproducible problems as issues; propose new functionality through a [WEP (Wails Enhancement Proposal)](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md) PR.
