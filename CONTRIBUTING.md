# Contributing to Wails

Thanks for your interest in contributing! Wails has two active tracks. Pick the one that matches the change you want to make.

## v2 (stable)

- Source: the `v2/` directory on the `master` branch.
- Docs: [wails.io](https://wails.io/).
- Detailed contributor guide: [wails.io/community-guide](https://wails.io/community-guide).
- Changelog: update `website/src/pages/changelog.mdx` under the `[Unreleased]` heading.

## v3

- Source: the `v3/` directory on `master`. Target ordinary contributions at `master`; release branches are maintainer-managed.
- Release scope: prioritise bug fixes by severity and deliver them through nightly patch updates. v3.1.0 is reserved for new features.
- Docs: [v3.wails.io](https://v3.wails.io/).
- Detailed contributor guide: [v3.wails.io/contributing](https://v3.wails.io/contributing/).
- **Documentation contributions:** [download M-Press and follow the contributor walkthrough](docs/mpress/content/contributing/documentation.md). You can preview and edit the docs without building Wails.
- Changelog: v3 entries are generated automatically after merge from the PR title and CodeRabbit walkthrough when available. Describe the user-visible impact clearly in your PR title and description. Do not edit `v3/UNRELEASED_CHANGELOG.md` or `docs/mpress/content/changelog.md` in ordinary contributions; changelog automation maintenance and explicitly requested release-note corrections are separate maintainer work.
- **New functionality and changes to public behaviour require a [Wails Enhancement Proposal (WEP)](v3/wep/README.md)** before implementation. Open the WEP as a draft PR; do not open a feature-request issue. [Ideas](https://github.com/wailsapp/wails/discussions/categories/ideas) is available for optional early discussion.

## Pull request checklist

- One feature or fix per PR.
- For v2, update its changelog. For v3, describe the change in the PR and leave the generated changelog files to automation.
- Follow the existing coding style in the area you're touching.
- Add tests where it makes sense.

## Licence and provenance

Wails is released under the [MIT Licence](LICENSE), and contributions come in
under the same terms. By opening a pull request you confirm that:

- You wrote the contribution yourself, or you have the right to submit it.
- You are licensing it to the project under the MIT Licence.
- Any third-party code included in it is compatibly licensed, and its origin
  and licence are stated in the pull request.

There is no CLA to sign and no sign-off trailer to add. If you are unsure
whether something you want to include is safe to contribute, ask in the pull
request before merging rather than after: sorting out the provenance of code
that has already shipped is considerably harder than sorting it out up front.
